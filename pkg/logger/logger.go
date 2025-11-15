package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func New(level string) (*zap.Logger, error) {
	var cfg zap.Config

	switch level {
	case "debug":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	case "warn":
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)

	case "error":
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("invalid log level: %s", level)
	}

	cfg.DisableStacktrace = true

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger, nil
}
