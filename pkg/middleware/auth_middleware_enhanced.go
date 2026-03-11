package middleware

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth creates a JWT authentication middleware
func JWTAuth(tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("Authorization header required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid authorization format. Use: Bearer <token>"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate access token
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			switch err {
			case service.ErrExpiredToken:
				c.JSON(http.StatusUnauthorized, response.Unauthorized("Access token has expired"))
			case service.ErrInvalidToken:
				c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid token"))
			case service.ErrInvalidTokenType:
				c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid token type"))
			default:
				c.JSON(http.StatusUnauthorized, response.Unauthorized("Token validation failed"))
			}
			c.Abort()
			return
		}

		// Set user claims in context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalJWTAuth creates an optional JWT authentication middleware
// Continues even if authentication fails (useful for public endpoints that behave differently for logged-in users)
func OptionalJWTAuth(tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			// Don't abort, just continue without user context
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(roles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleValue, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, response.Forbidden("Authentication required"))
			c.Abort()
			return
		}

		userRole, ok := userRoleValue.(model.UserRole)
		if !ok {
			// Try to handle if role is stored as string
			if roleStr, ok := userRoleValue.(string); ok {
				// Convert string to model.UserRole for comparison
				for _, allowedRole := range roles {
					if roleStr == string(allowedRole) {
						c.Next()
						return
					}
				}
			}
			c.JSON(http.StatusForbidden, response.Forbidden("Invalid role format"))
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, response.Forbidden("Insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole creates a middleware that requires at least one of the specified roles
func RequireAnyRole(roles ...model.UserRole) gin.HandlerFunc {
	return RequireRole(roles...)
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(model.RoleUserAdmin)
}

// RequireSeller creates a middleware that requires seller role
func RequireSeller() gin.HandlerFunc {
	return RequireRole(model.RoleSeller)
}

// RequireCustomer creates a middleware that requires customer role
func RequireCustomer() gin.HandlerFunc {
	return RequireRole(model.RoleCustomer)
}

// RequireSellerOrAdmin creates a middleware that requires seller or admin role
func RequireSellerOrAdmin() gin.HandlerFunc {
	return RequireRole(model.RoleSeller, model.RoleUserAdmin)
}

// RequireOwnerOrAdmin creates a middleware that checks if user is owner or admin
// Expects resource owner ID to be set in context or as URL parameter
func RequireOwnerOrAdmin(getOwnerID func(*gin.Context) uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, userExists := c.Get("user_id")
		userRoleValue, roleExists := c.Get("user_role")

		if !userExists || !roleExists {
			c.JSON(http.StatusForbidden, response.Forbidden("Authentication required"))
			c.Abort()
			return
		}

		userID := userIDValue.(uint)
		userRole := userRoleValue.(model.UserRole)

		// Admins can access everything
		if userRole == model.RoleUserAdmin {
			c.Next()
			return
		}

		// Check if user is the owner
		ownerID := getOwnerID(c)
		if userID != ownerID {
			c.JSON(http.StatusForbidden, response.Forbidden("You can only access your own resources"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// Permission defines a permission type
type Permission string

const (
	PermissionCreate Permission = "create"
	PermissionRead   Permission = "read"
	PermissionUpdate Permission = "update"
	PermissionDelete Permission = "delete"
)

// RequirePermission creates a middleware that requires specific permissions
// This is for more fine-grained access control beyond roles
func RequirePermission(resource string, permissions ...Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleValue, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, response.Forbidden("Authentication required"))
			c.Abort()
			return
		}

		userRole := userRoleValue.(model.UserRole)

		// Check permissions based on role and resource
		// In production, this would check against a permission matrix/database
		hasPermission := checkPermission(userRole, resource, permissions...)

		if !hasPermission {
			c.JSON(http.StatusForbidden, response.Forbidden("Insufficient permissions for this action"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkPermission checks if a role has the required permissions for a resource
func checkPermission(role model.UserRole, resource string, permissions ...Permission) bool {
	// Define permission matrix
	permissionMatrix := map[model.UserRole]map[string][]Permission{
		model.RoleUserAdmin: {
			"products": {PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete},
			"orders":   {PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete},
			"users":    {PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete},
		},
		model.RoleSeller: {
			"products": {PermissionCreate, PermissionRead, PermissionUpdate},
			"orders":   {PermissionRead, PermissionUpdate},
		},
		model.RoleCustomer: {
			"products": {PermissionRead},
			"orders":   {PermissionCreate, PermissionRead},
			"cart":     {PermissionCreate, PermissionRead, PermissionUpdate, PermissionDelete},
		},
	}

	rolePermissions, exists := permissionMatrix[role]
	if !exists {
		return false
	}

	resourcePermissions, exists := rolePermissions[resource]
	if !exists {
		return false
	}

	// Check if all required permissions are granted
	for _, requiredPerm := range permissions {
		found := false
		for _, allowedPerm := range resourcePermissions {
			if requiredPerm == allowedPerm {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// RateLimitByUserID creates a rate limiting key based on user ID
func RateLimitByUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if exists {
		return "user:" + string(rune(userID.(uint)))
	}
	return "ip:" + c.ClientIP()
}

// ExtractToken extracts token from request
func ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// IsAuthenticated checks if the request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// GetCurrentUserRole gets the current user's role from context
func GetCurrentUserRole(c *gin.Context) (model.UserRole, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(model.UserRole), true
}

// IsAdmin checks if current user is admin
func IsAdmin(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == model.RoleUserAdmin
}

// IsSeller checks if current user is seller
func IsSeller(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == model.RoleSeller
}

// IsCustomer checks if current user is customer
func IsCustomer(c *gin.Context) bool {
	role, exists := GetCurrentUserRole(c)
	return exists && role == model.RoleCustomer
}
