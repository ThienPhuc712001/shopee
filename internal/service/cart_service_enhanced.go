package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"fmt"
)

// CartServiceEnhanced defines the enhanced cart service interface
type CartServiceEnhanced interface {
	// Cart Management
	GetCart(userID uint) (*model.Cart, error)
	GetOrCreateCart(userID uint) (*model.Cart, error)
	ClearCart(userID uint) error

	// Cart Item Operations
	AddToCart(userID, productID uint, variantID *uint, quantity int) (*model.Cart, error)
	UpdateCartItem(userID, cartItemID uint, quantity int) (*model.Cart, error)
	RemoveFromCart(userID, cartItemID uint) (*model.Cart, error)
	RemoveItemByProduct(userID, productID uint, variantID *uint) error

	// Cart Calculations
	CalculateCartTotal(cartID uint) (float64, error)
	GetCartSummary(userID uint) (*model.CartSummary, error)

	// Cart Validation
	ValidateCart(userID uint) ([]model.CartItem, []error, error)
	ValidateCartForCheckout(userID uint) error

	// Checkout Preparation
	PrepareForCheckout(userID uint) (*model.CartCheckoutSummary, error)
	ReserveCartStock(userID uint) error
	ReleaseCartStock(userID uint) error

	// Cart Statistics
	GetCartStats(userID uint) (*CartStats, error)
}

// CartStats represents cart statistics
type CartStats struct {
	TotalItems      int     `json:"total_items"`
	TotalProducts   int     `json:"total_products"`
	Subtotal        float64 `json:"subtotal"`
	EstimatedTotal  float64 `json:"estimated_total"`
	HasOutOfStock   bool    `json:"has_out_of_stock"`
	HasPriceChanged bool    `json:"has_price_changed"`
}

type cartServiceEnhanced struct {
	cartRepo    repository.CartRepositoryEnhanced
	productRepo repository.ProductRepositoryEnhanced
}

// NewCartServiceEnhanced creates a new enhanced cart service
func NewCartServiceEnhanced(
	cartRepo repository.CartRepositoryEnhanced,
	productRepo repository.ProductRepositoryEnhanced,
) CartServiceEnhanced {
	return &cartServiceEnhanced{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// ==================== CART MANAGEMENT ====================

func (s *cartServiceEnhanced) GetCart(userID uint) (*model.Cart, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Recalculate totals to ensure accuracy
	cart.RecalculateTotals()

	return cart, nil
}

func (s *cartServiceEnhanced) GetOrCreateCart(userID uint) (*model.Cart, error) {
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	return cart, nil
}

func (s *cartServiceEnhanced) ClearCart(userID uint) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return ErrCartNotFound
	}

	// Release any reserved stock
	_ = s.cartRepo.ReleaseCartItems(cart.ID)

	// Clear all items
	return s.cartRepo.ClearCart(cart.ID)
}

// ==================== CART ITEM OPERATIONS ====================

func (s *cartServiceEnhanced) AddToCart(userID, productID uint, variantID *uint, quantity int) (*model.Cart, error) {
	// Validate quantity
	if quantity < 1 || quantity > 999 {
		return nil, ErrInvalidQuantity
	}

	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Validate product exists and is available
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.Status != model.ProductStatusActive {
		return nil, ErrProductNotAvailable
	}

	// Check stock availability
	availableStock := product.Stock
	currentPrice := product.Price

	if variantID != nil {
		variant, err := s.productRepo.FindVariantBySKU("")
		if err != nil {
			// Try to find variant by ID
			variants, _ := s.productRepo.FindVariantsByProductID(productID)
			for _, v := range variants {
				if uint(v.ID) == *variantID {
					variant = &v
					break
				}
			}
		}

		if variant == nil {
			return nil, ErrVariantNotFound
		}

		availableStock = variant.Stock
		if variant.Price > 0 {
			currentPrice = variant.Price
		}
	}

	// Check if sufficient stock
	if availableStock < quantity {
		return nil, ErrInsufficientStock
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.FindItemByCartAndProduct(cart.ID, productID, variantID)
	if err == nil && existingItem != nil {
		// Item exists - update quantity
		newQuantity := existingItem.Quantity + quantity
		if newQuantity > availableStock {
			return nil, ErrInsufficientStock
		}

		existingItem.Quantity = newQuantity
		if err := s.cartRepo.UpdateItem(existingItem); err != nil {
			return nil, err
		}
	} else {
		// Create new cart item
		productImage := ""
		if len(product.Images) > 0 {
			productImage = product.Images[0].URL
		}

		cartItem := &model.CartItem{
			CartID:       cart.ID,
			ProductID:    productID,
			VariantID:    variantID,
			Quantity:     quantity,
			Price:        currentPrice,
			ProductName:  product.Name,
			ProductImage: productImage,
			ShopID:       product.ShopID,
		}

		if err := s.cartRepo.AddItem(cartItem); err != nil {
			return nil, err
		}
	}

	// Update cart totals
	if err := s.cartRepo.UpdateCartTotals(cart.ID); err != nil {
		return nil, err
	}

	// Reload cart with items
	return s.cartRepo.GetCartByUserID(userID)
}

func (s *cartServiceEnhanced) UpdateCartItem(userID, cartItemID uint, quantity int) (*model.Cart, error) {
	// Validate quantity
	if quantity < 1 || quantity > 999 {
		return nil, ErrInvalidQuantity
	}

	// Get cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Get cart item
	item, err := s.cartRepo.FindItemByID(cartItemID)
	if err != nil {
		return nil, ErrCartItemNotFound
	}

	// Verify item belongs to user's cart
	if item.CartID != cart.ID {
		return nil, ErrCartItemNotFound
	}

	// Check stock availability
	availableStock := item.Product.Stock
	if item.VariantID != nil && item.Variant != nil {
		availableStock = item.Variant.Stock
	}

	if quantity > availableStock {
		return nil, ErrInsufficientStock
	}

	// Update quantity
	item.Quantity = quantity
	if err := s.cartRepo.UpdateItem(item); err != nil {
		return nil, err
	}

	// Update cart totals
	if err := s.cartRepo.UpdateCartTotals(cart.ID); err != nil {
		return nil, err
	}

	// Reload cart
	return s.cartRepo.GetCartByUserID(userID)
}

func (s *cartServiceEnhanced) RemoveFromCart(userID, cartItemID uint) (*model.Cart, error) {
	// Get cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	// Get cart item
	item, err := s.cartRepo.FindItemByID(cartItemID)
	if err != nil {
		return nil, ErrCartItemNotFound
	}

	// Verify item belongs to user's cart
	if item.CartID != cart.ID {
		return nil, ErrCartItemNotFound
	}

	// Remove item
	if err := s.cartRepo.RemoveItem(cartItemID); err != nil {
		return nil, err
	}

	// Update cart totals
	if err := s.cartRepo.UpdateCartTotals(cart.ID); err != nil {
		return nil, err
	}

	// Reload cart
	return s.cartRepo.GetCartByUserID(userID)
}

func (s *cartServiceEnhanced) RemoveItemByProduct(userID, productID uint, variantID *uint) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return ErrCartNotFound
	}

	item, err := s.cartRepo.FindItemByCartAndProduct(cart.ID, productID, variantID)
	if err != nil {
		return ErrCartItemNotFound
	}

	return s.cartRepo.RemoveItem(item.ID)
}

// ==================== CART CALCULATIONS ====================

func (s *cartServiceEnhanced) CalculateCartTotal(cartID uint) (float64, error) {
	return s.cartRepo.GetCartTotal(cartID)
}

func (s *cartServiceEnhanced) GetCartSummary(userID uint) (*model.CartSummary, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	return &model.CartSummary{
		TotalItems: cart.TotalItems,
		Subtotal:   cart.Subtotal,
		Discount:   cart.Discount,
		Total:      cart.Total,
		Currency:   cart.Currency,
	}, nil
}

// ==================== CART VALIDATION ====================

func (s *cartServiceEnhanced) ValidateCart(userID uint) ([]model.CartItem, []error, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, nil, ErrCartNotFound
	}

	return s.cartRepo.ValidateCartItems(cart.ID)
}

func (s *cartServiceEnhanced) ValidateCartForCheckout(userID uint) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return ErrCartNotFound
	}

	// Check if cart is empty
	if len(cart.Items) == 0 {
		return ErrCartEmpty
	}

	// Validate all items
	_, errors, err := s.cartRepo.ValidateCartItems(cart.ID)
	if err != nil {
		return err
	}

	if len(errors) > 0 {
		return fmt.Errorf("cart validation failed: %v", errors)
	}

	return nil
}

// ==================== CHECKOUT PREPARATION ====================

func (s *cartServiceEnhanced) PrepareForCheckout(userID uint) (*model.CartCheckoutSummary, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	if len(cart.Items) == 0 {
		return nil, ErrCartEmpty
	}

	// Validate cart
	_, validationErrors, err := s.cartRepo.ValidateCartItems(cart.ID)
	if err != nil {
		return nil, err
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("cart has invalid items: %v", validationErrors)
	}

	// Build checkout summary
	var checkoutItems []model.CartCheckoutItem
	var subtotal float64

	for _, item := range cart.Items {
		availableStock := item.Product.Stock
		if item.VariantID != nil && item.Variant != nil {
			availableStock = item.Variant.Stock
		}

		checkoutItem := model.CartCheckoutItem{
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			Quantity:     item.Quantity,
			Price:        item.Price,
			Subtotal:     item.Subtotal,
			ShopID:       item.ShopID,
			ShopName:     item.Shop.Name,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			IsAvailable:  availableStock >= item.Quantity,
			StockStatus:  model.GetStockStatus(availableStock, item.Quantity),
		}

		checkoutItems = append(checkoutItems, checkoutItem)
		subtotal += item.Subtotal
	}

	// Calculate shipping fee (simplified - in production, calculate based on items, location, etc.)
	shippingFee := 0.0
	if subtotal < 500000 {
		shippingFee = 30000 // Free shipping for orders over 500,000
	}

	return &model.CartCheckoutSummary{
		Items:       checkoutItems,
		Subtotal:    subtotal,
		ShippingFee: shippingFee,
		Discount:    cart.Discount,
		Total:       subtotal + shippingFee - cart.Discount,
		Currency:    cart.Currency,
	}, nil
}

func (s *cartServiceEnhanced) ReserveCartStock(userID uint) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return ErrCartNotFound
	}

	return s.cartRepo.ReserveCartItems(cart.ID)
}

func (s *cartServiceEnhanced) ReleaseCartStock(userID uint) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return ErrCartNotFound
	}

	return s.cartRepo.ReleaseCartItems(cart.ID)
}

// ==================== CART STATISTICS ====================

func (s *cartServiceEnhanced) GetCartStats(userID uint) (*CartStats, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartNotFound
	}

	stats := &CartStats{
		TotalItems:     cart.TotalItems,
		TotalProducts:  len(cart.Items),
		Subtotal:       cart.Subtotal,
		EstimatedTotal: cart.Total,
	}

	// Check for out of stock items and price changes
	for _, item := range cart.Items {
		availableStock := item.Product.Stock
		if item.VariantID != nil && item.Variant != nil {
			availableStock = item.Variant.Stock
		}

		if availableStock < item.Quantity {
			stats.HasOutOfStock = true
		}

		// Check if price has changed
		currentPrice := item.Product.Price
		if item.VariantID != nil && item.Variant != nil && item.Variant.Price > 0 {
			currentPrice = item.Variant.Price
		}

		if currentPrice != item.Price {
			stats.HasPriceChanged = true
		}
	}

	return stats, nil
}

// ==================== HELPER METHODS ====================

// GetCartWithItems gets cart with all items loaded
func (s *cartServiceEnhanced) GetCartWithItems(userID uint) (*model.Cart, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Ensure items are loaded
	if len(cart.Items) == 0 {
		cart.Items, _ = s.cartRepo.GetItemsByCartID(cart.ID)
	}

	return cart, nil
}

// SyncCartWithInventory syncs cart prices with current product prices
func (s *cartServiceEnhanced) SyncCartWithInventory(userID uint) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	items, err := s.cartRepo.GetItemsByCartID(cart.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			continue
		}

		// Update price if changed
		newPrice := product.Price
		if item.VariantID != nil {
			variants, _ := s.productRepo.FindVariantsByProductID(item.ProductID)
			for _, v := range variants {
				if uint(v.ID) == *item.VariantID && v.Price > 0 {
					newPrice = v.Price
					break
				}
			}
		}

		if newPrice != item.Price {
			item.Price = newPrice
			s.cartRepo.UpdateItem(&item)
		}
	}

	// Recalculate totals
	return s.cartRepo.UpdateCartTotals(cart.ID)
}

// GetExpiredCarts gets carts that haven't been updated in specified days
func (s *cartServiceEnhanced) GetExpiredCarts(days int) ([]model.Cart, error) {
	// Simplified - in production, implement proper repository method
	return []model.Cart{}, nil
}

// CleanExpiredCarts removes expired cart items
func (s *cartServiceEnhanced) CleanExpiredCarts(days int) (int64, error) {
	carts, err := s.GetExpiredCarts(days)
	if err != nil {
		return 0, err
	}

	var clearedCount int64
	for _, cart := range carts {
		count, _ := s.cartRepo.GetCartItemCount(cart.ID)
		if count > 0 {
			s.cartRepo.ClearCart(cart.ID)
			clearedCount += int64(count)
		}
	}

	return clearedCount, nil
}
