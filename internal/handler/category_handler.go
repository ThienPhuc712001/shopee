package handler

import (
	"ecommerce/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory handles POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input service.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNameRequired) || errors.Is(err, service.ErrCategoryNameLength) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrParentCategoryNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create category",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Category created successfully",
		"data":    category,
	})
}

// UpdateCategory handles PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category ID",
		})
		return
	}

	var input service.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), uint(id), input)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrCircularHierarchy) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrParentCategoryNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update category",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
		"data":    category,
	})
}

// DeleteCategory handles DELETE /api/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category ID",
		})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrCategoryHasChildren) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrCategoryHasProducts) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete category",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted successfully",
	})
}

// GetCategoryByID handles GET /api/categories/:id
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category ID",
		})
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
	})
}

// GetAllCategories handles GET /api/categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	activeParam := c.Query("active")
	
	var isActive *bool
	if activeParam != "" {
		active := activeParam == "true"
		isActive = &active
	}

	categories, err := h.categoryService.GetAllCategories(c.Request.Context(), isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get categories",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"categories": categories,
			"count":      len(categories),
		},
	})
}

// GetCategoryTree handles GET /api/categories/tree
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.categoryService.GetCategoryTreeWithCounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get category tree",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"tree": tree,
		},
	})
}

// GetCategoryBreadcrumb handles GET /api/categories/:id/breadcrumb
func (h *CategoryHandler) GetCategoryBreadcrumb(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category ID",
		})
		return
	}

	breadcrumbs, err := h.categoryService.GetCategoryPath(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"breadcrumb": breadcrumbs,
		},
	})
}

// GetProductsByCategory handles GET /api/categories/:id/products
func (h *CategoryHandler) GetProductsByCategory(c *gin.Context) {
	idStr := c.Param("id")
	categoryID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid category ID",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	products, total, err := h.categoryService.GetProductsByCategory(c.Request.Context(), uint(categoryID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get products",
			"error":   err.Error(),
		})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"products": products,
			"pagination": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// GetFeaturedCategories handles GET /api/categories/featured
func (h *CategoryHandler) GetFeaturedCategories(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "8")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 20 {
		limit = 8
	}

	categories, err := h.categoryService.GetFeaturedCategories(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get featured categories",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"categories": categories,
			"count":      len(categories),
		},
	})
}

// SearchCategories handles GET /api/categories/search
func (h *CategoryHandler) SearchCategories(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Search query is required",
		})
		return
	}

	categories, err := h.categoryService.SearchCategories(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to search categories",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"categories": categories,
			"count":      len(categories),
		},
	})
}
