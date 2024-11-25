package crowdfunding_usecase

import (
	"context"

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

func (f *FindCrowdfundingByIdUseCase) Execute(ctx context.Context, input *FindCrowdfundingByIdInputDTO) (*FindCrowdfundingOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindCrowdfundingById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Order
	for _, order := range res.Orders {
		orders = append(orders, &entity.Order{
			Id:             order.Id,
			CrowdfundingId: order.CrowdfundingId,
			Investor:       order.Investor,
			Amount:         order.Amount,
			InterestRate:   order.InterestRate,
			State:          order.State,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		})
	}
	return &FindCrowdfundingOutputDTO{
		Id:              res.Id,
		Creator:         res.Creator,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		TotalObligation: res.TotalObligation,
		State:           string(res.State),
		Orders:          orders,
		ClosesAt:        res.ClosesAt,
		MaturityAt:      res.MaturityAt,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}
