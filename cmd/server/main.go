package main

import (
	"ecommerce/api"
	"ecommerce/internal/domain/model"
	"ecommerce/internal/handler"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/config"
	"ecommerce/pkg/database"
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

func main() {
	// Load configuration
	cfg := config.Load()

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
	_ = repository.NewShopRepositoryEnhanced(db)
	paymentRepo := repository.NewPaymentRepositoryEnhanced(db)

	// Initialize token service with separate secrets for access and refresh tokens
	tokenService := service.NewTokenService(service.TokenServiceConfig{
		AccessSecret:  cfg.JWT.Secret,
		RefreshSecret: cfg.JWT.Secret + "-refresh-token-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, 15*time.Minute)
	productService := service.NewProductServiceEnhanced(productRepo)
	cartService := service.NewCartServiceEnhanced(cartRepo, productRepo)
	orderService := service.NewOrderServiceEnhanced(orderRepo, cartRepo, productRepo)
	paymentService := service.NewPaymentServiceEnhanced(paymentRepo, orderRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandlerEnhanced(productService)
	cartHandler := handler.NewCartHandlerEnhanced(cartService)
	orderHandler := handler.NewOrderHandlerEnhanced(orderService)
	_ = handler.NewPaymentHandlerEnhanced(paymentService)

	// Setup router with enhanced authentication
	router := routes.SetupEnhancedRouter(
		authHandler,
		productHandler,
		cartHandler,
		orderHandler,
		tokenService,
		cfg.CORS.AllowedOrigins,
	)

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
