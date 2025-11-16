package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
	"github.com/chilly266futon/reviewer-assignment-service/mocks"
	pkgErrors "github.com/chilly266futon/reviewer-assignment-service/pkg/errors"
)

func TestUserService_SetIsActive_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	user := &domain.User{
		ID:       "u1",
		Username: "Alice",
		IsActive: false,
	}

	userRepo.On("UpdateIsActive", mock.Anything, "u1", false).Return(nil)
	userRepo.On("GetByID", mock.Anything, "u1").Return(user, nil)

	result, err := service.SetIsActive(context.Background(), "u1", false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "u1", result.ID)
	assert.False(t, result.IsActive)
}

func TestUserService_SetIsActive_UserNotFound(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	userRepo.On("UpdateIsActive", mock.Anything, "u999", true).Return(repository.ErrNotFound)

	result, err := service.SetIsActive(context.Background(), "u999", true)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
}

func TestUserService_SetIsActive_EmptyUserID(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	result, err := service.SetIsActive(context.Background(), "", true)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
}

func TestUserService_GetUserReviews_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	expectedPRs := []*domain.PullRequest{
		{ID: "pr1", Name: "Fix bug 1", Status: domain.StatusOpen},
		{ID: "pr2", Name: "Fix bug 2", Status: domain.StatusMerged},
	}

	prRepo.On("GetByReviewerID", mock.Anything, "u1").Return(expectedPRs, nil)

	prs, err := service.GetUserReviews(context.Background(), "u1")

	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Len(t, prs, 2)
}

func TestUserService_GetUserReviews_EmptyUserID(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	prs, err := service.GetUserReviews(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, prs)
	assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
}

func TestUserService_GetUserReviews_NoPRs(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	prRepo.On("GetByReviewerID", mock.Anything, "u1").Return([]*domain.PullRequest{}, nil)

	prs, err := service.GetUserReviews(context.Background(), "u1")

	assert.NoError(t, err)
	assert.NotNil(t, prs)
	assert.Empty(t, prs)
}
