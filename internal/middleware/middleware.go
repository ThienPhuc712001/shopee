// Package middleware provides initialization for all middleware
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"ecommerce/pkg/logger"
)

// Config holds middleware configuration
type Config struct {
	// EnableRecovery enables panic recovery middleware
	EnableRecovery bool

	// EnableRequestLog enables request logging middleware
	EnableRequestLog bool

	// EnableMetrics enables metrics collection middleware
	EnableMetrics bool

	// SlowRequestThreshold is the threshold for slow request logging
	SlowRequestThreshold time.Duration

	// EnableSlowRequestLog enables slow request logging
	EnableSlowRequestLog bool
}

// DefaultConfig returns default middleware configuration
func DefaultConfig() *Config {
	return &Config{
		EnableRecovery:       true,
		EnableRequestLog:     true,
		EnableMetrics:        false,
		SlowRequestThreshold: 5 * time.Second,
		EnableSlowRequestLog: true,
	}
}

// Setup initializes and returns all configured middleware
func Setup(log *logger.Logger, config *Config) []gin.HandlerFunc {
	if config == nil {
		config = DefaultConfig()
	}

	var middlewares []gin.HandlerFunc

	// Recovery middleware (should be first)
	if config.EnableRecovery {
		middlewares = append(middlewares, RecoveryMiddleware(log))
	}

	// Request logging middleware
	if config.EnableRequestLog {
		middlewares = append(middlewares, RequestLogMiddleware(log))
	}

	// Metrics middleware
	if config.EnableMetrics {
		middlewares = append(middlewares, MetricsMiddleware())
	}

	// Slow request logging middleware
	if config.EnableSlowRequestLog {
		middlewares = append(middlewares, SlowRequestMiddleware(log, config.SlowRequestThreshold))
	}

	// Error handling middleware (should be after handlers)
	middlewares = append(middlewares, ErrorMiddleware(log))

	return middlewares
}

// SetupMinimal returns minimal middleware for production (recovery + error handling)
func SetupMinimal(log *logger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		RecoveryMiddleware(log),
		ErrorMiddleware(log),
	}
}

// SetupWithAll returns all middleware with default configuration
func SetupWithAll(log *logger.Logger) []gin.HandlerFunc {
	return Setup(log, &Config{
		EnableRecovery:       true,
		EnableRequestLog:     true,
		EnableMetrics:        true,
		SlowRequestThreshold: 5 * time.Second,
		EnableSlowRequestLog: true,
	})
}
