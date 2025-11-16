package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
)

type UserRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewUserRepository(pool *pgxpool.Pool, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		pool:   pool,
		logger: logger,
	}
}

// Create создает или обновляет пользователя
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, team_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET 
		   username = EXCLUDED.username,
		   team_id = EXCLUDED.team_id,
		   is_active = EXCLUDED.is_active,
		   updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Username,
		user.TeamID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("failed to create user",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.team_id, t.name as team_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN teams t ON u.team_id = t.id
		WHERE u.id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.TeamID,
		&user.TeamName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		r.logger.Error("failed to get user",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil
}

// UpdateIsActive изменяет статус активности пользователя
func (r *UserRepository) UpdateIsActive(ctx context.Context, id string, isActive bool) error {
	query := `
		UPDATE users
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id, isActive)
	if err != nil {
		r.logger.Error("failed to update user activity status",
			zap.String("user_id", id),
			zap.Bool("is_active", isActive),
			zap.Error(err),
		)
		return fmt.Errorf("update user active status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// GetActiveByTeamID возвращает активных пользователей команды, исключая указанных
func (r *UserRepository) GetActiveUsersByTeamID(ctx context.Context, teamID int, excludeUserIDs []string) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.team_id, t.name as team_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN teams t ON u.team_id = t.id
		WHERE u.team_id = $1 
		  AND u.is_active = TRUE
		  AND u.id != ALL($2)
		ORDER BY u.id
	`

	rows, err := r.pool.Query(ctx, query, teamID, excludeUserIDs)
	if err != nil {
		r.logger.Error("failed to get active users by team",
			zap.Int("team_id", teamID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get active users by team: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.TeamID,
			&user.TeamName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan user row", zap.Error(err))
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}
