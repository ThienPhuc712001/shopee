package repository

import (
	"ecommerce/internal/domain/model"
	"context"

	"gorm.io/gorm"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// UpdateCategory updates an existing category
func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory deletes a category (soft delete)
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// GetCategoryByID retrieves a category by ID
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (r *CategoryRepository) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories retrieves all categories (with optional filters)
func (r *CategoryRepository) GetAllCategories(ctx context.Context, isActive *bool) ([]model.Category, error) {
	var categories []model.Category
	query := r.db.WithContext(ctx)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("sort_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

// GetRootCategories retrieves all root categories (no parent)
func (r *CategoryRepository) GetRootCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetSubCategories retrieves all subcategories of a parent category
func (r *CategoryRepository) GetSubCategories(ctx context.Context, parentID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryTree retrieves categories in tree structure
func (r *CategoryRepository) GetCategoryTree(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL").
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoriesByLevel retrieves categories at a specific level
func (r *CategoryRepository) GetCategoriesByLevel(ctx context.Context, level int) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).
		Where("level = ?", level).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryWithProductCount retrieves category with product count
func (r *CategoryRepository) GetCategoryWithProductCount(ctx context.Context, id uint) (*model.CategoryWithCount, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}

	var count int64
	err = r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("category_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return nil, err
	}

	return &model.CategoryWithCount{
		Category:     category,
		ProductCount: count,
	}, nil
}

// GetAllCategoriesWithProductCount retrieves all categories with product counts
func (r *CategoryRepository) GetAllCategoriesWithProductCount(ctx context.Context) ([]model.CategoryWithCount, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	result := make([]model.CategoryWithCount, 0, len(categories))
	for _, cat := range categories {
		var count int64
		err = r.db.WithContext(ctx).
			Model(&model.Product{}).
			Where("category_id = ? AND deleted_at IS NULL", cat.ID).
			Count(&count).Error
		if err != nil {
			return nil, err
		}

		result = append(result, model.CategoryWithCount{
			Category:     cat,
			ProductCount: count,
		})
	}

	return result, nil
}

// GetProductsByCategory retrieves all products in a category
func (r *CategoryRepository) GetProductsByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Get total count
	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get products with pagination
	err = r.db.WithContext(ctx).
		Preload("Images").
		Preload("Category").
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

// GetCategoryPath retrieves the full path from root to category (for breadcrumbs)
func (r *CategoryRepository) GetCategoryPath(ctx context.Context, id uint) ([]model.Category, error) {
	var path []model.Category
	
	current, err := r.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for current != nil {
		path = append([]model.Category{*current}, path...)
		if current.ParentID == nil {
			break
		}
		current, err = r.GetCategoryByID(ctx, *current.ParentID)
		if err != nil {
			break
		}
	}

	return path, nil
}

// UpdateCategoryLevel updates the level of a category and its children
func (r *CategoryRepository) UpdateCategoryLevel(ctx context.Context, id uint, level int) error {
	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", id).
		Update("level", level).Error
}

// HasChildren checks if a category has subcategories
func (r *CategoryRepository) HasChildren(ctx context.Context, parentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Count(&count).Error
	return count > 0, err
}

// HasProducts checks if a category has products
func (r *CategoryRepository) HasProducts(ctx context.Context, categoryID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Count(&count).Error
	return count > 0, err
}

// SearchCategories searches categories by name or description
func (r *CategoryRepository) SearchCategories(ctx context.Context, query string) ([]model.Category, error) {
	var categories []model.Category
	searchTerm := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Where("name LIKE ? OR description LIKE ?", searchTerm, searchTerm).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetFeaturedCategories retrieves featured/active categories for homepage
func (r *CategoryRepository) GetFeaturedCategories(ctx context.Context, limit int) ([]model.CategoryWithCount, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND level = ?", true, 0).
		Order("sort_order ASC").
		Limit(limit).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	result := make([]model.CategoryWithCount, 0, len(categories))
	for _, cat := range categories {
		var count int64
		err = r.db.WithContext(ctx).
			Model(&model.Product{}).
			Where("category_id = ? AND deleted_at IS NULL", cat.ID).
			Count(&count).Error
		if err != nil {
			return nil, err
		}

		result = append(result, model.CategoryWithCount{
			Category:     cat,
			ProductCount: count,
		})
	}

	return result, nil
}
