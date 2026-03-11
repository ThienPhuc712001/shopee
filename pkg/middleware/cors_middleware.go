package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// CORS creates a CORS middleware
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Allow all origins in development
		if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// Check if origin is allowed
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// If origin not in list, still allow it (for development)
				if origin != "" {
					c.Header("Access-Control-Allow-Origin", origin)
				}
			}
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Forwarded-For")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter creates a rate limiting middleware
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(requestsPerMinute),
		burst:    burst,
	}
}

// RateLimit creates a rate limiting middleware
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			limiter := rate.NewLimiter(rl.rate, rl.burst)
			v = &visitor{limiter: limiter}
			rl.visitors[ip] = v
		}
		v.lastSeen = time.Now()
		rl.mu.Unlock()

		if !v.limiter.Allow() {
			c.JSON(429, gin.H{
				"success": false,
				"error":   "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Cleanup removes old visitors
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

// StartCleanup starts the cleanup goroutine
func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			rl.Cleanup()
		}
	}()
}
