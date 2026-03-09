package model

import (
	"time"

	"gorm.io/gorm"
)

// Shop represents a seller's shop
type Shop struct {
	ID                uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID            uint        `gorm:"type:int;not null;uniqueIndex" json:"user_id"`
	User              *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name              string      `gorm:"type:varchar(255);not null" json:"name"`
	Slug              string      `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Description       string      `gorm:"type:text" json:"description"`
	Logo              string      `gorm:"type:varchar(500)" json:"logo"`
	CoverImage        string      `gorm:"type:varchar(500)" json:"cover_image"`
	Phone             string      `gorm:"type:varchar(20)" json:"phone"`
	Email             string      `gorm:"type:varchar(255)" json:"email"`
	Address           string      `gorm:"type:varchar(500)" json:"address"`
	Status            ShopStatus  `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	VerificationStatus string     `gorm:"type:varchar(20);default:'unverified'" json:"verification_status"`
	Rating            float64     `gorm:"type:decimal(3,2);default:0" json:"rating"`
	RatingCount       int         `gorm:"type:int;default:0" json:"rating_count"`
	FollowerCount     int         `gorm:"type:int;default:0" json:"follower_count"`
	ProductCount      int         `gorm:"type:int;default:0" json:"product_count"`
	TotalSales        int64       `gorm:"type:int;default:0" json:"total_sales"`
	TotalRevenue      float64     `gorm:"type:decimal(18,2);default:0" json:"total_revenue"`

	// Relationships
	Products  []Product  `gorm:"foreignKey:ShopID" json:"products,omitempty"`
	Reviews   []Review   `gorm:"foreignKey:ShopID" json:"reviews,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ShopStatus defines the status of a shop
type ShopStatus string

const (
	ShopStatusActive     ShopStatus = "active"
	ShopStatusInactive   ShopStatus = "inactive"
	ShopStatusSuspended  ShopStatus = "suspended"
	ShopStatusPending    ShopStatus = "pending"
)

// Review represents a product or shop review
type Review struct {
	ID        uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"type:int;not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID uint      `gorm:"type:int;index" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ShopID    uint      `gorm:"type:int;index" json:"shop_id"`
	Shop      *Shop     `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	OrderID   uint      `gorm:"type:int;index" json:"order_id"`
	Rating    int       `gorm:"type:int;not null" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`
	Images    string    `gorm:"type:text" json:"images"` // JSON array
	IsApproved bool     `gorm:"type:bit;default:false" json:"is_approved"`
	HelpfulCount int    `gorm:"type:int;default:0" json:"helpful_count"`

	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Address represents a user's address
type Address struct {
	ID        uint   `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID    uint   `gorm:"type:int;not null;index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Phone     string `gorm:"type:varchar(20);not null" json:"phone"`
	Street    string `gorm:"type:varchar(500);not null" json:"street"`
	Ward      string `gorm:"type:varchar(200)" json:"ward"`
	District  string `gorm:"type:varchar(200);not null" json:"district"`
	City      string `gorm:"type:varchar(200);not null" json:"city"`
	Country   string `gorm:"type:varchar(100);default:'Vietnam'" json:"country"`
	IsDefault bool   `gorm:"type:bit;default:false" json:"is_default"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	slug := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			slug += string(r)
		} else if r >= 'A' && r <= 'Z' {
			slug += string(r + 32) // Convert to lowercase
		} else {
			slug += "-"
		}
	}
	return slug
}
