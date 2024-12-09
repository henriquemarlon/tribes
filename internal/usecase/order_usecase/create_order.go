package order_usecase

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateOrderInputDTO struct {
	CrowndfundingId uint         `json:"crowdfunding_id"`
	InterestRate    *uint256.Int `json:"interest_rate"`
}

type CreateOrderOutputDTO struct {
	Id             uint                `json:"id"`
	CrowdfundingId uint                `json:"crowdfunding_id"`
	Investor       custom_type.Address `json:"investor"`
	Amount         *uint256.Int        `json:"amount"`
	InterestRate   *uint256.Int        `json:"interest_rate"`
	State          string              `json:"state"`
	CreatedAt      int64               `json:"created_at"`
}

type CreateOrderUseCase struct {
	UserRepository         entity.UserRepository
	OrderRepository        entity.OrderRepository
	ContractRepository     entity.ContractRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCreateOrderUseCase(userRepository entity.UserRepository, orderRepository entity.OrderRepository, contractRepository entity.ContractRepository, crowdfundingRepository entity.CrowdfundingRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		UserRepository:         userRepository,
		OrderRepository:        orderRepository,
		ContractRepository:     contractRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (c *CreateOrderUseCase) Execute(ctx context.Context, input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type provided for order creation: %T", deposit)
	}

	user, err := c.UserRepository.FindUserByAddress(ctx, custom_type.Address(erc20Deposit.Sender))
	if user == nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// According with the CVM Resolution 88
	depositAmount := uint256.MustFromBig(erc20Deposit.Amount)
	if user.InvestmentLimit.Cmp(depositAmount) < 0 {
		return nil, fmt.Errorf("investor limit exceeded, cannot create order")
	}

	// According with the CVM Resolution 88
	if user.Role != entity.UserRoleNonQualifiedInvestor && user.Role != entity.UserRoleQualifiedInvestor {
		return nil, fmt.Errorf("user role not allowed to create order: %v", user.Role)
	}

	crowdfunding, err := c.CrowdfundingRepository.FindCrowdfundingById(ctx, input.CrowndfundingId)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaigns: %w", err)
	}

	if crowdfunding.ClosesAt-crowdfunding.FundraisingDuration > metadata.BlockTimestamp {
		return nil, fmt.Errorf("crowdfunding campaign not yet open, order cannot be placed")
	}

	if crowdfunding.ClosesAt < metadata.BlockTimestamp {
		return nil, fmt.Errorf("crowdfunding campaign closed, order cannot be placed")
	}

	stablecoin, err := c.ContractRepository.FindContractBySymbol(ctx, "STABLECOIN")
	if err != nil {
		return nil, fmt.Errorf("error finding stablecoin contract: %w", err)
	}
	if custom_type.Address(erc20Deposit.Token) != stablecoin.Address {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(crowdfunding.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active crowdfunding max interest rate")
	}

	order, err := entity.NewOrder(crowdfunding.Id, custom_type.Address(erc20Deposit.Sender), depositAmount, input.InterestRate, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	user.InvestmentLimit.Sub(user.InvestmentLimit, order.Amount)
	_, err = c.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error decreasing creator investment limit: %w", err)
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
