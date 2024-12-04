package social_account_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateSocialAccountInputDTO struct {
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Followers uint   `json:"followers"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
}

type CreateSocialAccountOutputDTO struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Followers uint   `json:"followers"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
}

type CreateSocialAccountUseCase struct {
	SocialAccountRepository entity.SocialAccountRepository
}

func NewCreateSocialAccountUseCase(socialAccountRepository entity.SocialAccountRepository) *CreateSocialAccountUseCase {
	return &CreateSocialAccountUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *CreateSocialAccountUseCase) Execute(ctx context.Context, input *CreateSocialAccountInputDTO) (*CreateSocialAccountOutputDTO, error) {
	socialAccount, err := entity.NewSocialAccount(input.UserId, input.Username, input.Followers, input.Platform, int64(input.CreatedAt))
	if err != nil {
		return nil, err
	}
	socialAccount, err = s.SocialAccountRepository.CreateSocialAccount(ctx, socialAccount)
	if err != nil {
		return nil, err
	}
	return &CreateSocialAccountOutputDTO{
		Id:        socialAccount.Id,
		UserId:    socialAccount.UserId,
		Username:  socialAccount.Username,
		Followers: socialAccount.Followers,
		Platform:  string(socialAccount.Platform),
		CreatedAt: socialAccount.CreatedAt,
	}, nil
}
