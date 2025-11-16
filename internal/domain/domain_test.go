package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
)

func TestTeam_Methods(t *testing.T) {
	team := &domain.Team{
		ID:   1,
		Name: "Test Team",
	}

	assert.Equal(t, 1, team.ID)
	assert.Equal(t, "Test Team", team.Name)
	assert.NotNil(t, team)
}

func TestUser_Methods(t *testing.T) {
	user := &domain.User{
		ID:       "user1",
		Username: "testuser",
		TeamID:   1,
		TeamName: "Team 1",
		IsActive: true,
	}

	assert.Equal(t, "user1", user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 1, user.TeamID)
	assert.Equal(t, "Team 1", user.TeamName)
	assert.True(t, user.IsActive)
}
