package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"ecommerce/pkg/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadService handles image upload business logic
type UploadService struct {
	repo    *repository.ImageRepository
	baseDir string
}

// UploadResult contains the result of an upload operation
type UploadResult struct {
	URL        string `json:"url"`
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	ImageType  string `json:"image_type"`
}

// NewUploadService creates a new upload service
func NewUploadService(repo *repository.ImageRepository, baseDir string) *UploadService {
	return &UploadService{
		repo:    repo,
		baseDir: baseDir,
	}
}

// ========== Product Image Upload ==========

// UploadProductImage uploads a product image
func (s *UploadService) UploadProductImage(
	file *multipart.FileHeader,
	productID int64,
	userID int64,
	ipAddress string,
	isPrimary bool,
) (*UploadResult, error) {
	// Get upload configuration
	config := model.GetUploadConfig(model.ImageTypeProduct)

	// Validate file
	if err := s.validateFile(file, config); err != nil {
		return nil, err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Validate MIME type
	if err := utils.ValidateImageMIMEType(src); err != nil {
		return nil, err
	}

	// Generate unique filename
	uniqueFilename := utils.GenerateUniqueFilename(file.Filename)

	// Create full file path
	uploadDir := filepath.Join(s.baseDir, config.UploadPath)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save file to storage
	written, err := utils.SaveFile(src, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("/%s/%s", config.UploadPath, uniqueFilename)

	// Create product image record
	productImage := &model.ProductImage{
		ProductID: uint(productID),
		URL:       url,
		AltText:   file.Filename,
		IsPrimary: isPrimary,
		SortOrder: 0,
	}

	if err := s.repo.SaveProductImage(productImage); err != nil {
		// Rollback: delete saved file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	// Save upload log
	s.saveUploadLog(userID, model.ImageTypeProduct, file.Filename, uniqueFilename, written, ipAddress)

	return &UploadResult{
		URL:       url,
		Filename:  uniqueFilename,
		Size:      written,
		ImageType: string(model.ImageTypeProduct),
	}, nil
}

// UploadMultipleProductImages uploads multiple product images
func (s *UploadService) UploadMultipleProductImages(
	files []*multipart.FileHeader,
	productID int64,
	userID int64,
	ipAddress string,
) ([]*UploadResult, error) {
	results := make([]*UploadResult, 0, len(files))

	for i, file := range files {
		isPrimary := i == 0 // First image is primary
		result, err := s.UploadProductImage(file, productID, userID, ipAddress, isPrimary)
		if err != nil {
			return results, fmt.Errorf("failed to upload image %d: %w", i+1, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// ========== Review Image Upload ==========

// UploadReviewImage uploads a review image
func (s *UploadService) UploadReviewImage(
	file *multipart.FileHeader,
	reviewID int64,
	userID int64,
	ipAddress string,
) (*UploadResult, error) {
	// Get upload configuration
	config := model.GetUploadConfig(model.ImageTypeReview)

	// Validate file
	if err := s.validateFile(file, config); err != nil {
		return nil, err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Validate MIME type
	if err := utils.ValidateImageMIMEType(src); err != nil {
		return nil, err
	}

	// Generate unique filename
	uniqueFilename := utils.GenerateUniqueFilename(file.Filename)

	// Create full file path
	uploadDir := filepath.Join(s.baseDir, config.UploadPath)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save file to storage
	written, err := utils.SaveFile(src, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("/%s/%s", config.UploadPath, uniqueFilename)

	// Create review image record
	reviewImage := &model.ReviewImage{
		ReviewID: reviewID,
		URL:      url,
	}

	if err := s.repo.SaveReviewImage(reviewImage); err != nil {
		// Rollback: delete saved file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	// Save upload log
	s.saveUploadLog(userID, model.ImageTypeReview, file.Filename, uniqueFilename, written, ipAddress)

	return &UploadResult{
		URL:       url,
		Filename:  uniqueFilename,
		Size:      written,
		ImageType: string(model.ImageTypeReview),
	}, nil
}

// ========== User Avatar Upload ==========

// UploadUserAvatar uploads a user avatar
func (s *UploadService) UploadUserAvatar(
	file *multipart.FileHeader,
	userID int64,
	ipAddress string,
) (*UploadResult, error) {
	// Get upload configuration
	config := model.GetUploadConfig(model.ImageTypeAvatar)

	// Validate file
	if err := s.validateFile(file, config); err != nil {
		return nil, err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Validate MIME type
	if err := utils.ValidateImageMIMEType(src); err != nil {
		return nil, err
	}

	// Generate unique filename
	uniqueFilename := utils.GenerateUniqueFilename(file.Filename)

	// Create full file path
	uploadDir := filepath.Join(s.baseDir, config.UploadPath)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save file to storage
	written, err := utils.SaveFile(src, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("/%s/%s", config.UploadPath, uniqueFilename)

	// Check if user has existing avatar and delete old file
	existingAvatar, err := s.repo.GetUserAvatar(userID)
	if err == nil && existingAvatar != nil {
		oldFilePath := filepath.Join(s.baseDir, existingAvatar.URL)
		os.Remove(oldFilePath)
	}

	// Create or update user avatar record
	userAvatar := &model.UserAvatar{
		UserID: userID,
		URL:    url,
	}

	if err := s.repo.SaveUserAvatar(userAvatar); err != nil {
		// Rollback: delete saved file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save avatar record: %w", err)
	}

	// Save upload log
	s.saveUploadLog(userID, model.ImageTypeAvatar, file.Filename, uniqueFilename, written, ipAddress)

	return &UploadResult{
		URL:       url,
		Filename:  uniqueFilename,
		Size:      written,
		ImageType: string(model.ImageTypeAvatar),
	}, nil
}

// ========== Helper Methods ==========

// validateFile validates the uploaded file
func (s *UploadService) validateFile(file *multipart.FileHeader, config model.ImageUpload) error {
	return utils.ValidateFile(file, config.MaxSize, config.AllowedExts)
}

// saveUploadLog saves an upload log entry
func (s *UploadService) saveUploadLog(
	userID int64,
	imageType model.ImageType,
	originalName string,
	storedName string,
	fileSize int64,
	ipAddress string,
) {
	log := &model.ImageUploadLog{
		UserID:       userID,
		ImageType:    string(imageType),
		OriginalName: originalName,
		StoredName:   storedName,
		FileSize:     fileSize,
		IPAddress:    ipAddress,
	}
	_ = s.repo.SaveUploadLog(log)
}

// DeleteProductImage deletes a product image
func (s *UploadService) DeleteProductImage(imageID int64) error {
	// Get image record
	img, err := s.repo.GetProductImage(imageID)
	if err != nil {
		return err
	}

	// Delete file from storage
	filePath := filepath.Join(s.baseDir, img.URL)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete record from database
	return s.repo.DeleteProductImage(imageID)
}

// DeleteReviewImage deletes a review image
func (s *UploadService) DeleteReviewImage(imageID int64) error {
	// Get image record
	img, err := s.repo.GetReviewImage(imageID)
	if err != nil {
		return err
	}

	// Delete file from storage
	filePath := filepath.Join(s.baseDir, img.URL)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete record from database
	return s.repo.DeleteReviewImage(imageID)
}

// GetProductImages retrieves all images for a product
func (s *UploadService) GetProductImages(productID int64) ([]model.ProductImage, error) {
	return s.repo.GetProductImages(productID)
}

// GetReviewImages retrieves all images for a review
func (s *UploadService) GetReviewImages(reviewID int64) ([]model.ReviewImage, error) {
	return s.repo.GetReviewImages(reviewID)
}

// GetUserAvatar retrieves a user's avatar
func (s *UploadService) GetUserAvatar(userID int64) (*model.UserAvatar, error) {
	return s.repo.GetUserAvatar(userID)
}

// ========== Error Definitions ==========

var (
	ErrInvalidProductID = errors.New("invalid product ID")
	ErrInvalidReviewID  = errors.New("invalid review ID")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrUploadFailed     = errors.New("upload failed")
)
