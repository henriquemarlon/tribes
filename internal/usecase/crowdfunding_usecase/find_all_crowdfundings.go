package crowdfunding_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindAllCrowdfundingsOutputDTO []*FindCrowdfundingOutputDTO

type FindAllCrowdfundingsUseCase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewFindAllCrowdfundingsUseCase(crowdfundingRepository entity.CrowdfundingRepository) *FindAllCrowdfundingsUseCase {
	return &FindAllCrowdfundingsUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindAllCrowdfundingsUseCase) Execute(ctx context.Context) (*FindAllCrowdfundingsOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindAllCrowdfundings(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllCrowdfundingsOutputDTO, len(res))
	for i, crowdfunding := range res {
		orders := make([]*entity.Order, len(crowdfunding.Orders))
		for j, order := range crowdfunding.Orders {
			orders[j] = &entity.Order{
				Id:             order.Id,
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         order.Amount,
				InterestRate:   order.InterestRate,
				State:          order.State,
				CreatedAt:      order.CreatedAt,
				UpdatedAt:      order.UpdatedAt,
			}
		}
		output[i] = &FindCrowdfundingOutputDTO{
			Id:              crowdfunding.Id,
			Creator:         crowdfunding.Creator,
			DebtIssued:      crowdfunding.DebtIssued,
			MaxInterestRate: crowdfunding.MaxInterestRate,
			TotalObligation: crowdfunding.TotalObligation,
			State:           string(crowdfunding.State),
			Orders:          orders,
			ExpiresAt:       crowdfunding.ExpiresAt,
			MaturityAt:      crowdfunding.MaturityAt,
			CreatedAt:       crowdfunding.CreatedAt,
			UpdatedAt:       crowdfunding.UpdatedAt,
		}
	}
	return &output, nil
}
