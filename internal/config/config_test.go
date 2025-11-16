package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chilly266futon/reviewer-assignment-service/internal/config"
)

func TestLoad_Success(t *testing.T) {
	// Устанавливаем необходимые переменные окружения
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "testpass", cfg.DBPassword)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "customuser")
	os.Setenv("DB_PASSWORD", "custompass")
	os.Setenv("DB_NAME", "customdb")
	os.Setenv("DB_SSL_MODE", "require")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSL_MODE")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "custom-host", cfg.DBHost)
	assert.Equal(t, "3306", cfg.DBPort)
	assert.Equal(t, "customuser", cfg.DBUser)
	assert.Equal(t, "custompass", cfg.DBPassword)
	assert.Equal(t, "customdb", cfg.DBName)
	assert.Equal(t, "require", cfg.DBSSLMode)
	assert.Equal(t, "9090", cfg.ServerPort)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestLoad_MissingRequired(t *testing.T) {
	// Очищаем все переменные окружения
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	_, err := config.Load()
	assert.Error(t, err)
}
