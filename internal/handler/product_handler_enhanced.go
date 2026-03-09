package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProductHandlerEnhanced handles product requests
type ProductHandlerEnhanced struct {
	productService service.ProductServiceEnhanced
}

// NewProductHandlerEnhanced creates a new enhanced product handler
func NewProductHandlerEnhanced(productService service.ProductServiceEnhanced) *ProductHandlerEnhanced {
	return &ProductHandlerEnhanced{
		productService: productService,
	}
}

// ==================== REQUEST/RESPONSE STRUCTS ====================

// CreateProductRequest represents a create product request
type CreateProductRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Price            float64 `json:"price" binding:"required,gt=0"`
	OriginalPrice    float64 `json:"original_price"`
	Stock            int     `json:"stock" binding:"required,gte=0"`
	CategoryID       uint    `json:"category_id" binding:"required"`
	Brand            string  `json:"brand"`
	Weight           float64 `json:"weight"`
	Dimensions       string  `json:"dimensions"`
	WarrantyPeriod   string  `json:"warranty_period"`
	ReturnDays       int     `json:"return_days"`
	Tags             []string `json:"tags"`
	MetaTitle        string  `json:"meta_title"`
	MetaDescription  string  `json:"meta_description"`
}

// UpdateProductRequest represents an update product request
type UpdateProductRequest struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Price            float64 `json:"price"`
	OriginalPrice    float64 `json:"original_price"`
	Stock            int     `json:"stock"`
	CategoryID       uint    `json:"category_id"`
	Brand            string  `json:"brand"`
	Weight           float64 `json:"weight"`
	Status           string  `json:"status"`
}

// CreateVariantRequest represents a create variant request
type CreateVariantRequest struct {
	SKU         string             `json:"sku"`
	Name        string             `json:"name"`
	Price       float64            `json:"price"`
	OriginalPrice float64          `json:"original_price"`
	Stock       int                `json:"stock" binding:"gte=0"`
	Attributes  map[string]string  `json:"attributes"`
	ImageURL    string             `json:"image_url"`
}

// UpdateStockRequest represents an update stock request
type UpdateStockRequest struct {
	Quantity  int    `json:"quantity" binding:"required,gte=0"`
	VariantID *uint  `json:"variant_id"`
}

// ProductFiltersRequest represents filter parameters
type ProductFiltersRequest struct {
	CategoryID  uint     `form:"category_id"`
	ShopID      uint     `form:"shop_id"`
	MinPrice    float64  `form:"min_price"`
	MaxPrice    float64  `form:"max_price"`
	MinRating   float64  `form:"min_rating"`
	Brands      string   `form:"brands"` // comma-separated
	SortBy      string   `form:"sort_by"`
	SortOrder   string   `form:"sort_order"`
	Page        int      `form:"page"`
	Limit       int      `form:"limit"`
}

// ==================== PRODUCT CRUD ====================

// CreateProduct handles product creation
// @Summary Create a product
// @Description Create a new product (Seller only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductRequest true "Product data"
// @Success 201 {object} response.Response
// @Router /api/products [post]
func (h *ProductHandlerEnhanced) CreateProduct(c *gin.Context) {
	// Get shop ID from context (set by auth middleware)
	shopIDValue, exists := c.Get("shop_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("Shop ID not found"))
		return
	}
	shopID := shopIDValue.(uint)

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Convert request to model
	product := &model.Product{
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		OriginalPrice:    req.OriginalPrice,
		Stock:            req.Stock,
		CategoryID:       req.CategoryID,
		Brand:            req.Brand,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		WarrantyPeriod:   req.WarrantyPeriod,
		ReturnDays:       req.ReturnDays,
		Status:           model.ProductStatusDraft,
	}

	// Create product
	createdProduct, err := h.productService.CreateProduct(product, shopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"product": createdProduct,
	}, "Product created successfully"))
}

// UpdateProduct handles product update
// @Summary Update product
// @Description Update an existing product (Seller only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body UpdateProductRequest true "Product data"
// @Success 200 {object} response.Response
// @Router /api/products/{id} [put]
func (h *ProductHandlerEnhanced) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	shopIDValue, exists := c.Get("shop_id")
	isAdmin := c.GetBool("is_admin")
	shopID := uint(0)
	if exists {
		shopID = shopIDValue.(uint)
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	product := &model.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		CategoryID:    req.CategoryID,
		Status:        model.ProductStatus(req.Status),
	}

	updatedProduct, err := h.productService.UpdateProduct(uint(id), product, shopID, isAdmin)
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, response.NotFound("Product not found"))
			return
		}
		if err == service.ErrUnauthorizedSeller {
			c.JSON(http.StatusForbidden, response.Forbidden("You can only update your own products"))
			return
		}
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"product": updatedProduct,
	}, "Product updated successfully"))
}

// DeleteProduct handles product deletion
// @Summary Delete product
// @Description Delete a product (Admin or owner)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/products/{id} [delete]
func (h *ProductHandlerEnhanced) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	shopIDValue, exists := c.Get("shop_id")
	isAdmin := c.GetBool("is_admin")
	shopID := uint(0)
	if exists {
		shopID = shopIDValue.(uint)
	}

	if err := h.productService.DeleteProduct(uint(id), shopID, isAdmin); err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, response.NotFound("Product not found"))
			return
		}
		if err == service.ErrUnauthorizedSeller {
			c.JSON(http.StatusForbidden, response.Forbidden("You can only delete your own products"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to delete product"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Product deleted successfully"))
}

// GetProduct handles getting a product by ID
// @Summary Get product by ID
// @Description Get product details
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/products/{id} [get]
func (h *ProductHandlerEnhanced) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Product not found"))
		return
	}

	// Increment view count (async in production)
	_ = h.productService.IncrementViewCount(uint(id))

	c.JSON(http.StatusOK, response.Success(gin.H{
		"product": product,
	}, ""))
}

// ==================== PRODUCT LISTING ====================

// GetProducts handles getting all products
// @Summary Get all products
// @Description Get all products with pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/products [get]
func (h *ProductHandlerEnhanced) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productService.GetProducts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get products"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"products": products,
	}, total, page, limit, ""))
}

// SearchProducts handles product search
// @Summary Search products
// @Description Search products with filters
// @Tags products
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param category_id query int false "Category ID"
// @Param min_price query float64 false "Minimum price"
// @Param max_price query float64 false "Maximum price"
// @Param sort_by query string false "Sort by (price, rating, sold_count, created_at)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/products/search [get]
func (h *ProductHandlerEnhanced) SearchProducts(c *gin.Context) {
	var req ProductFiltersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	keyword := c.Query("keyword")

	// Build filters
	filters := service.ProductFilters{
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.CategoryID > 0 {
		filters.CategoryID = &req.CategoryID
	}
	if req.ShopID > 0 {
		filters.ShopID = &req.ShopID
	}
	if req.MinPrice > 0 {
		filters.MinPrice = &req.MinPrice
	}
	if req.MaxPrice > 0 {
		filters.MaxPrice = &req.MaxPrice
	}
	if req.MinRating > 0 {
		filters.MinRating = &req.MinRating
	}
	if req.Brands != "" {
		filters.Brands = strings.Split(req.Brands, ",")
	}

	products, total, err := h.productService.SearchProducts(keyword, filters, req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to search products"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"products": products,
	}, total, req.Page, req.Limit, ""))
}

// GetProductsByCategory handles getting products by category
// @Summary Get products by category
// @Description Get products in a category
// @Tags products
// @Produce json
// @Param id path int true "Category ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/products/category/{id} [get]
func (h *ProductHandlerEnhanced) GetProductsByCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid category ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productService.GetProductsByCategoryID(uint(id), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get products"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"products": products,
	}, total, page, limit, ""))
}

// GetBestSellers handles getting best selling products
// @Summary Get best sellers
// @Description Get best selling products
// @Tags products
// @Produce json
// @Param limit query int false "Number of products"
// @Success 200 {object} response.Response
// @Router /api/products/best-sellers [get]
func (h *ProductHandlerEnhanced) GetBestSellers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.productService.GetBestSellers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get best sellers"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"products": products,
	}, ""))
}

// GetFeaturedProducts handles getting featured products
// @Summary Get featured products
// @Description Get featured products
// @Tags products
// @Produce json
// @Param limit query int false "Number of products"
// @Success 200 {object} response.Response
// @Router /api/products/featured [get]
func (h *ProductHandlerEnhanced) GetFeaturedProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get featured products"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"products": products,
	}, ""))
}

// ==================== PRODUCT IMAGES ====================

// UploadProductImages handles image upload
// @Summary Upload product images
// @Description Upload images for a product (Seller only)
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param images formData []file true "Image files"
// @Success 200 {object} response.Response
// @Router /api/products/{id}/images [post]
func (h *ProductHandlerEnhanced) UploadProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	_, exists := c.Get("shop_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("Shop ID not found"))
		return
	}

	// Get uploaded files
	files := c.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, response.BadRequest("No images uploaded"))
		return
	}

	// Upload images
	uploadDir := "./uploads/products"
	images, err := h.productService.UploadProductImages(uint(id), files, uploadDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"images": images,
	}, "Images uploaded successfully"))
}

// ==================== PRODUCT VARIANTS ====================

// CreateVariant handles variant creation
// @Summary Create product variant
// @Description Create a variant for a product (Seller only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body CreateVariantRequest true "Variant data"
// @Success 201 {object} response.Response
// @Router /api/products/{id}/variants [post]
func (h *ProductHandlerEnhanced) CreateVariant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	var req CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Convert attributes to JSON
	// In production, properly serialize the map
	attributesJSON := "{}"

	variant := &model.ProductVariant{
		SKU:         req.SKU,
		Name:        req.Name,
		Price:       req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:       req.Stock,
		Attributes:  attributesJSON,
		ImageURL:    req.ImageURL,
	}

	createdVariant, err := h.productService.CreateVariant(uint(id), variant)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"variant": createdVariant,
	}, "Variant created successfully"))
}

// GetVariants handles getting product variants
// @Summary Get product variants
// @Description Get all variants for a product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/products/{id}/variants [get]
func (h *ProductHandlerEnhanced) GetVariants(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	variants, err := h.productService.GetVariantsByProductID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get variants"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"variants": variants,
	}, ""))
}

// ==================== INVENTORY ====================

// UpdateStock handles stock update
// @Summary Update product stock
// @Description Update stock for a product or variant (Seller only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body UpdateStockRequest true "Stock data"
// @Success 200 {object} response.Response
// @Router /api/products/{id}/stock [put]
func (h *ProductHandlerEnhanced) UpdateStock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	shopIDValue, exists := c.Get("shop_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("Shop ID not found"))
		return
	}
	shopID := shopIDValue.(uint)

	var req UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if err := h.productService.UpdateProductStock(uint(id), req.VariantID, req.Quantity, shopID); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Stock updated successfully"))
}

// ==================== CATEGORIES ====================

// GetCategories handles getting all categories
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/categories [get]
func (h *ProductHandlerEnhanced) GetCategories(c *gin.Context) {
	categories, err := h.productService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get categories"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"categories": categories,
	}, ""))
}

// GetCategoryByID handles getting a category
// @Summary Get category by ID
// @Description Get category details
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response
// @Router /api/categories/{id} [get]
func (h *ProductHandlerEnhanced) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid category ID"))
		return
	}

	category, err := h.productService.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Category not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"category": category,
	}, ""))
}

// ==================== PRODUCT STATUS ====================

// PublishProduct handles product publishing
// @Summary Publish product
// @Description Publish a product (Seller only)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/products/{id}/publish [patch]
func (h *ProductHandlerEnhanced) PublishProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	shopIDValue, exists := c.Get("shop_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("Shop ID not found"))
		return
	}
	shopID := shopIDValue.(uint)

	if err := h.productService.PublishProduct(uint(id), shopID); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Product published successfully"))
}

// UnpublishProduct handles product unpublishing
// @Summary Unpublish product
// @Description Unpublish a product (Seller only)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/products/{id}/unpublish [patch]
func (h *ProductHandlerEnhanced) UnpublishProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	shopIDValue, exists := c.Get("shop_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("Shop ID not found"))
		return
	}
	shopID := shopIDValue.(uint)

	if err := h.productService.UnpublishProduct(uint(id), shopID); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Product unpublished successfully"))
}
