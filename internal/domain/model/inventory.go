package model

import (
	"time"
)

// InventoryChangeType represents the type of inventory change
type InventoryChangeType string

const (
	InventoryChangeRestock   InventoryChangeType = "restock"
	InventoryChangeReserve   InventoryChangeType = "reserve"
	InventoryChangeDeduct    InventoryChangeType = "deduct"
	InventoryChangeRelease   InventoryChangeType = "release"
	InventoryChangeReturn    InventoryChangeType = "return"
	InventoryChangeAdjust    InventoryChangeType = "adjust"
	InventoryChangeDamaged   InventoryChangeType = "damaged"
)

// Inventory represents stock information for a product
type Inventory struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID        uint       `gorm:"type:int;not null;uniqueIndex" json:"product_id"`
	StockQuantity    int        `gorm:"type:int;not null;default:0" json:"stock_quantity"`
	ReservedQuantity int        `gorm:"type:int;not null;default:0" json:"reserved_quantity"`
	AvailableQuantity int       `gorm:"type:int;not null;default:0" json:"available_quantity"` // Computed: Stock - Reserved
	WarehouseLocation string    `gorm:"type:varchar(100)" json:"warehouse_location"`
	LastStockCheck   *time.Time `gorm:"type:datetime" json:"last_stock_check"`
	ReorderLevel     int        `gorm:"type:int;default:0" json:"reorder_level"` // Alert when stock below this
	ReorderQuantity  int        `gorm:"type:int;default:0" json:"reorder_quantity"` // Quantity to reorder
	CreatedAt        time.Time  `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:datetime;not null" json:"updated_at"`

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for Inventory
func (Inventory) TableName() string {
	return "inventory"
}

// InventoryLog tracks all inventory changes
type InventoryLog struct {
	ID            uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID     uint                `gorm:"type:int;not null;index" json:"product_id"`
	InventoryID   uint                `gorm:"type:int;index" json:"inventory_id"`
	ChangeType    InventoryChangeType `gorm:"type:varchar(50);not null" json:"change_type"`
	Quantity      int                 `gorm:"type:int;not null" json:"quantity"` // Positive for add, negative for deduct
	StockBefore   int                 `gorm:"type:int;not null" json:"stock_before"`
	StockAfter    int                 `gorm:"type:int;not null" json:"stock_after"`
	ReservedBefore int                `gorm:"type:int;not null" json:"reserved_before"`
	ReservedAfter int                 `gorm:"type:int;not null" json:"reserved_after"`
	ReferenceType string              `gorm:"type:varchar(50)" json:"reference_type"` // order, restock, adjustment, etc.
	ReferenceID   string              `gorm:"type:varchar(100);index" json:"reference_id"` // Order ID, Restock ID, etc.
	UserID        *uint               `gorm:"type:int" json:"user_id"` // Who made the change
	Reason        string              `gorm:"type:text" json:"reason"`
	IPAddress     string              `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt     time.Time           `gorm:"type:datetime;not null;index" json:"created_at"`
}

// TableName specifies the table name for InventoryLog
func (InventoryLog) TableName() string {
	return "inventory_logs"
}

// StockAlert represents low stock alerts
type StockAlert struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint       `gorm:"type:int;not null;uniqueIndex" json:"product_id"`
	AlertType string     `gorm:"type:varchar(50)" json:"alert_type"` // low_stock, out_of_stock, overstock
	Threshold int        `gorm:"type:int" json:"threshold"`
	CurrentStock int     `gorm:"type:int" json:"current_stock"`
	IsResolved bool      `gorm:"type:bit;default:false" json:"is_resolved"`
	ResolvedAt *time.Time `gorm:"type:datetime" json:"resolved_at"`
	ResolvedBy *uint      `gorm:"type:int" json:"resolved_by"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// TableName specifies the table name for StockAlert
func (StockAlert) TableName() string {
	return "stock_alerts"
}

// InventorySummary represents inventory statistics
type InventorySummary struct {
	TotalProducts       int64  `json:"total_products"`
	InStockProducts     int64  `json:"in_stock_products"`
	LowStockProducts    int64  `json:"low_stock_products"`
	OutOfStockProducts  int64  `json:"out_of_stock_products"`
	TotalStockValue     float64 `json:"total_stock_value"`
	TotalReservedValue  float64 `json:"total_reserved_value"`
}

// StockStatus represents the stock status of a product
type StockStatus string

const (
	StockStatusInStock     StockStatus = "in_stock"
	StockStatusLowStock    StockStatus = "low_stock"
	StockStatusOutOfStock  StockStatus = "out_of_stock"
	StockStatusPreOrder    StockStatus = "pre_order"
)

// GetStockStatus returns the stock status based on quantities
func (i *Inventory) GetStockStatus() StockStatus {
	if i.AvailableQuantity <= 0 {
		if i.StockQuantity > 0 {
			return StockStatusPreOrder
		}
		return StockStatusOutOfStock
	}
	
	if i.AvailableQuantity <= i.ReorderLevel {
		return StockStatusLowStock
	}
	
	return StockStatusInStock
}

// IsAvailable checks if the product is available for purchase
func (i *Inventory) IsAvailable() bool {
	return i.AvailableQuantity > 0
}

// CanFulfillQuantity checks if a specific quantity can be fulfilled
func (i *Inventory) CanFulfillQuantity(quantity int) bool {
	return i.AvailableQuantity >= quantity
}
