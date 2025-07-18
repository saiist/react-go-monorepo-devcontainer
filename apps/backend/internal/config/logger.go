package config

import (
	"github.com/react-go-monorepo/backend/internal/logger"
)

// NewLogger creates a logger instance based on configuration
func NewLogger(cfg *Config) *logger.Logger {
	return logger.New(logger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
}