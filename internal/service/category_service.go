package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// CategoryService handles category business logic
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategoryInput represents input for creating a category
type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	ParentID    *uint  `json:"parent_id"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	ImageURL    string `json:"image_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
	Attributes  string `json:"attributes"`
}

// UpdateCategoryInput represents input for updating a category
type UpdateCategoryInput struct {
	Name        string `json:"name"`
	ParentID    *uint  `json:"parent_id"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	ImageURL    string `json:"image_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
	Attributes  string `json:"attributes"`
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, input CreateCategoryInput) (*model.Category, error) {
	// Validate input
	if err := s.validateCategoryInput(input); err != nil {
		return nil, err
	}

	// Check if parent category exists
	if input.ParentID != nil {
		parent, err := s.repo.GetCategoryByID(ctx, *input.ParentID)
		if err != nil {
			return nil, ErrParentCategoryNotFound
		}
		if parent.DeletedAt.Valid {
			return nil, ErrParentCategoryDeleted
		}
	}

	// Generate unique slug
	slug := s.generateSlug(input.Name)

	// Calculate level
	level := 0
	if input.ParentID != nil {
		parent, _ := s.repo.GetCategoryByID(ctx, *input.ParentID)
		if parent != nil {
			level = parent.Level + 1
		}
	}

	// Create category
	category := &model.Category{
		Name:        input.Name,
		ParentID:    input.ParentID,
		Description: input.Description,
		IconURL:     input.IconURL,
		ImageURL:    input.ImageURL,
		Slug:        slug,
		SortOrder:   input.SortOrder,
		IsActive:    input.IsActive,
		Attributes:  input.Attributes,
		Level:       level,
	}

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, input UpdateCategoryInput) (*model.Category, error) {
	// Get existing category
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// Validate parent if changed
	if input.ParentID != nil && category.ParentID != nil && *input.ParentID != *category.ParentID {
		// Prevent circular hierarchy
		if *input.ParentID == id {
			return nil, ErrCircularHierarchy
		}

		// Check if new parent is a descendant
		isDescendant, err := s.isDescendant(ctx, *input.ParentID, id)
		if err != nil {
			return nil, err
		}
		if isDescendant {
			return nil, ErrCircularHierarchy
		}

		// Check if parent exists
		parent, err := s.repo.GetCategoryByID(ctx, *input.ParentID)
		if err != nil {
			return nil, ErrParentCategoryNotFound
		}
		if parent.DeletedAt.Valid {
			return nil, ErrParentCategoryDeleted
		}
	}

	// Update fields
	if input.Name != "" {
		category.Name = input.Name
		category.Slug = s.generateSlug(input.Name)
	}
	if input.Description != "" {
		category.Description = input.Description
	}
	if input.IconURL != "" {
		category.IconURL = input.IconURL
	}
	if input.ImageURL != "" {
		category.ImageURL = input.ImageURL
	}
	if input.SortOrder != 0 {
		category.SortOrder = input.SortOrder
	}
	category.IsActive = input.IsActive
	if input.Attributes != "" {
		category.Attributes = input.Attributes
	}
	if input.ParentID != nil {
		category.ParentID = input.ParentID
		// Update level
		parent, _ := s.repo.GetCategoryByID(ctx, *input.ParentID)
		if parent != nil {
			category.Level = parent.Level + 1
		}
	}

	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	// Check if category exists
	_, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	// Check if has children
	hasChildren, err := s.repo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return ErrCategoryHasChildren
	}

	// Check if has products
	hasProducts, err := s.repo.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return ErrCategoryHasProducts
	}

	return s.repo.DeleteCategory(ctx, id)
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*model.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error) {
	category, err := s.repo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories(ctx context.Context, isActive *bool) ([]model.Category, error) {
	return s.repo.GetAllCategories(ctx, isActive)
}

// GetCategoryTree retrieves categories in tree structure
func (s *CategoryService) GetCategoryTree(ctx context.Context) ([]model.CategoryTree, error) {
	categories, err := s.repo.GetCategoryTree(ctx)
	if err != nil {
		return nil, err
	}

	return s.buildCategoryTree(categories), nil
}

// GetCategoryTreeWithCounts retrieves categories with product counts
func (s *CategoryService) GetCategoryTreeWithCounts(ctx context.Context) ([]model.CategoryTree, error) {
	categoriesWithCount, err := s.repo.GetAllCategoriesWithProductCount(ctx)
	if err != nil {
		return nil, err
	}

	return s.buildCategoryTreeWithCounts(categoriesWithCount), nil
}

// GetProductsByCategory retrieves products in a category
func (s *CategoryService) GetProductsByCategory(ctx context.Context, categoryID uint, page, limit int) ([]model.Product, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetProductsByCategory(ctx, categoryID, limit, offset)
}

// GetCategoryPath retrieves breadcrumb path for a category
func (s *CategoryService) GetCategoryPath(ctx context.Context, id uint) ([]model.CategoryBreadcrumb, error) {
	path, err := s.repo.GetCategoryPath(ctx, id)
	if err != nil {
		return nil, err
	}

	breadcrumbs := make([]model.CategoryBreadcrumb, len(path))
	for i, cat := range path {
		breadcrumbs[i] = model.CategoryBreadcrumb{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		}
	}

	return breadcrumbs, nil
}

// GetFeaturedCategories retrieves featured categories for homepage
func (s *CategoryService) GetFeaturedCategories(ctx context.Context, limit int) ([]model.CategoryWithCount, error) {
	return s.repo.GetFeaturedCategories(ctx, limit)
}

// SearchCategories searches categories
func (s *CategoryService) SearchCategories(ctx context.Context, query string) ([]model.Category, error) {
	return s.repo.SearchCategories(ctx, query)
}

// ========== Helper Methods ==========

// validateCategoryInput validates category input
func (s *CategoryService) validateCategoryInput(input CreateCategoryInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return ErrCategoryNameRequired
	}
	if len(input.Name) < 2 || len(input.Name) > 255 {
		return ErrCategoryNameLength
	}
	return nil
}

// generateSlug generates a URL-friendly slug from name
func (s *CategoryService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "--", "-")
	return fmt.Sprintf("%s-%s", slug, uuid.New().String()[:6])
}

// ensureUniqueSlug ensures slug is unique (simplified - just returns nil for now)
func (s *CategoryService) ensureUniqueSlug(ctx context.Context, slug string, id uint) error {
	return nil
}

// buildCategoryTree builds a tree structure from flat categories
func (s *CategoryService) buildCategoryTree(categories []model.Category) []model.CategoryTree {
	tree := make([]model.CategoryTree, 0)

	for _, cat := range categories {
		node := model.CategoryTree{
			ID:           cat.ID,
			ParentID:     cat.ParentID,
			Name:         cat.Name,
			Slug:         cat.Slug,
			Description:  cat.Description,
			IconURL:      cat.IconURL,
			ImageURL:     cat.ImageURL,
			Level:        cat.Level,
			SortOrder:    cat.SortOrder,
			IsActive:     cat.IsActive,
			ProductCount: 0,
			Children:     s.buildCategoryTree(cat.Children),
		}
		tree = append(tree, node)
	}

	return tree
}

// buildCategoryTreeWithCounts builds tree with product counts
func (s *CategoryService) buildCategoryTreeWithCounts(categories []model.CategoryWithCount) []model.CategoryTree {
	// Group by parent
	childrenMap := make(map[uint][]model.CategoryWithCount)
	var rootIDs []uint

	for _, cat := range categories {
		if cat.ParentID == nil {
			rootIDs = append(rootIDs, cat.ID)
		} else {
			childrenMap[*cat.ParentID] = append(childrenMap[*cat.ParentID], cat)
		}
	}

	// Build tree recursively
	tree := make([]model.CategoryTree, 0)
	for _, id := range rootIDs {
		var cat model.CategoryWithCount
		for _, c := range categories {
			if c.ID == id {
				cat = c
				break
			}
		}

		node := model.CategoryTree{
			ID:           cat.ID,
			ParentID:     cat.ParentID,
			Name:         cat.Name,
			Slug:         cat.Slug,
			Description:  cat.Description,
			IconURL:      cat.IconURL,
			ImageURL:     cat.ImageURL,
			Level:        cat.Level,
			SortOrder:    cat.SortOrder,
			IsActive:     cat.IsActive,
			ProductCount: cat.ProductCount,
			Children:     s.buildChildrenTree(categories, childrenMap, cat.ID),
		}
		tree = append(tree, node)
	}

	return tree
}

// buildChildrenTree builds children for a category
func (s *CategoryService) buildChildrenTree(
	allCategories []model.CategoryWithCount,
	childrenMap map[uint][]model.CategoryWithCount,
	parentID uint,
) []model.CategoryTree {
	children := childrenMap[parentID]
	if len(children) == 0 {
		return nil
	}

	tree := make([]model.CategoryTree, 0, len(children))
	for _, cat := range children {
		node := model.CategoryTree{
			ID:           cat.ID,
			ParentID:     cat.ParentID,
			Name:         cat.Name,
			Slug:         cat.Slug,
			Description:  cat.Description,
			IconURL:      cat.IconURL,
			ImageURL:     cat.ImageURL,
			Level:        cat.Level,
			SortOrder:    cat.SortOrder,
			IsActive:     cat.IsActive,
			ProductCount: cat.ProductCount,
			Children:     s.buildChildrenTree(allCategories, childrenMap, cat.ID),
		}
		tree = append(tree, node)
	}

	return tree
}

// isDescendant checks if a category is a descendant of another
func (s *CategoryService) isDescendant(ctx context.Context, potentialDescendantID, potentialAncestorID uint) (bool, error) {
	current, err := s.repo.GetCategoryByID(ctx, potentialDescendantID)
	if err != nil {
		return false, err
	}

	for current.ParentID != nil {
		if *current.ParentID == potentialAncestorID {
			return true, nil
		}
		current, err = s.repo.GetCategoryByID(ctx, *current.ParentID)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

// ========== Error Definitions ==========

var (
	ErrParentCategoryNotFound = errors.New("parent category not found")
	ErrParentCategoryDeleted  = errors.New("parent category is deleted")
	ErrCategoryHasChildren    = errors.New("cannot delete category with subcategories")
	ErrCategoryHasProducts    = errors.New("cannot delete category with products")
	ErrCircularHierarchy      = errors.New("circular category hierarchy detected")
	ErrCategoryNameRequired   = errors.New("category name is required")
	ErrCategoryNameLength     = errors.New("category name must be between 2 and 255 characters")
)
