package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes configures all notification routes
func SetupNotificationRoutes(
	rg *gin.RouterGroup,
	notificationHandler *handler.NotificationHandler,
	tokenService service.TokenService,
) {
	// User routes (require authentication)
	user := rg.Group("")
	user.Use(middleware.JWTAuth(tokenService))
	{
		// Get notifications
		user.GET("/notifications", notificationHandler.GetUserNotifications)
		
		// Get notification summary
		user.GET("/notifications/summary", notificationHandler.GetNotificationSummary)
		
		// Get unread count
		user.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
		
		// Get notification statistics
		user.GET("/notifications/stats", notificationHandler.GetNotificationStats)
		
		// Mark notification as read
		user.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		
		// Mark all notifications as read
		user.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
		
		// Delete notification
		user.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
		
		// Notification preferences
		user.GET("/notifications/preferences", notificationHandler.GetNotificationPreference)
		user.PUT("/notifications/preferences", notificationHandler.UpdateNotificationPreference)
	}

	// Admin routes (require authentication + admin role)
	admin := rg.Group("/admin")
	admin.Use(middleware.JWTAuth(tokenService))
	admin.Use(middleware.RequireAdmin())
	{
		// Create single notification
		admin.POST("/notifications", notificationHandler.CreateNotification)
		
		// Create batch notifications
		admin.POST("/notifications/batch", notificationHandler.CreateNotificationBatch)
		
		// Send promotion notification
		admin.POST("/notifications/promotion", notificationHandler.SendPromotionNotification)
		
		// Get delivery statistics
		admin.GET("/notifications/delivery-stats", notificationHandler.GetDeliveryStats)
		
		// Cleanup old notifications
		admin.POST("/notifications/cleanup", notificationHandler.CleanupOldNotifications)
	}
}
