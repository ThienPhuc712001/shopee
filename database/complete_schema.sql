-- ============================================
-- E-Commerce Database - Complete Schema
-- SQL Server - Production Ready
-- Aligned with Go GORM Models
-- ============================================

-- Create database if not exists
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'ecommerce')
BEGIN
    CREATE DATABASE ecommerce;
    PRINT 'Database ecommerce created successfully.';
END
ELSE
BEGIN
    PRINT 'Database ecommerce already exists.';
END
GO

USE ecommerce;
GO

PRINT '============================================';
PRINT 'Starting schema creation...';
PRINT '============================================';
GO

-- ============================================
-- USERS TABLE
-- ============================================
IF OBJECT_ID('dbo.users', 'U') IS NULL
BEGIN
    CREATE TABLE users (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        email NVARCHAR(255) NOT NULL UNIQUE,
        password NVARCHAR(255) NOT NULL,
        first_name NVARCHAR(100),
        last_name NVARCHAR(100),
        phone NVARCHAR(20) UNIQUE,
        avatar NVARCHAR(500),
        role NVARCHAR(20) DEFAULT 'customer',
        status NVARCHAR(20) DEFAULT 'active',
        email_verified BIT DEFAULT 0,
        last_login DATETIME,
        failed_login_attempts INT DEFAULT 0,
        locked_until DATETIME,
        refresh_token NVARCHAR(500),
        refresh_token_expiry DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME
    );

    CREATE INDEX IX_users_email ON users(email);
    CREATE INDEX IX_users_phone ON users(phone);
    CREATE INDEX IX_users_role ON users(role);
    CREATE INDEX IX_users_status ON users(status);
    CREATE INDEX IX_users_deleted_at ON users(deleted_at);

    PRINT '✓ Table [users] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [users] already exists.';
END
GO

-- ============================================
-- ADDRESSES TABLE
-- ============================================
IF OBJECT_ID('dbo.addresses', 'U') IS NULL
BEGIN
    CREATE TABLE addresses (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL,
        name NVARCHAR(100) NOT NULL,
        phone NVARCHAR(20) NOT NULL,
        street NVARCHAR(500) NOT NULL,
        ward NVARCHAR(200),
        district NVARCHAR(200) NOT NULL,
        city NVARCHAR(200) NOT NULL,
        country NVARCHAR(100) DEFAULT 'Vietnam',
        is_default BIT DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_addresses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_addresses_user_id ON addresses(user_id);
    CREATE INDEX IX_addresses_deleted_at ON addresses(deleted_at);

    PRINT '✓ Table [addresses] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [addresses] already exists.';
END
GO

-- ============================================
-- SHOPS TABLE
-- ============================================
IF OBJECT_ID('dbo.shops', 'U') IS NULL
BEGIN
    CREATE TABLE shops (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL UNIQUE,
        name NVARCHAR(255) NOT NULL,
        slug NVARCHAR(255) UNIQUE,
        description NVARCHAR(MAX),
        logo NVARCHAR(500),
        cover_image NVARCHAR(500),
        phone NVARCHAR(20),
        email NVARCHAR(255),
        address NVARCHAR(500),
        status NVARCHAR(20) DEFAULT 'pending',
        verification_status NVARCHAR(20) DEFAULT 'unverified',
        rating DECIMAL(3,2) DEFAULT 0,
        rating_count INT DEFAULT 0,
        follower_count INT DEFAULT 0,
        product_count INT DEFAULT 0,
        total_sales BIGINT DEFAULT 0,
        total_revenue DECIMAL(18,2) DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_shops_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_shops_user_id ON shops(user_id);
    CREATE INDEX IX_shops_status ON shops(status);
    CREATE INDEX IX_shops_deleted_at ON shops(deleted_at);

    PRINT '✓ Table [shops] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [shops] already exists.';
END
GO

-- ============================================
-- CATEGORIES TABLE
-- ============================================
IF OBJECT_ID('dbo.categories', 'U') IS NULL
BEGIN
    CREATE TABLE categories (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        parent_id BIGINT,
        name NVARCHAR(255) NOT NULL,
        slug NVARCHAR(255) UNIQUE NOT NULL,
        description NVARCHAR(MAX),
        icon_url NVARCHAR(500),
        image_url NVARCHAR(500),
        level INT DEFAULT 0,
        sort_order INT DEFAULT 0,
        is_active BIT DEFAULT 1,
        attributes NVARCHAR(MAX),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_categories_parent FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE NO ACTION
    );

    CREATE INDEX IX_categories_parent_id ON categories(parent_id);
    CREATE INDEX IX_categories_slug ON categories(slug);
    CREATE INDEX IX_categories_deleted_at ON categories(deleted_at);

    PRINT '✓ Table [categories] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [categories] already exists.';
END
GO

-- ============================================
-- PRODUCTS TABLE
-- ============================================
IF OBJECT_ID('dbo.products', 'U') IS NULL
BEGIN
    CREATE TABLE products (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        shop_id BIGINT NOT NULL,
        category_id BIGINT NOT NULL,
        name NVARCHAR(500) NOT NULL,
        slug NVARCHAR(500) NOT NULL,
        description NVARCHAR(MAX),
        short_description NVARCHAR(1000),
        sku NVARCHAR(100),
        brand NVARCHAR(200),
        price DECIMAL(18,2) NOT NULL,
        original_price DECIMAL(18,2),
        discount_percent INT DEFAULT 0,
        cost DECIMAL(18,2),
        stock INT NOT NULL DEFAULT 0,
        reserved_stock INT DEFAULT 0,
        available_stock INT,
        sold_count BIGINT DEFAULT 0,
        view_count BIGINT DEFAULT 0,
        rating_avg DECIMAL(3,2) DEFAULT 0,
        rating_count INT DEFAULT 0,
        review_count INT DEFAULT 0,
        status NVARCHAR(20) DEFAULT 'draft',
        is_featured BIT DEFAULT 0,
        is_flash_sale BIT DEFAULT 0,
        flash_sale_price DECIMAL(18,2),
        flash_sale_start DATETIME,
        flash_sale_end DATETIME,
        weight DECIMAL(10,2),
        dimensions NVARCHAR(50),
        warranty_period NVARCHAR(50),
        return_days INT DEFAULT 7,
        tags NVARCHAR(MAX),
        meta_title NVARCHAR(255),
        meta_description NVARCHAR(500),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_products_shop FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
        CONSTRAINT FK_products_category FOREIGN KEY (category_id) REFERENCES categories(id)
    );

    CREATE INDEX IX_products_shop_id ON products(shop_id);
    CREATE INDEX IX_products_category_id ON products(category_id);
    CREATE INDEX IX_products_slug ON products(slug);
    CREATE INDEX IX_products_sku ON products(sku);
    CREATE INDEX IX_products_status ON products(status);
    CREATE INDEX IX_products_sold_count ON products(sold_count);
    CREATE INDEX IX_products_deleted_at ON products(deleted_at);

    PRINT '✓ Table [products] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [products] already exists.';
END
GO

-- ============================================
-- PRODUCT IMAGES TABLE
-- ============================================
IF OBJECT_ID('dbo.product_images', 'U') IS NULL
BEGIN
    CREATE TABLE product_images (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        product_id BIGINT NOT NULL,
        url NVARCHAR(500) NOT NULL,
        alt_text NVARCHAR(255),
        is_primary BIT DEFAULT 0,
        sort_order INT DEFAULT 0,
        width INT,
        height INT,
        size_bytes BIGINT,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_product_images_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_product_images_product_id ON product_images(product_id);
    CREATE INDEX IX_product_images_deleted_at ON product_images(deleted_at);

    PRINT '✓ Table [product_images] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [product_images] already exists.';
END
GO

-- ============================================
-- PRODUCT VARIANTS TABLE
-- ============================================
IF OBJECT_ID('dbo.product_variants', 'U') IS NULL
BEGIN
    CREATE TABLE product_variants (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        product_id BIGINT NOT NULL,
        sku NVARCHAR(100),
        name NVARCHAR(200),
        price DECIMAL(18,2),
        original_price DECIMAL(18,2),
        stock INT NOT NULL DEFAULT 0,
        reserved_stock INT DEFAULT 0,
        attributes NVARCHAR(MAX) NOT NULL,
        image_url NVARCHAR(500),
        sort_order INT DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_product_variants_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_product_variants_product_id ON product_variants(product_id);
    CREATE INDEX IX_product_variants_sku ON product_variants(sku);
    CREATE INDEX IX_product_variants_deleted_at ON product_variants(deleted_at);

    PRINT '✓ Table [product_variants] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [product_variants] already exists.';
END
GO

-- ============================================
-- CARTS TABLE
-- ============================================
IF OBJECT_ID('dbo.carts', 'U') IS NULL
BEGIN
    CREATE TABLE carts (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL UNIQUE,
        total_items INT NOT NULL DEFAULT 0,
        subtotal DECIMAL(18,2) NOT NULL DEFAULT 0,
        discount DECIMAL(18,2) DEFAULT 0,
        total DECIMAL(18,2) NOT NULL DEFAULT 0,
        currency NVARCHAR(3) DEFAULT 'USD',
        last_activity DATETIME NOT NULL DEFAULT GETDATE(),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_carts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_carts_user_id ON carts(user_id);
    CREATE INDEX IX_carts_deleted_at ON carts(deleted_at);

    PRINT '✓ Table [carts] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [carts] already exists.';
END
GO

-- ============================================
-- CART ITEMS TABLE
-- ============================================
IF OBJECT_ID('dbo.cart_items', 'U') IS NULL
BEGIN
    CREATE TABLE cart_items (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        cart_id BIGINT NOT NULL,
        product_id BIGINT NOT NULL,
        variant_id BIGINT,
        quantity INT NOT NULL DEFAULT 1,
        price DECIMAL(18,2) NOT NULL,
        original_price DECIMAL(18,2),
        discount DECIMAL(18,2) DEFAULT 0,
        subtotal DECIMAL(18,2) NOT NULL,
        product_name NVARCHAR(500),
        product_image NVARCHAR(500),
        shop_id BIGINT,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_cart_items_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
        CONSTRAINT FK_cart_items_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT FK_cart_items_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id),
        CONSTRAINT FK_cart_items_shop FOREIGN KEY (shop_id) REFERENCES shops(id)
    );

    CREATE INDEX IX_cart_items_cart_id ON cart_items(cart_id);
    CREATE INDEX IX_cart_items_product_id ON cart_items(product_id);
    CREATE INDEX IX_cart_items_variant_id ON cart_items(variant_id);
    CREATE INDEX IX_cart_items_shop_id ON cart_items(shop_id);
    CREATE INDEX IX_cart_items_deleted_at ON cart_items(deleted_at);

    PRINT '✓ Table [cart_items] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [cart_items] already exists.';
END
GO

-- ============================================
-- ORDERS TABLE
-- ============================================
IF OBJECT_ID('dbo.orders', 'U') IS NULL
BEGIN
    CREATE TABLE orders (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        order_number NVARCHAR(50) NOT NULL UNIQUE,
        user_id BIGINT NOT NULL,
        shop_id BIGINT NOT NULL,
        parent_order_id BIGINT,
        status NVARCHAR(20) DEFAULT 'pending',
        payment_status NVARCHAR(20) DEFAULT 'pending',
        fulfillment_status NVARCHAR(20) DEFAULT 'unfulfilled',
        subtotal DECIMAL(18,2) NOT NULL,
        shipping_fee DECIMAL(18,2) DEFAULT 0,
        shipping_discount DECIMAL(18,2) DEFAULT 0,
        product_discount DECIMAL(18,2) DEFAULT 0,
        voucher_discount DECIMAL(18,2) DEFAULT 0,
        tax_amount DECIMAL(18,2) DEFAULT 0,
        total_amount DECIMAL(18,2) NOT NULL,
        paid_amount DECIMAL(18,2) DEFAULT 0,
        shipping_name NVARCHAR(200) NOT NULL,
        shipping_phone NVARCHAR(20) NOT NULL,
        shipping_address NVARCHAR(500) NOT NULL,
        shipping_ward NVARCHAR(200),
        shipping_district NVARCHAR(200),
        shipping_city NVARCHAR(200),
        shipping_state NVARCHAR(200),
        shipping_country NVARCHAR(100) DEFAULT 'Vietnam',
        shipping_postal_code NVARCHAR(20),
        shipping_method NVARCHAR(100),
        shipping_carrier NVARCHAR(100),
        tracking_number NVARCHAR(100),
        estimated_delivery DATETIME,
        buyer_note NVARCHAR(500),
        seller_note NVARCHAR(500),
        cancel_reason NVARCHAR(500),
        paid_at DATETIME,
        confirmed_at DATETIME,
        shipped_at DATETIME,
        delivered_at DATETIME,
        cancelled_at DATETIME,
        completed_at DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_orders_user FOREIGN KEY (user_id) REFERENCES users(id),
        CONSTRAINT FK_orders_shop FOREIGN KEY (shop_id) REFERENCES shops(id),
        CONSTRAINT FK_orders_parent FOREIGN KEY (parent_order_id) REFERENCES orders(id)
    );

    CREATE INDEX IX_orders_user_id ON orders(user_id);
    CREATE INDEX IX_orders_shop_id ON orders(shop_id);
    CREATE INDEX IX_orders_order_number ON orders(order_number);
    CREATE INDEX IX_orders_status ON orders(status);
    CREATE INDEX IX_orders_payment_status ON orders(payment_status);
    CREATE INDEX IX_orders_parent_order_id ON orders(parent_order_id);
    CREATE INDEX IX_orders_deleted_at ON orders(deleted_at);

    PRINT '✓ Table [orders] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [orders] already exists.';
END
GO

-- ============================================
-- ORDER ITEMS TABLE
-- ============================================
IF OBJECT_ID('dbo.order_items', 'U') IS NULL
BEGIN
    CREATE TABLE order_items (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        order_id BIGINT NOT NULL,
        product_id BIGINT NOT NULL,
        variant_id BIGINT,
        product_name NVARCHAR(500) NOT NULL,
        product_image NVARCHAR(500),
        product_sku NVARCHAR(100),
        quantity INT NOT NULL,
        price DECIMAL(18,2) NOT NULL,
        original_price DECIMAL(18,2),
        discount DECIMAL(18,2) DEFAULT 0,
        subtotal DECIMAL(18,2) NOT NULL,
        tax_amount DECIMAL(18,2) DEFAULT 0,
        final_amount DECIMAL(18,2) NOT NULL,
        shop_id BIGINT NOT NULL,
        fulfillment_status NVARCHAR(20) DEFAULT 'pending',
        tracking_number NVARCHAR(100),
        shipped_at DATETIME,
        delivered_at DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
        CONSTRAINT FK_order_items_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT FK_order_items_variant FOREIGN KEY (variant_id) REFERENCES product_variants(id),
        CONSTRAINT FK_order_items_shop FOREIGN KEY (shop_id) REFERENCES shops(id)
    );

    CREATE INDEX IX_order_items_order_id ON order_items(order_id);
    CREATE INDEX IX_order_items_product_id ON order_items(product_id);
    CREATE INDEX IX_order_items_variant_id ON order_items(variant_id);
    CREATE INDEX IX_order_items_shop_id ON order_items(shop_id);
    CREATE INDEX IX_order_items_deleted_at ON order_items(deleted_at);

    PRINT '✓ Table [order_items] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [order_items] already exists.';
END
GO

-- ============================================
-- PAYMENTS TABLE
-- ============================================
IF OBJECT_ID('dbo.payments', 'U') IS NULL
BEGIN
    CREATE TABLE payments (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        order_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        transaction_id NVARCHAR(255) UNIQUE,
        payment_method NVARCHAR(50) NOT NULL,
        payment_provider NVARCHAR(50),
        amount DECIMAL(18,2) NOT NULL,
        currency NVARCHAR(3) DEFAULT 'USD',
        status NVARCHAR(20) DEFAULT 'pending',
        gateway_response NVARCHAR(MAX),
        metadata NVARCHAR(MAX),
        paid_at DATETIME,
        failed_at DATETIME,
        failure_reason NVARCHAR(500),
        refunded_amount DECIMAL(18,2) DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_payments_order FOREIGN KEY (order_id) REFERENCES orders(id),
        CONSTRAINT FK_payments_user FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE INDEX IX_payments_order_id ON payments(order_id);
    CREATE INDEX IX_payments_user_id ON payments(user_id);
    CREATE INDEX IX_payments_transaction_id ON payments(transaction_id);
    CREATE INDEX IX_payments_status ON payments(status);
    CREATE INDEX IX_payments_deleted_at ON payments(deleted_at);

    PRINT '✓ Table [payments] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [payments] already exists.';
END
GO

-- ============================================
-- PAYMENT TRANSACTIONS TABLE
-- ============================================
IF OBJECT_ID('dbo.payment_transactions', 'U') IS NULL
BEGIN
    CREATE TABLE payment_transactions (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        payment_id BIGINT NOT NULL,
        type NVARCHAR(50) NOT NULL,
        amount DECIMAL(18,2) NOT NULL,
        status NVARCHAR(20) NOT NULL,
        gateway_id NVARCHAR(255),
        gateway_response NVARCHAR(MAX),
        processed_at DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_payment_transactions_payment FOREIGN KEY (payment_id) REFERENCES payments(id)
    );

    CREATE INDEX IX_payment_transactions_payment_id ON payment_transactions(payment_id);
    CREATE INDEX IX_payment_transactions_created_at ON payment_transactions(created_at);

    PRINT '✓ Table [payment_transactions] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [payment_transactions] already exists.';
END
GO

-- ============================================
-- REFUNDS TABLE
-- ============================================
IF OBJECT_ID('dbo.refunds', 'U') IS NULL
BEGIN
    CREATE TABLE refunds (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        payment_id BIGINT NOT NULL,
        order_id BIGINT NOT NULL,
        refund_number NVARCHAR(50) UNIQUE,
        amount DECIMAL(18,2) NOT NULL,
        reason NVARCHAR(500) NOT NULL,
        status NVARCHAR(20) DEFAULT 'pending',
        type NVARCHAR(20) NOT NULL,
        requested_by BIGINT,
        approved_by BIGINT,
        approved_at DATETIME,
        processed_at DATETIME,
        gateway_refund_id NVARCHAR(255),
        notes NVARCHAR(500),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_refunds_payment FOREIGN KEY (payment_id) REFERENCES payments(id),
        CONSTRAINT FK_refunds_order FOREIGN KEY (order_id) REFERENCES orders(id)
    );

    CREATE INDEX IX_refunds_payment_id ON refunds(payment_id);
    CREATE INDEX IX_refunds_order_id ON refunds(order_id);
    CREATE INDEX IX_refunds_status ON refunds(status);

    PRINT '✓ Table [refunds] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [refunds] already exists.';
END
GO

-- ============================================
-- REVIEWS TABLE
-- ============================================
IF OBJECT_ID('dbo.reviews', 'U') IS NULL
BEGIN
    CREATE TABLE reviews (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL,
        product_id BIGINT,
        shop_id BIGINT,
        order_id BIGINT,
        rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
        comment NVARCHAR(MAX),
        images NVARCHAR(MAX),
        is_approved BIT DEFAULT 0,
        helpful_count INT DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_reviews_user FOREIGN KEY (user_id) REFERENCES users(id),
        CONSTRAINT FK_reviews_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT FK_reviews_shop FOREIGN KEY (shop_id) REFERENCES shops(id),
        CONSTRAINT FK_reviews_order FOREIGN KEY (order_id) REFERENCES orders(id)
    );

    CREATE INDEX IX_reviews_user_id ON reviews(user_id);
    CREATE INDEX IX_reviews_product_id ON reviews(product_id);
    CREATE INDEX IX_reviews_shop_id ON reviews(shop_id);
    CREATE INDEX IX_reviews_order_id ON reviews(order_id);
    CREATE INDEX IX_reviews_rating ON reviews(rating);
    CREATE INDEX IX_reviews_deleted_at ON reviews(deleted_at);

    PRINT '✓ Table [reviews] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [reviews] already exists.';
END
GO

-- ============================================
-- ADMIN ROLES TABLE
-- ============================================
IF OBJECT_ID('dbo.admin_roles', 'U') IS NULL
BEGIN
    CREATE TABLE admin_roles (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(50) NOT NULL UNIQUE,
        description NVARCHAR(255),
        permissions NVARCHAR(MAX),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE()
    );

    CREATE INDEX IX_admin_roles_name ON admin_roles(name);

    PRINT '✓ Table [admin_roles] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [admin_roles] already exists.';
END
GO

-- ============================================
-- ADMIN USERS TABLE
-- ============================================
IF OBJECT_ID('dbo.admin_users', 'U') IS NULL
BEGIN
    CREATE TABLE admin_users (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        email NVARCHAR(255) NOT NULL UNIQUE,
        password NVARCHAR(255) NOT NULL,
        role_id BIGINT NOT NULL,
        first_name NVARCHAR(100),
        last_name NVARCHAR(100),
        phone NVARCHAR(20),
        avatar_url NVARCHAR(500),
        status NVARCHAR(20) DEFAULT 'active',
        last_login_at DATETIME,
        last_login_ip NVARCHAR(45),
        failed_login_attempts INT DEFAULT 0,
        locked_until DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_admin_users_role FOREIGN KEY (role_id) REFERENCES admin_roles(id)
    );

    CREATE INDEX IX_admin_users_email ON admin_users(email);
    CREATE INDEX IX_admin_users_status ON admin_users(status);
    CREATE INDEX IX_admin_users_deleted_at ON admin_users(deleted_at);

    PRINT '✓ Table [admin_users] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [admin_users] already exists.';
END
GO

-- ============================================
-- ADMIN PERMISSIONS TABLE
-- ============================================
IF OBJECT_ID('dbo.admin_permissions', 'U') IS NULL
BEGIN
    CREATE TABLE admin_permissions (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(100) NOT NULL UNIQUE,
        description NVARCHAR(255),
        module NVARCHAR(50),
        created_at DATETIME NOT NULL DEFAULT GETDATE()
    );

    CREATE INDEX IX_admin_permissions_name ON admin_permissions(name);
    CREATE INDEX IX_admin_permissions_module ON admin_permissions(module);

    PRINT '✓ Table [admin_permissions] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [admin_permissions] already exists.';
END
GO

-- ============================================
-- AUDIT LOGS TABLE
-- ============================================
IF OBJECT_ID('dbo.audit_logs', 'U') IS NULL
BEGIN
    CREATE TABLE audit_logs (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        admin_id BIGINT NOT NULL,
        action NVARCHAR(100) NOT NULL,
        entity_type NVARCHAR(50) NOT NULL,
        entity_id BIGINT,
        old_values NVARCHAR(MAX),
        new_values NVARCHAR(MAX),
        ip_address NVARCHAR(45),
        user_agent NVARCHAR(500),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_audit_logs_admin FOREIGN KEY (admin_id) REFERENCES admin_users(id)
    );

    CREATE INDEX IX_audit_logs_admin_id ON audit_logs(admin_id);
    CREATE INDEX IX_audit_logs_action ON audit_logs(action);
    CREATE INDEX IX_audit_logs_entity_type ON audit_logs(entity_type);
    CREATE INDEX IX_audit_logs_entity_id ON audit_logs(entity_id);
    CREATE INDEX IX_audit_logs_created_at ON audit_logs(created_at);

    PRINT '✓ Table [audit_logs] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [audit_logs] already exists.';
END
GO

-- ============================================
-- SYSTEM SETTINGS TABLE
-- ============================================
IF OBJECT_ID('dbo.system_settings', 'U') IS NULL
BEGIN
    CREATE TABLE system_settings (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        key_name NVARCHAR(100) NOT NULL UNIQUE,
        value NVARCHAR(MAX) NOT NULL,
        type NVARCHAR(50) NOT NULL,
        description NVARCHAR(500),
        is_public BIT DEFAULT 0,
        updated_by BIGINT,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_system_settings_admin FOREIGN KEY (updated_by) REFERENCES admin_users(id)
    );

    CREATE INDEX IX_system_settings_key ON system_settings(key_name);

    PRINT '✓ Table [system_settings] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [system_settings] already exists.';
END
GO

-- ============================================
-- MIGRATION: Update existing tables to match new schema
-- ============================================
PRINT '============================================';
PRINT 'Running migrations for existing tables...';
PRINT '============================================';

-- Update admin_users table
IF COL_LENGTH('dbo.admin_users', 'username') IS NOT NULL
BEGIN
    ALTER TABLE admin_users DROP COLUMN username;
    PRINT '✓ Removed [username] column from [admin_users]';
END

IF COL_LENGTH('dbo.admin_users', 'first_name') IS NULL
BEGIN
    ALTER TABLE admin_users ADD first_name NVARCHAR(100);
    PRINT '✓ Added [first_name] column to [admin_users]';
END

IF COL_LENGTH('dbo.admin_users', 'last_name') IS NULL
BEGIN
    ALTER TABLE admin_users ADD last_name NVARCHAR(100);
    PRINT '✓ Added [last_name] column to [admin_users]';
END

IF COL_LENGTH('dbo.admin_users', 'failed_login_attempts') IS NULL
BEGIN
    ALTER TABLE admin_users ADD failed_login_attempts INT DEFAULT 0;
    PRINT '✓ Added [failed_login_attempts] column to [admin_users]';
END

IF COL_LENGTH('dbo.admin_users', 'locked_until') IS NULL
BEGIN
    ALTER TABLE admin_users ADD locked_until DATETIME;
    PRINT '✓ Added [locked_until] column to [admin_users]';
END

-- Fix admin_users.role_id data type to match admin_roles.id
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('dbo.admin_users') AND name = 'role_id' AND user_type_id != TYPE_ID('BIGINT'))
BEGIN
    -- First drop the foreign key constraint if it exists
    IF EXISTS (SELECT 1 FROM sys.foreign_keys WHERE name = 'FK_admin_users_role')
    BEGIN
        ALTER TABLE admin_users DROP CONSTRAINT FK_admin_users_role;
        PRINT '✓ Dropped foreign key [FK_admin_users_role]';
    END
    
    -- Alter column to BIGINT
    ALTER TABLE admin_users ALTER COLUMN role_id BIGINT NOT NULL;
    PRINT '✓ Changed [role_id] to BIGINT in [admin_users]';
    
    -- Recreate foreign key
    ALTER TABLE admin_users ADD CONSTRAINT FK_admin_users_role FOREIGN KEY (role_id) REFERENCES admin_roles(id);
    PRINT '✓ Recreated foreign key [FK_admin_users_role]';
END

-- Update audit_logs table
IF COL_LENGTH('dbo.audit_logs', 'admin_user_id') IS NOT NULL
BEGIN
    IF COL_LENGTH('dbo.audit_logs', 'admin_id') IS NULL
    BEGIN
        EXEC sp_rename 'dbo.audit_logs.admin_user_id', 'admin_id', 'COLUMN';
        PRINT '✓ Renamed [admin_user_id] to [admin_id] in [audit_logs]';
    END
    ELSE
    BEGIN
        ALTER TABLE audit_logs DROP COLUMN admin_user_id;
        PRINT '✓ Removed [admin_user_id] column from [audit_logs]';
    END
END

-- Fix audit_logs.admin_id data type
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('dbo.audit_logs') AND name = 'admin_id' AND user_type_id != TYPE_ID('BIGINT'))
BEGIN
    ALTER TABLE audit_logs ALTER COLUMN admin_id BIGINT NOT NULL;
    PRINT '✓ Changed [admin_id] to BIGINT in [audit_logs]';
END

-- Update system_settings table
IF COL_LENGTH('dbo.system_settings', 'type') IS NULL
BEGIN
    ALTER TABLE system_settings ADD type NVARCHAR(50) NOT NULL DEFAULT 'string';
    PRINT '✓ Added [type] column to [system_settings]';
END

IF COL_LENGTH('dbo.system_settings', 'is_public') IS NULL
BEGIN
    ALTER TABLE system_settings ADD is_public BIT DEFAULT 0;
    PRINT '✓ Added [is_public] column to [system_settings]';
END

-- Fix system_settings.updated_by data type
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('dbo.system_settings') AND name = 'updated_by' AND user_type_id != TYPE_ID('BIGINT'))
BEGIN
    ALTER TABLE system_settings ALTER COLUMN updated_by BIGINT;
    PRINT '✓ Changed [updated_by] to BIGINT in [system_settings]';
END

-- Update reviews table
IF COL_LENGTH('dbo.reviews', 'title') IS NOT NULL
BEGIN
    ALTER TABLE reviews DROP COLUMN title;
    PRINT '✓ Removed [title] column from [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'content') IS NOT NULL
BEGIN
    ALTER TABLE reviews DROP COLUMN content;
    PRINT '✓ Removed [content] column from [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'comment') IS NULL
BEGIN
    ALTER TABLE reviews ADD comment NVARCHAR(MAX);
    PRINT '✓ Added [comment] column to [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'images') IS NULL
BEGIN
    ALTER TABLE reviews ADD images NVARCHAR(MAX);
    PRINT '✓ Added [images] column to [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'is_approved') IS NULL
BEGIN
    ALTER TABLE reviews ADD is_approved BIT DEFAULT 0;
    PRINT '✓ Added [is_approved] column to [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'is_verified_purchase') IS NOT NULL
BEGIN
    ALTER TABLE reviews DROP COLUMN is_verified_purchase;
    PRINT '✓ Removed [is_verified_purchase] column from [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'seller_response') IS NOT NULL
BEGIN
    ALTER TABLE reviews DROP COLUMN seller_response;
    PRINT '✓ Removed [seller_response] column from [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'seller_response_at') IS NOT NULL
BEGIN
    ALTER TABLE reviews DROP COLUMN seller_response_at;
    PRINT '✓ Removed [seller_response_at] column from [reviews]';
END

IF COL_LENGTH('dbo.reviews', 'order_item_id') IS NOT NULL
BEGIN
    ALTER TABLE reviews DROP COLUMN order_item_id;
    PRINT '✓ Removed [order_item_id] column from [reviews]';
END

-- Update products table
IF COL_LENGTH('dbo.products', 'weight') IS NULL
BEGIN
    ALTER TABLE products ADD weight DECIMAL(10,2);
    PRINT '✓ Added [weight] column to [products]';
END

IF COL_LENGTH('dbo.products', 'dimensions') IS NULL
BEGIN
    ALTER TABLE products ADD dimensions NVARCHAR(50);
    PRINT '✓ Added [dimensions] column to [products]';
END

IF COL_LENGTH('dbo.products', 'warranty_period') IS NULL
BEGIN
    ALTER TABLE products ADD warranty_period NVARCHAR(50);
    PRINT '✓ Added [warranty_period] column to [products]';
END

IF COL_LENGTH('dbo.products', 'return_days') IS NULL
BEGIN
    ALTER TABLE products ADD return_days INT DEFAULT 7;
    PRINT '✓ Added [return_days] column to [products]';
END

IF COL_LENGTH('dbo.products', 'tags') IS NULL
BEGIN
    ALTER TABLE products ADD tags NVARCHAR(MAX);
    PRINT '✓ Added [tags] column to [products]';
END

-- Update product_images table
IF COL_LENGTH('dbo.product_images', 'width') IS NULL
BEGIN
    ALTER TABLE product_images ADD width INT;
    PRINT '✓ Added [width] column to [product_images]';
END

IF COL_LENGTH('dbo.product_images', 'height') IS NULL
BEGIN
    ALTER TABLE product_images ADD height INT;
    PRINT '✓ Added [height] column to [product_images]';
END

IF COL_LENGTH('dbo.product_images', 'size_bytes') IS NULL
BEGIN
    ALTER TABLE product_images ADD size_bytes BIGINT;
    PRINT '✓ Added [size_bytes] column to [product_images]';
END

-- Update product_variants table
IF COL_LENGTH('dbo.product_variants', 'original_price') IS NULL
BEGIN
    ALTER TABLE product_variants ADD original_price DECIMAL(18,2);
    PRINT '✓ Added [original_price] column to [product_variants]';
END

IF COL_LENGTH('dbo.product_variants', 'image_url') IS NULL
BEGIN
    ALTER TABLE product_variants ADD image_url NVARCHAR(500);
    PRINT '✓ Added [image_url] column to [product_variants]';
END

IF COL_LENGTH('dbo.product_variants', 'sort_order') IS NULL
BEGIN
    ALTER TABLE product_variants ADD sort_order INT DEFAULT 0;
    PRINT '✓ Added [sort_order] column to [product_variants]';
END

-- Update order_items table
IF COL_LENGTH('dbo.order_items', 'product_sku') IS NULL
BEGIN
    ALTER TABLE order_items ADD product_sku NVARCHAR(100);
    PRINT '✓ Added [product_sku] column to [order_items]';
END

IF COL_LENGTH('dbo.order_items', 'tax_amount') IS NULL
BEGIN
    ALTER TABLE order_items ADD tax_amount DECIMAL(18,2) DEFAULT 0;
    PRINT '✓ Added [tax_amount] column to [order_items]';
END

IF COL_LENGTH('dbo.order_items', 'final_amount') IS NULL
BEGIN
    ALTER TABLE order_items ADD final_amount DECIMAL(18,2) NOT NULL DEFAULT 0;
    PRINT '✓ Added [final_amount] column to [order_items]';
END

IF COL_LENGTH('dbo.order_items', 'fulfillment_status') IS NULL
BEGIN
    ALTER TABLE order_items ADD fulfillment_status NVARCHAR(20) DEFAULT 'pending';
    PRINT '✓ Added [fulfillment_status] column to [order_items]';
END

IF COL_LENGTH('dbo.order_items', 'tracking_number') IS NULL
BEGIN
    ALTER TABLE order_items ADD tracking_number NVARCHAR(100);
    PRINT '✓ Added [tracking_number] column to [order_items]';
END

IF COL_LENGTH('dbo.order_items', 'shipped_at') IS NULL
BEGIN
    ALTER TABLE order_items ADD shipped_at DATETIME;
    PRINT '✓ Added [shipped_at] column to [order_items]';
END

IF COL_LENGTH('dbo.order_items', 'delivered_at') IS NULL
BEGIN
    ALTER TABLE order_items ADD delivered_at DATETIME;
    PRINT '✓ Added [delivered_at] column to [order_items]';
END

-- Update orders table
IF COL_LENGTH('dbo.orders', 'confirmed_at') IS NULL
BEGIN
    ALTER TABLE orders ADD confirmed_at DATETIME;
    PRINT '✓ Added [confirmed_at] column to [orders]';
END

-- Update payments table
IF COL_LENGTH('dbo.payments', 'refund_amount') IS NOT NULL
BEGIN
    IF COL_LENGTH('dbo.payments', 'refunded_amount') IS NULL
    BEGIN
        EXEC sp_rename 'dbo.payments.refund_amount', 'refunded_amount', 'COLUMN';
        PRINT '✓ Renamed [refund_amount] to [refunded_amount] in [payments]';
    END
    ELSE
    BEGIN
        ALTER TABLE payments DROP COLUMN refund_amount;
        PRINT '✓ Removed [refund_amount] column from [payments]';
    END
END

-- Update refunds table
IF COL_LENGTH('dbo.refunds', 'requested_by') IS NULL
BEGIN
    ALTER TABLE refunds ADD requested_by BIGINT;
    PRINT '✓ Added [requested_by] column to [refunds]';
END

IF COL_LENGTH('dbo.refunds', 'approved_by') IS NULL
BEGIN
    ALTER TABLE refunds ADD approved_by BIGINT;
    PRINT '✓ Added [approved_by] column to [refunds]';
END

IF COL_LENGTH('dbo.refunds', 'approved_at') IS NULL
BEGIN
    ALTER TABLE refunds ADD approved_at DATETIME;
    PRINT '✓ Added [approved_at] column to [refunds]';
END

IF COL_LENGTH('dbo.refunds', 'notes') IS NULL
BEGIN
    ALTER TABLE refunds ADD notes NVARCHAR(500);
    PRINT '✓ Added [notes] column to [refunds]';
END

-- Update addresses table
IF COL_LENGTH('dbo.addresses', 'full_name') IS NOT NULL
BEGIN
    IF COL_LENGTH('dbo.addresses', 'name') IS NULL
    BEGIN
        EXEC sp_rename 'dbo.addresses.full_name', 'name', 'COLUMN';
        PRINT '✓ Renamed [full_name] to [name] in [addresses]';
    END
    ELSE
    BEGIN
        ALTER TABLE addresses DROP COLUMN full_name;
        PRINT '✓ Removed [full_name] column from [addresses]';
    END
END

IF COL_LENGTH('dbo.addresses', 'address') IS NOT NULL
BEGIN
    IF COL_LENGTH('dbo.addresses', 'street') IS NULL
    BEGIN
        EXEC sp_rename 'dbo.addresses.address', 'street', 'COLUMN';
        PRINT '✓ Renamed [address] to [street] in [addresses]';
    END
    ELSE
    BEGIN
        ALTER TABLE addresses DROP COLUMN address;
        PRINT '✓ Removed [address] column from [addresses]';
    END
END

IF COL_LENGTH('dbo.addresses', 'address_type') IS NOT NULL
BEGIN
    ALTER TABLE addresses DROP COLUMN address_type;
    PRINT '✓ Removed [address_type] column from [addresses]';
END

IF COL_LENGTH('dbo.addresses', 'latitude') IS NOT NULL
BEGIN
    ALTER TABLE addresses DROP COLUMN latitude;
    PRINT '✓ Removed [latitude] column from [addresses]';
END

IF COL_LENGTH('dbo.addresses', 'longitude') IS NOT NULL
BEGIN
    ALTER TABLE addresses DROP COLUMN longitude;
    PRINT '✓ Removed [longitude] column from [addresses]';
END

IF COL_LENGTH('dbo.addresses', 'state') IS NOT NULL
BEGIN
    ALTER TABLE addresses DROP COLUMN state;
    PRINT '✓ Removed [state] column from [addresses]';
END

IF COL_LENGTH('dbo.addresses', 'postal_code') IS NOT NULL
BEGIN
    ALTER TABLE addresses DROP COLUMN postal_code;
    PRINT '✓ Removed [postal_code] column from [addresses]';
END

PRINT '✓ Migrations completed';
PRINT '';

-- ============================================
-- INSERT DEFAULT DATA
-- ============================================
PRINT '============================================';
PRINT 'Inserting default data...';
PRINT '============================================';

-- Insert default admin roles
IF NOT EXISTS (SELECT 1 FROM admin_roles WHERE name = 'super_admin')
BEGIN
    INSERT INTO admin_roles (name, description, permissions) VALUES
    ('super_admin', 'Super Administrator with full access', '["*"]'),
    ('admin', 'Administrator with limited access', '["users:read", "products:*", "orders:*"]'),
    ('support_agent', 'Support agent with basic access', '["orders:read", "orders:update", "refunds:read"]');
    PRINT '✓ Default admin roles created';
END
ELSE
BEGIN
    PRINT '○ Default admin roles already exist';
END

-- Insert default admin user (password: Admin@123)
IF NOT EXISTS (SELECT 1 FROM admin_users WHERE email = 'admin@ecommerce.com')
BEGIN
    INSERT INTO admin_users (email, password, role_id, first_name, last_name, status)
    VALUES ('admin@ecommerce.com', '$2a$10$X.vXKZOF4nM1bq3jY8qN5.vLqKZ8xJ9yH5nF2mR8tP6wQ3sL4uV2G', 1, 'System', 'Administrator', 'active');
    PRINT '✓ Default admin user created (email: admin@ecommerce.com, password: Admin@123)';
END
ELSE
BEGIN
    PRINT '○ Default admin user already exists';
END

-- Insert system admin (password: 712001lL)
IF NOT EXISTS (SELECT 1 FROM admin_users WHERE email = 'thienphuc71@gmail.com')
BEGIN
    INSERT INTO admin_users (email, password, role_id, first_name, last_name, status)
    VALUES ('thienphuc71@gmail.com', '$2a$10$KZbFtpL4JuO5UzsogELZGOQJHBgEGkHbJN/DLADe40AkxBC9hEdO', 1, 'Thien Phuc', 'Admin', 'active');
    PRINT '✓ Admin user thienphuc71@gmail.com created (password: 712001lL)';
END
ELSE
BEGIN
    PRINT '○ Admin user thienphuc71@gmail.com already exists';
END

-- Insert default system settings
IF NOT EXISTS (SELECT 1 FROM system_settings WHERE key_name = 'site_name')
BEGIN
    INSERT INTO system_settings (key_name, value, type, description) VALUES
    ('site_name', 'E-Commerce Store', 'string', 'Website name'),
    ('site_url', 'http://localhost:8080', 'string', 'Website URL'),
    ('currency', 'USD', 'string', 'Default currency'),
    ('tax_rate', '0', 'number', 'Default tax rate'),
    ('shipping_fee', '0', 'number', 'Default shipping fee');
    PRINT '✓ Default system settings created';
END
ELSE
BEGIN
    PRINT '○ Default system settings already exist';
END

-- Insert default categories
IF NOT EXISTS (SELECT 1 FROM categories WHERE slug = 'electronics')
BEGIN
    INSERT INTO categories (name, slug, description, level, sort_order, is_active) VALUES
    ('Electronics', 'electronics', 'Electronic devices and accessories', 0, 1, 1),
    ('Fashion', 'fashion', 'Clothing and fashion accessories', 0, 2, 1),
    ('Home & Living', 'home-living', 'Home decor and living essentials', 0, 3, 1),
    ('Books', 'books', 'Books and publications', 0, 4, 1),
    ('Sports', 'sports', 'Sports equipment and accessories', 0, 5, 1);
    PRINT '✓ Default categories created';
END
ELSE
BEGIN
    PRINT '○ Default categories already exist';
END

PRINT '============================================';
PRINT 'Schema creation completed successfully!';
PRINT '============================================';
PRINT '';
PRINT 'Admin Accounts:';
PRINT '  1. admin@ecommerce.com / Admin@123';
PRINT '  2. thienphuc71@gmail.com / 712001lL';
PRINT '';
PRINT 'Tables created:';
PRINT '  - users, addresses, shops, categories';
PRINT '  - products, product_images, product_variants';
PRINT '  - carts, cart_items';
PRINT '  - orders, order_items';
PRINT '  - payments, payment_transactions, refunds';
PRINT '  - reviews';
PRINT '  - admin_roles, admin_users, admin_permissions, audit_logs, system_settings';
PRINT '============================================';
GO
