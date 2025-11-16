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

func TestTeamService_CreateTeam_Success(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	input := &CreateTeamInput{
		TeamName: "backend",
		Members: []TeamMemberInput{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	teamRepo.On("GetByName", mock.Anything, "backend").Return(nil, repository.ErrNotFound)
	teamRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Team")).
		Run(func(args mock.Arguments) {
			team := args.Get(1).(*domain.Team)
			team.ID = 1
		}).Return(nil)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Twice()

	team, err := service.CreateTeam(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, team)
	assert.Equal(t, "backend", team.Name)
	assert.Len(t, team.Members, 2)
}

func TestTeamService_CreateTeam_AlreadyExists(t *testing.T) {
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

	existingTeam := &domain.Team{ID: 1, Name: "backend"}
	teamRepo.On("GetByName", mock.Anything, "backend").Return(existingTeam, nil)

	team, err := service.CreateTeam(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, team)
	assert.ErrorIs(t, err, pkgErrors.ErrTeamExists)
}

func TestTeamService_CreateTeam_InvalidInput(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	tests := []struct {
		name  string
		input *CreateTeamInput
	}{
		{
			name: "empty team name",
			input: &CreateTeamInput{
				TeamName: "",
				Members:  []TeamMemberInput{{UserID: "u1", Username: "Alice"}},
			},
		},
		{
			name: "empty user id",
			input: &CreateTeamInput{
				TeamName: "backend",
				Members:  []TeamMemberInput{{UserID: "", Username: "Alice"}},
			},
		},
		{
			name: "empty username",
			input: &CreateTeamInput{
				TeamName: "backend",
				Members:  []TeamMemberInput{{UserID: "u1", Username: ""}},
			},
		},
		{
			name: "duplicate user id",
			input: &CreateTeamInput{
				TeamName: "backend",
				Members: []TeamMemberInput{
					{UserID: "u1", Username: "Alice"},
					{UserID: "u1", Username: "Bob"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := service.CreateTeam(context.Background(), tt.input)
			assert.Error(t, err)
			assert.Nil(t, team)
			assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
		})
	}
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	expectedTeam := &domain.Team{
		ID:   1,
		Name: "backend",
		Members: []*domain.User{
			{ID: "u1", Username: "Alice"},
		},
	}

	teamRepo.On("GetByName", mock.Anything, "backend").Return(expectedTeam, nil)

	team, err := service.GetTeam(context.Background(), "backend")

	assert.NoError(t, err)
	assert.NotNil(t, team)
	assert.Equal(t, "backend", team.Name)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	teamRepo.On("GetByName", mock.Anything, "nonexistent").Return(nil, repository.ErrNotFound)

	team, err := service.GetTeam(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, team)
	assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
}

func TestTeamService_GetTeam_EmptyName(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	userRepo := mocks.NewUserRepository(t)
	logger := zap.NewNop()

	service := NewTeamService(teamRepo, userRepo, logger)

	team, err := service.GetTeam(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, team)
	assert.ErrorIs(t, err, pkgErrors.ErrInvalidInput)
}
