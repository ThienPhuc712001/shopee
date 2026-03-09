package repository

import (
	"ecommerce/internal/domain/model"
)

// CartRepository defines the interface for cart data operations
type CartRepository interface {
	// Create creates a new cart
	Create(cart *model.Cart) error

	// FindByUserID finds a cart by user ID
	FindByUserID(userID uint) (*model.Cart, error)

	// FindOrCreate finds or creates a cart for user
	FindOrCreate(userID uint) (*model.Cart, error)

	// Update updates an existing cart
	Update(cart *model.Cart) error

	// Delete deletes a cart by ID
	Delete(id uint) error

	// AddItem adds an item to cart
	AddItem(item *model.CartItem) error

	// UpdateItem updates a cart item
	UpdateItem(item *model.CartItem) error

	// DeleteItem deletes a cart item by ID
	DeleteItem(id uint) error

	// FindItemByCartAndProduct finds a cart item by cart ID and product ID
	FindItemByCartAndProduct(cartID, productID uint) (*model.CartItem, error)

	// GetItemsByCartID gets all items in a cart
	GetItemsByCartID(cartID uint) ([]model.CartItem, error)

	// ClearItems removes all items from a cart
	ClearItems(cartID uint) error

	// UpdateTotals updates cart totals
	UpdateTotals(cartID uint) error
}
