# Best Practices Guide

This document outlines best practices for developing, testing, and scaling this e-commerce backend.

## 1. Project Structure

### Package Organization

```
ecommerce/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go         # Main application
├── internal/               # Private application code
│   ├── domain/             # Domain layer (entities, interfaces)
│   │   ├── model/          # Domain entities
│   │   └── repository/     # Repository interfaces
│   ├── repository/         # Repository implementations
│   ├── service/            # Business logic
│   └── handler/            # HTTP handlers
├── pkg/                    # Public library code
│   ├── config/             # Configuration
│   ├── database/           # Database connection
│   ├── middleware/         # HTTP middleware
│   ├── response/           # API response helpers
│   └── utils/              # Utility functions
└── api/                    # Route definitions
```

### Key Principles

1. **Dependency Rule**: Dependencies point inward. Outer layers depend on inner layers.
2. **Interface Segregation**: Define interfaces in the domain layer, implement in outer layers.
3. **Single Responsibility**: Each package has one clear responsibility.

## 2. Testing Services

### Unit Testing

```go
// service/auth_service_test.go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    authService := service.NewAuthService(mockRepo, "secret", time.Hour)
    
    mockRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    // Act
    user, err := authService.Register("test@example.com", "password123", "John", "Doe")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    
    mockRepo.AssertExpectations(t)
}
```

### Integration Testing

```go
// tests/integration/order_test.go
func TestCreateOrder(t *testing.T) {
    // Setup test database
    db := setupTestDatabase()
    defer teardownTestDatabase(db)
    
    // Create test data
    user := createTestUser(db)
    product := createTestProduct(db)
    addToCart(db, user.ID, product.ID, 2)
    
    // Make API request
    req := httptest.NewRequest("POST", "/api/orders", orderBody)
    req.Header.Set("Authorization", "Bearer "+generateToken(user))
    
    // Assert response
    assert.Equal(t, 201, recorder.Code)
}
```

### Test Coverage Goals

- **Services**: 80%+ coverage
- **Handlers**: 70%+ coverage
- **Repositories**: 60%+ coverage (focus on complex queries)

## 3. Scaling the System

### Horizontal Scaling

1. **Stateless Design**: All services are stateless, enabling horizontal scaling.
2. **Session Management**: Use JWT tokens instead of server-side sessions.
3. **Database Connection Pooling**: Configure appropriate pool sizes.

```go
// Connection pool configuration
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Database Optimization

1. **Indexing**: Add indexes on frequently queried columns.
2. **Query Optimization**: Use `Preload` wisely to avoid N+1 queries.
3. **Read Replicas**: Separate read and write operations.

```go
// Example: Add indexes to models
type Product struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"index"`
    ShopID   uint   `gorm:"index"`
    CategoryID uint `gorm:"index"`
}
```

### Caching Strategy

1. **Redis for Sessions**: Store JWT blacklist in Redis.
2. **Product Cache**: Cache frequently accessed products.
3. **Cart Cache**: Store cart data in Redis for fast access.

### Message Queue

For async operations:
- Order confirmation emails
- Inventory updates
- Analytics events

```go
// Example: Publish order event
func (s *orderService) CreateOrder(...) (*model.Order, error) {
    order, err := s.orderRepo.Create(order)
    if err != nil {
        return nil, err
    }
    
    // Publish event
    s.eventPublisher.Publish("order.created", order)
    
    return order, nil
}
```

## 4. Security Best Practices

### Authentication

1. **Password Hashing**: Use bcrypt with cost factor 10+.
2. **JWT Expiry**: Short-lived access tokens (15-60 min).
3. **Refresh Tokens**: Implement refresh token rotation.

### Input Validation

```go
// Always validate input
type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required,max=50"`
}
```

### Rate Limiting

```go
// Apply rate limiting to prevent abuse
rateLimiter := middleware.NewRateLimiter(100, 200) // 100 req/min
router.Use(rateLimiter.RateLimit())
```

### SQL Injection Prevention

- Always use GORM's parameterized queries.
- Never concatenate user input into SQL strings.

## 5. API Design

### Response Format

```json
{
    "success": true,
    "data": { ... },
    "message": "Operation successful"
}
```

### Error Handling

```json
{
    "success": false,
    "error": "Invalid credentials"
}
```

### Pagination

```json
{
    "success": true,
    "data": { "items": [...] },
    "meta": {
        "current_page": 1,
        "per_page": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

## 6. Logging and Monitoring

### Structured Logging

```go
log.Printf("[INFO] User %d logged in from %s", userID, clientIP)
```

### Key Metrics to Monitor

- Request latency (p50, p95, p99)
- Error rates
- Database query performance
- Cache hit rates
- Active users

## 7. Deployment

### Environment Variables

```bash
# Production
APP_ENV=production
APP_PORT=8080
DB_HOST=prod-db.example.com
JWT_SECRET=<strong-random-secret>
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]
```

### Health Checks

```bash
curl http://localhost:8080/health
# Response: {"status": "healthy"}
```

## 8. Code Review Checklist

- [ ] Input validation implemented
- [ ] Error handling complete
- [ ] Database transactions used where needed
- [ ] No sensitive data in logs
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Security considerations addressed

## 9. Common Pitfalls

### N+1 Query Problem

```go
// Bad: Causes N+1 queries
products, _ := repo.FindAll()
for _, p := range products {
    shop, _ := shopRepo.FindByID(p.ShopID)
}

// Good: Use Preload
products, _ := repo.FindAllWithPreload()
```

### Memory Leaks

```go
// Bad: Goroutine leak
go func() {
    for {
        doSomething()
    }
}()

// Good: Use context for cancellation
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doSomething()
        }
    }
}(ctx)
```

### Transaction Management

```go
// Always handle transaction rollback
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&order).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```
