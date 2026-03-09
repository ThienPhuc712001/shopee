package repository

import (
	"ecommerce/internal/domain/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ProductRepositoryEnhanced defines the enhanced interface for product data operations
type ProductRepositoryEnhanced interface {
	// Basic CRUD
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id uint) error
	FindByID(id uint) (*model.Product, error)

	// Product Queries
	FindAll(limit, offset int) ([]model.Product, int64, error)
	FindByShopID(shopID uint, limit, offset int) ([]model.Product, int64, error)
	FindByCategoryID(categoryID uint, limit, offset int) ([]model.Product, int64, error)
	FindFeatured(limit int) ([]model.Product, error)
	FindBestSellers(limit int) ([]model.Product, error)
	FindLatest(limit int) ([]model.Product, error)

	// Search & Filter
	Search(keyword string, filters ProductFilters, limit, offset int) ([]model.Product, int64, error)
	FilterByPrice(minPrice, maxPrice float64, limit, offset int) ([]model.Product, int64, error)
	FilterByRating(minRating float64, limit, offset int) ([]model.Product, int64, error)

	// Product Images
	CreateImage(image *model.ProductImage) error
	UpdateImage(image *model.ProductImage) error
	DeleteImage(id uint) error
	FindImagesByProductID(productID uint) ([]model.ProductImage, error)

	// Product Variants
	CreateVariant(variant *model.ProductVariant) error
	UpdateVariant(variant *model.ProductVariant) error
	DeleteVariant(id uint) error
	FindVariantsByProductID(productID uint) ([]model.ProductVariant, error)
	FindVariantBySKU(sku string) (*model.ProductVariant, error)

	// Product Attributes
	CreateAttribute(attr *model.ProductAttribute) error
	UpdateAttribute(attr *model.ProductAttribute) error
	DeleteAttribute(id uint) error
	FindAttributesByProductID(productID uint) ([]model.ProductAttribute, error)

	// Inventory Management
	UpdateStock(productID uint, variantID *uint, quantity int) error
	ReserveStock(productID uint, variantID *uint, quantity int) error
	ReleaseStock(productID uint, variantID *uint, quantity int) error
	DecreaseStock(productID uint, variantID *uint, quantity int) error
	GetInventory(productID uint, variantID *uint) (*model.ProductInventory, error)

	// Bulk Operations
	BulkCreate(products []*model.Product) error
	BulkUpdateStatus(ids []uint, status model.ProductStatus) error

	// Analytics
	IncrementViewCount(id uint) error
	IncrementSoldCount(id uint, quantity int) error

	// Categories
	CreateCategory(category *model.Category) error
	UpdateCategory(category *model.Category) error
	DeleteCategory(id uint) error
	FindCategoryByID(id uint) (*model.Category, error)
	FindAllCategories() ([]model.Category, error)
	FindCategoriesByParentID(parentID *uint) ([]model.Category, error)
}

// ProductFilters represents search and filter parameters
type ProductFilters struct {
	CategoryID   *uint
	ShopID       *uint
	MinPrice     *float64
	MaxPrice     *float64
	MinRating    *float64
	Brands       []string
	Status       *model.ProductStatus
	IsFeatured   *bool
	IsFlashSale  *bool
	SortBy       string // price, rating, sold_count, created_at
	SortOrder    string // asc, desc
}

type productRepositoryEnhanced struct {
	db *gorm.DB
}

// NewProductRepositoryEnhanced creates a new enhanced product repository
func NewProductRepositoryEnhanced(db *gorm.DB) ProductRepositoryEnhanced {
	return &productRepositoryEnhanced{db: db}
}

// ==================== BASIC CRUD ====================

func (r *productRepositoryEnhanced) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepositoryEnhanced) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepositoryEnhanced) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepositoryEnhanced) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Shop").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Variants").
		Preload("Attributes").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// ==================== PRODUCT QUERIES ====================

func (r *productRepositoryEnhanced) FindAll(limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if err := r.db.Model(&model.Product{}).Where("status = ?", model.ProductStatusActive).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("status = ?", model.ProductStatusActive).
		Preload("Shop").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true).Order("sort_order ASC")
		}).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepositoryEnhanced) FindByShopID(shopID uint, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if err := r.db.Model(&model.Product{}).Where("shop_id = ? AND status = ?", shopID, model.ProductStatusActive).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("shop_id = ? AND status = ?", shopID, model.ProductStatusActive).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepositoryEnhanced) FindByCategoryID(categoryID uint, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if err := r.db.Model(&model.Product{}).Where("category_id = ? AND status = ?", categoryID, model.ProductStatusActive).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("category_id = ? AND status = ?", categoryID, model.ProductStatusActive).
		Preload("Shop").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepositoryEnhanced) FindFeatured(limit int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("status = ? AND is_featured = ?", model.ProductStatusActive, true).
		Preload("Shop").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

func (r *productRepositoryEnhanced) FindBestSellers(limit int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("status = ?", model.ProductStatusActive).
		Preload("Shop").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Limit(limit).
		Order("sold_count DESC").
		Find(&products).Error
	return products, err
}

func (r *productRepositoryEnhanced) FindLatest(limit int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("status = ?", model.ProductStatusActive).
		Preload("Shop").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// ==================== SEARCH & FILTER ====================

func (r *productRepositoryEnhanced) Search(keyword string, filters ProductFilters, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("status = ?", model.ProductStatusActive)

	// Keyword search
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR short_description LIKE ? OR brand LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply filters
	if filters.CategoryID != nil && *filters.CategoryID > 0 {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if filters.ShopID != nil && *filters.ShopID > 0 {
		query = query.Where("shop_id = ?", *filters.ShopID)
	}

	if filters.MinPrice != nil && *filters.MinPrice > 0 {
		query = query.Where("price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil && *filters.MaxPrice > 0 {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	if filters.MinRating != nil && *filters.MinRating > 0 {
		query = query.Where("rating_avg >= ?", *filters.MinRating)
	}

	if filters.Brands != nil && len(filters.Brands) > 0 {
		query = query.Where("brand IN ?", filters.Brands)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filters.IsFeatured)
	}

	if filters.IsFlashSale != nil {
		query = query.Where("is_flash_sale = ?", *filters.IsFlashSale)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filters.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	query = query.Preload("Shop").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true)
		}).
		Limit(limit).Offset(offset)

	// Validate sort field to prevent SQL injection
	validSortFields := map[string]bool{
		"price": true, "rating_avg": true, "sold_count": true,
		"created_at": true, "view_count": true, "name": true,
	}
	if validSortFields[sortBy] {
		query = query.Order(fmt.Sprintf("%s %s", sortBy, strings.ToUpper(sortOrder)))
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Find(&products).Error
	return products, total, err
}

func (r *productRepositoryEnhanced) FilterByPrice(minPrice, maxPrice float64, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("status = ?", model.ProductStatusActive)

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Shop").Preload("Images").
		Limit(limit).Offset(offset).
		Order("price ASC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepositoryEnhanced) FilterByRating(minRating float64, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).
		Where("status = ? AND rating_avg >= ?", model.ProductStatusActive, minRating)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Shop").Preload("Images").
		Limit(limit).Offset(offset).
		Order("rating_avg DESC").
		Find(&products).Error

	return products, total, err
}

// ==================== PRODUCT IMAGES ====================

func (r *productRepositoryEnhanced) CreateImage(image *model.ProductImage) error {
	return r.db.Create(image).Error
}

func (r *productRepositoryEnhanced) UpdateImage(image *model.ProductImage) error {
	return r.db.Save(image).Error
}

func (r *productRepositoryEnhanced) DeleteImage(id uint) error {
	return r.db.Delete(&model.ProductImage{}, id).Error
}

func (r *productRepositoryEnhanced) FindImagesByProductID(productID uint) ([]model.ProductImage, error) {
	var images []model.ProductImage
	err := r.db.Where("product_id = ?", productID).
		Order("is_primary DESC, sort_order ASC").
		Find(&images).Error
	return images, err
}

// ==================== PRODUCT VARIANTS ====================

func (r *productRepositoryEnhanced) CreateVariant(variant *model.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *productRepositoryEnhanced) UpdateVariant(variant *model.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *productRepositoryEnhanced) DeleteVariant(id uint) error {
	return r.db.Delete(&model.ProductVariant{}, id).Error
}

func (r *productRepositoryEnhanced) FindVariantsByProductID(productID uint) ([]model.ProductVariant, error) {
	var variants []model.ProductVariant
	err := r.db.Where("product_id = ? AND deleted_at IS NULL", productID).
		Preload("Inventory").
		Order("sort_order ASC").
		Find(&variants).Error
	return variants, err
}

func (r *productRepositoryEnhanced) FindVariantBySKU(sku string) (*model.ProductVariant, error) {
	var variant model.ProductVariant
	err := r.db.Where("sku = ?", sku).
		Preload("Product").
		Preload("Inventory").
		First(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

// ==================== PRODUCT ATTRIBUTES ====================

func (r *productRepositoryEnhanced) CreateAttribute(attr *model.ProductAttribute) error {
	return r.db.Create(attr).Error
}

func (r *productRepositoryEnhanced) UpdateAttribute(attr *model.ProductAttribute) error {
	return r.db.Save(attr).Error
}

func (r *productRepositoryEnhanced) DeleteAttribute(id uint) error {
	return r.db.Delete(&model.ProductAttribute{}, id).Error
}

func (r *productRepositoryEnhanced) FindAttributesByProductID(productID uint) ([]model.ProductAttribute, error) {
	var attributes []model.ProductAttribute
	err := r.db.Where("product_id = ?", productID).
		Preload("AttributeValues").
		Order("sort_order ASC").
		Find(&attributes).Error
	return attributes, err
}

// ==================== INVENTORY MANAGEMENT ====================

func (r *productRepositoryEnhanced) UpdateStock(productID uint, variantID *uint, quantity int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if variantID != nil {
		// Update variant stock
		if err := tx.Model(&model.ProductVariant{}).Where("id = ?", *variantID).Update("stock", quantity).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Update product stock
		if err := tx.Model(&model.Product{}).Where("id = ?", productID).Update("stock", quantity).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update inventory record
	var inventory model.ProductInventory
	if variantID != nil {
		if err := tx.Where("variant_id = ?", *variantID).First(&inventory).Error; err == nil {
			tx.Model(&inventory).Updates(map[string]interface{}{
				"quantity":  quantity,
				"available": quantity - inventory.Reserved,
			})
		}
	}

	return tx.Commit().Error
}

func (r *productRepositoryEnhanced) ReserveStock(productID uint, variantID *uint, quantity int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if variantID != nil {
		// Reserve variant stock
		var variant model.ProductVariant
		if err := tx.Where("id = ? AND stock >= ?", *variantID, quantity).First(&variant).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("insufficient stock for variant")
		}

		if err := tx.Model(&model.ProductVariant{}).Where("id = ?", *variantID).
			UpdateColumn("reserved_stock", gorm.Expr("reserved_stock + ?", quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Reserve product stock
		var product model.Product
		if err := tx.Where("id = ? AND stock >= ?", productID, quantity).First(&product).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("insufficient stock")
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", productID).
			UpdateColumn("reserved_stock", gorm.Expr("reserved_stock + ?", quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *productRepositoryEnhanced) ReleaseStock(productID uint, variantID *uint, quantity int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if variantID != nil {
		if err := tx.Model(&model.ProductVariant{}).Where("id = ?", *variantID).
			UpdateColumn("reserved_stock", gorm.Expr("reserved_stock - ?", quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Model(&model.Product{}).Where("id = ?", productID).
			UpdateColumn("reserved_stock", gorm.Expr("reserved_stock - ?", quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *productRepositoryEnhanced) DecreaseStock(productID uint, variantID *uint, quantity int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if variantID != nil {
		var variant model.ProductVariant
		availableStock := variant.Stock - variant.ReservedStock
		if availableStock < quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient available stock")
		}

		if err := tx.Model(&model.ProductVariant{}).Where("id = ?", *variantID).
			Updates(map[string]interface{}{
				"stock":           gorm.Expr("stock - ?", quantity),
				"reserved_stock":  gorm.Expr("reserved_stock - ?", quantity),
			}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		var product model.Product
		availableStock := product.Stock - product.ReservedStock
		if availableStock < quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient available stock")
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", productID).
			Updates(map[string]interface{}{
				"stock":          gorm.Expr("stock - ?", quantity),
				"reserved_stock": gorm.Expr("reserved_stock - ?", quantity),
			}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *productRepositoryEnhanced) GetInventory(productID uint, variantID *uint) (*model.ProductInventory, error) {
	var inventory model.ProductInventory

	var err error
	if variantID != nil {
		err = r.db.Where("variant_id = ?", *variantID).First(&inventory).Error
	} else {
		err = r.db.Where("product_id = ? AND variant_id IS NULL", productID).First(&inventory).Error
	}

	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

// ==================== BULK OPERATIONS ====================

func (r *productRepositoryEnhanced) BulkCreate(products []*model.Product) error {
	return r.db.CreateInBatches(products, 100).Error
}

func (r *productRepositoryEnhanced) BulkUpdateStatus(ids []uint, status model.ProductStatus) error {
	return r.db.Model(&model.Product{}).Where("id IN ?", ids).Update("status", status).Error
}

// ==================== ANALYTICS ====================

func (r *productRepositoryEnhanced) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *productRepositoryEnhanced) IncrementSoldCount(id uint, quantity int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).
		UpdateColumn("sold_count", gorm.Expr("sold_count + ?", quantity)).Error
}

// ==================== CATEGORIES ====================

func (r *productRepositoryEnhanced) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *productRepositoryEnhanced) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *productRepositoryEnhanced) DeleteCategory(id uint) error {
	// Check if category has products
	var count int64
	r.db.Model(&model.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with products")
	}
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *productRepositoryEnhanced) FindCategoryByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.Preload("Parent").Preload("Children").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *productRepositoryEnhanced) FindAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("is_active = ? AND deleted_at IS NULL", true).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Order("level ASC, sort_order ASC").
		Find(&categories).Error
	return categories, err
}

func (r *productRepositoryEnhanced) FindCategoriesByParentID(parentID *uint) ([]model.Category, error) {
	var categories []model.Category
	query := r.db.Where("is_active = ? AND deleted_at IS NULL", true)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}
