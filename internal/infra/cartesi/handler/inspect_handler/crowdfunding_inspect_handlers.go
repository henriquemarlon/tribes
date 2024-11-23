package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
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

func (h *CrowdfundingInspectHandlers) FindCrowdfundingByIdHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findCrowdfundingById := crowdfunding_usecase.NewFindCrowdfundingByIdUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingById.Execute(ctx, &crowdfunding_usecase.FindCrowdfundingByIdInputDTO{
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

func (h *CrowdfundingInspectHandlers) FindAllCrowdfundingsHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	findAllCrowdfundingsUseCase := crowdfunding_usecase.NewFindAllCrowdfundingsUseCase(h.CrowdfundingRepository)
	res, err := findAllCrowdfundingsUseCase.Execute(ctx)
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

func (h *CrowdfundingInspectHandlers) FindCrowdfundingsByInvestorHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	findCrowdfundingsByInvestor := crowdfunding_usecase.NewFindCrowdfundingsByInvestorUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingsByInvestor.Execute(ctx, &crowdfunding_usecase.FindCrowdfundingsByInvestorInputDTO{
		Investor: common.HexToAddress(router.PathValue(ctx, "address")),
	})
	if err != nil {
		return fmt.Errorf("failed to find crowdfundings by investor: %w", err)
	}
	crowdfundings, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal crowdfundings: %w", err)
	}
	env.Report(crowdfundings)
	return nil
}

func (h *CrowdfundingInspectHandlers) FindCrowdfundingsByCreatorHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	findCrowdfundingsByCreator := crowdfunding_usecase.NewFindCrowdfundingsByCreatorUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingsByCreator.Execute(ctx, &crowdfunding_usecase.FindCrowdfundingsByCreatorInputDTO{
		Creator: common.HexToAddress(router.PathValue(ctx, "address")),
	})
	if err != nil {
		return fmt.Errorf("failed to find crowdfundings by creator: %w", err)
	}
	crowdfundings, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal crowdfundings: %w", err)
	}
	env.Report(crowdfundings)
	return nil
}
