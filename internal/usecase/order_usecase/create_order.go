package order_usecase

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateOrderInputDTO struct {
	Creator      common.Address `json:"creator"`
	InterestRate *uint256.Int   `json:"interest_rate"`
}

type CreateOrderOutputDTO struct {
	Id             uint           `json:"id"`
	CrowdfundingId uint           `json:"crowdfunding_id"`
	Investor       common.Address `json:"investor"`
	Amount         *uint256.Int   `json:"amount"`
	InterestRate   *uint256.Int   `json:"interest_rate"`
	State          string         `json:"state"`
	CreatedAt      int64          `json:"created_at"`
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

	user, err := c.UserRepository.FindUserByAddress(ctx, erc20Deposit.Sender)
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

	crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaigns: %w", err)
	}

	var activeCrowdfunding *entity.Crowdfunding
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State == entity.CrowdfundingStateOngoing {
			if metadata.BlockTimestamp > crowdfunding.ClosesAt {
				return nil, fmt.Errorf("active crowdfunding expired, cannot create order")
			}
			activeCrowdfunding = crowdfunding
			break
		}
	}
	if activeCrowdfunding == nil {
		return nil, fmt.Errorf("no active crowdfunding found for creator: %v", input.Creator)
	}

	stablecoin, err := c.ContractRepository.FindContractBySymbol(ctx, "STABLECOIN")
	if err != nil {
		return nil, fmt.Errorf("error finding stablecoin contract: %w", err)
	}
	if erc20Deposit.Token != stablecoin.Address {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(activeCrowdfunding.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active crowdfunding max interest rate")
	}

	order, err := entity.NewOrder(activeCrowdfunding.Id, erc20Deposit.Sender, depositAmount, input.InterestRate, metadata.BlockTimestamp)
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
