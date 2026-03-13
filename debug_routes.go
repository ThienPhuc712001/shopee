package main

import (
	"ecommerce/api"
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/config"
	"ecommerce/pkg/database"
	"ecommerce/pkg/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	// Create mock handlers
	mockHandler := &handler.ProductHandlerEnhanced{}
	mockCartHandler := &handler.CartHandlerEnhanced{}
	mockOrderHandler := &handler.OrderHandlerEnhanced{}
	
	tokenService := &service.TokenService{}
	
	authHandler := &handler.AuthHandler{}

	// Setup router
	router := routes.SetupEnhancedRouter(
		authHandler,
		mockHandler,
		mockCartHandler,
		mockOrderHandler,
		tokenService,
		cfg.CORS.AllowedOrigins,
	)

	// Add a debug route to list all routes
	router.GET("/debug/routes", func(c *gin.Context) {
		routes := []string{}
		router.Routes()
		c.JSON(http.StatusOK, gin.H{
			"routes": routes,
			"message": "Check server logs for registered routes",
		})
	})

	// Start server
	go func() {
		fmt.Println("Server starting on :8080")
		router.Run(":8080")
	}()

	// Wait for routes to be registered
	time.Sleep(2 * time.Second)

	// Print all registered routes
	fmt.Println("\n=== Registered Routes ===")
	router.Routes()
	
	// Keep running
	select {}
}
