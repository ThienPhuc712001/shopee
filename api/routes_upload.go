package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupUploadRoutes configures all image upload routes
func SetupUploadRoutes(
	rg *gin.RouterGroup,
	uploadHandler *handler.UploadHandler,
	tokenService service.TokenService,
) {
	// Protected routes (require authentication)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Product image upload
		protected.POST("/upload/product", uploadHandler.UploadProductImage)
		protected.POST("/upload/product/multiple", uploadHandler.UploadMultipleProductImages)
		protected.DELETE("/upload/product/:id", uploadHandler.DeleteProductImage)
		protected.GET("/upload/product/images", uploadHandler.GetProductImages)

		// Review image upload
		protected.POST("/upload/review", uploadHandler.UploadReviewImage)

		// User avatar upload
		protected.POST("/upload/avatar", uploadHandler.UploadAvatar)
	}
}

// SetupStaticRoutes configures static file serving for uploaded images
func SetupStaticRoutes(router *gin.Engine, uploadDir string) {
	// Serve uploaded files statically
	// URL pattern: /uploads/{type}/{filename}
	// Maps to: ./uploads/{type}/{filename}
	router.Static("/uploads", uploadDir)

	// Health check for upload service
	router.GET("/health/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "upload",
		})
	})
}

// ApiResponse represents a standard API response
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// UploadResponse represents an upload response
type UploadResponse struct {
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	ImageType string `json:"image_type"`
}

// ImageListResponse represents a list of images response
type ImageListResponse struct {
	Images interface{} `json:"images"`
	Count  int         `json:"count"`
}
