package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Farukcoder/go-fiber-clean-arch-auth-rbac/internal/domain"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Table("refresh_tokens").Create(token).Error
}

func (r *RefreshTokenRepository) FindByJTI(ctx context.Context, jti string) (*domain.RefreshToken, error) {
	token := &domain.RefreshToken{}
	err := r.db.WithContext(ctx).Table("refresh_tokens").Where("jti = ?", jti).Take(token).Error
	return token, err
}

func (r *RefreshTokenRepository) RevokeByJTI(ctx context.Context, jti string) error {
	result := r.db.WithContext(ctx).Table("refresh_tokens").Where("jti = ? AND revoked_at IS NULL", jti).
		Updates(map[string]any{"revoked_at": time.Now().UTC(), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeUserTokens(ctx context.Context, userID int64) error {
	result := r.db.WithContext(ctx).Table("refresh_tokens").Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]any{"revoked_at": time.Now().UTC(), "updated_at": time.Now().UTC()})
	return result.Error
}

func NewRefreshToken(userID int64, jti string, expiresAt time.Time) *domain.RefreshToken {
	return &domain.RefreshToken{
		UserID:    userID,
		JTI:       jti,
		ExpiresAt: expiresAt,
	}
}
