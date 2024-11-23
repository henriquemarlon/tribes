package order_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindOrdersByCrowdfundingIdInputDTO struct {
	CrowdfundingId uint `json:"crowdfunding_id"`
}

type FindOrdersByCrowdfundingIdOutputDTO []*FindOrderOutputDTO

type FindOrdersByCrowdfundingIdUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewFindOrdersByCrowdfundingIdUseCase(orderRepository entity.OrderRepository) *FindOrdersByCrowdfundingIdUseCase {
	return &FindOrdersByCrowdfundingIdUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *FindOrdersByCrowdfundingIdUseCase) Execute(ctx context.Context, input *FindOrdersByCrowdfundingIdInputDTO) (*FindOrdersByCrowdfundingIdOutputDTO, error) {
	res, err := c.OrderRepository.FindOrdersByCrowdfundingId(ctx, input.CrowdfundingId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByCrowdfundingIdOutputDTO, len(res))
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
	return &output, nil
}
