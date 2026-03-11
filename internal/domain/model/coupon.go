package model

import (
	"time"
)

// DiscountType represents the type of discount
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
	DiscountTypeShipping   DiscountType = "free_shipping"
)

// CouponStatus represents the status of a coupon
type CouponStatus string

const (
	CouponStatusActive   CouponStatus = "active"
	CouponStatusInactive CouponStatus = "inactive"
	CouponStatusExpired  CouponStatus = "expired"
	CouponStatusUsedUp   CouponStatus = "used_up"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	DiscountType    DiscountType   `gorm:"type:varchar(50);not null" json:"discount_type"`
	DiscountValue   float64        `gorm:"type:decimal(18,2);not null" json:"discount_value"`
	MaxDiscount     float64        `gorm:"type:decimal(18,2)" json:"max_discount"` // For percentage discounts
	MinOrderValue   float64        `gorm:"type:decimal(18,2)" json:"min_order_value"`
	MaxOrderValue   float64        `gorm:"type:decimal(18,2)" json:"max_order_value"`
	UsageLimit      int            `gorm:"type:int;default:0" json:"usage_limit"` // 0 = unlimited
	UsedCount       int            `gorm:"type:int;default:0" json:"used_count"`
	UsageLimitPerUser int          `gorm:"type:int;default:1" json:"usage_limit_per_user"`
	StartDate       time.Time      `gorm:"type:datetime" json:"start_date"`
	EndDate         time.Time      `gorm:"type:datetime;index" json:"end_date"`
	IsActive        bool           `gorm:"type:bit;default:true" json:"is_active"`
	Status          CouponStatus   `gorm:"type:varchar(50);default:'active'" json:"status"`
	ApplicableCategories string    `gorm:"type:varchar(max)" json:"applicable_categories"` // JSON array of category IDs
	ApplicableProducts   string    `gorm:"type:varchar(max)" json:"applicable_products"`   // JSON array of product IDs
	ExcludedCategories   string    `gorm:"type:varchar(max)" json:"excluded_categories"`   // JSON array of category IDs
	ExcludedProducts     string    `gorm:"type:varchar(max)" json:"excluded_products"`     // JSON array of product IDs
	UserRestricted       bool      `gorm:"type:bit;default:false" json:"user_restricted"`
	RestrictedUsers      string    `gorm:"type:varchar(max)" json:"restricted_users"` // JSON array of user IDs
	CreatedBy           uint       `gorm:"type:int" json:"created_by"`
	CreatedAt           time.Time  `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"index" json:"-"`
}

// TableName specifies the table name for Coupon
func (Coupon) TableName() string {
	return "coupons"
}

// CouponUsage tracks coupon usage by users
type CouponUsage struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CouponID  uint       `gorm:"type:int;not null;index" json:"coupon_id"`
	UserID    uint       `gorm:"type:int;not null;index" json:"user_id"`
	OrderID   uint       `gorm:"type:int;index" json:"order_id"`
	DiscountAmount float64 `gorm:"type:decimal(18,2);not null" json:"discount_amount"`
	UsedAt    time.Time  `gorm:"type:datetime;not null" json:"used_at"`
	
	// Relationships
	Coupon *Coupon `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order  *Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName specifies the table name for CouponUsage
func (CouponUsage) TableName() string {
	return "coupon_usages"
}

// CouponDiscountResult represents the result of applying a coupon
type CouponDiscountResult struct {
	Success        bool    `json:"success"`
	Message        string  `json:"message"`
	Coupon         *Coupon `json:"coupon,omitempty"`
	DiscountAmount float64 `json:"discount_amount"`
	OriginalTotal  float64 `json:"original_total"`
	FinalTotal     float64 `json:"final_total"`
}

// IsExpired checks if the coupon is expired
func (c *Coupon) IsExpired() bool {
	return time.Now().After(c.EndDate)
}

// IsNotStarted checks if the coupon has not started yet
func (c *Coupon) IsNotStarted() bool {
	return time.Now().Before(c.StartDate)
}

// IsAvailable checks if the coupon is available for use
func (c *Coupon) IsAvailable() bool {
	if !c.IsActive {
		return false
	}
	if c.IsExpired() {
		return false
	}
	if c.IsNotStarted() {
		return false
	}
	if c.UsageLimit > 0 && c.UsedCount >= c.UsageLimit {
		return false
	}
	return true
}

// GetStatus returns the current status of the coupon
func (c *Coupon) GetStatus() CouponStatus {
	if !c.IsActive {
		return CouponStatusInactive
	}
	if c.IsExpired() {
		return CouponStatusExpired
	}
	if c.UsageLimit > 0 && c.UsedCount >= c.UsageLimit {
		return CouponStatusUsedUp
	}
	return CouponStatusActive
}

// CanBeUsedByUser checks if a user can use this coupon
func (c *Coupon) CanBeUsedByUser(userID uint, usageCount int) bool {
	if c.UsageLimitPerUser <= 0 {
		return true // Unlimited per user
	}
	return usageCount < c.UsageLimitPerUser
}

// CalculateDiscount calculates the discount amount for an order
func (c *Coupon) CalculateDiscount(orderTotal float64) float64 {
	var discount float64
	
	switch c.DiscountType {
	case DiscountTypePercentage:
		discount = orderTotal * (c.DiscountValue / 100)
		if c.MaxDiscount > 0 && discount > c.MaxDiscount {
			discount = c.MaxDiscount
		}
	case DiscountTypeFixed:
		discount = c.DiscountValue
		if discount > orderTotal {
			discount = orderTotal // Don't exceed order total
		}
	case DiscountTypeShipping:
		// Free shipping - discount will be calculated based on shipping fee
		discount = 0 // Will be handled separately
	}
	
	// Ensure discount doesn't exceed order total
	if discount > orderTotal {
		discount = orderTotal
	}
	
	return discount
}

// UserCouponUsage represents a user's usage of a coupon
type UserCouponUsage struct {
	UserID      uint `json:"user_id"`
	CouponID    uint `json:"coupon_id"`
	UsageCount  int  `json:"usage_count"`
}

// CouponStats represents coupon statistics
type CouponStats struct {
	TotalCoupons     int64   `json:"total_coupons"`
	ActiveCoupons    int64   `json:"active_coupons"`
	ExpiredCoupons   int64   `json:"expired_coupons"`
	TotalUsage       int64   `json:"total_usage"`
	TotalDiscount    float64 `json:"total_discount"`
	AverageDiscount  float64 `json:"average_discount"`
}

// CouponFilter represents filters for coupon search
type CouponFilter struct {
	Status     *CouponStatus `json:"status"`
	DiscountType *DiscountType `json:"discount_type"`
	IsActive   *bool         `json:"is_active"`
	Search     string        `json:"search"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}
