package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorConstants(t *testing.T) {
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrTeamExists)
	assert.NotNil(t, ErrPRExists)
	assert.NotNil(t, ErrPRMerged)
	assert.NotNil(t, ErrNotAssigned)
	assert.NotNil(t, ErrNoCandidate)
	assert.NotNil(t, ErrInvalidInput)
}

func TestErrorCodes(t *testing.T) {
	assert.Equal(t, "NOT_FOUND", CodeNotFound)
	assert.Equal(t, "TEAM_EXISTS", CodeTeamExists)
	assert.Equal(t, "PR_EXISTS", CodePRExists)
	assert.Equal(t, "PR_MERGED", CodePRMerged)
	assert.Equal(t, "NOT_ASSIGNED", CodeNotAssigned)
	assert.Equal(t, "NO_CANDIDATE", CodeNoCandidate)
}

func TestMapErrorToCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"not found", ErrNotFound, CodeNotFound},
		{"team exists", ErrTeamExists, CodeTeamExists},
		{"pr exists", ErrPRExists, CodePRExists},
		{"pr merged", ErrPRMerged, CodePRMerged},
		{"not assigned", ErrNotAssigned, CodeNotAssigned},
		{"no candidate", ErrNoCandidate, CodeNoCandidate},
		{"unknown error", assert.AnError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := MapErrorToCode(tt.err)
			assert.Equal(t, tt.expected, code)
		})
	}
}
