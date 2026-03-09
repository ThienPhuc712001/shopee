package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Security SecurityConfig
	RateLimit RateLimitConfig
	CORS     CORSConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	BcryptCost int
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file (ignore error if not present in production)
	_ = godotenv.Load()

	config := &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "1433"),
			Name:     getEnv("DB_NAME", "ecommerce"),
			User:     getEnv("DB_USER", "sa"),
			Password: getEnv("DB_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-key"),
			Expiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),
		},
		Security: SecurityConfig{
			BcryptCost: parseInt(getEnv("BCRYPT_COST", "10")),
		},
		RateLimit: RateLimitConfig{
			Requests: parseInt(getEnv("RATE_LIMIT_REQUESTS", "100")),
			Duration: parseDuration(getEnv("RATE_LIMIT_DURATION", "1m")),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "*")),
		},
	}

	log.Printf("Configuration loaded for environment: %s", config.App.Env)
	return config
}

// getEnv retrieves environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// parseInt parses an integer from string with default fallback
func parseInt(value string) int {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return 10
	}
	return intVal
}

// parseDuration parses a duration string with default fallback
func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 24 * time.Hour
	}
	return duration
}

// parseOrigins parses comma-separated origins into a slice
func parseOrigins(value string) []string {
	if value == "*" {
		return []string{"*"}
	}
	origins := []string{}
	for _, origin := range splitString(value, ",") {
		trimmed := trimSpace(origin)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

// splitString splits a string by delimiter (simple implementation)
func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

// trimSpace trims whitespace from a string
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
