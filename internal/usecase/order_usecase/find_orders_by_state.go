package order_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindOrdersByStateInputDTO struct {
	CrowdfundingId uint   `json:"crowdfunding_id"`
	State          string `json:"state"`
}

type FindOrdersByStateOutputDTO []*FindOrderOutputDTO

type FindOrdersByStateUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewFindOrdersByStateUseCase(orderRepository entity.OrderRepository) *FindOrdersByStateUseCase {
	return &FindOrdersByStateUseCase{
		OrderRepository: orderRepository,
	}
}

func (f *FindOrdersByStateUseCase) Execute(ctx context.Context, input *FindOrdersByStateInputDTO) (FindOrdersByStateOutputDTO, error) {
	res, err := f.OrderRepository.FindOrdersByState(ctx, input.CrowdfundingId, input.State)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByStateOutputDTO, len(res))
	for i, order := range res {
		output[i] = &FindOrderOutputDTO{
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
	return output, nil
}
