package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ShopService handles shop-related business logic
type ShopService interface {
	CreateShop(userID uint, input *CreateShopInput) (*model.Shop, error)
	GetShopByID(id uint) (*model.Shop, error)
	GetShopByUserID(userID uint) (*model.Shop, error)
	UpdateShop(id uint, input *UpdateShopInput) (*model.Shop, error)
}

type shopService struct {
	repo repository.ShopRepositoryEnhanced
}

type CreateShopInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Logo        string `json:"logo"`
	CoverImage  string `json:"cover_image"`
	Address     string `json:"address"`
}

type UpdateShopInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Logo        string `json:"logo"`
	CoverImage  string `json:"cover_image"`
	Address     string `json:"address"`
}

// NewShopService creates a new shop service
func NewShopService(repo repository.ShopRepositoryEnhanced) ShopService {
	return &shopService{repo: repo}
}

func (s *shopService) CreateShop(userID uint, input *CreateShopInput) (*model.Shop, error) {
	// Check if user already has a shop
	existingShop, _ := s.repo.FindByUserID(userID)
	if existingShop != nil {
		return nil, errors.New("user already has a shop")
	}

	// Generate slug from shop name
	slug := generateSlug(input.Name)

	shop := &model.Shop{
		UserID:      userID,
		Name:        input.Name,
		Slug:        slug,
		Description: input.Description,
		Phone:       input.Phone,
		Email:       input.Email,
		Logo:        input.Logo,
		CoverImage:  input.CoverImage,
		Address:     input.Address,
		Status:      model.ShopStatusActive,
	}

	if err := s.repo.Create(shop); err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *shopService) GetShopByID(id uint) (*model.Shop, error) {
	return s.repo.FindByID(id)
}

func (s *shopService) GetShopByUserID(userID uint) (*model.Shop, error) {
	return s.repo.FindByUserID(userID)
}

func (s *shopService) UpdateShop(id uint, input *UpdateShopInput) (*model.Shop, error) {
	shop, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("shop not found")
	}

	// Update fields
	if input.Name != "" {
		shop.Name = input.Name
		shop.Slug = generateSlug(input.Name)
	}
	if input.Description != "" {
		shop.Description = input.Description
	}
	if input.Phone != "" {
		shop.Phone = input.Phone
	}
	if input.Email != "" {
		shop.Email = input.Email
	}
	if input.Logo != "" {
		shop.Logo = input.Logo
	}
	if input.CoverImage != "" {
		shop.CoverImage = input.CoverImage
	}
	if input.Address != "" {
		shop.Address = input.Address
	}

	if err := s.repo.Update(shop); err != nil {
		return nil, err
	}

	return shop, nil
}

// generateSlug creates a URL-friendly slug from a string
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	// Add timestamp to ensure uniqueness
	return fmt.Sprintf("%s-%d", slug, time.Now().UnixNano())
}
