package middleware

import (
	"ecommerce/internal/errors"
	"ecommerce/pkg/logger"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// RATE LIMITER
// ============================================================================

// RateLimiter implements a sliding window rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	// Limit is the maximum number of requests allowed
	Limit int

	// Window is the time window for rate limiting
	Window time.Duration

	// Log is the logger for rate limit events
	Log *logger.Logger
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.Limit <= 0 {
		config.Limit = 100
	}
	if config.Window <= 0 {
		config.Window = time.Minute
	}

	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    config.Limit,
		window:   config.Window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request is allowed for the given key (IP)
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get existing requests for this key
	requests := rl.requests[key]

	// Filter out old requests outside the window
	valid := make([]time.Time, 0, len(requests))
	for _, t := range requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	// Check if limit exceeded
	if len(valid) >= rl.limit {
		rl.requests[key] = valid
		return false
	}

	// Add current request
	rl.requests[key] = append(valid, now)
	return true
}

// GetRemaining returns the remaining requests for a key
func (rl *RateLimiter) GetRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	requests := rl.requests[key]
	count := 0
	for _, t := range requests {
		if t.After(windowStart) {
			count++
		}
	}

	remaining := rl.limit - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetResetTime returns the time when the rate limit resets
func (rl *RateLimiter) GetResetTime(key string) time.Time {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	requests := rl.requests[key]
	if len(requests) == 0 {
		return time.Now().Add(rl.window)
	}

	// The oldest request in the window determines when it expires
	oldest := requests[0]
	for _, t := range requests {
		if t.Before(oldest) {
			oldest = t
		}
	}

	return oldest.Add(rl.window)
}

// cleanup periodically removes old entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window / 2)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for key, requests := range rl.requests {
			valid := make([]time.Time, 0, len(requests))
			for _, t := range requests {
				if t.After(windowStart) {
					valid = append(valid, t)
				}
			}

			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Check rate limit
		if !limiter.Allow(clientIP) {
			resetTime := limiter.GetResetTime(clientIP)
			remaining := limiter.GetRemaining(clientIP)

			// Log rate limit exceeded
			if log != nil {
				log.WithFields(logger.Fields{
					"client_ip":    clientIP,
					"path":         c.Request.URL.Path,
					"method":       c.Request.Method,
					"retry_after":  resetTime.Unix() - time.Now().Unix(),
					"remaining":    remaining,
				}).Warn("Rate limit exceeded")
			}

			// Return 429 Too Many Requests
			c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Header("Retry-After", strconv.FormatInt(resetTime.Unix()-time.Now().Unix(), 10))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, errors.
				TooManyRequests(int(resetTime.Unix() - time.Now().Unix())).
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		// Set rate limit headers
		remaining := limiter.GetRemaining(clientIP)
		resetTime := limiter.GetResetTime(clientIP)

		c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}

// ============================================================================
// ENDPOINT-SPECIFIC RATE LIMITING
// ============================================================================

// EndpointRateLimiter holds rate limiters for specific endpoints
type EndpointRateLimiter struct {
	limiters map[string]*RateLimiter
	mu       sync.RWMutex
}

// EndpointRateLimitConfig defines rate limits for an endpoint
type EndpointRateLimitConfig struct {
	Path   string
	Limit  int
	Window time.Duration
}

// NewEndpointRateLimiter creates rate limiters for specific endpoints
func NewEndpointRateLimiter(configs []EndpointRateLimitConfig) *EndpointRateLimiter {
	erl := &EndpointRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}

	for _, config := range configs {
		erl.limiters[config.Path] = NewRateLimiter(RateLimiterConfig{
			Limit:  config.Limit,
			Window: config.Window,
		})
	}

	return erl
}

// GetLimiter returns the rate limiter for a path
func (erl *EndpointRateLimiter) GetLimiter(path string) *RateLimiter {
	erl.mu.RLock()
	defer erl.mu.RUnlock()

	// Check for exact match
	if limiter, exists := erl.limiters[path]; exists {
		return limiter
	}

	// Check for prefix matches
	for endpointPath, limiter := range erl.limiters {
		if len(path) >= len(endpointPath) && path[:len(endpointPath)] == endpointPath {
			return limiter
		}
	}

	return nil
}

// EndpointRateLimitMiddleware creates middleware for endpoint-specific rate limiting
func EndpointRateLimitMiddleware(erl *EndpointRateLimiter, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		limiter := erl.GetLimiter(path)

		if limiter != nil {
			clientIP := c.ClientIP()

			if !limiter.Allow(clientIP) {
				resetTime := limiter.GetResetTime(clientIP)

				if log != nil {
					log.WithFields(logger.Fields{
						"client_ip":   clientIP,
						"path":        path,
						"method":      c.Request.Method,
						"retry_after": resetTime.Unix() - time.Now().Unix(),
					}).Warn("Endpoint rate limit exceeded")
				}

				c.Header("Retry-After", strconv.FormatInt(resetTime.Unix()-time.Now().Unix(), 10))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, errors.
					TooManyRequests(int(resetTime.Unix() - time.Now().Unix())).
					WithPath(path).
					ToResponse())
				return
			}
		}

		c.Next()
	}
}

// ============================================================================
// PREDEFINED RATE LIMIT CONFIGS
// ============================================================================

// DefaultRateLimitConfigs returns default rate limit configurations
func DefaultRateLimitConfigs() []EndpointRateLimitConfig {
	return []EndpointRateLimitConfig{
		{Path: "/api/auth/login", Limit: 5, Window: time.Minute},
		{Path: "/api/auth/register", Limit: 3, Window: time.Minute},
		{Path: "/api/auth/forgot-password", Limit: 3, Window: time.Minute},
		{Path: "/api/auth/reset-password", Limit: 5, Window: time.Minute},
		{Path: "/api/auth/refresh", Limit: 10, Window: time.Minute},
		{Path: "/api/admin", Limit: 50, Window: time.Minute},
		{Path: "/api", Limit: 100, Window: time.Minute},
	}
}

// CreateDefaultRateLimiter creates a rate limiter with default configs
func CreateDefaultRateLimiter(log *logger.Logger) *EndpointRateLimiter {
	return NewEndpointRateLimiter(DefaultRateLimitConfigs())
}

// ============================================================================
// IP-BASED KEY GENERATION
// ============================================================================

// GetClientIP extracts the client IP from the request
// It handles X-Forwarded-For and X-Real-IP headers
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (for proxied requests)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}
