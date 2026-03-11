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

	// Admin routes (require authentication + admin role)
	admin := rg.Group("")
	admin.Use(middleware.JWTAuth(tokenService))
	admin.Use(middleware.RequireAdmin())
	{
		// Get inventory summary
		admin.GET("/inventory/summary", inventoryHandler.GetInventorySummary)
		
		// Get low stock products
		admin.GET("/inventory/low-stock", inventoryHandler.GetLowStockProducts)
		
		// Get out of stock products
		admin.GET("/inventory/out-of-stock", inventoryHandler.GetOutOfStockProducts)
		
		// Restock product
		admin.POST("/inventory/restock", inventoryHandler.RestockProduct)
	}
}
