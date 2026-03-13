package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupShippingRoutes configures all shipping routes
func SetupShippingRoutes(
	rg *gin.RouterGroup,
	shippingHandler *handler.ShippingHandler,
	tokenService service.TokenService,
) {
	// Public routes (no authentication required)
	public := rg.Group("")
	{
		// Get active carriers
		public.GET("/shipping/carriers", shippingHandler.GetActiveCarriers)

		// Get shipping methods
		public.GET("/shipping/methods", shippingHandler.GetShippingMethods)
		
		// Calculate shipping cost
		public.POST("/shipping/calculate", shippingHandler.CalculateShipping)
	}

	// User routes (require authentication)
	user := rg.Group("")
	user.Use(middleware.JWTAuth(tokenService))
	{
		// Shipping addresses
		user.POST("/shipping/addresses", shippingHandler.CreateAddress)
		user.GET("/shipping/addresses", shippingHandler.GetAddresses)
		user.GET("/shipping/addresses/:id", shippingHandler.GetAddressByID)
		user.PUT("/shipping/addresses/:id", shippingHandler.UpdateAddress)
		user.DELETE("/shipping/addresses/:id", shippingHandler.DeleteAddress)
		user.POST("/shipping/addresses/:id/set-default", shippingHandler.SetDefaultAddress)

		// Order tracking (user can track their own orders)
		user.GET("/orders/:id/tracking", shippingHandler.GetOrderTracking)
		user.GET("/orders/:id/shipment", shippingHandler.GetShipmentByOrderID)
	}

	// Admin routes (require authentication + admin role)
	admin := rg.Group("")
	admin.Use(middleware.JWTAuth(tokenService))
	admin.Use(middleware.RequireAdmin())
	{
		// Shipments management
		admin.POST("/shipments", shippingHandler.CreateShipment)
		admin.GET("/shipments", shippingHandler.GetShipmentsByStatus)
		admin.GET("/shipments/:id", shippingHandler.GetShipmentByID)
		admin.GET("/shipments/:id/tracking", shippingHandler.GetTrackingTimeline)
		admin.PUT("/shipments/:id/status", shippingHandler.UpdateShipmentStatus)
		admin.POST("/shipments/:id/tracking", shippingHandler.AddTrackingEvent)

		// Statistics
		admin.GET("/shipments/stats", shippingHandler.GetShipmentStats)
	}
}
