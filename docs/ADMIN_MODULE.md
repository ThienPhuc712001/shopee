# Admin Management System - Business Flow & Implementation

## PART 1 — Admin Business Flow

### Administrator Responsibilities

```
┌─────────────────────────────────────────────────────────────────┐
│                    ADMIN RESPONSIBILITIES                        │
└─────────────────────────────────────────────────────────────────┘

1. MANAGE USERS
   ┌─────────────────────────────────────────────────────────┐
   │ • View all user accounts                                │
   │ • Search and filter users                               │
   │ • View user details and order history                   │
   │ • Ban/suspend problematic users                         │
   │ • Reactivate suspended accounts                         │
   │ • Verify user identity                                  │
   │ • Handle user complaints                                │
   └─────────────────────────────────────────────────────────┘

2. APPROVE SELLERS
   ┌─────────────────────────────────────────────────────────┐
   │ • Review seller applications                            │
   │ • Verify business documents                             │
   │ • Approve/reject seller registration                    │
   │ • Suspend non-compliant sellers                         │
   │ • Monitor seller performance                            │
   │ • Handle seller disputes                                │
   └─────────────────────────────────────────────────────────┘

3. MODERATE PRODUCTS
   ┌─────────────────────────────────────────────────────────┐
   │ • Review new product listings                           │
   │ • Remove prohibited items                               │
   │ • Verify product authenticity                           │
   │ • Check pricing compliance                              │
   │ • Handle product reports                                │
   │ • Manage product categories                             │
   └─────────────────────────────────────────────────────────┘

4. MONITOR ORDERS
   ┌─────────────────────────────────────────────────────────┐
   │ • View all platform orders                              │
   │ • Monitor order fulfillment                             │
   │ • Handle order disputes                                 │
   │ • Process order cancellations                           │
   │ • Track delivery issues                                 │
   │ • Generate order reports                                │
   └─────────────────────────────────────────────────────────┘

5. HANDLE DISPUTES
   ┌─────────────────────────────────────────────────────────┐
   │ • Review customer complaints                            │
   │ • Mediate buyer-seller disputes                         │
   │ • Process refund requests                               │
   │ • Investigate fraud cases                               │
   │ • Enforce platform policies                             │
   │ • Document resolution outcomes                          │
   └─────────────────────────────────────────────────────────┘

6. MANAGE PROMOTIONS
   ┌─────────────────────────────────────────────────────────┐
   │ • Create platform-wide promotions                       │
   │ • Manage voucher codes                                  │
   │ • Set up flash sales                                    │
   │ • Monitor promotion performance                         │
   │ • Approve seller promotions                             │
   │ • Track promotion ROI                                   │
   └─────────────────────────────────────────────────────────┘

7. VIEW ANALYTICS
   ┌─────────────────────────────────────────────────────────┐
   │ • Sales dashboard                                       │
   │ • User growth metrics                                   │
   │ • Product performance                                   │
   │ • Revenue reports                                       │
   │ • Traffic analytics                                     │
   │ • Conversion rates                                      │
   └─────────────────────────────────────────────────────────┘
```

---

## PART 2 — Admin Roles

### Role Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                      ADMIN ROLES                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ SUPER_ADMIN                                                  │
│ ───────────                                                  │
│ Permissions: FULL ACCESS                                     │
│ • All admin permissions                                      │
│ • Create/delete admin accounts                               │
│ • System configuration                                       │
│ • Access to all analytics                                    │
│ • Override any decision                                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ ADMIN                                                        │
│ ─────                                                        │
│ Permissions: STANDARD ACCESS                                 │
│ • Manage users (ban, suspend, reactivate)                    │
│ • Approve/reject sellers                                     │
│ • Moderate products                                          │
│ • Process refunds                                            │
│ • View standard analytics                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ SUPPORT_AGENT                                                │
│ ─────────────                                                │
│ Permissions: LIMITED ACCESS                                  │
│ • View orders                                                │
│ • Handle basic disputes                                      │
│ • Respond to tickets                                         │
│ • Process standard refunds                                   │
│ • View basic analytics                                       │
└─────────────────────────────────────────────────────────────┘
```

### Permission Matrix

| Permission | Super Admin | Admin | Support Agent |
|------------|-------------|-------|---------------|
| Manage Admins | ✅ | ❌ | ❌ |
| System Settings | ✅ | ❌ | ❌ |
| Ban Users | ✅ | ✅ | ❌ |
| Approve Sellers | ✅ | ✅ | ❌ |
| Delete Products | ✅ | ✅ | ❌ |
| Process Refunds | ✅ | ✅ | ✅ |
| View All Orders | ✅ | ✅ | ✅ |
| View Analytics | ✅ | ✅ | ⚠️ Limited |
| Audit Logs | ✅ | ✅ | ❌ |

---

## PART 3 — Database Tables

### AdminUsers Table

```sql
CREATE TABLE [dbo].[AdminUsers] (
    [id]                    BIGINT         IDENTITY(1,1) PRIMARY KEY,
    [email]                 NVARCHAR(255)  NOT NULL UNIQUE,
    [password_hash]         NVARCHAR(255)  NOT NULL,
    [role_id]               INT            NOT NULL,
    [first_name]            NVARCHAR(100),
    [last_name]             NVARCHAR(100),
    [phone]                 NVARCHAR(20),
    [avatar_url]            NVARCHAR(500),
    [status]                NVARCHAR(20)   NOT NULL DEFAULT 'active',
    [last_login_at]         DATETIME,
    [last_login_ip]         NVARCHAR(45),
    [failed_login_attempts] INT            DEFAULT 0,
    [locked_until]          DATETIME,
    [created_at]            DATETIME       NOT NULL DEFAULT GETDATE(),
    [updated_at]            DATETIME       NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([role_id]) REFERENCES [AdminRoles]([id])
);

CREATE INDEX [IX_AdminUsers_Email] ON [AdminUsers]([email]);
CREATE INDEX [IX_AdminUsers_RoleID] ON [AdminUsers]([role_id]);
CREATE INDEX [IX_AdminUsers_Status] ON [AdminUsers]([status]);
```

### AdminRoles Table

```sql
CREATE TABLE [dbo].[AdminRoles] (
    [id]          INT           IDENTITY(1,1) PRIMARY KEY,
    [name]        NVARCHAR(50)  NOT NULL UNIQUE,
    [description] NVARCHAR(255),
    [permissions] NVARCHAR(MAX), -- JSON array
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE()
);

-- Default roles
INSERT INTO [AdminRoles] ([name], [description]) VALUES
('super_admin', 'Full system access'),
('admin', 'Standard admin access'),
('support_agent', 'Limited support access');
```

### AuditLogs Table

```sql
CREATE TABLE [dbo].[AuditLogs] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [admin_id]    BIGINT        NOT NULL,
    [action]      NVARCHAR(100) NOT NULL,
    [entity_type] NVARCHAR(50)  NOT NULL, -- user, product, order, etc.
    [entity_id]   BIGINT,
    [old_values]  NVARCHAR(MAX), -- JSON
    [new_values]  NVARCHAR(MAX), -- JSON
    [ip_address]  NVARCHAR(45),
    [user_agent]  NVARCHAR(500),
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([admin_id]) REFERENCES [AdminUsers]([id])
);

CREATE INDEX [IX_AuditLogs_AdminID] ON [AuditLogs]([admin_id]);
CREATE INDEX [IX_AuditLogs_Action] ON [AuditLogs]([action]);
CREATE INDEX [IX_AuditLogs_EntityType] ON [AuditLogs]([entity_type]);
CREATE INDEX [IX_AuditLogs_EntityID] ON [AuditLogs]([entity_id]);
CREATE INDEX [IX_AuditLogs_CreatedAt] ON [AuditLogs]([created_at]);
```

### SystemSettings Table

```sql
CREATE TABLE [dbo].[SystemSettings] (
    [id]          BIGINT        IDENTITY(1,1) PRIMARY KEY,
    [key]         NVARCHAR(100) NOT NULL UNIQUE,
    [value]       NVARCHAR(MAX) NOT NULL,
    [type]        NVARCHAR(50)  NOT NULL, -- string, number, boolean, json
    [description] NVARCHAR(500),
    [is_public]   BIT           DEFAULT 0,
    [updated_by]  BIGINT,
    [created_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    [updated_at]  DATETIME      NOT NULL DEFAULT GETDATE(),
    
    FOREIGN KEY ([updated_by]) REFERENCES [AdminUsers]([id])
);

CREATE INDEX [IX_SystemSettings_Key] ON [SystemSettings]([key]);
```

---

## PART 4-11 — Implementation

The complete implementation follows in the code files below.

### Files Created:

1. **Models** - AdminUser, AdminRole, AuditLog, SystemSetting
2. **Repository** - Admin repository with all CRUD operations
3. **Service** - Admin service with business logic
4. **Handler** - REST API handlers
5. **Routes** - Admin route definitions

### Key Features:

✅ **Role-Based Access Control** - Super Admin, Admin, Support Agent
✅ **Audit Logging** - All admin actions tracked
✅ **User Management** - Ban, suspend, reactivate users
✅ **Seller Management** - Approve/reject sellers
✅ **Product Moderation** - Delete inappropriate products
✅ **Order Management** - View all orders, process refunds
✅ **Analytics Dashboard** - Platform metrics and reports
✅ **System Settings** - Platform configuration
✅ **Security** - JWT auth, IP tracking, action logging

---

## Admin Action Audit Trail

```
┌─────────────────────────────────────────────────────────────────┐
│                    AUDIT LOGGING FLOW                            │
└─────────────────────────────────────────────────────────────────┘

ADMIN ACTION
    │
    ▼
┌─────────────────┐
│ Authenticate    │
│ Admin User      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Check           │
│ Permissions     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Execute         │
│ Action          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Log to          │
│ AuditLogs       │
│ - Admin ID      │
│ - Action        │
│ - Entity        │
│ - IP Address    │
│ - Timestamp     │
└─────────────────┘
```

---

This admin system is production-ready with full RBAC, audit logging, and platform management capabilities.
