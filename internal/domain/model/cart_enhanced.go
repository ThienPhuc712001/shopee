package model

import (
	"time"

	"gorm.io/gorm"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID           uint          `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID       uint          `gorm:"type:int;uniqueIndex;not null" json:"user_id"`
	User         *User         `gorm:"foreignKey:UserID" json:"-"`
	TotalItems   int           `gorm:"type:int;not null;default:0" json:"total_items"`
	Subtotal     float64       `gorm:"type:decimal(18,2);not null;default:0" json:"subtotal"`
	Discount     float64       `gorm:"type:decimal(18,2);default:0" json:"discount"`
	Total        float64       `gorm:"type:decimal(18,2);not null;default:0" json:"total"`
	Currency     string        `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	LastActivity time.Time     `gorm:"type:datetime;not null;default:GETDATE()" json:"last_activity"`

	// Relationships
	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID            uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	CartID        uint        `gorm:"type:int;not null;index" json:"cart_id"`
	Cart          *Cart       `gorm:"foreignKey:CartID" json:"-"`
	ProductID     uint        `gorm:"type:int;not null;index" json:"product_id"`
	Product       *Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	VariantID     *uint       `gorm:"type:int;index" json:"variant_id"`
	Variant       *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Quantity      int         `gorm:"type:int;not null;default:1" json:"quantity"`
	Price         float64     `gorm:"type:decimal(18,2);not null" json:"price"`
	OriginalPrice float64     `gorm:"type:decimal(18,2)" json:"original_price"`
	Discount      float64     `gorm:"type:decimal(18,2);default:0" json:"discount"`
	Subtotal      float64     `gorm:"type:decimal(18,2);not null" json:"subtotal"`
	ProductName   string      `gorm:"type:varchar(500)" json:"product_name"`
	ProductImage  string      `gorm:"type:varchar(500)" json:"product_image"`
	ShopID        uint        `gorm:"type:int;index" json:"shop_id"`
	Shop          *Shop       `gorm:"foreignKey:ShopID" json:"shop,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook for CartItem - calculates subtotal
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	ci.Subtotal = ci.Price * float64(ci.Quantity)
	return nil
}

// BeforeSave hook for CartItem - recalculates subtotal
func (ci *CartItem) BeforeSave(tx *gorm.DB) error {
	ci.Subtotal = ci.Price * float64(ci.Quantity)
	return nil
}

// BeforeCreate hook for Cart - initializes timestamps
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.LastActivity = time.Now()
	return nil
}

// GetTotalPrice returns the total price of the cart
func (c *Cart) GetTotalPrice() float64 {
	return c.Subtotal - c.Discount
}

// GetTotalItems returns the total number of items in the cart
func (c *Cart) GetTotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// RecalculateTotals recalculates cart totals based on items
func (c *Cart) RecalculateTotals() {
	var subtotal float64
	var totalItems int

	for _, item := range c.Items {
		if item.DeletedAt.Valid {
			continue // Skip soft-deleted items
		}
		subtotal += item.Subtotal
		totalItems += item.Quantity
	}

	c.Subtotal = subtotal
	c.Total = subtotal - c.Discount
	c.TotalItems = totalItems
	c.LastActivity = time.Now()
}

// GetTotalPrice returns the subtotal for a cart item
func (ci *CartItem) GetTotalPrice() float64 {
	return ci.Subtotal
}

// GetDiscountedPrice returns the price after discount
func (ci *CartItem) GetDiscountedPrice() float64 {
	if ci.Discount > 0 {
		return ci.Price * (1 - ci.Discount/100)
	}
	return ci.Price
}

// IsAvailable checks if the cart item product is available
func (ci *CartItem) IsAvailable() bool {
	if ci.Product == nil {
		return false
	}
	return ci.Product.Status == ProductStatusActive && ci.Product.Stock >= ci.Quantity
}

// CartWithItems represents a cart with loaded items
type CartWithItems struct {
	Cart
	Items []CartItemWithProduct `json:"items"`
}

// CartItemWithProduct represents a cart item with loaded product
type CartItemWithProduct struct {
	CartItem
	Product *Product `json:"product"`
	Variant *ProductVariant `json:"variant"`
}

// CartSummary represents a summary of the cart
type CartSummary struct {
	TotalItems int     `json:"total_items"`
	Subtotal   float64 `json:"subtotal"`
	Discount   float64 `json:"discount"`
	Total      float64 `json:"total"`
	Currency   string  `json:"currency"`
}

// CartItemInput represents input for adding/updating cart item
type CartItemInput struct {
	ProductID uint  `json:"product_id" binding:"required"`
	VariantID *uint `json:"variant_id"`
	Quantity  int   `json:"quantity" binding:"required,min=1,max=999"`
}

// CartUpdateInput represents input for updating cart item
type CartUpdateInput struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=999"`
}

// CartCheckoutItem represents a cart item ready for checkout
type CartCheckoutItem struct {
	ProductID   uint    `json:"product_id"`
	VariantID   *uint   `json:"variant_id"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
	ShopID      uint    `json:"shop_id"`
	ShopName    string  `json:"shop_name"`
	ProductName string  `json:"product_name"`
	ProductImage string `json:"product_image"`
	IsAvailable bool    `json:"is_available"`
	StockStatus string  `json:"stock_status"` // in_stock, low_stock, out_of_stock
}

// CartCheckoutSummary represents cart summary for checkout
type CartCheckoutSummary struct {
	Items       []CartCheckoutItem `json:"items"`
	Subtotal    float64            `json:"subtotal"`
	ShippingFee float64            `json:"shipping_fee"`
	Discount    float64            `json:"discount"`
	Total       float64            `json:"total"`
	Currency    string             `json:"currency"`
}

// IsValid checks if cart input is valid
func (input *CartItemInput) IsValid() bool {
	return input.ProductID > 0 && input.Quantity > 0 && input.Quantity <= 999
}

// GetStockStatus returns the stock status for a product
func GetStockStatus(available, requested int) string {
	if available <= 0 {
		return "out_of_stock"
	}
	if available < requested {
		return "low_stock"
	}
	return "in_stock"
}
