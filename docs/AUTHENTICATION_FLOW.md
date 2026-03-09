# Authentication & Authorization System

## Overview

This document explains the complete authentication flow for the e-commerce platform.

---

## PART 1 — Authentication Flow

### Step-by-Step Authentication Process

```
┌─────────────────────────────────────────────────────────────────┐
│                    AUTHENTICATION FLOW                          │
└─────────────────────────────────────────────────────────────────┘

1. USER REGISTRATION
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Client  │────────>│  Server  │────────>│ Database │
   │          │  POST   │          │  INSERT │          │
   │          │ /register│          │  user   │          │
   └──────────┘         └──────────┘         └──────────┘

2. PASSWORD HASHING (Server-side)
   Plain Password: "mypassword123"
        │
        ▼
   bcrypt.GenerateFromPassword(cost=10)
        │
        ▼
   Hashed: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
        │
        ▼
   Store in database (NEVER store plain text!)

3. LOGIN
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │  Client  │────────>│  Server  │────────>│ Database │
   │          │  POST   │          │  SELECT │          │
   │          │  /login │          │  user   │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   4. VERIFY PASSWORD
      bcrypt.CompareHashAndPassword(hash, input)
                              │
                              ▼
   5. GENERATE JWT TOKENS
      - Access Token (15 min)
      - Refresh Token (7 days)
                              │
                              ▼
   6. RETURN TO CLIENT
      {
        "access_token": "eyJhbGc...",
        "refresh_token": "eyJhbGc...",
        "expires_in": 900
      }

7. PROTECTED API REQUEST
   ┌──────────┐         ┌──────────┐
   │  Client  │────────>│  Server  │
   │          │  GET    │          │
   │          │ /orders │          │
   │          │ +Header │          │
   └──────────┘         └──────────┘
   Authorization: Bearer eyJhbGc...
                              │
                              ▼
   8. JWT MIDDLEWARE
      - Extract token from header
      - Verify signature
      - Check expiration
      - Extract user claims
      - Inject into context
                              │
                              ▼
   9. RETURN PROTECTED DATA
      { "orders": [...] }

10. TOKEN REFRESH (when access token expires)
    ┌──────────┐         ┌──────────┐
    │  Client  │────────>│  Server  │
    │          │  POST   │          │
    │          │ /refresh│          │
    │          │ +refresh│          │
    │          │  token  │          │
    └──────────┘         └──────────┘
                              │
                              ▼
    Validate refresh token → Generate new access token
```

### Token Lifecycle

```
Access Token (15 minutes)          Refresh Token (7 days)
        │                                  │
        ▼                                  ▼
   ┌─────────┐                       ┌─────────┐
   │ Short   │                       │  Long   │
   │ Lived   │                       │  Lived  │
   └─────────┘                       └─────────┘
        │                                  │
        ▼                                  ▼
   Used for API                       Used to get
   Authentication                     new Access Token
```

---

## PART 2 — Security Best Practices

### Password Requirements
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character

### Token Security
- Access tokens stored in memory (client-side)
- Refresh tokens stored in httpOnly cookies
- Tokens signed with HS256 algorithm
- Secret key minimum 32 characters

### Rate Limiting
- Login: 5 attempts per minute per IP
- Register: 3 attempts per minute per IP
- Refresh: 10 attempts per minute per IP

---

## PART 3 — Common Attack Prevention

| Attack Type | Prevention |
|-------------|------------|
| SQL Injection | Parameterized queries (GORM) |
| XSS | httpOnly cookies, input sanitization |
| CSRF | CSRF tokens, SameSite cookies |
| Brute Force | Rate limiting, account lockout |
| Token Theft | Short expiry, refresh token rotation |
| Password Attacks | bcrypt hashing, salt |
