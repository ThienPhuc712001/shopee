package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// TEST SETUP
// ============================================================================

func setupJWTService() *Service {
	return NewService(Config{
		Secret:       "test-secret-key-min-32-characters-long",
		AccessExpiry: 15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
		Issuer:       "test-api",
	})
}

// ============================================================================
// TOKEN GENERATION TESTS
// ============================================================================

func TestGenerateAccessToken_Success(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	token, err := service.GenerateAccessToken(1, "test@example.com", "customer")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken_Success(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	token, err := service.GenerateRefreshToken(1)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateTokenPair_Success(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	accessToken, refreshToken, err := service.GenerateTokenPair(1, "test@example.com", "customer")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

// ============================================================================
// TOKEN VALIDATION TESTS
// ============================================================================

func TestValidateAccessToken_Success(t *testing.T) {
	// Arrange
	service := setupJWTService()
	token, _ := service.GenerateAccessToken(1, "test@example.com", "customer")

	// Act
	claims, err := service.ValidateAccessToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "customer", claims.Role)
}

func TestValidateRefreshToken_Success(t *testing.T) {
	// Arrange
	service := setupJWTService()
	token, _ := service.GenerateRefreshToken(1)

	// Act
	claims, err := service.ValidateRefreshToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, int64(1), claims.UserID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	claims, err := service.ValidateToken("invalid_token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	claims, err := service.ValidateToken("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	// Arrange
	service1 := NewService(Config{
		Secret: "secret-one-min-32-characters-long",
	})
	service2 := NewService(Config{
		Secret: "secret-two-min-32-characters-long",
	})

	token, _ := service1.GenerateAccessToken(1, "test@example.com", "customer")

	// Act
	claims, err := service2.ValidateAccessToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// ============================================================================
// TOKEN EXPIRATION TESTS
// ============================================================================

func TestValidateToken_Expired(t *testing.T) {
	// Arrange
	service := NewService(Config{
		Secret:       "test-secret-key-min-32-characters-long",
		AccessExpiry: -1 * time.Minute, // Already expired
	})

	token, _ := service.GenerateAccessToken(1, "test@example.com", "customer")

	// Act
	claims, err := service.ValidateAccessToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrExpiredToken, err)
}

func TestGetTokenExpiry(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	expiry := service.GetTokenExpiry()

	// Assert
	assert.Equal(t, 15*time.Minute, expiry)
}

func TestGetRefreshTokenExpiry(t *testing.T) {
	// Arrange
	service := setupJWTService()

	// Act
	expiry := service.GetRefreshTokenExpiry()

	// Assert
	assert.Equal(t, 7*24*time.Hour, expiry)
}

// ============================================================================
// TOKEN EXTRACTION TESTS
// ============================================================================

func TestExtractFromAuthHeader_Success(t *testing.T) {
	// Arrange
	authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

	// Act
	token := ExtractFromAuthHeader(authHeader)

	// Assert
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", token)
}

func TestExtractFromAuthHeader_NoBearer(t *testing.T) {
	// Arrange
	authHeader := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

	// Act
	token := ExtractFromAuthHeader(authHeader)

	// Assert
	assert.Empty(t, token)
}

func TestExtractFromAuthHeader_Empty(t *testing.T) {
	// Arrange
	authHeader := ""

	// Act
	token := ExtractFromAuthHeader(authHeader)

	// Assert
	assert.Empty(t, token)
}

func TestExtractFromAuthHeader_Malformed(t *testing.T) {
	// Arrange
	authHeader := "Bearer"

	// Act
	token := ExtractFromAuthHeader(authHeader)

	// Assert
	assert.Empty(t, token)
}

// ============================================================================
// TOKEN BLACKLIST TESTS
// ============================================================================

func TestInMemoryBlacklist_Add(t *testing.T) {
	// Arrange
	blacklist := NewInMemoryBlacklist()
	token := "test_token"
	expiry := time.Now().Add(15 * time.Minute)

	// Act
	err := blacklist.Add(token, expiry)

	// Assert
	assert.NoError(t, err)
}

func TestInMemoryBlacklist_IsBlacklisted(t *testing.T) {
	// Arrange
	blacklist := NewInMemoryBlacklist()
	token := "test_token"
	expiry := time.Now().Add(15 * time.Minute)
	blacklist.Add(token, expiry)

	// Act
	isBlacklisted, err := blacklist.IsBlacklisted(token)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isBlacklisted)
}

func TestInMemoryBlacklist_IsBlacklisted_NotFound(t *testing.T) {
	// Arrange
	blacklist := NewInMemoryBlacklist()

	// Act
	isBlacklisted, err := blacklist.IsBlacklisted("nonexistent_token")

	// Assert
	assert.NoError(t, err)
	assert.False(t, isBlacklisted)
}

func TestInMemoryBlacklist_Remove(t *testing.T) {
	// Arrange
	blacklist := NewInMemoryBlacklist()
	token := "test_token"
	expiry := time.Now().Add(15 * time.Minute)
	blacklist.Add(token, expiry)

	// Act
	err := blacklist.Remove(token)
	isBlacklisted, _ := blacklist.IsBlacklisted(token)

	// Assert
	assert.NoError(t, err)
	assert.False(t, isBlacklisted)
}

func TestInMemoryBlacklist_Cleanup(t *testing.T) {
	// Arrange
	blacklist := NewInMemoryBlacklist()
	expiredToken := "expired_token"
	validToken := "valid_token"

	blacklist.Add(expiredToken, time.Now().Add(-1*time.Minute))
	blacklist.Add(validToken, time.Now().Add(15*time.Minute))

	// Act
	blacklist.Cleanup()

	// Assert
	isBlacklisted, _ := blacklist.IsBlacklisted(expiredToken)
	assert.False(t, isBlacklisted) // Should be removed

	isBlacklisted, _ = blacklist.IsBlacklisted(validToken)
	assert.True(t, isBlacklisted) // Should still be there
}

// ============================================================================
// CLAIMS TESTS
// ============================================================================

func TestClaims_Struct(t *testing.T) {
	// Arrange
	claims := Claims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "customer",
	}

	// Assert
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "customer", claims.Role)
}

// ============================================================================
// ERROR TESTS
// ============================================================================

func TestErrors(t *testing.T) {
	// Test error messages
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "token has expired", ErrExpiredToken.Error())
	assert.Equal(t, "invalid token claims", ErrInvalidClaims.Error())
}

// ============================================================================
// CONFIGURATION TESTS
// ============================================================================

func TestNewService_Defaults(t *testing.T) {
	// Arrange
	config := Config{
		Secret: "test-secret-key-min-32-characters-long",
	}

	// Act
	service := NewService(config)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, 15*time.Minute, service.accessExpiry)    // Default
	assert.Equal(t, 7*24*time.Hour, service.refreshExpiry)   // Default
	assert.Equal(t, "ecommerce-api", service.issuer)         // Default
}

func TestNewService_CustomConfig(t *testing.T) {
	// Arrange
	config := Config{
		Secret:        "test-secret-key-min-32-characters-long",
		AccessExpiry:  30 * time.Minute,
		RefreshExpiry: 14 * 24 * time.Hour,
		Issuer:        "custom-api",
	}

	// Act
	service := NewService(config)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, 30*time.Minute, service.accessExpiry)
	assert.Equal(t, 14*24*time.Hour, service.refreshExpiry)
	assert.Equal(t, "custom-api", service.issuer)
}
