package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"context"
	"errors"
)

// InventoryService handles inventory business logic
type InventoryService struct {
	repo *repository.InventoryRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// StockCheckInput represents input for stock check
type StockCheckInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// ReserveStockInput represents input for stock reservation
type ReserveStockInput struct {
	ProductID   uint   `json:"product_id"`
	Quantity    int    `json:"quantity"`
	ReferenceID string `json:"reference_id"` // Order ID
}

// DeductStockInput represents input for stock deduction
type DeductStockInput struct {
	ProductID   uint   `json:"product_id"`
	Quantity    int    `json:"quantity"`
	ReferenceID string `json:"reference_id"` // Order ID
}

// RestockInput represents input for restocking
type RestockInput struct {
	ProductID   uint   `json:"product_id"`
	Quantity    int    `json:"quantity"`
	ReferenceID string `json:"reference_id"` // Restock ID
	UserID      uint   `json:"user_id"`
	Reason      string `json:"reason"`
}

// BatchReserveItem represents an item for batch reservation
type BatchReserveItem struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// BatchReserveInput represents input for batch reservation
type BatchReserveInput struct {
	Items       []BatchReserveItem `json:"items"`
	ReferenceID string             `json:"reference_id"` // Order ID
}

// CheckStock checks if stock is available for a product
func (s *InventoryService) CheckStock(ctx context.Context, productID uint, quantity int) (*model.Inventory, bool, error) {
	inventory, err := s.repo.GetInventoryByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, false, ErrInventoryNotFound
		}
		return nil, false, err
	}
	
	// Check if enough stock
	canFulfill := inventory.CanFulfillQuantity(quantity)
	
	return inventory, canFulfill, nil
}

// ReserveStock reserves stock for an order
func (s *InventoryService) ReserveStock(ctx context.Context, input ReserveStockInput) (*model.Inventory, error) {
	// Validate input
	if input.ProductID == 0 {
		return nil, ErrInvalidProductID
	}
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if input.ReferenceID == "" {
		return nil, ErrInvalidReferenceID
	}
	
	// Reserve stock
	inventory, err := s.repo.ReserveStock(ctx, input.ProductID, input.Quantity, input.ReferenceID)
	if err != nil {
		if err.Error() == "insufficient stock available" {
			return nil, ErrInsufficientStock
		}
		return nil, err
	}
	
	return inventory, nil
}

// ReleaseStock releases reserved stock
func (s *InventoryService) ReleaseStock(ctx context.Context, input ReserveStockInput) (*model.Inventory, error) {
	// Validate input
	if input.ProductID == 0 {
		return nil, ErrInvalidProductID
	}
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if input.ReferenceID == "" {
		return nil, ErrInvalidReferenceID
	}
	
	// Release stock
	inventory, err := s.repo.ReleaseStock(ctx, input.ProductID, input.Quantity, input.ReferenceID)
	if err != nil {
		if err.Error() == "insufficient reserved stock" {
			return nil, ErrInsufficientReservedStock
		}
		return nil, err
	}
	
	return inventory, nil
}

// DeductStock deducts stock after order confirmation
func (s *InventoryService) DeductStock(ctx context.Context, input DeductStockInput) (*model.Inventory, error) {
	// Validate input
	if input.ProductID == 0 {
		return nil, ErrInvalidProductID
	}
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if input.ReferenceID == "" {
		return nil, ErrInvalidReferenceID
	}
	
	// Deduct stock
	inventory, err := s.repo.DeductStock(ctx, input.ProductID, input.Quantity, input.ReferenceID)
	if err != nil {
		if err.Error() == "insufficient reserved stock to deduct" {
			return nil, ErrInsufficientReservedStock
		}
		return nil, err
	}
	
	return inventory, nil
}

// RestockProduct adds stock to a product
func (s *InventoryService) RestockProduct(ctx context.Context, input RestockInput) (*model.Inventory, error) {
	// Validate input
	if input.ProductID == 0 {
		return nil, ErrInvalidProductID
	}
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	
	// Restock
	inventory, err := s.repo.RestockProduct(ctx, input.ProductID, input.Quantity, input.ReferenceID, &input.UserID, input.Reason)
	if err != nil {
		return nil, err
	}
	
	return inventory, nil
}

// BatchReserveStock reserves stock for multiple products
func (s *InventoryService) BatchReserveStock(ctx context.Context, input BatchReserveInput) error {
	// Validate input
	if input.ReferenceID == "" {
		return ErrInvalidReferenceID
	}
	if len(input.Items) == 0 {
		return ErrInvalidQuantity
	}
	
	// Validate all items
	for _, item := range input.Items {
		if item.ProductID == 0 {
			return ErrInvalidProductID
		}
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}
	}
	
	// Convert to repository format
	repoItems := make([]struct {
		ProductID uint
		Quantity  int
	}, len(input.Items))
	
	for i, item := range input.Items {
		repoItems[i].ProductID = item.ProductID
		repoItems[i].Quantity = item.Quantity
	}
	
	// Batch reserve
	return s.repo.BatchReserveStock(ctx, repoItems, input.ReferenceID)
}

// ReturnStock returns stock from cancelled/returned order
func (s *InventoryService) ReturnStock(ctx context.Context, productID uint, quantity int, referenceID string, reason string) (*model.Inventory, error) {
	// Validate input
	if productID == 0 {
		return nil, ErrInvalidProductID
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	
	// Return stock
	inventory, err := s.repo.ReturnStock(ctx, productID, quantity, referenceID, reason)
	if err != nil {
		return nil, err
	}
	
	return inventory, nil
}

// GetInventoryByProductID retrieves inventory for a product
func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID uint) (*model.Inventory, error) {
	inventory, err := s.repo.GetInventoryByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, ErrInventoryNotFound
		}
		return nil, err
	}
	
	return inventory, nil
}

// GetInventoryLogs retrieves inventory logs for a product
func (s *InventoryService) GetInventoryLogs(ctx context.Context, productID uint, page, limit int) ([]model.InventoryLog, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetInventoryLogs(ctx, productID, limit, offset)
}

// GetLowStockProducts retrieves products with low stock
func (s *InventoryService) GetLowStockProducts(ctx context.Context, threshold int) ([]model.Inventory, error) {
	return s.repo.GetLowStockProducts(ctx, threshold)
}

// GetOutOfStockProducts retrieves products that are out of stock
func (s *InventoryService) GetOutOfStockProducts(ctx context.Context) ([]model.Inventory, error) {
	return s.repo.GetOutOfStockProducts(ctx)
}

// GetInventorySummary retrieves inventory statistics
func (s *InventoryService) GetInventorySummary(ctx context.Context) (*model.InventorySummary, error) {
	return s.repo.GetInventorySummary(ctx)
}

// CheckAndCreateAlerts checks stock levels and creates alerts
func (s *InventoryService) CheckAndCreateAlerts(ctx context.Context) error {
	return s.repo.CheckAndCreateAlerts(ctx)
}

// ProcessOrderStock handles stock operations for order lifecycle
func (s *InventoryService) ProcessOrderStock(ctx context.Context, orderID string, items []ReserveStockInput, action string) error {
	switch action {
	case "reserve":
		// Reserve stock when order is created
		for _, item := range items {
			item.ReferenceID = orderID
			_, err := s.ReserveStock(ctx, item)
			if err != nil {
				return err
			}
		}
	case "deduct":
		// Deduct stock when payment is confirmed
		for _, item := range items {
			_, err := s.DeductStock(ctx, DeductStockInput{
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				ReferenceID: orderID,
			})
			if err != nil {
				return err
			}
		}
	case "release":
		// Release stock when order is cancelled
		for _, item := range items {
			item.ReferenceID = orderID
			_, err := s.ReleaseStock(ctx, item)
			if err != nil {
				return err
			}
		}
	default:
		return ErrInvalidAction
	}
	
	return nil
}

// Error definitions - only define errors not already in errors.go
var (
	ErrInventoryNotFound         = errors.New("inventory not found")
	ErrInsufficientReservedStock = errors.New("insufficient reserved stock")
	ErrInvalidReferenceID        = errors.New("invalid reference ID")
	ErrInvalidAction             = errors.New("invalid action")
	ErrRecordNotFound            = errors.New("record not found")
)
