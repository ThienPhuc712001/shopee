# Security Implementation Guide

## Overview

This document describes the security implementation for the e-commerce platform, including authentication, token management, password protection, and API rate limiting.

---

## Table of Contents

1. [Security Overview](#1-security-overview)
2. [JWT Authentication](#2-jwt-authentication)
3. [Access Token & Refresh Token](#3-access-token--refresh-token)
4. [Database Tables](#4-database-tables)
5. [Password Security](#5-password-security)
6. [Rate Limiting](#6-rate-limiting)
7. [Rate Limit Middleware](#7-rate-limit-middleware)
8. [Token Validation](#8-token-validation)
9. [Refresh Token API](#9-refresh-token-api)
10. [Login Security](#10-login-security)
11. [Security Best Practices](#11-security-best-practices)
12. [Example Implementation](#12-example-implementation)

---

## 1. Security Overview

### Common Security Risks in E-commerce

| Risk | Description | Impact |
|------|-------------|--------|
| **Brute Force Attacks** | Automated login attempts to guess passwords | Account compromise |
| **Token Hijacking** | Stealing JWT tokens for unauthorized access | Session takeover |
| **API Abuse** | Excessive API calls causing service degradation | Service outage |
| **Password Leaks** | Exposed passwords in logs or database breaches | Mass account compromise |
| **SQL Injection** | Malicious SQL in queries | Data breach |
| **Cross-Site Scripting (XSS)** | Injected scripts in user input | Session theft |

### Why Security Measures Are Necessary

1. **Protect User Data**: Personal information, payment details
2. **Prevent Fraud**: Unauthorized purchases, account takeover
3. **Maintain Trust**: Security breaches damage reputation
4. **Compliance**: GDPR, PCI-DSS requirements
5. **Business Continuity**: Prevent service disruption from attacks

---

## 2. JWT Authentication

### How JWT Works

```
┌─────────┐      ┌─────────┐      ┌─────────┐
│  User   │      │ Server  │      │ Database│
└────┬────┘      └────┬────┘      └────┬────┘
     │                │                │
     │ 1. Login       │                │
     │───────────────>│                │
     │                │                │
     │                │ 2. Verify      │
     │                │───────────────>│
     │                │                │
     │                │ 3. Generate JWT│
     │                │                │
     │ 4. JWT Token   │                │
     │<───────────────│                │
     │                │                │
     │ 5. Request +   │                │
     │    JWT Token   │                │
     │───────────────>│                │
     │                │                │
     │                │ 6. Validate JWT│
     │                │                │
     │ 7. Response    │                │
     │<───────────────│                │
     └────────────────┴────────────────┘
```

### JWT Structure

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.    ← Header (algorithm)
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.    ← Payload (claims)
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c    ← Signature
```

### Example Authorization Header

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Advantages of JWT

| Advantage | Description |
|-----------|-------------|
| **Stateless** | No server-side session storage needed |
| **Scalable** | Works across multiple servers |
| **Self-contained** | Contains all user information |
| **Secure** | Cryptographically signed |
| **Flexible** | Custom claims for application needs |

---

## 3. Access Token & Refresh Token

### Token Strategy

| Token Type | Lifetime | Purpose | Storage |
|------------|----------|---------|---------|
| **Access Token** | 15-30 minutes | API requests | Memory/LocalStorage |
| **Refresh Token** | 7-30 days | Generate new access tokens | HttpOnly Cookie/DB |

### Refresh Flow

```
┌─────────┐      ┌─────────┐      ┌─────────┐
│  Client │      │ Server  │      │   DB    │
└────┬────┘      └────┬────┘      └────┬────┘
     │                │                │
     │ Access Token   │                │
     │ Expired        │                │
     │                │                │
     │ 1. POST /refresh
     │    + Refresh Token
     │───────────────>│                │
     │                │                │
     │                │ 2. Validate    │
     │                │ Refresh Token  │
     │                │───────────────>│
     │                │                │
     │                │ 3. Valid       │
     │                │<───────────────│
     │                │                │
     │                │ 4. Generate    │
     │                │ New Access Token
     │                │                │
     │ 5. New Access  │                │
     │    Token       │                │
     │<───────────────│                │
     │                │                │
     │ 6. API Request │                │
     │    + New Token │                │
     │───────────────>│                │
     └────────────────┴────────────────┘
```

---

## 4. Database Tables

### RefreshTokens Table

```sql
CREATE TABLE refresh_tokens (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token NVARCHAR(500) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    revoked BIT DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    revoked_at DATETIME,

    CONSTRAINT FK_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IX_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IX_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IX_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

### Relationships

- **refresh_tokens.user_id** → **users.id** (Many-to-One)
- One user can have multiple refresh tokens (multiple devices)
- Deleting a user cascades to delete all their refresh tokens

---

## 5. Password Security

### Password Hashing with bcrypt

```go
// Hash password
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(plainPassword), 
    bcrypt.DefaultCost, // 10
)

// Verify password
err := bcrypt.CompareHashAndPassword(
    []byte(hashedPassword), 
    []byte(plainPassword),
)
```

### Security Rules

| Rule | Description |
|------|-------------|
| **Never store plain text** | Always hash before storing |
| **Use bcrypt** | Designed for password hashing |
| **Minimum cost 10** | Higher = more secure but slower |
| **Minimum length 8** | Enforce strong passwords |
| **Check complexity** | Require letters, numbers, symbols |

### Verification Process

1. User submits password
2. Retrieve hashed password from database
3. Use `bcrypt.CompareHashAndPassword`
4. Return success/failure (never reveal why)

---

## 6. Rate Limiting

### What is Rate Limiting?

Rate limiting controls the number of requests a client can make in a time window.

### Example Configuration

```
Limit: 100 requests per minute per IP
```

### Benefits

| Benefit | Description |
|---------|-------------|
| **Prevent API Abuse** | Stop excessive requests |
| **Prevent DDoS** | Mitigate distributed attacks |
| **Protect Resources** | Ensure fair usage |
| **Security** | Slow down brute force attacks |

### Rate Limits by Endpoint

| Endpoint | Limit | Window |
|----------|-------|--------|
| `/api/auth/login` | 5 | 1 minute |
| `/api/auth/register` | 3 | 1 minute |
| `/api/*` (general) | 100 | 1 minute |
| `/api/admin/*` | 50 | 1 minute |

---

## 7. Rate Limit Middleware

### Algorithm: Sliding Window

```
┌─────────────────────────────────────────┐
│  Time Window: 1 minute                  │
│  Limit: 100 requests                    │
└─────────────────────────────────────────┘

Current Time: 10:01:30
Requests from IP 192.168.1.1:
- 10:00:35 ✓
- 10:00:40 ✓
- 10:00:45 ✓
...
- 10:01:25 ✓
- 10:01:30 ← Current request

Count requests in last 60 seconds
If count >= limit → Reject (429 Too Many Requests)
```

### Implementation

```go
type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int
    window   time.Duration
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    windowStart := now.Add(-rl.window)
    
    // Filter old requests
    var valid []time.Time
    for _, t := range rl.requests[ip] {
        if t.After(windowStart) {
            valid = append(valid, t)
        }
    }
    
    if len(valid) >= rl.limit {
        rl.requests[ip] = valid
        return false
    }
    
    rl.requests[ip] = append(valid, now)
    return true
}
```

---

## 8. Token Validation

### JWT Middleware Responsibilities

1. Extract token from Authorization header
2. Validate token format
3. Verify signature
4. Check expiration
5. Extract user_id from claims
6. Store user_id in context

### Example Middleware

```go
func JWTMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token
        authHeader := c.GetHeader("Authorization")
        token := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Parse and validate
        claims, err := jwt.Parse(token, secret)
        if err != nil {
            c.AbortWithStatusJSON(401, ErrorResponse{
                Message: "Invalid token",
            })
            return
        }
        
        // Check expiration
        if claims.ExpiresAt < time.Now().Unix() {
            c.AbortWithStatusJSON(401, ErrorResponse{
                Message: "Token expired",
            })
            return
        }
        
        // Store user_id in context
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

---

## 9. Refresh Token API

### Endpoint

```
POST /api/auth/refresh
```

### Request

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Response (Success)

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  },
  "message": "Token refreshed successfully"
}
```

### Response (Error)

```json
{
  "error": true,
  "message": "Invalid or expired refresh token",
  "code": "INVALID_TOKEN"
}
```

---

## 10. Login Security

### Protection Measures

| Measure | Implementation |
|---------|----------------|
| **Rate Limiting** | Max 5 login attempts per minute |
| **Account Lockout** | Lock after 5 failed attempts |
| **Cooldown Period** | 15 minute lock duration |
| **Logging** | Log all failed attempts |
| **Generic Errors** | "Invalid credentials" (not specific) |

### Failed Login Tracking

```go
type User struct {
    FailedLoginAttempts int
    LockedUntil         *time.Time
}

// After failed login
user.FailedLoginAttempts++
if user.FailedLoginAttempts >= 5 {
    lockTime := time.Now().Add(15 * time.Minute)
    user.LockedUntil = &lockTime
}

// After successful login
user.FailedLoginAttempts = 0
user.LockedUntil = nil
```

---

## 11. Security Best Practices

### HTTPS Only

- **Always use HTTPS** in production
- **Redirect HTTP to HTTPS**
- **HSTS header** to enforce HTTPS

### Secure Cookies

```go
c.SetCookie("refresh_token", token, 3600*24*7, "/", 
    "yourdomain.com",
    true,  // Secure (HTTPS only)
    true,  // HttpOnly (no JavaScript access)
    "Strict", // SameSite
)
```

### Strong Password Requirements

```
Minimum 8 characters
At least 1 uppercase letter
At least 1 lowercase letter
At least 1 number
At least 1 special character
```

### Additional Best Practices

| Practice | Description |
|----------|-------------|
| **Input Validation** | Validate all user input |
| **Output Encoding** | Prevent XSS |
| **Parameterized Queries** | Prevent SQL injection |
| **Security Headers** | CSP, X-Frame-Options, etc. |
| **Regular Updates** | Keep dependencies updated |
| **Audit Logging** | Log security events |
| **Secret Management** | Never commit secrets to code |

---

## 12. Example Implementation

### Project Structure

```
internal/
├── domain/
│   └── model/
│       └── refresh_token.go
├── repository/
│   └── refresh_token_repository.go
├── service/
│   └── auth_service.go
└── handler/
    └── auth_handler.go

pkg/
├── jwt/
│   └── jwt.go
├── password/
│   └── password.go
└── ratelimit/
    └── ratelimit.go

middleware/
├── jwt_auth.go
└── rate_limiter.go
```

### Quick Start

```go
// 1. Setup
log, _ := logger.New(logger.DefaultConfig())
db, _ := gorm.Open(sqlserver.Open(dsn))

// 2. Initialize services
jwtService := jwt.NewService(jwt.Config{
    Secret:         os.Getenv("JWT_SECRET"),
    AccessExpiry:   15 * time.Minute,
    RefreshExpiry:  7 * 24 * time.Hour,
})
authHandler := handler.NewAuthHandler(db, jwtService, log)

// 3. Setup routes
r := gin.New()
r.Use(middleware.RateLimitMiddleware(100, time.Minute))
r.POST("/api/auth/login", authHandler.Login)
r.POST("/api/auth/refresh", authHandler.Refresh)
r.POST("/api/auth/logout", 
    middleware.JWTAuth(jwtService), 
    authHandler.Logout)

// 4. Protected routes
api := r.Group("/api")
api.Use(middleware.JWTAuth(jwtService))
{
    api.GET("/users/me", userHandler.GetProfile)
}
```

---

## Quick Reference

### JWT Claims

```go
type Claims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_CREDENTIALS` | 401 | Wrong email/password |
| `ACCOUNT_LOCKED` | 403 | Too many failed attempts |
| `INVALID_TOKEN` | 401 | Malformed JWT |
| `TOKEN_EXPIRED` | 401 | JWT expired |
| `REFRESH_TOKEN_REVOKED` | 401 | Token was revoked |
| `TOO_MANY_REQUESTS` | 429 | Rate limit exceeded |

### Environment Variables

```bash
JWT_SECRET=your-super-secret-key-min-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h  # 7 days
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m
RATE_LIMIT=100
RATE_WINDOW=1m
```

---

## Support

For questions or issues, refer to:
- Package documentation in `pkg/jwt/`, `pkg/password/`
- Middleware documentation in `internal/middleware/`
- This guide: `docs/SECURITY_GUIDE.md`
