package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// Server
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	// Database
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`

	//Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
