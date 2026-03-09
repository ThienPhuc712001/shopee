package repository

import (
	"ecommerce/internal/domain/model"
)

// ShopRepository defines the interface for shop data operations
type ShopRepository interface {
	// Create creates a new shop
	Create(shop *model.Shop) error

	// Update updates an existing shop
	Update(shop *model.Shop) error

	// Delete deletes a shop by ID
	Delete(id uint) error

	// FindByID finds a shop by ID
	FindByID(id uint) (*model.Shop, error)

	// FindByUserID finds a shop by user ID
	FindByUserID(userID uint) (*model.Shop, error)

	// FindBySlug finds a shop by slug
	FindBySlug(slug string) (*model.Shop, error)

	// FindAll retrieves all shops with pagination
	FindAll(limit, offset int) ([]model.Shop, int64, error)

	// Search searches shops by keyword
	Search(keyword string, limit, offset int) ([]model.Shop, int64, error)

	// UpdateRating updates shop rating
	UpdateRating(id uint, rating float64, totalReviews int) error

	// IncrementTotalProducts increments total products count
	IncrementTotalProducts(id uint) error

	// IncrementTotalSales increments total sales count
	IncrementTotalSales(id uint, quantity int) error

	// FindTopShops finds top shops by rating or sales
	FindTopShops(limit int) ([]model.Shop, error)
}
