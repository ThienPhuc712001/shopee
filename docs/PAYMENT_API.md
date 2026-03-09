# Payment System API Documentation

## Base URL
```
http://localhost:8080/api
```

---

## Payment Endpoints

### 1. Create Payment
**POST** `/api/payments/create`

Create a new payment for an order.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "order_id": 1,
  "payment_method": "credit_card",
  "provider": "stripe",
  "save_method": false
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "payment": {
      "id": 1,
      "order_id": 1,
      "transaction_id": "TXN-20240115103000-abc123",
      "payment_method": "credit_card",
      "payment_provider": "stripe",
      "amount": 629.98,
      "currency": "USD",
      "status": "pending",
      "created_at": "2024-01-15T10:30:00Z"
    }
  },
  "message": "Payment created successfully"
}
```

---

### 2. Get Payment by Order
**GET** `/api/payments/order/:order_id`

Get payment information for a specific order.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "payment": {
      "id": 1,
      "order_id": 1,
      "transaction_id": "TXN-20240115103000-abc123",
      "payment_method": "credit_card",
      "amount": 629.98,
      "currency": "USD",
      "status": "paid",
      "paid_at": "2024-01-15T10:35:00Z",
      "transactions": [
        {
          "id": 1,
          "type": "charge",
          "amount": 629.98,
          "status": "completed",
          "gateway_id": "ch_xxx",
          "processed_at": "2024-01-15T10:35:00Z"
        }
      ]
    }
  }
}
```

---

### 3. Confirm Payment
**POST** `/api/payments/confirm`

Confirm a payment (for manual confirmation flows like bank transfer).

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "transaction_id": "TXN-20240115103000-abc123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "payment": {
      "id": 1,
      "transaction_id": "TXN-20240115103000-abc123",
      "status": "paid",
      "paid_at": "2024-01-15T10:35:00Z"
    }
  },
  "message": "Payment confirmed successfully"
}
```

---

### 4. Request Refund
**POST** `/api/payments/refund`

Request a refund for a payment.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body (Full Refund):**
```json
{
  "payment_id": 1,
  "reason": "Product defective",
  "type": "full"
}
```

**Request Body (Partial Refund):**
```json
{
  "payment_id": 1,
  "amount": 100.00,
  "reason": "Partial return",
  "type": "partial"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "refund": {
      "id": 1,
      "refund_number": "REF-20240115-abc123",
      "payment_id": 1,
      "order_id": 1,
      "amount": 629.98,
      "reason": "Product defective",
      "status": "pending",
      "type": "full",
      "created_at": "2024-01-15T11:00:00Z"
    }
  },
  "message": "Refund requested successfully"
}
```

---

### 5. Get User Payments
**GET** `/api/payments`

Get all payments for the current user.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| limit | int | 10 | Items per page |

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "payments": [
      {
        "id": 1,
        "order_id": 1,
        "transaction_id": "TXN-20240115103000-abc123",
        "amount": 629.98,
        "currency": "USD",
        "status": "paid",
        "payment_method": "credit_card",
        "paid_at": "2024-01-15T10:35:00Z"
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

### 6. Save Payment Method
**POST** `/api/payments/methods`

Save a payment method for future use.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "type": "credit_card",
  "provider": "stripe",
  "name": "Visa ending 4242",
  "last_four": "4242",
  "expiry_month": 12,
  "expiry_year": 2025,
  "token": "tok_visa_4242",
  "is_default": true
}
```

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "method": {
      "id": 1,
      "type": "credit_card",
      "provider": "stripe",
      "name": "Visa ending 4242",
      "last_four": "4242",
      "is_default": true
    }
  },
  "message": "Payment method saved successfully"
}
```

---

### 7. Get Payment Methods
**GET** `/api/payments/methods`

Get all saved payment methods.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "methods": [
      {
        "id": 1,
        "type": "credit_card",
        "provider": "stripe",
        "name": "Visa ending 4242",
        "last_four": "4242",
        "is_default": true
      },
      {
        "id": 2,
        "type": "e_wallet",
        "provider": "paypal",
        "name": "PayPal",
        "is_default": false
      }
    ]
  }
}
```

---

### 8. Delete Payment Method
**DELETE** `/api/payments/methods/:id`

Delete a saved payment method.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Payment method deleted successfully"
}
```

---

### 9. Set Default Payment Method
**POST** `/api/payments/methods/:id/default`

Set a payment method as default.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Default payment method updated"
}
```

---

### 10. Get Payment Statistics
**GET** `/api/payments/statistics`

Get statistics for user's payments.

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
      "total_payments": 25,
      "total_amount": 15000.00,
      "successful_payments": 23,
      "failed_payments": 2,
      "pending_payments": 0,
      "total_refunds": 1,
      "refund_amount": 500.00
    }
  }
}
```

---

### 11. Payment Webhook
**POST** `/api/payments/webhook`

Handle payment gateway webhook notifications (no auth required).

**Headers:**
```
Content-Type: application/json
X-Webhook-Signature: <signature>
```

**Request Body:**
```json
{
  "event": "payment.completed",
  "transaction_id": "TXN-20240115103000-abc123",
  "payment_id": "1",
  "status": "paid",
  "amount": 629.98,
  "currency": "USD",
  "signature": "hmac_signature_here",
  "metadata": {
    "order_id": "1",
    "user_id": "123"
  },
  "timestamp": 1705312500
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Webhook processed successfully"
}
```

---

## Payment Methods Reference

| Method | Code | Provider Examples |
|--------|------|-------------------|
| Credit Card | `credit_card` | Stripe, PayPal |
| Debit Card | `debit_card` | Stripe, VNPay |
| Bank Transfer | `bank_transfer` | Local banks |
| E-Wallet | `e_wallet` | MoMo, ZaloPay |
| Cash on Delivery | `cod` | N/A |
| PayPal | `paypal` | PayPal |

---

## Payment Status Reference

| Status | Description |
|--------|-------------|
| `pending` | Payment created, awaiting processing |
| `processing` | Payment being processed |
| `paid` | Payment successful |
| `failed` | Payment failed |
| `cancelled` | Payment cancelled |
| `refunded` | Payment fully refunded |
| `partial_refund` | Payment partially refunded |

---

## Refund Status Reference

| Status | Description |
|--------|-------------|
| `pending` | Refund requested, awaiting approval |
| `approved` | Refund approved |
| `processed` | Refund processed by gateway |
| `rejected` | Refund rejected |
| `failed` | Refund failed |

---

## Payment Flow Examples

### Credit Card Payment (Stripe)

```
1. Create Payment → Get client_secret
2. Client confirms with Stripe.js
3. Stripe sends webhook
4. System updates payment status
5. Order status updated to "paid"
```

### Bank Transfer

```
1. Create Payment → Get bank details
2. Customer transfers money
3. Bank confirms receipt
4. System confirms payment
5. Order status updated to "paid"
```

### Cash on Delivery

```
1. Create Payment → Status = "pending"
2. Order shipped
3. Customer pays on delivery
4. Seller confirms payment
5. Payment status = "paid"
```

---

## Error Responses

```json
// Payment not found (404)
{
  "success": false,
  "error": "Payment not found"
}

// Invalid payment method (400)
{
  "success": false,
  "error": "Invalid payment method"
}

// Payment already paid (400)
{
  "success": false,
  "error": "Payment already completed"
}

// Invalid signature (401)
{
  "success": false,
  "error": "Invalid webhook signature"
}

// Refund exceeds amount (400)
{
  "success": false,
  "error": "Refund amount exceeds payment"
}
```

---

## Security Measures

### Webhook Security
1. Verify HMAC signature
2. Validate gateway IP
3. Check transaction uniqueness
4. Idempotency handling

### Payment Security
1. Never store raw card data
2. Use tokenization
3. PCI DSS compliance
4. Encrypted transmission

---

## cURL Examples

### Create Payment
```bash
curl -X POST http://localhost:8080/api/payments/create \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 1,
    "payment_method": "credit_card",
    "provider": "stripe"
  }'
```

### Request Refund
```bash
curl -X POST http://localhost:8080/api/payments/refund \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": 1,
    "reason": "Product defective",
    "type": "full"
  }'
```

### Save Payment Method
```bash
curl -X POST http://localhost:8080/api/payments/methods \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "credit_card",
    "provider": "stripe",
    "name": "Visa ending 4242",
    "last_four": "4242",
    "token": "tok_visa_4242"
  }'
```

---

This Payment API is production-ready with full gateway integration and security measures.
