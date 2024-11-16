package order_usecase

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateOrderInputDTO struct {
	Creator common.Address `json:"creator"`
	Price   *uint256.Int    `json:"interest_rate"`
}

type CreateOrderOutputDTO struct {
	Id             uint           `json:"id"`
	CrowdfundingId uint           `json:"crowdfunding_id"`
	Investor       common.Address `json:"investor"`
	Amount         *uint256.Int    `json:"amount"`
	InterestRate   *uint256.Int    `json:"interest_rate"`
	State          string         `json:"state"`
	CreatedAt      int64          `json:"created_at"`
}

type CreateOrderUseCase struct {
	OrderRepository        entity.OrderRepository
	ContractRepository     entity.ContractRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepository, contractRepository entity.ContractRepository, crowdfundingRepository entity.CrowdfundingRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:        orderRepository,
		ContractRepository:     contractRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (c *CreateOrderUseCase) Execute(input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(input.Creator)
	if err != nil {
		return nil, err
	}
	var activeCrowdfunding *entity.Crowdfunding
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State == entity.CrowdfundingStateOngoing {
			activeCrowdfunding = crowdfunding
		}
	}
	if activeCrowdfunding == nil {
		return nil, fmt.Errorf("no active crowdfunding found, cannot create order for creator: %v", input.Creator)
	}

	if metadata.BlockTimestamp > activeCrowdfunding.ExpiresAt {
		return nil, fmt.Errorf("active crowdfunding expired, cannot create order")
	}
	stablecoin, err := c.ContractRepository.FindContractBySymbol("STABLECOIN")
	if err != nil {
		return nil, err
	}
	if deposit.(*rollmelette.ERC20Deposit).Token != stablecoin.Address {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", deposit.(*rollmelette.ERC20Deposit).Token)
	}

	if input.Price.Gt(activeCrowdfunding.MaxInterestRate) {
		return nil, fmt.Errorf("order price exceeds active crowdfunding max interest rate")
	}

	order, err := entity.NewOrder(activeCrowdfunding.Id, deposit.(*rollmelette.ERC20Deposit).Sender, uint256.MustFromBig(deposit.(*rollmelette.ERC20Deposit).Amount), input.Price, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.OrderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	return &CreateOrderOutputDTO{
		Id:             res.Id,
		CrowdfundingId: res.CrowdfundingId,
		Investor:       res.Investor,
		Amount:         res.Amount,
		InterestRate:   res.InterestRate,
		State:          string(res.State),
		CreatedAt:      res.CreatedAt,
	}, nil
}
