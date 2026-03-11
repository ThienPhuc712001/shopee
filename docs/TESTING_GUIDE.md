# Testing Guide for E-Commerce Platform

## Overview

This document describes the testing strategy and implementation for the e-commerce platform. The testing system ensures backend logic works correctly and prevents regressions.

---

## Table of Contents

1. [Why Testing is Important](#1-why-testing-is-important)
2. [Types of Tests](#2-types-of-tests)
3. [Test Directory Structure](#3-test-directory-structure)
4. [Test Naming Convention](#4-test-naming-convention)
5. [Mocking Dependencies](#5-mocking-dependencies)
6. [Service Test Example](#6-service-test-example)
7. [API Handler Test](#7-api-handler-test)
8. [Auth Testing](#8-auth-testing)
9. [Order Testing](#9-order-testing)
10. [Test Data Setup](#10-test-data-setup)
11. [Running Tests](#11-running-tests)
12. [Example Implementation](#12-example-implementation)

---

## 1. Why Testing is Important

Testing is critical for building reliable backend systems. Here's why:

### Benefits of Testing

| Benefit | Description |
|---------|-------------|
| **Prevent Regression Bugs** | Catch breaking changes before they reach production |
| **Validate Business Logic** | Ensure calculations and rules work correctly |
| **Ensure API Reliability** | Verify endpoints return expected responses |
| **Document Behavior** | Tests serve as living documentation |
| **Enable Refactoring** | Confidence to improve code without breaking features |
| **Reduce Debugging Time** | Catch issues early in development |

### Consequences of Missing Tests

1. **Silent Bugs**: Logic errors go undetected until production
2. **Regression**: New features break existing functionality
3. **Manual Testing Overhead**: Every change requires manual verification
4. **Production Incidents**: Critical bugs reach customers
5. **Technical Debt**: Fear of changing code leads to stagnation

### Example: Cost of Bug by Stage

| Stage | Cost to Fix |
|-------|-------------|
| Development (with tests) | 1x |
| QA/Testing | 5x |
| Production | 100x |

---

## 2. Types of Tests

### Testing Pyramid

```
         /\
        /  \
       / E2E \      ← End-to-End Tests (few)
      /______\
     /        \
    / Integration\  ← Integration Tests (some)
   /______________\
  /                \
 /     Unit Tests    \ ← Unit Tests (many)
/______________________\
```

### Test Types

| Type | Scope | Speed | Coverage |
|------|-------|-------|----------|
| **Unit Tests** | Single function/method | Fast (<10ms) | High |
| **Service Tests** | Business logic layer | Fast (<50ms) | High |
| **Handler Tests** | HTTP endpoints | Medium (<100ms) | Medium |
| **Integration Tests** | Multiple components | Slow (>500ms) | Medium |
| **E2E Tests** | Full system | Very Slow (>1s) | Low |

### Unit Tests

Test individual functions in isolation.

```go
func TestCalculateTotal(t *testing.T) {
    items := []CartItem{
        {Price: 100, Quantity: 2},
        {Price: 50, Quantity: 1},
    }
    
    total := CalculateTotal(items)
    
    if total != 250 {
        t.Errorf("Expected 250, got %d", total)
    }
}
```

### Service Tests

Test business logic layer with mocked dependencies.

```go
func TestCreateProduct(t *testing.T) {
    mockRepo := &MockProductRepository{}
    service := NewProductService(mockRepo)
    
    product, err := service.CreateProduct(&Product{
        Name: "Test Product",
        Price: 99.99,
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, product.ID)
}
```

### Handler Tests

Test HTTP endpoints with simulated requests.

```go
func TestGetProducts(t *testing.T) {
    router := setupRouter()
    req, _ := http.NewRequest("GET", "/api/products", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

---

## 3. Test Directory Structure

### Recommended Structure

```
ecommerce/
├── internal/
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── auth_service_test.go
│   │   ├── product_service.go
│   │   └── product_service_test.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── auth_handler_test.go
│   │   ├── product_handler.go
│   │   └── product_handler_test.go
│   └── repository/
│       ├── product_repository.go
│       └── product_repository_test.go
├── pkg/
│   ├── jwt/
│   │   ├── jwt.go
│   │   └── jwt_test.go
│   └── password/
│       ├── password.go
│       └── password_test.go
├── tests/
│   ├── integration/
│   │   ├── auth_integration_test.go
│   │   └── order_integration_test.go
│   └── testdata/
│       ├── users.json
│       └── products.json
└── testutil/
    ├── testutil.go
    └── mocks.go
```

### Organization Principles

1. **Co-location**: Test files next to source files
2. **Integration tests**: Separate directory for cross-component tests
3. **Test data**: Centralized test data in `testdata/`
4. **Utilities**: Shared test helpers in `testutil/`

---

## 4. Test Naming Convention

### File Naming

```
{component}_{layer}_test.go
```

Examples:
- `auth_service_test.go`
- `product_handler_test.go`
- `order_repository_test.go`
- `jwt_test.go`

### Test Function Naming

```
Test{Component}_{Action}_{ExpectedResult}
```

Examples:
- `TestAuthService_Login_Success`
- `TestAuthService_Login_InvalidCredentials`
- `TestProductService_CreateProduct_ValidData`
- `TestProductService_GetProduct_NotFound`

### Benefits of Clear Naming

| Benefit | Description |
|---------|-------------|
| **Discoverability** | Easy to find relevant tests |
| **Readability** | Understand test purpose at a glance |
| **Debugging** | Quickly locate failing tests |
| **Documentation** | Test names describe expected behavior |

---

## 5. Mocking Dependencies

### Why Mocking is Necessary

1. **Isolation**: Test one component at a time
2. **Speed**: Avoid slow database/network calls
3. **Reliability**: No dependency on external services
4. **Control**: Simulate edge cases and errors

### Mock Database Repository

```go
type MockProductRepository struct {
    Products map[uint]*model.Product
    NextID   uint
}

func (m *MockProductRepository) Create(ctx context.Context, product *model.Product) error {
    product.ID = m.NextID
    m.NextID++
    m.Products[product.ID] = product
    return nil
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
    product, exists := m.Products[id]
    if !exists {
        return nil, gorm.ErrRecordNotFound
    }
    return product, nil
}
```

### Mock External Services

```go
type MockEmailService struct {
    SentEmails []Email
}

func (m *MockEmailService) Send(to, subject, body string) error {
    m.SentEmails = append(m.SentEmails, Email{
        To:      to,
        Subject: subject,
        Body:    body,
    })
    return nil
}
```

---

## 6. Service Test Example

### ProductService Tests

```go
package service

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "ecommerce/internal/repository"
    "ecommerce/internal/domain/model"
)

func TestProductService_CreateProduct(t *testing.T) {
    // Arrange
    mockRepo := &repository.MockProductRepository{
        Products: make(map[uint]*model.Product),
        NextID:   1,
    }
    service := NewProductService(mockRepo)
    
    input := &CreateProductRequest{
        Name:  "Test Product",
        Price: 99.99,
        Stock: 100,
    }
    
    // Act
    product, err := service.CreateProduct(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, product)
    assert.Equal(t, "Test Product", product.Name)
    assert.Equal(t, 99.99, product.Price)
}

func TestProductService_GetProduct(t *testing.T) {
    // Arrange
    mockRepo := &repository.MockProductRepository{
        Products: map[uint]*model.Product{
            1: {ID: 1, Name: "Product 1", Price: 50.00},
        },
    }
    service := NewProductService(mockRepo)
    
    // Act
    product, err := service.GetProduct(context.Background(), 1)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "Product 1", product.Name)
}

func TestProductService_UpdateProduct(t *testing.T) {
    // Arrange
    mockRepo := &repository.MockProductRepository{
        Products: map[uint]*model.Product{
            1: {ID: 1, Name: "Original", Price: 50.00},
        },
    }
    service := NewProductService(mockRepo)
    
    input := &UpdateProductRequest{
        Name:  "Updated",
        Price: 75.00,
    }
    
    // Act
    err := service.UpdateProduct(context.Background(), 1, input)
    
    // Assert
    assert.NoError(t, err)
    product, _ := service.GetProduct(context.Background(), 1)
    assert.Equal(t, "Updated", product.Name)
    assert.Equal(t, 75.00, product.Price)
}
```

---

## 7. API Handler Test

### Testing HTTP Handlers

```go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestProductHandler_GetProducts(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    handler := NewProductHandler(mockService)
    
    router.GET("/api/products", handler.ListProducts)
    
    // Request
    req, _ := http.NewRequest("GET", "/api/products", nil)
    w := httptest.NewRecorder()
    
    // Execute
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "success")
}

func TestProductHandler_CreateProduct(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    handler := NewProductHandler(mockService)
    
    router.POST("/api/products", handler.CreateProduct)
    
    // Request body
    body := map[string]interface{}{
        "name":  "New Product",
        "price": 99.99,
        "stock": 50,
    }
    jsonBody, _ := json.Marshal(body)
    
    req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // Execute
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

---

## 8. Auth Testing

### Authentication Test Cases

```go
func TestAuthService_Login_Success(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{
        UsersByEmail: map[string]*model.User{
            "test@example.com": {
                ID:       1,
                Email:    "test@example.com",
                Password: "$2a$10$hashedpassword",
                Role:     model.RoleCustomer,
            },
        },
    }
    service := NewAuthService(mockRepo, mockTokenService)
    
    req := &LoginRequest{
        Email:    "test@example.com",
        Password: "correctpassword",
    }
    
    // Act
    response, err := service.Login(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, response.AccessToken)
    assert.NotNil(t, response.RefreshToken)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{
        UsersByEmail: map[string]*model.User{
            "test@example.com": {
                ID:       1,
                Email:    "test@example.com",
                Password: "$2a$10$hashedpassword",
            },
        },
    }
    service := NewAuthService(mockRepo, mockTokenService)
    
    req := &LoginRequest{
        Email:    "test@example.com",
        Password: "wrongpassword",
    }
    
    // Act
    response, err := service.Login(context.Background(), req)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, response)
    assert.Equal(t, ErrInvalidCredentials, err)
}

func TestAuthService_Login_AccountLocked(t *testing.T) {
    // Arrange
    lockedUntil := time.Now().Add(10 * time.Minute)
    mockRepo := &MockUserRepository{
        UsersByEmail: map[string]*model.User{
            "test@example.com": {
                ID:            1,
                Email:         "test@example.com",
                Password:      "$2a$10$hashedpassword",
                LockedUntil:   &lockedUntil,
                FailedLoginAttempts: 5,
            },
        },
    }
    service := NewAuthService(mockRepo, mockTokenService)
    
    req := &LoginRequest{
        Email:    "test@example.com",
        Password: "correctpassword",
    }
    
    // Act
    response, err := service.Login(context.Background(), req)
    
    // Assert
    assert.Error(t, err)
    assert.Equal(t, ErrAccountLocked, err)
}
```

---

## 9. Order Testing

### Order Workflow Tests

```go
func TestOrderService_CreateOrder(t *testing.T) {
    // Arrange
    mockRepo := &MockOrderRepository{}
    mockCartRepo := &MockCartRepository{}
    mockProductRepo := &MockProductRepository{}
    
    service := NewOrderService(mockRepo, mockCartRepo, mockProductRepo)
    
    req := &CreateOrderRequest{
        UserID:     1,
        Items: []OrderItemRequest{
            {ProductID: 1, Quantity: 2, Price: 50.00},
            {ProductID: 2, Quantity: 1, Price: 100.00},
        },
    }
    
    // Act
    order, err := service.CreateOrder(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, order.ID)
    assert.Equal(t, 200.00, order.TotalAmount)
}

func TestOrderService_CalculateTotal(t *testing.T) {
    // Arrange
    items := []OrderItem{
        {Price: 50.00, Quantity: 2},
        {Price: 100.00, Quantity: 1},
    }
    
    // Act
    total := calculateTotal(items)
    
    // Assert
    assert.Equal(t, 200.00, total)
}

func TestOrderService_ApplyCoupon(t *testing.T) {
    // Arrange
    mockCouponRepo := &MockCouponRepository{
        CouponsByCode: map[string]*model.Coupon{
            "SAVE10": {
                Code:       "SAVE10",
                DiscountType: "percentage",
                DiscountValue: 10,
                MinOrderAmount: 50,
            },
        },
    }
    service := NewOrderService(mockCouponRepo)
    
    // Act
    discount, err := service.ApplyCoupon(context.Background(), "SAVE10", 100.00)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 10.00, discount)
}
```

---

## 10. Test Data Setup

### Creating Test Data

```go
package testutil

import (
    "ecommerce/internal/domain/model"
    "time"
)

// CreateTestUser creates a test user
func CreateTestUser(id uint, email string) *model.User {
    return &model.User{
        ID:        id,
        Email:     email,
        Password:  "$2a$10$hashedpassword",
        FirstName: "Test",
        LastName:  "User",
        Role:      model.RoleCustomer,
        Status:    model.StatusActive,
        CreatedAt: time.Now(),
    }
}

// CreateTestProduct creates a test product
func CreateTestProduct(id uint, shopID uint, price float64) *model.Product {
    return &model.Product{
        ID:          id,
        ShopID:      shopID,
        CategoryID:  1,
        Name:        "Test Product",
        Slug:        "test-product",
        Price:       price,
        Stock:       100,
        Status:      "active",
        CreatedAt:   time.Now(),
    }
}

// CreateTestOrder creates a test order
func CreateTestOrder(id uint, userID uint, total float64) *model.Order {
    return &model.Order{
        ID:           id,
        UserID:       userID,
        OrderNumber:  generateOrderNumber(),
        Status:       "pending",
        TotalAmount:  total,
        CreatedAt:    time.Now(),
    }
}

func generateOrderNumber() string {
    return "ORD-" + time.Now().Format("20060102150405")
}
```

---

## 11. Running Tests

### Commands

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/service/

# Run specific test
go test -v ./internal/service/ -run TestAuthService_Login

# Run with coverage
go test -cover ./...

# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detector
go test -race ./...
```

### Understanding Output

```
=== RUN   TestAuthService_Login_Success
--- PASS: TestAuthService_Login_Success (0.00s)
=== RUN   TestAuthService_Login_InvalidPassword
--- PASS: TestAuthService_Login_InvalidPassword (0.00s)
=== RUN   TestAuthService_Login_AccountLocked
--- PASS: TestAuthService_Login_AccountLocked (0.00s)
PASS
ok      ecommerce/internal/service    0.015s
```

| Symbol | Meaning |
|--------|---------|
| `=== RUN` | Test started |
| `--- PASS` | Test passed |
| `--- FAIL` | Test failed |
| `--- SKIP` | Test skipped |
| `ok` | Package passed |
| `FAIL` | Package failed |

---

## 12. Example Implementation

### Test File Template

```go
package service

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "ecommerce/internal/testutil"
)

func TestService_Method_Success(t *testing.T) {
    // Arrange - Setup test data and mocks
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)
    
    input := &InputRequest{
        Field: "value",
    }
    
    // Act - Call the method being tested
    result, err := service.Method(context.Background(), input)
    
    // Assert - Verify the result
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "expected", result.Field)
}

func TestService_Method_Error(t *testing.T) {
    // Arrange
    mockRepo := &MockRepository{
        ErrorToReturn: ErrSomeError,
    }
    service := NewService(mockRepo)
    
    input := &InputRequest{
        Field: "invalid",
    }
    
    // Act
    result, err := service.Method(context.Background(), input)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Equal(t, ErrSomeError, err)
}
```

### Best Practices Checklist

- [ ] Test happy path (success scenarios)
- [ ] Test error cases
- [ ] Test edge cases (empty input, max values)
- [ ] Use descriptive test names
- [ ] Keep tests independent (no shared state)
- [ ] Clean up after tests (defer cleanup)
- [ ] Mock external dependencies
- [ ] Test both valid and invalid input
- [ ] Aim for >80% code coverage
- [ ] Run tests before every commit

---

## Quick Reference

### Test Structure

```go
func TestComponent_Method_Condition(t *testing.T) {
    // Arrange
    // Set up test data and mocks
    
    // Act
    // Call the function/method being tested
    
    // Assert
    // Verify the results
}
```

### Common Assertions

```go
assert.NoError(t, err)
assert.Error(t, err)
assert.Nil(t, value)
assert.NotNil(t, value)
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)
assert.True(t, condition)
assert.False(t, condition)
assert.Empty(t, slice)
assert.NotEmpty(t, slice)
assert.Len(t, slice, 5)
assert.Contains(t, str, substr)
```

---

## Support

For questions about testing:
- Review existing test files for patterns
- Check `internal/testutil/` for helpers
- Refer to this guide: `docs/TESTING_GUIDE.md`
