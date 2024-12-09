package repository

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/gorm"
)

type SocialAccountRepositorySqlite struct {
	Db *gorm.DB
}

func NewSocialAccountRepositorySqlite(db *gorm.DB) *SocialAccountRepositorySqlite {
	return &SocialAccountRepositorySqlite{Db: db}
}

func (r *SocialAccountRepositorySqlite) CreateSocialAccount(ctx context.Context, input *entity.SocialAccount) (*entity.SocialAccount, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create social account: %w", err)
	}
	return input, nil
}

func (r *SocialAccountRepositorySqlite) FindSocialAccountById(ctx context.Context, id uint) (*entity.SocialAccount, error) {
	var account entity.SocialAccount
	if err := r.Db.WithContext(ctx).First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("social account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find social account by ID: %w", err)
	}
	return &account, nil
}

func (r *SocialAccountRepositorySqlite) FindSocialAccountsByUserId(ctx context.Context, userID uint) ([]*entity.SocialAccount, error) {
	var accounts []*entity.SocialAccount
	if err := r.Db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to find social accounts by user ID: %w", err)
	}
	return accounts, nil
}

func (r *SocialAccountRepositorySqlite) DeleteSocialAccount(ctx context.Context, id uint) error {
	res := r.Db.WithContext(ctx).Delete(&entity.SocialAccount{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete social account: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("social account not found")
	}
	return nil
}
