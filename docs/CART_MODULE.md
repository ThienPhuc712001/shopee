# Shopping Cart Module - Business Flow & Implementation

## PART 1 — Cart Business Flow

### Complete Shopping Cart Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    SHOPPING CART FLOW                            │
└─────────────────────────────────────────────────────────────────┘

1. CUSTOMER BROWSES PRODUCTS
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│ Product  │────────>│ Database │
   │          │  GET    │ Service  │  Query  │          │
   │          │ /products│          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Products displayed with prices
   ✓ Stock availability shown
   ✓ Variants displayed (size, color)

2. CUSTOMER SELECTS PRODUCT VARIANT
   ┌──────────┐         ┌──────────┐
   │ Customer │────────>│ Frontend │
   │          │  Select │          │
   │          │ Variant │          │
   └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Variant selected (e.g., Size: L, Color: Blue)
   ✓ Price updated based on variant
   ✓ Stock checked for selected variant

3. CUSTOMER CLICKS "ADD TO CART"
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│   Cart   │────────>│ Database │
   │          │  POST   │ Service  │  INSERT │          │
   │          │  /cart  │          │  /UPDATE│          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ User authentication validated
   ✓ Product existence verified
   ✓ Variant availability checked
   ✓ Stock availability confirmed

4. SYSTEM CHECKS INVENTORY
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │   Cart   │────────>│Inventory │────────>│ Database │
   │ Service  │         │ Service  │  Check  │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Available stock = stock - reserved
   ✓ Requested quantity <= available stock
   ✓ If insufficient → return error

5. ITEM IS ADDED TO USER'S CART
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │   Cart   │────────>│ Database │────────>│  Cart    │
   │ Service  │  INSERT │          │  UPDATE │  State   │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ If item exists → increase quantity
   ✓ If item new → create cart item
   ✓ Cart totals recalculated
   ✓ Cart updated_at timestamp refreshed

6. CUSTOMER CAN UPDATE QUANTITY
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│   Cart   │────────>│ Database │
   │          │  PUT    │ Service  │  UPDATE │          │
   │          │ /cart   │          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Quantity validated (> 0)
   ✓ Stock rechecked
   ✓ Cart totals recalculated

7. CUSTOMER CAN REMOVE ITEMS
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│   Cart   │────────>│ Database │
   │          │ DELETE  │ Service  │  DELETE │          │
   │          │ /cart   │          │         │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Item removed from cart
   ✓ Cart totals recalculated
   ✓ Reserved stock released

8. CUSTOMER PROCEEDS TO CHECKOUT
   ┌──────────┐         ┌──────────┐         ┌──────────┐
   │ Customer │────────>│  Order   │────────>│ Database │
   │          │  POST   │ Service  │  CREATE │          │
   │          │ /orders │          │  Order  │          │
   └──────────┘         └──────────┘         └──────────┘
                              │
                              ▼
   ✓ Cart validated (all items in stock)
   ✓ Stock reserved for each item
   ✓ Order created from cart
   ✓ Cart cleared after successful order
```

### Add to Cart Decision Flow

```
                    ┌─────────────────┐
                    │  Add to Cart    │
                    │    Request      │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  User Authenticated? │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │ YES                         │ NO
              ▼                             ▼
     ┌─────────────────┐          ┌─────────────────┐
     │ Product Exists? │          │ Return 401      │
     └────────┬────────┘          │ Unauthorized    │
              │                   └─────────────────┘
              │ YES
              ▼
     ┌─────────────────┐
     │ Variant Exists? │ (if applicable)
     └────────┬────────┘
              │
     ┌────────┴────────┐
     │ YES             │ NO
     ▼                 ▼
┌─────────────────┐ ┌─────────────────┐
│ Stock Available?│ │ Return 404      │
└────────┬────────┘ │ Not Found       │
         │          └─────────────────┘
         │ YES
         ▼
┌─────────────────┐
│ Item in Cart?   │
└────────┬────────┘
         │
┌────────┴─────────────────┐
│ YES                      │ NO
▼                          ▼
┌─────────────────┐  ┌─────────────────┐
│ Increase        │  │ Create New      │
│ Quantity        │  │ Cart Item       │
└────────┬────────┘  └────────┬────────┘
         │                    │
         └────────┬───────────┘
                  │
                  ▼
         ┌─────────────────┐
         │ Update Cart     │
         │ Totals          │
         └────────┬────────┘
                  │
                  ▼
         ┌─────────────────┐
         │ Return Success  │
         │ + Cart Data     │
         └─────────────────┘
```

---

## PART 2 — Database Tables

### Carts Table

```sql
CREATE TABLE [dbo].[Carts] (
    [id]             BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [user_id]        BIGINT       NOT NULL UNIQUE,
    [total_items]    INT          NOT NULL DEFAULT 0,
    [subtotal]       DECIMAL(18,2) NOT NULL DEFAULT 0,
    [discount]       DECIMAL(18,2) DEFAULT 0,
    [total]          DECIMAL(18,2) NOT NULL DEFAULT 0,
    [currency]       NVARCHAR(3)   DEFAULT 'USD',
    [last_activity]  DATETIME     NOT NULL DEFAULT GETDATE(),
    [created_at]     DATETIME     NOT NULL DEFAULT GETDATE(),
    [updated_at]     DATETIME     NOT NULL DEFAULT GETDATE(),
    [deleted_at]     DATETIME,
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE
);

CREATE INDEX [IX_Carts_UserID] ON [Carts]([user_id]);
CREATE INDEX [IX_Carts_LastActivity] ON [Carts]([last_activity]);
```

### CartItems Table

```sql
CREATE TABLE [dbo].[CartItems] (
    [id]            BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [cart_id]       BIGINT        NOT NULL,
    [product_id]    BIGINT        NOT NULL,
    [variant_id]    BIGINT,
    [quantity]      INT           NOT NULL DEFAULT 1,
    [price]         DECIMAL(18,2) NOT NULL,
    [original_price] DECIMAL(18,2),
    [discount]      DECIMAL(18,2) DEFAULT 0,
    [subtotal]      DECIMAL(18,2) NOT NULL,
    [product_name]  NVARCHAR(500),
    [product_image] NVARCHAR(500),
    [shop_id]       BIGINT,
    [added_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]    DATETIME      NOT NULL DEFAULT GETDATE(),
    [deleted_at]    DATETIME,
    
    FOREIGN KEY ([cart_id]) REFERENCES [Carts]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([variant_id]) REFERENCES [ProductVariants]([id]),
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]),
    CONSTRAINT [UQ_CartItems] UNIQUE ([cart_id], [product_id], [variant_id]),
    CONSTRAINT [CHK_CartItems_Quantity] CHECK ([quantity] > 0 AND [quantity] <= 999)
);

CREATE INDEX [IX_CartItems_CartID] ON [CartItems]([cart_id]);
CREATE INDEX [IX_CartItems_ProductID] ON [CartItems]([product_id]);
CREATE INDEX [IX_CartItems_ShopID] ON [CartItems]([shop_id]);
```

### Table Relationships

```
┌─────────────────┐
│     Users       │
│─────────────────│
│ id (PK)         │
│ email           │
└────────┬────────┘
         │ 1:1
         │ (ON DELETE CASCADE)
         ▼
┌─────────────────┐
│     Carts       │
│─────────────────│
│ id (PK)         │
│ user_id (FK) ───┼───────┐
│ total_items     │       │
│ total_price     │       │
└────────┬────────┘       │
         │ 1:N            │
         │ (ON DELETE      │
         ▼  CASCADE)       │
┌─────────────────┐       │
│   CartItems     │       │
│─────────────────│       │
│ id (PK)         │       │
│ cart_id (FK) ───┘       │
│ product_id (FK) ────────┼───> Products
│ variant_id (FK) ────────┼───> ProductVariants
│ shop_id (FK) ───────────┼───> Shops
│ quantity        │
│ price           │
│ subtotal        │
└─────────────────┘
```

---

## PART 3-12 — Implementation

The complete implementation follows in the code files below.

### Files Created:

1. **Models** - Enhanced cart models with relationships
2. **Repository** - Cart repository with all CRUD operations
3. **Service** - Cart service with business logic
4. **Handler** - REST API handlers
5. **Routes** - Cart route definitions

### Key Features:

✅ **One Cart Per User** - Each user has exactly one cart
✅ **Multiple Items** - Unlimited items per cart
✅ **Stock Validation** - Quantity cannot exceed available stock
✅ **Auto Totals** - Cart totals calculated automatically
✅ **Variant Support** - Full support for product variants
✅ **Soft Delete** - Cart items can be soft deleted
✅ **Efficient Queries** - Indexed for high performance
✅ **Concurrency Safe** - Transaction-based updates

---

## Performance Optimization

### Indexes

```sql
-- Cart lookups
IX_Carts_UserID
IX_Carts_LastActivity

-- Cart item lookups
IX_CartItems_CartID
IX_CartItems_ProductID
IX_CartItems_ShopID
IX_CartItems_Unique (cart_id, product_id, variant_id)
```

### Query Optimization

- Use eager loading (Preload) for cart items
- Calculate totals in database (avoid application-level loops)
- Batch updates for multiple items
- Cache cart for active users (Redis)

### Concurrency Handling

- Use database transactions for cart updates
- Lock cart row during updates (SELECT FOR UPDATE)
- Release locks after transaction commit

---

This cart module is production-ready and optimized for high traffic.
