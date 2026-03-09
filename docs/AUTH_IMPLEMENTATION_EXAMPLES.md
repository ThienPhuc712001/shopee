# Authentication Implementation Examples

## Complete Working Examples

### 1. Password Hashing Example

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "mypassword123"
    
    // Hash password with default cost (10)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Hashed Password:", string(hashedPassword))
    // Output: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
    
    // Verify password
    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err == nil {
        fmt.Println("Password matches!")
    } else {
        fmt.Println("Password does not match!")
    }
    
    // Hash with custom cost (higher = more secure but slower)
    hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        panic(err)
    }
    fmt.Println("Hashed with cost 12:", string(hashedPassword))
}
```

### 2. JWT Token Generation Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
    UserID    uint   `json:"user_id"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    TokenType string `json:"token_type"`
    jwt.RegisteredClaims
}

func main() {
    secretKey := []byte("your-secret-key-min-32-characters-long")
    
    // Create claims
    claims := &CustomClaims{
        UserID:    1,
        Email:     "user@example.com",
        Role:      "customer",
        TokenType: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "ecommerce-api",
            Subject:   "1",
        },
    }
    
    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign token
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Access Token:", tokenString)
    // Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
    
    // Parse and validate token
    parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    
    if err != nil {
        panic(err)
    }
    
    if claims, ok := parsedToken.Claims.(*CustomClaims); ok && parsedToken.Valid {
        fmt.Printf("User ID: %d, Email: %s, Role: %s\n", claims.UserID, claims.Email, claims.Role)
    }
}
```

### 3. Complete Register Handler Example

```go
package handler

import (
    "ecommerce/internal/domain/model"
    "ecommerce/internal/domain/repository"
    "ecommerce/internal/service"
    "ecommerce/pkg/response"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    Phone     string `json:"phone" binding:"required"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

func Register(c *gin.Context) {
    var req RegisterRequest
    
    // 1. Validate input
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
        return
    }
    
    // 2. Validate password strength
    if err := validatePassword(req.Password); err != nil {
        c.JSON(http.StatusBadRequest, response.Error(err.Error()))
        return
    }
    
    // 3. Check if email already exists
    db := getDatabase() // Your database connection
    var existingUser model.User
    result := db.Where("email = ?", req.Email).First(&existingUser)
    if result.Error == nil {
        c.JSON(http.StatusConflict, response.Error("Email already registered"))
        return
    }
    
    // 4. Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.InternalError("Failed to process password"))
        return
    }
    
    // 5. Create user
    user := model.User{
        Email:     req.Email,
        Password:  string(hashedPassword),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Phone:     req.Phone,
        Role:      model.RoleCustomer,
        Status:    model.StatusActive,
    }
    
    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, response.InternalError("Failed to create user"))
        return
    }
    
    // 6. Generate JWT token
    token, err := generateToken(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.InternalError("Failed to generate token"))
        return
    }
    
    // 7. Return success response
    c.JSON(http.StatusCreated, response.Success(gin.H{
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
            "name":  user.GetFullName(),
            "role":  user.Role,
        },
        "token": token,
    }, "Registration successful"))
}

func validatePassword(password string) error {
    if len(password) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
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
        case char == '!' || char == '@' || char == '#' || char == '$':
            hasSpecial = true
        }
    }
    
    if !hasUpper {
        return fmt.Errorf("password must contain at least one uppercase letter")
    }
    if !hasLower {
        return fmt.Errorf("password must contain at least one lowercase letter")
    }
    if !hasNumber {
        return fmt.Errorf("password must contain at least one number")
    }
    if !hasSpecial {
        return fmt.Errorf("password must contain at least one special character")
    }
    
    return nil
}

func generateToken(user *model.User) (string, error) {
    secret := "your-secret-key-min-32-characters-long"
    
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### 4. Complete Login Handler Example

```go
package handler

import (
    "ecommerce/internal/domain/model"
    "ecommerce/pkg/response"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var req LoginRequest
    
    // 1. Validate input
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
        return
    }
    
    db := getDatabase()
    
    // 2. Find user by email
    var user model.User
    result := db.Where("email = ?", req.Email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid email or password"))
        return
    }
    
    // 3. Check if account is locked
    if user.Status == model.StatusLocked {
        c.JSON(http.StatusLocked, response.Error("Account is locked. Please contact support."))
        return
    }
    
    // 4. Check if account is active
    if user.Status != model.StatusActive {
        c.JSON(http.StatusForbidden, response.Error("Account is inactive"))
        return
    }
    
    // 5. Verify password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        // Increment failed login attempts
        incrementFailedLogin(user.ID)
        c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid email or password"))
        return
    }
    
    // 6. Reset failed login attempts
    resetFailedLogin(user.ID)
    
    // 7. Update last login
    db.Model(&user).Update("last_login", time.Now())
    
    // 8. Generate JWT tokens
    accessToken, err := generateAccessToken(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.InternalError("Failed to generate token"))
        return
    }
    
    refreshToken, err := generateRefreshToken(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.InternalError("Failed to generate token"))
        return
    }
    
    // 9. Save refresh token to database
    saveRefreshToken(user.ID, refreshToken)
    
    // 10. Return success response
    c.JSON(http.StatusOK, response.Success(gin.H{
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
            "name":  user.GetFullName(),
            "role":  user.Role,
            "avatar": user.Avatar,
        },
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "expires_in":    900, // 15 minutes in seconds
        "token_type":    "Bearer",
    }, "Login successful"))
}

func generateAccessToken(user *model.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id":    user.ID,
        "email":      user.Email,
        "role":       user.Role,
        "token_type": "access",
        "exp":        time.Now().Add(15 * time.Minute).Unix(),
        "iat":        time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("access-token-secret-key"))
}

func generateRefreshToken(user *model.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id":    user.ID,
        "email":      user.Email,
        "token_type": "refresh",
        "exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
        "iat":        time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("refresh-token-secret-key"))
}

func incrementFailedLogin(userID uint) {
    db := getDatabase()
    db.Model(&model.User{}).Where("id = ?", userID).
        UpdateColumn("failed_login_attempts", gorm.Expr("failed_login_attempts + 1"))
}

func resetFailedLogin(userID uint) {
    db := getDatabase()
    db.Model(&model.User{}).Where("id = ?", userID).
        Updates(map[string]interface{}{
            "failed_login_attempts": 0,
            "locked_until":          nil,
        })
}

func saveRefreshToken(userID uint, token string) {
    db := getDatabase()
    expiry := time.Now().Add(7 * 24 * time.Hour)
    db.Model(&model.User{}).Where("id = ?", userID).
        Updates(map[string]interface{}{
            "refresh_token":       token,
            "refresh_token_expiry": expiry,
        })
}
```

### 5. JWT Middleware Example

```go
package middleware

import (
    "ecommerce/pkg/response"
    "net/http"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
    secretKey := []byte("your-secret-key-min-32-characters-long")
    
    return func(c *gin.Context) {
        // 1. Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, response.Unauthorized("Authorization header required"))
            c.Abort()
            return
        }
        
        // 2. Check Bearer format
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid authorization format"))
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // 3. Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return secretKey, nil
        })
        
        // 4. Handle parsing errors
        if err != nil {
            if errors.Is(err, jwt.ErrTokenExpired) {
                c.JSON(http.StatusUnauthorized, response.Unauthorized("Token has expired"))
            } else {
                c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid token"))
            }
            c.Abort()
            return
        }
        
        // 5. Validate token
        if !token.Valid {
            c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid token"))
            c.Abort()
            return
        }
        
        // 6. Extract claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, response.Unauthorized("Invalid token claims"))
            c.Abort()
            return
        }
        
        // 7. Check expiration
        if exp, ok := claims["exp"].(float64); ok {
            if time.Now().Unix() > int64(exp) {
                c.JSON(http.StatusUnauthorized, response.Unauthorized("Token expired"))
                c.Abort()
                return
            }
        }
        
        // 8. Set user info in context
        c.Set("user_id", uint(claims["user_id"].(float64)))
        c.Set("user_email", claims["email"].(string))
        c.Set("user_role", claims["role"].(string))
        
        c.Next()
    }
}

// Usage in routes:
func SetupRoutes(r *gin.Engine) {
    // Public routes
    r.POST("/api/auth/register", Register)
    r.POST("/api/auth/login", Login)
    
    // Protected routes (require authentication)
    protected := r.Group("/api")
    protected.Use(JWTMiddleware())
    {
        protected.GET("/users/me", GetCurrentUser)
        protected.PUT("/users/profile", UpdateProfile)
        protected.POST("/orders", CreateOrder)
    }
}
```

### 6. Role-Based Access Control Example

```go
package middleware

import (
    "ecommerce/pkg/response"
    "net/http"
)

func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(http.StatusForbidden, response.Forbidden("Role not found"))
            c.Abort()
            return
        }
        
        roleStr := userRole.(string)
        
        // Check if user role is in allowed roles
        for _, role := range roles {
            if roleStr == role {
                c.Next()
                return
            }
        }
        
        c.JSON(http.StatusForbidden, response.Forbidden("Insufficient permissions"))
        c.Abort()
    }
}

// Usage in routes:
func SetupRBACRoutes(r *gin.Engine) {
    api := r.Group("/api")
    api.Use(JWTMiddleware())
    {
        // Only sellers can create products
        api.POST("/products", RequireRole("seller", "admin"), CreateProduct)
        
        // Only admins can delete products
        api.DELETE("/products/:id", RequireRole("admin"), DeleteProduct)
        
        // Only admins can manage users
        admin := api.Group("/admin")
        admin.Use(RequireRole("admin"))
        {
            admin.GET("/users", GetAllUsers)
            admin.PUT("/users/:id/status", UpdateUserStatus)
        }
    }
}
```

### 7. Standard API Response Format

```go
// All responses follow this format:

// Success Response:
{
    "success": true,
    "data": {
        "user": {
            "id": 1,
            "email": "user@example.com"
        },
        "token": "eyJhbGc..."
    },
    "message": "Login successful"
}

// Error Response:
{
    "success": false,
    "error": "Invalid email or password"
}

// Paginated Response:
{
    "success": true,
    "data": {
        "products": [...]
    },
    "meta": {
        "current_page": 1,
        "per_page": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

### 8. HTTP Status Codes Reference

| Code | Meaning | When to Use |
|------|---------|-------------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful resource creation (POST) |
| 204 | No Content | Successful deletion |
| 400 | Bad Request | Invalid input, validation errors |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Valid token but insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Email already exists |
| 422 | Unprocessable Entity | Validation errors |
| 423 | Locked | Account locked |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
