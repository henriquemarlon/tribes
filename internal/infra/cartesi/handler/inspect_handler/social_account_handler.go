package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/social_account_usecase"
	"github.com/tribeshq/tribes/pkg/router"
)

type SocialAccountInspectHandlers struct {
	SocialAccountRepository entity.SocialAccountRepository
}

func NewSocialAccountInspectHandlers(socialAccountRepository entity.SocialAccountRepository) *SocialAccountInspectHandlers {
	return &SocialAccountInspectHandlers{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (h *SocialAccountInspectHandlers) FindSocialAccountById(ctx context.Context, env rollmelette.EnvInspector) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findSocialAccountById := social_account_usecase.NewFindSocialAccountByIdUseCase(h.SocialAccountRepository)
	res, err := findSocialAccountById.Execute(ctx, &social_account_usecase.FindSocialAccountByIDInputDTO{
		SocialAccountId: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find social account: %w", err)
	}
	socialAccount, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("")
	}
	env.Report(socialAccount)
	return nil
}

func (h *SocialAccountInspectHandlers) FindSocialAccountsByUserId(ctx context.Context, env rollmelette.EnvInspector) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findSocialAccountsByUserId := social_account_usecase.NewFindSocialAccountsByUserIdUseCase(h.SocialAccountRepository)
	res, err := findSocialAccountsByUserId.Execute(ctx, &social_account_usecase.FindSocialAccountByUserIdInputDTO{
		UserId: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find social accounts %w for this id %v", err, id)
	}
	socialAccounts, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal social accounts: %w", err)
	}
	env.Report(socialAccounts)
	return nil
}
