package repository

import (
	"ecommerce/internal/domain/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// ShippingRepository handles database operations for shipping
type ShippingRepository struct {
	db *gorm.DB
}

// NewShippingRepository creates a new shipping repository
func NewShippingRepository(db *gorm.DB) *ShippingRepository {
	return &ShippingRepository{db: db}
}

// ==================== SHIPPING ADDRESS ====================

// CreateAddress creates a new shipping address
func (r *ShippingRepository) CreateAddress(ctx context.Context, address *model.ShippingAddress) error {
	return r.db.WithContext(ctx).Create(address).Error
}

// UpdateAddress updates an existing shipping address
func (r *ShippingRepository) UpdateAddress(ctx context.Context, address *model.ShippingAddress) error {
	return r.db.WithContext(ctx).Save(address).Error
}

// DeleteAddress deletes a shipping address (soft delete)
func (r *ShippingRepository) DeleteAddress(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ShippingAddress{}, id).Error
}

// GetAddressByID retrieves a shipping address by ID
func (r *ShippingRepository) GetAddressByID(ctx context.Context, id uint) (*model.ShippingAddress, error) {
	var address model.ShippingAddress
	err := r.db.WithContext(ctx).First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetAddressesByUser retrieves all shipping addresses for a user
func (r *ShippingRepository) GetAddressesByUser(ctx context.Context, userID uint) ([]model.ShippingAddress, error) {
	var addresses []model.ShippingAddress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// GetDefaultAddress retrieves the default shipping address for a user
func (r *ShippingRepository) GetDefaultAddress(ctx context.Context, userID uint) (*model.ShippingAddress, error) {
	var address model.ShippingAddress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ? AND deleted_at IS NULL", userID, true).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// SetDefaultAddress sets a shipping address as default for a user
func (r *ShippingRepository) SetDefaultAddress(ctx context.Context, userID, addressID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all default addresses for this user
		if err := tx.Model(&model.ShippingAddress{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set the specified address as default
		return tx.Model(&model.ShippingAddress{}).
			Where("id = ? AND user_id = ?", addressID, userID).
			Update("is_default", true).Error
	})
}

// ==================== SHIPMENT ====================

// CreateShipment creates a new shipment
func (r *ShippingRepository) CreateShipment(ctx context.Context, shipment *model.Shipment) error {
	return r.db.WithContext(ctx).Create(shipment).Error
}

// UpdateShipment updates an existing shipment
func (r *ShippingRepository) UpdateShipment(ctx context.Context, shipment *model.Shipment) error {
	return r.db.WithContext(ctx).Save(shipment).Error
}

// DeleteShipment deletes a shipment (soft delete)
func (r *ShippingRepository) DeleteShipment(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Shipment{}, id).Error
}

// GetShipmentByID retrieves a shipment by ID
func (r *ShippingRepository) GetShipmentByID(ctx context.Context, id uint) (*model.Shipment, error) {
	var shipment model.Shipment
	err := r.db.WithContext(ctx).
		Preload("Tracking").
		First(&shipment, id).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

// GetShipmentByOrderID retrieves a shipment by order ID
func (r *ShippingRepository) GetShipmentByOrderID(ctx context.Context, orderID uint) (*model.Shipment, error) {
	var shipment model.Shipment
	err := r.db.WithContext(ctx).
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("occurred_at DESC")
		}).
		Where("order_id = ?", orderID).
		First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

// GetShipmentsByStatus retrieves shipments by status
func (r *ShippingRepository) GetShipmentsByStatus(ctx context.Context, status model.ShipmentStatus, page, limit int) ([]model.Shipment, int64, error) {
	var shipments []model.Shipment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Shipment{}).
		Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Get shipments
	err := query.
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("occurred_at DESC").Limit(10)
		}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error

	return shipments, total, err
}

// GetShipmentsByCarrier retrieves shipments by carrier
func (r *ShippingRepository) GetShipmentsByCarrier(ctx context.Context, carrierID uint, page, limit int) ([]model.Shipment, int64, error) {
	var shipments []model.Shipment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Shipment{}).
		Where("carrier_id = ?", carrierID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Get shipments
	err := query.
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("occurred_at DESC").Limit(10)
		}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error

	return shipments, total, err
}

// UpdateShipmentStatus updates the status of a shipment
func (r *ShippingRepository) UpdateShipmentStatus(ctx context.Context, shipmentID uint, status model.ShipmentStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Set timestamp based on status
	now := time.Now()
	switch status {
	case model.ShipmentStatusShipped:
		updates["shipped_at"] = now
	case model.ShipmentStatusDelivered:
		updates["delivered_at"] = now
	case model.ShipmentStatusFailed:
		updates["failed_at"] = now
	}

	return r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Where("id = ?", shipmentID).
		Updates(updates).Error
}

// ==================== SHIPMENT TRACKING ====================

// AddTrackingEvent adds a tracking event to a shipment
func (r *ShippingRepository) AddTrackingEvent(ctx context.Context, event *model.ShipmentTracking) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetTrackingByShipmentID retrieves all tracking events for a shipment
func (r *ShippingRepository) GetTrackingByShipmentID(ctx context.Context, shipmentID uint) ([]model.ShipmentTracking, error) {
	var events []model.ShipmentTracking
	err := r.db.WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("occurred_at DESC").
		Find(&events).Error
	return events, err
}

// GetTrackingTimeline retrieves tracking timeline for a shipment
func (r *ShippingRepository) GetTrackingTimeline(ctx context.Context, shipmentID uint) (*model.TrackingTimeline, error) {
	var shipment model.Shipment
	err := r.db.WithContext(ctx).
		Where("id = ?", shipmentID).
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("occurred_at DESC")
		}).
		First(&shipment).Error
	if err != nil {
		return nil, err
	}

	return shipment.ToTrackingTimeline(), nil
}

// ==================== SHIPPING CARRIER ====================

// CreateCarrier creates a new shipping carrier
func (r *ShippingRepository) CreateCarrier(ctx context.Context, carrier *model.ShippingCarrier) error {
	return r.db.WithContext(ctx).Create(carrier).Error
}

// UpdateCarrier updates an existing shipping carrier
func (r *ShippingRepository) UpdateCarrier(ctx context.Context, carrier *model.ShippingCarrier) error {
	return r.db.WithContext(ctx).Save(carrier).Error
}

// DeleteCarrier deletes a shipping carrier
func (r *ShippingRepository) DeleteCarrier(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ShippingCarrier{}, id).Error
}

// GetCarrierByID retrieves a shipping carrier by ID
func (r *ShippingRepository) GetCarrierByID(ctx context.Context, id uint) (*model.ShippingCarrier, error) {
	var carrier model.ShippingCarrier
	err := r.db.WithContext(ctx).First(&carrier, id).Error
	if err != nil {
		return nil, err
	}
	return &carrier, nil
}

// GetCarrierByCode retrieves a shipping carrier by code
func (r *ShippingRepository) GetCarrierByCode(ctx context.Context, code string) (*model.ShippingCarrier, error) {
	var carrier model.ShippingCarrier
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&carrier).Error
	if err != nil {
		return nil, err
	}
	return &carrier, nil
}

// GetActiveCarriers retrieves all active shipping carriers
func (r *ShippingRepository) GetActiveCarriers(ctx context.Context) ([]model.ShippingCarrier, error) {
	var carriers []model.ShippingCarrier
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&carriers).Error
	return carriers, err
}

// GetAllCarriers retrieves all shipping carriers
func (r *ShippingRepository) GetAllCarriers(ctx context.Context) ([]model.ShippingCarrier, error) {
	var carriers []model.ShippingCarrier
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&carriers).Error
	return carriers, err
}

// ==================== STATISTICS ====================

// GetShipmentStats retrieves shipment statistics
type ShipmentStats struct {
	TotalShipments   int64   `json:"total_shipments"`
	PendingShipments int64   `json:"pending_shipments"`
	ShippedShipments int64   `json:"shipped_shipments"`
	DeliveredShipments int64 `json:"delivered_shipments"`
	FailedShipments  int64   `json:"failed_shipments"`
	AverageDeliveryDays float64 `json:"average_delivery_days"`
}

// GetShipmentStats retrieves shipment statistics
func (r *ShippingRepository) GetShipmentStats(ctx context.Context) (*ShipmentStats, error) {
	stats := &ShipmentStats{}

	// Total shipments
	if err := r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Count(&stats.TotalShipments).Error; err != nil {
		return nil, err
	}

	// Pending shipments
	if err := r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Where("status = ?", model.ShipmentStatusPending).
		Count(&stats.PendingShipments).Error; err != nil {
		return nil, err
	}

	// Shipped shipments (in transit)
	if err := r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Where("status IN ?", []model.ShipmentStatus{
			model.ShipmentStatusShipped,
			model.ShipmentStatusInTransit,
			model.ShipmentStatusOutForDelivery,
		}).
		Count(&stats.ShippedShipments).Error; err != nil {
		return nil, err
	}

	// Delivered shipments
	if err := r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Where("status = ?", model.ShipmentStatusDelivered).
		Count(&stats.DeliveredShipments).Error; err != nil {
		return nil, err
	}

	// Failed shipments
	if err := r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Where("status = ?", model.ShipmentStatusFailed).
		Count(&stats.FailedShipments).Error; err != nil {
		return nil, err
	}

	// Average delivery days (for delivered shipments)
	type DeliveryDays struct {
		Days float64
	}
	var deliveryDays DeliveryDays
	r.db.WithContext(ctx).
		Model(&model.Shipment{}).
		Select("AVG(DATEDIFF(day, shipped_at, delivered_at)) as days").
		Where("status = ? AND shipped_at IS NOT NULL AND delivered_at IS NOT NULL", model.ShipmentStatusDelivered).
		Scan(&deliveryDays)
	stats.AverageDeliveryDays = deliveryDays.Days

	return stats, nil
}

// GetShipmentsByDateRange retrieves shipments within a date range
func (r *ShippingRepository) GetShipmentsByDateRange(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]model.Shipment, int64, error) {
	var shipments []model.Shipment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Shipment{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Get shipments
	err := query.
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("occurred_at DESC").Limit(5)
		}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error

	return shipments, total, err
}
