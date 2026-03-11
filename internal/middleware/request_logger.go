package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ecommerce/pkg/logger"
)

// responseWriter is a custom response writer that captures
// the status code and response size
type responseWriter struct {
	gin.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// RequestLogMiddleware returns a middleware that logs every incoming request
// with detailed information for debugging and monitoring
func RequestLogMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Capture request details
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		// Get request ID for tracing
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Header("X-Request-ID", requestID)
		}

		// Wrap response writer to capture status and size
		w := &responseWriter{
			ResponseWriter: c.Writer,
			status:         http.StatusOK, // Default status
		}
		c.Writer = w

		// Execute request
		c.Next()

		// Calculate response time
		duration := time.Since(start)

		// Get user ID if available
		var userID interface{}
		if id, exists := c.Get("user_id"); exists {
			userID = id
		}

		// Prepare log fields
		fields := logger.Fields{
			"request_id":    requestID,
			"method":        method,
			"path":          path,
			"query":         query,
			"status":        w.status,
			"response_time": duration.Milliseconds(),
			"response_size": w.size,
			"client_ip":     clientIP,
			"user_agent":    userAgent,
			"referer":       referer,
		}

		if userID != nil {
			fields["user_id"] = userID
		}

		// Log based on status code
		switch {
		case w.status >= http.StatusInternalServerError:
			log.WithFields(fields).Error("Server error response")
		case w.status >= http.StatusBadRequest:
			log.WithFields(fields).Warning("Client error response")
		default:
			log.WithFields(fields).Info("Request processed")
		}
	}
}

// SlowRequestMiddleware returns a middleware that logs requests
// that take longer than a threshold
func SlowRequestMiddleware(log *logger.Logger, threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		if duration > threshold {
			log.WithFields(logger.Fields{
				"method":        c.Request.Method,
				"path":          c.Request.URL.Path,
				"duration_ms":   duration.Milliseconds(),
				"threshold_ms":  threshold.Milliseconds(),
				"status":        c.Writer.Status(),
				"client_ip":     c.ClientIP(),
			}).Warning("Slow request detected")
		}
	}
}

// MetricsMiddleware returns a middleware that collects request metrics
// Can be extended to integrate with Prometheus or other monitoring tools
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		path := c.Request.URL.Path

		// TODO: Integrate with metrics collection system
		// Example: metrics.RequestDuration.WithLabelValues(path, status).Observe(duration.Seconds())
		// Example: metrics.RequestCount.WithLabelValues(path, status).Inc()

		_ = duration
		_ = status
		_ = path
	}
}

// generateRequestID generates a unique request ID for tracing
// In production, consider using uuid or xid library
func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
