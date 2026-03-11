// Package errors provides custom error types and error handling utilities
// for the e-commerce platform.
package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// ============================================================================
// ERROR CODES
// ============================================================================

// ErrorCode represents a unique error code for the application
type ErrorCode string

const (
	// General errors
	ErrUnknown          ErrorCode = "UNKNOWN"
	ErrInternal         ErrorCode = "INTERNAL"
	ErrNotImplemented   ErrorCode = "NOT_IMPLEMENTED"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// Client errors (4xx)
	ErrBadRequest       ErrorCode = "BAD_REQUEST"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrForbidden        ErrorCode = "FORBIDDEN"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrConflict         ErrorCode = "CONFLICT"
	ErrValidationError  ErrorCode = "VALIDATION_ERROR"
	ErrTooManyRequests  ErrorCode = "TOO_MANY_REQUESTS"

	// Database errors
	ErrDatabase         ErrorCode = "DATABASE_ERROR"
	ErrRecordNotFound   ErrorCode = "RECORD_NOT_FOUND"
	ErrDuplicateKey     ErrorCode = "DUPLICATE_KEY"
	ErrConnectionFailed ErrorCode = "CONNECTION_FAILED"

	// Authentication errors
	ErrInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrTokenInvalid       ErrorCode = "TOKEN_INVALID"
	ErrSessionExpired     ErrorCode = "SESSION_EXPIRED"

	// Resource errors
	ErrResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ErrResourceExists   ErrorCode = "RESOURCE_EXISTS"
	ErrQuotaExceeded    ErrorCode = "QUOTA_EXCEEDED"
)

// ============================================================================
// APP ERROR STRUCT
// ============================================================================

// AppError represents a structured application error
type AppError struct {
	// Code is a unique error code for programmatic handling
	Code ErrorCode `json:"code"`

	// Message is a user-friendly error message
	Message string `json:"message"`

	// Details contains additional error information (optional)
	Details interface{} `json:"details,omitempty"`

	// HTTPStatus is the HTTP status code to return
	HTTPStatus int `json:"-"`

	// Err is the underlying error (for logging, not exposed to client)
	Err error `json:"-"`

	// Timestamp when the error occurred
	Timestamp time.Time `json:"timestamp"`

	// Path where the error occurred
	Path string `json:"path,omitempty"`

	// Method is the HTTP method that caused the error
	Method string `json:"method,omitempty"`

	// Stack trace (for internal debugging)
	stack []uintptr `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error for errors.Is/As
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// WithPath adds the request path to the error
func (e *AppError) WithPath(path string) *AppError {
	e.Path = path
	return e
}

// WithMethod adds the HTTP method to the error
func (e *AppError) WithMethod(method string) *AppError {
	e.Method = method
	return e
}

// WithError sets the underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// captureStack captures the current call stack
func (e *AppError) captureStack() {
	e.stack = make([]uintptr, 32)
	runtime.Callers(2, e.stack)
}

// ============================================================================
// ERROR RESPONSE
// ============================================================================

// ErrorResponse is the standard API error response sent to clients
type ErrorResponse struct {
	Error     bool        `json:"error"`
	Message   string      `json:"message"`
	Code      ErrorCode   `json:"code"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Path      string      `json:"path,omitempty"`
}

// ToResponse converts AppError to ErrorResponse for API responses
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Error:     true,
		Message:   e.Message,
		Code:      e.Code,
		Details:   e.Details,
		Timestamp: e.Timestamp,
		Path:      e.Path,
	}
}

// ============================================================================
// ERROR FACTORY FUNCTIONS
// ============================================================================

// New creates a new AppError with the given code and message
func New(code ErrorCode, message string) *AppError {
	err := &AppError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		HTTPStatus: getHTTPStatus(code),
	}
	err.captureStack()
	return err
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *AppError {
	appErr := &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		Timestamp: time.Now(),
		HTTPStatus: getHTTPStatus(code),
	}
	appErr.captureStack()
	return appErr
}

// ============================================================================
// CONVENIENCE ERROR CREATORS
// ============================================================================

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *AppError {
	return New(ErrBadRequest, message)
}

// InvalidInput creates a validation error for invalid input
func InvalidInput(field, reason string) *AppError {
	return New(ErrValidationError, fmt.Sprintf("Invalid %s: %s", field, reason)).
		WithDetails(map[string]string{"field": field, "reason": reason})
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *AppError {
	return New(ErrUnauthorized, message)
}

// InvalidCredentials creates an authentication error
func InvalidCredentials() *AppError {
	return New(ErrInvalidCredentials, "Invalid email or password")
}

// TokenExpired creates a token expiration error
func TokenExpired() *AppError {
	return New(ErrTokenExpired, "Authentication token has expired")
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *AppError {
	return New(ErrForbidden, message)
}

// NotFound creates a 404 Not Found error
func NotFound(resource string) *AppError {
	return New(ErrNotFound, fmt.Sprintf("%s not found", resource))
}

// RecordNotFound creates a database record not found error
func RecordNotFound(resource string) *AppError {
	return New(ErrRecordNotFound, fmt.Sprintf("%s record not found", resource))
}

// Conflict creates a 409 Conflict error
func Conflict(message string) *AppError {
	return New(ErrConflict, message)
}

// DuplicateKey creates a duplicate key error
func DuplicateKey(field string) *AppError {
	return New(ErrDuplicateKey, fmt.Sprintf("A record with this %s already exists", field)).
		WithDetails(map[string]string{"field": field})
}

// Internal creates a 500 Internal Server Error
func Internal(message string) *AppError {
	return New(ErrInternal, message)
}

// Database creates a database error
func Database(message string) *AppError {
	return New(ErrDatabase, message)
}

// ConnectionFailed creates a connection error
func ConnectionFailed(service string) *AppError {
	return New(ErrConnectionFailed, fmt.Sprintf("Failed to connect to %s", service))
}

// TooManyRequests creates a 429 rate limit error
func TooManyRequests(retryAfter int) *AppError {
	return New(ErrTooManyRequests, "Too many requests").
		WithDetails(map[string]int{"retry_after": retryAfter})
}

// ServiceUnavailable creates a 503 error
func ServiceUnavailable(service string) *AppError {
	return New(ErrServiceUnavailable, fmt.Sprintf("%s is currently unavailable", service))
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// getHTTPStatus returns the HTTP status code for an error code
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrBadRequest, ErrValidationError:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrInvalidCredentials, ErrTokenExpired, ErrTokenInvalid:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound, ErrRecordNotFound, ErrResourceNotFound:
		return http.StatusNotFound
	case ErrConflict, ErrDuplicateKey, ErrResourceExists:
		return http.StatusConflict
	case ErrTooManyRequests:
		return http.StatusTooManyRequests
	case ErrDatabase, ErrConnectionFailed:
		return http.StatusInternalServerError
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// Is checks if an error is of a specific error code
func Is(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// GetCode extracts the error code from an error
func GetCode(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrUnknown
}

// GetHTTPStatus extracts the HTTP status from an error
func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	return Is(err, ErrNotFound) || Is(err, ErrRecordNotFound) || Is(err, ErrResourceNotFound)
}

// IsUnauthorized checks if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	return Is(err, ErrUnauthorized) || Is(err, ErrInvalidCredentials) ||
		Is(err, ErrTokenExpired) || Is(err, ErrTokenInvalid)
}

// IsBadRequest checks if the error is a bad request error
func IsBadRequest(err error) bool {
	return Is(err, ErrBadRequest) || Is(err, ErrValidationError)
}
