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
		// User Management
		protected.GET("/users", adminHandler.GetUsers)
		protected.POST("/users/ban", adminHandler.BanUser)

		// Seller Management
		protected.GET("/sellers/pending", adminHandler.GetUsers) // Simplified
		protected.POST("/sellers/approve", adminHandler.ApproveSeller)

		// Product Management
		protected.GET("/products", adminHandler.GetUsers) // Simplified
		protected.DELETE("/products/:id", adminHandler.DeleteProduct)

		// Order Management
		protected.GET("/orders", adminHandler.GetOrders)
		protected.POST("/orders/refund", adminHandler.RefundOrder)

		// Analytics
		protected.GET("/analytics/stats", adminHandler.GetAdminStats)
		protected.GET("/analytics/sales", adminHandler.GetSalesAnalytics)
		protected.GET("/analytics/users", adminHandler.GetUserAnalytics)
		protected.GET("/analytics/products", adminHandler.GetProductAnalytics)

		// Audit Logs
		protected.GET("/audit-logs", adminHandler.GetAuditLogs)

		// System Settings
		protected.GET("/settings/:key", adminHandler.GetSystemSetting)
		protected.PUT("/settings/:key", adminHandler.UpdateSystemSetting)
	}
}
