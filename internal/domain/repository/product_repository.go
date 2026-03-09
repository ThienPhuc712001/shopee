package repository

import (
	"ecommerce/internal/domain/model"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Create creates a new product
	Create(product *model.Product) error

	// Update updates an existing product
	Update(product *model.Product) error

	// Delete deletes a product by ID
	Delete(id uint) error

	// FindByID finds a product by ID
	FindByID(id uint) (*model.Product, error)

	// FindByShopID finds all products by shop ID
	FindByShopID(shopID uint, limit, offset int) ([]model.Product, int64, error)

	// FindAll retrieves all products with pagination
	FindAll(limit, offset int) ([]model.Product, int64, error)

	// Search searches products by keyword
	Search(keyword string, categoryID *uint, shopID *uint, minPrice, maxPrice *float64, limit, offset int) ([]model.Product, int64, error)

	// FindByCategory finds products by category
	FindByCategory(categoryID uint, limit, offset int) ([]model.Product, int64, error)

	// UpdateStock updates product stock
	UpdateStock(id uint, quantity int) error

	// IncrementSold increments the sold count
	IncrementSold(id uint, quantity int) error

	// IncrementViewCount increments the view count
	IncrementViewCount(id uint) error

	// FindBestSellers finds best selling products
	FindBestSellers(limit int) ([]model.Product, error)

	// FindByIDs finds products by IDs
	FindByIDs(ids []uint) ([]model.Product, error)
}
