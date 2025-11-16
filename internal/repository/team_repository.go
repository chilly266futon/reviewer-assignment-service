package repository

import (
	"context"
	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByName(ctx context.Context, name string) (*domain.Team, error)
	GetByID(ctx context.Context, id int) (*domain.Team, error)
}
