package service

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
	"github.com/chilly266futon/reviewer-assignment-service/mocks"
	pkgErrors "github.com/chilly266futon/reviewer-assignment-service/pkg/errors"
)

func TestPRService_CreatePR_Success(t *testing.T) {
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
		{ID: "u3", Username: "Charlie", TeamID: 1, IsActive: true},
	}

	userRepo.On("GetByID", mock.Anything, "u1").Return(author, nil)
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u1"}).Return(candidates, nil)
	prRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PullRequest"), mock.Anything).Return(nil)

	pr, err := service.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, "pr1", pr.ID)
	assert.Equal(t, "Fix bug", pr.Name)
	assert.Equal(t, "u1", pr.AuthorID)
	assert.Equal(t, domain.StatusOpen, pr.Status)
	assert.LessOrEqual(t, len(pr.AssignedReviewers), 2)
}

func TestPRService_CreatePR_AuthorNotFound(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	userRepo.On("GetByID", mock.Anything, "u1").Return(nil, repository.ErrNotFound)

	pr, err := service.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	assert.Error(t, err)
	assert.Nil(t, pr)
	assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
}

func TestPRService_CreatePR_InvalidInput(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	tests := []struct {
		name     string
		prID     string
		prName   string
		authorID string
	}{
		{"empty prID", "", "name", "author"},
		{"empty name", "pr1", "", "author"},
		{"empty authorID", "pr1", "name", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr, err := service.CreatePR(context.Background(), tt.prID, tt.prName, tt.authorID)
			assert.Error(t, err)
			assert.Nil(t, pr)
			assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
		})
	}
}

func TestPRService_CreatePR_NoReviewers(t *testing.T) {
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
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u1"}).Return([]*domain.User{}, nil)
	prRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.PullRequest"), mock.Anything).Return(nil)

	pr, err := service.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Empty(t, pr.AssignedReviewers)
}

func TestPRService_MergePR_Success(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:     "pr1",
		Name:   "Fix bug",
		Status: domain.StatusOpen,
	}

	mergedPR := &domain.PullRequest{
		ID:     "pr1",
		Name:   "Fix bug",
		Status: domain.StatusMerged,
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil).Once()
	prRepo.On("UpdateStatus", mock.Anything, "pr1", domain.StatusMerged, mock.Anything).Return(nil)
	prRepo.On("GetByID", mock.Anything, "pr1").Return(mergedPR, nil).Once()

	result, err := service.MergePR(context.Background(), "pr1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.StatusMerged, result.Status)
}

func TestPRService_MergePR_AlreadyMerged(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:     "pr1",
		Name:   "Fix bug",
		Status: domain.StatusMerged,
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil)

	result, err := service.MergePR(context.Background(), "pr1")

	assert.NoError(t, err) // Должна быть идемпотентность
	assert.NotNil(t, result)
	assert.Equal(t, domain.StatusMerged, result.Status)
}

func TestPRService_MergePR_NotFound(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	prRepo.On("GetByID", mock.Anything, "pr1").Return(nil, repository.ErrNotFound)

	result, err := service.MergePR(context.Background(), "pr1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
}

func TestPRService_ReassignReviewer_Success(t *testing.T) {
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
		{ID: "u4", Username: "David", TeamID: 1, IsActive: true},
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil).Once()
	userRepo.On("GetByID", mock.Anything, "u2").Return(oldReviewer, nil)
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u2", "u1"}).Return(candidates, nil)
	prRepo.On("ReplaceReviewer", mock.Anything, "pr1", "u2", mock.Anything).Return(nil)
	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil).Once()

	newReviewerID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.NoError(t, err)
	assert.NotNil(t, updatedPR)
	assert.NotEmpty(t, newReviewerID)
}

func TestPRService_ReassignReviewer_PRMerged(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:                "pr1",
		Name:              "Fix bug",
		Status:            domain.StatusMerged,
		AssignedReviewers: []string{"u2"},
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil)

	newReviewerID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Nil(t, updatedPR)
	assert.Empty(t, newReviewerID)
	assert.ErrorIs(t, err, pkgErrors.ErrPRMerged)
}

func TestPRService_ReassignReviewer_ReviewerNotAssigned(t *testing.T) {
	prRepo := mocks.NewPullRequestRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewPRService(prRepo, userRepo, logger)

	pr := &domain.PullRequest{
		ID:                "pr1",
		Name:              "Fix bug",
		Status:            domain.StatusOpen,
		AssignedReviewers: []string{"u3"},
	}

	prRepo.On("GetByID", mock.Anything, "pr1").Return(pr, nil)

	newReviewerID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Nil(t, updatedPR)
	assert.Empty(t, newReviewerID)
	assert.ErrorIs(t, err, pkgErrors.ErrNotAssigned)
}

func TestPRService_ReassignReviewer_NoCandidates(t *testing.T) {
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
	userRepo.On("GetActiveUsersByTeamID", mock.Anything, 1, []string{"u2", "u1"}).Return([]*domain.User{}, nil)

	newReviewerID, updatedPR, err := service.ReassignReviewer(context.Background(), "pr1", "u2")

	assert.Error(t, err)
	assert.Nil(t, updatedPR)
	assert.Empty(t, newReviewerID)
	assert.ErrorIs(t, err, pkgErrors.ErrNoCandidate)
}

// Тесты для вспомогательных функций
func TestSelectRandomReviewers_ReturnsUpToMaxCount(t *testing.T) {
	service := &PRService{
		rng: rand.New(rand.NewSource(42)),
	}

	tests := []struct {
		name           string
		candidates     []*domain.User
		maxCount       int
		expectedLength int
	}{
		{
			name: "returns maxCount reviewers when enough candidates",
			candidates: []*domain.User{
				{ID: "u1"}, {ID: "u2"}, {ID: "u3"}, {ID: "u4"},
			},
			maxCount:       2,
			expectedLength: 2,
		},
		{
			name: "returns all candidates when maxCount exceeds available",
			candidates: []*domain.User{
				{ID: "u1"},
			},
			maxCount:       2,
			expectedLength: 1,
		},
		{
			name:           "returns empty slice for empty candidates",
			candidates:     []*domain.User{},
			maxCount:       2,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.selectRandomReviewers(tt.candidates, tt.maxCount)
			assert.Len(t, result, tt.expectedLength)
		})
	}
}

func TestSelectRandomReviewer_SingleSelection(t *testing.T) {
	service := &PRService{
		rng: rand.New(rand.NewSource(42)),
	}

	tests := []struct {
		name       string
		candidates []*domain.User
		shouldFind bool
	}{
		{
			name: "selects one reviewer",
			candidates: []*domain.User{
				{ID: "u1"}, {ID: "u2"}, {ID: "u3"},
			},
			shouldFind: true,
		},
		{
			name:       "returns empty string for empty candidates",
			candidates: []*domain.User{},
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.selectRandomReviewer(tt.candidates)
			if tt.shouldFind {
				assert.NotEmpty(t, result)
			} else {
				assert.Empty(t, result)
			}
		})
	}
}

func TestSelectRandomReviewers_Distribution(t *testing.T) {
	service := &PRService{
		rng: rand.New(rand.NewSource(42)),
	}

	candidates := []*domain.User{
		{ID: "u1"}, {ID: "u2"}, {ID: "u3"}, {ID: "u4"}, {ID: "u5"},
	}

	// Проверяем, что при многократном вызове выбираются разные комбинации
	selections := make(map[string]int)
	for i := 0; i < 100; i++ {
		reviewers := service.selectRandomReviewers(candidates, 2)
		key := reviewers[0] + "," + reviewers[1]
		selections[key]++
	}

	// Должно быть несколько разных комбинаций
	assert.Greater(t, len(selections), 1)
}
