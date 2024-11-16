package order_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindOrderByIdInputDTO struct {
	Id uint `json:"id"`
}

type FindOrderByIdUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewFindOrderByIdUseCase(orderRepository entity.OrderRepository) *FindOrderByIdUseCase {
	return &FindOrderByIdUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *FindOrderByIdUseCase) Execute(input *FindOrderByIdInputDTO) (*FindOrderOutputDTO, error) {
	res, err := c.OrderRepository.FindOrderById(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindOrderOutputDTO{
		Id:             res.Id,
		CrowdfundingId: res.CrowdfundingId,
		Investor:       res.Investor,
		Amount:         res.Amount,
		InterestRate:   res.InterestRate,
		State:          string(res.State),
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}, nil
}
