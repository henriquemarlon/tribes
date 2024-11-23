package crowdfunding_usecase

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindCrowdfundingsByCreatorInputDTO struct {
	Creator common.Address `json:"creator"`
}

type FindCrowdfundingsByCreatorOutputDTO []*FindCrowdfundingOutputDTO

type FindCrowdfundingsByCreatorUseCase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewFindCrowdfundingsByCreatorUseCase(crowdfundingRepository entity.CrowdfundingRepository) *FindCrowdfundingsByCreatorUseCase {
	return &FindCrowdfundingsByCreatorUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindCrowdfundingsByCreatorUseCase) Execute(ctx context.Context, input *FindCrowdfundingsByCreatorInputDTO) (*FindCrowdfundingsByCreatorOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}
	output := make(FindCrowdfundingsByCreatorOutputDTO, len(res))
	for i, crowdfunding := range res {
		orders := make([]*FindCrowdfundingOutputSubDTO, len(crowdfunding.Orders))
		for j, order := range crowdfunding.Orders {
			orders[j] = &FindCrowdfundingOutputSubDTO{
				Id:             order.Id,
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         order.Amount,
				InterestRate:   order.InterestRate,
				State:          string(order.State),
				CreatedAt:      order.CreatedAt,
				UpdatedAt:      order.UpdatedAt,
			}
		}
		output[i] = &FindCrowdfundingOutputDTO{
			Id:              crowdfunding.Id,
			DebtIssued:      crowdfunding.DebtIssued,
			MaxInterestRate: crowdfunding.MaxInterestRate,
			State:           string(crowdfunding.State),
			Orders:          orders,
			ExpiresAt:       crowdfunding.ExpiresAt,
			CreatedAt:       crowdfunding.CreatedAt,
			UpdatedAt:       crowdfunding.UpdatedAt,
		}
	}
	return &output, nil
}
