package repository

import (
	"ecommerce/internal/domain/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// ==================== NOTIFICATION ====================

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// CreateNotifications creates multiple notifications in a transaction
func (r *NotificationRepository) CreateNotifications(ctx context.Context, notifications []*model.Notification) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	for _, notification := range notifications {
		if err := tx.Create(notification).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetNotificationByID retrieves a notification by ID
func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id uint) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.WithContext(ctx).First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetUserNotifications retrieves notifications for a user with pagination
func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID uint, filter model.NotificationFilter) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ?", userID)

	// Apply filters
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.IsRead != nil {
		query = query.Where("is_read = ?", *filter.IsRead)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Get notifications ordered by created_at DESC
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, total, err
}

// GetUnreadCount retrieves the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uint, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// DeleteNotification deletes a notification (soft delete)
func (r *NotificationRepository) DeleteNotification(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Notification{}).Error
}

// DeleteNotificationsOlderThan deletes notifications older than a specific date
func (r *NotificationRepository) DeleteNotificationsOlderThan(ctx context.Context, userID uint, date time.Time) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND created_at < ?", userID, date).
		Delete(&model.Notification{}).Error
}

// GetNotificationsByType retrieves notifications by type
func (r *NotificationRepository) GetNotificationsByType(ctx context.Context, notificationType model.NotificationType, limit int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.WithContext(ctx).
		Where("type = ?", notificationType).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetNotificationStats retrieves notification statistics for a user
func (r *NotificationRepository) GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	stats := &model.NotificationStats{}

	// Total notifications
	if err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalNotifications).Error; err != nil {
		return nil, err
	}

	// Unread count
	if err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&stats.UnreadCount).Error; err != nil {
		return nil, err
	}

	// Read count
	stats.ReadCount = stats.TotalNotifications - stats.UnreadCount

	// Count by type
	r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND type = ?", userID, model.NotificationTypeOrder).
		Count(&stats.OrderNotifications)

	r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND type = ?", userID, model.NotificationTypePayment).
		Count(&stats.PaymentNotifications)

	r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND type = ?", userID, model.NotificationTypeShipping).
		Count(&stats.ShippingNotifications)

	r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND type = ?", userID, model.NotificationTypePromotion).
		Count(&stats.PromotionNotifications)

	return stats, nil
}

// GetRecentNotifications retrieves recent notifications for a user
func (r *NotificationRepository) GetRecentNotifications(ctx context.Context, userID uint, limit int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// ==================== NOTIFICATION PREFERENCE ====================

// GetNotificationPreference retrieves notification preference for a user
func (r *NotificationRepository) GetNotificationPreference(ctx context.Context, userID uint) (*model.NotificationPreference, error) {
	var preference model.NotificationPreference
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&preference).Error
	if err != nil {
		return nil, err
	}
	return &preference, nil
}

// CreateNotificationPreference creates notification preference for a user
func (r *NotificationRepository) CreateNotificationPreference(ctx context.Context, preference *model.NotificationPreference) error {
	return r.db.WithContext(ctx).Create(preference).Error
}

// UpdateNotificationPreference updates notification preference for a user
func (r *NotificationRepository) UpdateNotificationPreference(ctx context.Context, preference *model.NotificationPreference) error {
	return r.db.WithContext(ctx).Save(preference).Error
}

// GetOrCreateNotificationPreference gets or creates notification preference for a user
func (r *NotificationRepository) GetOrCreateNotificationPreference(ctx context.Context, userID uint) (*model.NotificationPreference, error) {
	preference, err := r.GetNotificationPreference(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		// Create default preference
		preference = &model.NotificationPreference{
			UserID:           userID,
			EmailEnabled:     true,
			PushEnabled:      true,
			SMSEnabled:       false,
			OrderEnabled:     true,
			PaymentEnabled:   true,
			ShippingEnabled:  true,
			PromotionEnabled: true,
			SystemEnabled:    true,
		}
		if err := r.CreateNotificationPreference(ctx, preference); err != nil {
			return nil, err
		}
		return preference, nil
	}
	if err != nil {
		return nil, err
	}
	return preference, nil
}

// ==================== NOTIFICATION LOG ====================

// CreateNotificationLog creates a notification log entry
func (r *NotificationRepository) CreateNotificationLog(ctx context.Context, log *model.NotificationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetNotificationLogsByNotification retrieves logs for a notification
func (r *NotificationRepository) GetNotificationLogsByNotification(ctx context.Context, notificationID uint) ([]model.NotificationLog, error) {
	var logs []model.NotificationLog
	err := r.db.WithContext(ctx).
		Where("notification_id = ?", notificationID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// UpdateNotificationLogStatus updates the status of a notification log
func (r *NotificationRepository) UpdateNotificationLogStatus(ctx context.Context, logID uint, status string, errorMessage string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMessage != "" {
		updates["error_message"] = errorMessage
	}
	if status == "sent" || status == "delivered" {
		now := time.Now()
		if status == "sent" {
			updates["sent_at"] = now
		} else {
			updates["delivered_at"] = now
		}
	}
	return r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Where("id = ?", logID).
		Updates(updates).Error
}

// GetNotificationDeliveryStats retrieves delivery statistics
type NotificationDeliveryStats struct {
	TotalSent       int64 `json:"total_sent"`
	TotalDelivered  int64 `json:"total_delivered"`
	TotalFailed     int64 `json:"total_failed"`
	EmailSent       int64 `json:"email_sent"`
	EmailDelivered  int64 `json:"email_delivered"`
	EmailFailed     int64 `json:"email_failed"`
}

// GetNotificationDeliveryStats retrieves notification delivery statistics
func (r *NotificationRepository) GetNotificationDeliveryStats(ctx context.Context) (*NotificationDeliveryStats, error) {
	stats := &NotificationDeliveryStats{}

	// Total stats
	r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Select("COUNT(*) as total_sent").
		Where("status = ?", "sent").
		Scan(&stats.TotalSent)

	r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Select("COUNT(*) as total_delivered").
		Where("status = ?", "delivered").
		Scan(&stats.TotalDelivered)

	r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Select("COUNT(*) as total_failed").
		Where("status = ?", "failed").
		Scan(&stats.TotalFailed)

	// Email stats
	r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Select("COUNT(*) as email_sent").
		Where("channel = ? AND status = ?", model.NotificationChannelEmail, "sent").
		Scan(&stats.EmailSent)

	r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Select("COUNT(*) as email_delivered").
		Where("channel = ? AND status = ?", model.NotificationChannelEmail, "delivered").
		Scan(&stats.EmailDelivered)

	r.db.WithContext(ctx).
		Model(&model.NotificationLog{}).
		Select("COUNT(*) as email_failed").
		Where("channel = ? AND status = ?", model.NotificationChannelEmail, "failed").
		Scan(&stats.EmailFailed)

	return stats, nil
}

// CleanupOldNotifications deletes notifications older than specified days
func (r *NotificationRepository) CleanupOldNotifications(ctx context.Context, days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	result := r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&model.Notification{})
	return result.RowsAffected, result.Error
}
