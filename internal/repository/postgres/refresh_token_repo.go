package postgres

import (
	"context"
	"time"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// RefreshTokenRepo is refresh-token repository.
type RefreshTokenRepo struct{ db *gorm.DB }

// NewRefreshTokenRepo creates refresh-token repository.
func NewRefreshTokenRepo(db *gorm.DB) *RefreshTokenRepo { return &RefreshTokenRepo{db: db} }

// Create inserts refresh token.
func (r *RefreshTokenRepo) Create(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByHash returns refresh token by hash.
func (r *RefreshTokenRepo) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// Revoke marks refresh token revoked.
func (r *RefreshTokenRepo) Revoke(ctx context.Context, token *domain.RefreshToken) error {
	now := time.Now()
	token.RevokedAt = &now
	return r.db.WithContext(ctx).Save(token).Error
}
