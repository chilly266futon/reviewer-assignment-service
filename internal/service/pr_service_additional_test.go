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

func TestPRService_CreatePR_PRAlreadyExists(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	author := &domain.User{
		ID:       "u1",
		Username: "Alice",
		TeamID:   1,
		IsActive: true,
	}

	candidates := []*domain.User{
		{ID: "u2", Username: "Bob", TeamID: 1, IsActive: true},
	}

	userRepo.On("GetByID", mock.Anything, "u1").Return(author, nil)
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u1"}).Return(candidates, nil)
	prRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PullRequest"), mock.Anything).Return(repository.ErrAlreadyExists)

	pr, err := service.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.ErrorIs(t, err, pkgErrors.ErrPRExists)
}

func TestPRService_CreatePR_GetCandidatesError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	author := &domain.User{
		ID:       "u1",
		Username: "Alice",
		TeamID:   1,
		IsActive: true,
	}

	userRepo.On("GetByID", mock.Anything, "u1").Return(author, nil)
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u1"}).Return(nil, assert.AnError)

	pr, err := service.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	assert.Error(t, err)
	assert.Nil(t, pr)
}

func TestPRService_CreatePR_GetAuthorError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	userRepo.On("GetByID", mock.Anything, "u1").Return(nil, assert.AnError)

	pr, err := service.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	assert.Error(t, err)
	assert.Nil(t, pr)
}

func TestPRService_MergePR_EmptyID(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	result, err := service.MergePR(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
}

func TestPRService_MergePR_GetPRError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	prRepo.On("GetByID", mock.Anything, "pr1").Return(nil, assert.AnError)

	result, err := service.MergePR(context.Background(), "pr1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPRService_MergePR_UpdateStatusError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:     "pr1",
		Name:   "Fix bug",
		Status: domain.StatusOpen,
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil).Once()
	prRepo.On("UpdateStatus", mock.Anything, "pr1", domain.StatusMerged, mock.Anything).Return(assert.AnError)

	result, err := service.MergePR(context.Background(), "pr1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPRService_ReassignReviewer_EmptyPRID(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	newID, pr, err := service.ReassignReviewer(context.Background(), "", "u2")

	assert.Error(t, err)
	assert.Empty(t, newID)
	assert.Nil(t, pr)
	assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
}

func TestPRService_ReassignReviewer_EmptyReviewerID(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	newID, pr, err := service.ReassignReviewer(context.Background(), "pr1", "")

	assert.Error(t, err)
	assert.Empty(t, newID)
	assert.Nil(t, pr)
	assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
}

func TestPRService_ReassignReviewer_GetPRError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	prRepo.On("GetByID", mock.Anything, "pr1").Return(nil, assert.AnError)

	newID, pr, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Empty(t, newID)
	assert.Nil(t, pr)
}

func TestPRService_ReassignReviewer_GetOldReviewerError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:                "pr1",
		Name:              "Fix bug",
		Status:            domain.StatusOpen,
		AuthorID:          "u1",
		AssignedReviewers: []string{"u2"},
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil)
	userRepo.On("GetByID", mock.Anything, "u2").Return(nil, assert.AnError)

	newID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Empty(t, newID)
	assert.Nil(t, updatedPR)
}

func TestPRService_ReassignReviewer_GetCandidatesError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:                "pr1",
		Name:              "Fix bug",
		Status:            domain.StatusOpen,
		AuthorID:          "u1",
		AssignedReviewers: []string{"u2"},
	}

	oldReviewer := &domain.User{
		ID:       "u2",
		Username: "Bob",
		TeamID:   1,
		IsActive: true,
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil)
	userRepo.On("GetByID", mock.Anything, "u2").Return(oldReviewer, nil)
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u2", "u1"}).Return(nil, assert.AnError)

	newID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Empty(t, newID)
	assert.Nil(t, updatedPR)
}

func TestPRService_ReassignReviewer_ReplaceError(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:                "pr1",
		Name:              "Fix bug",
		Status:            domain.StatusOpen,
		AuthorID:          "u1",
		AssignedReviewers: []string{"u2"},
	}

	oldReviewer := &domain.User{
		ID:       "u2",
		Username: "Bob",
		TeamID:   1,
		IsActive: true,
	}

	candidates := []*domain.User{
		{ID: "u3", Username: "Charlie", TeamID: 1, IsActive: true},
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil).Once()
	userRepo.On("GetByID", mock.Anything, "u2").Return(oldReviewer, nil)
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u2", "u1"}).Return(candidates, nil)
	prRepo.On("ReplaceReviewer", mock.Anything, "pr1", "u2", mock.Anything).Return(assert.AnError)

	newID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Empty(t, newID)
	assert.Nil(t, updatedPR)
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{"item exists", []string{"a", "b", "c"}, "b", true},
		{"item not exists", []string{"a", "b", "c"}, "d", false},
		{"empty slice", []string{}, "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}
