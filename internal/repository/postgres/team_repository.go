package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
)

type TeamRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewTeamRepository(pool *pgxpool.Pool, logger *zap.Logger) *TeamRepository {
	return &TeamRepository{
		pool:   pool,
		logger: logger,
	}
}

// Create создаёт новую команду
func (r *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	query := `
		INSERT INTO teams (name, created_at)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query, team.Name, time.Now()).Scan(&team.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrAlreadyExists
		}
		r.logger.Error("failed to create team",
			zap.String("team_name", team.Name),
			zap.Error(err),
		)
		return fmt.Errorf("create team: %w", err)
	}

	return nil
}

// GetByName возвращает команду по имени с участниками
func (r *TeamRepository) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	query := `
		SELECT 
		    t.id, t.name, t.created_at,
		    u.id, u.username, u.team_id, u.is_active, u.created_at, u.updated_at
		FROM teams t
		LEFT JOIN users u ON u.team_id = t.id
		WHERE t.name = $1
		ORDER BY u.id
 	`

	rows, err := r.pool.Query(ctx, query, name)
	if err != nil {
		r.logger.Error("failed to get team with members",
			zap.String("team_name", name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get team: %w", err)
	}
	defer rows.Close()

	var (
		team    *domain.Team
		members []*domain.User

		teamID        int
		teamName      string
		teamCreatedAt time.Time

		userID        *string
		username      *string
		userTeamID    *int
		isActive      *bool
		userCreatedAt *time.Time
		userUpdatedAt *time.Time
	)

	for rows.Next() {
		err := rows.Scan(
			&teamID, &teamName, &teamCreatedAt,
			&userID, &username, &userTeamID, &isActive, &userCreatedAt, &userUpdatedAt,
		)
		if err != nil {
			r.logger.Error("failed to scan team row", zap.Error(err))
			return nil, fmt.Errorf("scan team: %w", err)
		}

		if team == nil {
			team = &domain.Team{
				ID:        teamID,
				Name:      teamName,
				CreatedAt: teamCreatedAt,
			}
		}
		if userID != nil {
			members = append(members, &domain.User{
				ID:        *userID,
				Username:  *username,
				TeamID:    *userTeamID,
				TeamName:  teamName,
				IsActive:  *isActive,
				CreatedAt: *userCreatedAt,
				UpdatedAt: *userUpdatedAt,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	if team == nil {
		return nil, repository.ErrNotFound
	}

	team.Members = members
	return team, nil
}

// GetByID возвращает команду по ID
func (r *TeamRepository) GetByID(ctx context.Context, id int) (*domain.Team, error) {
	query := `
		SELECT id, name, created_at
		FROM teams
		WHERE id = $1
	`

	var team domain.Team
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&team.ID,
		&team.Name,
		&team.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		r.logger.Error("failed to get team by ID",
			zap.Int("team_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get team: %w", err)
	}

	return &team, nil
}

func isUniqueViolation(err error) bool {
	// pgx возвращает ошибку с кодом 23505 для нарушения уникальности
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
