package crowdfunding_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindAllCrowdfundingsOutputDTO []*FindCrowdfundingOutputDTO

type FindAllCrowdfundingsUseCase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewFindAllCrowdfundingsUseCase(crowdfundingRepository entity.CrowdfundingRepository) *FindAllCrowdfundingsUseCase {
	return &FindAllCrowdfundingsUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindAllCrowdfundingsUseCase) Execute() (*FindAllCrowdfundingsOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindAllCrowdfundings()
	if err != nil {
		return nil, err
	}
	output := make(FindAllCrowdfundingsOutputDTO, len(res))
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
