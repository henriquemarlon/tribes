package social_account_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type DeleteSocialAccountInputDTO struct {
	SocialAccountId uint `json:"social_account_id"`
}

type DeleteSocialAccountUseCase struct {
	SocialAccountRepository entity.SocialAccountRepository
}

func NewDeleteSocialAccountUseCase(socialAccountRepository entity.SocialAccountRepository) *DeleteSocialAccountUseCase {
	return &DeleteSocialAccountUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *DeleteSocialAccountUseCase) Execute(ctx context.Context, input *DeleteSocialAccountInputDTO) error {
	return s.SocialAccountRepository.DeleteSocialAccount(ctx, input.SocialAccountId)
}
