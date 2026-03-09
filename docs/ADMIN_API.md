# Admin Management System API Documentation

## Base URL
```
http://localhost:8080/api/admin
```

---

## Admin Authentication

### Admin Login
**POST** `/api/admin/auth/login`

**Request Body:**
```json
{
  "email": "admin@example.com",
  "password": "admin_password"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "admin": {
      "id": 1,
      "email": "admin@example.com",
      "first_name": "Admin",
      "last_name": "User",
      "role": {
        "id": 1,
        "name": "super_admin"
      }
    },
    "token": "admin_jwt_token_here"
  },
  "message": "Login successful"
}
```

---

## User Management

### Get All Users
**GET** `/api/admin/users`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| limit | int | 20 | Items per page |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "users": [...],
    "total": 1000
  }
}
```

### Ban User
**POST** `/api/admin/users/ban`

**Request Body:**
```json
{
  "user_id": 123,
  "reason": "Violation of terms of service"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "User banned successfully"
}
```

---

## Seller Management

### Get Pending Sellers
**GET** `/api/admin/sellers/pending`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "sellers": [
      {
        "id": 1,
        "name": "My Shop",
        "status": "pending",
        "user": {...}
      }
    ]
  }
}
```

### Approve Seller
**POST** `/api/admin/sellers/approve`

**Request Body:**
```json
{
  "shop_id": 1,
  "notes": "Documents verified"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Seller approved successfully"
}
```

---

## Product Management

### Get Products for Moderation
**GET** `/api/admin/products`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": [...]
  }
}
```

### Delete Product
**DELETE** `/api/admin/products/{id}`

**Request Body:**
```json
{
  "reason": "Prohibited item"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

---

## Order Management

### Get All Orders
**GET** `/api/admin/orders`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "orders": [...]
  }
}
```

### Get Order Details
**GET** `/api/admin/orders/{id}`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "order": {...}
  }
}
```

### Refund Order
**POST** `/api/admin/orders/refund`

**Request Body:**
```json
{
  "order_id": 1,
  "amount": 100.00,
  "reason": "Product defective"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "refund": {...}
  },
  "message": "Refund processed successfully"
}
```

---

## Analytics

### Get Admin Stats
**GET** `/api/admin/analytics/stats`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "stats": {
      "total_users": 10000,
      "total_sellers": 500,
      "total_products": 50000,
      "total_orders": 25000,
      "total_revenue": 5000000.00,
      "pending_orders": 100,
      "pending_refunds": 5,
      "active_users_24h": 2000,
      "new_users_today": 150
    }
  }
}
```

### Get Sales Analytics
**GET** `/api/admin/analytics/sales?start_date=2024-01-01&end_date=2024-01-31`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "sales": {
      "total_sales": 500000.00,
      "total_orders": 2500,
      "average_order_value": 200.00,
      "today_sales": 15000.00,
      "week_sales": 100000.00,
      "month_sales": 500000.00
    }
  }
}
```

### Get User Analytics
**GET** `/api/admin/analytics/users`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "users": {
      "total_users": 10000,
      "new_users_today": 150,
      "new_users_week": 1000,
      "new_users_month": 5000,
      "active_users": 3000,
      "banned_users": 50
    }
  }
}
```

### Get Product Analytics
**GET** `/api/admin/analytics/products?limit=10`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "products": {
      "total_products": 50000,
      "active_products": 48000,
      "out_of_stock": 500,
      "top_products": [
        {
          "id": 1,
          "name": "Popular Product",
          "sold_count": 5000,
          "revenue": 500000.00
        }
      ],
      "low_stock_products": [...]
    }
  }
}
```

---

## Audit Logs

### Get Audit Logs
**GET** `/api/admin/audit-logs`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Page number |
| limit | int | Items per page |
| admin_id | int | Filter by admin |
| action | string | Filter by action |
| entity_type | string | Filter by entity type |
| start_date | string | Start date |
| end_date | string | End date |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 1,
        "admin": {
          "id": 1,
          "email": "admin@example.com"
        },
        "action": "ban",
        "entity_type": "user",
        "entity_id": 123,
        "new_values": "{\"status\": \"banned\", \"reason\": \"TOS violation\"}",
        "ip_address": "192.168.1.1",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 5000
  }
}
```

---

## System Settings

### Get System Setting
**GET** `/api/admin/settings/{key}`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "setting": {
      "key": "site_name",
      "value": "E-Commerce Platform",
      "type": "string",
      "description": "Platform name"
    }
  }
}
```

### Update System Setting
**PUT** `/api/admin/settings/{key}`

**Request Body:**
```json
{
  "value": "New Platform Name"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Setting updated successfully"
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## cURL Examples

### Admin Login
```bash
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

### Ban User
```bash
curl -X POST http://localhost:8080/api/admin/users/ban \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "reason": "Violation of terms"
  }'
```

### Approve Seller
```bash
curl -X POST http://localhost:8080/api/admin/sellers/approve \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shop_id": 1,
    "notes": "Verified"
  }'
```

### Get Analytics
```bash
curl -X GET "http://localhost:8080/api/admin/analytics/stats" \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

This Admin API is production-ready with full RBAC, audit logging, and platform management capabilities.
