package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePRInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   CreatePRInput
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   CreatePRInput{PullRequestID: "pr1", PullRequestName: "Fix", AuthorID: "u1"},
			wantErr: false,
		},
		{
			name:    "empty pr id",
			input:   CreatePRInput{PullRequestID: "", PullRequestName: "Fix", AuthorID: "u1"},
			wantErr: true,
		},
		{
			name:    "empty pr name",
			input:   CreatePRInput{PullRequestID: "pr1", PullRequestName: "", AuthorID: "u1"},
			wantErr: true,
		},
		{
			name:    "empty author id",
			input:   CreatePRInput{PullRequestID: "pr1", PullRequestName: "Fix", AuthorID: ""},
			wantErr: true,
		},
		{
			name:    "pr id too long",
			input:   CreatePRInput{PullRequestID: string(make([]byte, 101)), PullRequestName: "Fix", AuthorID: "u1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetIsActiveInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   SetIsActiveInput
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   SetIsActiveInput{UserID: "u1", IsActive: true},
			wantErr: false,
		},
		{
			name:    "empty user id",
			input:   SetIsActiveInput{UserID: "", IsActive: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReassignReviewerInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   ReassignReviewerInput
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   ReassignReviewerInput{PullRequestID: "pr1", OldReviewerID: "u1"},
			wantErr: false,
		},
		{
			name:    "empty pr id",
			input:   ReassignReviewerInput{PullRequestID: "", OldReviewerID: "u1"},
			wantErr: true,
		},
		{
			name:    "empty reviewer id",
			input:   ReassignReviewerInput{PullRequestID: "pr1", OldReviewerID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMergePRInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   MergePRInput
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   MergePRInput{PullRequestID: "pr1"},
			wantErr: false,
		},
		{
			name:    "empty pr id",
			input:   MergePRInput{PullRequestID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
