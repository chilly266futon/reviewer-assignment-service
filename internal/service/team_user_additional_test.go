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

func TestTeamService_CreateTeam_CreateError(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	input := &CreateTeamInput{
		TeamName: "backend",
		Members: []TeamMemberInput{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	teamRepo.On("GetByName", mock.Anything, "backend").Return(nil, repository.ErrNotFound)
	teamRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Team")).Return(assert.AnError)

	team, err := service.CreateTeam(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, team)
}

func TestTeamService_CreateTeam_UserCreateError(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	input := &CreateTeamInput{
		TeamName: "backend",
		Members: []TeamMemberInput{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	teamRepo.On("GetByName", mock.Anything, "backend").Return(nil, repository.ErrNotFound)
	teamRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Team")).
		Run(func(args mock.Arguments) {
			team := args.Get(1).(*domain.Team)
			team.ID = 1
		}).Return(nil)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(assert.AnError)

	team, err := service.CreateTeam(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, team)
}

func TestTeamService_CreateTeam_GetByNameError(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	input := &CreateTeamInput{
		TeamName: "backend",
		Members: []TeamMemberInput{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	teamRepo.On("GetByName", mock.Anything, "backend").Return(nil, assert.AnError)

	team, err := service.CreateTeam(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, team)
}

func TestTeamService_CreateTeam_CreateReturnsAlreadyExists(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	input := &CreateTeamInput{
		TeamName: "backend",
		Members: []TeamMemberInput{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	teamRepo.On("GetByName", mock.Anything, "backend").Return(nil, repository.ErrNotFound)
	teamRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Team")).Return(repository.ErrAlreadyExists)

	team, err := service.CreateTeam(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, team)
	assert.ErrorIs(t, err, pkgErrors.ErrTeamExists)
}

func TestTeamService_GetTeam_GetByNameError(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	teamRepo.On("GetByName", mock.Anything, "backend").Return(nil, assert.AnError)

	team, err := service.GetTeam(context.Background(), "backend")

	assert.Error(t, err)
	assert.Nil(t, team)
}

func TestUserService_SetIsActive_GetByIDError(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	userRepo.On("UpdateIsActive", mock.Anything, "u1", false).Return(nil)
	userRepo.On("GetByID", mock.Anything, "u1").Return(nil, assert.AnError)

	result, err := service.SetIsActive(context.Background(), "u1", false)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserService_SetIsActive_UpdateError(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	userRepo.On("UpdateIsActive", mock.Anything, "u1", true).Return(assert.AnError)

	result, err := service.SetIsActive(context.Background(), "u1", true)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserService_GetUserReviews_Error(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	prRepo := mocks.NewPullRequestRepository(t)
	logger := zap.NewNop()

	service := NewUserService(userRepo, prRepo, logger)

	prRepo.On("GetByReviewerID", mock.Anything, "u1").Return(nil, assert.AnError)

	prs, err := service.GetUserReviews(context.Background(), "u1")

	assert.Error(t, err)
	assert.Nil(t, prs)
}
