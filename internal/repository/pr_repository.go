package repository

import (
	"context"
	"time"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest, reviewerIDs []string) error
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
	UpdateStatus(ctx context.Context, id string, status string, mergedAt *time.Time) error
	ReplaceReviewer(ctx context.Context, prID, oldUserID, newUserID string) error
	GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error)
}
