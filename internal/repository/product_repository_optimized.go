package repository

import (
	"context"
	"ecommerce/internal/domain/model"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// OptimizedProductRepository provides optimized product queries
type OptimizedProductRepository interface {
	// GetProductsOptimized returns paginated products with optimized query
	GetProductsOptimized(page, limit int) ([]model.Product, int64, error)

	// SearchProductsOptimized performs optimized search with filters
	SearchProductsOptimized(keyword string, filters ProductFilters, page, limit int) ([]model.Product, int64, error)

	// GetProductDetailOptimized returns product detail with optimized preloading
	GetProductDetailOptimized(id uint) (*model.Product, error)

	// GetProductsByCategoryOptimized returns products by category with optimization
	GetProductsByCategoryOptimized(categoryID uint, page, limit int) ([]model.Product, int64, error)

	// GetBestSellersOptimized returns best selling products with optimization
	GetBestSellersOptimized(limit int) ([]model.Product, error)

	// GetFeaturedProductsOptimized returns featured products with optimization
	GetFeaturedProductsOptimized(limit int) ([]model.Product, error)
}

type optimizedProductRepository struct {
	db *gorm.DB
}

// NewOptimizedProductRepository creates a new optimized product repository
func NewOptimizedProductRepository(db *gorm.DB) OptimizedProductRepository {
	return &optimizedProductRepository{db: db}
}

// ==================== OPTIMIZED PRODUCT QUERIES ====================

// GetProductsOptimized returns paginated products with optimized query
// Uses:
// - Specific column selection (covering index)
// - Eager loading for related data
// - Proper pagination
// - Indexed sorting
func (r *optimizedProductRepository) GetProductsOptimized(page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Use context with timeout for long-running queries
	ctx, cancel := createContext()
	defer cancel()

	// Build base query - use filtered index on status
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("status = ?", model.ProductStatusActive)

	// Count total (separate query for accuracy)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Select specific columns (uses covering index)
	// Only select columns needed for product listing
	err := query.
		Select([]string{
			"products.id",
			"products.shop_id",
			"products.category_id",
			"products.name",
			"products.slug",
			"products.short_description",
			"products.price",
			"products.original_price",
			"products.discount_percent",
			"products.stock",
			"products.sold_count",
			"products.rating_avg",
			"products.rating_count",
			"products.status",
			"products.is_featured",
			"products.created_at",
		}).
		Preload("Category", "is_active = ?", true).
		Preload("Shop", "is_active = ?", true).
		Preload("Images", "is_primary = ?", true).
		Offset((page - 1) * limit).
		Limit(limit).
		Order("products.created_at DESC").
		Find(&products).Error

	return products, total, err
}

// SearchProductsOptimized performs optimized search with filters
// Optimization techniques:
// - Apply most selective filters first
// - Use indexed columns for filtering
// - Keyword search applied last (least selective)
// - Validate sort fields to prevent SQL injection
func (r *optimizedProductRepository) SearchProductsOptimized(keyword string, filters ProductFilters, page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	ctx, cancel := createContext()
	defer cancel()

	// Build base query
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("status = ?", model.ProductStatusActive)

	// Apply most selective filters first (uses indexes)
	// Category filter (highly selective - uses IX_products_category_id_status)
	if filters.CategoryID != nil && *filters.CategoryID > 0 {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	// Shop filter (uses IX_products_shop_id)
	if filters.ShopID != nil && *filters.ShopID > 0 {
		query = query.Where("shop_id = ?", *filters.ShopID)
	}

	// Price range filter (uses IX_products_price or IX_products_category_price)
	if filters.MinPrice != nil && *filters.MinPrice > 0 {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil && *filters.MaxPrice > 0 {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	// Rating filter (uses IX_products_rating_avg)
	if filters.MinRating != nil && *filters.MinRating > 0 {
		query = query.Where("rating_avg >= ?", *filters.MinRating)
	}

	// Brand filter
	if filters.Brands != nil && len(filters.Brands) > 0 {
		query = query.Where("brand IN ?", filters.Brands)
	}

	// Status filter
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	// Featured filter (uses IX_products_is_featured)
	if filters.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filters.IsFeatured)
	}

	// Flash sale filter (uses IX_products_flash_sale)
	if filters.IsFlashSale != nil {
		query = query.Where("is_flash_sale = ?", *filters.IsFlashSale)
	}

	// Keyword search last (least selective - uses IX_products_name)
	if keyword != "" {
		searchPattern := "%" + strings.TrimSpace(keyword) + "%"
		query = query.Where(
			"name LIKE ? OR short_description LIKE ? OR brand LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Validate and apply sorting
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(filters.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Whitelist allowed sort fields to prevent SQL injection
	allowedSortFields := map[string]bool{
		"price": true, "rating_avg": true, "sold_count": true,
		"created_at": true, "view_count": true, "name": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	// Select specific columns
	err := query.
		Select([]string{
			"products.id",
			"products.shop_id",
			"products.category_id",
			"products.name",
			"products.slug",
			"products.price",
			"products.original_price",
			"products.discount_percent",
			"products.stock",
			"products.sold_count",
			"products.rating_avg",
			"products.rating_count",
			"products.status",
			"products.is_featured",
			"products.created_at",
		}).
		Preload("Category", "is_active = ?", true).
		Preload("Shop", "is_active = ?", true).
		Preload("Images", "is_primary = ?", true).
		Offset((page - 1) * limit).
		Limit(limit).
		Order(fmt.Sprintf("products.%s %s", sortBy, sortOrder)).
		Find(&products).Error

	return products, total, err
}

// GetProductDetailOptimized returns product detail with optimized preloading
// Uses eager loading to avoid N+1 queries
func (r *optimizedProductRepository) GetProductDetailOptimized(id uint) (*model.Product, error) {
	var product model.Product

	ctx, cancel := createContext()
	defer cancel()

	// Single query with all necessary preloads (avoids N+1)
	err := r.db.WithContext(ctx).
		Preload("Shop").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Variants", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("sort_order ASC")
		}).
		Preload("Attributes", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_visible = ?", true).Order("sort_order ASC")
		}).
		Preload("Attributes.AttributeValues").
		First(&product, id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProductsByCategoryOptimized returns products by category with optimization
func (r *optimizedProductRepository) GetProductsByCategoryOptimized(categoryID uint, page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	ctx, cancel := createContext()
	defer cancel()

	// Use composite index IX_products_category_id_status
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("category_id = ? AND status = ?", categoryID, model.ProductStatusActive)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Select specific columns
	err := query.
		Select([]string{
			"products.id",
			"products.shop_id",
			"products.category_id",
			"products.name",
			"products.slug",
			"products.price",
			"products.original_price",
			"products.discount_percent",
			"products.stock",
			"products.sold_count",
			"products.rating_avg",
			"products.rating_count",
			"products.status",
			"products.created_at",
		}).
		Preload("Shop", "is_active = ?", true).
		Preload("Images", "is_primary = ?", true).
		Offset((page - 1) * limit).
		Limit(limit).
		Order("products.created_at DESC").
		Find(&products).Error

	return products, total, err
}

// GetBestSellersOptimized returns best selling products with optimization
// Uses IX_products_sold_count index for sorting
func (r *optimizedProductRepository) GetBestSellersOptimized(limit int) ([]model.Product, error) {
	var products []model.Product

	ctx, cancel := createContext()
	defer cancel()

	// Use IX_products_sold_count for efficient sorting
	err := r.db.WithContext(ctx).
		Where("status = ?", model.ProductStatusActive).
		Select([]string{
			"products.id",
			"products.shop_id",
			"products.category_id",
			"products.name",
			"products.slug",
			"products.price",
			"products.original_price",
			"products.discount_percent",
			"products.sold_count",
			"products.rating_avg",
			"products.rating_count",
		}).
		Preload("Shop", "is_active = ?", true).
		Preload("Images", "is_primary = ?", true).
		Limit(limit).
		Order("products.sold_count DESC").
		Find(&products).Error

	return products, err
}

// GetFeaturedProductsOptimized returns featured products with optimization
// Uses IX_products_is_featured filtered index
func (r *optimizedProductRepository) GetFeaturedProductsOptimized(limit int) ([]model.Product, error) {
	var products []model.Product

	ctx, cancel := createContext()
	defer cancel()

	// Use filtered index IX_products_is_featured
	err := r.db.WithContext(ctx).
		Where("status = ? AND is_featured = ?", model.ProductStatusActive, true).
		Select([]string{
			"products.id",
			"products.shop_id",
			"products.category_id",
			"products.name",
			"products.slug",
			"products.price",
			"products.original_price",
			"products.discount_percent",
			"products.sold_count",
			"products.rating_avg",
		}).
		Preload("Shop", "is_active = ?", true).
		Preload("Images", "is_primary = ?", true).
		Limit(limit).
		Order("products.created_at DESC").
		Find(&products).Error

	return products, err
}

// ==================== HELPER FUNCTIONS ====================

// createContext creates a context with timeout for database queries
func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
