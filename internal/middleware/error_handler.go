// Package middleware provides Gin middleware for error handling and logging
package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	apperrors "ecommerce/internal/errors"
	"ecommerce/pkg/logger"
)

// RecoveryMiddleware returns a middleware that recovers from panics
// and returns a proper error response
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.WithFields(logger.Fields{
					"error":      err,
					"stack":      string(debug.Stack()),
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"client_ip":  c.ClientIP(),
				}).Error("Panic recovered")

				// Return internal server error response
				c.JSON(http.StatusInternalServerError, apperrors.
					Internal("An unexpected error occurred").
					WithPath(c.Request.URL.Path).
					WithMethod(c.Request.Method).
					ToResponse())

				c.Abort()
			}
		}()

		c.Next()
	}
}

// ErrorMiddleware returns a middleware that handles errors returned
// by handlers and formats them consistently
func ErrorMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors in the context
		if len(c.Errors) > 0 {
			handleContextErrors(c, log)
			return
		}
	}
}

// handleContextErrors processes errors stored in Gin context
func handleContextErrors(c *gin.Context, log *logger.Logger) {
	for _, ginErr := range c.Errors {
		var appErr *apperrors.AppError

		// Try to extract AppError from the error
		if errors.As(ginErr.Err, &appErr) {
			// Log the error with context
			logError(log, c, appErr)

			// Return error response
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
			c.Abort()
			return
		}

		// Handle non-AppError errors
		log.WithFields(logger.Fields{
			"error":     ginErr.Err,
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"client_ip": c.ClientIP(),
		}).Error("Unhandled error in context")

		// Return generic internal error
		c.JSON(http.StatusInternalServerError, apperrors.
			Internal("An unexpected error occurred").
			WithPath(c.Request.URL.Path).
			WithMethod(c.Request.Method).
			ToResponse())
		c.Abort()
		return
	}
}

// logError logs the error with appropriate level based on error type
func logError(log *logger.Logger, c *gin.Context, appErr *apperrors.AppError) {
	fields := logger.Fields{
		"error_code":  appErr.Code,
		"error_msg":   appErr.Message,
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"client_ip":   c.ClientIP(),
		"user_agent":  c.Request.UserAgent(),
	}

	// Add user ID if available
	if userID, exists := c.Get("user_id"); exists {
		fields["user_id"] = userID
	}

	// Add details if present
	if appErr.Details != nil {
		fields["details"] = appErr.Details
	}

	// Add underlying error if present (for debugging)
	if appErr.Err != nil {
		fields["underlying_error"] = appErr.Err.Error()
	}

	// Log with appropriate level
	switch appErr.Code {
	case apperrors.ErrInternal, apperrors.ErrDatabase, apperrors.ErrConnectionFailed, apperrors.ErrServiceUnavailable:
		log.WithFields(fields).Error("Server error occurred")
	case apperrors.ErrUnauthorized, apperrors.ErrForbidden, apperrors.ErrInvalidCredentials:
		log.WithFields(fields).Warning("Authentication/Authorization error")
	default:
		log.WithFields(fields).Info("Client error occurred")
	}
}

// ErrorHandlerFunc is a handler function that returns an error
type ErrorHandlerFunc func(c *gin.Context) error

// Handle wraps an ErrorHandlerFunc and handles errors automatically
func (h ErrorHandlerFunc) Handle(c *gin.Context) {
	if err := h(c); err != nil {
		c.Error(err)
	}
}

// HandleError is a helper to handle and return error in handlers
func HandleError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
		c.Abort()
	}
}

// AbortWithError aborts the request with an error
func AbortWithError(c *gin.Context, err *apperrors.AppError) {
	c.AbortWithStatusJSON(err.HTTPStatus, err.ToResponse())
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    Pagination  `json:"meta"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Paginated sends a paginated success response
func Paginated(c *gin.Context, data interface{}, page, limit int, total int64) {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
