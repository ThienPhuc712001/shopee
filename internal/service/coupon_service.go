package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

// CouponService handles coupon business logic
type CouponService struct {
	repo *repository.CouponRepository
}

// NewCouponService creates a new coupon service
func NewCouponService(repo *repository.CouponRepository) *CouponService {
	return &CouponService{repo: repo}
}

// CreateCouponInput represents input for creating a coupon
type CreateCouponInput struct {
	Code              string  `json:"code" binding:"required,uppercase"`
	Name              string  `json:"name" binding:"required,min=2,max=255"`
	Description       string  `json:"description"`
	DiscountType      string  `json:"discount_type" binding:"required,oneof=percentage fixed free_shipping"`
	DiscountValue     float64 `json:"discount_value" binding:"required,gt=0"`
	MaxDiscount       float64 `json:"max_discount"`
	MinOrderValue     float64 `json:"min_order_value"`
	MaxOrderValue     float64 `json:"max_order_value"`
	UsageLimit        int     `json:"usage_limit"`
	UsageLimitPerUser int     `json:"usage_limit_per_user"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date" binding:"required"`
	IsActive          bool    `json:"is_active"`
	ApplicableCategories string `json:"applicable_categories"`
	ApplicableProducts   string `json:"applicable_products"`
	ExcludedCategories   string `json:"excluded_categories"`
	ExcludedProducts     string `json:"excluded_products"`
	UserRestricted    bool    `json:"user_restricted"`
	RestrictedUsers   string  `json:"restricted_users"`
}

// UpdateCouponInput represents input for updating a coupon
type UpdateCouponInput struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	DiscountValue     float64 `json:"discount_value"`
	MaxDiscount       float64 `json:"max_discount"`
	MinOrderValue     float64 `json:"min_order_value"`
	MaxOrderValue     float64 `json:"max_order_value"`
	UsageLimit        int     `json:"usage_limit"`
	UsageLimitPerUser int     `json:"usage_limit_per_user"`
	EndDate           string  `json:"end_date"`
	IsActive          *bool   `json:"is_active"`
}

// ApplyCouponInput represents input for applying a coupon
type ApplyCouponInput struct {
	Code        string  `json:"code" binding:"required"`
	OrderTotal  float64 `json:"order_total" binding:"required,gt=0"`
	UserID      uint    `json:"user_id"`
	ProductIDs  []uint  `json:"product_ids"`
	CategoryIDs []uint  `json:"category_ids"`
}

// CreateCoupon creates a new coupon
func (s *CouponService) CreateCoupon(ctx context.Context, input CreateCouponInput, createdBy uint) (*model.Coupon, error) {
	// Validate input
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Check if code already exists
	_, err := s.repo.GetCouponByCode(ctx, input.Code)
	if err == nil {
		return nil, ErrCouponCodeExists
	}

	// Parse dates
	var startDate, endDate time.Time
	if input.StartDate != "" {
		startDate, err = time.Parse(time.RFC3339, input.StartDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
	} else {
		startDate = time.Now()
	}

	endDate, err = time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	// Check date range
	if endDate.Before(startDate) {
		return nil, ErrInvalidDateRange
	}

	// Create coupon
	coupon := &model.Coupon{
		Code:                 input.Code,
		Name:                 input.Name,
		Description:          input.Description,
		DiscountType:         model.DiscountType(input.DiscountType),
		DiscountValue:        input.DiscountValue,
		MaxDiscount:          input.MaxDiscount,
		MinOrderValue:        input.MinOrderValue,
		MaxOrderValue:        input.MaxOrderValue,
		UsageLimit:           input.UsageLimit,
		UsageLimitPerUser:    input.UsageLimitPerUser,
		StartDate:            startDate,
		EndDate:              endDate,
		IsActive:             input.IsActive,
		ApplicableCategories: input.ApplicableCategories,
		ApplicableProducts:   input.ApplicableProducts,
		ExcludedCategories:   input.ExcludedCategories,
		ExcludedProducts:     input.ExcludedProducts,
		UserRestricted:       input.UserRestricted,
		RestrictedUsers:      input.RestrictedUsers,
		CreatedBy:            createdBy,
		Status:               model.CouponStatusActive,
	}

	if err := s.repo.CreateCoupon(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

// UpdateCoupon updates an existing coupon
func (s *CouponService) UpdateCoupon(ctx context.Context, id uint, input UpdateCouponInput) (*model.Coupon, error) {
	// Get existing coupon
	coupon, err := s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, ErrCouponNotFound
	}

	// Update fields
	if input.Name != "" {
		coupon.Name = input.Name
	}
	if input.Description != "" {
		coupon.Description = input.Description
	}
	if input.DiscountValue > 0 {
		coupon.DiscountValue = input.DiscountValue
	}
	if input.MaxDiscount > 0 {
		coupon.MaxDiscount = input.MaxDiscount
	}
	if input.MinOrderValue > 0 {
		coupon.MinOrderValue = input.MinOrderValue
	}
	if input.MaxOrderValue > 0 {
		coupon.MaxOrderValue = input.MaxOrderValue
	}
	if input.UsageLimit >= 0 {
		coupon.UsageLimit = input.UsageLimit
	}
	if input.UsageLimitPerUser > 0 {
		coupon.UsageLimitPerUser = input.UsageLimitPerUser
	}
	if input.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, input.EndDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		coupon.EndDate = endDate
	}
	if input.IsActive != nil {
		coupon.IsActive = *input.IsActive
	}

	// Update status
	coupon.Status = coupon.GetStatus()

	if err := s.repo.UpdateCoupon(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

// DeleteCoupon deletes a coupon
func (s *CouponService) DeleteCoupon(ctx context.Context, id uint) error {
	_, err := s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return ErrCouponNotFound
	}

	return s.repo.DeleteCoupon(ctx, id)
}

// GetCouponByID retrieves a coupon by ID
func (s *CouponService) GetCouponByID(ctx context.Context, id uint) (*model.Coupon, error) {
	coupon, err := s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, ErrCouponNotFound
	}
	return coupon, nil
}

// GetCouponByCode retrieves a coupon by code
func (s *CouponService) GetCouponByCode(ctx context.Context, code string) (*model.Coupon, error) {
	coupon, err := s.repo.GetCouponByCode(ctx, code)
	if err != nil {
		return nil, ErrCouponNotFound
	}
	return coupon, nil
}

// GetCoupons retrieves coupons with filters
func (s *CouponService) GetCoupons(ctx context.Context, filter model.CouponFilter) ([]model.Coupon, int64, error) {
	return s.repo.GetCoupons(ctx, filter)
}

// GetActiveCoupons retrieves all active coupons
func (s *CouponService) GetActiveCoupons(ctx context.Context) ([]model.Coupon, error) {
	return s.repo.GetActiveCoupons(ctx)
}

// ApplyCoupon validates and applies a coupon to an order
func (s *CouponService) ApplyCoupon(ctx context.Context, input ApplyCouponInput) (*model.CouponDiscountResult, error) {
	result := &model.CouponDiscountResult{
		OriginalTotal: input.OrderTotal,
		Success:       false,
	}

	// Get coupon by code
	coupon, err := s.repo.GetCouponByCode(ctx, input.Code)
	if err != nil {
		result.Message = "Invalid coupon code"
		return result, nil
	}

	result.Coupon = coupon

	// Validate coupon
	if err := s.ValidateCoupon(ctx, coupon, input.OrderTotal, input.UserID, input.ProductIDs, input.CategoryIDs); err != nil {
		result.Message = err.Error()
		return result, nil
	}

	// Calculate discount
	discount := coupon.CalculateDiscount(input.OrderTotal)
	result.DiscountAmount = discount
	result.FinalTotal = input.OrderTotal - discount
	result.Success = true
	result.Message = "Coupon applied successfully"

	return result, nil
}

// ValidateCoupon validates if a coupon can be applied
func (s *CouponService) ValidateCoupon(ctx context.Context, coupon *model.Coupon, orderTotal float64, userID uint, productIDs, categoryIDs []uint) error {
	// Check if coupon exists
	if coupon == nil {
		return ErrCouponNotFound
	}

	// Check if coupon is active
	if !coupon.IsActive {
		return ErrCouponInactive
	}

	// Check if coupon is expired
	if coupon.IsExpired() {
		return ErrCouponExpired
	}

	// Check if coupon has started
	if coupon.IsNotStarted() {
		return ErrCouponNotStarted
	}

	// Check usage limit
	if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
		return ErrCouponUsageLimitReached
	}

	// Check minimum order value
	if coupon.MinOrderValue > 0 && orderTotal < coupon.MinOrderValue {
		return fmt.Errorf("minimum order value is $%.2f", coupon.MinOrderValue)
	}

	// Check maximum order value
	if coupon.MaxOrderValue > 0 && orderTotal > coupon.MaxOrderValue {
		return fmt.Errorf("maximum order value is $%.2f", coupon.MaxOrderValue)
	}

	// Check user usage limit
	if userID > 0 && coupon.UsageLimitPerUser > 0 {
		usageCount, err := s.repo.GetUserCouponUsageCount(ctx, coupon.ID, userID)
		if err != nil {
			return err
		}
		if usageCount >= coupon.UsageLimitPerUser {
			return ErrCouponAlreadyUsed
		}
	}

	// Check user restriction
	if coupon.UserRestricted && coupon.RestrictedUsers != "" {
		// In production, parse JSON and check if userID is in the list
		// For now, simplified check
		if userID == 0 {
			return ErrCouponRestricted
		}
	}

	// Check applicable products (if specified)
	if coupon.ApplicableProducts != "" && len(productIDs) > 0 {
		// In production, parse JSON and check if any product is applicable
		// Simplified for now
	}

	// Check excluded products
	if coupon.ExcludedProducts != "" && len(productIDs) > 0 {
		// In production, parse JSON and check if any product is excluded
		// Simplified for now
	}

	return nil
}

// CalculateDiscount calculates the discount for a coupon
func (s *CouponService) CalculateDiscount(coupon *model.Coupon, orderTotal float64) float64 {
	return coupon.CalculateDiscount(orderTotal)
}

// UseCoupon marks a coupon as used by a user
func (s *CouponService) UseCoupon(ctx context.Context, couponID, userID, orderID uint, discountAmount float64) error {
	return s.repo.IncrementUsageCount(ctx, couponID)
}

// CreateCouponUsage creates a coupon usage record
func (s *CouponService) CreateCouponUsage(ctx context.Context, couponID, userID, orderID uint, discountAmount float64) error {
	usage := &model.CouponUsage{
		CouponID:       couponID,
		UserID:         userID,
		OrderID:        orderID,
		DiscountAmount: discountAmount,
		UsedAt:         time.Now(),
	}

	return s.repo.CreateCouponUsage(ctx, usage)
}

// GetUserCouponUsages retrieves coupon usages for a user
func (s *CouponService) GetUserCouponUsages(ctx context.Context, userID uint, page, limit int) ([]model.CouponUsage, int64, error) {
	return s.repo.GetCouponUsagesByUser(ctx, userID, page, limit)
}

// GetCouponStats retrieves coupon statistics
func (s *CouponService) GetCouponStats(ctx context.Context) (*model.CouponStats, error) {
	return s.repo.GetCouponStats(ctx)
}

// GetExpiringCoupons retrieves coupons expiring soon
func (s *CouponService) GetExpiringCoupons(ctx context.Context, days int) ([]model.Coupon, error) {
	return s.repo.GetExpiringCoupons(ctx, days)
}

// DeactivateExpiredCoupons deactivates all expired coupons
func (s *CouponService) DeactivateExpiredCoupons(ctx context.Context) (int64, error) {
	return s.repo.DeactivateExpiredCoupons(ctx)
}

// GenerateCouponCode generates a unique coupon code
func (s *CouponService) GenerateCouponCode(prefix string, length int) string {
	// Simple implementation - in production, use a more secure method
	code := prefix
	for len(code) < length {
		code += string(rune('A' + (time.Now().Nanosecond() % 26)))
	}
	return code[:length]
}

// validateCreateInput validates coupon creation input
func (s *CouponService) validateCreateInput(input CreateCouponInput) error {
	if input.Code == "" {
		return ErrCouponCodeRequired
	}
	if len(input.Code) < 3 {
		return ErrCouponCodeTooShort
	}
	if input.Name == "" {
		return ErrCouponNameRequired
	}
	if input.DiscountValue <= 0 {
		return ErrInvalidDiscountValue
	}
	if input.DiscountType == string(model.DiscountTypePercentage) {
		if input.DiscountValue > 100 {
			return ErrDiscountPercentageExceeded
		}
	}
	if input.EndDate == "" {
		return ErrEndDateRequired
	}
	return nil
}

// Error definitions
var (
	ErrCouponNotFound            = errors.New("coupon not found")
	ErrCouponCodeExists          = errors.New("coupon code already exists")
	ErrCouponInactive            = errors.New("coupon is inactive")
	ErrCouponExpired             = errors.New("coupon has expired")
	ErrCouponNotStarted          = errors.New("coupon has not started yet")
	ErrCouponUsageLimitReached   = errors.New("coupon usage limit reached")
	ErrCouponAlreadyUsed         = errors.New("coupon already used by this user")
	ErrCouponRestricted          = errors.New("coupon is restricted to specific users")
	ErrInvalidDiscountValue      = errors.New("invalid discount value")
	ErrDiscountPercentageExceeded = errors.New("discount percentage cannot exceed 100")
	ErrInvalidDateRange          = errors.New("end date must be after start date")
	ErrInvalidDateFormat         = errors.New("invalid date format, use RFC3339")
	ErrCouponCodeRequired        = errors.New("coupon code is required")
	ErrCouponCodeTooShort        = errors.New("coupon code must be at least 3 characters")
	ErrCouponNameRequired        = errors.New("coupon name is required")
	ErrEndDateRequired           = errors.New("end date is required")
	ErrMinOrderValueNotMet       = errors.New("minimum order value not met")
	ErrMaxOrderValueExceeded     = errors.New("maximum order value exceeded")
)
