package social_account_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindSocialAccountByUserIdInputDTO struct {
	UserId uint `json:"user_id"`
}

type FindSocialAccountsByUserIdOutputDTO []*FindSocialAccountOutputDTO

type FindSocialAccountsByUserIdUseCase struct {
	SocialAccountRepository entity.SocialAccountRepository
}

func NewFindSocialAccountsByUserIdUseCase(socialAccountRepository entity.SocialAccountRepository) *FindSocialAccountsByUserIdUseCase {
	return &FindSocialAccountsByUserIdUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *FindSocialAccountsByUserIdUseCase) Execute(ctx context.Context, input *FindSocialAccountByUserIdInputDTO) (*FindSocialAccountsByUserIdOutputDTO, error) {
	socialAccounts, err := s.SocialAccountRepository.FindSocialAccountsByUserId(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	output := make(FindSocialAccountsByUserIdOutputDTO, len(socialAccounts))
	for i, socialAccount := range socialAccounts {
		output[i] = &FindSocialAccountOutputDTO{
			Id:        socialAccount.Id,
			UserId:    socialAccount.UserId,
			Username:  socialAccount.Username,
			Followers: socialAccount.Followers,
			Platform:  string(socialAccount.Platform),
			CreatedAt: socialAccount.CreatedAt,
			UpdatedAt: socialAccount.UpdatedAt,
		}
	}
	return &output, nil
}
