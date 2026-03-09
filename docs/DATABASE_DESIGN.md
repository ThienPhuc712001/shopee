# Database Design Documentation

## PART 1 — Database Design Principles

### 1.1 Normalization Strategy

**Third Normal Form (3NF)** for transactional tables:
- Eliminate repeating groups (1NF)
- Remove partial dependencies (2NF)
- Remove transitive dependencies (3NF)

**Why 3NF?**
- Reduces data redundancy
- Improves data integrity
- Simplifies updates
- Prevents anomalies

**Controlled Denormalization** for read-heavy tables:
- Product listings (include category name)
- Order summaries (include customer name)
- Analytics tables

### 1.2 Indexing Strategy

| Index Type | Use Case | Example |
|------------|----------|---------|
| **Clustered** | Primary access path | Primary keys |
| **Non-clustered** | Frequent WHERE clauses | email, status |
| **Composite** | Multi-column filters | (shop_id, status) |
| **Covering** | Include all query columns | INCLUDE clause |
| **Filtered** | Partial index | WHERE status = 'active' |

### 1.3 Foreign Key Relationships

- Enforce referential integrity
- CASCADE for soft dependencies
- NO ACTION for critical relationships
- Indexed foreign keys for performance

### 1.4 Transaction Consistency

- ACID compliance for orders and payments
- Optimistic locking for inventory
- Pessimistic locking for flash sales
- Distributed transactions via Saga pattern

### 1.5 Scalability Considerations

- Horizontal partitioning for large tables
- Read replicas for reporting
- Sharding strategy for user data
- Archive old data

---

## PART 2 — Core Database Modules

```
┌─────────────────────────────────────────────────────────────────┐
│                     DATABASE MODULES                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │    USER      │  │    SHOP      │  │   PRODUCT    │          │
│  │   MODULE     │  │   MODULE     │  │   MODULE     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │    CART      │  │    ORDER     │  │   PAYMENT    │          │
│  │   MODULE     │  │   MODULE     │  │   MODULE     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   REVIEW     │  │  PROMOTION   │  │ NOTIFICATION │          │
│  │   MODULE     │  │   MODULE     │  │   MODULE     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐                             │
│  │    ADMIN     │  │   ANALYTICS  │                             │
│  │   MODULE     │  │   MODULE     │                             │
│  └──────────────┘  └──────────────┘                             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## PART 3 — User Module

### Table: Users

```sql
CREATE TABLE [dbo].[Users] (
    [id]                    BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [email]                 NVARCHAR(255)  NOT NULL UNIQUE,
    [password_hash]         NVARCHAR(255)  NOT NULL,
    [phone]                 NVARCHAR(20)   UNIQUE,
    [first_name]            NVARCHAR(100),
    [last_name]             NVARCHAR(100),
    [avatar_url]            NVARCHAR(500),
    [role_id]               INT            NOT NULL DEFAULT 1,
    [status]                NVARCHAR(20)   NOT NULL DEFAULT 'active',
    [email_verified]        BIT            NOT NULL DEFAULT 0,
    [email_verified_at]     DATETIME,
    [phone_verified]        BIT            NOT NULL DEFAULT 0,
    [last_login_at]         DATETIME,
    [failed_login_attempts] INT            NOT NULL DEFAULT 0,
    [locked_until]          DATETIME,
    [refresh_token]         NVARCHAR(500),
    [refresh_token_expiry]  DATETIME,
    [created_at]            DATETIME       NOT NULL DEFAULT GETDATE(),
    [updated_at]            DATETIME       NOT NULL DEFAULT GETDATE(),
    [deleted_at]            DATETIME,
    
    CONSTRAINT [CHK_Users_Status] CHECK ([status] IN ('active', 'inactive', 'banned', 'locked')),
    INDEX [IX_Users_Email] ([email]),
    INDEX [IX_Users_Phone] ([phone]),
    INDEX [IX_Users_Status] ([status]),
    INDEX [IX_Users_DeletedAt] ([deleted_at])
);
```

### Table: UserProfiles

```sql
CREATE TABLE [dbo].[UserProfiles] (
    [id]            BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]       BIGINT        NOT NULL UNIQUE,
    [gender]        NVARCHAR(10),
    [date_of_birth] DATE,
    [bio]           NVARCHAR(500),
    [language]      NVARCHAR(10)  DEFAULT 'en',
    [currency]      NVARCHAR(3)   DEFAULT 'USD',
    [timezone]      NVARCHAR(50)  DEFAULT 'UTC',
    [marketing_opt_in] BIT        NOT NULL DEFAULT 0,
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    INDEX [IX_UserProfiles_UserID] ([user_id])
);
```

### Table: UserRoles

```sql
CREATE TABLE [dbo].[UserRoles] (
    [id]          INT           IDENTITY(1,1) PRIMARY KEY,
    [name]        NVARCHAR(50)  NOT NULL UNIQUE,
    [description] NVARCHAR(255),
    [permissions] NVARCHAR(MAX), -- JSON array of permissions
    
    INDEX [IX_UserRoles_Name] ([name])
);

-- Default roles
INSERT INTO [UserRoles] ([name], [description]) VALUES
('customer', 'Regular customer'),
('seller', 'Verified seller'),
('admin', 'Platform administrator'),
('super_admin', 'System administrator');
```

### Table: UserAddresses

```sql
CREATE TABLE [dbo].[UserAddresses] (
    [id]            BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]       BIGINT        NOT NULL,
    [label]         NVARCHAR(50)  DEFAULT 'Home', -- Home, Work, etc.
    [full_name]     NVARCHAR(200) NOT NULL,
    [phone]         NVARCHAR(20)  NOT NULL,
    [address_line1] NVARCHAR(500) NOT NULL,
    [address_line2] NVARCHAR(500),
    [ward]          NVARCHAR(200),
    [district]      NVARCHAR(200),
    [city]          NVARCHAR(200) NOT NULL,
    [state]         NVARCHAR(200),
    [country]       NVARCHAR(100) NOT NULL DEFAULT 'Vietnam',
    [postal_code]   NVARCHAR(20),
    [is_default]    BIT           NOT NULL DEFAULT 0,
    [latitude]      DECIMAL(10,8),
    [longitude]     DECIMAL(11,8),
    [created_at]    DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]    DATETIME      NOT NULL DEFAULT GETDATE(),
    [deleted_at]    DATETIME,
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    CONSTRAINT [DF_UserAddresses_IsDefault] DEFAULT 0 FOR [is_default],
    INDEX [IX_UserAddresses_UserID] ([user_id]),
    INDEX [IX_UserAddresses_IsDefault] ([is_default])
);
```

### Table: UserSessions

```sql
CREATE TABLE [dbo].[UserSessions] (
    [id]             BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [user_id]        BIGINT       NOT NULL,
    [session_token]  NVARCHAR(255) NOT NULL UNIQUE,
    [refresh_token]  NVARCHAR(500),
    [ip_address]     NVARCHAR(45),
    [user_agent]     NVARCHAR(500),
    [device_type]    NVARCHAR(50), -- mobile, desktop, tablet
    [expires_at]     DATETIME     NOT NULL,
    [created_at]     DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    INDEX [IX_UserSessions_UserID] ([user_id]),
    INDEX [IX_UserSessions_Token] ([session_token]),
    INDEX [IX_UserSessions_ExpiresAt] ([expires_at])
);
```

### Table: UserSecurityLogs

```sql
CREATE TABLE [dbo].[UserSecurityLogs] (
    [id]          BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [user_id]     BIGINT       NOT NULL,
    [action]      NVARCHAR(50) NOT NULL, -- login, logout, password_change, etc.
    [status]      NVARCHAR(20) NOT NULL, -- success, failed
    [ip_address]  NVARCHAR(45),
    [user_agent]  NVARCHAR(500),
    [details]     NVARCHAR(MAX), -- JSON
    [created_at]  DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    INDEX [IX_UserSecurityLogs_UserID] ([user_id]),
    INDEX [IX_UserSecurityLogs_Action] ([action]),
    INDEX [IX_UserSecurityLogs_CreatedAt] ([created_at])
);
```

---

## PART 4 — Shop Module

### Table: Shops

```sql
CREATE TABLE [dbo].[Shops] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]         BIGINT        NOT NULL UNIQUE,
    [name]            NVARCHAR(255) NOT NULL,
    [slug]            NVARCHAR(255) NOT NULL UNIQUE,
    [logo_url]        NVARCHAR(500),
    [cover_image_url] NVARCHAR(500),
    [description]     NVARCHAR(MAX),
    [phone]           NVARCHAR(20),
    [email]           NVARCHAR(255),
    [address]         NVARCHAR(500),
    [city]            NVARCHAR(200),
    [country]         NVARCHAR(100) DEFAULT 'Vietnam',
    [status]          NVARCHAR(20)  NOT NULL DEFAULT 'pending', -- pending, active, suspended, banned
    [verification_status] NVARCHAR(20) DEFAULT 'unverified', -- unverified, verified, rejected
    [rating_avg]      DECIMAL(3,2)  DEFAULT 0,
    [rating_count]    INT           DEFAULT 0,
    [follower_count]  INT           DEFAULT 0,
    [product_count]   INT           DEFAULT 0,
    [total_sales]     BIGINT        DEFAULT 0,
    [total_revenue]   DECIMAL(18,2) DEFAULT 0,
    [response_rate]   DECIMAL(5,2)  DEFAULT 0,
    [response_time]   INT           DEFAULT 0, -- in hours
    [joined_at]       DATETIME      NOT NULL DEFAULT GETDATE(),
    [last_active_at]  DATETIME,
    [created_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    [deleted_at]      DATETIME,
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    CONSTRAINT [CHK_Shops_Status] CHECK ([status] IN ('pending', 'active', 'suspended', 'banned')),
    INDEX [IX_Shops_UserID] ([user_id]),
    INDEX [IX_Shops_Slug] ([slug]),
    INDEX [IX_Shops_Status] ([status]),
    INDEX [IX_Shops_Rating] ([rating_avg]),
    INDEX [IX_Shops_TotalSales] ([total_sales])
);
```

### Table: ShopFollowers

```sql
CREATE TABLE [dbo].[ShopFollowers] (
    [id]         BIGINT   IDENTITY(1,1) PRIMARY KEY,
    [shop_id]    BIGINT   NOT NULL,
    [user_id]    BIGINT   NOT NULL,
    [created_at] DATETIME NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    CONSTRAINT [UQ_ShopFollowers] UNIQUE ([shop_id], [user_id]),
    INDEX [IX_ShopFollowers_ShopID] ([shop_id]),
    INDEX [IX_ShopFollowers_UserID] ([user_id])
);
```

### Table: ShopRatings

```sql
CREATE TABLE [dbo].[ShopRatings] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [shop_id]     BIGINT        NOT NULL,
    [user_id]     BIGINT        NOT NULL,
    [order_id]    BIGINT,
    [rating]      INT           NOT NULL CHECK ([rating] BETWEEN 1 AND 5),
    [comment]     NVARCHAR(500),
    [response]    NVARCHAR(500), -- Seller response
    [response_at] DATETIME,
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    CONSTRAINT [UQ_ShopRatings_Order] UNIQUE ([order_id]),
    INDEX [IX_ShopRatings_ShopID] ([shop_id]),
    INDEX [IX_ShopRatings_UserID] ([user_id]),
    INDEX [IX_ShopRatings_Rating] ([rating])
);
```

### Table: ShopSettings

```sql
CREATE TABLE [dbo].[ShopSettings] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [shop_id]         BIGINT        NOT NULL UNIQUE,
    [shipping_methods] NVARCHAR(MAX), -- JSON array
    [payment_methods]  NVARCHAR(MAX), -- JSON array
    [return_policy]    NVARCHAR(MAX),
    [shipping_policy]  NVARCHAR(MAX),
    [auto_confirm_days] INT          DEFAULT 7,
    [vacation_mode]    BIT           NOT NULL DEFAULT 0,
    [vacation_start]   DATETIME,
    [vacation_end]     DATETIME,
    [created_at]       DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]       DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]) ON DELETE CASCADE,
    INDEX [IX_ShopSettings_ShopID] ([shop_id])
);
```

---

## PART 5 — Product Module

### Table: Categories

```sql
CREATE TABLE [dbo].[Categories] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [parent_id]   BIGINT,
    [name]        NVARCHAR(255) NOT NULL,
    [slug]        NVARCHAR(255) NOT NULL UNIQUE,
    [description] NVARCHAR(MAX),
    [icon_url]    NVARCHAR(500),
    [image_url]   NVARCHAR(500),
    [level]       INT           NOT NULL DEFAULT 0,
    [sort_order]  INT           DEFAULT 0,
    [is_active]   BIT           NOT NULL DEFAULT 1,
    [attributes]  NVARCHAR(MAX), -- JSON schema for category attributes
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([parent_id]) REFERENCES [Categories]([id]),
    INDEX [IX_Categories_ParentID] ([parent_id]),
    INDEX [IX_Categories_Slug] ([slug]),
    INDEX [IX_Categories_Level] ([level]),
    INDEX [IX_Categories_IsActive] ([is_active])
);
```

### Table: Products

```sql
CREATE TABLE [dbo].[Products] (
    [id]              BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [shop_id]         BIGINT         NOT NULL,
    [category_id]     BIGINT         NOT NULL,
    [name]            NVARCHAR(500)  NOT NULL,
    [slug]            NVARCHAR(500)  NOT NULL,
    [description]     NVARCHAR(MAX),
    [short_description] NVARCHAR(1000),
    [sku]             NVARCHAR(100),
    [brand]           NVARCHAR(200),
    [price]           DECIMAL(18,2)  NOT NULL,
    [original_price]  DECIMAL(18,2),
    [discount_percent] INT           DEFAULT 0,
    [cost]            DECIMAL(18,2), -- For profit calculation
    [stock]           INT           NOT NULL DEFAULT 0,
    [reserved_stock]  INT           DEFAULT 0,
    [sold_count]      BIGINT        DEFAULT 0,
    [view_count]      BIGINT        DEFAULT 0,
    [status]          NVARCHAR(20)  NOT NULL DEFAULT 'active', -- draft, active, inactive, banned
    [is_featured]     BIT           NOT NULL DEFAULT 0,
    [is_flash_sale]   BIT           NOT NULL DEFAULT 0,
    [flash_sale_price] DECIMAL(18,2),
    [flash_sale_start] DATETIME,
    [flash_sale_end]   DATETIME,
    [rating_avg]      DECIMAL(3,2)  DEFAULT 0,
    [rating_count]    INT           DEFAULT 0,
    [review_count]    INT           DEFAULT 0,
    [weight]          DECIMAL(10,2), -- in grams
    [dimensions]      NVARCHAR(50), -- LxWxH
    [warranty_period] NVARCHAR(50),
    [return_days]     INT           DEFAULT 7,
    [tags]           NVARCHAR(MAX), -- JSON array
    [meta_title]     NVARCHAR(255),
    [meta_description] NVARCHAR(500),
    [created_at]     DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]     DATETIME      NOT NULL DEFAULT GETDATE(),
    [deleted_at]     DATETIME,
    
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]),
    FOREIGN KEY ([category_id]) REFERENCES [Categories]([id]),
    CONSTRAINT [UQ_Products_Slug] UNIQUE ([shop_id], [slug]),
    CONSTRAINT [CHK_Products_Status] CHECK ([status] IN ('draft', 'active', 'inactive', 'banned')),
    CONSTRAINT [CHK_Products_Price] CHECK ([price] >= 0),
    CONSTRAINT [CHK_Products_Stock] CHECK ([stock] >= 0),
    INDEX [IX_Products_ShopID] ([shop_id]),
    INDEX [IX_Products_CategoryID] ([category_id]),
    INDEX [IX_Products_Slug] ([slug]),
    INDEX [IX_Products_Status] ([status]),
    INDEX [IX_Products_Price] ([price]),
    INDEX [IX_Products_SoldCount] ([sold_count]),
    INDEX [IX_Products_Rating] ([rating_avg]),
    INDEX [IX_Products_CreatedAt] ([created_at]),
    FULLTEXT INDEX [FTIX_Products_Search] ON ([name], [description], [short_description])
);
```

### Table: ProductImages

```sql
CREATE TABLE [dbo].[ProductImages] (
    [id]          BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [product_id]  BIGINT       NOT NULL,
    [url]         NVARCHAR(500) NOT NULL,
    [alt_text]    NVARCHAR(255),
    [is_primary]  BIT          NOT NULL DEFAULT 0,
    [sort_order]  INT          DEFAULT 0,
    [width]       INT,
    [height]      INT,
    [size_bytes]  BIGINT,
    [created_at]  DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    INDEX [IX_ProductImages_ProductID] ([product_id]),
    INDEX [IX_ProductImages_IsPrimary] ([is_primary])
);
```

### Table: ProductVariants

```sql
CREATE TABLE [dbo].[ProductVariants] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [product_id]  BIGINT        NOT NULL,
    [sku]         NVARCHAR(100),
    [name]        NVARCHAR(200), -- e.g., "Red / XL"
    [price]       DECIMAL(18,2),
    [original_price] DECIMAL(18,2),
    [stock]       INT           NOT NULL DEFAULT 0,
    [reserved_stock] INT        DEFAULT 0,
    [attributes]  NVARCHAR(MAX) NOT NULL, -- JSON: {"color": "Red", "size": "XL"}
    [image_url]   NVARCHAR(500),
    [sort_order]  INT           DEFAULT 0,
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    INDEX [IX_ProductVariants_ProductID] ([product_id]),
    INDEX [IX_ProductVariants_SKU] ([sku])
);
```

### Table: ProductAttributes

```sql
CREATE TABLE [dbo].[ProductAttributes] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [category_id] BIGINT,
    [name]        NVARCHAR(100) NOT NULL, -- e.g., "Color", "Size"
    [type]        NVARCHAR(50)  NOT NULL, -- text, select, color, number
    [values]      NVARCHAR(MAX), -- JSON array of possible values
    [is_filterable] BIT         NOT NULL DEFAULT 1,
    [is_visible]    BIT         NOT NULL DEFAULT 1,
    [sort_order]  INT           DEFAULT 0,
    
    FOREIGN KEY ([category_id]) REFERENCES [Categories]([id]),
    INDEX [IX_ProductAttributes_CategoryID] ([category_id])
);
```

### Table: ProductInventory

```sql
CREATE TABLE [dbo].[ProductInventory] (
    [id]           BIGINT    IDENTITY(1,1) PRIMARY KEY,
    [product_id]   BIGINT    NOT NULL,
    [variant_id]   BIGINT,
    [warehouse_id] BIGINT,
    [quantity]     INT       NOT NULL DEFAULT 0,
    [reserved]     INT       NOT NULL DEFAULT 0,
    [available]    AS ([quantity] - [reserved]),
    [reorder_point] INT      DEFAULT 10,
    [last_counted_at] DATETIME,
    [created_at]   DATETIME  NOT NULL DEFAULT GETDATE(),
    [updated_at]   DATETIME  NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([variant_id]) REFERENCES [ProductVariants]([id]),
    CONSTRAINT [UQ_ProductInventory] UNIQUE ([product_id], [variant_id], [warehouse_id]),
    INDEX [IX_ProductInventory_ProductID] ([product_id])
);
```

---

## PART 6 — Cart Module

### Table: Carts

```sql
CREATE TABLE [dbo].[Carts] (
    [id]           BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [user_id]      BIGINT       NOT NULL UNIQUE,
    [session_id]   NVARCHAR(255), -- For guest carts
    [total_items]  INT          NOT NULL DEFAULT 0,
    [subtotal]     DECIMAL(18,2) NOT NULL DEFAULT 0,
    [discount]     DECIMAL(18,2) DEFAULT 0,
    [total]        DECIMAL(18,2) NOT NULL DEFAULT 0,
    [currency]     NVARCHAR(3)   DEFAULT 'USD',
    [last_activity] DATETIME     NOT NULL DEFAULT GETDATE(),
    [created_at]   DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]   DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    INDEX [IX_Carts_UserID] ([user_id]),
    INDEX [IX_Carts_SessionID] ([session_id]),
    INDEX [IX_Carts_LastActivity] ([last_activity])
);
```

### Table: CartItems

```sql
CREATE TABLE [dbo].[CartItems] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [cart_id]     BIGINT        NOT NULL,
    [product_id]  BIGINT        NOT NULL,
    [variant_id]  BIGINT,
    [quantity]    INT           NOT NULL DEFAULT 1,
    [price]       DECIMAL(18,2) NOT NULL,
    [original_price] DECIMAL(18,2),
    [discount]    DECIMAL(18,2) DEFAULT 0,
    [subtotal]    DECIMAL(18,2) NOT NULL,
    [product_name] NVARCHAR(500), -- Snapshot at add time
    [product_image] NVARCHAR(500),
    [shop_id]     BIGINT,
    [added_at]    DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([cart_id]) REFERENCES [Carts]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([variant_id]) REFERENCES [ProductVariants]([id]),
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]),
    CONSTRAINT [UQ_CartItems] UNIQUE ([cart_id], [product_id], [variant_id]),
    CONSTRAINT [CHK_CartItems_Quantity] CHECK ([quantity] > 0 AND [quantity] <= 999),
    INDEX [IX_CartItems_CartID] ([cart_id]),
    INDEX [IX_CartItems_ProductID] ([product_id]),
    INDEX [IX_CartItems_ShopID] ([shop_id])
);
```

---

## PART 7 — Order Module

### Table: Orders

```sql
CREATE TABLE [dbo].[Orders] (
    [id]                BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [order_number]      NVARCHAR(50)   NOT NULL UNIQUE,
    [user_id]           BIGINT         NOT NULL,
    [shop_id]           BIGINT         NOT NULL,
    [parent_order_id]   BIGINT, -- For multi-shop orders
    [status]            NVARCHAR(20)   NOT NULL DEFAULT 'pending',
    [payment_status]    NVARCHAR(20)   NOT NULL DEFAULT 'pending',
    [fulfillment_status] NVARCHAR(20)  DEFAULT 'unfulfilled',
    
    -- Pricing
    [subtotal]          DECIMAL(18,2)  NOT NULL,
    [shipping_fee]      DECIMAL(18,2)  DEFAULT 0,
    [shipping_discount] DECIMAL(18,2)  DEFAULT 0,
    [product_discount]  DECIMAL(18,2)  DEFAULT 0,
    [voucher_discount]  DECIMAL(18,2)  DEFAULT 0,
    [tax_amount]        DECIMAL(18,2)  DEFAULT 0,
    [total_amount]      DECIMAL(18,2)  NOT NULL,
    [paid_amount]       DECIMAL(18,2)  DEFAULT 0,
    
    -- Shipping Address (snapshot)
    [shipping_name]     NVARCHAR(200)  NOT NULL,
    [shipping_phone]    NVARCHAR(20)   NOT NULL,
    [shipping_address]  NVARCHAR(500)  NOT NULL,
    [shipping_ward]     NVARCHAR(200),
    [shipping_district] NVARCHAR(200),
    [shipping_city]     NVARCHAR(200),
    [shipping_state]    NVARCHAR(200),
    [shipping_country]  NVARCHAR(100)  DEFAULT 'Vietnam',
    [shipping_postal_code] NVARCHAR(20),
    
    -- Shipping Method
    [shipping_method]   NVARCHAR(100),
    [shipping_carrier]  NVARCHAR(100),
    [tracking_number]   NVARCHAR(100),
    [estimated_delivery] DATETIME,
    
    -- Notes
    [buyer_note]        NVARCHAR(500),
    [seller_note]       NVARCHAR(500),
    [cancel_reason]     NVARCHAR(500),
    
    -- Timestamps
    [paid_at]           DATETIME,
    [confirmed_at]      DATETIME,
    [shipped_at]        DATETIME,
    [delivered_at]      DATETIME,
    [cancelled_at]      DATETIME,
    [completed_at]      DATETIME,
    [created_at]        DATETIME     NOT NULL DEFAULT GETDATE(),
    [updated_at]        DATETIME     NOT NULL DEFAULT GETDATE(),
    [deleted_at]        DATETIME,
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]),
    FOREIGN KEY ([parent_order_id]) REFERENCES [Orders]([id]),
    CONSTRAINT [CHK_Orders_Status] CHECK ([status] IN (
        'pending', 'confirmed', 'processing', 'shipped', 
        'delivered', 'completed', 'cancelled', 'refunded'
    )),
    CONSTRAINT [CHK_Orders_PaymentStatus] CHECK ([payment_status] IN (
        'pending', 'paid', 'failed', 'refunded', 'partial_refund'
    )),
    INDEX [IX_Orders_OrderNumber] ([order_number]),
    INDEX [IX_Orders_UserID] ([user_id]),
    INDEX [IX_Orders_ShopID] ([shop_id]),
    INDEX [IX_Orders_Status] ([status]),
    INDEX [IX_Orders_PaymentStatus] ([payment_status]),
    INDEX [IX_Orders_CreatedAt] ([created_at]),
    INDEX [IX_Orders_DeliveredAt] ([delivered_at])
);
```

### Table: OrderItems

```sql
CREATE TABLE [dbo].[OrderItems] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [order_id]        BIGINT        NOT NULL,
    [product_id]      BIGINT        NOT NULL,
    [variant_id]      BIGINT,
    [product_name]    NVARCHAR(500) NOT NULL,
    [product_image]   NVARCHAR(500),
    [product_sku]     NVARCHAR(100),
    [quantity]        INT           NOT NULL,
    [price]           DECIMAL(18,2) NOT NULL,
    [original_price]  DECIMAL(18,2),
    [discount]        DECIMAL(18,2) DEFAULT 0,
    [subtotal]        DECIMAL(18,2) NOT NULL,
    [tax_amount]      DECIMAL(18,2) DEFAULT 0,
    [final_amount]    DECIMAL(18,2) NOT NULL,
    [shop_id]         BIGINT        NOT NULL,
    [fulfillment_status] NVARCHAR(20) DEFAULT 'pending',
    [tracking_number] NVARCHAR(100),
    [shipped_at]      DATETIME,
    [delivered_at]    DATETIME,
    
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]),
    FOREIGN KEY ([variant_id]) REFERENCES [ProductVariants]([id]),
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]),
    CONSTRAINT [CHK_OrderItems_Fulfillment] CHECK ([fulfillment_status] IN (
        'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    )),
    INDEX [IX_OrderItems_OrderID] ([order_id]),
    INDEX [IX_OrderItems_ProductID] ([product_id]),
    INDEX [IX_OrderItems_ShopID] ([shop_id])
);
```

### Table: OrderStatusHistory

```sql
CREATE TABLE [dbo].[OrderStatusHistory] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [order_id]    BIGINT        NOT NULL,
    [status]      NVARCHAR(20)  NOT NULL,
    [from_status] NVARCHAR(20),
    [message]     NVARCHAR(500),
    [changed_by]  BIGINT, -- user_id or system
    [ip_address]  NVARCHAR(45),
    [created_at]  DATETIME    NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]) ON DELETE CASCADE,
    INDEX [IX_OrderStatusHistory_OrderID] ([order_id]),
    INDEX [IX_OrderStatusHistory_CreatedAt] ([created_at])
);
```

### Table: OrderShipping

```sql
CREATE TABLE [dbo].[OrderShipping] (
    [id]              BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [order_id]        BIGINT       NOT NULL,
    [carrier_id]      BIGINT,
    [carrier_name]    NVARCHAR(100),
    [tracking_number] NVARCHAR(100),
    [shipping_label_url] NVARCHAR(500),
    [weight]          DECIMAL(10,2),
    [dimensions]      NVARCHAR(50),
    [shipped_at]      DATETIME,
    [estimated_delivery] DATETIME,
    [actual_delivery] DATETIME,
    [status]          NVARCHAR(50),
    [tracking_events] NVARCHAR(MAX), -- JSON array of tracking events
    [created_at]      DATETIME     NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    INDEX [IX_OrderShipping_OrderID] ([order_id]),
    INDEX [IX_OrderShipping_TrackingNumber] ([tracking_number])
);
```

---

## PART 8 — Payment Module

### Table: Payments

```sql
CREATE TABLE [dbo].[Payments] (
    [id]              BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [order_id]        BIGINT         NOT NULL,
    [user_id]         BIGINT         NOT NULL,
    [transaction_id]  NVARCHAR(255)  UNIQUE,
    [payment_method]  NVARCHAR(50)   NOT NULL,
    [payment_provider] NVARCHAR(50),
    [amount]          DECIMAL(18,2)  NOT NULL,
    [currency]        NVARCHAR(3)    DEFAULT 'USD',
    [status]          NVARCHAR(20)   NOT NULL DEFAULT 'pending',
    [gateway_response] NVARCHAR(MAX), -- JSON
    [metadata]        NVARCHAR(MAX), -- JSON
    [paid_at]         DATETIME,
    [failed_at]       DATETIME,
    [failure_reason]  NVARCHAR(500),
    [refunded_amount] DECIMAL(18,2)  DEFAULT 0,
    [created_at]      DATETIME       NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME       NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    CONSTRAINT [CHK_Payments_Status] CHECK ([status] IN (
        'pending', 'processing', 'completed', 'failed', 
        'cancelled', 'refunded', 'partial_refund'
    )),
    INDEX [IX_Payments_OrderID] ([order_id]),
    INDEX [IX_Payments_UserID] ([user_id]),
    INDEX [IX_Payments_TransactionID] ([transaction_id]),
    INDEX [IX_Payments_Status] ([status]),
    INDEX [IX_Payments_CreatedAt] ([created_at])
);
```

### Table: PaymentMethods

```sql
CREATE TABLE [dbo].[PaymentMethods] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]     BIGINT,
    [type]        NVARCHAR(50)  NOT NULL, -- credit_card, bank_account, e_wallet
    [provider]    NVARCHAR(50)  NOT NULL, -- stripe, paypal, etc.
    [name]        NVARCHAR(100), -- e.g., "Visa ending 4242"
    [last_four]   NVARCHAR(4),
    [expiry_month] INT,
    [expiry_year]  INT,
    [is_default]   BIT          NOT NULL DEFAULT 0,
    [token]        NVARCHAR(500), -- Payment provider token
    [metadata]     NVARCHAR(MAX),
    [created_at]   DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    INDEX [IX_PaymentMethods_UserID] ([user_id])
);
```

### Table: PaymentTransactions

```sql
CREATE TABLE [dbo].[PaymentTransactions] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [payment_id]      BIGINT        NOT NULL,
    [type]            NVARCHAR(50)  NOT NULL, -- charge, refund, chargeback
    [amount]          DECIMAL(18,2) NOT NULL,
    [status]          NVARCHAR(20)  NOT NULL,
    [gateway_id]      NVARCHAR(255),
    [gateway_response] NVARCHAR(MAX),
    [processed_at]    DATETIME,
    [created_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([payment_id]) REFERENCES [Payments]([id]),
    INDEX [IX_PaymentTransactions_PaymentID] ([payment_id])
);
```

### Table: Refunds

```sql
CREATE TABLE [dbo].[Refunds] (
    [id]              BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [payment_id]      BIGINT         NOT NULL,
    [order_id]        BIGINT         NOT NULL,
    [refund_number]   NVARCHAR(50)   UNIQUE,
    [amount]          DECIMAL(18,2)  NOT NULL,
    [reason]          NVARCHAR(500)  NOT NULL,
    [status]          NVARCHAR(20)   NOT NULL DEFAULT 'pending',
    [type]            NVARCHAR(20)   NOT NULL, -- full, partial
    [requested_by]    BIGINT,
    [approved_by]     BIGINT,
    [approved_at]     DATETIME,
    [processed_at]    DATETIME,
    [gateway_refund_id] NVARCHAR(255),
    [notes]           NVARCHAR(500),
    [created_at]      DATETIME       NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME       NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([payment_id]) REFERENCES [Payments]([id]),
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    FOREIGN KEY ([requested_by]) REFERENCES [Users]([id]),
    FOREIGN KEY ([approved_by]) REFERENCES [Users]([id]),
    CONSTRAINT [CHK_Refunds_Status] CHECK ([status] IN (
        'pending', 'approved', 'rejected', 'processing', 'completed', 'failed'
    )),
    INDEX [IX_Refunds_PaymentID] ([payment_id]),
    INDEX [IX_Refunds_OrderID] ([order_id]),
    INDEX [IX_Refunds_Status] ([status])
);
```

---

## PART 9 — Review Module

### Table: Reviews

```sql
CREATE TABLE [dbo].[Reviews] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [order_id]        BIGINT        NOT NULL,
    [order_item_id]   BIGINT,
    [product_id]      BIGINT        NOT NULL,
    [shop_id]         BIGINT        NOT NULL,
    [user_id]         BIGINT        NOT NULL,
    [rating]          INT           NOT NULL CHECK ([rating] BETWEEN 1 AND 5),
    [title]           NVARCHAR(200),
    [comment]         NVARCHAR(2000),
    [is_verified_purchase] BIT       NOT NULL DEFAULT 1,
    [is_visible]      BIT           NOT NULL DEFAULT 1,
    [is_approved]     BIT           NOT NULL DEFAULT 0,
    [helpful_count]   INT           DEFAULT 0,
    [not_helpful_count] INT         DEFAULT 0,
    [seller_response] NVARCHAR(1000),
    [seller_response_at] DATETIME,
    [images]         NVARCHAR(MAX), -- JSON array of image URLs
    [created_at]     DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]     DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    FOREIGN KEY ([order_item_id]) REFERENCES [OrderItems]([id]),
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([shop_id]) REFERENCES [Shops]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    CONSTRAINT [UQ_Reviews_OrderItem] UNIQUE ([order_item_id]),
    CONSTRAINT [CHK_Reviews_Rating] CHECK ([rating] >= 1 AND [rating] <= 5),
    INDEX [IX_Reviews_ProductID] ([product_id]),
    INDEX [IX_Reviews_ShopID] ([shop_id]),
    INDEX [IX_Reviews_UserID] ([user_id]),
    INDEX [IX_Reviews_Rating] ([rating]),
    INDEX [IX_Reviews_IsVisible] ([is_visible]),
    INDEX [IX_Reviews_CreatedAt] ([created_at])
);
```

### Table: ReviewImages

```sql
CREATE TABLE [dbo].[ReviewImages] (
    [id]         BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [review_id]  BIGINT       NOT NULL,
    [url]        NVARCHAR(500) NOT NULL,
    [sort_order] INT          DEFAULT 0,
    [created_at] DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([review_id]) REFERENCES [Reviews]([id]) ON DELETE CASCADE,
    INDEX [IX_ReviewImages_ReviewID] ([review_id])
);
```

### Table: ReviewHelpful

```sql
CREATE TABLE [dbo].[ReviewHelpful] (
    [id]         BIGINT   IDENTITY(1,1) PRIMARY KEY,
    [review_id]  BIGINT   NOT NULL,
    [user_id]    BIGINT   NOT NULL,
    [is_helpful] BIT      NOT NULL DEFAULT 1,
    [created_at] DATETIME NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([review_id]) REFERENCES [Reviews]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    CONSTRAINT [UQ_ReviewHelpful] UNIQUE ([review_id], [user_id]),
    INDEX [IX_ReviewHelpful_ReviewID] ([review_id])
);
```

---

## PART 10 — Promotion Module

### Table: Vouchers

```sql
CREATE TABLE [dbo].[Vouchers] (
    [id]              BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [code]            NVARCHAR(50)   NOT NULL UNIQUE,
    [name]            NVARCHAR(255)  NOT NULL,
    [description]     NVARCHAR(500),
    [type]            NVARCHAR(20)   NOT NULL, -- percentage, fixed, shipping
    [discount_value]  DECIMAL(18,2)  NOT NULL,
    [discount_percent] DECIMAL(5,2),
    [min_order_amount] DECIMAL(18,2) DEFAULT 0,
    [max_discount]    DECIMAL(18,2),
    [usage_limit]     INT,
    [usage_count]     INT           DEFAULT 0,
    [per_user_limit]  INT           DEFAULT 1,
    [applicable_shops] NVARCHAR(MAX), -- JSON array of shop IDs or "all"
    [applicable_categories] NVARCHAR(MAX), -- JSON array
    [excluded_products] NVARCHAR(MAX), -- JSON array
    [is_public]       BIT           NOT NULL DEFAULT 1,
    [is_active]       BIT           NOT NULL DEFAULT 1,
    [valid_from]      DATETIME      NOT NULL,
    [valid_until]     DATETIME      NOT NULL,
    [created_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT [CHK_Vouchers_Type] CHECK ([type] IN ('percentage', 'fixed', 'shipping')),
    CONSTRAINT [CHK_Vouchers_DiscountPercent] CHECK ([discount_percent] <= 100),
    INDEX [IX_Vouchers_Code] ([code]),
    INDEX [IX_Vouchers_IsActive] ([is_active]),
    INDEX [IX_Vouchers_ValidFrom] ([valid_from]),
    INDEX [IX_Vouchers_ValidUntil] ([valid_until])
);
```

### Table: UserVouchers

```sql
CREATE TABLE [dbo].[UserVouchers] (
    [id]          BIGINT   IDENTITY(1,1) PRIMARY KEY,
    [voucher_id]  BIGINT   NOT NULL,
    [user_id]     BIGINT   NOT NULL,
    [status]      NVARCHAR(20) NOT NULL DEFAULT 'available',
    [obtained_at] DATETIME NOT NULL DEFAULT GETDATE(),
    [used_at]     DATETIME,
    [order_id]    BIGINT,
    [expires_at]  DATETIME,
    
    FOREIGN KEY ([voucher_id]) REFERENCES [Vouchers]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([order_id]) REFERENCES [Orders]([id]),
    CONSTRAINT [CHK_UserVouchers_Status] CHECK ([status] IN ('available', 'used', 'expired', 'revoked')),
    INDEX [IX_UserVouchers_UserID] ([user_id]),
    INDEX [IX_UserVouchers_VoucherID] ([voucher_id]),
    INDEX [IX_UserVouchers_Status] ([status])
);
```

### Table: FlashSales

```sql
CREATE TABLE [dbo].[FlashSales] (
    [id]              BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [name]            NVARCHAR(255) NOT NULL,
    [description]     NVARCHAR(500),
    [banner_url]      NVARCHAR(500),
    [start_time]      DATETIME      NOT NULL,
    [end_time]        DATETIME      NOT NULL,
    [is_active]       BIT           NOT NULL DEFAULT 1,
    [created_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    
    CONSTRAINT [CHK_FlashSales_Time] CHECK ([end_time] > [start_time]),
    INDEX [IX_FlashSales_StartTime] ([start_time]),
    INDEX [IX_FlashSales_IsActive] ([is_active])
);
```

### Table: FlashSaleProducts

```sql
CREATE TABLE [dbo].[FlashSaleProducts] (
    [id]           BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [flash_sale_id] BIGINT       NOT NULL,
    [product_id]   BIGINT        NOT NULL,
    [flash_price]  DECIMAL(18,2) NOT NULL,
    [original_price] DECIMAL(18,2) NOT NULL,
    [quantity]     INT           NOT NULL,
    [sold_count]   INT           DEFAULT 0,
    [sort_order]   INT           DEFAULT 0,
    
    FOREIGN KEY ([flash_sale_id]) REFERENCES [FlashSales]([id]) ON DELETE CASCADE,
    FOREIGN KEY ([product_id]) REFERENCES [Products]([id]) ON DELETE CASCADE,
    CONSTRAINT [UQ_FlashSaleProducts] UNIQUE ([flash_sale_id], [product_id]),
    INDEX [IX_FlashSaleProducts_FlashSaleID] ([flash_sale_id])
);
```

---

## PART 11 — Notification Module

### Table: Notifications

```sql
CREATE TABLE [dbo].[Notifications] (
    [id]           BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]      BIGINT        NOT NULL,
    [type]         NVARCHAR(50)  NOT NULL,
    [title]        NVARCHAR(255) NOT NULL,
    [message]      NVARCHAR(1000) NOT NULL,
    [data]         NVARCHAR(MAX), -- JSON
    [image_url]    NVARCHAR(500),
    [action_url]   NVARCHAR(500),
    [is_read]      BIT           NOT NULL DEFAULT 0,
    [read_at]      DATETIME,
    [sent_at]      DATETIME      NOT NULL DEFAULT GETDATE(),
    [expires_at]   DATETIME,
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]) ON DELETE CASCADE,
    INDEX [IX_Notifications_UserID] ([user_id]),
    INDEX [IX_Notifications_IsRead] ([is_read]),
    INDEX [IX_Notifications_SentAt] ([sent_at])
);
```

### Table: NotificationTemplates

```sql
CREATE TABLE [dbo].[NotificationTemplates] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [name]        NVARCHAR(100) NOT NULL UNIQUE,
    [type]        NVARCHAR(50)  NOT NULL, -- email, sms, push, in_app
    [subject]     NVARCHAR(255),
    [body]        NVARCHAR(MAX) NOT NULL,
    [variables]   NVARCHAR(MAX), -- JSON array of variable names
    [is_active]   BIT           NOT NULL DEFAULT 1,
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    INDEX [IX_NotificationTemplates_Type] ([type])
);
```

---

## PART 12 — Admin Module

### Table: AdminUsers

```sql
CREATE TABLE [dbo].[AdminUsers] (
    [id]              BIGINT       IDENTITY(1,1) PRIMARY KEY,
    [user_id]         BIGINT       NOT NULL UNIQUE,
    [employee_id]     NVARCHAR(50) UNIQUE,
    [department]      NVARCHAR(100),
    [position]        NVARCHAR(100),
    [permissions]     NVARCHAR(MAX), -- JSON
    [last_password_change] DATETIME,
    [require_password_change] BIT NOT NULL DEFAULT 0,
    [created_at]      DATETIME     NOT NULL DEFAULT GETDATE(),
    [updated_at]      DATETIME     NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    INDEX [IX_AdminUsers_UserID] ([user_id])
);
```

### Table: AuditLogs

```sql
CREATE TABLE [dbo].[AuditLogs] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [user_id]     BIGINT,
    [action]      NVARCHAR(100) NOT NULL,
    [entity_type] NVARCHAR(50)  NOT NULL,
    [entity_id]   BIGINT,
    [old_values]  NVARCHAR(MAX), -- JSON
    [new_values]  NVARCHAR(MAX), -- JSON
    [ip_address]  NVARCHAR(45),
    [user_agent]  NVARCHAR(500),
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([user_id]) REFERENCES [Users]([id]),
    INDEX [IX_AuditLogs_UserID] ([user_id]),
    INDEX [IX_AuditLogs_Action] ([action]),
    INDEX [IX_AuditLogs_EntityType] ([entity_type]),
    INDEX [IX_AuditLogs_EntityID] ([entity_id]),
    INDEX [IX_AuditLogs_CreatedAt] ([created_at])
);
```

### Table: SystemSettings

```sql
CREATE TABLE [dbo].[SystemSettings] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [key]         NVARCHAR(100) NOT NULL UNIQUE,
    [value]       NVARCHAR(MAX) NOT NULL,
    [type]        NVARCHAR(50)  NOT NULL, -- string, number, boolean, json
    [description] NVARCHAR(500),
    [is_public]   BIT           NOT NULL DEFAULT 0,
    [updated_by]  BIGINT,
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([updated_by]) REFERENCES [AdminUsers]([id]),
    INDEX [IX_SystemSettings_Key] ([key])
);
```

---

## PART 13 — Database Relationship Summary

```
Users (1) ──────< (1) Shop
Users (1) ──────< (N) Orders
Users (1) ──────< (1) Cart
Users (1) ──────< (N) Reviews
Users (1) ──────< (N) UserAddresses
Users (1) ──────< (N) Payments

Shop (1) ──────< (N) Products
Shop (1) ──────< (N) OrderItems
Shop (1) ──────< (N) Reviews

Categories (1) ──────< (N) Products

Products (1) ──────< (N) ProductImages
Products (1) ──────< (N) ProductVariants
Products (1) ──────< (N) CartItems
Products (1) ──────< (N) OrderItems
Products (1) ──────< (N) Reviews

Cart (1) ──────< (N) CartItems

Orders (1) ──────< (N) OrderItems
Orders (1) ──────< (N) Payments
Orders (1) ──────< (N) OrderStatusHistory
Orders (1) ──────< (1) OrderShipping

Payments (1) ──────< (N) PaymentTransactions
Payments (1) ──────< (N) Refunds

Vouchers (1) ──────< (N) UserVouchers

FlashSales (1) ──────< (N) FlashSaleProducts
```

---

## PART 14 — Indexing Strategy

### High Priority Indexes

| Table | Columns | Type | Purpose |
|-------|---------|------|---------|
| Users | email | Unique | Login lookup |
| Users | phone | Unique | Phone lookup |
| Products | shop_id, status | Composite | Shop products |
| Products | category_id, status | Composite | Category browse |
| Products | status, created_at | Composite | New products |
| Products | status, sold_count | Composite | Best sellers |
| Orders | user_id, created_at | Composite | User orders |
| Orders | shop_id, created_at | Composite | Shop orders |
| Orders | status, created_at | Composite | Order management |
| CartItems | cart_id | Simple | Cart retrieval |
| Reviews | product_id, is_visible | Composite | Product reviews |

### Full-Text Search Indexes

```sql
-- Product search
CREATE FULLTEXT CATALOG FTCatalog AS DEFAULT;
CREATE FULLTEXT INDEX ON Products(name, description, short_description)
KEY INDEX PK_Products ON FTCatalog;

-- Shop search
CREATE FULLTEXT INDEX ON Shops(name, description)
KEY INDEX PK_Shops ON FTCatalog;
```

---

This database design provides a solid foundation for a scalable e-commerce platform capable of handling millions of users and products.
