-- ============================================
-- Migration: Create Coupon Tables
-- E-Commerce Database - Coupon System
-- SQL Server
-- ============================================

USE ecommerce;
GO

PRINT '============================================';
PRINT 'Creating Coupon Tables...';
PRINT '============================================';

-- ============================================
-- COUPONS TABLE
-- ============================================
IF OBJECT_ID('dbo.coupons', 'U') IS NULL
BEGIN
    CREATE TABLE coupons (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        code NVARCHAR(50) NOT NULL UNIQUE,
        name NVARCHAR(255) NOT NULL,
        description NVARCHAR(MAX),
        discount_type NVARCHAR(50) NOT NULL, -- 'percentage', 'fixed', 'free_shipping'
        discount_value DECIMAL(18,2) NOT NULL,
        max_discount DECIMAL(18,2) DEFAULT 0, -- For percentage discounts
        min_order_value DECIMAL(18,2) DEFAULT 0,
        max_order_value DECIMAL(18,2) DEFAULT 0,
        usage_limit INT DEFAULT 0, -- 0 = unlimited
        used_count INT DEFAULT 0,
        usage_limit_per_user INT DEFAULT 1,
        start_date DATETIME,
        end_date DATETIME NOT NULL,
        is_active BIT DEFAULT 1,
        status NVARCHAR(50) DEFAULT 'active', -- 'active', 'inactive', 'expired', 'used_up'
        applicable_categories NVARCHAR(MAX), -- JSON array of category IDs
        applicable_products NVARCHAR(MAX),   -- JSON array of product IDs
        excluded_categories NVARCHAR(MAX),   -- JSON array of category IDs
        excluded_products NVARCHAR(MAX),     -- JSON array of product IDs
        user_restricted BIT DEFAULT 0,
        restricted_users NVARCHAR(MAX),      -- JSON array of user IDs
        created_by BIGINT,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME
    );

    CREATE INDEX IX_coupons_code ON coupons(code);
    CREATE INDEX IX_coupons_end_date ON coupons(end_date);
    CREATE INDEX IX_coupons_status ON coupons(status);
    CREATE INDEX IX_coupons_is_active ON coupons(is_active);
    CREATE INDEX IX_coupons_deleted_at ON coupons(deleted_at);

    PRINT '✓ Table [coupons] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [coupons] already exists.';
END
GO

-- ============================================
-- COUPON USAGES TABLE
-- ============================================
IF OBJECT_ID('dbo.coupon_usages', 'U') IS NULL
BEGIN
    CREATE TABLE coupon_usages (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        coupon_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        order_id BIGINT NOT NULL,
        discount_amount DECIMAL(18,2) NOT NULL,
        used_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_coupon_usages_coupon FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE,
        CONSTRAINT FK_coupon_usages_user FOREIGN KEY (user_id) REFERENCES users(id),
        CONSTRAINT FK_coupon_usages_order FOREIGN KEY (order_id) REFERENCES orders(id)
    );

    CREATE INDEX IX_coupon_usages_coupon_id ON coupon_usages(coupon_id);
    CREATE INDEX IX_coupon_usages_user_id ON coupon_usages(user_id);
    CREATE INDEX IX_coupon_usages_order_id ON coupon_usages(order_id);
    CREATE INDEX IX_coupon_usages_used_at ON coupon_usages(used_at);

    PRINT '✓ Table [coupon_usages] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [coupon_usages] already exists.';
END
GO

-- ============================================
-- MIGRATION: Add coupon_id to orders table
-- ============================================
IF COL_LENGTH('dbo.orders', 'coupon_id') IS NULL
BEGIN
    ALTER TABLE orders ADD coupon_id BIGINT NULL;
    ALTER TABLE orders ADD coupon_discount DECIMAL(18,2) DEFAULT 0;
    ALTER TABLE orders ADD coupon_code NVARCHAR(50);
    
    CREATE INDEX IX_orders_coupon_id ON orders(coupon_id);
    
    PRINT '✓ Added [coupon_id], [coupon_discount], and [coupon_code] columns to [orders]';
END
ELSE
BEGIN
    PRINT '○ Column [coupon_id] already exists in [orders]';
END
GO

-- Add foreign key constraint for coupon_id if not exists
IF NOT EXISTS (SELECT 1 FROM sys.foreign_keys WHERE name = 'FK_orders_coupon')
BEGIN
    ALTER TABLE orders ADD CONSTRAINT FK_orders_coupon 
    FOREIGN KEY (coupon_id) REFERENCES coupons(id);
    
    PRINT '✓ Added foreign key constraint [FK_orders_coupon]';
END
GO

PRINT '============================================';
PRINT 'Coupon Tables Migration Complete!';
PRINT '============================================';
