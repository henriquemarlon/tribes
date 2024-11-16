package order_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type DeleteOrderInputDTO struct {
	Id uint `json:"id"`
}

type DeleteOrderUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewDeleteOrderUseCase(orderRepository entity.OrderRepository) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *DeleteOrderUseCase) Execute(input *DeleteOrderInputDTO) error {
	return c.OrderRepository.DeleteOrder(input.Id)
}
