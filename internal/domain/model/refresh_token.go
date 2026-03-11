// Package model contains domain models
package model

import (
	"time"
)

// RefreshToken represents a JWT refresh token stored in the database
type RefreshToken struct {
	// ID is the primary key
	ID int64 `json:"id" gorm:"primaryKey"`

	// UserID is the user who owns this token
	UserID int64 `json:"user_id" gorm:"not null;index"`

	// Token is the JWT refresh token string
	Token string `json:"-" gorm:"type:nvarchar(500);not null;uniqueIndex"`

	// ExpiresAt is when the token expires
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`

	// Revoked indicates if the token has been revoked
	Revoked bool `json:"revoked" gorm:"default:false;index"`

	// RevokedAt is when the token was revoked
	RevokedAt *time.Time `json:"revoked_at,omitempty"`

	// CreatedAt is when the token was created
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked checks if the token is revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.Revoked
}

// IsValid checks if the token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}

// Revoke marks the token as revoked
func (rt *RefreshToken) Revoke() {
	rt.Revoked = true
	now := time.Now()
	rt.RevokedAt = &now
}
