package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
	pkgErrors "github.com/chilly266futon/reviewer-assignment-service/pkg/errors"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
	logger   *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PullRequestRepository, logger *zap.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
		logger:   logger,
	}
}

// SetIsActive изменяет статус активности пользователя и возвращает обновленного пользователя
func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if userID == "" {
		return nil, pkgErrors.ErrInvalidInput
	}

	// Обновляем статус пользователя
	if err := s.userRepo.UpdateIsActive(ctx, userID, isActive); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("failed to update user active status",
			zap.String("user_id", userID),
			zap.Bool("isActive", isActive),
			zap.Error(err),
		)
		return nil, fmt.Errorf("update user active status: %w", err)
	}

	s.logger.Info("user active status updated",
		zap.String("user_id", userID),
		zap.Bool("isActive", isActive),
	)

	// Получаем обновленного пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get updated user: %w", err)
	}

	return user, nil
}

// GetUserReviews возвращает все PRs, в которых пользователь назначен ревьюером
func (s *UserService) GetUserReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	if userID == "" {
		return nil, pkgErrors.ErrInvalidInput
	}

	prs, err := s.prRepo.GetByReviewerID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user reviews",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get user reviews: %w", err)
	}

	return prs, nil
}
