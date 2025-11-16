package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantError bool
	}{
		{"debug logger", "debug", false},
		{"info logger", "info", false},
		{"warn logger", "warn", false},
		{"error logger", "error", false},
		{"invalid logger", "invalid", true},
		{"empty level", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}
