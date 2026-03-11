// Package dberrors provides database error handling utilities
package dberrors

import (
	"strings"

	"ecommerce/internal/errors"
	"gorm.io/gorm"
)

// ============================================================================
// DATABASE ERROR TYPES
// ============================================================================

// DBErrorType represents the type of database error
type DBErrorType string

const (
	// ErrNotFound indicates a record was not found
	ErrNotFound DBErrorType = "NOT_FOUND"

	// ErrDuplicate indicates a duplicate key violation
	ErrDuplicate DBErrorType = "DUPLICATE"

	// ErrForeignKey indicates a foreign key violation
	ErrForeignKey DBErrorType = "FOREIGN_KEY"

	// ErrNotNull indicates a not null constraint violation
	ErrNotNull DBErrorType = "NOT_NULL"

	// ErrCheck indicates a check constraint violation
	ErrCheck DBErrorType = "CHECK"

	// ErrUnique indicates a unique constraint violation
	ErrUnique DBErrorType = "UNIQUE"

	// ErrConnection indicates a connection error
	ErrConnection DBErrorType = "CONNECTION"

	// ErrTimeout indicates a query timeout
	ErrTimeout DBErrorType = "TIMEOUT"

	// ErrUnknown indicates an unknown database error
	ErrUnknown DBErrorType = "UNKNOWN"
)

// ============================================================================
// ERROR HANDLER
// ============================================================================

// Handler handles database errors and converts them to application errors
type Handler struct {
	// CustomMessages allows customizing error messages for specific tables
	CustomMessages map[string]map[DBErrorType]string
}

// NewHandler creates a new database error handler
func NewHandler() *Handler {
	return &Handler{
		CustomMessages: make(map[string]map[DBErrorType]string),
	}
}

// WithCustomMessage adds a custom error message for a table and error type
func (h *Handler) WithCustomMessage(table string, errType DBErrorType, message string) *Handler {
	if h.CustomMessages[table] == nil {
		h.CustomMessages[table] = make(map[DBErrorType]string)
	}
	h.CustomMessages[table][errType] = message
	return h
}

// HandleError converts a database error to an application error
func (h *Handler) HandleError(err error, resource string) error {
	if err == nil {
		return nil
	}

	// Handle gorm.ErrRecordNotFound
	if err == gorm.ErrRecordNotFound {
		return errors.RecordNotFound(resource)
	}

	// Get the error type
	errType := h.getErrorType(err)

	// Check for custom message
	if messages, ok := h.CustomMessages[resource]; ok {
		if msg, ok := messages[errType]; ok {
			return errors.Database(msg).WithError(err)
		}
	}

	// Return appropriate error based on type
	switch errType {
	case ErrNotFound:
		return errors.RecordNotFound(resource)
	case ErrDuplicate, ErrUnique:
		return errors.DuplicateKey(resource)
	case ErrForeignKey:
		return errors.BadRequest("Related record not found")
	case ErrNotNull:
		return errors.BadRequest("Required field is missing")
	case ErrConnection:
		return errors.ConnectionFailed("database")
	case ErrTimeout:
		return errors.ServiceUnavailable("database")
	default:
		return errors.Database("A database error occurred").WithError(err)
	}
}

// getErrorType determines the type of database error
func (h *Handler) getErrorType(err error) DBErrorType {
	errStr := strings.ToLower(err.Error())

	// Check for common SQL Server error patterns
	switch {
	case strings.Contains(errStr, "primary key constraint"):
		return ErrDuplicate
	case strings.Contains(errStr, "unique constraint"):
		return ErrUnique
	case strings.Contains(errStr, "duplicate key"):
		return ErrDuplicate
	case strings.Contains(errStr, "foreign key constraint"):
		return ErrForeignKey
	case strings.Contains(errStr, "cannot insert null"):
		return ErrNotNull
	case strings.Contains(errStr, "check constraint"):
		return ErrCheck
	case strings.Contains(errStr, "connection"):
		return ErrConnection
	case strings.Contains(errStr, "timeout"):
		return ErrTimeout
	case strings.Contains(errStr, "lock"):
		return ErrTimeout
	default:
		return ErrUnknown
	}
}

// ============================================================================
// CONVENIENCE FUNCTIONS
// ============================================================================

// HandleNotFound returns a standardized not found error
func HandleNotFound(resource string) error {
	return errors.RecordNotFound(resource)
}

// HandleDuplicate returns a standardized duplicate error
func HandleDuplicate(field string) error {
	return errors.DuplicateKey(field)
}

// HandleConnection returns a standardized connection error
func HandleConnection() error {
	return errors.ConnectionFailed("database")
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound || errors.IsNotFound(err)
}

// IsDuplicate checks if an error is a duplicate key error
func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "primary key constraint")
}

// ============================================================================
// GORM ERROR WRAPPER
// ============================================================================

// WrapGormError wraps a GORM error with context
func WrapGormError(err error, operation string, resource string) error {
	if err == nil {
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		return errors.RecordNotFound(resource)
	}

	return errors.Database("Failed to " + operation + " " + resource).
		WithError(err)
}

// CheckGormError checks GORM query result and returns appropriate error
func CheckGormError(result *gorm.DB) error {
	if result.Error == nil {
		return nil
	}

	if result.Error == gorm.ErrRecordNotFound {
		return errors.RecordNotFound("record")
	}

	return errors.Database("Database query failed").
		WithError(result.Error)
}

// ============================================================================
// SQL SERVER SPECIFIC ERRORS
// ============================================================================

// SQL Server error numbers (for reference)
// 2627 - Unique constraint violation
// 2601 - Duplicate key error
// 547 - Foreign key constraint violation
// 515 - Cannot insert NULL
// 513 - Invalid value for column
// 1205 - Deadlock victim
// 1222 - Lock request timeout

// HandleSQLServerError handles SQL Server specific errors
func HandleSQLServerError(err error, resource string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Check for SQL Server specific error patterns
	switch {
	case strings.Contains(errStr, "Violation of UNIQUE KEY constraint"):
		return errors.DuplicateKey(resource)
	case strings.Contains(errStr, "Cannot insert duplicate key"):
		return errors.DuplicateKey(resource)
	case strings.Contains(errStr, "The INSERT statement conflicted with the FOREIGN KEY constraint"):
		return errors.BadRequest("Related record does not exist")
	case strings.Contains(errStr, "Cannot insert the value NULL"):
		return errors.BadRequest("Required field cannot be null")
	case strings.Contains(errStr, "deadlock"):
		return errors.ServiceUnavailable("database")
	default:
		return errors.Database("Database operation failed").
			WithError(err)
	}
}
