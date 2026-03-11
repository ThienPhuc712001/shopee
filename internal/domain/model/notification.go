package model

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeOrder      NotificationType = "order"
	NotificationTypePayment    NotificationType = "payment"
	NotificationTypeShipping   NotificationType = "shipping"
	NotificationTypePromotion  NotificationType = "promotion"
	NotificationTypeSystem     NotificationType = "system"
	NotificationTypeReview     NotificationType = "review"
)

// NotificationPriority represents the priority level of a notification
type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// Notification represents a user notification
type Notification struct {
	ID        uint               `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID    uint               `gorm:"type:int;not null;index" json:"user_id"`
	Title     string             `gorm:"type:varchar(255);not null" json:"title"`
	Message   string             `gorm:"type:text;not null" json:"message"`
	Type      NotificationType   `gorm:"type:varchar(50);not null;index" json:"type"`
	Priority  NotificationPriority `gorm:"type:varchar(20);default:'normal'" json:"priority"`
	IsRead    bool               `gorm:"type:bit;default:false;index" json:"is_read"`
	ReadAt    *time.Time         `gorm:"type:datetime" json:"read_at"`
	Data      string             `gorm:"type:text" json:"data"` // JSON data for additional info
	ActionURL string             `gorm:"type:varchar(500)" json:"action_url"`
	ImageURL  string             `gorm:"type:varchar(500)" json:"image_url"`
	CreatedAt time.Time          `gorm:"type:datetime;not null;index" json:"created_at"`
	UpdatedAt time.Time          `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt *time.Time         `gorm:"index" json:"-"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}

// NotificationChannel represents a notification delivery channel
type NotificationChannel string

const (
	NotificationChannelInApp NotificationChannel = "in_app"
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelSMS   NotificationChannel = "sms"
	NotificationChannelPush  NotificationChannel = "push"
)

// NotificationPreference represents user's notification preferences
type NotificationPreference struct {
	ID                uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID            uint      `gorm:"type:int;not null;uniqueIndex" json:"user_id"`
	EmailEnabled      bool      `gorm:"type:bit;default:true" json:"email_enabled"`
	SMSEnabled        bool      `gorm:"type:bit;default:false" json:"sms_enabled"`
	PushEnabled       bool      `gorm:"type:bit;default:true" json:"push_enabled"`
	OrderEnabled      bool      `gorm:"type:bit;default:true" json:"order_enabled"`
	PaymentEnabled    bool      `gorm:"type:bit;default:true" json:"payment_enabled"`
	ShippingEnabled   bool      `gorm:"type:bit;default:true" json:"shipping_enabled"`
	PromotionEnabled  bool      `gorm:"type:bit;default:true" json:"promotion_enabled"`
	SystemEnabled     bool      `gorm:"type:bit;default:true" json:"system_enabled"`
	CreatedAt         time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:datetime;not null" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for NotificationPreference
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// NotificationLog tracks notification delivery attempts
type NotificationLog struct {
	ID              uint                `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	NotificationID  uint                `gorm:"type:int;not null;index" json:"notification_id"`
	Channel         NotificationChannel `gorm:"type:varchar(50);not null" json:"channel"`
	Status          string              `gorm:"type:varchar(50);not null" json:"status"` // sent, delivered, failed
	ErrorMessage    string              `gorm:"type:text" json:"error_message"`
	SentAt          *time.Time          `gorm:"type:datetime" json:"sent_at"`
	DeliveredAt     *time.Time          `gorm:"type:datetime" json:"delivered_at"`
	CreatedAt       time.Time           `gorm:"type:datetime;not null" json:"created_at"`

	// Relationships
	Notification *Notification `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
}

// TableName specifies the table name for NotificationLog
func (NotificationLog) TableName() string {
	return "notification_logs"
}

// NotificationInput represents input for creating a notification
type NotificationInput struct {
	UserID    uint               `json:"user_id" binding:"required"`
	Title     string             `json:"title" binding:"required"`
	Message   string             `json:"message" binding:"required"`
	Type      NotificationType   `json:"type"`
	Priority  NotificationPriority `json:"priority"`
	Data      string             `json:"data"`
	ActionURL string             `json:"action_url"`
	ImageURL  string             `json:"image_url"`
}

// NotificationBatchInput represents input for creating multiple notifications
type NotificationBatchInput struct {
	UserIDs   []uint             `json:"user_ids" binding:"required"`
	Title     string             `json:"title" binding:"required"`
	Message   string             `json:"message" binding:"required"`
	Type      NotificationType   `json:"type"`
	Priority  NotificationPriority `json:"priority"`
	Data      string             `json:"data"`
	ActionURL string             `json:"action_url"`
}

// NotificationFilter represents filters for notification queries
type NotificationFilter struct {
	UserID    uint               `json:"user_id"`
	Type      *NotificationType  `json:"type"`
	IsRead    *bool              `json:"is_read"`
	Priority  *NotificationPriority `json:"priority"`
	StartDate *time.Time         `json:"start_date"`
	EndDate   *time.Time         `json:"end_date"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalNotifications int64 `json:"total_notifications"`
	UnreadCount        int64 `json:"unread_count"`
	ReadCount          int64 `json:"read_count"`
	OrderNotifications int64 `json:"order_notifications"`
	PaymentNotifications int64 `json:"payment_notifications"`
	ShippingNotifications int64 `json:"shipping_notifications"`
	PromotionNotifications int64 `json:"promotion_notifications"`
}

// NotificationSummary represents a summary of notifications for a user
type NotificationSummary struct {
	TotalCount   int64                    `json:"total_count"`
	UnreadCount  int64                    `json:"unread_count"`
	RecentNotifications []NotificationView `json:"recent_notifications"`
}

// NotificationView represents a notification view for API responses
type NotificationView struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	Type      NotificationType   `json:"type"`
	Priority  NotificationPriority `json:"priority"`
	IsRead    bool               `json:"is_read"`
	CreatedAt time.Time          `json:"created_at"`
	ActionURL string             `json:"action_url"`
	ImageURL  string             `json:"image_url"`
}

// ToView converts a Notification to NotificationView
func (n *Notification) ToView() NotificationView {
	return NotificationView{
		ID:        n.ID,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type,
		Priority:  n.Priority,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
		ActionURL: n.ActionURL,
		ImageURL:  n.ImageURL,
	}
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}

// IsUrgent checks if the notification is urgent
func (n *Notification) IsUrgent() bool {
	return n.Priority == NotificationPriorityUrgent || n.Priority == NotificationPriorityHigh
}

// GetNotificationTypeLabel returns a human-readable label for the notification type
func (n *Notification) GetNotificationTypeLabel() string {
	labels := map[NotificationType]string{
		NotificationTypeOrder:     "Order Update",
		NotificationTypePayment:   "Payment Update",
		NotificationTypeShipping:  "Shipping Update",
		NotificationTypePromotion: "Promotion",
		NotificationTypeSystem:    "System",
		NotificationTypeReview:    "Review",
	}
	if label, ok := labels[n.Type]; ok {
		return label
	}
	return "Notification"
}

// EmailTemplate represents an email notification template
type EmailTemplate struct {
	TemplateName string            `json:"template_name"`
	Subject      string            `json:"subject"`
	Body         string            `json:"body"`
	Data         map[string]string `json:"data"`
}

// EmailConfig represents SMTP email configuration
type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	FromName string `json:"from_name"`
	FromEmail string `json:"from_email"`
	UseTLS   bool   `json:"use_tls"`
}

// EmailRequest represents an email sending request
type EmailRequest struct {
	To          []string          `json:"to"`
	ToName      string            `json:"to_name"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTML        bool              `json:"html"`
	CC          []string          `json:"cc"`
	BCC         []string          `json:"bcc"`
	Attachments []string          `json:"attachments"`
	Data        map[string]string `json:"data"`
}
