package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
)

func TestReassignReviewerRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     ReassignReviewerRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     ReassignReviewerRequest{PullRequestID: "pr1", OldUserID: "u1"},
			wantErr: false,
		},
		{
			name:    "empty pr id",
			req:     ReassignReviewerRequest{PullRequestID: "", OldUserID: "u1"},
			wantErr: true,
		},
		{
			name:    "empty reviewer id",
			req:     ReassignReviewerRequest{PullRequestID: "pr1", OldUserID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToPRShort(t *testing.T) {
	pr := &domain.PullRequest{
		ID:       "pr1",
		Name:     "Fix bug",
		AuthorID: "u1",
		Status:   domain.StatusOpen,
	}

	short := ToPRShort(pr)

	assert.NotNil(t, short)
	assert.Equal(t, "pr1", short.PullRequestID)
	assert.Equal(t, "Fix bug", short.PullRequestName)
	assert.Equal(t, "u1", short.AuthorID)
	assert.Equal(t, domain.StatusOpen, short.Status)
}

func TestCreateTeamRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTeamRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateTeamRequest{
				TeamName: "backend",
				Members:  []TeamMemberRequest{{UserID: "u1", Username: "Alice", IsActive: true}},
			},
			wantErr: false,
		},
		{
			name: "empty team name",
			req: CreateTeamRequest{
				TeamName: "",
				Members:  []TeamMemberRequest{{UserID: "u1", Username: "Alice", IsActive: true}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetIsActiveRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     SetIsActiveRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     SetIsActiveRequest{UserID: "u1", IsActive: true},
			wantErr: false,
		},
		{
			name:    "empty user id",
			req:     SetIsActiveRequest{UserID: "", IsActive: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreatePRRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreatePRRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     CreatePRRequest{PullRequestID: "pr1", PullRequestName: "Fix", AuthorID: "u1"},
			wantErr: false,
		},
		{
			name:    "empty pr id",
			req:     CreatePRRequest{PullRequestID: "", PullRequestName: "Fix", AuthorID: "u1"},
			wantErr: true,
		},
		{
			name:    "empty pr name",
			req:     CreatePRRequest{PullRequestID: "pr1", PullRequestName: "", AuthorID: "u1"},
			wantErr: true,
		},
		{
			name:    "empty author id",
			req:     CreatePRRequest{PullRequestID: "pr1", PullRequestName: "Fix", AuthorID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMergePRRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     MergePRRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     MergePRRequest{PullRequestID: "pr1"},
			wantErr: false,
		},
		{
			name:    "empty pr id",
			req:     MergePRRequest{PullRequestID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
