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
	rg.POST("/webhook", paymentHandler.WebhookHandler)

	// Protected routes (auth required)
	protected := rg.Group("")
	protected.Use(middleware.JWTAuth(tokenService))
	{
		// Payment creation and retrieval
		protected.POST("/create", paymentHandler.CreatePayment)
		protected.POST("/confirm", paymentHandler.ConfirmPayment)
		protected.GET("/order/:order_id", paymentHandler.GetPaymentByOrderID)
		protected.GET("", paymentHandler.GetUserPayments)

		// Refunds
		protected.POST("/refund", paymentHandler.RequestRefund)

		// Payment methods
		protected.POST("/methods", paymentHandler.SavePaymentMethod)
		protected.GET("/methods", paymentHandler.GetPaymentMethods)
		protected.DELETE("/methods/:id", paymentHandler.DeletePaymentMethod)
		protected.POST("/methods/:id/default", paymentHandler.SetDefaultPaymentMethod)

		// Statistics
		protected.GET("/statistics", paymentHandler.GetPaymentStats)
	}
}
