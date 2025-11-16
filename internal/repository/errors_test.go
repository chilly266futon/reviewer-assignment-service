package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", repository.ErrNotFound},
		{"ErrAlreadyExists", repository.ErrAlreadyExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.Error(t, tt.err)
		})
	}
}
