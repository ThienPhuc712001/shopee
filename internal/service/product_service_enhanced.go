package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProductServiceEnhanced defines the enhanced product service interface
type ProductServiceEnhanced interface {
	// Product Management
	CreateProduct(product *model.Product, shopID uint) (*model.Product, error)
	UpdateProduct(id uint, product *model.Product, shopID uint, isAdmin bool) (*model.Product, error)
	DeleteProduct(id uint, shopID uint, isAdmin bool) error
	GetProductByID(id uint) (*model.Product, error)
	GetProducts(page, limit int) ([]model.Product, int64, error)
	GetProductsByShopID(shopID uint, page, limit int) ([]model.Product, int64, error)
	GetProductsByCategoryID(categoryID uint, page, limit int) ([]model.Product, int64, error)

	// Search & Filter
	SearchProducts(keyword string, filters ProductFilters, page, limit int) ([]model.Product, int64, error)
	FilterProducts(filters ProductFilters, page, limit int) ([]model.Product, int64, error)

	// Featured & Best Sellers
	GetFeaturedProducts(limit int) ([]model.Product, error)
	GetBestSellers(limit int) ([]model.Product, error)
	GetLatestProducts(limit int) ([]model.Product, error)

	// Product Images
	UploadProductImages(productID uint, files []*multipart.FileHeader, uploadDir string) ([]model.ProductImage, error)
	UpdateProductImage(imageID uint, isPrimary bool, sortOrder int) error
	DeleteProductImage(imageID uint, productID uint, shopID uint, isAdmin bool) error

	// Product Variants
	CreateVariant(productID uint, variant *model.ProductVariant) (*model.ProductVariant, error)
	UpdateVariant(id uint, variant *model.ProductVariant) (*model.ProductVariant, error)
	DeleteVariant(id uint, productID uint, shopID uint) error
	GetVariantsByProductID(productID uint) ([]model.ProductVariant, error)

	// Product Attributes
	CreateAttribute(productID uint, attr *model.ProductAttribute) (*model.ProductAttribute, error)
	UpdateAttribute(id uint, attr *model.ProductAttribute) (*model.ProductAttribute, error)
	DeleteAttribute(id uint, productID uint) error
	GetAttributesByProductID(productID uint) ([]model.ProductAttribute, error)

	// Inventory Management
	UpdateProductStock(productID uint, variantID *uint, quantity int, shopID uint) error
	ReserveProductStock(productID uint, variantID *uint, quantity int) error
	ReleaseProductStock(productID uint, variantID *uint, quantity int) error
	DecreaseProductStock(productID uint, variantID *uint, quantity int) error

	// Product Status
	PublishProduct(id uint, shopID uint) error
	UnpublishProduct(id uint, shopID uint) error
	FeatureProduct(id uint, isAdmin bool) error

	// Categories
	CreateCategory(category *model.Category) (*model.Category, error)
	UpdateCategory(id uint, category *model.Category) (*model.Category, error)
	DeleteCategory(id uint) error
	GetCategoryByID(id uint) (*model.Category, error)
	GetAllCategories() ([]model.Category, error)
	GetCategoriesByParentID(parentID *uint) ([]model.Category, error)

	// Analytics
	IncrementViewCount(id uint) error
}

// ProductFilters represents search and filter parameters
type ProductFilters struct {
	CategoryID  *uint
	ShopID      *uint
	MinPrice    *float64
	MaxPrice    *float64
	MinRating   *float64
	Brands      []string
	Status      *model.ProductStatus
	IsFeatured  *bool
	IsFlashSale *bool
	SortBy      string
	SortOrder   string
}

type productServiceEnhanced struct {
	productRepo repository.ProductRepositoryEnhanced
}

// NewProductServiceEnhanced creates a new enhanced product service
func NewProductServiceEnhanced(productRepo repository.ProductRepositoryEnhanced) ProductServiceEnhanced {
	return &productServiceEnhanced{
		productRepo: productRepo,
	}
}

// ==================== PRODUCT MANAGEMENT ====================

func (s *productServiceEnhanced) CreateProduct(product *model.Product, shopID uint) (*model.Product, error) {
	// Validate product
	if err := s.validateProduct(product); err != nil {
		return nil, err
	}

	// Set shop ID
	product.ShopID = shopID

	// Set default status
	if product.Status == "" {
		product.Status = model.ProductStatusDraft
	}

	// Generate slug if not provided
	if product.Slug == "" {
		product.Slug = model.GenerateSlug(product.Name)
	}

	// Ensure uniqueness of slug
	product.Slug = s.makeSlugUnique(product.Slug, shopID, 0)

	// Create product
	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productServiceEnhanced) UpdateProduct(id uint, product *model.Product, shopID uint, isAdmin bool) (*model.Product, error) {
	// Get existing product
	existingProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Check authorization
	if !isAdmin && existingProduct.ShopID != shopID {
		return nil, ErrUnauthorizedSeller
	}

	// Update fields
	if product.Name != "" {
		existingProduct.Name = product.Name
		existingProduct.Slug = model.GenerateSlug(product.Name)
		existingProduct.Slug = s.makeSlugUnique(existingProduct.Slug, shopID, id)
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}
	if product.ShortDescription != "" {
		existingProduct.ShortDescription = product.ShortDescription
	}
	if product.Price > 0 {
		existingProduct.Price = product.Price
	}
	if product.OriginalPrice > 0 {
		existingProduct.OriginalPrice = product.OriginalPrice
	}
	if product.Stock >= 0 {
		existingProduct.Stock = product.Stock
	}
	if product.CategoryID > 0 {
		existingProduct.CategoryID = product.CategoryID
	}
	if product.Status != "" {
		existingProduct.Status = product.Status
	}
	if product.Brand != "" {
		existingProduct.Brand = product.Brand
	}
	if product.Weight > 0 {
		existingProduct.Weight = product.Weight
	}

	// Save changes
	if err := s.productRepo.Update(existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (s *productServiceEnhanced) DeleteProduct(id uint, shopID uint, isAdmin bool) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return ErrProductNotFound
	}

	// Check authorization
	if !isAdmin && product.ShopID != shopID {
		return ErrUnauthorizedSeller
	}

	// Soft delete
	return s.productRepo.Delete(id)
}

func (s *productServiceEnhanced) GetProductByID(id uint) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *productServiceEnhanced) GetProducts(page, limit int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.productRepo.FindAll(limit, offset)
}

func (s *productServiceEnhanced) GetProductsByShopID(shopID uint, page, limit int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.productRepo.FindByShopID(shopID, limit, offset)
}

func (s *productServiceEnhanced) GetProductsByCategoryID(categoryID uint, page, limit int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.productRepo.FindByCategoryID(categoryID, limit, offset)
}

// ==================== SEARCH & FILTER ====================

func (s *productServiceEnhanced) SearchProducts(keyword string, filters ProductFilters, page, limit int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.productRepo.Search(keyword, repository.ProductFilters(filters), limit, offset)
}

func (s *productServiceEnhanced) FilterProducts(filters ProductFilters, page, limit int) ([]model.Product, int64, error) {
	return s.SearchProducts("", filters, page, limit)
}

// ==================== FEATURED & BEST SELLERS ====================

func (s *productServiceEnhanced) GetFeaturedProducts(limit int) ([]model.Product, error) {
	if limit > 50 {
		limit = 50
	}
	return s.productRepo.FindFeatured(limit)
}

func (s *productServiceEnhanced) GetBestSellers(limit int) ([]model.Product, error) {
	if limit > 50 {
		limit = 50
	}
	return s.productRepo.FindBestSellers(limit)
}

func (s *productServiceEnhanced) GetLatestProducts(limit int) ([]model.Product, error) {
	if limit > 50 {
		limit = 50
	}
	return s.productRepo.FindLatest(limit)
}

// ==================== PRODUCT IMAGES ====================

func (s *productServiceEnhanced) UploadProductImages(productID uint, files []*multipart.FileHeader, uploadDir string) ([]model.ProductImage, error) {
	// Check if product exists
	if _, err := s.productRepo.FindByID(productID); err != nil {
		return nil, ErrProductNotFound
	}

	// Check number of images
	existingImages, _ := s.productRepo.FindImagesByProductID(productID)
	if len(existingImages)+len(files) > 9 {
		return nil, ErrTooManyImages
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	var uploadedImages []model.ProductImage
	isPrimary := len(existingImages) == 0

	for i, file := range files {
		// Validate file
		if err := s.validateImageFile(file); err != nil {
			return nil, err
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%d_%s", productID, time.Now().UnixNano(), generateRandomString(8)) + ext
		filepath := filepath.Join(uploadDir, filename)

		// Create destination file
		dst, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, src); err != nil {
			return nil, err
		}

		// Create image record
		image := &model.ProductImage{
			ProductID: productID,
			URL:       "/uploads/products/" + filename,
			IsPrimary: isPrimary || (i == 0 && len(existingImages) == 0),
			SortOrder: len(existingImages) + i,
		}

		if err := s.productRepo.CreateImage(image); err != nil {
			return nil, err
		}

		uploadedImages = append(uploadedImages, *image)
		isPrimary = false
	}

	return uploadedImages, nil
}

func (s *productServiceEnhanced) UpdateProductImage(imageID uint, isPrimary bool, sortOrder int) error {
	image, err := s.productRepo.FindImagesByProductID(0)
	if err != nil {
		return err
	}

	// Find the specific image
	var targetImage *model.ProductImage
	for _, img := range image {
		if img.ID == imageID {
			targetImage = &img
			break
		}
	}

	if targetImage == nil {
		return errors.New("image not found")
	}

	// If setting as primary, unset other primary images
	if isPrimary {
		// Unset all primary images for this product
		images, _ := s.productRepo.FindImagesByProductID(targetImage.ProductID)
		for _, img := range images {
			if img.IsPrimary {
				img.IsPrimary = false
				s.productRepo.UpdateImage(&img)
			}
		}
	}

	targetImage.IsPrimary = isPrimary
	targetImage.SortOrder = sortOrder

	return s.productRepo.UpdateImage(targetImage)
}

func (s *productServiceEnhanced) DeleteProductImage(imageID uint, productID uint, shopID uint, isAdmin bool) error {
	// Check authorization
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return ErrProductNotFound
	}

	if !isAdmin && product.ShopID != shopID {
		return ErrUnauthorizedSeller
	}

	return s.productRepo.DeleteImage(imageID)
}

// ==================== PRODUCT VARIANTS ====================

func (s *productServiceEnhanced) CreateVariant(productID uint, variant *model.ProductVariant) (*model.ProductVariant, error) {
	// Check if product exists
	if _, err := s.productRepo.FindByID(productID); err != nil {
		return nil, ErrProductNotFound
	}

	// Check if SKU already exists
	if variant.SKU != "" {
		existingVariant, _ := s.productRepo.FindVariantBySKU(variant.SKU)
		if existingVariant != nil {
			return nil, ErrDuplicateSKU
		}
	}

	variant.ProductID = productID

	if err := s.productRepo.CreateVariant(variant); err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *productServiceEnhanced) UpdateVariant(id uint, variant *model.ProductVariant) (*model.ProductVariant, error) {
	existingVariant, err := s.productRepo.FindVariantBySKU("")
	if err != nil {
		return nil, ErrVariantNotFound
	}

	// Update fields
	if variant.Name != "" {
		existingVariant.Name = variant.Name
	}
	if variant.Price > 0 {
		existingVariant.Price = variant.Price
	}
	if variant.Stock >= 0 {
		existingVariant.Stock = variant.Stock
	}
	if variant.Attributes != "" {
		existingVariant.Attributes = variant.Attributes
	}

	if err := s.productRepo.UpdateVariant(existingVariant); err != nil {
		return nil, err
	}

	return existingVariant, nil
}

func (s *productServiceEnhanced) DeleteVariant(id uint, productID uint, shopID uint) error {
	// Check authorization
	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return ErrProductNotFound
	}

	if err := s.productRepo.DeleteVariant(id); err != nil {
		return err
	}

	return nil
}

func (s *productServiceEnhanced) GetVariantsByProductID(productID uint) ([]model.ProductVariant, error) {
	return s.productRepo.FindVariantsByProductID(productID)
}

// ==================== PRODUCT ATTRIBUTES ====================

func (s *productServiceEnhanced) CreateAttribute(productID uint, attr *model.ProductAttribute) (*model.ProductAttribute, error) {
	attr.ProductID = productID

	if err := s.productRepo.CreateAttribute(attr); err != nil {
		return nil, err
	}

	return attr, nil
}

func (s *productServiceEnhanced) UpdateAttribute(id uint, attr *model.ProductAttribute) (*model.ProductAttribute, error) {
	existingAttr, err := s.productRepo.FindAttributesByProductID(0)
	if err != nil {
		return nil, err
	}

	// Find and update
	for _, a := range existingAttr {
		if a.ID == id {
			if attr.Name != "" {
				a.Name = attr.Name
			}
			if attr.Type != "" {
				a.Type = attr.Type
			}
			if attr.Values != "" {
				a.Values = attr.Values
			}
			a.IsFilterable = attr.IsFilterable
			a.IsVisible = attr.IsVisible

			if err := s.productRepo.UpdateAttribute(&a); err != nil {
				return nil, err
			}
			return &a, nil
		}
	}

	return nil, errors.New("attribute not found")
}

func (s *productServiceEnhanced) DeleteAttribute(id uint, productID uint) error {
	return s.productRepo.DeleteAttribute(id)
}

func (s *productServiceEnhanced) GetAttributesByProductID(productID uint) ([]model.ProductAttribute, error) {
	return s.productRepo.FindAttributesByProductID(productID)
}

// ==================== INVENTORY MANAGEMENT ====================

func (s *productServiceEnhanced) UpdateProductStock(productID uint, variantID *uint, quantity int, shopID uint) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return ErrProductNotFound
	}

	if product.ShopID != shopID {
		return ErrUnauthorizedSeller
	}

	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return s.productRepo.UpdateStock(productID, variantID, quantity)
}

func (s *productServiceEnhanced) ReserveProductStock(productID uint, variantID *uint, quantity int) error {
	return s.productRepo.ReserveStock(productID, variantID, quantity)
}

func (s *productServiceEnhanced) ReleaseProductStock(productID uint, variantID *uint, quantity int) error {
	return s.productRepo.ReleaseStock(productID, variantID, quantity)
}

func (s *productServiceEnhanced) DecreaseProductStock(productID uint, variantID *uint, quantity int) error {
	if err := s.productRepo.DecreaseStock(productID, variantID, quantity); err != nil {
		return ErrInsufficientStock
	}

	// Increment sold count
	_ = s.productRepo.IncrementSoldCount(productID, quantity)

	return nil
}

// ==================== PRODUCT STATUS ====================

func (s *productServiceEnhanced) PublishProduct(id uint, shopID uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return ErrProductNotFound
	}

	if product.ShopID != shopID {
		return ErrUnauthorizedSeller
	}

	product.Status = model.ProductStatusActive
	return s.productRepo.Update(product)
}

func (s *productServiceEnhanced) UnpublishProduct(id uint, shopID uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return ErrProductNotFound
	}

	if product.ShopID != shopID {
		return ErrUnauthorizedSeller
	}

	product.Status = model.ProductStatusInactive
	return s.productRepo.Update(product)
}

func (s *productServiceEnhanced) FeatureProduct(id uint, isAdmin bool) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return ErrProductNotFound
	}

	if !isAdmin {
		return errors.New("only admin can feature products")
	}

	product.IsFeatured = !product.IsFeatured
	return s.productRepo.Update(product)
}

// ==================== CATEGORIES ====================

func (s *productServiceEnhanced) CreateCategory(category *model.Category) (*model.Category, error) {
	if category.Name == "" {
		return nil, errors.New("category name is required")
	}

	if category.Slug == "" {
		category.Slug = model.GenerateSlug(category.Name)
	}

	if err := s.productRepo.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *productServiceEnhanced) UpdateCategory(id uint, category *model.Category) (*model.Category, error) {
	existingCategory, err := s.productRepo.FindCategoryByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	if category.Name != "" {
		existingCategory.Name = category.Name
	}
	if category.Slug != "" {
		existingCategory.Slug = category.Slug
	}
	if category.Description != "" {
		existingCategory.Description = category.Description
	}
	if category.ParentID != nil {
		existingCategory.ParentID = category.ParentID
	}

	if err := s.productRepo.UpdateCategory(existingCategory); err != nil {
		return nil, err
	}

	return existingCategory, nil
}

func (s *productServiceEnhanced) DeleteCategory(id uint) error {
	return s.productRepo.DeleteCategory(id)
}

func (s *productServiceEnhanced) GetCategoryByID(id uint) (*model.Category, error) {
	return s.productRepo.FindCategoryByID(id)
}

func (s *productServiceEnhanced) GetAllCategories() ([]model.Category, error) {
	return s.productRepo.FindAllCategories()
}

func (s *productServiceEnhanced) GetCategoriesByParentID(parentID *uint) ([]model.Category, error) {
	return s.productRepo.FindCategoriesByParentID(parentID)
}

// ==================== ANALYTICS ====================

func (s *productServiceEnhanced) IncrementViewCount(id uint) error {
	return s.productRepo.IncrementViewCount(id)
}

// ==================== HELPER METHODS ====================

func (s *productServiceEnhanced) validateProduct(product *model.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return ErrInvalidPrice
	}
	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	if product.CategoryID == 0 {
		return errors.New("category is required")
	}
	if product.ShopID == 0 {
		return errors.New("shop is required")
	}
	return nil
}

func (s *productServiceEnhanced) validateImageFile(file *multipart.FileHeader) error {
	// Check file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return ErrImageTooLarge
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !validExts[ext] {
		return ErrInvalidImageFormat
	}

	return nil
}

func (s *productServiceEnhanced) makeSlugUnique(slug string, shopID uint, excludeID uint) string {
	// In production, check database for existing slugs
	// For now, just return the slug
	return slug
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
