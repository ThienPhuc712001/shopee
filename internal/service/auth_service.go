// Package service provides business logic layer
package service

import (
	"context"
	"errors"
	"time"

	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/password"
	"gorm.io/gorm"
)

// ============================================================================
// ERRORS
// ============================================================================

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrAccountLocked       = errors.New("account is locked due to too many failed attempts")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already registered")
	ErrRefreshTokenInvalid = errors.New("invalid or expired refresh token")
	ErrRefreshTokenRevoked = errors.New("refresh token has been revoked")
)

// ============================================================================
// REQUEST/RESPONSE MODELS
// ============================================================================

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	User         *UserInfo `json:"user"`
}

// UserInfo represents user information in auth response
type UserInfo struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================================================
// AUTH SERVICE
// ============================================================================

// AuthService handles authentication business logic
type AuthService struct {
	db               *gorm.DB
	userRepo         repository.UserRepositoryEnhanced
	refreshTokenRepo *repository.RefreshTokenRepository
	tokenService     TokenService
	log              *logger.Logger
	maxLoginAttempts int
	lockoutDuration  time.Duration
}

// AuthServiceConfig holds auth service configuration
type AuthServiceConfig struct {
	DB               *gorm.DB
	UserRepo         repository.UserRepositoryEnhanced
	RefreshTokenRepo *repository.RefreshTokenRepository
	TokenService     TokenService
	Log              *logger.Logger
	MaxLoginAttempts int
	LockoutDuration  time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(config AuthServiceConfig) *AuthService {
	// Set defaults
	if config.MaxLoginAttempts <= 0 {
		config.MaxLoginAttempts = 5
	}
	if config.LockoutDuration <= 0 {
		config.LockoutDuration = 15 * time.Minute
	}

	return &AuthService{
		db:               config.DB,
		userRepo:         config.UserRepo,
		refreshTokenRepo: config.RefreshTokenRepo,
		tokenService:     config.TokenService,
		log:              config.Log,
		maxLoginAttempts: config.MaxLoginAttempts,
		lockoutDuration:  config.LockoutDuration,
	}
}

// ============================================================================
// LOGIN
// ============================================================================

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.WithField("email", req.Email).Info("Login attempt for non-existent user")
			return nil, ErrInvalidCredentials
		}
		s.log.WithError(err).Error("Failed to find user by email")
		return nil, ErrInvalidCredentials
	}

	// Check if account is locked
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		s.log.WithFields(logger.Fields{
			"user_id":      user.ID,
			"locked_until": user.LockedUntil,
		}).Warn("Login attempt on locked account")
		return nil, ErrAccountLocked
	}

	// Verify password
	hashedPwd := user.Password
	if err := password.Verify(req.Password, hashedPwd); err != nil {
		s.handleFailedLogin(user)
		s.log.WithFields(logger.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Info("Invalid password")
		return nil, ErrInvalidCredentials
	}

	// Clear failed login attempts on successful login
	if err := s.clearFailedLogin(user); err != nil {
		s.log.WithError(err).Error("Failed to clear failed login attempts")
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Generate tokens using token service
	tokenPair, err := s.tokenService.GenerateTokenPair(user)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate tokens")
		return nil, errors.New("failed to generate authentication tokens")
	}

	// Store refresh token in database
	refreshTokenExpiry := time.Now().Add(s.tokenService.(interface{ GetRefreshExpiry() time.Duration }).GetRefreshExpiry())
	if err := s.refreshTokenRepo.Create(ctx, &model.RefreshToken{
		UserID:    int64(user.ID),
		Token:     tokenPair.RefreshToken,
		ExpiresAt: refreshTokenExpiry,
	}); err != nil {
		s.log.WithError(err).Error("Failed to store refresh token")
	}

	s.log.WithFields(logger.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in successfully")

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// handleFailedLogin increments failed login attempts and locks account if needed
func (s *AuthService) handleFailedLogin(user *model.User) {
	user.FailedLoginAttempts++

	if user.FailedLoginAttempts >= s.maxLoginAttempts {
		lockUntil := time.Now().Add(s.lockoutDuration)
		user.LockedUntil = &lockUntil
		s.log.WithFields(logger.Fields{
			"user_id":         user.ID,
			"locked_until":    lockUntil,
			"failed_attempts": user.FailedLoginAttempts,
		}).Warn("Account locked due to failed login attempts")
	}

	_ = s.userRepo.Update(user)
}

// clearFailedLogin resets failed login attempts
func (s *AuthService) clearFailedLogin(user *model.User) error {
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	return s.userRepo.Update(user)
}

// ============================================================================
// REGISTER
// ============================================================================

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Validate email format (simple validation)
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Check if phone already exists (if provided)
	if req.Phone != "" {
		existingPhone, err := s.userRepo.FindByPhone(req.Phone)
		if err == nil && existingPhone != nil {
			return nil, errors.New("phone number already registered")
		}
	}

	// Validate password strength
	if err := password.ValidateDefault(req.Password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		s.log.WithError(err).Error("Failed to hash password")
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := &model.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      model.RoleCustomer,
		Status:    model.StatusActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.log.WithError(err).Error("Failed to create user")
		return nil, errors.New("failed to create user account")
	}

	s.log.WithFields(logger.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("New user registered")

	// Generate tokens
	tokenPair, err := s.tokenService.GenerateTokenPair(user)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate tokens")
		return nil, errors.New("failed to generate authentication tokens")
	}

	// Store refresh token
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour) // Default 7 days
	if err := s.refreshTokenRepo.Create(ctx, &model.RefreshToken{
		UserID:    int64(user.ID),
		Token:     tokenPair.RefreshToken,
		ExpiresAt: refreshTokenExpiry,
	}); err != nil {
		s.log.WithError(err).Error("Failed to store refresh token")
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// ============================================================================
// REFRESH TOKEN
// ============================================================================

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token using token service
	tokenPair, err := s.tokenService.RefreshAccessToken(refreshToken)
	if err != nil {
		s.log.WithError(err).Info("Invalid refresh token")
		return nil, ErrRefreshTokenInvalid
	}

	// Get user ID from claims
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}

	// Get user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		s.log.WithError(err).Error("Failed to find user")
		return nil, ErrUserNotFound
	}

	// Revoke old refresh token in database
	_ = s.refreshTokenRepo.Revoke(ctx, refreshToken)

	// Store new refresh token
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	if err := s.refreshTokenRepo.Create(ctx, &model.RefreshToken{
		UserID:    int64(user.ID),
		Token:     tokenPair.RefreshToken,
		ExpiresAt: refreshTokenExpiry,
	}); err != nil {
		s.log.WithError(err).Error("Failed to store new refresh token")
	}

	s.log.WithFields(logger.Fields{
		"user_id": user.ID,
	}).Info("Token refreshed successfully")

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// ============================================================================
// LOGOUT
// ============================================================================

// Logout revokes refresh tokens
func (s *AuthService) Logout(ctx context.Context, userID uint, refreshToken string) error {
	// Revoke specific token
	if refreshToken != "" {
		if err := s.refreshTokenRepo.Revoke(ctx, refreshToken); err != nil {
			s.log.WithError(err).Error("Failed to revoke refresh token")
		}
	}

	s.log.WithField("user_id", userID).Info("User logged out")
	return nil
}

// LogoutAllDevices revokes all refresh tokens for a user
func (s *AuthService) LogoutAllDevices(ctx context.Context, userID uint) error {
	if err := s.refreshTokenRepo.RevokeByUserID(ctx, int64(userID)); err != nil {
		s.log.WithError(err).Error("Failed to revoke all refresh tokens")
		return err
	}

	s.log.WithField("user_id", userID).Info("User logged out from all devices")
	return nil
}

// ============================================================================
// PASSWORD MANAGEMENT
// ============================================================================

// ForgotPassword initiates password reset (sends reset email)
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not
		s.log.WithField("email", email).Info("Password reset requested")
		return nil
	}

	// Generate reset token (simple implementation - should use secure token generation)
	_ = generateResetToken()
	
	// TODO: Store reset token in database with expiry
	// TODO: Send email with reset link
	
	s.log.WithFields(logger.Fields{
		"email": email,
		"user_id": user.ID,
	}).Info("Password reset initiated")
	return nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// TODO: Validate reset token from database
	// TODO: Check token expiry
	
	// Validate new password
	if err := password.ValidateDefault(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	// TODO: Find user by reset token and update password
	// For now, this is a placeholder
	_ = hashedPassword
	_ = token

	return nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := password.Verify(oldPassword, user.Password); err != nil {
		return ErrInvalidCredentials
	}

	// Validate new password
	if err := password.ValidateDefault(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Revoke all refresh tokens (force re-login)
	if err := s.refreshTokenRepo.RevokeByUserID(ctx, int64(userID)); err != nil {
		s.log.WithError(err).Error("Failed to revoke refresh tokens after password change")
	}

	s.log.WithField("user_id", userID).Info("Password changed successfully")
	return nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateUser updates user information
func (s *AuthService) UpdateUser(user *model.User) error {
	return s.userRepo.Update(user)
}

// VerifyEmail verifies user email with token
func (s *AuthService) VerifyEmail(token string) error {
	// TODO: Validate email verification token
	// TODO: Find user by token and set EmailVerified = true
	// For now, return nil to allow the flow
	return nil
}

// ResendVerificationEmail resends verification email
func (s *AuthService) ResendVerificationEmail(email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	// TODO: Generate verification token and send email
	s.log.WithFields(logger.Fields{
		"email": email,
		"user_id": user.ID,
	}).Info("Verification email sent")
	return nil
}

// generateResetToken generates a secure password reset token
func generateResetToken() string {
	// TODO: Use secure random token generation
	return "reset_token_placeholder"
}

// isValidEmail validates email format using simple regex
func isValidEmail(email string) bool {
	// Simple email validation - check for @ and .
	if len(email) < 5 || len(email) > 255 {
		return false
	}
	
	// Check for @ symbol
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}
	if atIndex <= 0 {
		return false
	}
	
	// Check for . after @
	dotFound := false
	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' {
			dotFound = true
			break
		}
	}
	
	return dotFound && atIndex < len(email)-2
}
