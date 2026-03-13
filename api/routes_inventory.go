package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupInventoryRoutes configures all inventory routes
func SetupInventoryRoutes(
	rg *gin.RouterGroup,
	inventoryHandler *handler.InventoryHandler,
	tokenService service.TokenService,
) {
	// Protected routes (require authentication)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Get inventory list (for sellers/admins)
		protected.GET("/inventory", inventoryHandler.GetInventoryList)

		// Get inventory summary
		protected.GET("/inventory/summary", inventoryHandler.GetInventorySummary)

		// Get low stock products
		protected.GET("/inventory/low-stock", inventoryHandler.GetLowStockProducts)

		// Get out of stock products
		protected.GET("/inventory/out-of-stock", inventoryHandler.GetOutOfStockProducts)

		// Restock product
		protected.POST("/inventory/restock", inventoryHandler.RestockProduct)
		
		// Create stock alert
		protected.POST("/inventory/alerts", inventoryHandler.CreateStockAlert)
	}

	// Public routes (no authentication required)
	public := rg.Group("")
	{
		// Check stock availability
		public.POST("/inventory/check", inventoryHandler.CheckStock)

		// Get inventory by product ID
		public.GET("/inventory/:product_id", inventoryHandler.GetInventory)

		// Get inventory logs
		public.GET("/inventory/:product_id/logs", inventoryHandler.GetInventoryLogs)
	}
}
