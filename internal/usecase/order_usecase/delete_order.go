package order_usecase

import (
	"context"

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

func (c *DeleteOrderUseCase) Execute(ctx context.Context, input *DeleteOrderInputDTO) error {
	return c.OrderRepository.DeleteOrder(ctx, input.Id)
}
