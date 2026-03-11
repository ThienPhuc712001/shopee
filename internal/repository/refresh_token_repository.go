// Package repository provides data access layer
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"ecommerce/internal/domain/model"
)

// RefreshTokenRepository handles refresh token data operations
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetByToken retrieves a refresh token by its token string
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.RefreshToken, error) {
	var tokens []*model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// GetValidByToken retrieves a valid (not expired, not revoked) refresh token
func (r *RefreshTokenRepository) GetValidByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// Revoke revokes a refresh token
func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"revoked":     true,
			"revoked_at":  time.Now(),
		})
	return result.Error
}

// RevokeByUserID revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID int64) error {
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":     true,
			"revoked_at":  time.Now(),
		})
	return result.Error
}

// DeleteExpired deletes expired tokens (for cleanup)
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	result := r.db.WithContext(ctx).
		Where("expires_at < ? OR (revoked = ? AND revoked_at < ?)", 
			time.Now(), true, cutoffTime).
		Delete(&model.RefreshToken{})
	return result.RowsAffected, result.Error
}

// DeleteByToken deletes a specific refresh token
func (r *RefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&model.RefreshToken{}).Error
}

// CountByUserID counts refresh tokens for a user
func (r *RefreshTokenRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// CleanupOldTokens periodically cleans up old revoked tokens
func (r *RefreshTokenRepository) CleanupOldTokens(ctx context.Context, retentionDays int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	result := r.db.WithContext(ctx).
		Where("(revoked = ? AND revoked_at < ?) OR expires_at < ?", 
			true, cutoffTime, time.Now()).
		Delete(&model.RefreshToken{})
	return result.RowsAffected, result.Error
}

// RevokeAllExceptCurrent revokes all tokens for a user except the current one
func (r *RefreshTokenRepository) RevokeAllExceptCurrent(ctx context.Context, userID int64, currentToken string) error {
	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND token != ? AND revoked = ?", userID, currentToken, false).
		Updates(map[string]interface{}{
			"revoked":     true,
			"revoked_at":  time.Now(),
		})
	return result.Error
}
