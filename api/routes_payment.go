package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupPaymentRoutes configures all payment routes
func SetupPaymentRoutes(
	rg *gin.RouterGroup,
	paymentHandler *handler.PaymentHandlerEnhanced,
	tokenService service.TokenService,
) {
	// Public routes (no auth required for webhook)
	public := rg.Group("")
	{
		public.POST("/payments/webhook", paymentHandler.WebhookHandler)
	}

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Payment creation and retrieval
		protected.POST("/payments/create", paymentHandler.CreatePayment)
		protected.POST("/payments/confirm", paymentHandler.ConfirmPayment)
		protected.GET("/payments/order/:order_id", paymentHandler.GetPaymentByOrderID)
		protected.GET("/payments", paymentHandler.GetUserPayments)

		// Refunds
		protected.POST("/payments/refund", paymentHandler.RequestRefund)

		// Payment methods
		protected.POST("/payments/methods", paymentHandler.SavePaymentMethod)
		protected.GET("/payments/methods", paymentHandler.GetPaymentMethods)
		protected.DELETE("/payments/methods/:id", paymentHandler.DeletePaymentMethod)
		protected.POST("/payments/methods/:id/default", paymentHandler.SetDefaultPaymentMethod)

		// Statistics
		protected.GET("/payments/statistics", paymentHandler.GetPaymentStats)
	}
}
