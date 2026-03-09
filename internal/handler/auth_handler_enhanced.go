package handler

import (
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

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

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account with email, password, and phone
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	user, err := h.authService.Register(&service.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, response.Error("Email already registered"))
			return
		}
		if err == service.ErrPhoneAlreadyExists {
			c.JSON(http.StatusConflict, response.Error("Phone number already registered"))
			return
		}
		// Handle password validation errors
		if strings.Contains(err.Error(), "password") {
			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to register user"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"user": user,
	}, "User registered successfully"))
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 423 {object} response.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	user, token, err := h.authService.Login(&service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid email or password"))
			return
		}
		if err == service.ErrAccountLocked {
			c.JSON(http.StatusLocked, response.Error("Account is locked. Please try again later or contact support."))
			return
		}
		if err == service.ErrAccountInactive {
			c.JSON(http.StatusForbidden, response.Error("Account is inactive. Please contact support."))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to login"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"user":  user,
		"token": token,
	}, "Login successful"))
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/refresh [post]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate refresh token and logout user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response
// @Router /api/auth/logout [post]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// Me handles get current user
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized("User not found"))
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("User not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"user": gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"first_name":     user.FirstName,
			"last_name":      user.LastName,
			"phone":          user.Phone,
			"avatar":         user.Avatar,
			"role":           user.Role,
			"email_verified": user.EmailVerified,
			"created_at":     user.CreatedAt,
		},
	}, ""))
}

// UpdateProfile handles updating user profile
// @Summary Update profile
// @Description Update current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile data"
// @Success 200 {object} response.Response
// @Router /api/auth/profile [put]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// ChangePassword handles changing user password
// @Summary Change password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/change-password [post]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// RequestPasswordReset handles password reset request
// @Summary Request password reset
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email"
// @Success 200 {object} response.Response
// @Router /api/auth/forgot-password [post]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPassword handles password reset with token
// @Summary Reset password
// @Description Reset password with token from email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset data"
// @Success 200 {object} response.Response
// @Router /api/auth/reset-password [post]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// VerifyEmail handles email verification
// @Summary Verify email
// @Description Verify user email with token
// @Tags auth
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.Response
// @Router /api/auth/verify-email [get]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// ResendVerificationEmail handles resending verification email
// @Summary Resend verification email
// @Description Resend email verification link
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendVerificationRequest true "Email"
// @Success 200 {object} response.Response
// @Router /api/auth/resend-verification [post]
// Note: This endpoint is not implemented in the current service version
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, response.Error("This endpoint is not implemented"))
}

// ResendVerificationRequest represents a resend verification request
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUserIDFromParam extracts user ID from URL parameter
func GetUserIDFromParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
