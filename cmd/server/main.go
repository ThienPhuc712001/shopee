package main

import (
	"ecommerce/api"
	"ecommerce/internal/domain/model"
	"ecommerce/internal/handler"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/config"
	"ecommerce/pkg/database"
	"ecommerce/pkg/email"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var appLog *logger.Logger

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	var err error
	appLog, err = logger.New(logger.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Setup Gin mode based on environment
	setupGinMode(cfg.App.Env)

	// Initialize database
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto migrations
	if err := runMigrations(); err != nil {
		log.Printf("Warning: Migration encountered an issue: %v", err)
		log.Println("If this is a unique constraint error on 'email', you may need to manually drop the existing constraint.")
		log.Println("SQL: DROP INDEX IX_Users_Email ON dbo.users; -- or use ALTER TABLE dbo.users DROP CONSTRAINT [constraint_name];")
		// Continue anyway - the schema may still be functional
	}

	// Initialize repositories
	userRepo := repository.NewUserRepositoryEnhanced(db)
	productRepo := repository.NewProductRepositoryEnhanced(db)
	cartRepo := repository.NewCartRepositoryEnhanced(db)
	orderRepo := repository.NewOrderRepositoryEnhanced(db)
	shopRepo := repository.NewShopRepositoryEnhanced(db)
	paymentRepo := repository.NewPaymentRepositoryEnhanced(db)
	imageRepo := repository.NewImageRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	adminRepo := repository.NewAdminRepositoryEnhanced(db)

	// Initialize token service with separate secrets for access and refresh tokens
	tokenService := service.NewTokenService(service.TokenServiceConfig{
		AccessSecret:  cfg.JWT.Secret,
		RefreshSecret: cfg.JWT.Secret + "-refresh-token-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	// Initialize repositories
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Initialize auth service with new constructor
	authService := service.NewAuthService(service.AuthServiceConfig{
		DB:               db,
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
		TokenService:     tokenService,
		Log:              appLog,
		MaxLoginAttempts: 5,
		LockoutDuration:  15 * time.Minute,
	})
	productService := service.NewProductServiceEnhanced(productRepo)
	cartService := service.NewCartServiceEnhanced(cartRepo, productRepo)
	
	// Initialize email service
	emailConfig := model.EmailConfig{
		Host:      cfg.Email.Host,
		Port:      cfg.Email.Port,
		Username:  cfg.Email.Username,
		Password:  cfg.Email.Password,
		FromName:  cfg.Email.FromName,
		FromEmail: cfg.Email.FromEmail,
		UseTLS:    cfg.Email.UseTLS,
	}
	emailService := email.NewEmailService(emailConfig)
	
	couponService := service.NewCouponService(couponRepo)
	notificationRepo := repository.NewNotificationRepository(db)
	notificationService := service.NewNotificationService(notificationRepo, emailService, userRepo)
	orderService := service.NewOrderServiceEnhanced(orderRepo, cartRepo, productRepo, couponRepo, couponService, notificationService)
	paymentService := service.NewPaymentServiceEnhanced(paymentRepo, orderRepo)
	uploadService := service.NewUploadService(imageRepo, "./uploads")
	categoryService := service.NewCategoryService(categoryRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	shippingRepo := repository.NewShippingRepository(db)
	shippingService := service.NewShippingService(shippingRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, appLog)
	
	// Initialize shop service and handler first (needed by product handler)
	shopService := service.NewShopService(shopRepo)
	shopHandler := handler.NewShopHandler(shopService)
	
	// Initialize product handler with shop service
	productHandler := handler.NewProductHandlerEnhanced(productService)
	productHandler.SetShopService(shopService) // Set shop service for product handler
	
	cartHandler := handler.NewCartHandlerEnhanced(cartService)
	orderHandler := handler.NewOrderHandlerEnhanced(orderService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	couponHandler := handler.NewCouponHandler(couponService)
	shippingHandler := handler.NewShippingHandler(shippingService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	_ = handler.NewPaymentHandlerEnhanced(paymentService)
	
	// Initialize admin service and handler
	adminService := service.NewAdminServiceEnhanced(adminRepo, userRepo, shopRepo, productRepo, orderRepo)
	adminHandler := handler.NewAdminHandlerEnhanced(adminService)

	// Setup router with enhanced authentication
	router := routes.SetupEnhancedRouter(
		authHandler,
		productHandler,
		cartHandler,
		orderHandler,
		tokenService,
		cfg.CORS.AllowedOrigins,
	)

	// Setup upload routes
	apiRoutes := router.Group("/api")
	routes.SetupUploadRoutes(apiRoutes, uploadHandler, tokenService)

	// Setup category routes
	routes.SetupCategoryRoutes(apiRoutes, categoryHandler, tokenService)

	// Setup shop routes
	routes.SetupShopRoutes(apiRoutes, shopHandler, tokenService)

	// Setup product routes (search is included)
	routes.SetupProductRoutes(apiRoutes, productHandler, tokenService)

	// Setup inventory routes
	routes.SetupInventoryRoutes(apiRoutes, inventoryHandler, tokenService)

	// Setup coupon routes
	routes.SetupCouponRoutes(apiRoutes, couponHandler, tokenService)

	// Setup shipping routes
	routes.SetupShippingRoutes(apiRoutes, shippingHandler, tokenService)

	// Setup notification routes
	routes.SetupNotificationRoutes(apiRoutes, notificationHandler, tokenService)

	// Setup admin routes
	routes.SetupAdminRoutes(apiRoutes.Group("/admin"), adminHandler, tokenService)

	// Setup static file serving for uploaded images
	routes.SetupStaticRoutes(router, "./uploads")

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Requests*2)
	rateLimiter.StartCleanup(1 * time.Minute)

	// Apply rate limiting to API routes
	router.Use(rateLimiter.RateLimit())

	// Create HTTP server with secure settings
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("╔════════════════════════════════════════════════════════╗")
		log.Printf("║          E-Commerce API Server Starting                ║")
		log.Printf("╠════════════════════════════════════════════════════════╣")
		log.Printf("║  Port:           http://localhost:%s                   ║", cfg.App.Port)
		log.Printf("║  Environment:    %s                                    ║", cfg.App.Env)
		log.Printf("║  Health Check:   /health                               ║")
		log.Printf("╠════════════════════════════════════════════════════════╣")
		log.Printf("║  Product Search Endpoints:                             ║")
		log.Printf("║    GET  /api/products/search                           ║")
		log.Printf("║    GET  /api/products/featured                         ║")
		log.Printf("║    GET  /api/products/best-sellers                     ║")
		log.Printf("╠════════════════════════════════════════════════════════╣")
		log.Printf("║  Category Endpoints:                                   ║")
		log.Printf("║    GET  /api/categories                                ║")
		log.Printf("║    GET  /api/categories/tree                           ║")
		log.Printf("║    GET  /api/categories/:id/products                   ║")
		log.Printf("║    POST /api/categories (Admin)                        ║")
		log.Printf("╠════════════════════════════════════════════════════════╣")
		log.Printf("║  Image Upload Endpoints:                               ║")
		log.Printf("║    POST /api/upload/product                            ║")
		log.Printf("║    POST /api/upload/product/multiple                   ║")
		log.Printf("║    POST /api/upload/review                             ║")
		log.Printf("║    POST /api/upload/avatar                             ║")
		log.Printf("║    GET  /uploads/{filename}                            ║")
		log.Printf("╠════════════════════════════════════════════════════════╣")
		log.Printf("║  Authentication Endpoints:                             ║")
		log.Printf("║    POST /api/auth/register                             ║")
		log.Printf("║    POST /api/auth/login                                ║")
		log.Printf("║    POST /api/auth/refresh                              ║")
		log.Printf("║    GET  /api/auth/me                                   ║")
		log.Printf("║    POST /api/auth/logout                               ║")
		log.Printf("╚════════════════════════════════════════════════════════╝")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := database.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server exited gracefully")
}

// runMigrations runs database migrations for all models
func runMigrations() error {
	models := []interface{}{
		&model.User{},
		&model.Address{},
		&model.Shop{},
		&model.Category{},
		&model.Product{},
		&model.ProductImage{},
		&model.Cart{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
		&model.Review{},
		// Image upload tables
		&model.ReviewImage{},
		&model.UserAvatar{},
		&model.ImageUploadLog{},
		// Inventory tables
		&model.Inventory{},
		&model.InventoryLog{},
		&model.StockAlert{},
		// Coupon tables
		&model.Coupon{},
		&model.CouponUsage{},
		// Shipping tables
		&model.ShippingAddress{},
		&model.Shipment{},
		&model.ShipmentTracking{},
		&model.ShippingCarrier{},
		// Notification tables
		&model.Notification{},
		&model.NotificationPreference{},
		&model.NotificationLog{},
	}

	return database.AutoMigrate(models...)
}

// setupGinMode configures Gin mode based on environment
func setupGinMode(env string) {
	switch env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
