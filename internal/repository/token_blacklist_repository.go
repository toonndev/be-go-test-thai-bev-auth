package repository

import (
	"time"

	"be-go-test-thai-bev-auth/internal/model"

	"gorm.io/gorm"
)

type TokenBlacklistRepository interface {
	Add(token string, expiresAt time.Time) error
	IsBlacklisted(token string) (bool, error)
}

type tokenBlacklistRepository struct {
	db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) TokenBlacklistRepository {
	return &tokenBlacklistRepository{db: db}
}

func (r *tokenBlacklistRepository) Add(token string, expiresAt time.Time) error {
	return r.db.Create(&model.TokenBlacklist{
		Token:     token,
		ExpiresAt: expiresAt,
	}).Error
}

func (r *tokenBlacklistRepository) IsBlacklisted(token string) (bool, error) {
	// cleanup expired tokens
	r.db.Where("expires_at < ?", time.Now()).Delete(&model.TokenBlacklist{})

	var count int64
	err := r.db.Model(&model.TokenBlacklist{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
