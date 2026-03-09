package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaymentHandlerEnhanced handles payment-related requests
type PaymentHandlerEnhanced struct {
	paymentService service.PaymentServiceEnhanced
}

// NewPaymentHandlerEnhanced creates a new enhanced payment handler
func NewPaymentHandlerEnhanced(paymentService service.PaymentServiceEnhanced) *PaymentHandlerEnhanced {
	return &PaymentHandlerEnhanced{
		paymentService: paymentService,
	}
}

// CreatePaymentRequest represents a create payment request
type CreatePaymentRequest struct {
	OrderID       uint   `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Provider      string `json:"provider"`
	SaveMethod    bool   `json:"save_method"`
}

// ConfirmPaymentRequest represents a confirm payment request
type ConfirmPaymentRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID uint   `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount"`
	Reason    string `json:"reason" binding:"required"`
	Type      string `json:"type" binding:"required"` // full or partial
}

// SavePaymentMethodRequest represents a save payment method request
type SavePaymentMethodRequest struct {
	Type        string `json:"type" binding:"required"`
	Provider    string `json:"provider" binding:"required"`
	Name        string `json:"name"`
	LastFour    string `json:"last_four"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	Token       string `json:"token" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}

// CreatePayment handles creating a new payment
// @Summary Create payment
// @Description Create a new payment for an order
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePaymentRequest true "Payment data"
// @Success 201 {object} response.Response
// @Router /api/payments/create [post]
func (h *PaymentHandlerEnhanced) CreatePayment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	payment, err := h.paymentService.CreatePayment(userID.(uint), req.OrderID, &model.PaymentInput{
		OrderID:       req.OrderID,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
		Provider:      req.Provider,
		SaveMethod:    req.SaveMethod,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"payment": payment,
	}, "Payment created successfully"))
}

// GetPaymentByOrderID handles getting payment by order ID
// @Summary Get payment by order
// @Description Get payment information for an order
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.Response
// @Router /api/payments/order/{order_id} [get]
func (h *PaymentHandlerEnhanced) GetPaymentByOrderID(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid order ID"))
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Payment not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"payment": payment,
	}, ""))
}

// ConfirmPayment handles confirming a payment
// @Summary Confirm payment
// @Description Confirm a payment (for manual confirmation flows)
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConfirmPaymentRequest true "Confirmation data"
// @Success 200 {object} response.Response
// @Router /api/payments/confirm [post]
func (h *PaymentHandlerEnhanced) ConfirmPayment(c *gin.Context) {
	var req ConfirmPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	payment, err := h.paymentService.ConfirmPayment(req.TransactionID)
	if err != nil {
		if err == service.ErrPaymentNotFound {
			c.JSON(http.StatusNotFound, response.NotFound("Payment not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to confirm payment"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"payment": payment,
	}, "Payment confirmed successfully"))
}

// RequestRefund handles requesting a refund
// @Summary Request refund
// @Description Request a refund for a payment
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RefundRequest true "Refund data"
// @Success 201 {object} response.Response
// @Router /api/payments/refund [post]
func (h *PaymentHandlerEnhanced) RequestRefund(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	refund, err := h.paymentService.RequestRefund(userID.(uint), &model.RefundInput{
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Type:      req.Type,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"refund": refund,
	}, "Refund requested successfully"))
}

// GetUserPayments handles getting user's payments
// @Summary Get user payments
// @Description Get all payments for the current user
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/payments [get]
func (h *PaymentHandlerEnhanced) GetUserPayments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	payments, total, err := h.paymentService.GetUserPayments(userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get payments"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"payments": payments,
	}, total, page, limit, ""))
}

// SavePaymentMethod handles saving a payment method
// @Summary Save payment method
// @Description Save a payment method for future use
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SavePaymentMethodRequest true "Payment method data"
// @Success 201 {object} response.Response
// @Router /api/payments/methods [post]
func (h *PaymentHandlerEnhanced) SavePaymentMethod(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	var req SavePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	method, err := h.paymentService.SavePaymentMethod(userID.(uint), &model.PaymentMethodInput{
		Type:        req.Type,
		Provider:    req.Provider,
		Name:        req.Name,
		LastFour:    req.LastFour,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		Token:       req.Token,
		IsDefault:   req.IsDefault,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"method": gin.H{
			"id":         method.ID,
			"type":       method.Type,
			"provider":   method.Provider,
			"name":       method.Name,
			"last_four":  method.LastFour,
			"is_default": method.IsDefault,
		},
	}, "Payment method saved successfully"))
}

// GetPaymentMethods handles getting user's saved payment methods
// @Summary Get payment methods
// @Description Get all saved payment methods for the current user
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/payments/methods [get]
func (h *PaymentHandlerEnhanced) GetPaymentMethods(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	methods, err := h.paymentService.GetPaymentMethods(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get payment methods"))
		return
	}

	// Don't expose sensitive token data
	sanitizedMethods := make([]gin.H, len(methods))
	for i, method := range methods {
		sanitizedMethods[i] = gin.H{
			"id":         method.ID,
			"type":       method.Type,
			"provider":   method.Provider,
			"name":       method.Name,
			"last_four":  method.LastFour,
			"is_default": method.IsDefault,
		}
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"methods": sanitizedMethods,
	}, ""))
}

// DeletePaymentMethod handles deleting a payment method
// @Summary Delete payment method
// @Description Delete a saved payment method
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment Method ID"
// @Success 200 {object} response.Response
// @Router /api/payments/methods/{id} [delete]
func (h *PaymentHandlerEnhanced) DeletePaymentMethod(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	methodID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid payment method ID"))
		return
	}

	if err := h.paymentService.DeletePaymentMethod(userID.(uint), uint(methodID)); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Payment method deleted successfully"))
}

// SetDefaultPaymentMethod handles setting default payment method
// @Summary Set default payment method
// @Description Set a payment method as default
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment Method ID"
// @Success 200 {object} response.Response
// @Router /api/payments/methods/{id}/default [post]
func (h *PaymentHandlerEnhanced) SetDefaultPaymentMethod(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	methodID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid payment method ID"))
		return
	}

	if err := h.paymentService.SetDefaultPaymentMethod(userID.(uint), uint(methodID)); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Default payment method updated"))
}

// GetPaymentStats handles getting payment statistics
// @Summary Get payment statistics
// @Description Get statistics for user's payments
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/payments/statistics [get]
func (h *PaymentHandlerEnhanced) GetPaymentStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not authenticated"))
		return
	}

	stats, err := h.paymentService.GetPaymentStats(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get statistics"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"statistics": stats,
	}, ""))
}

// WebhookHandler handles payment gateway webhooks
// @Summary Payment webhook
// @Description Handle payment gateway webhook notifications
// @Tags payments
// @Accept json
// @Produce json
// @Param request body model.PaymentWebhookInput true "Webhook payload"
// @Success 200 {object} response.Response
// @Router /api/payments/webhook [post]
func (h *PaymentHandlerEnhanced) WebhookHandler(c *gin.Context) {
	var input model.PaymentWebhookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if err := h.paymentService.HandleWebhook(&input); err != nil {
		if err == service.ErrInvalidSignature {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid signature"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Webhook processing failed"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Webhook processed successfully"))
}
