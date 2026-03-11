package repository

import (
	"context"
	"ecommerce/internal/domain/model"
	"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// FindAll finds all users with pagination
func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByRole finds users by role
func (r *UserRepository) FindByRole(ctx context.Context, role string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.WithContext(ctx).Where("role = ?", role).Find(&users).Error
	return users, err
}

// UpdateLastLogin updates the last login time
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("last_login", gorm.Expr("GETDATE()")).Error
}

// UpdatePassword updates the password
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// UpdateStatus updates the user status
func (r *UserRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// ExistsByEmail checks if a user with the given email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// Search searches users by keyword
func (r *UserRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	searchPattern := "%" + keyword + "%"

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email LIKE ? OR first_name LIKE ? OR last_name LIKE ?", searchPattern, searchPattern, searchPattern).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Where("email LIKE ? OR first_name LIKE ? OR last_name LIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
