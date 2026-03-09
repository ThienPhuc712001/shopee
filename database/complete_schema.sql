-- ============================================
-- E-Commerce Database - Complete Schema
-- SQL Server - Production Ready
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
        id INT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL,
        full_name NVARCHAR(200) NOT NULL,
        phone NVARCHAR(20) NOT NULL,
        address NVARCHAR(500) NOT NULL,
        ward NVARCHAR(200),
        district NVARCHAR(200),
        city NVARCHAR(200),
        state NVARCHAR(200),
        country NVARCHAR(100) DEFAULT 'Vietnam',
        postal_code NVARCHAR(20),
        is_default BIT DEFAULT 0,
        address_type NVARCHAR(20) DEFAULT 'home',
        latitude DECIMAL(10, 8),
        longitude DECIMAL(11, 8),
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
        flash_sale_start DATETIME,
        flash_sale_end DATETIME,
        flash_sale_price DECIMAL(18,2),
        meta_title NVARCHAR(200),
        meta_description NVARCHAR(500),
        meta_keywords NVARCHAR(500),
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
        id INT IDENTITY(1,1) PRIMARY KEY,
        product_id BIGINT NOT NULL,
        image_url NVARCHAR(500) NOT NULL,
        sort_order INT DEFAULT 0,
        is_primary BIT DEFAULT 0,
        alt_text NVARCHAR(255),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

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
        stock INT DEFAULT 0,
        reserved_stock INT DEFAULT 0,
        attributes NVARCHAR(MAX),
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
        refund_reason NVARCHAR(500),
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
        sku NVARCHAR(100),
        quantity INT NOT NULL,
        price DECIMAL(18,2) NOT NULL,
        original_price DECIMAL(18,2),
        discount DECIMAL(18,2) DEFAULT 0,
        subtotal DECIMAL(18,2) NOT NULL,
        shop_id BIGINT,
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
        transaction_id NVARCHAR(100) UNIQUE,
        payment_number NVARCHAR(50) UNIQUE,
        amount DECIMAL(18,2) NOT NULL,
        currency NVARCHAR(3) DEFAULT 'USD',
        payment_method NVARCHAR(50),
        payment_provider NVARCHAR(50),
        payment_gateway NVARCHAR(50),
        status NVARCHAR(20) DEFAULT 'pending',
        gateway_transaction_id NVARCHAR(100),
        gateway_response NVARCHAR(MAX),
        payment_url NVARCHAR(500),
        client_secret NVARCHAR(200),
        paid_at DATETIME,
        failed_at DATETIME,
        failure_reason NVARCHAR(500),
        refund_amount DECIMAL(18,2) DEFAULT 0,
        metadata NVARCHAR(MAX),
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
-- REFUNDS TABLE
-- ============================================
IF OBJECT_ID('dbo.refunds', 'U') IS NULL
BEGIN
    CREATE TABLE refunds (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        payment_id BIGINT NOT NULL,
        order_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        refund_number NVARCHAR(50) UNIQUE,
        amount DECIMAL(18,2) NOT NULL,
        reason NVARCHAR(500),
        type NVARCHAR(50),
        status NVARCHAR(20) DEFAULT 'pending',
        gateway_refund_id NVARCHAR(100),
        processed_at DATETIME,
        processed_by BIGINT,
        metadata NVARCHAR(MAX),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_refunds_payment FOREIGN KEY (payment_id) REFERENCES payments(id),
        CONSTRAINT FK_refunds_order FOREIGN KEY (order_id) REFERENCES orders(id),
        CONSTRAINT FK_refunds_user FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE INDEX IX_refunds_payment_id ON refunds(payment_id);
    CREATE INDEX IX_refunds_order_id ON refunds(order_id);
    CREATE INDEX IX_refunds_user_id ON refunds(user_id);
    CREATE INDEX IX_refunds_status ON refunds(status);
    CREATE INDEX IX_refunds_deleted_at ON refunds(deleted_at);
    
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
        order_item_id BIGINT,
        rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
        title NVARCHAR(200),
        content NVARCHAR(MAX),
        is_verified_purchase BIT DEFAULT 0,
        is_visible BIT DEFAULT 1,
        helpful_count INT DEFAULT 0,
        images NVARCHAR(MAX),
        seller_response NVARCHAR(MAX),
        seller_response_at DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_reviews_user FOREIGN KEY (user_id) REFERENCES users(id),
        CONSTRAINT FK_reviews_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT FK_reviews_shop FOREIGN KEY (shop_id) REFERENCES shops(id),
        CONSTRAINT FK_reviews_order FOREIGN KEY (order_id) REFERENCES orders(id),
        CONSTRAINT FK_reviews_order_item FOREIGN KEY (order_item_id) REFERENCES order_items(id)
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
-- REVIEW IMAGES TABLE
-- ============================================
IF OBJECT_ID('dbo.review_images', 'U') IS NULL
BEGIN
    CREATE TABLE review_images (
        id INT IDENTITY(1,1) PRIMARY KEY,
        review_id BIGINT NOT NULL,
        image_url NVARCHAR(500) NOT NULL,
        sort_order INT DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_review_images_review FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_review_images_review_id ON review_images(review_id);
    
    PRINT '✓ Table [review_images] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [review_images] already exists.';
END
GO

-- ============================================
-- ADMIN USERS TABLE
-- ============================================
IF OBJECT_ID('dbo.admin_users', 'U') IS NULL
BEGIN
    CREATE TABLE admin_users (
        id INT IDENTITY(1,1) PRIMARY KEY,
        username NVARCHAR(100) NOT NULL UNIQUE,
        email NVARCHAR(255) NOT NULL UNIQUE,
        password NVARCHAR(255) NOT NULL,
        full_name NVARCHAR(200),
        avatar_url NVARCHAR(500),
        role_id INT,
        status NVARCHAR(20) DEFAULT 'active',
        last_login DATETIME,
        last_login_ip NVARCHAR(50),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME
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
-- AUDIT LOGS TABLE
-- ============================================
IF OBJECT_ID('dbo.audit_logs', 'U') IS NULL
BEGIN
    CREATE TABLE audit_logs (
        id INT IDENTITY(1,1) PRIMARY KEY,
        admin_user_id INT,
        action NVARCHAR(100) NOT NULL,
        entity_type NVARCHAR(50),
        entity_id INT,
        old_value NVARCHAR(MAX),
        new_value NVARCHAR(MAX),
        ip_address NVARCHAR(50),
        user_agent NVARCHAR(500),
        created_at DATETIME NOT NULL DEFAULT GETDATE()
    );

    CREATE INDEX IX_audit_logs_admin_user_id ON audit_logs(admin_user_id);
    CREATE INDEX IX_audit_logs_action ON audit_logs(action);
    CREATE INDEX IX_audit_logs_entity ON audit_logs(entity_type, entity_id);
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
        id INT IDENTITY(1,1) PRIMARY KEY,
        key_name NVARCHAR(100) NOT NULL UNIQUE,
        value NVARCHAR(MAX),
        description NVARCHAR(500),
        updated_by INT,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE()
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
-- INSERT DEFAULT DATA
-- ============================================
PRINT '============================================';
PRINT 'Inserting default data...';
PRINT '============================================';

-- Insert default admin user (password: Admin@123)
IF NOT EXISTS (SELECT 1 FROM admin_users WHERE username = 'admin')
BEGIN
    INSERT INTO admin_users (username, email, password, full_name, role_id, status)
    VALUES ('admin', 'admin@ecommerce.com', '$2a$10$X.vXKZOF4nM1bq3jY8qN5.vLqKZ8xJ9yH5nF2mR8tP6wQ3sL4uV2G', 'System Administrator', 1, 'active');
    PRINT '✓ Default admin user created (email: admin@ecommerce.com, password: Admin@123)';
END
ELSE
BEGIN
    PRINT '○ Default admin user already exists';
END

-- Insert system admin (password: 712001lL)
IF NOT EXISTS (SELECT 1 FROM admin_users WHERE email = 'thienphuc71@gmail.com')
BEGIN
    INSERT INTO admin_users (username, email, password, full_name, role_id, status)
    VALUES ('thienphuc71', 'thienphuc71@gmail.com', '$2a$10$KZbFtpL4JuO5UzsogELZGOQJHBgEGkHbJN/DLADe40AkxBC9hEdO', 'Thien Phuc Admin', 1, 'active');
    PRINT '✓ Admin user thienphuc71@gmail.com created (password: 712001lL)';
END
ELSE
BEGIN
    PRINT '○ Admin user thienphuc71@gmail.com already exists';
END

-- Insert default system settings
IF NOT EXISTS (SELECT 1 FROM system_settings WHERE key_name = 'site_name')
BEGIN
    INSERT INTO system_settings (key_name, value, description) VALUES
    ('site_name', 'E-Commerce Store', 'Website name'),
    ('site_url', 'http://localhost:8080', 'Website URL'),
    ('currency', 'USD', 'Default currency'),
    ('tax_rate', '0', 'Default tax rate'),
    ('shipping_fee', '0', 'Default shipping fee');
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
PRINT '  - payments, refunds';
PRINT '  - reviews, review_images';
PRINT '  - admin_users, audit_logs, system_settings';
PRINT '============================================';
GO
