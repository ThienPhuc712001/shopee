package handler

import (
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShopHandler handles shop-related HTTP requests
type ShopHandler struct {
	shopService service.ShopService
}

// NewShopHandler creates a new shop handler
func NewShopHandler(shopService service.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

// CreateShop handles shop creation
// @Summary Create a new shop
// @Description Create a new shop for the authenticated user
// @Tags shops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateShopInput true "Shop details"
// @Success 201 {object} response.Response
// @Router /api/shops [post]
func (h *ShopHandler) CreateShop(c *gin.Context) {
	// Get user ID from context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid user ID"))
		return
	}

	var req service.CreateShopInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	shop, err := h.shopService.CreateShop(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"shop": shop,
	}, "Shop created successfully"))
}

// GetShop handles getting a shop by ID
// @Summary Get shop details
// @Description Get details of a specific shop
// @Tags shops
// @Produce json
// @Param id path int true "Shop ID"
// @Success 200 {object} response.Response
// @Router /api/shops/:id [get]
func (h *ShopHandler) GetShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid shop ID"))
		return
	}

	shop, err := h.shopService.GetShopByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Shop not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"shop": shop,
	}, "Shop retrieved successfully"))
}

// GetMyShop handles getting the current user's shop
// @Summary Get my shop
// @Description Get the shop owned by the authenticated user
// @Tags shops
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/shops/my [get]
func (h *ShopHandler) GetMyShop(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid user ID"))
		return
	}

	shop, err := h.shopService.GetShopByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Shop not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"shop": shop,
	}, "Shop retrieved successfully"))
}
