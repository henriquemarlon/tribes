package crowdfunding_usecase

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type SettleCrowdfundingInputDTO struct {
	CrowdfundingId uint `json:"crowdfunding_id"`
}

type SettleCrowdfundingOutputDTO struct {
	Id                  uint                `json:"id"`
	Token               custom_type.Address `json:"token"`
	Amount              *uint256.Int        `json:"amount"`
	Creator             custom_type.Address `json:"creator"`
	DebtIssued          *uint256.Int        `json:"debt_issued"`
	MaxInterestRate     *uint256.Int        `json:"max_interest_rate"`
	TotalObligation     *uint256.Int        `json:"total_obligation"`
	Orders              []*entity.Order     `json:"orders"`
	State               string              `json:"state"`
	FundraisingDuration int64               `json:"fundraising_duration"`
	ClosesAt            int64               `json:"closes_at"`
	MaturityAt          int64               `json:"maturity_at"`
	CreatedAt           int64               `json:"created_at"`
	UpdatedAt           int64               `json:"updated_at"`
}

type SettleCrowdfundingUseCase struct {
	UserRepository         entity.UserRepository
	ContractRepository     entity.ContractRepository
	CrowdfundingRepository entity.CrowdfundingRepository
	OrderRepository        entity.OrderRepository
}

func NewSettleCrowdfundingUseCase(
	userRepository entity.UserRepository,
	crowdfundingRepository entity.CrowdfundingRepository,
	contractRepository entity.ContractRepository,
	orderRepository entity.OrderRepository,
) *SettleCrowdfundingUseCase {
	return &SettleCrowdfundingUseCase{
		UserRepository:         userRepository,
		ContractRepository:     contractRepository,
		CrowdfundingRepository: crowdfundingRepository,
		OrderRepository:        orderRepository,
	}
}

func (uc *SettleCrowdfundingUseCase) Execute(
	ctx context.Context,
	input *SettleCrowdfundingInputDTO,
	deposit rollmelette.Deposit,
	metadata rollmelette.Metadata,
) (*SettleCrowdfundingOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type: %T", deposit)
	}

	stablecoin, err := uc.ContractRepository.FindContractBySymbol(ctx, "STABLECOIN")
	if err != nil {
		return nil, fmt.Errorf("error finding stablecoin contract: %w", err)
	}

	if custom_type.Address(erc20Deposit.Token) != stablecoin.Address {
		return nil, fmt.Errorf("token deposit is not the stablecoin %v, cannot settle crowdfunding", stablecoin.Address)
	}

	crowdfunding, err := uc.CrowdfundingRepository.FindCrowdfundingById(ctx, input.CrowdfundingId)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaign: %w", err)
	}

	if metadata.BlockTimestamp > crowdfunding.MaturityAt {
		return nil, fmt.Errorf("the maturity date of the crowdfunding campaign has passed")
	}

	if crowdfunding.State == entity.CrowdfundingStateSettled {
		return nil, fmt.Errorf("crowdfunding campaign already settled")
	}

	if crowdfunding.State != entity.CrowdfundingStateClosed {
		return nil, fmt.Errorf("crowdfunding campaign not closed")
	}

	if erc20Deposit.Amount.Cmp(crowdfunding.TotalObligation.ToBig()) < 0 {
		return nil, fmt.Errorf("deposit amount is lower than the total obligation (sum of amount and interest of all orders)")
	}

	for _, order := range crowdfunding.Orders {
		if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
			order.State = entity.OrderStateSettled
			_, err := uc.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, fmt.Errorf("error updating order: %w", err)
			}
		}
	}

	crowdfunding.State = entity.CrowdfundingStateSettled
	crowdfunding.UpdatedAt = metadata.BlockTimestamp
	res, err := uc.CrowdfundingRepository.UpdateCrowdfunding(ctx, crowdfunding)
	if err != nil {
		return nil, fmt.Errorf("error updating crowdfunding campaign: %w", err)
	}

	creator, err := uc.UserRepository.FindUserByAddress(ctx, crowdfunding.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	creator.DebtIssuanceLimit = new(uint256.Int).Sub(creator.DebtIssuanceLimit, crowdfunding.DebtIssued)
	_, err = uc.UserRepository.UpdateUser(ctx, creator)
	if err != nil {
		return nil, fmt.Errorf("error updating creator's debt issuance limit: %w", err)
	}

	return &SettleCrowdfundingOutputDTO{
		Id:                  res.Id,
		Token:               res.Token,
		Amount:              res.Amount,
		Creator:             res.Creator,
		DebtIssued:          res.DebtIssued,
		MaxInterestRate:     res.MaxInterestRate,
		TotalObligation:     res.TotalObligation,
		Orders:              res.Orders,
		State:               string(res.State),
		FundraisingDuration: res.FundraisingDuration,
		ClosesAt:            res.ClosesAt,
		MaturityAt:          res.MaturityAt,
		CreatedAt:           res.CreatedAt,
		UpdatedAt:           res.UpdatedAt,
	}, nil
}
