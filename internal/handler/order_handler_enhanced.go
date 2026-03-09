package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderHandlerEnhanced handles order-related requests
type OrderHandlerEnhanced struct {
	orderService service.OrderServiceEnhanced
}

// NewOrderHandlerEnhanced creates a new enhanced order handler
func NewOrderHandlerEnhanced(orderService service.OrderServiceEnhanced) *OrderHandlerEnhanced {
	return &OrderHandlerEnhanced{
		orderService: orderService,
	}
}

// CheckoutRequest represents a checkout request
type CheckoutRequest struct {
	ShippingInfo     ShippingInfoRequest `json:"shipping_info" binding:"required"`
	PaymentMethod    string              `json:"payment_method" binding:"required"`
	BuyerNote        string              `json:"buyer_note"`
	VoucherCode      string              `json:"voucher_code"`
	ShippingMethod   string              `json:"shipping_method"`
}

// ShippingInfoRequest represents shipping info in request
type ShippingInfoRequest struct {
	Name        string `json:"name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Ward        string `json:"ward"`
	District    string `json:"district" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state"`
	Country     string `json:"country" default:"Vietnam"`
	PostalCode  string `json:"postal_code"`
}

// CancelOrderRequest represents a cancel order request
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// UpdateOrderStatusRequest represents an update order status request
type UpdateOrderStatusRequest struct {
	Status  string `json:"status" binding:"required"`
	Message string `json:"message"`
}

// Checkout handles the checkout process
// @Summary Checkout cart
// @Description Convert cart to order and proceed to payment
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckoutRequest true "Checkout data"
// @Success 201 {object} response.Response
// @Router /api/orders/checkout [post]
func (h *OrderHandlerEnhanced) Checkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Convert request to service input
	input := &model.OrderInput{
		ShippingInfo: model.ShippingInfo{
			Name:       req.ShippingInfo.Name,
			Phone:      req.ShippingInfo.Phone,
			Address:    req.ShippingInfo.Address,
			Ward:       req.ShippingInfo.Ward,
			District:   req.ShippingInfo.District,
			City:       req.ShippingInfo.City,
			State:      req.ShippingInfo.State,
			Country:    req.ShippingInfo.Country,
			PostalCode: req.ShippingInfo.PostalCode,
		},
		PaymentMethod:  req.PaymentMethod,
		BuyerNote:      req.BuyerNote,
		VoucherCode:    req.VoucherCode,
		ShippingMethod: req.ShippingMethod,
	}

	order, err := h.orderService.CheckoutCart(userID.(uint), input)
	if err != nil {
		switch err {
		case service.ErrCartEmpty:
			c.JSON(http.StatusBadRequest, response.BadRequest("Cart is empty"))
		case service.ErrInvalidShippingInfo:
			c.JSON(http.StatusBadRequest, response.BadRequest("Invalid shipping information"))
		case service.ErrInvalidPaymentMethod:
			c.JSON(http.StatusBadRequest, response.BadRequest("Invalid payment method"))
		case service.ErrInsufficientStock:
			c.JSON(http.StatusBadRequest, response.BadRequest("Insufficient stock available"))
		case service.ErrProductUnavailable:
			c.JSON(http.StatusBadRequest, response.BadRequest("Product is no longer available"))
		case service.ErrInventoryLockFailed:
			c.JSON(http.StatusInternalServerError, response.InternalError("Failed to lock inventory"))
		default:
			c.JSON(http.StatusInternalServerError, response.InternalError("Failed to create order"))
		}
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"order": order,
	}, "Order created successfully"))
}

// GetUserOrders handles getting user's orders
// @Summary Get user orders
// @Description Get all orders for the current user
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/orders [get]
func (h *OrderHandlerEnhanced) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.orderService.GetUserOrders(userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get orders"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"orders": orders,
	}, total, page, limit, ""))
}

// GetOrder handles getting an order by ID
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Router /api/orders/{id} [get]
func (h *OrderHandlerEnhanced) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid order ID"))
		return
	}

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Order not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"order": order,
	}, ""))
}

// CancelOrder handles order cancellation
// @Summary Cancel order
// @Description Cancel an order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body CancelOrderRequest true "Cancellation reason"
// @Success 200 {object} response.Response
// @Router /api/orders/{id}/cancel [post]
func (h *OrderHandlerEnhanced) CancelOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid order ID"))
		return
	}

	var req CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	order, err := h.orderService.CancelOrder(uint(id), req.Reason, userID.(uint))
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, response.NotFound("Order not found"))
			return
		}
		if err == service.ErrOrderCannotCancel {
			c.JSON(http.StatusBadRequest, response.BadRequest("Order cannot be cancelled at this stage"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to cancel order"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"order": order,
	}, "Order cancelled successfully"))
}

// UpdateOrderStatus handles updating order status (admin/seller)
// @Summary Update order status
// @Description Update order status (Seller/Admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body UpdateOrderStatusRequest true "Status data"
// @Success 200 {object} response.Response
// @Router /api/orders/{id}/status [put]
func (h *OrderHandlerEnhanced) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid order ID"))
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Get user info from context
	userID, _ := c.Get("user_id")
	ipAddress := c.ClientIP()

	order, err := h.orderService.UpdateOrderStatus(uint(id), model.OrderStatus(req.Status), userID.(uint), ipAddress)
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, response.NotFound("Order not found"))
			return
		}
		if err == service.ErrInvalidOrderStatus {
			c.JSON(http.StatusBadRequest, response.BadRequest("Invalid status transition"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to update order status"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"order": order,
	}, "Order status updated successfully"))
}

// GetOrderTracking handles getting order tracking
// @Summary Get order tracking
// @Description Get tracking information for an order
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Router /api/orders/{id}/tracking [get]
func (h *OrderHandlerEnhanced) GetOrderTracking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid order ID"))
		return
	}

	tracking, err := h.orderService.GetOrderTracking(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get tracking"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"tracking": tracking,
	}, ""))
}

// GetOrderStatistics handles getting order statistics
// @Summary Get order statistics
// @Description Get statistics for user's orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/orders/statistics [get]
func (h *OrderHandlerEnhanced) GetOrderStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	stats, err := h.orderService.GetOrderStatistics(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get statistics"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"statistics": stats,
	}, ""))
}
