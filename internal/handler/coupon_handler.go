package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CouponHandler handles coupon HTTP requests
type CouponHandler struct {
	couponService *service.CouponService
}

// NewCouponHandler creates a new coupon handler
func NewCouponHandler(couponService *service.CouponService) *CouponHandler {
	return &CouponHandler{couponService: couponService}
}

// CreateCoupon handles POST /api/coupons
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var input service.CreateCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context (admin)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	coupon, err := h.couponService.CreateCoupon(c.Request.Context(), input, userID.(uint))
	if err != nil {
		if errors.Is(err, service.ErrCouponCodeExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) || 
		   errors.Is(err, service.ErrInvalidDateRange) ||
		   errors.Is(err, service.ErrCouponCodeTooShort) ||
		   errors.Is(err, service.ErrDiscountPercentageExceeded) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create coupon",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Coupon created successfully",
		"data":    coupon,
	})
}

// UpdateCoupon handles PUT /api/coupons/:id
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid coupon ID",
		})
		return
	}

	var input service.UpdateCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	coupon, err := h.couponService.UpdateCoupon(c.Request.Context(), uint(id), input)
	if err != nil {
		if errors.Is(err, service.ErrCouponNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update coupon",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Coupon updated successfully",
		"data":    coupon,
	})
}

// DeleteCoupon handles DELETE /api/coupons/:id
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid coupon ID",
		})
		return
	}

	if err := h.couponService.DeleteCoupon(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrCouponNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete coupon",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Coupon deleted successfully",
	})
}

// GetCouponByID handles GET /api/coupons/:id
func (h *CouponHandler) GetCouponByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid coupon ID",
		})
		return
	}

	coupon, err := h.couponService.GetCouponByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Coupon not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    coupon,
	})
}

// GetCoupons handles GET /api/coupons
func (h *CouponHandler) GetCoupons(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	status := c.Query("status")
	isActive := c.Query("is_active")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := model.CouponFilter{
		Page:  page,
		Limit: limit,
	}

	if status != "" {
		s := model.CouponStatus(status)
		filter.Status = &s
	}
	if isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}

	coupons, total, err := h.couponService.GetCoupons(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get coupons",
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
			"coupons": coupons,
			"pagination": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// GetActiveCoupons handles GET /api/coupons/active
func (h *CouponHandler) GetActiveCoupons(c *gin.Context) {
	coupons, err := h.couponService.GetActiveCoupons(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get active coupons",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"coupons": coupons,
			"count":   len(coupons),
		},
	})
}

// ApplyCoupon handles POST /api/coupons/apply
func (h *CouponHandler) ApplyCoupon(c *gin.Context) {
	var input service.ApplyCouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context if available
	userID, exists := c.Get("user_id")
	if exists {
		input.UserID = userID.(uint)
	}

	result, err := h.couponService.ApplyCoupon(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to apply coupon",
			"error":   err.Error(),
		})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": result.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": result.Message,
		"data": gin.H{
			"coupon":          result.Coupon,
			"discount_amount": result.DiscountAmount,
			"original_total":  result.OriginalTotal,
			"final_total":     result.FinalTotal,
		},
	})
}

// GetCouponStats handles GET /api/coupons/stats
func (h *CouponHandler) GetCouponStats(c *gin.Context) {
	stats, err := h.couponService.GetCouponStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get coupon statistics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetUserCouponUsages handles GET /api/coupons/my-usages
func (h *CouponHandler) GetUserCouponUsages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	usages, total, err := h.couponService.GetUserCouponUsages(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get coupon usages",
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
			"usages": usages,
			"pagination": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}
