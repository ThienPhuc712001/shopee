# Order & Checkout API Documentation

## Base URL
```
http://localhost:8080/api
```

---

## Order Endpoints

### 1. Checkout (Create Order from Cart)
**POST** `/api/orders/checkout`

Convert cart items to order and proceed to payment.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "shipping_info": {
    "name": "John Doe",
    "phone": "0123456789",
    "address": "123 Main Street",
    "ward": "Ward 1",
    "district": "District 1",
    "city": "Ho Chi Minh City",
    "state": "HCMC",
    "country": "Vietnam",
    "postal_code": "700000"
  },
  "payment_method": "cod",
  "buyer_note": "Please deliver before 5 PM",
  "voucher_code": "SAVE10",
  "shipping_method": "standard"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "order": {
      "id": 1,
      "order_number": "ORD-20240115-00001",
      "user_id": 123,
      "shop_id": 1,
      "status": "pending",
      "payment_status": "pending",
      "subtotal": 599.98,
      "shipping_fee": 30000,
      "discount": 0,
      "total_amount": 629.98,
      "shipping_name": "John Doe",
      "shipping_phone": "0123456789",
      "shipping_address": "123 Main Street, Ward 1, District 1",
      "shipping_city": "Ho Chi Minh City",
      "shipping_country": "Vietnam",
      "items": [
        {
          "id": 1,
          "product_id": 101,
          "product_name": "Wireless Headphones",
          "quantity": 2,
          "price": 299.99,
          "subtotal": 599.98
        }
      ],
      "created_at": "2024-01-15T10:30:00Z"
    }
  },
  "message": "Order created successfully"
}
```

**Error Responses:**

```json
// Cart empty (400)
{
  "success": false,
  "error": "Cart is empty"
}

// Insufficient stock (400)
{
  "success": false,
  "error": "Insufficient stock available"
}

// Product unavailable (400)
{
  "success": false,
  "error": "Product is no longer available"
}

// Invalid shipping info (400)
{
  "success": false,
  "error": "Invalid shipping information"
}
```

---

### 2. Get User Orders
**GET** `/api/orders`

Get all orders for the current user.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| limit | int | 10 | Items per page (max: 100) |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "orders": [
      {
        "id": 1,
        "order_number": "ORD-20240115-00001",
        "status": "pending",
        "payment_status": "pending",
        "total_amount": 629.98,
        "total_items": 2,
        "created_at": "2024-01-15T10:30:00Z",
        "shop": {
          "id": 1,
          "name": "AudioTech Store"
        },
        "items": [
          {
            "id": 1,
            "product_name": "Wireless Headphones",
            "quantity": 2,
            "price": 299.99,
            "product_image": "/uploads/products/headphones.jpg"
          }
        ]
      }
    ]
  },
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

---

### 3. Get Order by ID
**GET** `/api/orders/:id`

Get detailed information about a specific order.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "order": {
      "id": 1,
      "order_number": "ORD-20240115-00001",
      "status": "processing",
      "payment_status": "paid",
      "fulfillment_status": "processing",
      "subtotal": 599.98,
      "shipping_fee": 30000,
      "total_amount": 629.98,
      "shipping_name": "John Doe",
      "shipping_phone": "0123456789",
      "shipping_address": "123 Main Street",
      "shipping_city": "Ho Chi Minh City",
      "shipping_country": "Vietnam",
      "tracking_number": "VN123456789",
      "shipping_carrier": "Vietnam Post",
      "estimated_delivery": "2024-01-20T00:00:00Z",
      "items": [
        {
          "id": 1,
          "product_id": 101,
          "product_name": "Wireless Headphones",
          "product_image": "/uploads/products/headphones.jpg",
          "quantity": 2,
          "price": 299.99,
          "subtotal": 599.98
        }
      ],
      "status_history": [
        {
          "status": "processing",
          "from_status": "paid",
          "message": "Order confirmed by seller",
          "created_at": "2024-01-15T11:00:00Z"
        },
        {
          "status": "paid",
          "from_status": "pending",
          "message": "Payment received",
          "created_at": "2024-01-15T10:35:00Z"
        },
        {
          "status": "pending",
          "message": "Order created",
          "created_at": "2024-01-15T10:30:00Z"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

---

### 4. Cancel Order
**POST** `/api/orders/:id/cancel`

Cancel an order (only pending or paid orders).

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "reason": "Changed my mind"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "order": {
      "id": 1,
      "order_number": "ORD-20240115-00001",
      "status": "cancelled",
      "cancel_reason": "Changed my mind",
      "cancelled_at": "2024-01-15T10:45:00Z"
    }
  },
  "message": "Order cancelled successfully"
}
```

**Error Responses:**

```json
// Order not found (404)
{
  "success": false,
  "error": "Order not found"
}

// Cannot cancel (400)
{
  "success": false,
  "error": "Order cannot be cancelled at this stage"
}
```

---

### 5. Update Order Status (Seller/Admin)
**PUT** `/api/orders/:id/status`

Update order status (Seller or Admin only).

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "status": "shipped",
  "message": "Order has been shipped"
}
```

**Valid Status Values:**
- `pending`
- `paid`
- `processing`
- `shipped`
- `delivered`
- `cancelled`
- `refunded`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "order": {
      "id": 1,
      "status": "shipped",
      "shipped_at": "2024-01-16T10:00:00Z",
      "tracking_number": "VN123456789"
    }
  },
  "message": "Order status updated successfully"
}
```

---

### 6. Get Order Tracking
**GET** `/api/orders/:id/tracking`

Get tracking information for an order.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "tracking": [
      {
        "id": 1,
        "tracking_number": "VN123456789",
        "status": "in_transit",
        "location": "Ho Chi Minh City Hub",
        "description": "Package departed from sorting center",
        "timestamp": "2024-01-17T08:00:00Z"
      },
      {
        "id": 2,
        "tracking_number": "VN123456789",
        "status": "shipped",
        "location": "Ho Chi Minh City",
        "description": "Package has been shipped",
        "timestamp": "2024-01-16T10:00:00Z"
      }
    ]
  }
}
```

---

### 7. Get Order Statistics
**GET** `/api/orders/statistics`

Get statistics for user's orders.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "statistics": {
      "total_orders": 25,
      "total_revenue": 15000000,
      "pending_orders": 2,
      "completed_orders": 20,
      "cancelled_orders": 3,
      "average_order_value": 600000
    }
  }
}
```

---

## Order Status Reference

### Status Values

| Status | Description | Can Cancel | Can Refund |
|--------|-------------|------------|------------|
| `pending` | Awaiting payment | Yes | No |
| `paid` | Payment confirmed | Yes | Yes |
| `processing` | Being prepared | No | Yes |
| `shipped` | In transit | No | Yes |
| `delivered` | Received | No | Yes |
| `cancelled` | Cancelled | N/A | No |
| `refunded` | Refunded | N/A | N/A |

### Payment Status Values

| Status | Description |
|--------|-------------|
| `pending` | Awaiting payment |
| `paid` | Payment successful |
| `failed` | Payment failed |
| `refunded` | Payment refunded |
| `partial_refund` | Partially refunded |

---

## Payment Methods

| Method | Code | Description |
|--------|------|-------------|
| Cash on Delivery | `cod` | Pay upon receipt |
| Bank Transfer | `bank_transfer` | Direct bank transfer |
| Credit Card | `credit_card` | Visa, Mastercard, etc. |
| E-Wallet | `e_wallet` | MoMo, ZaloPay, etc. |

---

## Shipping Methods

| Method | Code | Estimated Delivery | Fee |
|--------|------|-------------------|-----|
| Standard | `standard` | 3-5 days | 30,000 VND |
| Express | `express` | 1-2 days | 50,000 VND |
| Same Day | `same_day` | Same day | 100,000 VND |
| Free | `free` | 5-7 days | 0 VND (orders > 500,000) |

---

## Order Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                    ORDER LIFECYCLE                           │
└─────────────────────────────────────────────────────────────┘

1. CREATED (pending)
   - Order created from cart
   - Inventory locked
   - Awaiting payment (15-30 min timeout)

2. PAID
   - Payment confirmed
   - Seller notified
   - Inventory still locked

3. PROCESSING
   - Seller confirms order
   - Items being packed
   - Shipping label created

4. SHIPPED
   - Package with carrier
   - Tracking number active
   - Customer can track

5. DELIVERED
   - Customer received order
   - Inventory decreased
   - Order completed

Alternative flows:
- CANCELLED: At any point before shipped
- REFUNDED: After delivery with return
```

---

## cURL Examples

### Checkout
```bash
curl -X POST http://localhost:8080/api/orders/checkout \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_info": {
      "name": "John Doe",
      "phone": "0123456789",
      "address": "123 Main Street",
      "district": "District 1",
      "city": "Ho Chi Minh City",
      "country": "Vietnam"
    },
    "payment_method": "cod",
    "buyer_note": "Please deliver before 5 PM"
  }'
```

### Get Orders
```bash
curl -X GET "http://localhost:8080/api/orders?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Order Details
```bash
curl -X GET http://localhost:8080/api/orders/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Cancel Order
```bash
curl -X POST http://localhost:8080/api/orders/1/cancel \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Changed my mind"}'
```

### Update Order Status (Seller)
```bash
curl -X PUT http://localhost:8080/api/orders/1/status \
  -H "Authorization: Bearer SELLER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "shipped", "message": "Order shipped"}'
```

### Get Tracking
```bash
curl -X GET http://localhost:8080/api/orders/1/tracking \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Error Handling

### Common Errors

| Error | HTTP Code | Description |
|-------|-----------|-------------|
| Cart is empty | 400 | No items in cart |
| Insufficient stock | 400 | Product out of stock |
| Product unavailable | 400 | Product deleted/inactive |
| Invalid shipping info | 400 | Missing required fields |
| Invalid payment method | 400 | Unsupported payment method |
| Order not found | 404 | Order doesn't exist |
| Cannot cancel | 400 | Order already shipped |
| Invalid status transition | 400 | Invalid status change |

---

## Inventory Locking

When an order is created:
1. **Stock is reserved** (not decreased)
2. **15-30 minute timeout** for payment
3. **If payment succeeds**: Stock decreased on delivery
4. **If payment fails/times out**: Stock released, order cancelled

This prevents overselling while allowing payment processing time.

---

This Order API is production-ready with full lifecycle management and inventory control.
