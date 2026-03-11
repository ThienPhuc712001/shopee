# Inventory Management System - Documentation

## Overview

The Inventory Management System manages product stock quantities, handles stock reservation for orders, and prevents overselling through atomic database operations.

---

## Table of Contents

1. [Business Flow](#business-flow)
2. [Inventory Rules](#inventory-rules)
3. [Database Schema](#database-schema)
4. [Stock States](#stock-states)
5. [API Endpoints](#api-endpoints)
6. [Order Integration](#order-integration)
7. [Concurrency Control](#concurrency-control)
8. [Examples](#examples)

---

## Business Flow

### How Inventory Management Works

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  1. Admin   │────▶│  2. Users    │────▶│  3. Order   │
│  Sets Stock │     │  Browse      │     │  Placed     │
│  Quantity   │     │  Products    │     │             │
└─────────────┘     └──────────────┘     └─────────────┘
                                              │
                                              ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  8. Order   │◀────│  7. Stock    │◀────│  5. System  │
│  Completed  │     │  Deducted    │     │  Reserves   │
│             │     │              │     │  Stock      │
└─────────────┘     └──────────────┘     └─────────────┘
                                              │
                                      ┌───────┴───────┐
                                      ▼               ▼
                               ┌─────────────┐ ┌─────────────┐
                               │  6. Payment │ │  Payment   │
                               │  Success    │ │  Failed    │
                               └─────────────┘ └─────────────┘
                                      │               │
                                      ▼               ▼
                               ┌─────────────┐ ┌─────────────┐
                               │  7. Deduct  │ │  Release    │
                               │  Stock      │ │  Reserved   │
                               └─────────────┘ └─────────────┘
```

### Step-by-Step Flow

| Step | Actor | Action | Result |
|------|-------|--------|--------|
| 1 | Admin | Sets product stock quantity | Inventory created/updated |
| 2 | Users | Browse products | View available stock |
| 3 | User | Places order | Order created |
| 4 | System | Checks available stock | Validation |
| 5 | System | Reserves stock | Stock reserved for order |
| 6a | Payment | Success | Proceed to deduct |
| 6b | Payment | Failed | Release reserved stock |
| 7a | System | Deducts stock | Stock reduced |
| 7b | System | Releases stock | Stock available again |
| 8 | System | Order completed | Transaction complete |

---

## Inventory Rules

### Core Rules

| Rule | Description | Implementation |
|------|-------------|----------------|
| **Stock cannot be negative** | Prevents invalid inventory | Validation before update |
| **Order quantity ≤ stock** | Prevents overselling | Check before reserve |
| **Atomic updates** | Prevents race conditions | Database transactions with row locking |
| **Reserved stock tracked** | Tracks committed stock | Separate reserved_quantity field |
| **All changes logged** | Audit trail | InventoryLog table |

### Why Atomic Updates Are Important

When multiple users try to order the same product simultaneously:

```
Time  User A                  User B
│     │                       │
├─────┼───────────────────────┤
│     │                       │
▼     ▼                       ▼
T1    Read stock: 5           │
│     │                       │
├─────┼───────────────────────┤
│     │                       │
│     ▼                       ▼
T2    Reserve 3               Read stock: 5
│     │                       │
├─────┼───────────────────────┤
│     │                       │
│     ▼                       ▼
T3    (pending...)            Reserve 4
│     │                       │
├─────┼───────────────────────┤
│     │                       │
│     ▼                       ▼
T4    Success: 2 left         FAIL: Only 2 available!
```

**Without atomic updates:** Both users could reserve stock, leading to overselling.

**With atomic updates (SELECT FOR UPDATE):**
- User A locks the row at T1
- User B waits until User A completes
- User B sees updated stock (2) and fails gracefully

---

## Database Schema

### Inventory Table

```sql
CREATE TABLE inventory (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    product_id BIGINT NOT NULL UNIQUE,
    stock_quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    available_quantity INT NOT NULL DEFAULT 0,
    warehouse_location NVARCHAR(100),
    last_stock_check DATETIME,
    reorder_level INT DEFAULT 0,
    reorder_quantity INT DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT FK_inventory_product 
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    
    INDEX IX_inventory_product_id (product_id),
    INDEX IX_inventory_available_quantity (available_quantity)
);
```

### Inventory Logs Table

```sql
CREATE TABLE inventory_logs (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    product_id BIGINT NOT NULL,
    inventory_id BIGINT NOT NULL,
    change_type NVARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    stock_before INT NOT NULL,
    stock_after INT NOT NULL,
    reserved_before INT NOT NULL,
    reserved_after INT NOT NULL,
    reference_type NVARCHAR(50),
    reference_id NVARCHAR(100),
    user_id BIGINT,
    reason NVARCHAR(MAX),
    ip_address NVARCHAR(45),
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT FK_inventory_logs_product 
        FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT FK_inventory_logs_inventory 
        FOREIGN KEY (inventory_id) REFERENCES inventory(id),
    
    INDEX IX_inventory_logs_product_id (product_id),
    INDEX IX_inventory_logs_reference_id (reference_id),
    INDEX IX_inventory_logs_created_at (created_at)
);
```

### Stock Alerts Table

```sql
CREATE TABLE stock_alerts (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    product_id BIGINT NOT NULL UNIQUE,
    alert_type NVARCHAR(50),
    threshold INT,
    current_stock INT,
    is_resolved BIT DEFAULT 0,
    resolved_at DATETIME,
    resolved_by BIGINT,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT FK_stock_alerts_product 
        FOREIGN KEY (product_id) REFERENCES products(id)
);
```

---

## Stock States

### Stock Formula

```
available_stock = stock_quantity - reserved_quantity
```

### Stock States

| State | Formula | Description |
|-------|---------|-------------|
| **In Stock** | `available > reorder_level` | Normal availability |
| **Low Stock** | `0 < available ≤ reorder_level` | Reorder recommended |
| **Out of Stock** | `available ≤ 0 AND stock = 0` | Cannot fulfill orders |
| **Pre-Order** | `available ≤ 0 AND stock > 0` | All stock reserved |

### Stock Status Flow

```
In Stock → Low Stock → Out of Stock
    ↑         ↑              ↑
    │         │              │
Restock   Restock      Restock/Returns
```

---

## API Endpoints

### Public Endpoints

#### 1. Check Stock

```http
POST /api/inventory/check
Content-Type: application/json
```

**Request:**
```json
{
  "product_id": 123,
  "quantity": 2
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "inventory": {
      "id": 1,
      "product_id": 123,
      "stock_quantity": 100,
      "reserved_quantity": 20,
      "available_quantity": 80
    },
    "can_fulfill": true,
    "available": 80
  }
}
```

#### 2. Get Inventory

```http
GET /api/inventory/:product_id
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "product_id": 123,
    "stock_quantity": 100,
    "reserved_quantity": 20,
    "available_quantity": 80,
    "status": "in_stock"
  }
}
```

#### 3. Get Inventory Logs

```http
GET /api/inventory/:product_id/logs?page=1&limit=20
```

### Admin Endpoints (Require Authentication)

#### 1. Get Inventory Summary

```http
GET /api/inventory/summary
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_products": 500,
    "in_stock_products": 450,
    "low_stock_products": 30,
    "out_of_stock_products": 20,
    "total_stock_value": 125000.00,
    "total_reserved_value": 15000.00
  }
}
```

#### 2. Get Low Stock Products

```http
GET /api/inventory/low-stock?threshold=10
```

#### 3. Get Out of Stock Products

```http
GET /api/inventory/out-of-stock
```

#### 4. Restock Product

```http
POST /api/inventory/restock
Authorization: Bearer {admin_token}
Content-Type: application/json
```

**Request:**
```json
{
  "product_id": 123,
  "quantity": 50,
  "reference_id": "RESTOCK-2026-001",
  "reason": "Weekly restock from supplier"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Product restocked successfully",
  "data": {
    "id": 1,
    "product_id": 123,
    "stock_quantity": 150,
    "reserved_quantity": 20,
    "available_quantity": 130
  }
}
```

---

## Order Integration

### Order Lifecycle Stock Flow

```
┌─────────────────┐
│  Order Created  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Reserve Stock  │  ← Decreases available_quantity
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌─────────┐ ┌──────────┐
│ Payment │ │ Payment  │
│ Success │ │ Failed   │
└────┬────┘ └────┬─────┘
     │           │
     ▼           ▼
┌─────────┐ ┌──────────┐
│ Deduct  │ │ Release  │
│ Stock   │ │ Stock    │
└─────────┘ └──────────┘
```

### Integration Points

| Order Event | Inventory Action |
|-------------|------------------|
| Order created | Reserve stock |
| Payment confirmed | Deduct stock |
| Order cancelled | Release stock |
| Order returned | Return stock |

---

## Concurrency Control

### Database Transactions

All stock operations use database transactions:

```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 1. Lock the row
    if err := tx.Set("gorm:query_option", "WITH (UPDLOCK)").
        Where("product_id = ?", productID).
        First(&inventory).Error; err != nil {
        return err
    }
    
    // 2. Check stock
    if inventory.AvailableQuantity < quantity {
        return errors.New("insufficient stock")
    }
    
    // 3. Update stock
    inventory.ReservedQuantity += quantity
    inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
    
    // 4. Save and log
    tx.Save(&inventory)
    tx.Create(&log)
    
    return nil
})
```

### Row-Level Locking (SQL Server)

```sql
-- Lock row for update
SELECT * FROM inventory WITH (UPDLOCK)
WHERE product_id = 123;

-- Update is now safe from concurrent modifications
UPDATE inventory
SET reserved_quantity = reserved_quantity + 2
WHERE product_id = 123;
```

---

## Inventory Logging

### Change Types

| Type | Description | When Used |
|------|-------------|-----------|
| `restock` | Add stock | Admin restocks product |
| `reserve` | Reserve stock | Order created |
| `deduct` | Deduct stock | Payment confirmed |
| `release` | Release stock | Order cancelled |
| `return` | Return stock | Customer return |
| `adjust` | Manual adjustment | Inventory correction |
| `damaged` | Remove damaged | Stock damaged/lost |

### Why Logs Are Important

1. **Audit Trail** - Track who changed what and when
2. **Debugging** - Understand stock discrepancies
3. **Compliance** - Meet regulatory requirements
4. **Analytics** - Analyze stock movement patterns
5. **Reconciliation** - Match physical counts with system

---

## Examples

### cURL Examples

**Check stock:**
```bash
curl -X POST http://localhost:8080/api/inventory/check \
  -H "Content-Type: application/json" \
  -d '{"product_id": 123, "quantity": 2}'
```

**Get inventory:**
```bash
curl -X GET http://localhost:8080/api/inventory/123
```

**Restock product:**
```bash
curl -X POST http://localhost:8080/api/inventory/restock \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 123, "quantity": 50, "reason": "Weekly restock"}'
```

**Get inventory logs:**
```bash
curl -X GET "http://localhost:8080/api/inventory/123/logs?page=1&limit=20"
```

**Get inventory summary:**
```bash
curl -X GET http://localhost:8080/api/inventory/summary \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### Frontend Integration (React)

```jsx
function StockChecker({ productId }) {
  const [stock, setStock] = useState(null);
  const [loading, setLoading] = useState(false);

  const checkStock = async (quantity) => {
    setLoading(true);
    const response = await fetch('/api/inventory/check', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ product_id: productId, quantity })
    });
    const result = await response.json();
    
    if (result.success) {
      setStock(result.data);
    }
    setLoading(false);
  };

  useEffect(() => {
    checkStock(1);
  }, [productId]);

  if (!stock) return <LoadingSpinner />;

  return (
    <div className="stock-status">
      {stock.can_fulfill ? (
        <span className="in-stock">
          {stock.available} in stock
        </span>
      ) : (
        <span className="out-of-stock">Out of stock</span>
      )}
    </div>
  );
}
```

---

## Implementation Summary

### Files Created

| File | Purpose |
|------|---------|
| `internal/domain/model/inventory.go` | Inventory models |
| `internal/repository/inventory_repository.go` | Database operations |
| `internal/service/inventory_service.go` | Business logic |
| `internal/handler/inventory_handler.go` | HTTP handlers |
| `api/routes_inventory.go` | API routes |

### Key Functions

**Repository:**
- `ReserveStock()` - Reserve with locking
- `DeductStock()` - Deduct after payment
- `ReleaseStock()` - Release on cancel
- `RestockProduct()` - Add stock

**Service:**
- `CheckStock()` - Check availability
- `ProcessOrderStock()` - Handle order lifecycle
- `BatchReserveStock()` - Reserve multiple items

**Handler:**
- `CheckStock()` - POST /api/inventory/check
- `GetInventory()` - GET /api/inventory/:id
- `RestockProduct()` - POST /api/inventory/restock

---

## Best Practices

1. **Always use transactions** - Never update stock outside transactions
2. **Lock rows early** - Lock before checking stock
3. **Log all changes** - Every stock change must be logged
4. **Monitor low stock** - Set up alerts for low inventory
5. **Regular audits** - Compare physical counts with system
6. **Handle failures** - Always release stock on order cancellation

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Stock shows negative | Check for missing transaction |
| Overselling occurred | Verify row locking is enabled |
| Logs missing | Ensure logging in transaction |
| Deadlocks | Reduce transaction scope |

---

✅ **System Ready**

The Inventory Management System is fully implemented with:
- Atomic stock updates
- Stock reservation system
- Comprehensive logging
- Low stock alerts
- Order integration
- Concurrency control
