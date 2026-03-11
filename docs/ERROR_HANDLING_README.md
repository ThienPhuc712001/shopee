# Error Handling & Logging System - Quick Start

## Overview

This system provides consistent error handling and structured logging for the e-commerce API.

## Files Created

```
internal/
├── errors/
│   └── errors.go          # Custom error types and helpers
├── middleware/
│   ├── error_handler.go   # Error handling middleware
│   ├── request_logger.go  # Request logging middleware
│   └── middleware.go      # Middleware initialization
└── handler/
    └── example_handler.go # Example handlers

pkg/
├── logger/
│   └── logger.go          # Logrus-based logging system
└── dberrors/
    └── dberrors.go        # Database error utilities

cmd/
└── server/
    └── main_example.go    # Example main.go

docs/
└── ERROR_HANDLING_GUIDE.md # Full documentation
```

## Quick Setup

### 1. Initialize Logger

```go
log, err := logger.New(&logger.Config{
    Level:         "info",
    Format:        "text",
    FilePath:      "logs/app.log",
    ErrorFilePath: "logs/error.log",
    Stdout:        true,
    Service:       "ecommerce-api",
})
```

### 2. Setup Middleware

```go
r := gin.New()
r.Use(middleware.Setup(log, middleware.DefaultConfig())...)
```

### 3. Use in Handlers

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

## Error Types

| Function | HTTP Status | Use Case |
|----------|-------------|----------|
| `errors.BadRequest(msg)` | 400 | Invalid request |
| `errors.InvalidInput(field, reason)` | 400 | Validation error |
| `errors.Unauthorized(msg)` | 401 | Auth failed |
| `errors.NotFound(resource)` | 404 | Not found |
| `errors.Conflict(msg)` | 409 | Duplicate |
| `errors.Internal(msg)` | 500 | Server error |
| `errors.Database(msg)` | 500 | DB error |

## Log Levels

- `DEBUG` - Development debugging
- `INFO` - Normal operations
- `WARN` - Warnings
- `ERROR` - Errors (saved to error.log)
- `FATAL` - Critical failures

## Example Error Response

```json
{
  "error": true,
  "message": "Invalid email format",
  "code": "VALIDATION_ERROR",
  "details": {"field": "email", "reason": "Invalid format"},
  "timestamp": "2026-03-11T08:00:00Z",
  "path": "/api/users"
}
```

## Example Log Output

```
INFO 2026-03-11 08:00:00 | POST /api/users | 201 | 45ms | user_id=123
ERROR 2026-03-11 08:00:01 | Database error | code=DATABASE_ERROR | path=/api/users
```

## Documentation

See [ERROR_HANDLING_GUIDE.md](ERROR_HANDLING_GUIDE.md) for complete documentation.
