package routes

import (
	"ecommerce/internal/handler"
	"ecommerce/internal/service"
	"ecommerce/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupCouponRoutes configures all coupon routes
func SetupCouponRoutes(
	rg *gin.RouterGroup,
	couponHandler *handler.CouponHandler,
	tokenService service.TokenService,
) {
	// Public routes (no authentication required)
	public := rg.Group("")
	{
		// Get active coupons
		public.GET("/coupons/active", couponHandler.GetActiveCoupons)
		
		// Apply coupon (during checkout)
		public.POST("/coupons/apply", couponHandler.ApplyCoupon)
	}

	// Admin routes (require authentication + admin role)
	admin := rg.Group("")
	admin.Use(middleware.JWTAuth(tokenService))
	admin.Use(middleware.RequireAdmin())
	{
		// Get coupon statistics
		admin.GET("/coupons/stats", couponHandler.GetCouponStats)
		
		// Get all coupons
		admin.GET("/coupons", couponHandler.GetCoupons)
		
		// Get coupon by ID
		admin.GET("/coupons/:id", couponHandler.GetCouponByID)
		
		// Create coupon
		admin.POST("/coupons", couponHandler.CreateCoupon)
		
		// Update coupon
		admin.PUT("/coupons/:id", couponHandler.UpdateCoupon)
		
		// Delete coupon
		admin.DELETE("/coupons/:id", couponHandler.DeleteCoupon)
	}

	// User routes (require authentication)
	user := rg.Group("")
	user.Use(middleware.JWTAuth(tokenService))
	{
		// Get user's coupon usages
		user.GET("/coupons/my-usages", couponHandler.GetUserCouponUsages)
	}
}
