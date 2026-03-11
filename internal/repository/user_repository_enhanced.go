package repository

import (
	"ecommerce/internal/domain/model"
	"gorm.io/gorm"
)

// UserRepositoryEnhanced defines the enhanced interface for user data operations
type UserRepositoryEnhanced interface {
	// Basic CRUD
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)

	// Queries
	FindAll(limit, offset int) ([]model.User, int64, error)
	Search(keyword string, limit, offset int) ([]model.User, int64, error)

	// Status
	UpdatePassword(id uint, hashedPassword string) error
	UpdateLastLogin(id uint) error
	UpdateStatus(id uint, status model.UserStatus) error
}

type userRepositoryEnhanced struct {
	db *gorm.DB
}

// NewUserRepositoryEnhanced creates a new enhanced user repository
func NewUserRepositoryEnhanced(db *gorm.DB) UserRepositoryEnhanced {
	return &userRepositoryEnhanced{db: db}
}

func (r *userRepositoryEnhanced) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryEnhanced) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepositoryEnhanced) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepositoryEnhanced) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Cart").Preload("Addresses").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryEnhanced) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryEnhanced) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryEnhanced) FindAll(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepositoryEnhanced) Search(keyword string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	searchPattern := "%" + keyword + "%"

	if err := r.db.Model(&model.User{}).
		Where("email LIKE ? OR first_name LIKE ? OR last_name LIKE ?",
			searchPattern, searchPattern, searchPattern).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("email LIKE ? OR first_name LIKE ? OR last_name LIKE ?",
		searchPattern, searchPattern, searchPattern).
		Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *userRepositoryEnhanced) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *userRepositoryEnhanced) UpdateLastLogin(id uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("last_login", gorm.Expr("GETDATE()")).Error
}

func (r *userRepositoryEnhanced) UpdateStatus(id uint, status model.UserStatus) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}
