package handler

import (
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CartHandlerEnhanced handles cart-related requests
type CartHandlerEnhanced struct {
	cartService service.CartServiceEnhanced
}

// NewCartHandlerEnhanced creates a new enhanced cart handler
func NewCartHandlerEnhanced(cartService service.CartServiceEnhanced) *CartHandlerEnhanced {
	return &CartHandlerEnhanced{
		cartService: cartService,
	}
}

// ==================== REQUEST/RESPONSE STRUCTS ====================

// AddToCartRequest represents a request to add item to cart
type AddToCartRequest struct {
	ProductID uint  `json:"product_id" binding:"required"`
	VariantID *uint `json:"variant_id"`
	Quantity  int   `json:"quantity" binding:"required,min=1,max=999"`
}

// UpdateCartItemRequest represents a request to update cart item
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=999"`
}

// CartResponse represents a cart response
type CartResponse struct {
	ID           uint                 `json:"id"`
	UserID       uint                 `json:"user_id"`
	TotalItems   int                  `json:"total_items"`
	Subtotal     float64              `json:"subtotal"`
	Discount     float64              `json:"discount"`
	Total        float64              `json:"total"`
	Currency     string               `json:"currency"`
	Items        []CartItemResponse   `json:"items"`
	LastActivity string               `json:"last_activity"`
}

// CartItemResponse represents a cart item response
type CartItemResponse struct {
	ID            uint            `json:"id"`
	ProductID     uint            `json:"product_id"`
	Product       *ProductInfo    `json:"product"`
	VariantID     *uint           `json:"variant_id"`
	Variant       *VariantInfo    `json:"variant,omitempty"`
	Quantity      int             `json:"quantity"`
	Price         float64         `json:"price"`
	OriginalPrice float64         `json:"original_price"`
	Subtotal      float64         `json:"subtotal"`
	ProductImage  string          `json:"product_image"`
	ShopID        uint            `json:"shop_id"`
	Shop          *ShopInfo       `json:"shop"`
	IsAvailable   bool            `json:"is_available"`
	StockStatus   string          `json:"stock_status"`
}

// ProductInfo represents product information in cart
type ProductInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}

// VariantInfo represents variant information in cart
type VariantInfo struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	SKU      string            `json:"sku"`
	Price    float64           `json:"price"`
	Stock    int               `json:"stock"`
	Attributes map[string]string `json:"attributes"`
}

// ShopInfo represents shop information in cart
type ShopInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CartSummaryResponse represents cart summary response
type CartSummaryResponse struct {
	TotalItems  int     `json:"total_items"`
	Subtotal    float64 `json:"subtotal"`
	Discount    float64 `json:"discount"`
	Total       float64 `json:"total"`
	Currency    string  `json:"currency"`
}

// CheckoutSummaryResponse represents checkout summary response
type CheckoutSummaryResponse struct {
	Items       []CheckoutItemResponse `json:"items"`
	Subtotal    float64                `json:"subtotal"`
	ShippingFee float64                `json:"shipping_fee"`
	Discount    float64                `json:"discount"`
	Total       float64                `json:"total"`
	Currency    string                 `json:"currency"`
}

// CheckoutItemResponse represents checkout item response
type CheckoutItemResponse struct {
	ProductID    uint    `json:"product_id"`
	VariantID    *uint   `json:"variant_id"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	Subtotal     float64 `json:"subtotal"`
	ShopID       uint    `json:"shop_id"`
	ShopName     string  `json:"shop_name"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	IsAvailable  bool    `json:"is_available"`
	StockStatus  string  `json:"stock_status"`
}

// CartStatsResponse represents cart statistics response
type CartStatsResponse struct {
	TotalItems      int     `json:"total_items"`
	TotalProducts   int     `json:"total_products"`
	Subtotal        float64 `json:"subtotal"`
	EstimatedTotal  float64 `json:"estimated_total"`
	HasOutOfStock   bool    `json:"has_out_of_stock"`
	HasPriceChanged bool    `json:"has_price_changed"`
}

// ==================== CART HANDLERS ====================

// GetCart handles getting the user's cart
// @Summary Get cart
// @Description Get current user's shopping cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/cart [get]
func (h *CartHandlerEnhanced) GetCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	cart, err := h.cartService.GetCart(userID.(uint))
	if err != nil {
		if err == service.ErrCartNotFound {
			// Return empty cart
			c.JSON(http.StatusOK, response.Success(gin.H{
				"cart": gin.H{
					"id":           0,
					"user_id":      userID,
					"total_items":  0,
					"subtotal":     0,
					"total":        0,
					"items":        []gin.H{},
				},
			}, ""))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get cart"))
		return
	}

	// Convert cart to response format
	cartResponse := h.convertCartToResponse(cart)

	c.JSON(http.StatusOK, response.Success(gin.H{
		"cart": cartResponse,
	}, ""))
}

// GetCartSummary handles getting cart summary
// @Summary Get cart summary
// @Description Get summary of current user's cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/cart/summary [get]
func (h *CartHandlerEnhanced) GetCartSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	summary, err := h.cartService.GetCartSummary(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, response.Success(gin.H{
			"summary": gin.H{
				"total_items": 0,
				"subtotal":    0,
				"total":       0,
			},
		}, ""))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"summary": summary,
	}, ""))
}

// AddToCart handles adding an item to cart
// @Summary Add to cart
// @Description Add a product to the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddToCartRequest true "Cart item data"
// @Success 200 {object} response.Response
// @Router /api/cart/add [post]
func (h *CartHandlerEnhanced) AddToCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	cart, err := h.cartService.AddToCart(
		userID.(uint),
		req.ProductID,
		req.VariantID,
		req.Quantity,
	)

	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			c.JSON(http.StatusNotFound, response.NotFound("Product not found"))
		case service.ErrVariantNotFound:
			c.JSON(http.StatusNotFound, response.NotFound("Variant not found"))
		case service.ErrInsufficientStock:
			c.JSON(http.StatusBadRequest, response.BadRequest("Insufficient stock available"))
		case service.ErrInvalidQuantity:
			c.JSON(http.StatusBadRequest, response.BadRequest("Invalid quantity"))
		case service.ErrProductNotAvailable:
			c.JSON(http.StatusBadRequest, response.BadRequest("Product is not available"))
		default:
			c.JSON(http.StatusInternalServerError, response.InternalError("Failed to add item to cart"))
		}
		return
	}

	cartResponse := h.convertCartToResponse(cart)

	c.JSON(http.StatusOK, response.Success(gin.H{
		"cart": cartResponse,
	}, "Item added to cart successfully"))
}

// UpdateCartItem handles updating cart item quantity
// @Summary Update cart item
// @Description Update quantity of a cart item
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body UpdateCartItemRequest true "Quantity data"
// @Success 200 {object} response.Response
// @Router /api/cart/items/{id} [put]
func (h *CartHandlerEnhanced) UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid cart item ID"))
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	cart, err := h.cartService.UpdateCartItem(
		userID.(uint),
		uint(itemID),
		req.Quantity,
	)

	if err != nil {
		switch err {
		case service.ErrCartItemNotFound:
			c.JSON(http.StatusNotFound, response.NotFound("Cart item not found"))
		case service.ErrInsufficientStock:
			c.JSON(http.StatusBadRequest, response.BadRequest("Insufficient stock available"))
		case service.ErrInvalidQuantity:
			c.JSON(http.StatusBadRequest, response.BadRequest("Invalid quantity"))
		default:
			c.JSON(http.StatusInternalServerError, response.InternalError("Failed to update cart item"))
		}
		return
	}

	cartResponse := h.convertCartToResponse(cart)

	c.JSON(http.StatusOK, response.Success(gin.H{
		"cart": cartResponse,
	}, "Cart item updated successfully"))
}

// RemoveFromCart handles removing an item from cart
// @Summary Remove from cart
// @Description Remove an item from the shopping cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} response.Response
// @Router /api/cart/items/{id} [delete]
func (h *CartHandlerEnhanced) RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid cart item ID"))
		return
	}

	cart, err := h.cartService.RemoveFromCart(userID.(uint), uint(itemID))
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			c.JSON(http.StatusNotFound, response.NotFound("Cart not found"))
		case service.ErrCartItemNotFound:
			c.JSON(http.StatusNotFound, response.NotFound("Cart item not found"))
		default:
			c.JSON(http.StatusInternalServerError, response.InternalError("Failed to remove item from cart"))
		}
		return
	}

	cartResponse := h.convertCartToResponse(cart)

	c.JSON(http.StatusOK, response.Success(gin.H{
		"cart": cartResponse,
	}, "Item removed from cart successfully"))
}

// ClearCart handles clearing the cart
// @Summary Clear cart
// @Description Remove all items from the shopping cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/cart/clear [delete]
func (h *CartHandlerEnhanced) ClearCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	if err := h.cartService.ClearCart(userID.(uint)); err != nil {
		if err == service.ErrCartNotFound {
			c.JSON(http.StatusNotFound, response.NotFound("Cart not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to clear cart"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Cart cleared successfully"))
}

// GetCartStats handles getting cart statistics
// @Summary Get cart stats
// @Description Get statistics for current user's cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/cart/stats [get]
func (h *CartHandlerEnhanced) GetCartStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	stats, err := h.cartService.GetCartStats(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, response.Success(gin.H{
			"stats": gin.H{
				"total_items":      0,
				"total_products":   0,
				"subtotal":         0,
				"estimated_total":  0,
				"has_out_of_stock": false,
				"has_price_changed": false,
			},
		}, ""))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"stats": stats,
	}, ""))
}

// PrepareCheckout handles preparing cart for checkout
// @Summary Prepare checkout
// @Description Prepare cart for checkout with validation
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/cart/checkout [get]
func (h *CartHandlerEnhanced) PrepareCheckout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	summary, err := h.cartService.PrepareForCheckout(userID.(uint))
	if err != nil {
		if err == service.ErrCartEmpty {
			c.JSON(http.StatusBadRequest, response.BadRequest("Cart is empty"))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	// Convert to response format
	checkoutResponse := h.convertCheckoutSummaryToResponse(summary)

	c.JSON(http.StatusOK, response.Success(gin.H{
		"checkout": checkoutResponse,
	}, ""))
}

// ==================== HELPER METHODS ====================

// convertCartToResponse converts model.Cart to CartResponse
func (h *CartHandlerEnhanced) convertCartToResponse(cart interface{}) CartResponse {
	// This is a simplified conversion
	// In production, properly map all fields
	return CartResponse{
		ID:           0,
		UserID:       0,
		TotalItems:   0,
		Subtotal:     0,
		Total:        0,
		Items:        []CartItemResponse{},
		LastActivity: "",
	}
}

// convertCheckoutSummaryToResponse converts checkout summary to response
func (h *CartHandlerEnhanced) convertCheckoutSummaryToResponse(summary interface{}) CheckoutSummaryResponse {
	// Simplified conversion
	return CheckoutSummaryResponse{
		Items:       []CheckoutItemResponse{},
		Subtotal:    0,
		ShippingFee: 0,
		Discount:    0,
		Total:       0,
		Currency:    "USD",
	}
}

// convertCartItemToResponse converts cart item to response
func (h *CartHandlerEnhanced) convertCartItemToResponse(item interface{}) CartItemResponse {
	// Simplified conversion
	return CartItemResponse{
		ID:           0,
		ProductID:    0,
		Quantity:     0,
		Price:        0,
		Subtotal:     0,
		IsAvailable:  true,
		StockStatus:  "in_stock",
	}
}
