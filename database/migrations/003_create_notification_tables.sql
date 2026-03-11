-- ============================================
-- Migration: Create Notification Tables
-- E-Commerce Database - Notification System
-- SQL Server
-- ============================================

USE ecommerce;
GO

PRINT '============================================';
PRINT 'Creating Notification Tables...';
PRINT '============================================';

-- ============================================
-- NOTIFICATIONS TABLE
-- ============================================
IF OBJECT_ID('dbo.notifications', 'U') IS NULL
BEGIN
    CREATE TABLE notifications (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL,
        title NVARCHAR(255) NOT NULL,
        message NVARCHAR(MAX) NOT NULL,
        type NVARCHAR(50) NOT NULL, -- 'order', 'payment', 'shipping', 'promotion', 'system', 'review'
        priority NVARCHAR(20) DEFAULT 'normal', -- 'low', 'normal', 'high', 'urgent'
        is_read BIT DEFAULT 0,
        read_at DATETIME,
        data NVARCHAR(MAX), -- JSON data for additional info
        action_url NVARCHAR(500),
        image_url NVARCHAR(500),
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),
        deleted_at DATETIME,

        CONSTRAINT FK_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_notifications_user_id ON notifications(user_id);
    CREATE INDEX IX_notifications_type ON notifications(type);
    CREATE INDEX IX_notifications_is_read ON notifications(is_read);
    CREATE INDEX IX_notifications_created_at ON notifications(created_at);
    CREATE INDEX IX_notifications_deleted_at ON notifications(deleted_at);

    PRINT '✓ Table [notifications] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [notifications] already exists.';
END
GO

-- ============================================
-- NOTIFICATION PREFERENCES TABLE
-- ============================================
IF OBJECT_ID('dbo.notification_preferences', 'U') IS NULL
BEGIN
    CREATE TABLE notification_preferences (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL UNIQUE,
        email_enabled BIT DEFAULT 1,
        sms_enabled BIT DEFAULT 0,
        push_enabled BIT DEFAULT 1,
        order_enabled BIT DEFAULT 1,
        payment_enabled BIT DEFAULT 1,
        shipping_enabled BIT DEFAULT 1,
        promotion_enabled BIT DEFAULT 1,
        system_enabled BIT DEFAULT 1,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),
        updated_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_notification_preferences_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_notification_preferences_user_id ON notification_preferences(user_id);

    PRINT '✓ Table [notification_preferences] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [notification_preferences] already exists.';
END
GO

-- ============================================
-- NOTIFICATION LOGS TABLE
-- ============================================
IF OBJECT_ID('dbo.notification_logs', 'U') IS NULL
BEGIN
    CREATE TABLE notification_logs (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        notification_id BIGINT NOT NULL,
        channel NVARCHAR(50) NOT NULL, -- 'in_app', 'email', 'sms', 'push'
        status NVARCHAR(50) NOT NULL, -- 'sent', 'delivered', 'failed'
        error_message NVARCHAR(MAX),
        sent_at DATETIME,
        delivered_at DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_notification_logs_notification FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_notification_logs_notification_id ON notification_logs(notification_id);
    CREATE INDEX IX_notification_logs_channel ON notification_logs(channel);
    CREATE INDEX IX_notification_logs_status ON notification_logs(status);
    CREATE INDEX IX_notification_logs_created_at ON notification_logs(created_at);

    PRINT '✓ Table [notification_logs] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [notification_logs] already exists.';
END
GO

PRINT '============================================';
PRINT 'Notification Tables Migration Complete!';
PRINT '============================================';
