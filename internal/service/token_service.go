package service

import (
	"ecommerce/internal/domain/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token errors
var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("token has expired")
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrMissingClaim      = errors.New("missing required claim")
	ErrTokenNotActiveYet = errors.New("token not active yet")
)

// TokenType defines the type of token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
}

// TokenClaims represents JWT claims with custom fields
type TokenClaims struct {
	UserID    uint           `json:"user_id"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	TokenType TokenType      `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenService handles JWT token operations
type TokenService interface {
	// GenerateTokenPair generates both access and refresh tokens
	GenerateTokenPair(user *model.User) (*TokenPair, error)

	// GenerateAccessToken generates an access token
	GenerateAccessToken(user *model.User) (string, time.Time, error)

	// GenerateRefreshToken generates a refresh token
	GenerateRefreshToken(user *model.User) (string, time.Time, error)

	// ValidateToken validates a token and returns claims
	ValidateToken(tokenString string, expectedType TokenType) (*TokenClaims, error)

	// ValidateAccessToken validates an access token
	ValidateAccessToken(tokenString string) (*TokenClaims, error)

	// ValidateRefreshToken validates a refresh token
	ValidateRefreshToken(tokenString string) (*TokenClaims, error)

	// RefreshAccessToken generates a new access token using refresh token
	RefreshAccessToken(refreshToken string) (*TokenPair, error)

	// RevokeToken invalidates a refresh token
	RevokeToken(tokenString string) error

	// GetTokenExpiry returns the expiry time of a token
	GetTokenExpiry(tokenString string) (time.Time, error)

	// GetRefreshExpiry returns the refresh token expiry duration
	GetRefreshExpiry() time.Duration
}

type tokenService struct {
	accessSecret     []byte
	refreshSecret    []byte
	accessExpiry     time.Duration
	refreshExpiry    time.Duration
	issuer           string
}

// TokenServiceConfig holds configuration for token service
type TokenServiceConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

// NewTokenService creates a new token service
func NewTokenService(cfg TokenServiceConfig) TokenService {
	return &tokenService{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
		issuer:        cfg.Issuer,
	}
}

// DefaultTokenService creates a token service with default settings
func DefaultTokenService(accessSecret, refreshSecret string) TokenService {
	return &tokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
		issuer:        "ecommerce-api",
	}
}

func (s *tokenService) GenerateTokenPair(user *model.User) (*TokenPair, error) {
	accessToken, accessExpiry, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
		TokenType:    "Bearer",
		Expiry:       accessExpiry,
	}, nil
}

func (s *tokenService) GenerateAccessToken(user *model.User) (string, time.Time, error) {
	now := time.Now()
	expiry := now.Add(s.accessExpiry)

	claims := &TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.accessSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiry, nil
}

func (s *tokenService) GenerateRefreshToken(user *model.User) (string, time.Time, error) {
	now := time.Now()
	expiry := now.Add(s.refreshExpiry)

	claims := &TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   string(rune(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.refreshSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiry, nil
}

func (s *tokenService) ValidateToken(tokenString string, expectedType TokenType) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		claims, ok := token.Claims.(*TokenClaims)
		if !ok {
			return nil, ErrInvalidToken
		}

		// Use appropriate secret based on token type
		if claims.TokenType == TokenTypeAccess {
			return s.accessSecret, nil
		}
		return s.refreshSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotActiveYet
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Verify token type
	if claims.TokenType != expectedType {
		return nil, ErrInvalidTokenType
	}

	// Verify issuer
	if claims.Issuer != s.issuer {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *tokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return s.ValidateToken(tokenString, TokenTypeAccess)
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return s.ValidateToken(tokenString, TokenTypeRefresh)
}

func (s *tokenService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Create a mock user object from claims
	user := &model.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}

	// Generate new token pair
	return s.GenerateTokenPair(user)
}

func (s *tokenService) RevokeToken(tokenString string) error {
	// In production, you would:
	// 1. Add token to a blacklist (Redis/database)
	// 2. Clear refresh token from user record
	// For now, we just validate the token
	_, err := s.ValidateToken(tokenString, TokenTypeRefresh)
	return err
}

func (s *tokenService) GetTokenExpiry(tokenString string) (time.Time, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &TokenClaims{})
	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return time.Time{}, ErrInvalidToken
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, ErrMissingClaim
	}

	return claims.ExpiresAt.Time, nil
}

// GetRefreshExpiry returns the refresh token expiry duration
func (s *tokenService) GetRefreshExpiry() time.Duration {
	return s.refreshExpiry
}

// TokenBlacklist defines interface for token blacklist
type TokenBlacklist interface {
	// Add adds a token to the blacklist
	Add(token string, expiry time.Time) error

	// Contains checks if a token is blacklisted
	Contains(token string) (bool, error)

	// Remove removes a token from the blacklist
	Remove(token string) error

	// Cleanup removes expired tokens
	Cleanup() error
}
