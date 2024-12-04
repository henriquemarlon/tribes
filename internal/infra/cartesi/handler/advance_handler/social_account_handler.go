package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/social_account_usecase"
)

type SocialAccountAdvanceHandlers struct {
	SocialAccountRepository entity.SocialAccountRepository
}

func NewSocialAccountAdvanceHandlers(socialAccountRepository entity.SocialAccountRepository) *SocialAccountAdvanceHandlers {
	return &SocialAccountAdvanceHandlers{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *SocialAccountAdvanceHandlers) CreateSocialAccountHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input social_account_usecase.CreateSocialAccountInputDTO
	err := json.Unmarshal(payload, &input)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	ctx := context.Background()
	createSocialAccount := social_account_usecase.NewCreateSocialAccountUseCase(s.SocialAccountRepository)
	res, err := createSocialAccount.Execute(ctx, &input)
	if err != nil {
		return err
	}
	socialAccount, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("social account created - "), socialAccount...))
	return nil
}

func (s *SocialAccountAdvanceHandlers) DeleteSocialAccountHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input social_account_usecase.DeleteSocialAccountInputDTO
	err := json.Unmarshal(payload, &input)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	ctx := context.Background()
	deleteSocialAccount := social_account_usecase.NewDeleteSocialAccountUseCase(s.SocialAccountRepository)
	err = deleteSocialAccount.Execute(ctx, &input)
	if err != nil {
		return err
	}
	socialAccount, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("social account deleted - "), socialAccount...))
	return nil
}
