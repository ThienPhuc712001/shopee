package repository

import (
	"ecommerce/internal/domain/model"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(user *model.User) error

	// Update updates an existing user
	Update(user *model.User) error

	// Delete deletes a user by ID
	Delete(id uint) error

	// FindByID finds a user by ID
	FindByID(id uint) (*model.User, error)

	// FindByEmail finds a user by email
	FindByEmail(email string) (*model.User, error)

	// FindAll retrieves all users with pagination
	FindAll(limit, offset int) ([]model.User, int64, error)

	// Search searches users by keyword
	Search(keyword string, limit, offset int) ([]model.User, int64, error)

	// UpdatePassword updates user password
	UpdatePassword(id uint, hashedPassword string) error

	// UpdateLastLogin updates user's last login time
	UpdateLastLogin(id uint) error

	// UpdateStatus updates user status
	UpdateStatus(id uint, status model.UserStatus) error
}
