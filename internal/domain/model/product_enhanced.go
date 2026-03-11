package model

import (
	"time"

	"gorm.io/gorm"
)

// ProductStatus defines the status of a product
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPending   ProductStatus = "pending"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusBanned    ProductStatus = "banned"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
)

// Category represents a product category
type Category struct {
	ID          uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ParentID    *uint       `gorm:"type:int;index" json:"parent_id"`
	Parent      *Category   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Name        string      `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string      `gorm:"type:text" json:"description"`
	IconURL     string      `gorm:"type:varchar(500)" json:"icon_url"`
	ImageURL    string      `gorm:"type:varchar(500)" json:"image_url"`
	Level       int         `gorm:"type:int;default:0" json:"level"`
	SortOrder   int         `gorm:"type:int;default:0" json:"sort_order"`
	IsActive    bool        `gorm:"type:bit;default:true" json:"is_active"`
	Attributes  string      `gorm:"type:text" json:"attributes"` // JSON schema for category attributes

	// Relationships
	Products []Product `gorm:"foreignKey:CategoryID" json:"-"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CategoryWithCount includes product count for display
type CategoryWithCount struct {
	Category
	ProductCount int64 `json:"product_count"`
}

// CategoryTree represents a category with its full hierarchy
type CategoryTree struct {
	ID           uint           `json:"id"`
	ParentID     *uint          `json:"parent_id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Description  string         `json:"description"`
	IconURL      string         `json:"icon_url"`
	ImageURL     string         `json:"image_url"`
	Level        int            `json:"level"`
	SortOrder    int            `json:"sort_order"`
	IsActive     bool           `json:"is_active"`
	ProductCount int64          `json:"product_count"`
	Children     []CategoryTree `json:"children"`
}

// CategoryBreadcrumb represents a category in a breadcrumb trail
type CategoryBreadcrumb struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Product represents a product in the system
type Product struct {
	ID                uint          `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ShopID            uint          `gorm:"type:int;not null;index" json:"shop_id"`
	Shop              *Shop         `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	CategoryID        uint          `gorm:"type:int;not null;index" json:"category_id"`
	Category          *Category     `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name              string        `gorm:"type:varchar(500);not null" json:"name"`
	Slug              string        `gorm:"type:varchar(500);not null;index" json:"slug"`
	Description       string        `gorm:"type:text" json:"description"`
	ShortDescription  string        `gorm:"type:varchar(1000)" json:"short_description"`
	SKU               string        `gorm:"type:varchar(100);index" json:"sku"`
	Brand             string        `gorm:"type:varchar(200)" json:"brand"`

	// Pricing
	Price             float64       `gorm:"type:decimal(18,2);not null" json:"price"`
	OriginalPrice     float64       `gorm:"type:decimal(18,2)" json:"original_price"`
	DiscountPercent   int           `gorm:"type:int;default:0" json:"discount_percent"`
	Cost              float64       `gorm:"type:decimal(18,2)" json:"cost"` // For profit calculation

	// Inventory
	Stock             int           `gorm:"type:int;not null;default:0" json:"stock"`
	ReservedStock     int           `gorm:"type:int;default:0" json:"reserved_stock"`
	AvailableStock    int           `gorm:"type:int" json:"available_stock"` // Computed field

	// Sales & Analytics
	SoldCount         int64         `gorm:"type:int;default:0;index" json:"sold_count"`
	ViewCount         int64         `gorm:"type:int;default:0" json:"view_count"`
	RatingAvg         float64       `gorm:"type:decimal(3,2);default:0" json:"rating_avg"`
	RatingCount       int           `gorm:"type:int;default:0" json:"rating_count"`
	ReviewCount       int           `gorm:"type:int;default:0" json:"review_count"`

	// Status & Visibility
	Status            ProductStatus `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"`
	IsFeatured        bool          `gorm:"type:bit;default:false" json:"is_featured"`
	IsFlashSale       bool          `gorm:"type:bit;default:false" json:"is_flash_sale"`

	// Flash Sale
	FlashSalePrice    float64       `gorm:"type:decimal(18,2)" json:"flash_sale_price"`
	FlashSaleStart    *time.Time    `gorm:"type:datetime" json:"flash_sale_start"`
	FlashSaleEnd      *time.Time    `gorm:"type:datetime" json:"flash_sale_end"`

	// Physical Properties
	Weight            float64       `gorm:"type:decimal(10,2)" json:"weight"` // in grams
	Dimensions        string        `gorm:"type:varchar(50)" json:"dimensions"` // LxWxH
	WarrantyPeriod    string        `gorm:"type:varchar(50)" json:"warranty_period"`
	ReturnDays        int           `gorm:"type:int;default:7" json:"return_days"`

	// SEO
	Tags              string        `gorm:"type:text" json:"tags"` // JSON array
	MetaTitle         string        `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDescription   string        `gorm:"type:varchar(500)" json:"meta_description"`

	// Relationships
	Images     []ProductImage     `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Variants   []ProductVariant   `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Attributes []ProductAttribute `gorm:"foreignKey:ProductID" json:"attributes,omitempty"`
	Reviews    []Review           `gorm:"foreignKey:ProductID" json:"-"`
	CartItems  []CartItem         `gorm:"foreignKey:ProductID" json:"-"`
	OrderItems []OrderItem        `gorm:"foreignKey:ProductID" json:"-"`

	CreatedAt time.Time      `gorm:"type:datetime;not null;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProductImage represents an image for a product
type ProductImage struct {
	ID         uint   `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ProductID  uint   `gorm:"type:int;not null;index" json:"product_id"`
	Product    *Product `gorm:"foreignKey:ProductID" json:"-"`
	URL        string `gorm:"type:varchar(500);not null" json:"url"`
	AltText    string `gorm:"type:varchar(255)" json:"alt_text"`
	IsPrimary  bool   `gorm:"type:bit;default:false" json:"is_primary"`
	SortOrder  int    `gorm:"type:int;default:0" json:"sort_order"`
	Width      int    `gorm:"type:int" json:"width"`
	Height     int    `gorm:"type:int" json:"height"`
	SizeBytes  int64  `gorm:"type:bigint" json:"size_bytes"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// ProductVariant represents a variant of a product (e.g., size, color)
type ProductVariant struct {
	ID            uint            `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ProductID     uint            `gorm:"type:int;not null;index" json:"product_id"`
	Product       *Product        `gorm:"foreignKey:ProductID" json:"-"`
	SKU           string          `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	Name          string          `gorm:"type:varchar(200)" json:"name"` // e.g., "Red / XL"
	Price         float64         `gorm:"type:decimal(18,2)" json:"price"`
	OriginalPrice float64         `gorm:"type:decimal(18,2)" json:"original_price"`
	Stock         int             `gorm:"type:int;not null;default:0" json:"stock"`
	ReservedStock int             `gorm:"type:int;default:0" json:"reserved_stock"`
	Attributes    string          `gorm:"type:text;not null" json:"attributes"` // JSON: {"color": "Red", "size": "XL"}
	ImageURL      string          `gorm:"type:varchar(500)" json:"image_url"`
	SortOrder     int             `gorm:"type:int;default:0" json:"sort_order"`

	// Relationships
	Inventory *ProductInventory `gorm:"foreignKey:VariantID" json:"inventory,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProductAttribute represents an attribute definition for a product
type ProductAttribute struct {
	ID           uint   `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ProductID    uint   `gorm:"type:int;not null;index" json:"product_id"`
	Product      *Product `gorm:"foreignKey:ProductID" json:"-"`
	CategoryID   *uint  `gorm:"type:int;index" json:"category_id"`
	Name         string `gorm:"type:varchar(100);not null" json:"name"` // e.g., "Color", "Size"
	Type         string `gorm:"type:varchar(50);not null" json:"type"` // text, select, color, number
	Values       string `gorm:"type:text" json:"values"` // JSON array of possible values
	IsFilterable bool   `gorm:"type:bit;default:true" json:"is_filterable"`
	IsVisible    bool   `gorm:"type:bit;default:true" json:"is_visible"`
	SortOrder    int    `gorm:"type:int;default:0" json:"sort_order"`

	// Relationships
	AttributeValues []ProductAttributeValue `gorm:"foreignKey:AttributeID" json:"attribute_values,omitempty"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// ProductAttributeValue represents the value of an attribute for a product
type ProductAttributeValue struct {
	ID          uint   `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	AttributeID uint   `gorm:"type:int;not null;index" json:"attribute_id"`
	Attribute   *ProductAttribute `gorm:"foreignKey:AttributeID" json:"-"`
	ProductID   uint   `gorm:"type:int;not null;index" json:"product_id"`
	Product     *Product `gorm:"foreignKey:ProductID" json:"-"`
	Value       string `gorm:"type:varchar(255);not null" json:"value"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

// ProductInventory represents inventory for a product variant
type ProductInventory struct {
	ID           uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ProductID    uint      `gorm:"type:int;not null;index" json:"product_id"`
	Product      *Product  `gorm:"foreignKey:ProductID" json:"-"`
	VariantID    *uint     `gorm:"type:int;uniqueIndex" json:"variant_id"`
	Variant      *ProductVariant `gorm:"foreignKey:VariantID" json:"-"`
	WarehouseID  *uint     `gorm:"type:int" json:"warehouse_id"`
	Quantity     int       `gorm:"type:int;not null;default:0" json:"quantity"`
	Reserved     int       `gorm:"type:int;not null;default:0" json:"reserved"`
	Available    int       `gorm:"type:int" json:"available"` // Computed: quantity - reserved
	ReorderPoint int       `gorm:"type:int;default:10" json:"reorder_point"`
	LastCountedAt *time.Time `gorm:"type:datetime" json:"last_counted_at"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// ProductTag represents a tag for products
type ProductTag struct {
	ID        uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

// ProductVariantOption represents an option for product variants
type ProductVariantOption struct {
	ID        uint   `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ProductID uint   `gorm:"type:int;not null;index" json:"product_id"`
	Name      string `gorm:"type:varchar(50);not null" json:"name"` // e.g., "Color", "Size"
	Values    string `gorm:"type:text;not null" json:"values"` // JSON array: ["Red", "Blue", "Green"]

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// BeforeCreate hook for Product
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		p.Slug = GenerateSlug(p.Name)
	}
	p.AvailableStock = p.Stock - p.ReservedStock
	return nil
}

// BeforeSave hook for Product
func (p *Product) BeforeSave(tx *gorm.DB) error {
	p.AvailableStock = p.Stock - p.ReservedStock
	
	// Auto-update status based on stock
	if p.Stock <= 0 && p.Status != ProductStatusInactive {
		p.Status = ProductStatusOutOfStock
	} else if p.Stock > 0 && p.Status == ProductStatusOutOfStock {
		if p.Status == ProductStatusActive {
			p.Status = ProductStatusActive
		}
	}
	
	return nil
}

// GetDiscountPrice returns the price after discount
func (p *Product) GetDiscountPrice() float64 {
	if p.DiscountPercent > 0 {
		return p.Price * (1 - float64(p.DiscountPercent)/100)
	}
	return p.Price
}

// IsOnSale checks if product is on sale
func (p *Product) IsOnSale() bool {
	return p.DiscountPercent > 0 || (p.OriginalPrice > 0 && p.OriginalPrice > p.Price)
}

// IsAvailable checks if product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.Status == ProductStatusActive && p.AvailableStock > 0
}

// HasVariants checks if product has variants
func (p *Product) HasVariants() bool {
	return len(p.Variants) > 0
}

// GetPrimaryImage returns the primary image of the product
func (p *Product) GetPrimaryImage() *ProductImage {
	for _, img := range p.Images {
		if img.IsPrimary {
			return &img
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

// GetTotalStock calculates total stock from all variants
func (p *Product) GetTotalStock() int {
	if len(p.Variants) == 0 {
		return p.Stock
	}
	
	total := 0
	for _, variant := range p.Variants {
		total += variant.Stock
	}
	return total
}

// Category.GetFullPath returns the full path of category names
func (c *Category) GetFullPath() string {
	if c.Parent == nil {
		return c.Name
	}
	return c.Parent.GetFullPath() + " > " + c.Name
}
