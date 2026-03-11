package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes configures all category routes
func SetupCategoryRoutes(
	rg *gin.RouterGroup,
	categoryHandler *handler.CategoryHandler,
	tokenService service.TokenService,
) {
	// Public routes (no authentication required)
	public := rg.Group("")
	{
		// Get all categories
		public.GET("/categories", categoryHandler.GetAllCategories)

		// Get category tree (hierarchical structure)
		public.GET("/categories/tree", categoryHandler.GetCategoryTree)

		// Get featured categories (for homepage)
		public.GET("/categories/featured", categoryHandler.GetFeaturedCategories)

		// Get single category
		public.GET("/categories/:id", categoryHandler.GetCategoryByID)

		// Get category breadcrumb
		public.GET("/categories/:id/breadcrumb", categoryHandler.GetCategoryBreadcrumb)

		// Get products by category
		public.GET("/categories/:id/products", categoryHandler.GetProductsByCategory)

		// Search categories
		public.GET("/categories/search", categoryHandler.SearchCategories)
	}

	// Admin routes (require authentication + admin role)
	admin := rg.Group("")
	admin.Use(middleware.JWTAuth(tokenService))
	admin.Use(middleware.RequireAdmin())
	{
		// Create category
		admin.POST("/categories", categoryHandler.CreateCategory)

		// Update category
		admin.PUT("/categories/:id", categoryHandler.UpdateCategory)

		// Delete category
		admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	}
}

// CategoryResponse represents a category API response
type CategoryResponse struct {
	ID             int64   `json:"id"`
	ParentID       *int64  `json:"parent_id"`
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	Description    string  `json:"description"`
	IconURL        string  `json:"icon_url"`
	ImageURL       string  `json:"image_url"`
	Level          int     `json:"level"`
	SortOrder      int     `json:"sort_order"`
	IsActive       bool    `json:"is_active"`
	ProductCount   int64   `json:"product_count,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// CategoryTreeResponse represents a category in tree format
type CategoryTreeResponse struct {
	ID           int64                 `json:"id"`
	ParentID     *int64                `json:"parent_id"`
	Name         string                `json:"name"`
	Slug         string                `json:"slug"`
	Description  string                `json:"description"`
	IconURL      string                `json:"icon_url"`
	ImageURL     string                `json:"image_url"`
	Level        int                   `json:"level"`
	SortOrder    int                   `json:"sort_order"`
	IsActive     bool                  `json:"is_active"`
	ProductCount int64                 `json:"product_count"`
	Children     []CategoryTreeResponse `json:"children"`
}

// CategoryBreadcrumbResponse represents a breadcrumb item
type CategoryBreadcrumbResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// HealthCategoryCheck handles GET /health/category
func HealthCategoryCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "category",
	})
}
