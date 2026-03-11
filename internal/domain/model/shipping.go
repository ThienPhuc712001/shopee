package model

import (
	"time"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusPending      ShipmentStatus = "pending"
	ShipmentStatusConfirmed    ShipmentStatus = "confirmed"
	ShipmentStatusProcessing   ShipmentStatus = "processing"
	ShipmentStatusPacked       ShipmentStatus = "packed"
	ShipmentStatusShipped      ShipmentStatus = "shipped"
	ShipmentStatusInTransit    ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered    ShipmentStatus = "delivered"
	ShipmentStatusFailed       ShipmentStatus = "failed"
	ShipmentStatusReturned     ShipmentStatus = "returned"
	ShipmentStatusCancelled    ShipmentStatus = "cancelled"
)

// CarrierType represents the type of shipping carrier
type CarrierType string

const (
	CarrierTypeInternal   CarrierType = "internal"
	CarrierTypeThirdParty CarrierType = "third_party"
	CarrierTypeLocal      CarrierType = "local"
)

// ShippingAddress represents a user's saved shipping address
type ShippingAddress struct {
	ID             uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"type:int;not null;index" json:"user_id"`
	RecipientName  string    `gorm:"type:varchar(200);not null" json:"recipient_name"`
	Phone          string    `gorm:"type:varchar(20);not null" json:"phone"`
	AddressLine    string    `gorm:"type:varchar(500);not null" json:"address_line"`
	Ward           string    `gorm:"type:varchar(200)" json:"ward"`
	District       string    `gorm:"type:varchar(200)" json:"district"`
	City           string    `gorm:"type:varchar(200);not null" json:"city"`
	PostalCode     string    `gorm:"type:varchar(20)" json:"postal_code"`
	Country        string    `gorm:"type:varchar(100);default:'Vietnam'" json:"country"`
	IsDefault      bool      `gorm:"type:bit;default:false" json:"is_default"`
	Latitude       float64   `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude      float64   `gorm:"type:decimal(11,8)" json:"longitude"`
	Notes          string    `gorm:"type:varchar(500)" json:"notes"`
	CreatedAt      time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"-"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for ShippingAddress
func (ShippingAddress) TableName() string {
	return "shipping_addresses"
}

// Shipment represents a shipment for an order
type Shipment struct {
	ID              uint           `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderID         uint           `gorm:"type:int;not null;uniqueIndex" json:"order_id"`
	CarrierID       *uint          `gorm:"type:int" json:"carrier_id"`
	CarrierName     string         `gorm:"type:varchar(200)" json:"carrier_name"`
	CarrierType     CarrierType    `gorm:"type:varchar(50);default:'third_party'" json:"carrier_type"`
	TrackingNumber  string         `gorm:"type:varchar(100);index" json:"tracking_number"`
	Status          ShipmentStatus `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"`
	
	// Shipping details
	ShippingFrom    string         `gorm:"type:varchar(500)" json:"shipping_from"`
	ShippingTo      string         `gorm:"type:varchar(500)" json:"shipping_to"`
	Weight          float64        `gorm:"type:decimal(10,2)" json:"weight"`
	Dimensions      string         `gorm:"type:varchar(50)" json:"dimensions"` // LxWxH
	PackageCount    int            `gorm:"type:int;default:1" json:"package_count"`
	
	// Timestamps
	ShippedAt       *time.Time     `gorm:"type:datetime" json:"shipped_at"`
	EstimatedDelivery *time.Time   `gorm:"type:datetime" json:"estimated_delivery"`
	DeliveredAt     *time.Time     `gorm:"type:datetime" json:"delivered_at"`
	FailedAt        *time.Time     `gorm:"type:datetime" json:"failed_at"`
	FailureReason   string         `gorm:"type:varchar(500)" json:"failure_reason"`
	
	// Additional info
	ShippingFee     float64        `gorm:"type:decimal(18,2);default:0" json:"shipping_fee"`
	InsuranceAmount float64        `gorm:"type:decimal(18,2);default:0" json:"insurance_amount"`
	Notes           string         `gorm:"type:varchar(500)" json:"notes"`
	
	CreatedAt       time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt       *time.Time     `gorm:"index" json:"-"`

	// Relationships
	Order    *Order            `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Tracking []ShipmentTracking `gorm:"foreignKey:ShipmentID" json:"tracking,omitempty"`
}

// TableName specifies the table name for Shipment
func (Shipment) TableName() string {
	return "shipments"
}

// ShipmentTracking represents a tracking event for a shipment
type ShipmentTracking struct {
	ID           uint           `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	ShipmentID   uint           `gorm:"type:int;not null;index" json:"shipment_id"`
	Status       ShipmentStatus `gorm:"type:varchar(50);not null" json:"status"`
	Location     string         `gorm:"type:varchar(300)" json:"location"`
	Description  string         `gorm:"type:varchar(1000);not null" json:"description"`
	OccurredAt   time.Time      `gorm:"type:datetime;not null" json:"occurred_at"`
	CreatedAt    time.Time      `gorm:"type:datetime;not null" json:"created_at"`

	// Relationships
	Shipment *Shipment `gorm:"foreignKey:ShipmentID" json:"shipment,omitempty"`
}

// TableName specifies the table name for ShipmentTracking
func (ShipmentTracking) TableName() string {
	return "shipment_tracking"
}

// ShippingCarrier represents a shipping carrier/courier service
type ShippingCarrier struct {
	ID          uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(200);not null;uniqueIndex" json:"name"`
	Code        string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Type        CarrierType `gorm:"type:varchar(50);not null" json:"type"`
	ContactName string    `gorm:"type:varchar(200)" json:"contact_name"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`
	Email       string    `gorm:"type:varchar(255)" json:"email"`
	Website     string    `gorm:"type:varchar(500)" json:"website"`
	APIEndpoint string    `gorm:"type:varchar(500)" json:"api_endpoint"`
	APIKey      string    `gorm:"type:varchar(500)" json:"-"` // Encrypted in production
	IsActive    bool      `gorm:"type:bit;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// TableName specifies the table name for ShippingCarrier
func (ShippingCarrier) TableName() string {
	return "shipping_carriers"
}

// IsShipped checks if the shipment has been shipped
func (s *Shipment) IsShipped() bool {
	return s.Status == ShipmentStatusShipped || 
		   s.Status == ShipmentStatusInTransit || 
		   s.Status == ShipmentStatusOutForDelivery ||
		   s.Status == ShipmentStatusDelivered
}

// IsDelivered checks if the shipment has been delivered
func (s *Shipment) IsDelivered() bool {
	return s.Status == ShipmentStatusDelivered
}

// IsCancelled checks if the shipment has been cancelled
func (s *Shipment) IsCancelled() bool {
	return s.Status == ShipmentStatusCancelled
}

// CanCancel checks if shipment can be cancelled
func (s *Shipment) CanCancel() bool {
	return s.Status == ShipmentStatusPending || 
		   s.Status == ShipmentStatusConfirmed || 
		   s.Status == ShipmentStatusProcessing
}

// ShipmentInput represents input for creating a shipment
type ShipmentInput struct {
	OrderID         uint    `json:"order_id" binding:"required"`
	CarrierID       *uint   `json:"carrier_id"`
	CarrierName     string  `json:"carrier_name"`
	CarrierType     string  `json:"carrier_type"`
	TrackingNumber  string  `json:"tracking_number"`
	ShippingFrom    string  `json:"shipping_from"`
	ShippingTo      string  `json:"shipping_to"`
	Weight          float64 `json:"weight"`
	Dimensions      string  `json:"dimensions"`
	PackageCount    int     `json:"package_count"`
	EstimatedDelivery string `json:"estimated_delivery"` // RFC3339 format
	ShippingFee     float64 `json:"shipping_fee"`
	Notes           string  `json:"notes"`
}

// TrackingEventInput represents input for adding a tracking event
type TrackingEventInput struct {
	Status      string `json:"status" binding:"required"`
	Location    string `json:"location"`
	Description string `json:"description" binding:"required"`
	OccurredAt  string `json:"occurred_at"` // RFC3339 format, defaults to now
}

// ShippingAddressInput represents input for creating/updating a shipping address
type ShippingAddressInput struct {
	RecipientName string  `json:"recipient_name" binding:"required"`
	Phone         string  `json:"phone" binding:"required"`
	AddressLine   string  `json:"address_line" binding:"required"`
	Ward          string  `json:"ward"`
	District      string  `json:"district" binding:"required"`
	City          string  `json:"city" binding:"required"`
	PostalCode    string  `json:"postal_code"`
	Country       string  `json:"country"`
	IsDefault     bool    `json:"is_default"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Notes         string  `json:"notes"`
}

// TrackingTimeline represents a timeline view of tracking events
type TrackingTimeline struct {
	ShipmentID      uint                `json:"shipment_id"`
	TrackingNumber  string              `json:"tracking_number"`
	CurrentStatus   ShipmentStatus      `json:"current_status"`
	EstimatedDelivery *time.Time        `json:"estimated_delivery"`
	DeliveredAt     *time.Time          `json:"delivered_at"`
	Events          []TrackingEventView `json:"events"`
}

// TrackingEventView represents a single tracking event in the timeline
type TrackingEventView struct {
	ID          uint           `json:"id"`
	Status      ShipmentStatus `json:"status"`
	Location    string         `json:"location"`
	Description string         `json:"description"`
	OccurredAt  time.Time      `json:"occurred_at"`
}

// ToTrackingTimeline converts a shipment with tracking to a timeline view
func (s *Shipment) ToTrackingTimeline() *TrackingTimeline {
	timeline := &TrackingTimeline{
		ShipmentID:      s.ID,
		TrackingNumber:  s.TrackingNumber,
		CurrentStatus:   s.Status,
		EstimatedDelivery: s.EstimatedDelivery,
		DeliveredAt:     s.DeliveredAt,
		Events:          make([]TrackingEventView, 0),
	}

	for _, t := range s.Tracking {
		timeline.Events = append(timeline.Events, TrackingEventView{
			ID:          t.ID,
			Status:      t.Status,
			Location:    t.Location,
			Description: t.Description,
			OccurredAt:  t.OccurredAt,
		})
	}

	return timeline
}
