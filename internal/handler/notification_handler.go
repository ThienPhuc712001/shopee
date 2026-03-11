package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification HTTP requests
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// ==================== USER NOTIFICATION ENDPOINTS ====================

// GetUserNotifications handles GET /api/notifications
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	typeParam := c.Query("type")
	isReadParam := c.Query("is_read")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := model.NotificationFilter{
		Page:  page,
		Limit: limit,
	}

	if typeParam != "" {
		t := model.NotificationType(typeParam)
		filter.Type = &t
	}
	if isReadParam != "" {
		isRead := isReadParam == "true"
		filter.IsRead = &isRead
	}

	notifications, total, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID.(uint), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get notifications",
			"error":   err.Error(),
		})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	// Convert to views
	views := make([]model.NotificationView, 0, len(notifications))
	for _, n := range notifications {
		views = append(views, n.ToView())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"notifications": views,
			"pagination": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// GetNotificationSummary handles GET /api/notifications/summary
func (h *NotificationHandler) GetNotificationSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	summary, err := h.notificationService.GetNotificationSummary(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get notification summary",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetUnreadCount handles GET /api/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get unread count",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"unread_count": count,
		},
	})
}

// MarkAsRead handles PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.MarkAsRead(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to mark notification as read",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead handles PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to mark all notifications as read",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All notifications marked as read",
	})
}

// DeleteNotification handles DELETE /api/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.DeleteNotification(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete notification",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification deleted successfully",
	})
}

// GetNotificationStats handles GET /api/notifications/stats
func (h *NotificationHandler) GetNotificationStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	stats, err := h.notificationService.GetNotificationStats(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get notification statistics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ==================== NOTIFICATION PREFERENCES ====================

// GetNotificationPreference handles GET /api/notifications/preferences
func (h *NotificationHandler) GetNotificationPreference(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	preference, err := h.notificationService.GetNotificationPreference(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get notification preferences",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    preference,
	})
}

// UpdateNotificationPreference handles PUT /api/notifications/preferences
func (h *NotificationHandler) UpdateNotificationPreference(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	var input model.NotificationPreference
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	preference, err := h.notificationService.UpdateNotificationPreference(c.Request.Context(), userID.(uint), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update notification preferences",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification preferences updated successfully",
		"data":    preference,
	})
}

// ==================== ADMIN NOTIFICATION ENDPOINTS ====================

// CreateNotification handles POST /api/admin/notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var input model.NotificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	notification, err := h.notificationService.CreateNotification(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create notification",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Notification created successfully",
		"data":    notification,
	})
}

// CreateNotificationBatch handles POST /api/admin/notifications/batch
func (h *NotificationHandler) CreateNotificationBatch(c *gin.Context) {
	var input model.NotificationBatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	count, err := h.notificationService.CreateNotificationBatch(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create notifications",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Notifications created successfully",
		"data": gin.H{
			"count": count,
		},
	})
}

// SendPromotionNotification handles POST /api/admin/notifications/promotion
func (h *NotificationHandler) SendPromotionNotification(c *gin.Context) {
	var input struct {
		UserIDs         []uint  `json:"user_ids" binding:"required"`
		PromoTitle      string  `json:"promo_title" binding:"required"`
		PromoCode       string  `json:"promo_code" binding:"required"`
		DiscountPercent float64 `json:"discount_percent" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	count, err := h.notificationService.SendPromotionNotificationBatch(
		c.Request.Context(),
		input.UserIDs,
		input.PromoTitle,
		input.PromoCode,
		input.DiscountPercent,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send promotion notifications",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Promotion notifications sent successfully",
		"data": gin.H{
			"sent_count": count,
		},
	})
}

// GetDeliveryStats handles GET /api/admin/notifications/delivery-stats
func (h *NotificationHandler) GetDeliveryStats(c *gin.Context) {
	// This would require adding a method to the service
	// For now, return a placeholder
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Delivery statistics endpoint",
		},
	})
}

// CleanupOldNotifications handles POST /api/admin/notifications/cleanup
func (h *NotificationHandler) CleanupOldNotifications(c *gin.Context) {
	var input struct {
		Days int `json:"days" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	if input.Days < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Days must be at least 1",
		})
		return
	}

	deleted, err := h.notificationService.CleanupOldNotifications(c.Request.Context(), input.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to cleanup old notifications",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Old notifications cleaned up successfully",
		"data": gin.H{
			"deleted_count": deleted,
		},
	})
}
