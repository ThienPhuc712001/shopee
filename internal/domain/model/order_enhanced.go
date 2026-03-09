package model

import (
	"time"

	"gorm.io/gorm"
)

// Note: PaymentStatus is defined in payment_enhanced.go
// Note: Shop is referenced but defined separately

// OrderStatus defines the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// FulfillmentStatus defines the fulfillment status
type FulfillmentStatus string

const (
	FulfillmentUnfulfilled FulfillmentStatus = "unfulfilled"
	FulfillmentProcessing  FulfillmentStatus = "processing"
	FulfillmentPacked      FulfillmentStatus = "packed"
	FulfillmentShipped     FulfillmentStatus = "shipped"
	FulfillmentDelivered   FulfillmentStatus = "delivered"
)

// Order represents a customer order
type Order struct {
	ID                  uint            `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderNumber         string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID              uint            `gorm:"type:int;not null;index" json:"user_id"`
	User                *User           `gorm:"foreignKey:UserID" json:"-"`
	ShopID              uint            `gorm:"type:int;not null;index" json:"shop_id"`
	Shop                *Shop           `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	ParentOrderID       *uint           `gorm:"type:int;index" json:"parent_order_id"`
	ParentOrder         *Order          `gorm:"foreignKey:ParentOrderID" json:"-"`
	Status              OrderStatus     `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	PaymentStatus       PaymentStatus   `gorm:"type:varchar(20);not null;default:'pending';index" json:"payment_status"`
	FulfillmentStatus   FulfillmentStatus `gorm:"type:varchar(20);default:'unfulfilled'" json:"fulfillment_status"`

	// Pricing
	Subtotal            float64         `gorm:"type:decimal(18,2);not null" json:"subtotal"`
	ShippingFee         float64         `gorm:"type:decimal(18,2);default:0" json:"shipping_fee"`
	ShippingDiscount    float64         `gorm:"type:decimal(18,2);default:0" json:"shipping_discount"`
	ProductDiscount     float64         `gorm:"type:decimal(18,2);default:0" json:"product_discount"`
	VoucherDiscount     float64         `gorm:"type:decimal(18,2);default:0" json:"voucher_discount"`
	TaxAmount           float64         `gorm:"type:decimal(18,2);default:0" json:"tax_amount"`
	TotalAmount         float64         `gorm:"type:decimal(18,2);not null" json:"total_amount"`
	PaidAmount          float64         `gorm:"type:decimal(18,2);default:0" json:"paid_amount"`

	// Shipping Address (snapshot at order time)
	ShippingName        string          `gorm:"type:varchar(200);not null" json:"shipping_name"`
	ShippingPhone       string          `gorm:"type:varchar(20);not null" json:"shipping_phone"`
	ShippingAddress     string          `gorm:"type:varchar(500);not null" json:"shipping_address"`
	ShippingWard        string          `gorm:"type:varchar(200)" json:"shipping_ward"`
	ShippingDistrict    string          `gorm:"type:varchar(200)" json:"shipping_district"`
	ShippingCity        string          `gorm:"type:varchar(200)" json:"shipping_city"`
	ShippingState       string          `gorm:"type:varchar(200)" json:"shipping_state"`
	ShippingCountry     string          `gorm:"type:varchar(100);default:'Vietnam'" json:"shipping_country"`
	ShippingPostalCode  string          `gorm:"type:varchar(20)" json:"shipping_postal_code"`

	// Shipping Method
	ShippingMethod      string          `gorm:"type:varchar(100)" json:"shipping_method"`
	ShippingCarrier     string          `gorm:"type:varchar(100)" json:"shipping_carrier"`
	TrackingNumber      string          `gorm:"type:varchar(100)" json:"tracking_number"`
	EstimatedDelivery   *time.Time      `gorm:"type:datetime" json:"estimated_delivery"`

	// Notes
	BuyerNote           string          `gorm:"type:varchar(500)" json:"buyer_note"`
	SellerNote          string          `gorm:"type:varchar(500)" json:"seller_note"`
	CancelReason        string          `gorm:"type:varchar(500)" json:"cancel_reason"`

	// Timestamps
	PaidAt              *time.Time      `gorm:"type:datetime" json:"paid_at"`
	ConfirmedAt         *time.Time      `gorm:"type:datetime" json:"confirmed_at"`
	ShippedAt           *time.Time      `gorm:"type:datetime" json:"shipped_at"`
	DeliveredAt         *time.Time      `gorm:"type:datetime" json:"delivered_at"`
	CancelledAt         *time.Time      `gorm:"type:datetime" json:"cancelled_at"`
	CompletedAt         *time.Time      `gorm:"type:datetime" json:"completed_at"`

	// Relationships
	Items       []OrderItem       `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	Payments    []Payment         `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
	StatusHistory []OrderStatusHistory `gorm:"foreignKey:OrderID" json:"status_history,omitempty"`
	Shipping    *OrderShipping    `gorm:"foreignKey:OrderID" json:"shipping,omitempty"`
	Tracking    []OrderTracking   `gorm:"foreignKey:OrderID" json:"tracking,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID                uint            `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderID           uint            `gorm:"type:int;not null;index" json:"order_id"`
	Order             *Order          `gorm:"foreignKey:OrderID" json:"-"`
	ProductID         uint            `gorm:"type:int;not null;index" json:"product_id"`
	Product           *Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	VariantID         *uint           `gorm:"type:int;index" json:"variant_id"`
	Variant           *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	ProductName       string          `gorm:"type:varchar(500);not null" json:"product_name"`
	ProductImage      string          `gorm:"type:varchar(500)" json:"product_image"`
	ProductSKU        string          `gorm:"type:varchar(100)" json:"product_sku"`
	Quantity          int             `gorm:"type:int;not null" json:"quantity"`
	Price             float64         `gorm:"type:decimal(18,2);not null" json:"price"`
	OriginalPrice     float64         `gorm:"type:decimal(18,2)" json:"original_price"`
	Discount          float64         `gorm:"type:decimal(18,2);default:0" json:"discount"`
	Subtotal          float64         `gorm:"type:decimal(18,2);not null" json:"subtotal"`
	TaxAmount         float64         `gorm:"type:decimal(18,2);default:0" json:"tax_amount"`
	FinalAmount       float64         `gorm:"type:decimal(18,2);not null" json:"final_amount"`
	ShopID            uint            `gorm:"type:int;not null;index" json:"shop_id"`
	Shop              *Shop           `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	FulfillmentStatus FulfillmentStatus `gorm:"type:varchar(20);default:'pending'" json:"fulfillment_status"`
	TrackingNumber    string          `gorm:"type:varchar(100)" json:"tracking_number"`
	ShippedAt         *time.Time      `gorm:"type:datetime" json:"shipped_at"`
	DeliveredAt       *time.Time      `gorm:"type:datetime" json:"delivered_at"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
}

// OrderStatusHistory tracks order status changes
type OrderStatusHistory struct {
	ID          uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderID     uint        `gorm:"type:int;not null;index" json:"order_id"`
	Order       *Order      `gorm:"foreignKey:OrderID" json:"-"`
	Status      OrderStatus `gorm:"type:varchar(20);not null" json:"status"`
	FromStatus  OrderStatus `gorm:"type:varchar(20)" json:"from_status"`
	Message     string      `gorm:"type:varchar(500)" json:"message"`
	ChangedBy   *uint       `gorm:"type:int" json:"changed_by"` // user_id or nil for system
	IPAddress   string      `gorm:"type:varchar(45)" json:"ip_address"`

	CreatedAt time.Time `gorm:"type:datetime;not null;index" json:"created_at"`
}

// OrderShipping contains shipping information
type OrderShipping struct {
	ID                uint       `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderID           uint       `gorm:"type:int;not null;uniqueIndex" json:"order_id"`
	Order             *Order     `gorm:"foreignKey:OrderID" json:"-"`
	CarrierID         *uint      `gorm:"type:int" json:"carrier_id"`
	CarrierName       string     `gorm:"type:varchar(100)" json:"carrier_name"`
	TrackingNumber    string     `gorm:"type:varchar(100)" json:"tracking_number"`
	ShippingLabelURL  string     `gorm:"type:varchar(500)" json:"shipping_label_url"`
	Weight            float64    `gorm:"type:decimal(10,2)" json:"weight"`
	Dimensions        string     `gorm:"type:varchar(50)" json:"dimensions"` // LxWxH
	ShippedAt         *time.Time `gorm:"type:datetime" json:"shipped_at"`
	EstimatedDelivery *time.Time `gorm:"type:datetime" json:"estimated_delivery"`
	ActualDelivery    *time.Time `gorm:"type:datetime" json:"actual_delivery"`
	Status            string     `gorm:"type:varchar(50)" json:"status"`
	TrackingEvents    string     `gorm:"type:text" json:"tracking_events"` // JSON array

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// OrderTracking represents a tracking event
type OrderTracking struct {
	ID            uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderID       uint      `gorm:"type:int;not null;index" json:"order_id"`
	Order         *Order    `gorm:"foreignKey:OrderID" json:"-"`
	TrackingNumber string   `gorm:"type:varchar(100);not null" json:"tracking_number"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status"`
	Location      string    `gorm:"type:varchar(200)" json:"location"`
	Description   string    `gorm:"type:varchar(500);not null" json:"description"`
	Timestamp     time.Time `gorm:"type:datetime;not null" json:"timestamp"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

// BeforeCreate hook for Order - generates order number
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.OrderNumber == "" {
		o.OrderNumber = GenerateOrderNumber()
	}
	return nil
}

// BeforeCreate hook for OrderItem - calculates amounts
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	oi.Subtotal = oi.Price * float64(oi.Quantity)
	oi.FinalAmount = oi.Subtotal - oi.Discount + oi.TaxAmount
	return nil
}

// GetTotalPrice returns the total price of the order
func (o *Order) GetTotalPrice() float64 {
	return o.TotalAmount
}

// GetTotalItems returns the total number of items in the order
func (o *Order) GetTotalItems() int {
	total := 0
	for _, item := range o.Items {
		total += item.Quantity
	}
	return total
}

// CanCancel checks if order can be cancelled
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusPaid
}

// CanRefund checks if order can be refunded
func (o *Order) CanRefund() bool {
	return o.PaymentStatus == PaymentStatusPaid &&
		(o.Status == OrderStatusDelivered || o.Status == OrderStatusShipped)
}

// IsPending checks if order is pending payment
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsPaid checks if order is paid
func (o *Order) IsPaid() bool {
	return o.PaymentStatus == PaymentStatusPaid
}

// IsCompleted checks if order is completed
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusDelivered
}

// IsCancelled checks if order is cancelled
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	// Format: ORD-YYYYMMDD-XXXXX
	// In production, use a more robust unique generator
	return "ORD-" + time.Now().Format("20060102") + "-00001"
}

// OrderInput represents input for creating an order
type OrderInput struct {
	ShippingInfo     ShippingInfo `json:"shipping_info" binding:"required"`
	PaymentMethod    string       `json:"payment_method" binding:"required"`
	BuyerNote        string       `json:"buyer_note"`
	VoucherCode      string       `json:"voucher_code"`
	ShippingMethod   string       `json:"shipping_method"`
	EstimatedDelivery *time.Time  `json:"estimated_delivery"`
}

// ShippingInfo contains shipping address information
type ShippingInfo struct {
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Ward         string `json:"ward"`
	District     string `json:"district" binding:"required"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state"`
	Country      string `json:"country" default:"Vietnam"`
	PostalCode   string `json:"postal_code"`
}

// OrderSummary represents a summary of an order
type OrderSummary struct {
	ID              uint        `json:"id"`
	OrderNumber     string      `json:"order_number"`
	Status          OrderStatus `json:"status"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	TotalAmount     float64     `json:"total_amount"`
	TotalItems      int         `json:"total_items"`
	CreatedAt       time.Time   `json:"created_at"`
}

// OrderWithItems represents an order with items loaded
type OrderWithItems struct {
	Order
	Items []OrderItemWithProduct `json:"items"`
}

// OrderItemWithProduct represents an order item with product info
type OrderItemWithProduct struct {
	OrderItem
	Product *Product `json:"product"`
	Variant *ProductVariant `json:"variant"`
}

// OrderStatusChange represents a status change request
type OrderStatusChange struct {
	Status  OrderStatus `json:"status" binding:"required"`
	Message string      `json:"message"`
}

// CancelOrderInput represents input for cancelling an order
type CancelOrderInput struct {
	Reason string `json:"reason" binding:"required"`
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders     int64   `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	PendingOrders   int64   `json:"pending_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}
