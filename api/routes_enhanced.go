package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configures all product routes
func SetupProductRoutes(
	rg *gin.RouterGroup,
	productHandler *handler.ProductHandlerEnhanced,
	tokenService service.TokenService,
) {
	// Public routes (no auth required)
	public := rg.Group("")
	{
		public.GET("", productHandler.GetProducts)
		public.GET("/search", productHandler.SearchProducts)
		public.GET("/featured", productHandler.GetFeaturedProducts)
		public.GET("/best-sellers", productHandler.GetBestSellers)
		public.GET("/:id", productHandler.GetProduct)
		public.GET("/category/:id", productHandler.GetProductsByCategory)
		public.GET("/:id/variants", productHandler.GetVariants)
	}

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Seller routes
		seller := protected.Group("")
		seller.Use(middleware.RequireSeller())
		{
			seller.POST("", productHandler.CreateProduct)
			seller.PUT("/:id", productHandler.UpdateProduct)
			seller.DELETE("/:id", productHandler.DeleteProduct)
			seller.POST("/:id/images", productHandler.UploadProductImages)
			seller.POST("/:id/variants", productHandler.CreateVariant)
			seller.PUT("/:id/stock", productHandler.UpdateStock)
			seller.PATCH("/:id/publish", productHandler.PublishProduct)
			seller.PATCH("/:id/unpublish", productHandler.UnpublishProduct)
		}
	}
}

// SetupCategoryRoutes configures category routes
func SetupCategoryRoutes(rg *gin.RouterGroup, productHandler *handler.ProductHandlerEnhanced) {
	rg.GET("", productHandler.GetCategories)
	rg.GET("/:id", productHandler.GetCategoryByID)
}

// SetupAuthRoutes configures all authentication routes
func SetupAuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, tokenService service.TokenService) {
	// Public routes
	public := rg.Group("")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
		public.POST("/forgot-password", authHandler.RequestPasswordReset)
		public.POST("/reset-password", authHandler.ResetPassword)
		public.GET("/verify-email", authHandler.VerifyEmail)
		public.POST("/resend-verification", authHandler.ResendVerificationEmail)
	}

	// Protected routes (require authentication)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		protected.GET("/me", authHandler.Me)
		protected.PUT("/profile", authHandler.UpdateProfile)
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.POST("/logout", authHandler.Logout)
	}
}

// SetupRouter configures all routes for the application with enhanced auth
func SetupEnhancedRouter(
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandlerEnhanced,
	cartHandler *handler.CartHandlerEnhanced,
	orderHandler *handler.OrderHandlerEnhanced,
	tokenService service.TokenService,
	allowedOrigins []string,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS(allowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "ecommerce-api",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes
		SetupAuthRoutes(api.Group("/auth"), authHandler, tokenService)

		// Public product routes (no auth required)
		setupPublicProductRoutes(api.Group("/products"), productHandler)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(tokenService))
		{
			// Product routes (auth required)
			setupProtectedProductRoutes(protected.Group("/products"), productHandler, tokenService)

			// Cart routes
			setupCartRoutes(protected.Group("/cart"), cartHandler)

			// Order routes
			setupOrderRoutes(protected.Group("/orders"), orderHandler)

			// Admin routes
			setupAdminRoutes(protected.Group("/admin"), tokenService)
		}
	}

	return router
}

// setupPublicProductRoutes configures public product routes
func setupPublicProductRoutes(rg *gin.RouterGroup, handler *handler.ProductHandlerEnhanced) {
	rg.GET("", handler.GetProducts)
	rg.GET("/search", handler.SearchProducts)
	rg.GET("/best-sellers", handler.GetBestSellers)
	rg.GET("/:id", handler.GetProduct)
}

// setupProtectedProductRoutes configures protected product routes with RBAC
func setupProtectedProductRoutes(rg *gin.RouterGroup, handler *handler.ProductHandlerEnhanced, tokenService service.TokenService) {
	// Sellers can create products
	sellerOnly := rg.Group("")
	sellerOnly.Use(middleware.RequireSeller())
	{
		sellerOnly.POST("", handler.CreateProduct)
	}

	// Product owners can update/delete their products
	rg.PUT("/:id", handler.UpdateProduct)
	rg.DELETE("/:id", handler.DeleteProduct)
}

// setupCartRoutes configures cart routes
func setupCartRoutes(rg *gin.RouterGroup, handler *handler.CartHandlerEnhanced) {
	rg.GET("", handler.GetCart)
	rg.POST("/add", handler.AddToCart)
	rg.PUT("/items/:id", handler.UpdateCartItem)
	rg.DELETE("/items/:id", handler.RemoveFromCart)
	rg.DELETE("/clear", handler.ClearCart)
}

// setupOrderRoutes configures order routes
func setupOrderRoutes(rg *gin.RouterGroup, handler *handler.OrderHandlerEnhanced) {
	rg.POST("", handler.Checkout)
	rg.GET("", handler.GetUserOrders)
	rg.GET("/:id", handler.GetOrder)
	rg.POST("/:id/cancel", handler.CancelOrder)

	// Only sellers/admins can update order status
	sellerAdmin := rg.Group("")
	sellerAdmin.Use(middleware.RequireSellerOrAdmin())
	{
		sellerAdmin.PUT("/:id/status", handler.UpdateOrderStatus)
	}
}

// setupAdminRoutes configures admin-only routes
func setupAdminRoutes(rg *gin.RouterGroup, tokenService service.TokenService) {
	adminOnly := rg.Group("")
	adminOnly.Use(middleware.RequireAdmin())
	{
		// Admin-only endpoints go here
		rg.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Admin users list endpoint",
			})
		})

		rg.GET("/shops", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Admin shops list endpoint",
			})
		})

		rg.GET("/orders/all", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Admin all orders endpoint",
			})
		})
	}
}

// Example route with owner check
func SetupExampleOwnerRoute(rg *gin.RouterGroup, tokenService service.TokenService) {
	// Only owner or admin can access
	rg.GET("/orders/:id", func(c *gin.Context) {
		// Get owner ID from database based on order ID
		// For example purposes, we'll use a mock function
		getOwnerID := func(c *gin.Context) uint {
			// In production, fetch order and get user_id
			// This is just an example
			return 1
		}

		// Check if user is authenticated
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Check if admin
		userRole, _ := c.Get("user_role")
		if userRole == "admin" {
			c.Next()
			return
		}

		// Check if owner
		ownerID := getOwnerID(c)
		if userID.(uint) != ownerID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}, func(c *gin.Context) {
		// Handler logic here
		c.JSON(http.StatusOK, gin.H{"message": "Order details"})
	})
}
