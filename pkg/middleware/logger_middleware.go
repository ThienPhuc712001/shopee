package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger creates a logging middleware
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get method
		method := c.Request.Method

		// Log request
		log.Printf("[%s] %s %s %d %v %s",
			method,
			path,
			query,
			statusCode,
			latency,
			clientIP,
		)

		// Log errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %v", e)
			}
		}
	}
}

// Recovery creates a panic recovery middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Printf("Panic recovered: %v", err)

				// Abort the request
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
