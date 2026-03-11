# Notification System Documentation

## Overview

The Notification System is a comprehensive module for managing user notifications in the e-commerce platform. It supports in-app notifications, email notifications, and provides a flexible framework for adding SMS and push notifications.

---

## Table of Contents

1. [Notification Business Flow](#notification-business-flow)
2. [Notification Types](#notification-types)
3. [Database Schema](#database-schema)
4. [Go Models](#go-models)
5. [API Endpoints](#api-endpoints)
6. [Order Integration](#order-integration)
7. [Email Notifications](#email-notifications)
8. [Frontend Integration](#frontend-integration)
9. [Performance & Pagination](#performance--pagination)
10. [Usage Examples](#usage-examples)

---

## Notification Business Flow

### How Notifications Work

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   USER      │     │   SYSTEM    │     │ NOTIFICATION│
│  PERFORMS   │────▶│  GENERATES  │────▶│   CREATED   │
│   ACTION    │     │    EVENT    │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   USER      │◀────│   EMAIL     │◀────│  SAVED TO   │
│  RECEIVES   │     │   SENT      │     │  DATABASE   │
│ NOTIFICATION│     │ (OPTIONAL)  │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Step-by-Step Flow

| Step | Actor | Action | Description |
|------|-------|--------|-------------|
| 1 | User | Performs action | Places order, makes payment, etc. |
| 2 | System | Generates event | Order created, payment successful |
| 3 | System | Creates notification | Notification object created |
| 4 | Database | Saves notification | Stored in `notifications` table |
| 5 | Email Service | Sends email (optional) | Based on user preferences |
| 6 | User | Receives notification | Via in-app notification or email |

---

## Notification Types

### Supported Types

| Type | Description | Examples |
|------|-------------|----------|
| `order` | Order-related notifications | Order created, confirmed, shipped, delivered |
| `payment` | Payment-related notifications | Payment successful, failed, refunded |
| `shipping` | Shipping-related notifications | Order shipped, out for delivery, delivered |
| `promotion` | Marketing notifications | Discount codes, special offers, sales |
| `system` | System notifications | Account updates, security alerts |
| `review` | Review-related notifications | Review approved, helpful count increased |

### Notification Priorities

| Priority | Description | Use Cases |
|----------|-------------|-----------|
| `low` | Low priority | General updates, newsletters |
| `normal` | Normal priority | Standard notifications |
| `high` | High priority | Order updates, payment confirmations |
| `urgent` | Urgent | Security alerts, critical issues |

---

## Database Schema

### Notifications Table

```sql
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

    CONSTRAINT FK_notifications_user FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Notification Preferences Table

```sql
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

    CONSTRAINT FK_notification_preferences_user FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Notification Logs Table

```sql
CREATE TABLE notification_logs (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    notification_id BIGINT NOT NULL,
    channel NVARCHAR(50) NOT NULL, -- 'in_app', 'email', 'sms', 'push'
    status NVARCHAR(50) NOT NULL, -- 'sent', 'delivered', 'failed'
    error_message NVARCHAR(MAX),
    sent_at DATETIME,
    delivered_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_notification_logs_notification FOREIGN KEY (notification_id) REFERENCES notifications(id)
);
```

### Relationships

```
users (1) ──────< notifications (>
     │
     └────── (1) ──────< notification_preferences (>
     
notifications (1) ──────< notification_logs (>
```

---

## Go Models

### Notification

```go
type Notification struct {
    ID        uint               `gorm:"primaryKey" json:"id"`
    UserID    uint               `gorm:"not null;index" json:"user_id"`
    Title     string             `gorm:"varchar(255);not null" json:"title"`
    Message   string             `gorm:"text;not null" json:"message"`
    Type      NotificationType   `gorm:"varchar(50);not null" json:"type"`
    Priority  NotificationPriority `gorm:"varchar(20);default:'normal'" json:"priority"`
    IsRead    bool               `gorm:"default:false;index" json:"is_read"`
    ReadAt    *time.Time         `gorm:"type:datetime" json:"read_at"`
    Data      string             `gorm:"text" json:"data"` // JSON data
    ActionURL string             `gorm:"varchar(500)" json:"action_url"`
    ImageURL  string             `gorm:"varchar(500)" json:"image_url"`
    CreatedAt time.Time          `gorm:"not null" json:"created_at"`
    
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

### NotificationPreference

```go
type NotificationPreference struct {
    ID               uint      `gorm:"primaryKey" json:"id"`
    UserID           uint      `gorm:"not null;uniqueIndex" json:"user_id"`
    EmailEnabled     bool      `gorm:"default:true" json:"email_enabled"`
    SMSEnabled       bool      `gorm:"default:false" json:"sms_enabled"`
    PushEnabled      bool      `gorm:"default:true" json:"push_enabled"`
    OrderEnabled     bool      `gorm:"default:true" json:"order_enabled"`
    PaymentEnabled   bool      `gorm:"default:true" json:"payment_enabled"`
    ShippingEnabled  bool      `gorm:"default:true" json:"shipping_enabled"`
    PromotionEnabled bool      `gorm:"default:true" json:"promotion_enabled"`
    SystemEnabled    bool      `gorm:"default:true" json:"system_enabled"`
    
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

---

## API Endpoints

### User Endpoints

#### Get User Notifications
```http
GET /api/notifications?page=1&limit=20&type=order&is_read=false
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": 1,
        "title": "Order Created",
        "message": "Your order ORD-20240315-00001 has been created successfully.",
        "type": "order",
        "priority": "high",
        "is_read": false,
        "created_at": "2024-03-15T10:30:00Z",
        "action_url": "/orders/ORD-20240315-00001"
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "limit": 20,
      "total_pages": 3
    }
  }
}
```

#### Get Notification Summary
```http
GET /api/notifications/summary
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "total_count": 50,
    "unread_count": 5,
    "recent_notifications": [
      {
        "id": 1,
        "title": "Order Created",
        "type": "order",
        "is_read": false,
        "created_at": "2024-03-15T10:30:00Z"
      }
    ]
  }
}
```

#### Get Unread Count
```http
GET /api/notifications/unread-count
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "unread_count": 5
  }
}
```

#### Mark Notification as Read
```http
PUT /api/notifications/:id/read
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "message": "Notification marked as read"
}
```

#### Mark All as Read
```http
PUT /api/notifications/read-all
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "message": "All notifications marked as read"
}
```

#### Delete Notification
```http
DELETE /api/notifications/:id
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "message": "Notification deleted successfully"
}
```

#### Get Notification Preferences
```http
GET /api/notifications/preferences
Authorization: Bearer {token}

Response (200):
{
  "success": true,
  "data": {
    "email_enabled": true,
    "sms_enabled": false,
    "push_enabled": true,
    "order_enabled": true,
    "payment_enabled": true,
    "shipping_enabled": true,
    "promotion_enabled": true,
    "system_enabled": true
  }
}
```

#### Update Notification Preferences
```http
PUT /api/notifications/preferences
Authorization: Bearer {token}
Content-Type: application/json

Request:
{
  "email_enabled": true,
  "sms_enabled": false,
  "push_enabled": true,
  "order_enabled": true,
  "payment_enabled": true,
  "shipping_enabled": true,
  "promotion_enabled": false,
  "system_enabled": true
}

Response (200):
{
  "success": true,
  "message": "Notification preferences updated successfully",
  "data": { ...preference object... }
}
```

---

### Admin Endpoints

#### Create Notification
```http
POST /api/admin/notifications
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "user_id": 123,
  "title": "Special Offer",
  "message": "You have received a special discount!",
  "type": "promotion",
  "priority": "normal",
  "data": "{\"promo_code\": \"SAVE20\"}",
  "action_url": "/promotions/save20"
}

Response (201):
{
  "success": true,
  "message": "Notification created successfully",
  "data": { ...notification object... }
}
```

#### Create Batch Notifications
```http
POST /api/admin/notifications/batch
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "user_ids": [1, 2, 3, 4, 5],
  "title": "System Maintenance",
  "message": "Scheduled maintenance on March 20, 2024",
  "type": "system",
  "priority": "high"
}

Response (201):
{
  "success": true,
  "message": "Notifications created successfully",
  "data": {
    "count": 5
  }
}
```

#### Send Promotion Notification
```http
POST /api/admin/notifications/promotion
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "user_ids": [1, 2, 3, 4, 5],
  "promo_title": "Flash Sale",
  "promo_code": "FLASH50",
  "discount_percent": 50
}

Response (200):
{
  "success": true,
  "message": "Promotion notifications sent successfully",
  "data": {
    "sent_count": 5
  }
}
```

#### Cleanup Old Notifications
```http
POST /api/admin/notifications/cleanup
Authorization: Bearer {token}
Role: admin
Content-Type: application/json

Request:
{
  "days": 90
}

Response (200):
{
  "success": true,
  "message": "Old notifications cleaned up successfully",
  "data": {
    "deleted_count": 1500
  }
}
```

---

## Order Integration

### How Order System Triggers Notifications

```
┌─────────────────────────────────────────────────────────────┐
│            ORDER NOTIFICATION TRIGGERS                       │
├─────────────────────────────────────────────────────────────┤
│  1. Order Created → Send "Order Created" notification       │
│  2. Order Confirmed → Send "Order Confirmed" notification   │
│  3. Order Shipped → Send "Order Shipped" notification       │
│  4. Order Delivered → Send "Order Delivered" notification   │
│  5. Order Cancelled → Send "Order Cancelled" notification   │
└─────────────────────────────────────────────────────────────┘
```

### Code Integration

```go
// In order_service_enhanced.go
func (s *orderServiceEnhanced) CheckoutCart(...) (*model.Order, error) {
    // ... create order ...
    
    // Send order notification
    if s.notifSvc != nil {
        go s.notifSvc.SendOrderNotification(
            context.Background(), 
            userID, 
            order.OrderNumber, 
            "created", 
            order.TotalAmount,
        )
    }
    
    return order, nil
}

func (s *orderServiceEnhanced) CompleteOrder(orderID uint) (*model.Order, error) {
    // ... complete order ...
    
    // Send order delivered notification
    if s.notifSvc != nil {
        go s.notifSvc.SendOrderNotification(
            context.Background(), 
            order.UserID, 
            order.OrderNumber, 
            "delivered", 
            order.TotalAmount,
        )
    }
    
    return order, nil
}
```

### Payment Integration

```go
// In payment service
func (s *paymentService) ProcessPayment(...) error {
    // ... process payment ...
    
    // Send payment notification
    if s.notifSvc != nil {
        go s.notifSvc.SendPaymentNotification(
            context.Background(),
            userID,
            orderNumber,
            amount,
            "success",
            paymentMethod,
        )
    }
    
    return nil
}
```

---

## Email Notifications

### SMTP Configuration

```go
type EmailConfig struct {
    Host      string // e.g., "smtp.gmail.com"
    Port      int    // e.g., 587
    Username  string // SMTP username
    Password  string // SMTP password
    FromName  string // e.g., "E-Commerce Store"
    FromEmail string // e.g., "noreply@example.com"
    UseTLS    bool   // true for TLS
}
```

### Email Templates

The system includes pre-built HTML email templates for:

1. **Order Confirmation** - Sent when order is created
2. **Payment Receipt** - Sent when payment is successful
3. **Shipping Update** - Sent when order is shipped
4. **Promotion Email** - Sent for promotional campaigns
5. **Password Reset** - Sent for password reset requests
6. **Welcome Email** - Sent to new users

### Example Email Configuration

```bash
# .env file
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_FROM_NAME=E-Commerce Store
EMAIL_FROM_EMAIL=noreply@example.com
EMAIL_USE_TLS=true
```

---

## Frontend Integration

### React Notification Bell Component

```jsx
import React, { useState, useEffect } from 'react';

function NotificationBell() {
  const [unreadCount, setUnreadCount] = useState(0);
  const [notifications, setNotifications] = useState([]);
  const [isOpen, setIsOpen] = useState(false);

  // Fetch unread count
  useEffect(() => {
    fetchUnreadCount();
    const interval = setInterval(fetchUnreadCount, 30000); // Poll every 30s
    return () => clearInterval(interval);
  }, []);

  const fetchUnreadCount = async () => {
    const response = await fetch('/api/notifications/unread-count', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await response.json();
    setUnreadCount(data.data.unread_count);
  };

  const fetchNotifications = async () => {
    const response = await fetch('/api/notifications?limit=10', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await response.json();
    setNotifications(data.data.notifications);
  };

  const markAsRead = async (id) => {
    await fetch(`/api/notifications/${id}/read`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token}` }
    });
    fetchUnreadCount();
    fetchNotifications();
  };

  return (
    <div className="notification-bell">
      <button onClick={() => setIsOpen(!isOpen)}>
        🔔
        {unreadCount > 0 && (
          <span className="badge">{unreadCount}</span>
        )}
      </button>
      
      {isOpen && (
        <div className="notification-dropdown">
          <div className="header">
            <h3>Notifications</h3>
            <button onClick={markAllAsRead}>Mark all as read</button>
          </div>
          
          <div className="notifications-list">
            {notifications.map(notification => (
              <div 
                key={notification.id} 
                className={`notification-item ${notification.is_read ? 'read' : 'unread'}`}
                onClick={() => markAsRead(notification.id)}
              >
                <div className="notification-content">
                  <h4>{notification.title}</h4>
                  <p>{notification.message}</p>
                  <span className="time">{formatTime(notification.created_at)}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
```

### Notification Badge Styles

```css
.notification-bell {
  position: relative;
  display: inline-block;
}

.notification-bell .badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #ff4444;
  color: white;
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 10px;
  font-weight: bold;
}

.notification-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 350px;
  max-height: 400px;
  background: white;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  border-radius: 8px;
  overflow-y: auto;
  z-index: 1000;
}

.notification-item {
  padding: 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background 0.2s;
}

.notification-item:hover {
  background: #f5f5f5;
}

.notification-item.unread {
  background: #e3f2fd;
  border-left: 3px solid #2196f3;
}

.notification-item.read {
  opacity: 0.7;
}
```

---

## Performance & Pagination

### Pagination Implementation

All notification endpoints support pagination to limit the number of records returned:

```http
GET /api/notifications?page=1&limit=20
```

**Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

### Response Format

```json
{
  "success": true,
  "data": {
    "notifications": [...],
    "pagination": {
      "total": 150,
      "page": 1,
      "limit": 20,
      "total_pages": 8
    }
  }
}
```

### Database Indexes

The following indexes are created for optimal performance:

```sql
CREATE INDEX IX_notifications_user_id ON notifications(user_id);
CREATE INDEX IX_notifications_type ON notifications(type);
CREATE INDEX IX_notifications_is_read ON notifications(is_read);
CREATE INDEX IX_notifications_created_at ON notifications(created_at);
```

### Cleanup Old Notifications

To prevent database bloat, old notifications can be cleaned up:

```http
POST /api/admin/notifications/cleanup
{
  "days": 90
}
```

This deletes notifications older than 90 days.

---

## Usage Examples

### Example 1: Get User Notifications

```bash
curl -X GET "http://localhost:8080/api/notifications?page=1&limit=20" \
  -H "Authorization: Bearer {user_token}"
```

### Example 2: Mark Notification as Read

```bash
curl -X PUT "http://localhost:8080/api/notifications/123/read" \
  -H "Authorization: Bearer {user_token}"
```

### Example 3: Update Notification Preferences

```bash
curl -X PUT "http://localhost:8080/api/notifications/preferences" \
  -H "Authorization: Bearer {user_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "email_enabled": true,
    "promotion_enabled": false
  }'
```

### Example 4: Send Promotion (Admin)

```bash
curl -X POST "http://localhost:8080/api/admin/notifications/promotion" \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "user_ids": [1, 2, 3, 4, 5],
    "promo_title": "Summer Sale",
    "promo_code": "SUMMER30",
    "discount_percent": 30
  }'
```

### Example 5: Create Notification (Admin)

```bash
curl -X POST "http://localhost:8080/api/admin/notifications" \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "Account Verified",
    "message": "Your account has been verified successfully.",
    "type": "system",
    "priority": "normal",
    "action_url": "/account/settings"
  }'
```

---

## Environment Variables

```bash
# Email/SMTP Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_FROM_NAME=E-Commerce Store
EMAIL_FROM_EMAIL=noreply@example.com
EMAIL_USE_TLS=true
```

---

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Notifications not appearing | Database not migrated | Run migration script |
| Email not sent | SMTP not configured | Check .env settings |
| Too many notifications | No cleanup | Run cleanup endpoint |
| Slow performance | Missing indexes | Check database indexes |

### Debug Queries

```sql
-- Check user notifications
SELECT * FROM notifications WHERE user_id = 1 ORDER BY created_at DESC;

-- Check unread count
SELECT COUNT(*) FROM notifications WHERE user_id = 1 AND is_read = 0;

-- Check notification preferences
SELECT * FROM notification_preferences WHERE user_id = 1;

-- Check email delivery logs
SELECT * FROM notification_logs WHERE channel = 'email' ORDER BY created_at DESC;
```

---

## License

Internal use only - E-Commerce Platform
