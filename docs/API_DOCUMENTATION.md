# API Documentation

## Base URL
```
http://localhost:8080/api
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your_access_token>
```

---

## Authentication Endpoints

### 1. Register User
**POST** `/api/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "phone": "0123456789",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Password Requirements:**
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "phone": "0123456789",
      "role": "customer",
      "email_verified": false,
      "created_at": "2024-01-15T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900,
      "token_type": "Bearer",
      "expiry": "2024-01-15T10:45:00Z"
    }
  },
  "message": "Registration successful"
}
```

---

### 2. Login
**POST** `/api/auth/login`

Authenticate user and receive tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "avatar": "https://example.com/avatar.jpg",
      "role": "customer"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900,
      "token_type": "Bearer"
    }
  },
  "message": "Login successful"
}
```

**Error Responses:**
```json
// Invalid credentials (401)
{
  "success": false,
  "error": "Invalid email or password"
}

// Account locked (423)
{
  "success": false,
  "error": "Account is locked. Please try again later or contact support."
}

// Account inactive (403)
{
  "success": false,
  "error": "Account is inactive. Please contact support."
}
```

---

### 3. Refresh Token
**POST** `/api/auth/refresh`

Generate new access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900,
      "token_type": "Bearer"
    }
  },
  "message": "Token refreshed successfully"
}
```

---

### 4. Get Current User
**GET** `/api/auth/me`

Get current authenticated user information.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "phone": "0123456789",
      "avatar": "https://example.com/avatar.jpg",
      "role": "customer",
      "email_verified": false,
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

---

### 5. Update Profile
**PUT** `/api/auth/profile`

Update current user's profile information.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "phone": "0987654321"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "phone": "0987654321"
    }
  },
  "message": "Profile updated successfully"
}
```

---

### 6. Change Password
**POST** `/api/auth/change-password`

Change current user's password.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "old_password": "OldPass123!",
  "new_password": "NewSecurePass456!"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

---

### 7. Logout
**POST** `/api/auth/logout`

Invalidate refresh token and logout.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

---

### 8. Forgot Password
**POST** `/api/auth/forgot-password`

Request password reset email.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "If the email exists, a password reset link has been sent"
}
```

---

### 9. Reset Password
**POST** `/api/auth/reset-password`

Reset password with token from email.

**Request Body:**
```json
{
  "token": "reset-token-from-email",
  "new_password": "NewSecurePass123!"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

---

### 10. Verify Email
**GET** `/api/auth/verify-email?token=<token>`

Verify user email with token.

**Query Parameters:**
- `token` (required): Verification token from email

**Success Response (200):**
```json
{
  "success": true,
  "message": "Email verified successfully"
}
```

---

### 11. Resend Verification Email
**POST** `/api/auth/resend-verification`

Resend email verification link.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Verification email sent"
}
```

---

## Product Endpoints

### Get All Products
**GET** `/api/products`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": 1,
        "name": "Product Name",
        "slug": "product-name",
        "price": 99.99,
        "original_price": 149.99,
        "discount": 33,
        "stock": 100,
        "sold": 50,
        "rating": 4.5,
        "total_reviews": 20,
        "shop": {
          "id": 1,
          "name": "Shop Name",
          "rating": 4.8
        },
        "images": [
          {
            "url": "https://example.com/image.jpg",
            "is_primary": true
          }
        ]
      }
    ]
  },
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

### Get Product by ID
**GET** `/api/products/:id`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "product": {
      "id": 1,
      "name": "Product Name",
      "description": "Product description...",
      "price": 99.99,
      "stock": 100,
      "shop": {...},
      "category": {...},
      "images": [...]
    }
  }
}
```

---

### Search Products
**GET** `/api/products/search`

**Query Parameters:**
- `keyword`: Search term
- `category_id`: Filter by category
- `shop_id`: Filter by shop
- `min_price`: Minimum price
- `max_price`: Maximum price
- `page`: Page number
- `limit`: Items per page

---

### Create Product (Seller Only)
**POST** `/api/products`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 99.99,
  "original_price": 149.99,
  "stock": 100,
  "category_id": 1,
  "shop_id": 1
}
```

---

## Cart Endpoints

### Get Cart
**GET** `/api/cart`

**Headers:**
```
Authorization: Bearer <access_token>
```

---

### Add to Cart
**POST** `/api/cart/add`

**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2
}
```

---

### Update Cart Item
**PUT** `/api/cart/items/:id`

**Request Body:**
```json
{
  "quantity": 5
}
```

---

### Remove from Cart
**DELETE** `/api/cart/items/:id`

---

### Clear Cart
**DELETE** `/api/cart/clear`

---

## Order Endpoints

### Create Order
**POST** `/api/orders`

**Request Body:**
```json
{
  "shipping_info": {
    "name": "John Doe",
    "phone": "0123456789",
    "address": "123 Main St",
    "ward": "Ward 1",
    "district": "District 1",
    "city": "Ho Chi Minh City",
    "country": "Vietnam"
  },
  "note": "Please deliver before 5 PM"
}
```

---

### Get User Orders
**GET** `/api/orders`

**Query Parameters:**
- `page`: Page number
- `limit`: Items per page

---

### Get Order by ID
**GET** `/api/orders/:id`

---

### Cancel Order
**POST** `/api/orders/:id/cancel`

**Request Body:**
```json
{
  "reason": "Changed my mind"
}
```

---

### Update Order Status (Seller/Admin Only)
**PUT** `/api/orders/:id/status`

**Request Body:**
```json
{
  "status": "confirmed"
}
```

**Valid Statuses:**
- `pending`
- `confirmed`
- `processing`
- `shipped`
- `delivered`
- `cancelled`
- `refunded`

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 204 | No Content - Successful deletion |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation error |
| 423 | Locked - Account locked |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

---

## Rate Limiting

- Default: 100 requests per minute per IP
- Authenticated: Based on user ID
- Exceeded: Returns 429 status code

```json
{
  "success": false,
  "error": "Too many requests, please try again later"
}
```

---

## Error Handling

All errors follow this format:
```json
{
  "success": false,
  "error": "Error message here"
}
```

Common error messages:
- "Authorization header required"
- "Invalid authorization format"
- "Token has expired"
- "Invalid token"
- "Insufficient permissions"
- "Resource not found"
- "Validation failed"
