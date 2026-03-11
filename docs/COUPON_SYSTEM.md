# Coupon and Promotion System Documentation

## Overview

The Coupon and Promotion System is a comprehensive discount management module for the e-commerce platform. It allows administrators to create and manage discount coupons, and enables users to apply coupons during checkout.

---

## Table of Contents

1. [Business Flow](#business-flow)
2. [Coupon Types](#coupon-types)
3. [Database Schema](#database-schema)
4. [Coupon Rules](#coupon-rules)
5. [Discount Calculation](#discount-calculation)
6. [API Endpoints](#api-endpoints)
7. [Order Integration](#order-integration)
8. [Security Measures](#security-measures)
9. [Usage Examples](#usage-examples)

---

## Business Flow

### Complete Coupon Lifecycle

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   ADMIN     │     │    USER     │     │   SYSTEM    │
│  CREATES    │     │   ADDS TO   │     │  VALIDATES  │
│   COUPON    │────▶│    CART     │────▶│   COUPON    │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   USAGE     │◀────│   ORDER     │◀────│ CALCULATES  │
│  RECORDED   │     │ COMPLETED   │     │ DISCOUNT    │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Step-by-Step Flow

| Step | Actor | Action | Description |
|------|-------|--------|-------------|
| 1 | Admin | Create Coupon | Administrator creates a coupon with rules and expiration |
| 2 | User | Add to Cart | User adds products to shopping cart |
| 3 | User | Enter Code | User enters coupon code at checkout |
| 4 | System | Validate | System validates coupon against all rules |
| 5 | System | Calculate | System calculates discount amount |
| 6 | System | Apply | Discount applied to order total |
| 7 | User | Complete Order | User completes purchase |
| 8 | System | Record Usage | Usage count incremented, usage record created |

---

## Coupon Types

### 1. Percentage Discount (`percentage`)

Applies a percentage-based discount to the order total.

**Example:**
- Code: `SAVE10`
- Discount Value: `10`
- Order Total: `$100`
- **Discount: `$10` (10%)**

**Configuration:**
```json
{
  "discount_type": "percentage",
  "discount_value": 10,
  "max_discount": 50
}
```

**Notes:**
- `max_discount` limits the maximum discount amount
- Useful for sales events (e.g., "Up to 50% off")

---

### 2. Fixed Amount Discount (`fixed`)

Applies a fixed monetary discount to the order total.

**Example:**
- Code: `SAVE5`
- Discount Value: `5`
- Order Total: `$50`
- **Discount: `$5`**

**Configuration:**
```json
{
  "discount_type": "fixed",
  "discount_value": 5
}
```

**Notes:**
- Discount cannot exceed order total
- Ideal for promotional offers (e.g., "$5 off your next order")

---

### 3. Free Shipping (`free_shipping`)

Waives the shipping fee for the order.

**Example:**
- Code: `FREESHIP`
- Order Total: `$30`
- Shipping Fee: `$5`
- **Discount: `$5` (shipping fee waived)**

**Configuration:**
```json
{
  "discount_type": "free_shipping",
  "min_order_value": 30
}
```

**Notes:**
- Often combined with minimum order value
- Shipping discount applied separately from product discount

---

## Database Schema

### Coupons Table

```sql
CREATE TABLE coupons (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    code NVARCHAR(50) NOT NULL UNIQUE,
    name NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX),
    discount_type NVARCHAR(50) NOT NULL,
    discount_value DECIMAL(18,2) NOT NULL,
    max_discount DECIMAL(18,2) DEFAULT 0,
    min_order_value DECIMAL(18,2) DEFAULT 0,
    max_order_value DECIMAL(18,2) DEFAULT 0,
    usage_limit INT DEFAULT 0,
    used_count INT DEFAULT 0,
    usage_limit_per_user INT DEFAULT 1,
    start_date DATETIME,
    end_date DATETIME NOT NULL,
    is_active BIT DEFAULT 1,
    status NVARCHAR(50) DEFAULT 'active',
    applicable_categories NVARCHAR(MAX),
    applicable_products NVARCHAR(MAX),
    excluded_categories NVARCHAR(MAX),
    excluded_products NVARCHAR(MAX),
    user_restricted BIT DEFAULT 0,
    restricted_users NVARCHAR(MAX),
    created_by BIGINT,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    deleted_at DATETIME
);
```

### Coupon Usages Table

```sql
CREATE TABLE coupon_usages (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    coupon_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    discount_amount DECIMAL(18,2) NOT NULL,
    used_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_coupon_usages_coupon FOREIGN KEY (coupon_id) REFERENCES coupons(id),
    CONSTRAINT FK_coupon_usages_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_coupon_usages_order FOREIGN KEY (order_id) REFERENCES orders(id)
);
```

### Orders Table (Updated)

```sql
ALTER TABLE orders ADD coupon_id BIGINT NULL;
ALTER TABLE orders ADD coupon_code NVARCHAR(50);
ALTER TABLE orders ADD coupon_discount DECIMAL(18,2) DEFAULT 0;
```

### Relationships

```
coupons (1) ──────< coupon_usages (>────── users (1)
     │                                       │
     │                                       │
     └──────────────< orders (>──────────────┘
```

---

## Coupon Rules

### Validation Rules

| Rule | Description | Error Message |
|------|-------------|---------------|
| **Coupon Exists** | Code must exist in database | "Invalid coupon code" |
| **Not Expired** | Current date before `end_date` | "Coupon has expired" |
| **Active Status** | `is_active` must be true | "Coupon is inactive" |
| **Started** | Current date after `start_date` | "Coupon has not started yet" |
| **Usage Limit** | `used_count < usage_limit` | "Coupon usage limit reached" |
| **Per User Limit** | User usage < `usage_limit_per_user` | "Coupon already used by this user" |
| **Min Order Value** | `order_total >= min_order_value` | "Minimum order value is $X.XX" |
| **Max Order Value** | `order_total <= max_order_value` | "Maximum order value is $X.XX" |
| **Product Restrictions** | Products in cart match `applicable_products` | "Coupon not applicable to cart items" |
| **Category Restrictions** | Categories match `applicable_categories` | "Coupon not applicable to selected categories" |
| **User Restrictions** | User ID in `restricted_users` list | "Coupon is restricted to specific users" |

### Validation Flow

```go
func ValidateCoupon(coupon, orderTotal, userID, products) error {
    1. Check coupon exists
    2. Check is_active == true
    3. Check not expired (end_date > now)
    4. Check started (start_date <= now)
    5. Check usage_limit not reached
    6. Check user usage limit
    7. Check min_order_value
    8. Check max_order_value
    9. Check product/category restrictions
    10. Check user restrictions
}
```

---

## Discount Calculation

### Formula

```
if discount_type == 'percentage':
    discount = order_total × (discount_value / 100)
    if max_discount > 0:
        discount = min(discount, max_discount)

elif discount_type == 'fixed':
    discount = discount_value
    discount = min(discount, order_total)

elif discount_type == 'free_shipping':
    discount = shipping_fee

final_price = order_total - discount
```

### Examples

#### Example 1: Percentage Discount

```
Order Total: $100
Coupon: 10% off
Max Discount: $15

Calculation:
discount = 100 × (10/100) = $10
discount < max_discount, so discount = $10
Final Price: $100 - $10 = $90
```

#### Example 2: Fixed Discount

```
Order Total: $50
Coupon: $5 off

Calculation:
discount = $5
Final Price: $50 - $5 = $45
```

#### Example 3: Free Shipping

```
Order Total: $80
Shipping Fee: $10
Coupon: Free Shipping

Calculation:
discount = $10 (shipping fee)
Final Price: $80 + $0 (shipping) = $80
```

---

## API Endpoints

### Admin APIs

#### Create Coupon
```http
POST /api/coupons
Authorization: Bearer {token}
Role: admin

Request:
{
  "code": "SUMMER20",
  "name": "Summer Sale",
  "description": "20% off summer collection",
  "discount_type": "percentage",
  "discount_value": 20,
  "max_discount": 50,
  "min_order_value": 100,
  "usage_limit": 1000,
  "usage_limit_per_user": 1,
  "start_date": "2024-06-01T00:00:00Z",
  "end_date": "2024-08-31T23:59:59Z",
  "is_active": true
}

Response (201):
{
  "success": true,
  "message": "Coupon created successfully",
  "data": { ...coupon object... }
}
```

#### Get All Coupons
```http
GET /api/coupons?page=1&limit=20&status=active
Authorization: Bearer {token}
Role: admin

Response (200):
{
  "success": true,
  "data": {
    "coupons": [...],
    "pagination": {
      "total": 50,
      "page": 1,
      "limit": 20,
      "total_pages": 3
    }
  }
}
```

#### Delete Coupon
```http
DELETE /api/coupons/{id}
Authorization: Bearer {token}
Role: admin

Response (200):
{
  "success": true,
  "message": "Coupon deleted successfully"
}
```

#### Get Coupon Statistics
```http
GET /api/coupons/stats
Authorization: Bearer {token}
Role: admin

Response (200):
{
  "success": true,
  "data": {
    "total_coupons": 50,
    "active_coupons": 35,
    "expired_coupons": 10,
    "total_usage": 1250,
    "total_discount": 12500.00,
    "average_discount": 10.00
  }
}
```

---

### User APIs

#### Apply Coupon
```http
POST /api/coupons/apply
Content-Type: application/json

Request:
{
  "code": "SUMMER20",
  "order_total": 150,
  "user_id": 123,
  "product_ids": [1, 2, 3],
  "category_ids": [5, 10]
}

Response (200):
{
  "success": true,
  "message": "Coupon applied successfully",
  "data": {
    "coupon": { ... },
    "discount_amount": 30,
    "original_total": 150,
    "final_total": 120
  }
}

Response (400) - Invalid:
{
  "success": false,
  "message": "Coupon has expired"
}
```

#### Get My Coupon Usages
```http
GET /api/coupons/my-usages?page=1&limit=20
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "usages": [
      {
        "id": 1,
        "coupon_id": 5,
        "coupon": { "code": "SUMMER20", ... },
        "order_id": 100,
        "discount_amount": 30,
        "used_at": "2024-06-15T10:30:00Z"
      }
    ],
    "pagination": { ... }
  }
}
```

---

## Order Integration

### Checkout Flow with Coupon

```
┌─────────────────────────────────────────────────────────────┐
│                    CHECKOUT PROCESS                          │
├─────────────────────────────────────────────────────────────┤
│  1. User clicks "Checkout"                                  │
│  2. System loads cart items                                 │
│  3. System calculates subtotal                              │
│  4. User enters coupon code                                 │
│  5. System validates coupon                                 │
│  6. System calculates discount                              │
│  7. System updates order total                              │
│  8. User confirms order                                     │
│  9. System creates order with coupon_id                     │
│  10. System clears cart                                     │
│  11. Order delivered → System records coupon usage          │
└─────────────────────────────────────────────────────────────┘
```

### Order Model Integration

```go
type Order struct {
    // ... other fields ...
    
    // Coupon fields
    CouponID       *uint    `json:"coupon_id"`
    CouponCode     string   `json:"coupon_code"`
    CouponDiscount float64  `json:"coupon_discount"`
    
    // Pricing
    Subtotal       float64  `json:"subtotal"`
    ShippingFee    float64  `json:"shipping_fee"`
    TotalAmount    float64  `json:"total_amount"`
}
```

### Code Flow

```go
// In order_service_enhanced.go
func (s *orderServiceEnhanced) createSingleOrder(...) (*model.Order, error) {
    // 1. Calculate subtotal
    subtotal := calculateSubtotal(cartItems)
    
    // 2. Apply coupon if provided
    if input.CouponCode != "" {
        coupon, _ := s.couponRepo.GetCouponByCode(code)
        result, _ := s.couponSvc.ApplyCoupon(coupon, subtotal, userID)
        
        if result.Success {
            couponDiscount = result.DiscountAmount
            couponID = &coupon.ID
        }
    }
    
    // 3. Calculate final total
    totalAmount = subtotal + shippingFee - couponDiscount
    
    // 4. Create order with coupon info
    order := &model.Order{
        CouponID: couponID,
        CouponCode: couponCode,
        CouponDiscount: couponDiscount,
        TotalAmount: totalAmount,
    }
    
    return order, nil
}
```

---

## Security Measures

### 1. Usage Limit Enforcement

```go
// Atomic increment with row lock
func (r *CouponRepository) IncrementUsageCount(couponID uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Lock row for update
        var coupon Coupon
        tx.Set("gorm:query_option", "WITH (UPDLOCK)").
            First(&coupon, couponID)
        
        // Check limit
        if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
            return errors.New("usage limit reached")
        }
        
        // Increment
        coupon.UsedCount++
        tx.Save(&coupon)
        
        return nil
    })
}
```

### 2. Per-User Restrictions

- Track usage per user in `coupon_usages` table
- Validate before applying: `SELECT COUNT(*) FROM coupon_usages WHERE coupon_id = ? AND user_id = ?`
- Enforce `usage_limit_per_user`

### 3. Race Condition Prevention

- Use database transactions
- Row-level locking (`UPDLOCK` in SQL Server)
- Atomic operations for usage count

### 4. Validation at Multiple Layers

| Layer | Validation |
|-------|------------|
| **API Handler** | Input format, required fields |
| **Service** | Business rules, coupon validity |
| **Repository** | Database constraints, atomic operations |

### 5. Audit Trail

- All coupon usage recorded in `coupon_usages`
- Links to user and order
- Timestamp and discount amount tracked

---

## Usage Examples

### Example 1: Create a 10% Discount Coupon

```bash
curl -X POST http://localhost:8080/api/coupons \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SAVE10",
    "name": "10% Off Sale",
    "discount_type": "percentage",
    "discount_value": 10,
    "max_discount": 50,
    "min_order_value": 100,
    "usage_limit": 1000,
    "usage_limit_per_user": 1,
    "end_date": "2024-12-31T23:59:59Z",
    "is_active": true
  }'
```

### Example 2: Create a Free Shipping Coupon

```bash
curl -X POST http://localhost:8080/api/coupons \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "FREESHIP",
    "name": "Free Shipping",
    "discount_type": "free_shipping",
    "min_order_value": 50,
    "usage_limit": 500,
    "usage_limit_per_user": 3,
    "end_date": "2024-12-31T23:59:59Z",
    "is_active": true
  }'
```

### Example 3: Apply Coupon at Checkout

```bash
curl -X POST http://localhost:8080/api/coupons/apply \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SAVE10",
    "order_total": 150,
    "user_id": 123
  }'

# Response:
{
  "success": true,
  "message": "Coupon applied successfully",
  "data": {
    "discount_amount": 15,
    "original_total": 150,
    "final_total": 135
  }
}
```

### Example 4: Checkout with Coupon

```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Authorization: Bearer {user_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_info": {
      "name": "John Doe",
      "phone": "123456789",
      "address": "123 Main St",
      "city": "New York",
      "district": "Manhattan"
    },
    "payment_method": "credit_card",
    "coupon_code": "SAVE10"
  }'
```

---

## Implementation Checklist

### Completed ✅

- [x] Coupon model with all fields
- [x] CouponUsage model
- [x] Coupon repository with CRUD operations
- [x] Coupon service with validation logic
- [x] Coupon handler with HTTP endpoints
- [x] Coupon routes (admin and user)
- [x] Database migration script
- [x] Order integration (coupon fields in orders)
- [x] Checkout flow with coupon application
- [x] Coupon usage tracking
- [x] Free shipping support
- [x] Percentage and fixed discount support

### Future Enhancements 📋

- [ ] Coupon combination rules (stackable/non-stackable)
- [ ] Automatic coupon application (cart rules)
- [ ] Coupon recommendation engine
- [ ] A/B testing for coupon effectiveness
- [ ] Bulk coupon generation
- [ ] Coupon import/export
- [ ] Advanced analytics dashboard
- [ ] Email notifications for expiring coupons

---

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "Coupon not found" | Code typo, deleted coupon | Verify code, check soft delete |
| "Coupon expired" | Past end_date | Update end_date or create new coupon |
| "Usage limit reached" | used_count >= usage_limit | Increase usage_limit |
| "Already used" | User exceeded usage_limit_per_user | Use different account or wait |
| "Min order value not met" | Cart total too low | Add more items to cart |

### Debug Queries

```sql
-- Check coupon status
SELECT code, status, is_active, used_count, usage_limit, start_date, end_date
FROM coupons WHERE code = 'SAVE10';

-- Check user's usage
SELECT cu.*, c.code 
FROM coupon_usages cu 
JOIN coupons c ON cu.coupon_id = c.id 
WHERE cu.user_id = 123;

-- Check order's coupon
SELECT o.order_number, o.coupon_code, o.coupon_discount, c.name
FROM orders o
LEFT JOIN coupons c ON o.coupon_id = c.id
WHERE o.id = 100;
```

---

## License

Internal use only - E-Commerce Platform
