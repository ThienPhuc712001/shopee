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
	// Public routes
	public := rg.Group("/shops")
	{
		public.GET("/:id", shopHandler.GetShop)
	}

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Shop management routes
		protected.POST("/shops", shopHandler.CreateShop)
		protected.GET("/shops/my", shopHandler.GetMyShop)
	}
}
