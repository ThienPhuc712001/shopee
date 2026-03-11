package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ecommerce/internal/errors"
	"ecommerce/pkg/jwt"
	"ecommerce/pkg/logger"
)

// JWTAuthConfig holds JWT authentication configuration
type JWTAuthConfig struct {
	// JWTService is the JWT service for token validation
	JWTService *jwt.Service

	// Log is the logger for authentication events
	Log *logger.Logger

	// Blacklist is optional token blacklist for logout
	Blacklist jwt.TokenBlacklist

	// SkipPaths are paths that don't require authentication
	SkipPaths map[string]bool
}

// JWTAuth creates a JWT authentication middleware
func JWTAuth(config JWTAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		if config.SkipPaths != nil && config.SkipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if config.Log != nil {
				config.Log.WithFields(logger.Fields{
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
					"client_ip": c.ClientIP(),
				}).Info("Missing authorization header")
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.
				Unauthorized("Missing authorization header").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		// Extract Bearer token
		token := jwt.ExtractFromAuthHeader(authHeader)
		if token == "" {
			if config.Log != nil {
				config.Log.WithFields(logger.Fields{
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
					"client_ip": c.ClientIP(),
				}).Info("Invalid authorization format")
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.
				Unauthorized("Invalid authorization format. Use: Bearer <token>").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		// Check token blacklist
		if config.Blacklist != nil {
			blacklisted, err := config.Blacklist.IsBlacklisted(token)
			if err != nil {
				if config.Log != nil {
					config.Log.WithError(err).Error("Failed to check token blacklist")
				}
			}
			if blacklisted {
				if config.Log != nil {
					config.Log.WithFields(logger.Fields{
						"path":      c.Request.URL.Path,
						"method":    c.Request.Method,
						"client_ip": c.ClientIP(),
					}).Info("Token is blacklisted")
				}

				c.AbortWithStatusJSON(http.StatusUnauthorized, errors.
					Unauthorized("Token has been revoked").
					WithPath(c.Request.URL.Path).
					ToResponse())
				return
			}
		}

		// Validate token
		claims, err := config.JWTService.ValidateAccessToken(token)
		if err != nil {
			if config.Log != nil {
				config.Log.WithFields(logger.Fields{
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
					"client_ip": c.ClientIP(),
					"error":     err.Error(),
				}).Info("Invalid token")
			}

			if err == jwt.ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errors.
					TokenExpired().
					WithPath(c.Request.URL.Path).
					ToResponse())
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.
				Unauthorized("Invalid or expired token").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// ============================================================================
// ROLE-BASED AUTHORIZATION
// ============================================================================

// RequireRole creates a middleware that requires specific roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, errors.
				Forbidden("Role not found").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, errors.
				Forbidden("Invalid role format").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, role := range roles {
			if roleStr == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, errors.
				Forbidden("Insufficient permissions").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}

		c.Next()
	}
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin", "super_admin")
}

// RequireSeller creates a middleware that requires seller role
func RequireSeller() gin.HandlerFunc {
	return RequireRole("seller", "admin", "super_admin")
}

// ============================================================================
// OPTIONAL JWT AUTH
// ============================================================================

// OptionalJWTAuth creates a middleware that optionally authenticates if token is present
// Useful for endpoints that work for both authenticated and unauthenticated users
func OptionalJWTAuth(config JWTAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token := jwt.ExtractFromAuthHeader(authHeader)
		if token == "" {
			c.Next()
			return
		}

		claims, err := config.JWTService.ValidateAccessToken(token)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)
			c.Set("claims", claims)
		}

		c.Next()
	}
}

// ============================================================================
// TOKEN EXTRACTION HELPERS
// ============================================================================

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	if !ok {
		// Try to convert from float64 (JSON number)
		if floatID, ok := userID.(float64); ok {
			return int64(floatID), true
		}
		return 0, false
	}

	return id, true
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetClaimsFromContext extracts claims from context
func GetClaimsFromContext(c *gin.Context) (*jwt.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	claimsObj, ok := claims.(*jwt.Claims)
	return claimsObj, ok
}

// RequireAuth is a simple middleware that just checks if user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.
				Unauthorized("Authentication required").
				WithPath(c.Request.URL.Path).
				ToResponse())
			return
		}
		c.Next()
	}
}

// ============================================================================
// TOKEN REFRESH MIDDLEWARE
// ============================================================================

// RefreshTokenExtractor extracts refresh token from request
func RefreshTokenExtractor(c *gin.Context) string {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token := jwt.ExtractFromAuthHeader(authHeader)
		if token != "" {
			return token
		}
	}

	// Try form value
	refreshToken := c.PostForm("refresh_token")
	if refreshToken != "" {
		return refreshToken
	}

	// Try JSON body (requires manual parsing or binding)
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		return req.RefreshToken
	}

	return ""
}

// ExtractBearerToken extracts Bearer token from string
func ExtractBearerToken(authHeader string) string {
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}
