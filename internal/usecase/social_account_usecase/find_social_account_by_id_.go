package social_account_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindSocialAccountByIDInputDTO struct {
	SocialAccountId uint `json:"social_account_id"`
}

type FindSocialAccountByIdUseCase struct {
	SocialAccountRepository entity.SocialAccountRepository
}

func NewFindSocialAccountByIdUseCase(socialAccountRepository entity.SocialAccountRepository) *FindSocialAccountByIdUseCase {
	return &FindSocialAccountByIdUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *FindSocialAccountByIdUseCase) Execute(ctx context.Context, input *FindSocialAccountByIDInputDTO) (*FindSocialAccountOutputDTO, error) {
	socialAccount, err := s.SocialAccountRepository.FindSocialAccountById(ctx, input.SocialAccountId)
	if err != nil {
		return nil, err
	}
	return &FindSocialAccountOutputDTO{
		Id:        socialAccount.Id,
		UserId:    socialAccount.UserId,
		Username:  socialAccount.Username,
		Followers: socialAccount.Followers,
		Platform:  string(socialAccount.Platform),
		CreatedAt: socialAccount.CreatedAt,
		UpdatedAt: socialAccount.UpdatedAt,
	}, nil
}
