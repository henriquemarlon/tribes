package crowdfunding_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCrowdfundingsByCreatorInputDTO struct {
	Creator custom_type.Address `json:"creator"`
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
			Id:                  crowdfunding.Id,
			Token:               crowdfunding.Token,
			Amount:              crowdfunding.Amount,
			Creator:             crowdfunding.Creator,
			DebtIssued:          crowdfunding.DebtIssued,
			MaxInterestRate:     crowdfunding.MaxInterestRate,
			Orders:              orders,
			State:               string(crowdfunding.State),
			FundraisingDuration: crowdfunding.FundraisingDuration,
			ClosesAt:            crowdfunding.ClosesAt,
			MaturityAt:          crowdfunding.MaturityAt,
			CreatedAt:           crowdfunding.CreatedAt,
			UpdatedAt:           crowdfunding.UpdatedAt,
		}
	}
	return &output, nil
}
