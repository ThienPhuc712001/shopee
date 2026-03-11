package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

// ShippingService handles shipping business logic
type ShippingService struct {
	repo *repository.ShippingRepository
}

// NewShippingService creates a new shipping service
func NewShippingService(repo *repository.ShippingRepository) *ShippingService {
	return &ShippingService{repo: repo}
}

// ==================== SHIPPING ADDRESS ====================

// CreateAddress creates a new shipping address
func (s *ShippingService) CreateAddress(ctx context.Context, userID uint, input model.ShippingAddressInput) (*model.ShippingAddress, error) {
	// Validate input
	if err := s.validateAddressInput(input); err != nil {
		return nil, err
	}

	address := &model.ShippingAddress{
		UserID:        userID,
		RecipientName: input.RecipientName,
		Phone:         input.Phone,
		AddressLine:   input.AddressLine,
		Ward:          input.Ward,
		District:      input.District,
		City:          input.City,
		PostalCode:    input.PostalCode,
		Country:       input.Country,
		IsDefault:     input.IsDefault,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		Notes:         input.Notes,
	}

	// If this is set as default, unset other defaults
	if input.IsDefault {
		if err := s.repo.SetDefaultAddress(ctx, userID, 0); err != nil {
			return nil, err
		}
	}

	if err := s.repo.CreateAddress(ctx, address); err != nil {
		return nil, err
	}

	return address, nil
}

// UpdateAddress updates an existing shipping address
func (s *ShippingService) UpdateAddress(ctx context.Context, addressID uint, input model.ShippingAddressInput) (*model.ShippingAddress, error) {
	// Get existing address
	address, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}

	// Update fields
	address.RecipientName = input.RecipientName
	address.Phone = input.Phone
	address.AddressLine = input.AddressLine
	address.Ward = input.Ward
	address.District = input.District
	address.City = input.City
	address.PostalCode = input.PostalCode
	address.Country = input.Country
	address.Latitude = input.Latitude
	address.Longitude = input.Longitude
	address.Notes = input.Notes

	// If this is set as default, unset other defaults
	if input.IsDefault && !address.IsDefault {
		if err := s.repo.SetDefaultAddress(ctx, address.UserID, addressID); err != nil {
			return nil, err
		}
		address.IsDefault = true
	}

	if err := s.repo.UpdateAddress(ctx, address); err != nil {
		return nil, err
	}

	return address, nil
}

// DeleteAddress deletes a shipping address
func (s *ShippingService) DeleteAddress(ctx context.Context, addressID uint) error {
	_, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return ErrAddressNotFound
	}

	return s.repo.DeleteAddress(ctx, addressID)
}

// GetAddressByID retrieves a shipping address by ID
func (s *ShippingService) GetAddressByID(ctx context.Context, addressID uint) (*model.ShippingAddress, error) {
	address, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}
	return address, nil
}

// GetAddressesByUser retrieves all shipping addresses for a user
func (s *ShippingService) GetAddressesByUser(ctx context.Context, userID uint) ([]model.ShippingAddress, error) {
	return s.repo.GetAddressesByUser(ctx, userID)
}

// GetDefaultAddress retrieves the default shipping address for a user
func (s *ShippingService) GetDefaultAddress(ctx context.Context, userID uint) (*model.ShippingAddress, error) {
	return s.repo.GetDefaultAddress(ctx, userID)
}

// SetDefaultAddress sets a shipping address as default
func (s *ShippingService) SetDefaultAddress(ctx context.Context, userID, addressID uint) error {
	return s.repo.SetDefaultAddress(ctx, userID, addressID)
}

// ==================== SHIPMENT ====================

// CreateShipmentInput represents input for creating a shipment
type CreateShipmentInput struct {
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
	EstimatedDelivery string `json:"estimated_delivery"`
	ShippingFee     float64 `json:"shipping_fee"`
	Notes           string  `json:"notes"`
}

// CreateShipment creates a new shipment
func (s *ShippingService) CreateShipment(ctx context.Context, input CreateShipmentInput) (*model.Shipment, error) {
	// Validate input
	if err := s.validateShipmentInput(input); err != nil {
		return nil, err
	}

	// Check if order already has a shipment
	existing, err := s.repo.GetShipmentByOrderID(ctx, input.OrderID)
	if err == nil && existing != nil {
		return nil, ErrShipmentExists
	}

	// Parse estimated delivery date
	var estimatedDelivery *time.Time
	if input.EstimatedDelivery != "" {
		t, err := time.Parse(time.RFC3339, input.EstimatedDelivery)
		if err != nil {
			return nil, errors.New("invalid date format, use RFC3339")
		}
		estimatedDelivery = &t
	}

	// Determine carrier type
	carrierType := model.CarrierTypeThirdParty
	if input.CarrierType != "" {
		carrierType = model.CarrierType(input.CarrierType)
	}

	shipment := &model.Shipment{
		OrderID:         input.OrderID,
		CarrierID:       input.CarrierID,
		CarrierName:     input.CarrierName,
		CarrierType:     carrierType,
		TrackingNumber:  input.TrackingNumber,
		Status:          model.ShipmentStatusPending,
		ShippingFrom:    input.ShippingFrom,
		ShippingTo:      input.ShippingTo,
		Weight:          input.Weight,
		Dimensions:      input.Dimensions,
		PackageCount:    input.PackageCount,
		EstimatedDelivery: estimatedDelivery,
		ShippingFee:     input.ShippingFee,
		Notes:           input.Notes,
	}

	if err := s.repo.CreateShipment(ctx, shipment); err != nil {
		return nil, err
	}

	// Add initial tracking event
	initialEvent := &model.ShipmentTracking{
		ShipmentID:  shipment.ID,
		Status:      model.ShipmentStatusPending,
		Description: "Shipment created, awaiting confirmation",
		OccurredAt:  time.Now(),
	}
	_ = s.repo.AddTrackingEvent(ctx, initialEvent)

	return shipment, nil
}

// UpdateShipment updates an existing shipment
func (s *ShippingService) UpdateShipment(ctx context.Context, shipmentID uint, input CreateShipmentInput) (*model.Shipment, error) {
	shipment, err := s.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, ErrShipmentNotFound
	}

	// Update fields
	if input.CarrierID != nil {
		shipment.CarrierID = input.CarrierID
	}
	if input.CarrierName != "" {
		shipment.CarrierName = input.CarrierName
	}
	if input.CarrierType != "" {
		shipment.CarrierType = model.CarrierType(input.CarrierType)
	}
	if input.TrackingNumber != "" {
		shipment.TrackingNumber = input.TrackingNumber
	}
	if input.ShippingFrom != "" {
		shipment.ShippingFrom = input.ShippingFrom
	}
	if input.ShippingTo != "" {
		shipment.ShippingTo = input.ShippingTo
	}
	if input.Weight > 0 {
		shipment.Weight = input.Weight
	}
	if input.Dimensions != "" {
		shipment.Dimensions = input.Dimensions
	}
	if input.PackageCount > 0 {
		shipment.PackageCount = input.PackageCount
	}
	if input.EstimatedDelivery != "" {
		t, err := time.Parse(time.RFC3339, input.EstimatedDelivery)
		if err == nil {
			shipment.EstimatedDelivery = &t
		}
	}
	if input.Notes != "" {
		shipment.Notes = input.Notes
	}

	if err := s.repo.UpdateShipment(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

// UpdateShipmentStatus updates the status of a shipment
func (s *ShippingService) UpdateShipmentStatus(ctx context.Context, shipmentID uint, status model.ShipmentStatus, notes string) (*model.Shipment, error) {
	shipment, err := s.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, ErrShipmentNotFound
	}

	// Validate status transition
	if !s.isValidStatusTransition(shipment.Status, status) {
		return nil, ErrInvalidStatusTransition
	}

	// Update status
	oldStatus := shipment.Status
	shipment.Status = status

	// Set timestamps based on status
	now := time.Now()
	switch status {
	case model.ShipmentStatusShipped:
		shipment.ShippedAt = &now
	case model.ShipmentStatusDelivered:
		shipment.DeliveredAt = &now
	case model.ShipmentStatusFailed:
		shipment.FailedAt = &now
		shipment.FailureReason = notes
	}

	if err := s.repo.UpdateShipment(ctx, shipment); err != nil {
		return nil, err
	}

	// Add tracking event for status change
	event := &model.ShipmentTracking{
		ShipmentID:  shipment.ID,
		Status:      status,
		Description: fmt.Sprintf("Status changed from %s to %s. %s", oldStatus, status, notes),
		OccurredAt:  now,
	}
	_ = s.repo.AddTrackingEvent(ctx, event)

	return shipment, nil
}

// DeleteShipment deletes a shipment
func (s *ShippingService) DeleteShipment(ctx context.Context, shipmentID uint) error {
	_, err := s.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return ErrShipmentNotFound
	}

	return s.repo.DeleteShipment(ctx, shipmentID)
}

// GetShipmentByID retrieves a shipment by ID
func (s *ShippingService) GetShipmentByID(ctx context.Context, shipmentID uint) (*model.Shipment, error) {
	shipment, err := s.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, ErrShipmentNotFound
	}
	return shipment, nil
}

// GetShipmentByOrderID retrieves a shipment by order ID
func (s *ShippingService) GetShipmentByOrderID(ctx context.Context, orderID uint) (*model.Shipment, error) {
	shipment, err := s.repo.GetShipmentByOrderID(ctx, orderID)
	if err != nil {
		return nil, ErrShipmentNotFound
	}
	return shipment, nil
}

// GetShipmentsByStatus retrieves shipments by status
func (s *ShippingService) GetShipmentsByStatus(ctx context.Context, status model.ShipmentStatus, page, limit int) ([]model.Shipment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.GetShipmentsByStatus(ctx, status, page, limit)
}

// ==================== TRACKING ====================

// AddTrackingEvent adds a tracking event to a shipment
func (s *ShippingService) AddTrackingEvent(ctx context.Context, shipmentID uint, input model.TrackingEventInput) (*model.ShipmentTracking, error) {
	// Verify shipment exists
	_, err := s.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, ErrShipmentNotFound
	}

	// Parse occurred_at
	occurredAt := time.Now()
	if input.OccurredAt != "" {
		occurredAt, err = time.Parse(time.RFC3339, input.OccurredAt)
		if err != nil {
			return nil, errors.New("invalid date format, use RFC3339")
		}
	}

	// Parse status
	status := model.ShipmentStatus(input.Status)

	event := &model.ShipmentTracking{
		ShipmentID:  shipmentID,
		Status:      status,
		Location:    input.Location,
		Description: input.Description,
		OccurredAt:  occurredAt,
	}

	if err := s.repo.AddTrackingEvent(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

// GetTrackingByShipmentID retrieves all tracking events for a shipment
func (s *ShippingService) GetTrackingByShipmentID(ctx context.Context, shipmentID uint) ([]model.ShipmentTracking, error) {
	return s.repo.GetTrackingByShipmentID(ctx, shipmentID)
}

// GetTrackingTimeline retrieves tracking timeline for a shipment
func (s *ShippingService) GetTrackingTimeline(ctx context.Context, shipmentID uint) (*model.TrackingTimeline, error) {
	return s.repo.GetTrackingTimeline(ctx, shipmentID)
}

// GetOrderTracking retrieves tracking info for an order
func (s *ShippingService) GetOrderTracking(ctx context.Context, orderID uint) (*model.TrackingTimeline, error) {
	shipment, err := s.repo.GetShipmentByOrderID(ctx, orderID)
	if err != nil {
		return nil, ErrShipmentNotFound
	}

	return shipment.ToTrackingTimeline(), nil
}

// ==================== CARRIER ====================

// CreateCarrier creates a new shipping carrier
func (s *ShippingService) CreateCarrier(ctx context.Context, carrier *model.ShippingCarrier) error {
	return s.repo.CreateCarrier(ctx, carrier)
}

// GetActiveCarriers retrieves all active shipping carriers
func (s *ShippingService) GetActiveCarriers(ctx context.Context) ([]model.ShippingCarrier, error) {
	return s.repo.GetActiveCarriers(ctx)
}

// GetAllCarriers retrieves all shipping carriers
func (s *ShippingService) GetAllCarriers(ctx context.Context) ([]model.ShippingCarrier, error) {
	return s.repo.GetAllCarriers(ctx)
}

// GetShipmentStats retrieves shipment statistics
func (s *ShippingService) GetShipmentStats(ctx context.Context) (*repository.ShipmentStats, error) {
	return s.repo.GetShipmentStats(ctx)
}

// ==================== VALIDATION ====================

// validateAddressInput validates shipping address input
func (s *ShippingService) validateAddressInput(input model.ShippingAddressInput) error {
	if input.RecipientName == "" {
		return errors.New("recipient name is required")
	}
	if input.Phone == "" {
		return errors.New("phone number is required")
	}
	if input.AddressLine == "" {
		return errors.New("address line is required")
	}
	if input.District == "" {
		return errors.New("district is required")
	}
	if input.City == "" {
		return errors.New("city is required")
	}
	return nil
}

// validateShipmentInput validates shipment input
func (s *ShippingService) validateShipmentInput(input CreateShipmentInput) error {
	if input.OrderID == 0 {
		return errors.New("order ID is required")
	}
	if input.TrackingNumber == "" {
		return errors.New("tracking number is required")
	}
	if input.CarrierName == "" && input.CarrierID == nil {
		return errors.New("carrier name or carrier ID is required")
	}
	return nil
}

// isValidStatusTransition validates status transitions
func (s *ShippingService) isValidStatusTransition(from, to model.ShipmentStatus) bool {
	validTransitions := map[model.ShipmentStatus][]model.ShipmentStatus{
		model.ShipmentStatusPending: {
			model.ShipmentStatusConfirmed,
			model.ShipmentStatusCancelled,
		},
		model.ShipmentStatusConfirmed: {
			model.ShipmentStatusProcessing,
			model.ShipmentStatusCancelled,
		},
		model.ShipmentStatusProcessing: {
			model.ShipmentStatusPacked,
			model.ShipmentStatusCancelled,
		},
		model.ShipmentStatusPacked: {
			model.ShipmentStatusShipped,
		},
		model.ShipmentStatusShipped: {
			model.ShipmentStatusInTransit,
		},
		model.ShipmentStatusInTransit: {
			model.ShipmentStatusOutForDelivery,
			model.ShipmentStatusFailed,
		},
		model.ShipmentStatusOutForDelivery: {
			model.ShipmentStatusDelivered,
			model.ShipmentStatusFailed,
		},
		model.ShipmentStatusDelivered: {},
		model.ShipmentStatusFailed:    {},
		model.ShipmentStatusCancelled: {},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}

// Error definitions
var (
	ErrAddressNotFound         = errors.New("shipping address not found")
	ErrShipmentNotFound        = errors.New("shipment not found")
	ErrShipmentExists          = errors.New("shipment already exists for this order")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrCarrierNotFound         = errors.New("carrier not found")
)
