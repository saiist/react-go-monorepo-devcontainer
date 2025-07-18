package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string
	JWTSecret   string
	Environment string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:        getEnv("API_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/myapp_dev?sslmode=disable"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: getEnv("GO_ENV", "development"),
	}
}

// IsProduction returns true if the application is running in production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
