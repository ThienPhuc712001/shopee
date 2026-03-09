package repository

import (
	"ecommerce/internal/domain/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// CartRepositoryEnhanced defines the enhanced interface for cart data operations
type CartRepositoryEnhanced interface {
	// Cart Management
	GetCartByUserID(userID uint) (*model.Cart, error)
	CreateCart(userID uint) (*model.Cart, error)
	GetOrCreateCart(userID uint) (*model.Cart, error)
	UpdateCart(cart *model.Cart) error
	DeleteCart(cartID uint) error

	// Cart Item Management
	AddItem(item *model.CartItem) error
	UpdateItem(item *model.CartItem) error
	RemoveItem(itemID uint) error
	FindItemByID(itemID uint) (*model.CartItem, error)
	FindItemByCartAndProduct(cartID, productID uint, variantID *uint) (*model.CartItem, error)
	GetItemsByCartID(cartID uint) ([]model.CartItem, error)
	ClearCart(cartID uint) error

	// Cart Calculations
	UpdateCartTotals(cartID uint) error
	GetCartTotal(cartID uint) (float64, error)
	GetCartItemCount(cartID uint) (int, error)

	// Bulk Operations
	BulkAddItems(items []*model.CartItem) error
	BulkUpdateItems(items []*model.CartItem) error
	BulkRemoveItems(itemIDs []uint) error

	// Stock Management
	ReserveCartItems(cartID uint) error
	ReleaseCartItems(cartID uint) error

	// Cart Validation
	ValidateCartItems(cartID uint) ([]model.CartItem, []error, error)
	GetInvalidCartItems(cartID uint) ([]model.CartItem, error)

	// Analytics
	GetActiveCartsCount() (int64, error)
	GetCartAbandonmentRate() (float64, error)
}

type cartRepositoryEnhanced struct {
	db *gorm.DB
}

// NewCartRepositoryEnhanced creates a new enhanced cart repository
func NewCartRepositoryEnhanced(db *gorm.DB) CartRepositoryEnhanced {
	return &cartRepositoryEnhanced{db: db}
}

// ==================== CART MANAGEMENT ====================

func (r *cartRepositoryEnhanced) GetCartByUserID(userID uint) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("user_id = ?", userID).
		Preload("Items.Product").
		Preload("Items.Variant").
		Preload("Items.Shop").
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepositoryEnhanced) CreateCart(userID uint) (*model.Cart, error) {
	cart := &model.Cart{
		UserID:     userID,
		TotalItems: 0,
		Subtotal:   0,
		Total:      0,
		Currency:   "USD",
	}

	err := r.db.Create(cart).Error
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *cartRepositoryEnhanced) GetOrCreateCart(userID uint) (*model.Cart, error) {
	var cart model.Cart

	// Try to find existing cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error

	if err == gorm.ErrRecordNotFound {
		// Create new cart if not exists
		cart = model.Cart{
			UserID:     userID,
			TotalItems: 0,
			Subtotal:   0,
			Total:      0,
			Currency:   "USD",
		}

		if err := r.db.Create(&cart).Error; err != nil {
			return nil, err
		}

		return &cart, nil
	}

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryEnhanced) UpdateCart(cart *model.Cart) error {
	return r.db.Save(cart).Error
}

func (r *cartRepositoryEnhanced) DeleteCart(cartID uint) error {
	return r.db.Delete(&model.Cart{}, cartID).Error
}

// ==================== CART ITEM MANAGEMENT ====================

func (r *cartRepositoryEnhanced) AddItem(item *model.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepositoryEnhanced) UpdateItem(item *model.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepositoryEnhanced) RemoveItem(itemID uint) error {
	// Soft delete
	return r.db.Delete(&model.CartItem{}, itemID).Error
}

func (r *cartRepositoryEnhanced) FindItemByID(itemID uint) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Preload("Product").Preload("Variant").Preload("Shop").
		First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepositoryEnhanced) FindItemByCartAndProduct(cartID, productID uint, variantID *uint) (*model.CartItem, error) {
	var item model.CartItem

	query := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.First(&item).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *cartRepositoryEnhanced) GetItemsByCartID(cartID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	err := r.db.Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Preload("Product").
		Preload("Variant").
		Preload("Shop").
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *cartRepositoryEnhanced) ClearCart(cartID uint) error {
	// Soft delete all items in cart
	return r.db.Model(&model.CartItem{}).
		Where("cart_id = ?", cartID).
		Update("deleted_at", gorm.Expr("GETDATE()")).Error
}

// ==================== CART CALCULATIONS ====================

func (r *cartRepositoryEnhanced) UpdateCartTotals(cartID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Calculate total items and subtotal
	var totalItems int
	var subtotal float64

	err := tx.Model(&model.CartItem{}).
		Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalItems).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.CartItem{}).
		Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&subtotal).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update cart
	err = tx.Model(&model.Cart{}).
		Where("id = ?", cartID).
		Updates(map[string]interface{}{
			"total_items":  totalItems,
			"subtotal":     subtotal,
			"total":        subtotal, // Discount would be applied separately
			"last_activity": gorm.Expr("GETDATE()"),
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *cartRepositoryEnhanced) GetCartTotal(cartID uint) (float64, error) {
	var total float64
	err := r.db.Model(&model.CartItem{}).
		Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&total).Error
	return total, err
}

func (r *cartRepositoryEnhanced) GetCartItemCount(cartID uint) (int, error) {
	var count int
	err := r.db.Model(&model.CartItem{}).
		Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&count).Error
	return count, err
}

// ==================== BULK OPERATIONS ====================

func (r *cartRepositoryEnhanced) BulkAddItems(items []*model.CartItem) error {
	return r.db.CreateInBatches(items, 50).Error
}

func (r *cartRepositoryEnhanced) BulkUpdateItems(items []*model.CartItem) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	for _, item := range items {
		if err := tx.Save(item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *cartRepositoryEnhanced) BulkRemoveItems(itemIDs []uint) error {
	return r.db.Model(&model.CartItem{}).
		Where("id IN ?", itemIDs).
		Update("deleted_at", gorm.Expr("GETDATE()")).Error
}

// ==================== CART VALIDATION ====================

func (r *cartRepositoryEnhanced) ValidateCartItems(cartID uint) ([]model.CartItem, []error, error) {
	items, err := r.GetItemsByCartID(cartID)
	if err != nil {
		return nil, nil, err
	}

	var invalidItems []model.CartItem
	var validationErrors []error

	for _, item := range items {
		// Check if product exists and is active
		if item.Product == nil {
			invalidItems = append(invalidItems, item)
			validationErrors = append(validationErrors, fmt.Errorf("product %d not found", item.ProductID))
			continue
		}

		// Check product status
		if item.Product.Status != "active" {
			invalidItems = append(invalidItems, item)
			validationErrors = append(validationErrors, fmt.Errorf("product %s is not available", item.Product.Name))
			continue
		}

		// Check stock availability
		availableStock := item.Product.Stock
		if item.VariantID != nil && item.Variant != nil {
			availableStock = item.Variant.Stock
		}

		if availableStock < item.Quantity {
			invalidItems = append(invalidItems, item)
			validationErrors = append(validationErrors, fmt.Errorf("insufficient stock for %s", item.Product.Name))
		}
	}

	return items, validationErrors, nil
}

func (r *cartRepositoryEnhanced) GetInvalidCartItems(cartID uint) ([]model.CartItem, error) {
	items, errors, err := r.ValidateCartItems(cartID)
	if err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		return items, nil // Return items with validation errors
	}

	return nil, nil
}

// ==================== ANALYTICS ====================

func (r *cartRepositoryEnhanced) GetActiveCartsCount() (int64, error) {
	var count int64
	// Carts with activity in the last 30 days
	err := r.db.Model(&model.Cart{}).
		Where("last_activity >= DATEADD(day, -30, GETDATE())").
		Count(&count).Error
	return count, err
}

func (r *cartRepositoryEnhanced) GetCartAbandonmentRate() (float64, error) {
	// This is a simplified calculation
	// In production, you'd track cart creation vs order completion

	var totalCarts int64
	var cartsWithOrders int64

	err := r.db.Model(&model.Cart{}).
		Where("created_at >= DATEADD(day, -30, GETDATE())").
		Count(&totalCarts).Error
	if err != nil {
		return 0, err
	}

	// Count carts that resulted in orders (simplified)
	err = r.db.Model(&model.Cart{}).
		Where("created_at >= DATEADD(day, -30, GETDATE())").
		Where("total_items > 0").
		Count(&cartsWithOrders).Error
	if err != nil {
		return 0, err
	}

	if totalCarts == 0 {
		return 0, nil
	}

	// Abandonment rate = (1 - conversion rate) * 100
	conversionRate := float64(cartsWithOrders) / float64(totalCarts)
	abandonmentRate := (1 - conversionRate) * 100

	return abandonmentRate, nil
}

// ==================== HELPER METHODS ====================

// FindCartWithItems finds a cart with all items loaded
func (r *cartRepositoryEnhanced) FindCartWithItems(cartID uint) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("id = ?", cartID).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").
				Preload("Product").
				Preload("Variant").
				Preload("Shop")
		}).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// ItemExists checks if a cart item exists
func (r *cartRepositoryEnhanced) ItemExists(cartID, productID uint, variantID *uint) (bool, error) {
	var count int64

	query := r.db.Model(&model.CartItem{}).
		Where("cart_id = ? AND product_id = ? AND deleted_at IS NULL", cartID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetCartItemsByShop groups cart items by shop
func (r *cartRepositoryEnhanced) GetCartItemsByShop(cartID uint) (map[uint][]model.CartItem, error) {
	items, err := r.GetItemsByCartID(cartID)
	if err != nil {
		return nil, err
	}

	shopItems := make(map[uint][]model.CartItem)
	for _, item := range items {
		shopItems[item.ShopID] = append(shopItems[item.ShopID], item)
	}

	return shopItems, nil
}

// ReserveCartItems reserves stock for all cart items
func (r *cartRepositoryEnhanced) ReserveCartItems(cartID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	items, err := r.GetItemsByCartID(cartID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range items {
		// Reserve stock for product or variant
		if item.VariantID != nil {
			err = tx.Model(&model.ProductVariant{}).
				Where("id = ?", *item.VariantID).
				UpdateColumn("reserved_stock", gorm.Expr("reserved_stock + ?", item.Quantity)).Error
		} else {
			err = tx.Model(&model.Product{}).
				Where("id = ?", item.ProductID).
				UpdateColumn("reserved_stock", gorm.Expr("reserved_stock + ?", item.Quantity)).Error
		}

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// ReleaseCartItems releases reserved stock for all cart items
func (r *cartRepositoryEnhanced) ReleaseCartItems(cartID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	items, err := r.GetItemsByCartID(cartID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range items {
		// Release stock for product or variant
		if item.VariantID != nil {
			err = tx.Model(&model.ProductVariant{}).
				Where("id = ?", *item.VariantID).
				UpdateColumn("reserved_stock", gorm.Expr("reserved_stock - ?", item.Quantity)).Error
		} else {
			err = tx.Model(&model.Product{}).
				Where("id = ?", item.ProductID).
				UpdateColumn("reserved_stock", gorm.Expr("reserved_stock - ?", item.Quantity)).Error
		}

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Error definitions
var (
	ErrCartNotFound      = errors.New("cart not found")
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrItemAlreadyExists = errors.New("item already exists in cart")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInsufficientStock = errors.New("insufficient stock")
)
