package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullRequest_IsOpen(t *testing.T) {
	tests := []struct {
		name     string
		pr       PullRequest
		expected bool
	}{
		{
			name:     "open pr",
			pr:       PullRequest{ID: "pr1", Status: StatusOpen},
			expected: true,
		},
		{
			name:     "merged pr",
			pr:       PullRequest{ID: "pr2", Status: StatusMerged},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pr.IsOpen()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPullRequest_IsMerged(t *testing.T) {
	tests := []struct {
		name     string
		pr       PullRequest
		expected bool
	}{
		{
			name:     "merged pr",
			pr:       PullRequest{ID: "pr1", Status: StatusMerged},
			expected: true,
		},
		{
			name:     "open pr",
			pr:       PullRequest{ID: "pr2", Status: StatusOpen},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pr.IsMerged()
			assert.Equal(t, tt.expected, result)
		})
	}
}
