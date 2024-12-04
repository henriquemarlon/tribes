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
	return &SocialAccountRepositorySqlite{
		Db: db,
	}
}

func (r *SocialAccountRepositorySqlite) CreateSocialAccount(ctx context.Context, input *entity.SocialAccount) (*entity.SocialAccount, error) {
	err := r.Db.WithContext(ctx).Raw(`
		INSERT INTO social_accounts (user_id, username, followers, platform, created_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`, input.UserId, input.Username, input.Followers, input.Platform, input.CreatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create social account: %w", err)
	}
	return input, nil
}

func (r *SocialAccountRepositorySqlite) FindSocialAccountById(ctx context.Context, id uint) (*entity.SocialAccount, error) {
	var result map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, user_id, username, followers, platform
		FROM social_accounts
		WHERE id = ? LIMIT 1
	`, id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("social account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find social account by ID: %w", err)
	}

	return &entity.SocialAccount{
		Id:        uint(result["id"].(int64)),
		UserId:    uint(result["user_id"].(int64)),
		Username:  result["username"].(string),
		Followers: uint(result["followers"].(int64)),
		Platform:  entity.Platform(result["platform"].(string)),
	}, nil
}

func (r *SocialAccountRepositorySqlite) FindSocialAccountsByUserId(ctx context.Context, userID uint) ([]*entity.SocialAccount, error) {
	var results []map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, user_id, username, followers, platform
		FROM social_accounts
		WHERE user_id = ?
	`, userID).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find social accounts by user ID: %w", err)
	}

	var accounts []*entity.SocialAccount
	for _, data := range results {
		accounts = append(accounts, &entity.SocialAccount{
			Id:        uint(data["id"].(int64)),
			UserId:    uint(data["user_id"].(int64)),
			Username:  data["username"].(string),
			Followers: uint(data["followers"].(int64)),
			Platform:  entity.Platform(data["platform"].(string)),
		})
	}
	return accounts, nil
}

func (r *SocialAccountRepositorySqlite) DeleteSocialAccount(ctx context.Context, id uint) error {
	res := r.Db.WithContext(ctx).Delete(&entity.SocialAccount{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete social account: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("social account not found")
	}
	return nil
}
