# Error Handling and Logging System Documentation

## Overview

This document describes the Error Handling and Logging System for the e-commerce platform. The system provides consistent error responses and structured logging for debugging and monitoring.

---

## Table of Contents

1. [Why Proper Error Handling is Important](#1-why-proper-error-handling-is-important)
2. [Standard Error Response Format](#2-standard-error-response-format)
3. [Error Types](#3-error-types)
4. [Custom Error Structure](#4-custom-error-structure)
5. [Global Error Handler](#5-global-error-handler)
6. [Logging System](#6-logging-system)
7. [Request Logging](#7-request-logging)
8. [Database Error Handling](#8-database-error-handling)
9. [Log File Management](#9-log-file-management)
10. [Security Considerations](#10-security-considerations)
11. [Error Monitoring](#11-error-monitoring)
12. [Example Implementation](#12-example-implementation)

---

## 1. Why Proper Error Handling is Important

Proper error handling is crucial for building reliable and maintainable backend systems. Here's why:

### Importance

- **User Experience**: Clear error messages help users understand what went wrong
- **Debugging**: Structured errors make it easier to identify and fix issues
- **Security**: Proper error handling prevents information leakage
- **Monitoring**: Consistent errors enable better system health tracking
- **API Contract**: Standard error responses maintain API consistency

### Common Error Scenarios

| Error Type | Example | Impact |
|------------|---------|--------|
| **Database Errors** | Connection failure, query timeout, constraint violation | Service unavailable, data inconsistency |
| **Invalid User Input** | Missing required fields, invalid format, out-of-range values | Bad requests, validation failures |
| **Authentication Errors** | Invalid credentials, expired token, unauthorized access | Security breaches, unauthorized access |
| **Server Failures** | Out of memory, disk full, external service failure | Complete service outage |

### Consequences of Poor Error Handling

1. **Security Vulnerabilities**: Exposing stack traces or SQL errors
2. **Poor User Experience**: Confusing or generic error messages
3. **Debugging Nightmares**: No context or logs to trace issues
4. **System Instability**: Unhandled errors causing crashes
5. **Compliance Issues**: Logging sensitive data improperly

---

## 2. Standard Error Response Format

All API errors follow a consistent JSON structure:

```json
{
  "error": true,
  "message": "Invalid request",
  "code": "BAD_REQUEST",
  "details": {
    "field": "email",
    "reason": "Invalid format"
  },
  "timestamp": "2026-03-11T08:00:00Z",
  "path": "/api/users"
}
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `error` | boolean | Always `true` for error responses |
| `message` | string | User-friendly error message |
| `code` | string | Machine-readable error code |
| `details` | object | Optional additional context |
| `timestamp` | string | ISO 8601 timestamp of the error |
| `path` | string | Request path that caused the error |

---

## 3. Error Types

### Client Errors (4xx)

| Error Code | HTTP Status | When to Use |
|------------|-------------|-------------|
| `BAD_REQUEST` | 400 | Malformed request syntax |
| `VALIDATION_ERROR` | 400 | Invalid input validation |
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `INVALID_CREDENTIALS` | 401 | Wrong email/password |
| `TOKEN_EXPIRED` | 401 | JWT token has expired |
| `FORBIDDEN` | 403 | Valid auth but insufficient permissions |
| `NOT_FOUND` | 404 | Resource does not exist |
| `RECORD_NOT_FOUND` | 404 | Database record not found |
| `CONFLICT` | 409 | Resource conflict (e.g., duplicate email) |
| `DUPLICATE_KEY` | 409 | Unique constraint violation |
| `TOO_MANY_REQUESTS` | 429 | Rate limit exceeded |

### Server Errors (5xx)

| Error Code | HTTP Status | When to Use |
|------------|-------------|-------------|
| `INTERNAL` | 500 | Unexpected server error |
| `DATABASE_ERROR` | 500 | Database operation failed |
| `CONNECTION_FAILED` | 500 | Cannot connect to external service |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

---

## 4. Custom Error Structure

The `AppError` struct provides a flexible error representation:

```go
type AppError struct {
    Code       ErrorCode   `json:"code"`
    Message    string      `json:"message"`
    Details    interface{} `json:"details,omitempty"`
    HTTPStatus int         `json:"-"`
    Err        error       `json:"-"`
    Timestamp  time.Time   `json:"timestamp"`
    Path       string      `json:"path,omitempty"`
    Method     string      `json:"method,omitempty"`
}
```

### Creating Errors

```go
// Using convenience functions
errors.BadRequest("Invalid input")
errors.NotFound("User not found")
errors.Unauthorized("Invalid token")
errors.Internal("Something went wrong")

// With details
errors.InvalidInput("email", "Invalid format").
    WithDetails(map[string]string{"field": "email"})

// Wrapping underlying errors
errors.Wrap(dbErr, errors.ErrDatabase, "Failed to fetch user")
```

---

## 5. Global Error Handler

The error handler middleware automatically catches and formats errors:

### Middleware Setup

```go
// In main.go
r := gin.New()
r.Use(middleware.RecoveryMiddleware(log))
r.Use(middleware.ErrorMiddleware(log))
```

### Responsibilities

1. **Catch Panics**: Recover from panics and return 500 error
2. **Handle Errors**: Process errors returned by handlers
3. **Log Errors**: Log with appropriate level based on error type
4. **Format Response**: Return standard error JSON response

### Usage in Handlers

```go
func (h *Handler) GetUser(c *gin.Context) {
    user, err := h.repo.FindByID(id)
    if err != nil {
        middleware.AbortWithError(c, err)
        return
    }
    middleware.Success(c, user, "Success")
}
```

---

## 6. Logging System

The logging system uses Logrus with structured logging:

### Log Levels

| Level | When to Use |
|-------|-------------|
| `DEBUG` | Detailed diagnostic information (development) |
| `INFO` | Normal operational messages |
| `WARN` | Unexpected but handled situations |
| `ERROR` | Error conditions that need attention |
| `FATAL` | Critical errors that stop the application |
| `PANIC` | Unrecoverable errors |

### Usage Examples

```go
// Simple logging
log.Info("Server started on port 8080")
log.Error("Database connection failed")

// With fields
log.WithFields(logger.Fields{
    "user_id": userID,
    "action":  "login",
}).Info("User logged in")

// With error
log.WithError(err).Error("Query failed")

// Formatted logging
log.Infof("User %s created order %d", userID, orderID)
```

---

## 7. Request Logging

Every request is automatically logged with:

```json
{
  "request_id": "20260311080000.000000",
  "method": "POST",
  "path": "/api/users",
  "query": "include=profile",
  "status": 201,
  "response_time": 45,
  "response_size": 1024,
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "user_id": 123
}
```

### Benefits

- **Audit Trail**: Track all API access
- **Performance Monitoring**: Identify slow endpoints
- **Debugging**: Trace request flow
- **Security**: Detect suspicious patterns
- **Analytics**: Understand API usage patterns

---

## 8. Database Error Handling

The `dberrors` package handles common database errors:

### Error Types

```go
// Record not found
dberrors.HandleNotFound("User")

// Duplicate key
dberrors.HandleDuplicate("email")

// Connection failure
dberrors.HandleConnection()

// Wrap GORM errors
dberrors.WrapGormError(err, "fetch", "user")
```

### Usage Pattern

```go
func (r *Repository) FindByID(id int64) (*User, error) {
    var user User
    result := r.db.First(&user, id)
    
    if err := dberrors.CheckGormError(result); err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

---

## 9. Log File Management

### File Structure

```
logs/
├── app.log          # All application logs
└── error.log        # Error-level logs only
```

### Log Rotation

Logs are automatically rotated based on:

| Setting | Default | Description |
|---------|---------|-------------|
| `MaxSize` | 100 MB | Max file size before rotation |
| `MaxBackups` | 30 | Number of old files to keep |
| `MaxAge` | 90 days | Maximum age of log files |
| `Compress` | true | Compress rotated files |

### Configuration

```go
logConfig := &logger.Config{
    FilePath:      "logs/app.log",
    ErrorFilePath: "logs/error.log",
    MaxSize:       100,
    MaxBackups:    30,
    MaxAge:        90,
    Compress:      true,
}
```

---

## 10. Security Considerations

### What NOT to Expose

❌ **Never expose in API responses:**

- SQL error messages
- Stack traces
- Database schema details
- Internal IP addresses
- API keys or secrets
- User passwords or tokens

### Safe Error Messages

| Instead of | Use |
|------------|-----|
| "SQL: SELECT * FROM users WHERE..." | "Invalid request" |
| "panic: nil pointer dereference at..." | "Internal server error" |
| "Connection string: server=..." | "Database unavailable" |

### Implementation

```go
// Bad - exposes internal details
return errors.Internal(err.Error())

// Good - safe message
return errors.Internal("An unexpected error occurred")
```

---

## 11. Error Monitoring

### Using Logs for Monitoring

1. **Detect Frequent Errors**: Track error codes over time
2. **Track System Failures**: Monitor 5xx error rates
3. **Debug Production Issues**: Use request_id for tracing
4. **Performance Analysis**: Track slow requests
5. **Security Monitoring**: Detect auth failures, suspicious patterns

### Key Metrics to Track

- Error rate by endpoint
- Error rate by error code
- Average response time
- Slow request count
- Authentication failures
- Database error rate

---

## 12. Example Implementation

### Complete Handler Example

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/cntt.tts13/tmdt/internal/errors"
    "github.com/cntt.tts13/tmdt/internal/middleware"
    "github.com/cntt.tts13/tmdt/pkg/logger"
)

type UserHandler struct {
    log *logger.Logger
    repo *repository.UserRepository
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // Bind and validate
    if err := c.ShouldBindJSON(&req); err != nil {
        middleware.AbortWithError(c, 
            errors.BadRequest("Invalid request body"))
        return
    }
    
    // Validate email
    if req.Email == "" {
        middleware.AbortWithError(c, 
            errors.InvalidInput("email", "is required"))
        return
    }
    
    // Check existing user
    existing, err := h.repo.FindByEmail(req.Email)
    if err != nil && !errors.IsNotFound(err) {
        h.log.WithError(err).Error("Failed to check existing user")
        middleware.AbortWithError(c, 
            errors.Database("Failed to create user"))
        return
    }
    
    if existing != nil {
        middleware.AbortWithError(c, 
            errors.Conflict("Email already registered"))
        return
    }
    
    // Create user
    user, err := h.repo.Create(&User{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
    })
    
    if err != nil {
        h.log.WithFields(logger.Fields{
            "email": req.Email,
            "error": err,
        }).Error("Failed to create user")
        
        middleware.AbortWithError(c, 
            errors.Database("Failed to create user"))
        return
    }
    
    h.log.WithFields(logger.Fields{
        "user_id": user.ID,
        "email":   user.Email,
    }).Info("User created successfully")
    
    middleware.Created(c, user, "User created successfully")
}
```

### Main Application Setup

```go
func main() {
    // Initialize logger
    log, _ := logger.New(&logger.Config{
        Level:         "info",
        Format:        "json",
        FilePath:      "logs/app.log",
        ErrorFilePath: "logs/error.log",
        Service:       "ecommerce-api",
    })
    
    // Setup Gin with middleware
    r := gin.New()
    r.Use(middleware.Setup(log, middleware.DefaultConfig())...)
    
    // Register routes
    api := r.Group("/api")
    handler.RegisterRoutes(api, log)
    
    // Start server
    r.Run(":8080")
}
```

---

## Quick Reference

### Error Creation Cheat Sheet

```go
// Client errors
errors.BadRequest("message")
errors.InvalidInput("field", "reason")
errors.Unauthorized("message")
errors.Forbidden("message")
errors.NotFound("resource")
errors.Conflict("message")

// Database errors
errors.RecordNotFound("resource")
errors.DuplicateKey("field")
errors.Database("message")
errors.ConnectionFailed("service")

// Server errors
errors.Internal("message")
errors.ServiceUnavailable("service")

// With context
errors.BadRequest("message").
    WithDetails(map[string]interface{}{"key": "value"}).
    WithPath(c.Request.URL.Path).
    WithMethod(c.Request.Method)
```

### Logging Cheat Sheet

```go
// Basic logging
log.Debug("Debug message")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message")

// With fields
log.WithFields(logger.Fields{
    "user_id": 123,
    "action":  "create",
}).Info("User action")

// With error
log.WithError(err).Error("Operation failed")

// Formatted
log.Infof("Processed %d items in %dms", count, duration)
```

---

## Support

For questions or issues, refer to:
- Package documentation in `internal/errors/`
- Logger documentation in `pkg/logger/`
- Example code in `cmd/server/main_example.go`
