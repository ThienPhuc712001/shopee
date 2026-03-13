package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupShopRoutes configures all shop routes
func SetupShopRoutes(
	rg *gin.RouterGroup,
	shopHandler *handler.ShopHandler,
	tokenService service.TokenService,
) {
	// Public routes - register parameterized routes LAST to avoid conflicts
	public := rg.Group("")
	{
		// List all shops - use /shops-list to avoid conflict with /shops/:id
		public.GET("/shops-list", shopHandler.GetShops)
		
		// Parameterized route MUST be last
		public.GET("/shops/:id", shopHandler.GetShop)
	}

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Shop management routes
		protected.POST("/shops", shopHandler.CreateShop)
		protected.GET("/shops/my", shopHandler.GetMyShop)
		protected.GET("/shops/seller/me", shopHandler.GetShopsBySeller)
	}
}
