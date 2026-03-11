package repository

import (
	"ecommerce/internal/domain/model"
	"errors"

	"gorm.io/gorm"
)

// ImageRepository handles database operations for images
type ImageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new image repository
func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// ========== Product Image Methods ==========

// SaveProductImage saves a product image to the database
func (r *ImageRepository) SaveProductImage(img *model.ProductImage) error {
	return r.db.Create(img).Error
}

// GetProductImages retrieves all images for a product
func (r *ImageRepository) GetProductImages(productID int64) ([]model.ProductImage, error) {
	var images []model.ProductImage
	err := r.db.Where("product_id = ?", productID).Order("sort_order ASC, created_at ASC").Find(&images).Error
	return images, err
}

// GetProductImage retrieves a specific product image by ID
func (r *ImageRepository) GetProductImage(id int64) (*model.ProductImage, error) {
	var img model.ProductImage
	err := r.db.First(&img, id).Error
	if err != nil {
		return nil, err
	}
	return &img, nil
}

// SetPrimaryImage sets one image as primary and others as non-primary
func (r *ImageRepository) SetPrimaryImage(productID, imageID int64) error {
	tx := r.db.Begin()

	// Set all images to non-primary
	if err := tx.Model(&model.ProductImage{}).Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set selected image as primary
	if err := tx.Model(&model.ProductImage{}).Where("id = ?", imageID).Update("is_primary", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DeleteProductImage deletes a product image from the database
func (r *ImageRepository) DeleteProductImage(id int64) error {
	return r.db.Delete(&model.ProductImage{}, id).Error
}

// DeleteProductImages deletes all images for a product
func (r *ImageRepository) DeleteProductImages(productID int64) error {
	return r.db.Where("product_id = ?", productID).Delete(&model.ProductImage{}).Error
}

// UpdateProductImageSortOrder updates the sort order of product images
func (r *ImageRepository) UpdateProductImageSortOrder(imageID int64, sortOrder int) error {
	return r.db.Model(&model.ProductImage{}).Where("id = ?", imageID).Update("sort_order", sortOrder).Error
}

// ========== Review Image Methods ==========

// SaveReviewImage saves a review image to the database
func (r *ImageRepository) SaveReviewImage(img *model.ReviewImage) error {
	return r.db.Create(img).Error
}

// GetReviewImages retrieves all images for a review
func (r *ImageRepository) GetReviewImages(reviewID int64) ([]model.ReviewImage, error) {
	var images []model.ReviewImage
	err := r.db.Where("review_id = ?", reviewID).Order("created_at ASC").Find(&images).Error
	return images, err
}

// GetReviewImage retrieves a specific review image by ID
func (r *ImageRepository) GetReviewImage(id int64) (*model.ReviewImage, error) {
	var img model.ReviewImage
	err := r.db.First(&img, id).Error
	if err != nil {
		return nil, err
	}
	return &img, nil
}

// DeleteReviewImage deletes a review image from the database
func (r *ImageRepository) DeleteReviewImage(id int64) error {
	return r.db.Delete(&model.ReviewImage{}, id).Error
}

// DeleteReviewImages deletes all images for a review
func (r *ImageRepository) DeleteReviewImages(reviewID int64) error {
	return r.db.Where("review_id = ?", reviewID).Delete(&model.ReviewImage{}).Error
}

// ========== User Avatar Methods ==========

// SaveUserAvatar saves or updates a user avatar
func (r *ImageRepository) SaveUserAvatar(avatar *model.UserAvatar) error {
	// Check if user already has an avatar
	existing := &model.UserAvatar{}
	err := r.db.Where("user_id = ?", avatar.UserID).First(existing).Error

	if err == nil {
		// Update existing avatar
		existing.URL = avatar.URL
		return r.db.Save(existing).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new avatar
		return r.db.Create(avatar).Error
	}

	return err
}

// GetUserAvatar retrieves a user's avatar
func (r *ImageRepository) GetUserAvatar(userID int64) (*model.UserAvatar, error) {
	var avatar model.UserAvatar
	err := r.db.Where("user_id = ?", userID).First(&avatar).Error
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

// DeleteUserAvatar deletes a user's avatar
func (r *ImageRepository) DeleteUserAvatar(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserAvatar{}).Error
}

// ========== Upload Log Methods ==========

// SaveUploadLog saves an upload log entry
func (r *ImageRepository) SaveUploadLog(log *model.ImageUploadLog) error {
	return r.db.Create(log).Error
}

// GetUploadLogs retrieves upload logs for a user
func (r *ImageRepository) GetUploadLogs(userID int64, limit int) ([]model.ImageUploadLog, error) {
	var logs []model.ImageUploadLog
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
