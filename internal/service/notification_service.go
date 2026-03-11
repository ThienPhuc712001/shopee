package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"ecommerce/pkg/email"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// NotificationService handles notification business logic
type NotificationService struct {
	repo           *repository.NotificationRepository
	emailService   *email.EmailService
	userRepository repository.UserRepositoryEnhanced
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	repo *repository.NotificationRepository,
	emailService *email.EmailService,
	userRepo repository.UserRepositoryEnhanced,
) *NotificationService {
	return &NotificationService{
		repo:           repo,
		emailService:   emailService,
		userRepository: userRepo,
	}
}

// ==================== NOTIFICATION CREATION ====================

// CreateNotification creates a single notification
func (s *NotificationService) CreateNotification(ctx context.Context, input model.NotificationInput) (*model.Notification, error) {
	notification := &model.Notification{
		UserID:    input.UserID,
		Title:     input.Title,
		Message:   input.Message,
		Type:      input.Type,
		Priority:  input.Priority,
		Data:      input.Data,
		ActionURL: input.ActionURL,
		ImageURL:  input.ImageURL,
	}

	if notification.Type == "" {
		notification.Type = model.NotificationTypeSystem
	}
	if notification.Priority == "" {
		notification.Priority = model.NotificationPriorityNormal
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// CreateNotificationBatch creates multiple notifications
func (s *NotificationService) CreateNotificationBatch(ctx context.Context, input model.NotificationBatchInput) (int, error) {
	notifications := make([]*model.Notification, 0, len(input.UserIDs))

	for _, userID := range input.UserIDs {
		notification := &model.Notification{
			UserID:    userID,
			Title:     input.Title,
			Message:   input.Message,
			Type:      input.Type,
			Priority:  input.Priority,
			Data:      input.Data,
			ActionURL: input.ActionURL,
		}
		notifications = append(notifications, notification)
	}

	if err := s.repo.CreateNotifications(ctx, notifications); err != nil {
		return 0, err
	}

	return len(notifications), nil
}

// ==================== ORDER NOTIFICATIONS ====================

// SendOrderNotification sends order-related notification
func (s *NotificationService) SendOrderNotification(ctx context.Context, userID uint, orderNumber string, orderStatus string, total float64) error {
	title, message := s.getOrderNotificationContent(orderStatus)

	data, _ := json.Marshal(map[string]interface{}{
		"order_number": orderNumber,
		"order_status": orderStatus,
		"total":        total,
	})

	notification := &model.Notification{
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      model.NotificationTypeOrder,
		Priority:  model.NotificationPriorityHigh,
		Data:      string(data),
		ActionURL: fmt.Sprintf("/orders/%s", orderNumber),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return err
	}

	// Send email notification
	go s.sendOrderEmail(ctx, userID, orderNumber, orderStatus, total)

	return nil
}

func (s *NotificationService) getOrderNotificationContent(status string) (string, string) {
	switch status {
	case "created":
		return "Order Created", fmt.Sprintf("Your order %s has been created successfully.", status)
	case "confirmed":
		return "Order Confirmed", "Your order has been confirmed and is being processed."
	case "processing":
		return "Order Processing", "Your order is being prepared for shipment."
	case "shipped":
		return "Order Shipped", "Your order has been shipped and is on its way."
	case "delivered":
		return "Order Delivered", "Your order has been delivered. Enjoy your purchase!"
	case "cancelled":
		return "Order Cancelled", "Your order has been cancelled."
	default:
		return "Order Update", fmt.Sprintf("Your order status has been updated to %s.", status)
	}
}

func (s *NotificationService) sendOrderEmail(ctx context.Context, userID uint, orderNumber, status string, total float64) {
	// Get user email
	user, err := s.userRepository.FindByID(userID)
	if err != nil || user.Email == "" {
		return
	}

	// Get user preference
	preference, err := s.repo.GetNotificationPreference(ctx, userID)
	if err != nil || !preference.EmailEnabled || !preference.OrderEnabled {
		return
	}

	switch status {
	case "created":
		_ = s.emailService.SendOrderConfirmation(user.Email, user.FirstName+" "+user.LastName, orderNumber, total)
	case "shipped":
		_ = s.emailService.SendShippingUpdate(user.Email, user.FirstName+" "+user.LastName, orderNumber, "", "shipped")
	}
}

// ==================== PAYMENT NOTIFICATIONS ====================

// SendPaymentNotification sends payment-related notification
func (s *NotificationService) SendPaymentNotification(ctx context.Context, userID uint, orderNumber string, amount float64, paymentStatus string, paymentMethod string) error {
	title, message := s.getPaymentNotificationContent(paymentStatus)

	data, _ := json.Marshal(map[string]interface{}{
		"order_number":   orderNumber,
		"amount":         amount,
		"payment_status": paymentStatus,
		"payment_method": paymentMethod,
	})

	notification := &model.Notification{
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      model.NotificationTypePayment,
		Priority:  model.NotificationPriorityHigh,
		Data:      string(data),
		ActionURL: fmt.Sprintf("/orders/%s", orderNumber),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return err
	}

	// Send email notification
	go s.sendPaymentEmail(ctx, userID, orderNumber, amount, paymentStatus, paymentMethod)

	return nil
}

func (s *NotificationService) getPaymentNotificationContent(status string) (string, string) {
	switch status {
	case "success":
		return "Payment Successful", "Your payment has been processed successfully."
	case "pending":
		return "Payment Pending", "Your payment is being processed."
	case "failed":
		return "Payment Failed", "Your payment could not be processed. Please try again."
	case "refunded":
		return "Payment Refunded", "Your payment has been refunded."
	default:
		return "Payment Update", fmt.Sprintf("Your payment status has been updated to %s.", status)
	}
}

func (s *NotificationService) sendPaymentEmail(ctx context.Context, userID uint, orderNumber string, amount float64, status, paymentMethod string) {
	// Get user email
	user, err := s.userRepository.FindByID(userID)
	if err != nil || user.Email == "" {
		return
	}

	// Get user preference
	preference, err := s.repo.GetNotificationPreference(ctx, userID)
	if err != nil || !preference.EmailEnabled || !preference.PaymentEnabled {
		return
	}

	if status == "success" {
		_ = s.emailService.SendPaymentConfirmation(user.Email, user.FirstName+" "+user.LastName, orderNumber, amount, paymentMethod)
	}
}

// ==================== SHIPPING NOTIFICATIONS ====================

// SendShippingNotification sends shipping-related notification
func (s *NotificationService) SendShippingNotification(ctx context.Context, userID uint, orderNumber string, trackingNumber string, status string) error {
	title, message := s.getShippingNotificationContent(status)

	data, _ := json.Marshal(map[string]interface{}{
		"order_number":    orderNumber,
		"tracking_number": trackingNumber,
		"shipping_status": status,
	})

	notification := &model.Notification{
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      model.NotificationTypeShipping,
		Priority:  model.NotificationPriorityNormal,
		Data:      string(data),
		ActionURL: fmt.Sprintf("/orders/%s/tracking", orderNumber),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return err
	}

	// Send email notification
	go s.sendShippingEmail(ctx, userID, orderNumber, trackingNumber, status)

	return nil
}

func (s *NotificationService) getShippingNotificationContent(status string) (string, string) {
	switch status {
	case "shipped":
		return "Order Shipped", "Your order has been shipped and is on its way."
	case "in_transit":
		return "Order In Transit", "Your order is currently in transit."
	case "out_for_delivery":
		return "Out for Delivery", "Your order is out for delivery today."
	case "delivered":
		return "Order Delivered", "Your order has been delivered successfully."
	case "failed":
		return "Delivery Failed", "Delivery attempt failed. We'll retry soon."
	default:
		return "Shipping Update", fmt.Sprintf("Your order shipping status: %s.", status)
	}
}

func (s *NotificationService) sendShippingEmail(ctx context.Context, userID uint, orderNumber, trackingNumber, status string) {
	// Get user email
	user, err := s.userRepository.FindByID(userID)
	if err != nil || user.Email == "" {
		return
	}

	// Get user preference
	preference, err := s.repo.GetNotificationPreference(ctx, userID)
	if err != nil || !preference.EmailEnabled || !preference.ShippingEnabled {
		return
	}

	_ = s.emailService.SendShippingUpdate(user.Email, user.FirstName+" "+user.LastName, orderNumber, trackingNumber, status)
}

// ==================== PROMOTION NOTIFICATIONS ====================

// SendPromotionNotification sends promotion-related notification
func (s *NotificationService) SendPromotionNotification(ctx context.Context, userID uint, promoTitle string, promoCode string, discountPercent float64) error {
	title := fmt.Sprintf("🎉 %s", promoTitle)
	message := fmt.Sprintf("Use code %s to get %.0f%% off your next order!", promoCode, discountPercent)

	data, _ := json.Marshal(map[string]interface{}{
		"promo_title":      promoTitle,
		"promo_code":       promoCode,
		"discount_percent": discountPercent,
	})

	notification := &model.Notification{
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      model.NotificationTypePromotion,
		Priority:  model.NotificationPriorityNormal,
		Data:      string(data),
		ActionURL: "/promotions",
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return err
	}

	// Send email notification
	go s.sendPromotionEmail(ctx, userID, promoTitle, promoCode, discountPercent)

	return nil
}

func (s *NotificationService) sendPromotionEmail(ctx context.Context, userID uint, promoTitle, promoCode string, discountPercent float64) {
	// Get user email
	user, err := s.userRepository.FindByID(userID)
	if err != nil || user.Email == "" {
		return
	}

	// Get user preference
	preference, err := s.repo.GetNotificationPreference(ctx, userID)
	if err != nil || !preference.EmailEnabled || !preference.PromotionEnabled {
		return
	}

	_ = s.emailService.SendPromotionEmail(user.Email, user.FirstName+" "+user.LastName, promoTitle, promoCode, discountPercent)
}

// SendPromotionNotificationBatch sends promotion to multiple users
func (s *NotificationService) SendPromotionNotificationBatch(ctx context.Context, userIDs []uint, promoTitle string, promoCode string, discountPercent float64) (int, error) {
	count := 0
	for _, userID := range userIDs {
		if err := s.SendPromotionNotification(ctx, userID, promoTitle, promoCode, discountPercent); err == nil {
			count++
		}
	}
	return count, nil
}

// ==================== NOTIFICATION QUERIES ====================

// GetUserNotifications retrieves notifications for a user
func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uint, filter model.NotificationFilter) ([]model.Notification, int64, error) {
	return s.repo.GetUserNotifications(ctx, userID, filter)
}

// GetUnreadCount retrieves unread notification count for a user
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

// GetNotificationSummary retrieves notification summary for a user
func (s *NotificationService) GetNotificationSummary(ctx context.Context, userID uint) (*model.NotificationSummary, error) {
	totalCount, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	recent, err := s.repo.GetRecentNotifications(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	views := make([]model.NotificationView, 0, len(recent))
	for _, n := range recent {
		views = append(views, n.ToView())
	}

	return &model.NotificationSummary{
		TotalCount:        totalCount,
		UnreadCount:       unreadCount,
		RecentNotifications: views,
	}, nil
}

// GetNotificationStats retrieves notification statistics for a user
func (s *NotificationService) GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	return s.repo.GetNotificationStats(ctx, userID)
}

// ==================== NOTIFICATION MANAGEMENT ====================

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uint, userID uint) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID uint, userID uint) error {
	return s.repo.DeleteNotification(ctx, notificationID, userID)
}

// GetNotificationByID retrieves a notification by ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, id uint) (*model.Notification, error) {
	return s.repo.GetNotificationByID(ctx, id)
}

// ==================== NOTIFICATION PREFERENCES ====================

// GetNotificationPreference retrieves notification preference for a user
func (s *NotificationService) GetNotificationPreference(ctx context.Context, userID uint) (*model.NotificationPreference, error) {
	return s.repo.GetOrCreateNotificationPreference(ctx, userID)
}

// UpdateNotificationPreference updates notification preference for a user
func (s *NotificationService) UpdateNotificationPreference(ctx context.Context, userID uint, input model.NotificationPreference) (*model.NotificationPreference, error) {
	preference, err := s.repo.GetOrCreateNotificationPreference(ctx, userID)
	if err != nil {
		return nil, err
	}

	preference.EmailEnabled = input.EmailEnabled
	preference.SMSEnabled = input.SMSEnabled
	preference.PushEnabled = input.PushEnabled
	preference.OrderEnabled = input.OrderEnabled
	preference.PaymentEnabled = input.PaymentEnabled
	preference.ShippingEnabled = input.ShippingEnabled
	preference.PromotionEnabled = input.PromotionEnabled
	preference.SystemEnabled = input.SystemEnabled

	if err := s.repo.UpdateNotificationPreference(ctx, preference); err != nil {
		return nil, err
	}

	return preference, nil
}

// ==================== CLEANUP ====================

// CleanupOldNotifications deletes old notifications
func (s *NotificationService) CleanupOldNotifications(ctx context.Context, days int) (int64, error) {
	return s.repo.CleanupOldNotifications(ctx, days)
}

// Error definitions
var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrPreferenceNotFound   = errors.New("notification preference not found")
)
