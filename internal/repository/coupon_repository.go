package repository

import (
	"ecommerce/internal/domain/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CouponRepository handles database operations for coupons
type CouponRepository struct {
	db *gorm.DB
}

// NewCouponRepository creates a new coupon repository
func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

// CreateCoupon creates a new coupon
func (r *CouponRepository) CreateCoupon(ctx context.Context, coupon *model.Coupon) error {
	return r.db.WithContext(ctx).Create(coupon).Error
}

// UpdateCoupon updates an existing coupon
func (r *CouponRepository) UpdateCoupon(ctx context.Context, coupon *model.Coupon) error {
	return r.db.WithContext(ctx).Save(coupon).Error
}

// DeleteCoupon deletes a coupon (soft delete)
func (r *CouponRepository) DeleteCoupon(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Coupon{}, id).Error
}

// GetCouponByID retrieves a coupon by ID
func (r *CouponRepository) GetCouponByID(ctx context.Context, id uint) (*model.Coupon, error) {
	var coupon model.Coupon
	err := r.db.WithContext(ctx).First(&coupon, id).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetCouponByCode retrieves a coupon by code
func (r *CouponRepository) GetCouponByCode(ctx context.Context, code string) (*model.Coupon, error) {
	var coupon model.Coupon
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetCoupons retrieves coupons with filters and pagination
func (r *CouponRepository) GetCoupons(ctx context.Context, filter model.CouponFilter) ([]model.Coupon, int64, error) {
	var coupons []model.Coupon
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Coupon{})

	// Apply filters
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.DiscountType != nil {
		query = query.Where("discount_type = ?", *filter.DiscountType)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("code LIKE ? OR name LIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Get coupons
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&coupons).Error

	return coupons, total, err
}

// GetActiveCoupons retrieves all active coupons
func (r *CouponRepository) GetActiveCoupons(ctx context.Context) ([]model.Coupon, error) {
	var coupons []model.Coupon
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND start_date <= ? AND end_date >= ?", true, now, now).
		Where("usage_limit = 0 OR used_count < usage_limit").
		Order("created_at DESC").
		Find(&coupons).Error
	return coupons, err
}

// IncrementUsageCount increments the usage count of a coupon
func (r *CouponRepository) IncrementUsageCount(ctx context.Context, couponID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var coupon model.Coupon
		
		// Lock the row for update
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			First(&coupon, couponID).Error; err != nil {
			return err
		}

		// Check if usage limit is reached
		if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
			return errors.New("coupon usage limit reached")
		}

		// Increment usage count
		coupon.UsedCount++
		
		// Update status if needed
		if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
			coupon.Status = model.CouponStatusUsedUp
		}

		return tx.Save(&coupon).Error
	})
}

// GetUserCouponUsageCount gets how many times a user has used a coupon
func (r *CouponRepository) GetUserCouponUsageCount(ctx context.Context, couponID, userID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.CouponUsage{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).
		Count(&count).Error
	return int(count), err
}

// CreateCouponUsage creates a coupon usage record
func (r *CouponRepository) CreateCouponUsage(ctx context.Context, usage *model.CouponUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

// GetCouponUsageByOrder retrieves coupon usage by order ID
func (r *CouponRepository) GetCouponUsageByOrder(ctx context.Context, orderID uint) (*model.CouponUsage, error) {
	var usage model.CouponUsage
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Preload("Coupon").
		First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

// GetCouponUsagesByUser retrieves all coupon usages by a user
func (r *CouponRepository) GetCouponUsagesByUser(ctx context.Context, userID uint, page, limit int) ([]model.CouponUsage, int64, error) {
	var usages []model.CouponUsage
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&model.CouponUsage{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Get usages
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Coupon").
		Preload("Order").
		Order("used_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&usages).Error

	return usages, total, err
}

// GetCouponStats retrieves coupon statistics
func (r *CouponRepository) GetCouponStats(ctx context.Context) (*model.CouponStats, error) {
	stats := &model.CouponStats{}

	// Total coupons
	if err := r.db.WithContext(ctx).
		Model(&model.Coupon{}).
		Count(&stats.TotalCoupons).Error; err != nil {
		return nil, err
	}

	// Active coupons
	if err := r.db.WithContext(ctx).
		Model(&model.Coupon{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveCoupons).Error; err != nil {
		return nil, err
	}

	// Expired coupons
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&model.Coupon{}).
		Where("end_date < ?", now).
		Count(&stats.ExpiredCoupons).Error; err != nil {
		return nil, err
	}

	// Total usage and discount
	type UsageStats struct {
		TotalUsage    int64
		TotalDiscount float64
	}
	var usageStats UsageStats
	if err := r.db.WithContext(ctx).
		Model(&model.CouponUsage{}).
		Select("COUNT(*) as total_usage, SUM(discount_amount) as total_discount").
		Scan(&usageStats).Error; err != nil {
		return nil, err
	}

	stats.TotalUsage = usageStats.TotalUsage
	stats.TotalDiscount = usageStats.TotalDiscount

	if stats.TotalUsage > 0 {
		stats.AverageDiscount = stats.TotalDiscount / float64(stats.TotalUsage)
	}

	return stats, nil
}

// GetExpiringCoupons retrieves coupons that are about to expire
func (r *CouponRepository) GetExpiringCoupons(ctx context.Context, days int) ([]model.Coupon, error) {
	var coupons []model.Coupon
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Where("is_active = ? AND end_date BETWEEN ? AND ?", now, futureDate).
		Where("usage_limit = 0 OR used_count < usage_limit").
		Find(&coupons).Error

	return coupons, err
}

// UpdateCouponStatus updates the status of a coupon
func (r *CouponRepository) UpdateCouponStatus(ctx context.Context, id uint, status model.CouponStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Coupon{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// DeactivateExpiredCoupons deactivates all expired coupons
func (r *CouponRepository) DeactivateExpiredCoupons(ctx context.Context) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.Coupon{}).
		Where("end_date < ? AND is_active = ?", now, true).
		Update("status", model.CouponStatusExpired).
		Update("is_active", false)

	return result.RowsAffected, result.Error
}

// CheckCouponAvailability checks if a coupon is still available
func (r *CouponRepository) CheckCouponAvailability(ctx context.Context, code string) (bool, error) {
	var coupon model.Coupon
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&coupon).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return coupon.IsAvailable(), nil
}

// GetCouponsByProduct retrieves coupons applicable to a product
func (r *CouponRepository) GetCouponsByProduct(ctx context.Context, productID, categoryID uint) ([]model.Coupon, error) {
	var coupons []model.Coupon
	now := time.Now()

	// This is a simplified query - in production, you'd parse the JSON arrays
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND start_date <= ? AND end_date >= ?", true, now, now).
		Where("usage_limit = 0 OR used_count < usage_limit").
		Where(
			"(applicable_products = '' OR applicable_products LIKE ?)",
			`%"`+string(rune(productID))+"\"%",
		).
		Find(&coupons).Error

	return coupons, err
}
