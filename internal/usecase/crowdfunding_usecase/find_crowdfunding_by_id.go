package crowdfunding_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindCrowdfundingByIdInputDTO struct {
	Id uint `json:"id"`
}

type FindCrowdfundingByIdUseCase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewFindCrowdfundingByIdUseCase(crowdfundingRepository entity.CrowdfundingRepository) *FindCrowdfundingByIdUseCase {
	return &FindCrowdfundingByIdUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindCrowdfundingByIdUseCase) Execute(input *FindCrowdfundingByIdInputDTO) (*FindCrowdfundingOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindCrowdfundingById(input.Id)
	if err != nil {
		return nil, err
	}
	var orders []*FindCrowdfundingOutputSubDTO
	for _, order := range res.Orders {
		orders = append(orders, &FindCrowdfundingOutputSubDTO{
			Id:             order.Id,
			CrowdfundingId: order.CrowdfundingId,
			Investor:       order.Investor,
			Amount:         order.Amount,
			InterestRate:   order.InterestRate,
			State:          string(order.State),
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		})
	}
	return &FindCrowdfundingOutputDTO{
		Id:              res.Id,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		State:           string(res.State),
		Orders:          orders,
		ExpiresAt:       res.ExpiresAt,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}
