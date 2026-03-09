# Authentication Quick Reference

## Token Configuration

| Token Type | Expiry | Secret | Storage |
|------------|--------|--------|---------|
| Access Token | 15 minutes | JWT_SECRET | Client memory |
| Refresh Token | 7 days | JWT_SECRET + suffix | httpOnly cookie |

## Environment Variables

```env
# Required
JWT_SECRET=your-super-secret-key-min-32-characters

# Optional
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h  # 7 days
BCRYPT_COST=10
```

## Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /api/auth/register | No | Register new user |
| POST | /api/auth/login | No | Login user |
| POST | /api/auth/refresh | No | Refresh access token |
| GET | /api/auth/me | Yes | Get current user |
| PUT | /api/auth/profile | Yes | Update profile |
| POST | /api/auth/change-password | Yes | Change password |
| POST | /api/auth/logout | Yes | Logout user |
| POST | /api/auth/forgot-password | No | Request reset |
| POST | /api/auth/reset-password | No | Reset password |
| GET | /api/auth/verify-email | No | Verify email |
| POST | /api/auth/resend-verification | No | Resend verification |

## Role-Based Access

| Role | Can Access |
|------|------------|
| customer | Public endpoints + own cart, orders, profile |
| seller | customer + create/manage own products, view shop orders |
| admin | all endpoints + user management, all orders |

## Middleware Usage

```go
// Require authentication
router.Use(middleware.JWTAuth(tokenService))

// Require specific role
router.Use(middleware.RequireRole(model.RoleSeller))

// Require admin
router.Use(middleware.RequireAdmin())

// Require seller or admin
router.Use(middleware.RequireSellerOrAdmin())

// Optional auth (continues even if not authenticated)
router.Use(middleware.OptionalJWTAuth(tokenService))
```

## Password Requirements

- Minimum 8 characters
- At least 1 uppercase letter (A-Z)
- At least 1 lowercase letter (a-z)
- At least 1 number (0-9)
- At least 1 special character (!@#$%^&*)

## Account Lockout

- 5 failed login attempts → Account locked for 30 minutes
- Lock status checked on every login attempt
- Successful login resets failed attempt counter

## Token Flow

```
┌─────────────────────────────────────────────────────────┐
│ 1. Client registers/logins                              │
│    → Receives access_token + refresh_token              │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 2. Client makes API requests                            │
│    → Includes: Authorization: Bearer <access_token>     │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 3. Access token expires (15 min)                        │
│    → Server returns 401 Unauthorized                    │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 4. Client uses refresh_token                            │
│    → POST /api/auth/refresh                             │
│    → Receives new access_token + refresh_token          │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│ 5. Refresh token expires (7 days)                       │
│    → Client must login again                            │
└─────────────────────────────────────────────────────────┘
```

## Example cURL Commands

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "phone": "0123456789",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

### Get Current User
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer <your_access_token>"
```

### Refresh Token
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<your_refresh_token>"
  }'
```

### Change Password
```bash
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer <your_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "OldPass123!",
    "new_password": "NewSecurePass456!"
  }'
```

## Security Checklist

- [ ] JWT_SECRET is at least 32 characters
- [ ] JWT_SECRET is different in production
- [ ] HTTPS enabled in production
- [ ] Passwords hashed with bcrypt (cost >= 10)
- [ ] Rate limiting enabled
- [ ] CORS configured properly
- [ ] Account lockout after failed attempts
- [ ] Token expiry set appropriately
- [ ] Refresh tokens rotated on use
- [ ] Logout invalidates refresh tokens

## Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| "Token has expired" | Use refresh token to get new access token |
| "Invalid token" | Ensure token is not modified, check secret key |
| "Authorization header required" | Include `Authorization: Bearer <token>` |
| "Account is locked" | Wait 30 minutes or contact support |
| "Invalid email or password" | Check credentials, verify caps lock |
| "Email already registered" | Use different email or login |

## Testing

```bash
# Health check
curl http://localhost:8080/health

# Register and capture token
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecurePass123!"}' \
  | jq -r '.data.tokens.access_token')

# Use token
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"
```
