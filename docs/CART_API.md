# Shopping Cart Module API Documentation

## Base URL
```
http://localhost:8080/api
```

---

## Cart Endpoints

### 1. Get Cart
**GET** `/api/cart`

Get the current user's shopping cart with all items.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "cart": {
      "id": 1,
      "user_id": 123,
      "total_items": 5,
      "subtotal": 599.98,
      "discount": 0,
      "total": 599.98,
      "currency": "USD",
      "last_activity": "2024-01-15T10:30:00Z",
      "items": [
        {
          "id": 1,
          "product_id": 101,
          "variant_id": null,
          "quantity": 2,
          "price": 299.99,
          "subtotal": 599.98,
          "product": {
            "id": 101,
            "name": "Wireless Headphones",
            "slug": "wireless-headphones",
            "image": "/uploads/products/headphones.jpg"
          },
          "shop": {
            "id": 1,
            "name": "AudioTech Store"
          },
          "is_available": true,
          "stock_status": "in_stock"
        }
      ]
    }
  }
}
```

**Empty Cart Response:**
```json
{
  "success": true,
  "data": {
    "cart": {
      "id": 1,
      "user_id": 123,
      "total_items": 0,
      "subtotal": 0,
      "total": 0,
      "items": []
    }
  }
}
```

---

### 2. Get Cart Summary
**GET** `/api/cart/summary`

Get a summary of the cart (totals only, without items).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "summary": {
      "total_items": 5,
      "subtotal": 599.98,
      "discount": 0,
      "total": 599.98,
      "currency": "USD"
    }
  }
}
```

---

### 3. Get Cart Statistics
**GET** `/api/cart/stats`

Get cart statistics including stock and price change alerts.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "stats": {
      "total_items": 5,
      "total_products": 3,
      "subtotal": 599.98,
      "estimated_total": 599.98,
      "has_out_of_stock": false,
      "has_price_changed": false
    }
  }
}
```

---

### 4. Add to Cart
**POST** `/api/cart/add`

Add a product to the shopping cart.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "product_id": 101,
  "variant_id": null,
  "quantity": 2
}
```

**With Variant:**
```json
{
  "product_id": 101,
  "variant_id": 5,
  "quantity": 1
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "cart": {
      "id": 1,
      "user_id": 123,
      "total_items": 5,
      "subtotal": 899.97,
      "total": 899.97,
      "items": [
        {
          "id": 1,
          "product_id": 101,
          "quantity": 2,
          "price": 299.99,
          "subtotal": 599.98,
          "product": {
            "id": 101,
            "name": "Wireless Headphones"
          }
        },
        {
          "id": 2,
          "product_id": 102,
          "quantity": 1,
          "price": 299.99,
          "subtotal": 299.99,
          "product": {
            "id": 102,
            "name": "USB-C Cable"
          }
        }
      ]
    }
  },
  "message": "Item added to cart successfully"
}
```

**Error Responses:**

```json
// Product not found (404)
{
  "success": false,
  "error": "Product not found"
}

// Insufficient stock (400)
{
  "success": false,
  "error": "Insufficient stock available"
}

// Invalid quantity (400)
{
  "success": false,
  "error": "Invalid quantity"
}

// Product not available (400)
{
  "success": false,
  "error": "Product is not available"
}
```

---

### 5. Update Cart Item
**PUT** `/api/cart/items/:id`

Update the quantity of a cart item.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "quantity": 5
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "cart": {
      "id": 1,
      "total_items": 8,
      "subtotal": 1499.95,
      "total": 1499.95,
      "items": [
        {
          "id": 1,
          "product_id": 101,
          "quantity": 5,
          "price": 299.99,
          "subtotal": 1499.95
        }
      ]
    }
  },
  "message": "Cart item updated successfully"
}
```

**Error Responses:**

```json
// Cart item not found (404)
{
  "success": false,
  "error": "Cart item not found"
}

// Insufficient stock (400)
{
  "success": false,
  "error": "Insufficient stock available"
}

// Invalid quantity (400)
{
  "success": false,
  "error": "Invalid quantity"
}
```

---

### 6. Remove from Cart
**DELETE** `/api/cart/items/:id`

Remove an item from the shopping cart.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "cart": {
      "id": 1,
      "total_items": 3,
      "subtotal": 299.99,
      "total": 299.99,
      "items": []
    }
  },
  "message": "Item removed from cart successfully"
}
```

**Error Responses:**

```json
// Cart item not found (404)
{
  "success": false,
  "error": "Cart item not found"
}

// Cart not found (404)
{
  "success": false,
  "error": "Cart not found"
}
```

---

### 7. Clear Cart
**DELETE** `/api/cart/clear`

Remove all items from the shopping cart.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Cart cleared successfully"
}
```

---

### 8. Prepare Checkout
**GET** `/api/cart/checkout`

Prepare cart for checkout with validation and shipping calculation.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "checkout": {
      "items": [
        {
          "product_id": 101,
          "variant_id": null,
          "quantity": 2,
          "price": 299.99,
          "subtotal": 599.98,
          "shop_id": 1,
          "shop_name": "AudioTech Store",
          "product_name": "Wireless Headphones",
          "product_image": "/uploads/products/headphones.jpg",
          "is_available": true,
          "stock_status": "in_stock"
        }
      ],
      "subtotal": 599.98,
      "shipping_fee": 30000,
      "discount": 0,
      "total": 629.98,
      "currency": "USD"
    }
  }
}
```

**Error Responses:**

```json
// Cart empty (400)
{
  "success": false,
  "error": "Cart is empty"
}

// Validation failed (400)
{
  "success": false,
  "error": "cart validation failed: [insufficient stock for Wireless Headphones]"
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 400 | Bad Request - Invalid input or validation failed |
| 401 | Unauthorized - Invalid or missing token |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

---

## Business Rules

### Cart Rules

1. **One Cart Per User**
   - Each user has exactly one cart
   - Cart is created automatically on first add
   - Cart persists across sessions

2. **Item Quantity**
   - Minimum: 1
   - Maximum: 999
   - Cannot exceed available stock

3. **Stock Validation**
   - Stock checked on add
   - Stock checked on quantity update
   - Stock reserved during checkout

4. **Price Updates**
   - Price captured at add time
   - Price synced on cart view
   - User notified of price changes

5. **Product Availability**
   - Inactive products cannot be added
   - Out of stock items flagged
   - Removed after prolonged unavailability

---

## Cart Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                    CART LIFECYCLE                            │
└─────────────────────────────────────────────────────────────┘

1. CREATE
   - User adds first item
   - Cart created automatically
   - User ID associated

2. UPDATE
   - Add items
   - Update quantities
   - Remove items
   - Totals recalculated

3. VALIDATE (Checkout)
   - Check stock availability
   - Check product status
   - Calculate shipping
   - Apply discounts

4. RESERVE (Checkout)
   - Stock reserved for items
   - 15-minute hold
   - Prevents overselling

5. CONVERT (Order)
   - Cart items → Order items
   - Stock decreased
   - Cart cleared

6. ABANDON
   - No activity for 30 days
   - Stock released
   - Cart archived
```

---

## cURL Examples

### Get Cart
```bash
curl -X GET http://localhost:8080/api/cart \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Add to Cart
```bash
curl -X POST http://localhost:8080/api/cart/add \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 101,
    "quantity": 2
  }'
```

### Add Variant to Cart
```bash
curl -X POST http://localhost:8080/api/cart/add \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 101,
    "variant_id": 5,
    "quantity": 1
  }'
```

### Update Cart Item
```bash
curl -X PUT http://localhost:8080/api/cart/items/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 5}'
```

### Remove from Cart
```bash
curl -X DELETE http://localhost:8080/api/cart/items/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Clear Cart
```bash
curl -X DELETE http://localhost:8080/api/cart/clear \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Prepare Checkout
```bash
curl -X GET http://localhost:8080/api/cart/checkout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Performance Optimization

### Indexes

```sql
-- Cart lookups
CREATE INDEX IX_Carts_UserID ON Carts(user_id);
CREATE INDEX IX_Carts_LastActivity ON Carts(last_activity);

-- Cart item lookups
CREATE INDEX IX_CartItems_CartID ON CartItems(cart_id);
CREATE INDEX IX_CartItems_ProductID ON CartItems(product_id);
CREATE UNIQUE INDEX IX_CartItems_Unique ON CartItems(cart_id, product_id, variant_id);
```

### Caching Strategy

- Cart data cached for 5 minutes
- Cache invalidated on cart changes
- Redis recommended for high traffic

### Concurrency Handling

- Database transactions for updates
- Row-level locking during checkout
- Optimistic locking for cart items

---

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "Cart not found" | User has no cart | Cart created on first add |
| "Cart item not found" | Item doesn't exist | Refresh cart, item may be removed |
| "Insufficient stock" | Product out of stock | Reduce quantity or remove item |
| "Invalid quantity" | Quantity < 1 or > 999 | Use valid quantity |
| "Product not available" | Product inactive/deleted | Remove from cart |

---

This cart API is production-ready and optimized for high-traffic e-commerce platforms.
