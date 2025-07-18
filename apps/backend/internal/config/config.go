package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string
	JWTSecret   string
	Environment string
	
	// Logging configuration
	LogLevel         string
	LogFormat        string
	LogSamplingRate  float64
}

// Load loads configuration from environment variables
func Load() *Config {
	env := getEnv("GO_ENV", "development")
	
	// Default log configuration based on environment
	defaultLogLevel := "debug"
	defaultLogFormat := "text"
	if env == "production" {
		defaultLogLevel = "info"
		defaultLogFormat = "json"
	}
	
	return &Config{
		Port:        getEnv("API_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/myapp_dev?sslmode=disable"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: env,
		
		// Logging configuration
		LogLevel:        getEnv("LOG_LEVEL", defaultLogLevel),
		LogFormat:       getEnv("LOG_FORMAT", defaultLogFormat),
		LogSamplingRate: getEnvFloat("LOG_SAMPLING_RATE", 1.0),
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

// getEnvFloat gets an environment variable as float64 or returns a default value
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}
