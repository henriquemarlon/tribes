package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/crowdfunding_usecase"
	"github.com/tribeshq/tribes/pkg/router"
)

type CrowdfundingInspectHandlers struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCrowdfundingInspectHandlers(crowdfundingRepository entity.CrowdfundingRepository) *CrowdfundingInspectHandlers {
	return &CrowdfundingInspectHandlers{
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (h *CrowdfundingInspectHandlers) FindCrowdfundingByIdHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findCrowdfundingById := crowdfunding_usecase.NewFindCrowdfundingByIdUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingById.Execute(&crowdfunding_usecase.FindCrowdfundingByIdInputDTO{
		Id: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find crowdfunding: %w", err)
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal crowdfunding: %w", err)
	}
	env.Report(crowdfunding)
	return nil
}

func (h *CrowdfundingInspectHandlers) FindAllCrowdfundingsHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	findAllCrowdfundingsUseCase := crowdfunding_usecase.NewFindAllCrowdfundingsUseCase(h.CrowdfundingRepository)
	res, err := findAllCrowdfundingsUseCase.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all crowdfundings: %w", err)
	}
	allCrowdfundings, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all crowdfundings: %w", err)
	}
	env.Report(allCrowdfundings)
	return nil
}

func (h *CrowdfundingInspectHandlers) FindCrowdfundingsByInvestorHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	return nil
}