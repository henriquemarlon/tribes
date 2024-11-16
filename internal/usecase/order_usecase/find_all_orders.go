package order_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindAllOrdersOutputDTO []*FindOrderOutputDTO

type FindAllOrdersUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewFindAllOrdersUseCase(orderRepository entity.OrderRepository) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (f *FindAllOrdersUseCase) Execute() (*FindAllOrdersOutputDTO, error) {
	res, err := f.OrderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}
	output := make(FindAllOrdersOutputDTO, len(res))
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
