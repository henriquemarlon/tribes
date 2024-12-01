package order_usecase

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CancelOrderInputDTO struct {
	Id uint
}

type CancelOrderOutputDTO struct {
	Id             uint
	CrowdfundingId uint
	Investor       common.Address
	Amount         *uint256.Int
	InterestRate   *uint256.Int
	State          string
	CreatedAt      int64
	UpdatedAt      int64
}

type CancelOrderUseCase struct {
	UserRepository         entity.UserRepository
	OrderRepository        entity.OrderRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCancelOrderUseCase(userRepository entity.UserRepository, orderRepository entity.OrderRepository, crowdfundingRepository entity.CrowdfundingRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		UserRepository:         userRepository,
		OrderRepository:        orderRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (c *CancelOrderUseCase) Execute(ctx context.Context, input *CancelOrderInputDTO, metadata rollmelette.Metadata) (*CancelOrderOutputDTO, error) {
	order, err := c.OrderRepository.FindOrderById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	if order.Investor != metadata.MsgSender {
		return nil, errors.New("only the investor can cancel the order")
	}
	crowdfunding, err := c.CrowdfundingRepository.FindCrowdfundingById(ctx, order.CrowdfundingId)
	if err != nil {
		return nil, err
	}
	if crowdfunding.ClosesAt < metadata.BlockTimestamp {
		return nil, errors.New("cannot cancel order after crowdfunding closes")
	}
	err = c.OrderRepository.DeleteOrder(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return &CancelOrderOutputDTO{
		Id:             order.Id,
		CrowdfundingId: order.CrowdfundingId,
		Investor:       order.Investor,
		Amount:         order.Amount,
		InterestRate:   order.InterestRate,
		State:          string(order.State),
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}, nil
}
