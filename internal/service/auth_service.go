package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Auth errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrPhoneAlreadyExists = errors.New("phone number already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrAccountLocked      = errors.New("account is locked")
	ErrAccountInactive    = errors.New("account is inactive")
)

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Phone     string `json:"phone" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthService handles authentication business logic
type AuthService interface {
	Register(req *RegisterRequest) (*model.User, error)
	Login(req *LoginRequest) (*model.User, string, error)
	GetUserByID(id uint) (*model.User, error)
}

type authService struct {
	userRepo   repository.UserRepositoryEnhanced
	jwtSecret  string
	jwtExpiry  time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepositoryEnhanced, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *authService) Register(req *RegisterRequest) (*model.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create new user
	user := &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      model.RoleCustomer,
		Status:    model.StatusActive,
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *LoginRequest) (*model.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, "", ErrInvalidCredentials
	}

	// Check user status
	if user.Status != model.StatusActive {
		return nil, "", ErrAccountInactive
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *authService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
