// Package jwt provides JWT token generation and validation
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ============================================================================
// ERRORS
// ============================================================================

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidClaims = errors.New("invalid token claims")
)

// ============================================================================
// CLAIMS
// ============================================================================

// Claims represents JWT claims with custom fields
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ============================================================================
// SERVICE
// ============================================================================

// Service handles JWT operations
type Service struct {
	secret          []byte
	accessExpiry    time.Duration
	refreshExpiry   time.Duration
	issuer          string
}

// Config holds JWT configuration
type Config struct {
	// Secret is the signing key (minimum 32 bytes recommended)
	Secret string

	// AccessExpiry is the access token lifetime (default: 15 minutes)
	AccessExpiry time.Duration

	// RefreshExpiry is the refresh token lifetime (default: 7 days)
	RefreshExpiry time.Duration

	// Issuer is the token issuer (default: "ecommerce-api")
	Issuer string
}

// NewService creates a new JWT service
func NewService(config Config) *Service {
	// Set defaults
	if config.AccessExpiry == 0 {
		config.AccessExpiry = 15 * time.Minute
	}
	if config.RefreshExpiry == 0 {
		config.RefreshExpiry = 7 * 24 * time.Hour
	}
	if config.Issuer == "" {
		config.Issuer = "ecommerce-api"
	}

	return &Service{
		secret:          []byte(config.Secret),
		accessExpiry:    config.AccessExpiry,
		refreshExpiry:   config.RefreshExpiry,
		issuer:          config.Issuer,
	}
}

// ============================================================================
// TOKEN GENERATION
// ============================================================================

// GenerateAccessToken generates a new access token
func (s *Service) GenerateAccessToken(userID int64, email, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   string(rune(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken generates a new refresh token
func (s *Service) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   string(rune(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateTokenPair generates both access and refresh tokens
func (s *Service) GenerateTokenPair(userID int64, email, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = s.GenerateAccessToken(userID, email, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ============================================================================
// TOKEN VALIDATION
// ============================================================================

// ValidateToken validates a token and returns claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verify issuer
	if claims.Issuer != s.issuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verify issuer
	if claims.Issuer != s.issuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ============================================================================
// TOKEN INFO
// ============================================================================

// GetTokenExpiry returns the access token expiry duration
func (s *Service) GetTokenExpiry() time.Duration {
	return s.accessExpiry
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (s *Service) GetRefreshTokenExpiry() time.Duration {
	return s.refreshExpiry
}

// ============================================================================
// TOKEN EXTRACTION
// ============================================================================

// ExtractFromAuthHeader extracts token from Authorization header
// Expected format: "Bearer <token>"
func ExtractFromAuthHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// ============================================================================
// TOKEN BLACKLIST (Optional - for logout)
// ============================================================================

// TokenBlacklist interface for token blacklist storage
type TokenBlacklist interface {
	// Add adds a token to the blacklist with expiry
	Add(token string, expiry time.Time) error
	// IsBlacklisted checks if a token is blacklisted
	IsBlacklisted(token string) (bool, error)
	// Remove removes a token from the blacklist
	Remove(token string) error
}

// InMemoryBlacklist is an in-memory token blacklist
type InMemoryBlacklist struct {
	tokens map[string]time.Time
}

// NewInMemoryBlacklist creates a new in-memory blacklist
func NewInMemoryBlacklist() *InMemoryBlacklist {
	return &InMemoryBlacklist{
		tokens: make(map[string]time.Time),
	}
}

// Add adds a token to the blacklist
func (b *InMemoryBlacklist) Add(token string, expiry time.Time) error {
	b.tokens[token] = expiry
	return nil
}

// IsBlacklisted checks if a token is blacklisted
func (b *InMemoryBlacklist) IsBlacklisted(token string) (bool, error) {
	expiry, exists := b.tokens[token]
	if !exists {
		return false, nil
	}

	// Remove expired entries
	if time.Now().After(expiry) {
		delete(b.tokens, token)
		return false, nil
	}

	return true, nil
}

// Remove removes a token from the blacklist
func (b *InMemoryBlacklist) Remove(token string) error {
	delete(b.tokens, token)
	return nil
}

// Cleanup removes expired tokens from the blacklist
func (b *InMemoryBlacklist) Cleanup() {
	now := time.Now()
	for token, expiry := range b.tokens {
		if now.After(expiry) {
			delete(b.tokens, token)
		}
	}
}

// StartCleanup starts periodic cleanup of expired tokens
func (b *InMemoryBlacklist) StartCleanup(interval time.Duration, stopChan chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				b.Cleanup()
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}
