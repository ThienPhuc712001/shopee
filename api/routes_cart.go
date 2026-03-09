package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupCartRoutes configures all cart routes
func SetupCartRoutes(
	rg *gin.RouterGroup,
	cartHandler *handler.CartHandlerEnhanced,
	tokenService service.TokenService,
) {
	// All cart routes require authentication
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Get cart
		protected.GET("", cartHandler.GetCart)
		protected.GET("/summary", cartHandler.GetCartSummary)
		protected.GET("/stats", cartHandler.GetCartStats)
		protected.GET("/checkout", cartHandler.PrepareCheckout)

		// Add to cart
		protected.POST("/add", cartHandler.AddToCart)

		// Update cart item
		protected.PUT("/items/:id", cartHandler.UpdateCartItem)

		// Remove from cart
		protected.DELETE("/items/:id", cartHandler.RemoveFromCart)

		// Clear cart
		protected.DELETE("/clear", cartHandler.ClearCart)
	}
}
