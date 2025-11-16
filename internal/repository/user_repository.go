package repository

import (
	"context"
	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	UpdateIsActive(ctx context.Context, id string, isActive bool) error
	GetActiveUsersByTeamID(ctx context.Context, teamID int, excludedUserIDs []string) ([]*domain.User, error)
}
