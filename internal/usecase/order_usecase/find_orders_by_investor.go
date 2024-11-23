package order_usecase

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FinsOrdersByInvestorInputDTO struct {
	Investor common.Address `json:"investor"`
}

type FindOrdersByInvestorOutputDTO []*FindOrderOutputDTO

type FindOrdersByInvestorUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewFindOrdersByInvestorUseCase(orderRepository entity.OrderRepository) *FindOrdersByInvestorUseCase {
	return &FindOrdersByInvestorUseCase{
		OrderRepository: orderRepository,
	}
}

func (o *FindOrdersByInvestorUseCase) Execute(ctx context.Context, input *FinsOrdersByInvestorInputDTO) (FindOrdersByInvestorOutputDTO, error) {
	res, err := o.OrderRepository.FindOrdersByInvestor(ctx, input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByInvestorOutputDTO, len(res))
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
