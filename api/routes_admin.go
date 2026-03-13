package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes configures all admin routes
func SetupAdminRoutes(
	rg *gin.RouterGroup,
	adminHandler *handler.AdminHandlerEnhanced,
	tokenService service.TokenService,
) {
	// Public routes (no auth)
	public := rg.Group("")
	{
		public.POST("/auth/login", adminHandler.AdminLogin)
	}

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Admin-only middleware
		adminOnly := protected.Group("")
		adminOnly.Use(middleware.RequireAdmin())
		{
			// User Management
			adminOnly.GET("/users", adminHandler.GetUsers)
			adminOnly.POST("/users/ban", adminHandler.BanUser)

			// Seller Management
			adminOnly.GET("/sellers/pending", adminHandler.GetPendingSellers)
			adminOnly.POST("/sellers/approve", adminHandler.ApproveSeller)

			// Product Management
			adminOnly.GET("/products", adminHandler.GetProductsForModeration)
			adminOnly.DELETE("/products/:id", adminHandler.DeleteProduct)

			// Order Management
			adminOnly.GET("/orders", adminHandler.GetOrders)
			adminOnly.POST("/orders/refund", adminHandler.RefundOrder)

			// Review Management
			adminOnly.GET("/reviews", adminHandler.GetAllReviews)

			// Analytics
			adminOnly.GET("/analytics/stats", adminHandler.GetAdminStats)
			adminOnly.GET("/analytics/sales", adminHandler.GetSalesAnalytics)
			adminOnly.GET("/analytics/users", adminHandler.GetUserAnalytics)
			adminOnly.GET("/analytics/products", adminHandler.GetProductAnalytics)

			// Audit Logs
			adminOnly.GET("/audit-logs", adminHandler.GetAuditLogs)

			// System Settings
			adminOnly.GET("/settings/:key", adminHandler.GetSystemSetting)
			adminOnly.PUT("/settings/:key", adminHandler.UpdateSystemSetting)
		}
	}
}
