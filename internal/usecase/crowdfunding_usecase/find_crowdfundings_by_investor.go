package crowdfunding_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindCrowdfundingsByInvestorInputDTO struct {
	Investor common.Address `json:"investor"`
}

type FindCrowdfundingsByInvestorOutputDTO []*FindCrowdfundingOutputDTO


type FindCrowdfundingByInvestorUseCase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewFindCrowdfundingByInvestorUseCase(crowdfundingRepository entity.CrowdfundingRepository) *FindCrowdfundingByInvestorUseCase {
	return &FindCrowdfundingByInvestorUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindCrowdfundingByInvestorUseCase) Execute(input *FindCrowdfundingsByInvestorInputDTO) (*FindCrowdfundingsByInvestorOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindCrowdfundingsByInvestor(input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindCrowdfundingsByInvestorOutputDTO, len(res))
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
