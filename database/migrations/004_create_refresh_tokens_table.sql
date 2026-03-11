-- ============================================
-- REFRESH TOKENS TABLE
-- ============================================
-- Purpose: Store JWT refresh tokens for token refresh flow
-- Security: Tokens are stored with expiration and revocation support
-- ============================================

USE ecommerce;
GO

IF OBJECT_ID('dbo.refresh_tokens', 'U') IS NULL
BEGIN
    CREATE TABLE refresh_tokens (
        id BIGINT IDENTITY(1,1) PRIMARY KEY,
        user_id BIGINT NOT NULL,
        token NVARCHAR(500) NOT NULL UNIQUE,
        expires_at DATETIME NOT NULL,
        revoked BIT DEFAULT 0,
        revoked_at DATETIME,
        created_at DATETIME NOT NULL DEFAULT GETDATE(),

        CONSTRAINT FK_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IX_refresh_tokens_user_id ON refresh_tokens(user_id);
    CREATE INDEX IX_refresh_tokens_token ON refresh_tokens(token);
    CREATE INDEX IX_refresh_tokens_expires_at ON refresh_tokens(expires_at);
    CREATE INDEX IX_refresh_tokens_revoked ON refresh_tokens(revoked);

    PRINT '✓ Table [refresh_tokens] created successfully.';
END
ELSE
BEGIN
    PRINT '○ Table [refresh_tokens] already exists.';
END
GO

-- ============================================
-- UPDATE USERS TABLE FOR LOGIN SECURITY
-- ============================================
-- Add columns for failed login tracking and account lockout

IF COL_LENGTH('dbo.users', 'failed_login_attempts') IS NULL
BEGIN
    ALTER TABLE users ADD failed_login_attempts INT DEFAULT 0;
    PRINT '✓ Added [failed_login_attempts] column to [users]';
END

IF COL_LENGTH('dbo.users', 'locked_until') IS NULL
BEGIN
    ALTER TABLE users ADD locked_until DATETIME;
    PRINT '✓ Added [locked_until] column to [users]';
END
GO

-- ============================================
-- CLEANUP EXPIRED REFRESH TOKENS
-- ============================================
-- This can be run periodically to clean up old tokens

/*
DELETE FROM refresh_tokens 
WHERE expires_at < GETDATE() 
   OR (revoked = 1 AND revoked_at < DATEADD(DAY, -7, GETDATE()));
*/
GO
