-- ============================================
-- Migration: Create Shipping Tables
-- E-Commerce Database - Shipping System
-- SQL Server
-- ============================================

USE ecommerce;
GO

PRINT '============================================';
PRINT 'Creating Shipping Tables...';
PRINT '============================================';

-- ============================================
-- SHIPPING ADDRESSES TABLE
-- ============================================
IF OBJECT_ID('dbo.shipping_addresses', 'U') IS NULL
BEGIN
    CREATE TABLE shipping_addresses (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL,
        recipient_name NVARCHAR(200) NOT NULL,
        phone NVARCHAR(20) NOT NULL,
        address_line NVARCHAR(500) NOT NULL,
        ward NVARCHAR(200),
        district NVARCHAR(200) NOT NULL,
        city NVARCHAR(200) NOT NULL,
        postal_code NVARCHAR(20),
        country NVARCHAR(100) DEFAULT 'Vietnam',
        is_default BIT DEFAULT 0,
        latitude DECIMAL(10,8),
        longitude DECIMAL(11,8),
        notes NVARCHAR(500),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_shipping_addresses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_shipping_addresses_user_id ON shipping_addresses(user_id);
    CREATE INDEX IX_shipping_addresses_is_default ON shipping_addresses(is_default);
    CREATE INDEX IX_shipping_addresses_deleted_at ON shipping_addresses(deleted_at);

    PRINT '✓ Table [shipping_addresses] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [shipping_addresses] already exists.';
END
GO

-- ============================================
-- SHIPPING CARRIERS TABLE
-- ============================================
IF OBJECT_ID('dbo.shipping_carriers', 'U') IS NULL
BEGIN
    CREATE TABLE shipping_carriers (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(200) NOT NULL UNIQUE,
        code NVARCHAR(50) NOT NULL UNIQUE,
        type NVARCHAR(50) NOT NULL, -- 'internal', 'third_party', 'local'
        contact_name NVARCHAR(200),
        phone NVARCHAR(20),
        email NVARCHAR(255),
        website NVARCHAR(500),
        api_endpoint NVARCHAR(500),
        api_key NVARCHAR(500),
        is_active BIT DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE()
    );

    CREATE INDEX IX_shipping_carriers_code ON shipping_carriers(code);
    CREATE INDEX IX_shipping_carriers_is_active ON shipping_carriers(is_active);

    PRINT '✓ Table [shipping_carriers] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [shipping_carriers] already exists.';
END
GO

-- ============================================
-- SHIPMENTS TABLE
-- ============================================
IF OBJECT_ID('dbo.shipments', 'U') IS NULL
BEGIN
    CREATE TABLE shipments (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        order_id BIGINT NOT NULL UNIQUE,
        carrier_id BIGINT,
        carrier_name NVARCHAR(200),
        carrier_type NVARCHAR(50) DEFAULT 'third_party',
        tracking_number NVARCHAR(100),
        status NVARCHAR(50) NOT NULL DEFAULT 'pending',
        
        -- Shipping details
        shipping_from NVARCHAR(500),
        shipping_to NVARCHAR(500),
        weight DECIMAL(10,2),
        dimensions NVARCHAR(50), -- LxWxH
        package_count INT DEFAULT 1,
        
        -- Timestamps
        shipped_at DATETIME,
        estimated_delivery DATETIME,
        delivered_at DATETIME,
        failed_at DATETIME,
        failure_reason NVARCHAR(500),
        
        -- Additional info
        shipping_fee DECIMAL(18,2) DEFAULT 0,
        insurance_amount DECIMAL(18,2) DEFAULT 0,
        notes NVARCHAR(500),
        
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_shipments_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
        CONSTRAINT FK_shipments_carrier FOREIGN KEY (carrier_id) REFERENCES shipping_carriers(id)
    );

    CREATE INDEX IX_shipments_order_id ON shipments(order_id);
    CREATE INDEX IX_shipments_tracking_number ON shipments(tracking_number);
    CREATE INDEX IX_shipments_status ON shipments(status);
    CREATE INDEX IX_shipments_deleted_at ON shipments(deleted_at);

    PRINT '✓ Table [shipments] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [shipments] already exists.';
END
GO

-- ============================================
-- SHIPMENT TRACKING TABLE
-- ============================================
IF OBJECT_ID('dbo.shipment_tracking', 'U') IS NULL
BEGIN
    CREATE TABLE shipment_tracking (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        shipment_id BIGINT NOT NULL,
        status NVARCHAR(50) NOT NULL,
        location NVARCHAR(300),
        description NVARCHAR(1000) NOT NULL,
        occurred_at DATETIME NOT NULL,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_shipment_tracking_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_shipment_tracking_shipment_id ON shipment_tracking(shipment_id);
    CREATE INDEX IX_shipment_tracking_status ON shipment_tracking(status);
    CREATE INDEX IX_shipment_tracking_occurred_at ON shipment_tracking(occurred_at);

    PRINT '✓ Table [shipment_tracking] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [shipment_tracking] already exists.';
END
GO

-- ============================================
-- SEED DATA: Shipping Carriers
-- ============================================
PRINT '============================================';
PRINT 'Seeding shipping carriers...';
PRINT '============================================';

-- Insert default carriers if not exist
IF NOT EXISTS (SELECT 1 FROM shipping_carriers WHERE code = 'INTERNAL')
BEGIN
    INSERT INTO shipping_carriers (name, code, type, is_active) 
    VALUES ('Internal Shipping', 'INTERNAL', 'internal', 1);
    PRINT '✓ Added Internal Shipping carrier';
END

IF NOT EXISTS (SELECT 1 FROM shipping_carriers WHERE code = 'GHN')
BEGIN
    INSERT INTO shipping_carriers (name, code, type, website, is_active) 
    VALUES ('Giao Hàng Nhanh', 'GHN', 'third_party', 'https://ghn.vn', 1);
    PRINT '✓ Added GHN carrier';
END

IF NOT EXISTS (SELECT 1 FROM shipping_carriers WHERE code = 'GHTK')
BEGIN
    INSERT INTO shipping_carriers (name, code, type, website, is_active) 
    VALUES ('Giao Hàng Tiết Kiệm', 'GHTK', 'third_party', 'https://ghtk.vn', 1);
    PRINT '✓ Added GHTK carrier';
END

IF NOT EXISTS (SELECT 1 FROM shipping_carriers WHERE code = 'VIETTEL')
BEGIN
    INSERT INTO shipping_carriers (name, code, type, website, is_active) 
    VALUES ('Viettel Post', 'VIETTEL', 'third_party', 'https://viettelpost.com.vn', 1);
    PRINT '✓ Added Viettel Post carrier';
END

IF NOT EXISTS (SELECT 1 FROM shipping_carriers WHERE code = 'VNPOST')
BEGIN
    INSERT INTO shipping_carriers (name, code, type, website, is_active) 
    VALUES ('Vietnam Post', 'VNPOST', 'third_party', 'https://vnpost.vn', 1);
    PRINT '✓ Added Vietnam Post carrier';
END

PRINT '============================================';
PRINT 'Shipping Tables Migration Complete!';
PRINT '============================================';
