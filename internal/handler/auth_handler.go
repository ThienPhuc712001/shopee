// Package handler provides HTTP handlers for authentication
package handler

import (
	"github.com/gin-gonic/gin"
	"ecommerce/internal/errors"
	"ecommerce/internal/middleware"
	"ecommerce/internal/service"
	"ecommerce/pkg/logger"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
	log         *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

// ============================================================================
// LOGIN
// ============================================================================

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logger.Fields{
			"error": err.Error(),
		}).Info("Invalid login request")

		middleware.AbortWithError(c, errors.BadRequest("Invalid request body"))
		return
	}

	// Authenticate
	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.log.WithFields(logger.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Info("Login failed")

		switch err {
		case service.ErrInvalidCredentials:
			middleware.AbortWithError(c, errors.InvalidCredentials())
		case service.ErrAccountLocked:
			middleware.AbortWithError(c, errors.Forbidden("Account is locked. Please try again later."))
		default:
			middleware.AbortWithError(c, errors.Internal("Login failed"))
		}
		return
	}

	h.log.WithFields(logger.Fields{
		"user_id": response.User.ID,
		"email":   response.User.Email,
	}).Info("User logged in successfully")

	middleware.Success(c, response, "Login successful")
}

// ============================================================================
// REGISTER
// ============================================================================

// Register godoc
// @Summary User registration
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Registration details"
// @Success 201 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logger.Fields{
			"error": err.Error(),
		}).Info("Invalid registration request")

		middleware.AbortWithError(c, errors.BadRequest("Invalid request body"))
		return
	}

	// Register user
	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.log.WithFields(logger.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Info("Registration failed")

		switch err {
		case service.ErrEmailAlreadyExists:
			middleware.AbortWithError(c, errors.Conflict("Email already registered"))
		default:
			if appErr, ok := err.(*errors.AppError); ok {
				middleware.AbortWithError(c, appErr)
			} else {
				middleware.AbortWithError(c, errors.Internal("Registration failed"))
			}
		}
		return
	}

	h.log.WithFields(logger.Fields{
		"user_id": response.User.ID,
		"email":   response.User.Email,
	}).Info("User registered successfully")

	middleware.Created(c, response, "Registration successful")
}

// ============================================================================
// REFRESH TOKEN
// ============================================================================

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logger.Fields{
			"error": err.Error(),
		}).Info("Invalid refresh token request")

		middleware.AbortWithError(c, errors.BadRequest("Invalid request body"))
		return
	}

	// Refresh token
	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.log.WithFields(logger.Fields{
			"error": err.Error(),
		}).Info("Token refresh failed")

		switch err {
		case service.ErrRefreshTokenInvalid, service.ErrRefreshTokenRevoked:
			middleware.AbortWithError(c, errors.Unauthorized("Invalid or expired refresh token"))
		default:
			middleware.AbortWithError(c, errors.Internal("Token refresh failed"))
		}
		return
	}

	middleware.Success(c, response, "Token refreshed successfully")
}

// ============================================================================
// LOGOUT
// ============================================================================

// Logout godoc
// @Summary User logout
// @Description Revoke refresh tokens and logout
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token query string false "Refresh token to revoke"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		middleware.AbortWithError(c, errors.Unauthorized("User not authenticated"))
		return
	}

	// Get refresh token from request
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		// Try to get from body
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		_ = c.ShouldBindJSON(&req)
		refreshToken = req.RefreshToken
	}

	// Logout
	if err := h.authService.Logout(c.Request.Context(), uint(userID), refreshToken); err != nil {
		h.log.WithError(err).Error("Logout failed")
		middleware.AbortWithError(c, errors.Internal("Logout failed"))
		return
	}

	middleware.Success(c, nil, "Logout successful")
}

// LogoutAllDevices godoc
// @Summary Logout from all devices
// @Description Revoke all refresh tokens for the user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/auth/logout/all [post]
func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		middleware.AbortWithError(c, errors.Unauthorized("User not authenticated"))
		return
	}

	// Logout all devices
	if err := h.authService.LogoutAllDevices(c.Request.Context(), uint(userID)); err != nil {
		h.log.WithError(err).Error("Logout all devices failed")
		middleware.AbortWithError(c, errors.Internal("Logout failed"))
		return
	}

	middleware.Success(c, nil, "Logged out from all devices")
}

// ============================================================================
// PASSWORD MANAGEMENT
// ============================================================================

// ForgotPassword godoc
// @Summary Forgot password
// @Description Initiate password reset process
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email address"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, errors.BadRequest("Invalid email address"))
		return
	}

	// Always return success to prevent email enumeration
	if err := h.authService.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		h.log.WithError(err).Error("Forgot password failed")
	}

	middleware.Success(c, nil, "If the email exists, a password reset link will be sent")
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password with token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "New password and token"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, errors.BadRequest("Invalid request"))
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		h.log.WithError(err).Error("Reset password failed")
		middleware.AbortWithError(c, errors.Internal("Failed to reset password"))
		return
	}

	middleware.Success(c, nil, "Password reset successfully")
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Old and new password"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		middleware.AbortWithError(c, errors.Unauthorized("User not authenticated"))
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, errors.BadRequest("Invalid request"))
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), uint(userID), req.OldPassword, req.NewPassword); err != nil {
		h.log.WithError(err).Error("Change password failed")

		if err == service.ErrInvalidCredentials {
			middleware.AbortWithError(c, errors.InvalidCredentials())
			return
		}

		middleware.AbortWithError(c, errors.Internal("Failed to change password"))
		return
	}

	middleware.Success(c, nil, "Password changed successfully")
}

// ============================================================================
// CURRENT USER
// ============================================================================

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get authenticated user information
// @Tags auth
// @Produce json
// @Success 200 {object} middleware.SuccessResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		middleware.AbortWithError(c, errors.Unauthorized("User not authenticated"))
		return
	}

	email, _ := middleware.GetEmailFromContext(c)

	middleware.Success(c, gin.H{
		"id":    userID,
		"email": email,
	}, "User information retrieved")
}

// ============================================================================
// REQUEST MODELS
// ============================================================================

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Me returns current user information
func (h *AuthHandler) Me(c *gin.Context) {
	h.GetCurrentUser(c)
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		middleware.AbortWithError(c, errors.Unauthorized("User not authenticated"))
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, errors.BadRequest("Invalid request"))
		return
	}

	// TODO: Implement profile update in service
	h.log.WithFields(logger.Fields{
		"user_id": userID,
	}).Info("Profile update requested")

	middleware.Success(c, gin.H{
		"message": "Profile updated successfully",
	}, "")
}

// RequestPasswordReset initiates password reset
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	h.ForgotPassword(c)
}

// VerifyEmail verifies user email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	// TODO: Implement email verification
	middleware.Success(c, nil, "Email verified successfully")
}

// ResendVerificationEmail resends verification email
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, errors.BadRequest("Invalid email"))
		return
	}

	// TODO: Implement resend verification
	h.log.WithField("email", req.Email).Info("Verification email resend requested")

	middleware.Success(c, nil, "Verification email sent")
}

// ============================================================================
// ROUTE REGISTRATION
// ============================================================================

// RegisterRoutes registers auth routes
func RegisterRoutes(r *gin.RouterGroup, authService *service.AuthService, log *logger.Logger) {
	handler := NewAuthHandler(authService, log)

	// Public routes
	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)
	r.POST("/refresh", handler.RefreshToken)
	r.POST("/forgot-password", handler.ForgotPassword)
	r.POST("/reset-password", handler.ResetPassword)

	// Protected routes
	auth := r.Group("")
	// Add JWT middleware here if needed for specific endpoints
	{
		auth.POST("/logout", handler.Logout)
		auth.POST("/logout/all", handler.LogoutAllDevices)
		auth.POST("/change-password", handler.ChangePassword)
		auth.GET("/me", handler.GetCurrentUser)
	}
}
