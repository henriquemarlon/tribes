package entity

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidSocialAccount  = errors.New("invalid social account")
	ErrSocialAccountNotFound = errors.New("social account not found")
)

type Platform string

const (
	PlatformTwitter   Platform = "twitter"
	PlatformInstagram Platform = "instagram"
)

type SocialAccountRepository interface {
	CreateSocialAccount(ctx context.Context, socialAccount *SocialAccount) (*SocialAccount, error)
	FindSocialAccountById(ctx context.Context, id uint) (*SocialAccount, error)
	FindSocialAccountsByUserId(ctx context.Context, userID uint) ([]*SocialAccount, error)
	DeleteSocialAccount(ctx context.Context, id uint) error
}

type SocialAccount struct {
	Id        uint     `json:"id" gorm:"primaryKey"`
	UserId    uint     `json:"user_id,omitempty" gorm:"not null"`
	Username  string   `json:"username,omitempty" gorm:"type:text;not null;uniqueIndex:idx_username_platform"`
	Followers uint     `json:"followers,omitempty" gorm:"not null"`
	Platform  Platform `json:"platform,omitempty" gorm:"not null;uniqueIndex:idx_username_platform"`
	CreatedAt int64    `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt int64    `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewSocialAccount(userID uint, username string, followers uint, platform string, createdAt int64) (*SocialAccount, error) {
	socialAccount := &SocialAccount{
		UserId:    userID,
		Username:  username,
		Followers: followers,
		Platform:  Platform(platform),
		CreatedAt: createdAt,
	}
	if err := socialAccount.Validate(); err != nil {
		return nil, err
	}
	return socialAccount, nil
}

func (s *SocialAccount) Validate() error {
	if s.UserId == 0 {
		return fmt.Errorf("%w: user ID cannot be zero", ErrInvalidSocialAccount)
	}
	if s.Username == "" {
		return fmt.Errorf("%w: username cannot be empty", ErrInvalidSocialAccount)
	}
	if s.Followers == 0 {
		return fmt.Errorf("%w: followers cannot be zero", ErrInvalidSocialAccount)
	}
	if s.Platform == "" {
		return fmt.Errorf("%w: platform cannot be empty", ErrInvalidSocialAccount)
	}
	if s.Platform != PlatformTwitter && s.Platform != PlatformInstagram {
		return fmt.Errorf("%w: platform must be 'twitter' or 'instagram'", ErrInvalidSocialAccount)
	}
	if s.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidSocialAccount)
	}
	return nil
}
