package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
)

type PullRequestRepository struct {
	pool      *pgxpool.Pool
	txManager *TxManager
	logger    *zap.Logger
}

func NewPRRepository(pool *pgxpool.Pool, txManager *TxManager, logger *zap.Logger) *PullRequestRepository {
	return &PullRequestRepository{
		pool:      pool,
		txManager: txManager,
		logger:    logger,
	}
}

// Create создает PR с назначенными ревьюерами (атомарно)
func (r PullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest, reviewerIDs []string) error {
	return r.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		// Создаем PR
		prQuery := `
			INSERT INTO pull_requests (id, name, author_id, status_id, created_at)
			VALUES ($1, $2, $3, (SELECT id FROM pr_statuses WHERE name = $4), $5)
		`

		_, err := tx.Exec(ctx, prQuery,
			pr.ID,
			pr.Name,
			pr.AuthorID,
			pr.Status,
			pr.CreatedAt,
		)

		if err != nil {
			if isUniqueViolation(err) {

				return repository.ErrAlreadyExists
			}

			r.logger.Error("failed to insert pull request",
				zap.String("pr_id", pr.ID),
				zap.Error(err),
			)
			return fmt.Errorf("insert pull request: %w", err)
		}

		// Назначаем ревьюеров
		if len(reviewerIDs) > 0 {
			reviewerQuery := `
			INSERT INTO pr_reviewers (pull_request_id, user_id, assigned_at)
			VALUES ($1, $2, now()) 
		`

			for _, reviewerID := range reviewerIDs {
				_, err := tx.Exec(ctx, reviewerQuery, pr.ID, reviewerID)
				if err != nil {
					r.logger.Error("failed to assign reviewer",
						zap.String("pr_id", pr.ID),
						zap.String("reviewer_id", reviewerID),
						zap.Error(err),
					)
					return fmt.Errorf("assign reviewer: %w", err)
				}
			}
		}

		return nil
	})
}

// GetByID возвращает по ID с назначенными ревьюерами
func (r PullRequestRepository) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	query := `
		SELECT 
		    pr.id,
		    pr.name,
		    pr.author_id,
		    ps.name as status,
		    pr.created_at,
		    pr.merged_at,
		    COALESCE(array_agg(rev.user_id) FILTER ( WHERE rev.user_id IS NOT NULL ), '{}') as reviewers
		FROM pull_requests pr
		INNER JOIN pr_statuses ps ON pr.status_id = ps.id
		LEFT JOIN pr_reviewers rev ON rev.pull_request_id = pr.id
		WHERE pr.id = $1
		GROUP BY pr.id, ps.name, pr.id, pr.name, pr.author_id, ps.name, pr.created_at, pr.merged_at
	`

	var pr domain.PullRequest
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
		&pr.AssignedReviewers,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		r.logger.Error("failed to get pull request",
			zap.String("pr_id", pr.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get pull request: %w", err)
	}

	return &pr, nil
}

// UpdateStatus обновляет статус PR, если он в статусе OPEN
func (r PullRequestRepository) UpdateStatus(ctx context.Context, id string, status string, mergedAt *time.Time) error {
	query := `
		UPDATE pull_requests
		SET 
		    status_id = (SELECT id FROM pr_statuses WHERE name = $2),
		    merged_at = $3
		WHERE id = $1 AND status_id = (SELECT id FROM pr_statuses WHERE name = 'OPEN')
	`

	result, err := r.pool.Exec(ctx, query, id, status, mergedAt)
	if err != nil {
		r.logger.Error("failed to update PR status",
			zap.String("pr_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("update PR status: %w", err)
	}

	if result.RowsAffected() == 0 {
		// Проверяем существование PR
		existsQuery := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)`
		var exists bool
		err := r.pool.QueryRow(ctx, existsQuery, id).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check PR exists: %w", err)
		}

		if !exists {
			return repository.ErrNotFound
		}

		// PR существует, но уже не в статусе OPEN (идемпотентность - ничего не делаем)
		return nil
	}

	return nil
}

// ReplaceReviewer заменяет одного ревьюера на другого
func (r PullRequestRepository) ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error {
	return r.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		// Проверяем, что PR в статусе OPEN
		statusQuery := `
			SELECT ps.name
			FROM pull_requests pr
			INNER JOIN pr_statuses ps ON pr.status_id = ps.id
			WHERE pr.id = $1
		`

		var status string
		err := tx.QueryRow(ctx, statusQuery, prID).Scan(&status)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return repository.ErrNotFound
			}
			return fmt.Errorf("check PR status: %w", err)
		}

		if status != domain.StatusOpen {
			return fmt.Errorf("PR is %s: %w", status, repository.ErrConflict)
		}

		// Удаляем старого ревьюера
		deleteQuery := `
			DELETE FROM pr_reviewers
			WHERE pull_request_id = $1 AND user_id = $2
		`

		result, err := tx.Exec(ctx, deleteQuery, prID, oldUserID)
		if err != nil {
			return fmt.Errorf("delete old reviewer: %w", err)
		}

		if result.RowsAffected() == 0 {
			// Старый ревьюер не был назначен
			return fmt.Errorf("reviewer not assigned: %w", repository.ErrNotFound)
		}

		// Добавляем нового ревьюера
		insertQuery := `
			INSERT INTO pr_reviewers(pull_request_id, user_id, assigned_at)
			VALUES ($1, $2, now())
		`

		_, err = tx.Exec(ctx, insertQuery, prID, newUserID)
		if err != nil {
			r.logger.Error("failed to insert new reviewer",
				zap.String("pr_id", prID),
				zap.String("new_user_id", newUserID),
				zap.Error(err),
			)
			return fmt.Errorf("insert new reviewer: %w", err)
		}

		return nil
	})
}

// GetByReviewerID возвращает все PR, где пользователь назначен ревьюером
func (r PullRequestRepository) GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error) {
	query := `
		SELECT 
		    pr.id,
		    pr.name,
		    pr.author_id,
		    ps.name as status
		FROM pull_requests pr
		INNER JOIN pr_statuses ps ON pr.status_id = ps.id
		INNER JOIN pr_reviewers rev ON rev.pull_request_id = pr.id
		WHERE rev.user_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, reviewerID)
	if err != nil {
		r.logger.Error("failed to get PRs by reviewer",
			zap.String("reviewer_id", reviewerID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get PRs by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
		)
		if err != nil {
			r.logger.Error("failed to scan PR row", zap.Error(err))
			return nil, fmt.Errorf("scan PR: %w", err)
		}
		prs = append(prs, &pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate PRs: %w", err)
	}

	return prs, nil
}
