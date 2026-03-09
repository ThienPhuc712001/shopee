package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole defines the role of a user
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleSeller   UserRole = "seller"
	RoleUserAdmin    UserRole = "admin" // Renamed to avoid conflict with AdminRoleType
)

// UserStatus defines the status of a user account
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
	StatusLocked   UserStatus = "locked"
)

// User represents a user in the system
type User struct {
	ID              uint       `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password        string     `gorm:"type:varchar(255);not null" json:"-"`
	FirstName       string     `gorm:"type:varchar(100)" json:"first_name"`
	LastName        string     `gorm:"type:varchar(100)" json:"last_name"`
	Phone           string     `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Avatar          string     `gorm:"type:varchar(500)" json:"avatar"`
	Role            UserRole   `gorm:"type:varchar(20);default:'customer';index" json:"role"`
	Status          UserStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`
	EmailVerified   bool       `gorm:"type:bit;default:false" json:"email_verified"`
	LastLogin       *time.Time `gorm:"type:datetime" json:"last_login"`
	FailedLoginAttempts int    `gorm:"type:int;default:0" json:"-"`
	LockedUntil     *time.Time `gorm:"type:datetime" json:"-"`
	RefreshToken    string     `gorm:"type:varchar(500)" json:"-"`
	RefreshTokenExpiry *time.Time `gorm:"type:datetime" json:"-"`

	// Relationships
	Shops     []Shop     `gorm:"foreignKey:UserID" json:"shops,omitempty"`
	Cart      *Cart      `gorm:"foreignKey:UserID" json:"cart,omitempty"`
	Orders    []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Reviews   []Review   `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
	Addresses []Address  `gorm:"foreignKey:UserID" json:"addresses,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &PasswordError{Message: "failed to hash password", Err: err}
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// LockAccount locks the user account for a specified duration
func (u *User) LockAccount(duration time.Duration) {
	until := time.Now().Add(duration)
	u.LockedUntil = &until
	u.Status = StatusLocked
}

// UnlockAccount unlocks the user account
func (u *User) UnlockAccount() {
	u.LockedUntil = nil
	u.FailedLoginAttempts = 0
	if u.Status == StatusLocked {
		u.Status = StatusActive
	}
}

// IncrementFailedLogin increments the failed login attempt counter
func (u *User) IncrementFailedLogin() {
	u.FailedLoginAttempts++
	// Lock account after 5 failed attempts for 30 minutes
	if u.FailedLoginAttempts >= 5 {
		u.LockAccount(30 * time.Minute)
	}
}

// ResetFailedLogin resets the failed login counter
func (u *User) ResetFailedLogin() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// HasAnyRole checks if user has any of the specified roles
func (u *User) HasAnyRole(roles ...UserRole) bool {
	for _, role := range roles {
		if u.Role == role {
			return true
		}
	}
	return false
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleUserAdmin
}

// IsSeller checks if user is a seller
func (u *User) IsSeller() bool {
	return u.Role == RoleSeller
}

// IsCustomer checks if user is a customer
func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}

// PasswordError represents a password-related error
type PasswordError struct {
	Message string
	Err     error
}

func (e *PasswordError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// ValidatePassword validates password strength
func ValidatePassword(password string) *PasswordError {
	if len(password) < 8 {
		return &PasswordError{Message: "password must be at least 8 characters"}
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' ||
			char == '%' || char == '^' || char == '&' || char == '*' ||
			char == '(' || char == ')' || char == '-' || char == '_':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return &PasswordError{Message: "password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return &PasswordError{Message: "password must contain at least one lowercase letter"}
	}
	if !hasNumber {
		return &PasswordError{Message: "password must contain at least one number"}
	}
	if !hasSpecial {
		return &PasswordError{Message: "password must contain at least one special character"}
	}

	return nil
}
