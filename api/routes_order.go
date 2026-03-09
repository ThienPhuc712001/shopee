package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes configures all order routes
func SetupOrderRoutes(
	rg *gin.RouterGroup,
	orderHandler *handler.OrderHandlerEnhanced,
	tokenService service.TokenService,
) {
	// All order routes require authentication
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Checkout
		protected.POST("/checkout", orderHandler.Checkout)

		// Get orders
		protected.GET("", orderHandler.GetUserOrders)
		protected.GET("/:id", orderHandler.GetOrder)
		protected.GET("/:id/tracking", orderHandler.GetOrderTracking)
		protected.GET("/statistics", orderHandler.GetOrderStatistics)

		// Order actions
		protected.POST("/:id/cancel", orderHandler.CancelOrder)

		// Seller/Admin routes
		sellerAdmin := protected.Group("")
		sellerAdmin.Use(middleware.RequireSellerOrAdmin())
		{
			sellerAdmin.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		}
	}
}
