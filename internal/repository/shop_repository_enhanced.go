package repository

import (
	"ecommerce/internal/domain/model"
	"gorm.io/gorm"
)

// ShopRepositoryEnhanced defines the enhanced interface for shop data operations
type ShopRepositoryEnhanced interface {
	// Basic CRUD
	Create(shop *model.Shop) error
	Update(shop *model.Shop) error
	Delete(id uint) error
	FindByID(id uint) (*model.Shop, error)
	FindByUserID(userID uint) (*model.Shop, error)
	FindBySlug(slug string) (*model.Shop, error)

	// Queries
	FindAll(limit, offset int) ([]model.Shop, int64, error)
	Search(keyword string, limit, offset int) ([]model.Shop, int64, error)

	// Statistics
	UpdateRating(id uint, rating float64, totalReviews int) error
	IncrementTotalProducts(id uint) error
	IncrementTotalSales(id uint, quantity int) error
	FindTopShops(limit int) ([]model.Shop, error)
}

type shopRepositoryEnhanced struct {
	db *gorm.DB
}

// NewShopRepositoryEnhanced creates a new enhanced shop repository
func NewShopRepositoryEnhanced(db *gorm.DB) ShopRepositoryEnhanced {
	return &shopRepositoryEnhanced{db: db}
}

func (r *shopRepositoryEnhanced) Create(shop *model.Shop) error {
	return r.db.Create(shop).Error
}

func (r *shopRepositoryEnhanced) Update(shop *model.Shop) error {
	return r.db.Save(shop).Error
}

func (r *shopRepositoryEnhanced) Delete(id uint) error {
	return r.db.Delete(&model.Shop{}, id).Error
}

func (r *shopRepositoryEnhanced) FindByID(id uint) (*model.Shop, error) {
	var shop model.Shop
	err := r.db.Preload("User").First(&shop, id).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepositoryEnhanced) FindByUserID(userID uint) (*model.Shop, error) {
	var shop model.Shop
	err := r.db.Where("user_id = ?", userID).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepositoryEnhanced) FindBySlug(slug string) (*model.Shop, error) {
	var shop model.Shop
	err := r.db.Where("slug = ?", slug).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepositoryEnhanced) FindAll(limit, offset int) ([]model.Shop, int64, error) {
	var shops []model.Shop
	var total int64

	if err := r.db.Model(&model.Shop{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("User").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&shops).Error

	return shops, total, err
}

func (r *shopRepositoryEnhanced) Search(keyword string, limit, offset int) ([]model.Shop, int64, error) {
	var shops []model.Shop
	var total int64

	searchPattern := "%" + keyword + "%"

	if err := r.db.Model(&model.Shop{}).
		Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern).
		Preload("User").
		Limit(limit).Offset(offset).
		Find(&shops).Error

	return shops, total, err
}

func (r *shopRepositoryEnhanced) UpdateRating(id uint, rating float64, totalReviews int) error {
	return r.db.Model(&model.Shop{}).Where("id = ?", id).Updates(map[string]interface{}{
		"rating":        rating,
		"total_reviews": totalReviews,
	}).Error
}

func (r *shopRepositoryEnhanced) IncrementTotalProducts(id uint) error {
	return r.db.Model(&model.Shop{}).Where("id = ?", id).UpdateColumn("total_products", gorm.Expr("total_products + 1")).Error
}

func (r *shopRepositoryEnhanced) IncrementTotalSales(id uint, quantity int) error {
	return r.db.Model(&model.Shop{}).Where("id = ?", id).UpdateColumn("total_sales", gorm.Expr("total_sales + ?", quantity)).Error
}

func (r *shopRepositoryEnhanced) FindTopShops(limit int) ([]model.Shop, error) {
	var shops []model.Shop
	err := r.db.Where("status = ?", model.ShopStatusActive).
		Order("rating DESC, total_sales DESC").
		Limit(limit).
		Find(&shops).Error
	return shops, err
}
