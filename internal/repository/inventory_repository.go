package repository

import (
	"ecommerce/internal/domain/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

// InventoryRepository handles database operations for inventory
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// GetInventoryByProductID retrieves inventory for a product
func (r *InventoryRepository) GetInventoryByProductID(ctx context.Context, productID uint) (*model.Inventory, error) {
	var inventory model.Inventory
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		First(&inventory).Error
	if err != nil {
		return nil, err
	}
	
	// Calculate available quantity
	inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
	
	return &inventory, nil
}

// GetOrCreateInventory retrieves or creates inventory for a product
func (r *InventoryRepository) GetOrCreateInventory(ctx context.Context, productID uint, initialStock int) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		FirstOrCreate(&inventory, model.Inventory{
			ProductID:         productID,
			StockQuantity:     initialStock,
			ReservedQuantity:  0,
			AvailableQuantity: initialStock,
		}).Error
	
	if err != nil {
		return nil, err
	}
	
	inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
	return &inventory, nil
}

// UpdateStock updates stock quantity atomically
func (r *InventoryRepository) UpdateStock(ctx context.Context, productID uint, quantity int) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			Where("product_id = ?", productID).
			First(&inventory).Error; err != nil {
			return err
		}
		
		// Calculate new stock
		newStock := inventory.StockQuantity + quantity
		if newStock < 0 {
			return errors.New("stock cannot be negative")
		}

		// Update stock
		inventory.StockQuantity = newStock
		inventory.AvailableQuantity = newStock - inventory.ReservedQuantity

		return tx.Save(&inventory).Error
	})
	
	if err != nil {
		return nil, err
	}
	
	return &inventory, nil
}

// ReserveStock reserves stock for an order (atomic operation)
func (r *InventoryRepository) ReserveStock(ctx context.Context, productID uint, quantity int, referenceID string) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update to prevent race conditions
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			Where("product_id = ?", productID).
			First(&inventory).Error; err != nil {
			return err
		}
		
		// Calculate available stock
		availableStock := inventory.StockQuantity - inventory.ReservedQuantity
		
		// Check if enough stock is available
		if availableStock < quantity {
			return errors.New("insufficient stock available")
		}
		
		// Reserve the stock
		oldReserved := inventory.ReservedQuantity
		inventory.ReservedQuantity += quantity
		inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
		
		// Create inventory log
		log := model.InventoryLog{
			ProductID:      productID,
			InventoryID:    inventory.ID,
			ChangeType:     model.InventoryChangeReserve,
			Quantity:       quantity,
			StockBefore:    inventory.StockQuantity,
			StockAfter:     inventory.StockQuantity,
			ReservedBefore: oldReserved,
			ReservedAfter:  inventory.ReservedQuantity,
			ReferenceType:  "order",
			ReferenceID:    referenceID,
		}
		
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &inventory, nil
}

// ReleaseStock releases reserved stock (when order is cancelled)
func (r *InventoryRepository) ReleaseStock(ctx context.Context, productID uint, quantity int, referenceID string) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			Where("product_id = ?", productID).
			First(&inventory).Error; err != nil {
			return err
		}
		
		// Check if enough reserved stock
		if inventory.ReservedQuantity < quantity {
			return errors.New("insufficient reserved stock")
		}
		
		// Release the stock
		oldReserved := inventory.ReservedQuantity
		inventory.ReservedQuantity -= quantity
		inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
		
		// Create inventory log
		log := model.InventoryLog{
			ProductID:      productID,
			InventoryID:    inventory.ID,
			ChangeType:     model.InventoryChangeRelease,
			Quantity:       -quantity,
			StockBefore:    inventory.StockQuantity,
			StockAfter:     inventory.StockQuantity,
			ReservedBefore: oldReserved,
			ReservedAfter:  inventory.ReservedQuantity,
			ReferenceType:  "order",
			ReferenceID:    referenceID,
		}
		
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &inventory, nil
}

// DeductStock deducts stock after order confirmation (atomic operation)
func (r *InventoryRepository) DeductStock(ctx context.Context, productID uint, quantity int, referenceID string) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			Where("product_id = ?", productID).
			First(&inventory).Error; err != nil {
			return err
		}
		
		// Check if enough reserved stock
		if inventory.ReservedQuantity < quantity {
			return errors.New("insufficient reserved stock to deduct")
		}
		
		// Deduct from both stock and reserved
		oldStock := inventory.StockQuantity
		oldReserved := inventory.ReservedQuantity
		
		inventory.StockQuantity -= quantity
		inventory.ReservedQuantity -= quantity
		inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
		
		// Create inventory log
		log := model.InventoryLog{
			ProductID:      productID,
			InventoryID:    inventory.ID,
			ChangeType:     model.InventoryChangeDeduct,
			Quantity:       -quantity,
			StockBefore:    oldStock,
			StockAfter:     inventory.StockQuantity,
			ReservedBefore: oldReserved,
			ReservedAfter:  inventory.ReservedQuantity,
			ReferenceType:  "order",
			ReferenceID:    referenceID,
		}
		
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &inventory, nil
}

// RestockProduct adds stock to a product
func (r *InventoryRepository) RestockProduct(ctx context.Context, productID uint, quantity int, referenceID string, userID *uint, reason string) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			Where("product_id = ?", productID).
			First(&inventory).Error; err != nil {
			// If not found, create new inventory
			inventory = model.Inventory{
				ProductID:         productID,
				StockQuantity:     0,
				ReservedQuantity:  0,
				AvailableQuantity: 0,
			}
			if err := tx.Create(&inventory).Error; err != nil {
				return err
			}
		}
		
		// Add stock
		oldStock := inventory.StockQuantity
		oldReserved := inventory.ReservedQuantity
		
		inventory.StockQuantity += quantity
		inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
		
		// Create inventory log
		log := model.InventoryLog{
			ProductID:      productID,
			InventoryID:    inventory.ID,
			ChangeType:     model.InventoryChangeRestock,
			Quantity:       quantity,
			StockBefore:    oldStock,
			StockAfter:     inventory.StockQuantity,
			ReservedBefore: oldReserved,
			ReservedAfter:  inventory.ReservedQuantity,
			ReferenceType:  "restock",
			ReferenceID:    referenceID,
			UserID:         userID,
			Reason:         reason,
		}
		
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &inventory, nil
}

// GetInventoryLogs retrieves inventory logs for a product
func (r *InventoryRepository) GetInventoryLogs(ctx context.Context, productID uint, limit, offset int) ([]model.InventoryLog, int64, error) {
	var logs []model.InventoryLog
	var total int64
	
	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&model.InventoryLog{}).
		Where("product_id = ?", productID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get logs with pagination
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	
	return logs, total, err
}

// GetLowStockProducts retrieves products with low stock
func (r *InventoryRepository) GetLowStockProducts(ctx context.Context, threshold int) ([]model.Inventory, error) {
	var inventories []model.Inventory
	err := r.db.WithContext(ctx).
		Where("stock_quantity - reserved_quantity <= ? AND stock_quantity - reserved_quantity > 0", threshold).
		Preload("Product").
		Find(&inventories).Error
	return inventories, err
}

// GetOutOfStockProducts retrieves products that are out of stock
func (r *InventoryRepository) GetOutOfStockProducts(ctx context.Context) ([]model.Inventory, error) {
	var inventories []model.Inventory
	err := r.db.WithContext(ctx).
		Where("stock_quantity - reserved_quantity <= 0").
		Preload("Product").
		Find(&inventories).Error
	return inventories, err
}

// GetInventorySummary retrieves inventory statistics
func (r *InventoryRepository) GetInventorySummary(ctx context.Context) (*model.InventorySummary, error) {
	summary := &model.InventorySummary{}
	
	// Total products
	if err := r.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Count(&summary.TotalProducts).Error; err != nil {
		return nil, err
	}
	
	// In stock products
	if err := r.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Where("stock_quantity - reserved_quantity > 0").
		Count(&summary.InStockProducts).Error; err != nil {
		return nil, err
	}
	
	// Low stock products (less than 10)
	if err := r.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Where("stock_quantity - reserved_quantity > 0 AND stock_quantity - reserved_quantity <= 10").
		Count(&summary.LowStockProducts).Error; err != nil {
		return nil, err
	}
	
	// Out of stock products
	if err := r.db.WithContext(ctx).
		Model(&model.Inventory{}).
		Where("stock_quantity - reserved_quantity <= 0").
		Count(&summary.OutOfStockProducts).Error; err != nil {
		return nil, err
	}
	
	return summary, nil
}

// CheckAndCreateAlerts checks stock levels and creates alerts
func (r *InventoryRepository) CheckAndCreateAlerts(ctx context.Context) error {
	var inventories []model.Inventory
	
	// Get all inventories
	if err := r.db.WithContext(ctx).
		Where("stock_quantity - reserved_quantity <= reorder_level AND reorder_level > 0").
		Find(&inventories).Error; err != nil {
		return err
	}
	
	// Create alerts for each
	for _, inv := range inventories {
		alert := model.StockAlert{
			ProductID:    inv.ProductID,
			AlertType:    "low_stock",
			Threshold:    inv.ReorderLevel,
			CurrentStock: inv.AvailableQuantity,
		}
		
		// Check if alert already exists
		var existing model.StockAlert
		if err := r.db.WithContext(ctx).
			Where("product_id = ? AND is_resolved = ?", inv.ProductID, false).
			First(&existing).Error; err == gorm.ErrRecordNotFound {
			// Create new alert
			r.db.WithContext(ctx).Create(&alert)
		}
	}
	
	return nil
}

// BatchReserveStock reserves stock for multiple products (for cart checkout)
func (r *InventoryRepository) BatchReserveStock(ctx context.Context, items []struct {
	ProductID uint
	Quantity  int
}, referenceID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			var inventory model.Inventory
			
			// Lock and get inventory
			if err := tx.
				Set("gorm:query_option", "WITH (UPDLOCK)").
				Where("product_id = ?", item.ProductID).
				First(&inventory).Error; err != nil {
				return err
			}
			
			// Check available stock
			availableStock := inventory.StockQuantity - inventory.ReservedQuantity
			if availableStock < item.Quantity {
				return errors.New("insufficient stock for product")
			}
			
			// Reserve stock
			oldReserved := inventory.ReservedQuantity
			inventory.ReservedQuantity += item.Quantity
			inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
			
			if err := tx.Save(&inventory).Error; err != nil {
				return err
			}
			
			// Create log
			log := model.InventoryLog{
				ProductID:      item.ProductID,
				InventoryID:    inventory.ID,
				ChangeType:     model.InventoryChangeReserve,
				Quantity:       item.Quantity,
				StockBefore:    inventory.StockQuantity,
				StockAfter:     inventory.StockQuantity,
				ReservedBefore: oldReserved,
				ReservedAfter:  inventory.ReservedQuantity,
				ReferenceType:  "order",
				ReferenceID:    referenceID,
			}
			
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// ReturnStock returns stock from cancelled/returned order
func (r *InventoryRepository) ReturnStock(ctx context.Context, productID uint, quantity int, referenceID string, reason string) (*model.Inventory, error) {
	var inventory model.Inventory
	
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.
			Set("gorm:query_option", "WITH (UPDLOCK)").
			Where("product_id = ?", productID).
			First(&inventory).Error; err != nil {
			return err
		}
		
		// Add stock back
		oldStock := inventory.StockQuantity
		oldReserved := inventory.ReservedQuantity
		
		inventory.StockQuantity += quantity
		inventory.AvailableQuantity = inventory.StockQuantity - inventory.ReservedQuantity
		
		// Create inventory log
		log := model.InventoryLog{
			ProductID:      productID,
			InventoryID:    inventory.ID,
			ChangeType:     model.InventoryChangeReturn,
			Quantity:       quantity,
			StockBefore:    oldStock,
			StockAfter:     inventory.StockQuantity,
			ReservedBefore: oldReserved,
			ReservedAfter:  inventory.ReservedQuantity,
			ReferenceType:  "return",
			ReferenceID:    referenceID,
			Reason:         reason,
		}
		
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &inventory, nil
}

// UpdateInventorySync syncs inventory with product stock fields
func (r *InventoryRepository) UpdateInventorySync(ctx context.Context, productID uint) error {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return err
	}
	
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory model.Inventory
		if err := tx.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
			return err
		}
		
		// Sync values
		inventory.StockQuantity = product.Stock
		inventory.ReservedQuantity = product.ReservedStock
		inventory.AvailableQuantity = product.Stock - product.ReservedStock
		
		return tx.Save(&inventory).Error
	})
}
