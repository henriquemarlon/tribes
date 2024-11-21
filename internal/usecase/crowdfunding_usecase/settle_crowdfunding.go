package crowdfunding_usecase

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type SettleCrowdfundingInputDTO struct {
	CrowdfundingId uint `json:"crowdfunding_id"`
}

type SettleCrowdfundingOutputDTO struct {
	Id              uint            `json:"id"`
	Creator         common.Address  `json:"creator"`
	DebtIssued      *uint256.Int    `json:"debt_issued"`
	MaxInterestRate *uint256.Int    `json:"max_interest_rate"`
	TotalObligation *uint256.Int    `json:"total_obligation"`
	State           string          `json:"state"`
	Orders          []*entity.Order `json:"orders"`
	ExpiresAt       int64           `json:"expires_at"`
	MaturityAt      int64           `json:"maturity_at"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`
}

type SettleCrowdfundingUseCase struct {
	UserRepository         entity.UserRepository
	ContractRepository     entity.ContractRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewSettleCrowdfundingUseCase(userRepository entity.UserRepository, crowdfundingRepository entity.CrowdfundingRepository, contractRepository entity.ContractRepository) *SettleCrowdfundingUseCase {
	return &SettleCrowdfundingUseCase{
		UserRepository:         userRepository,
		ContractRepository:     contractRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (uc *SettleCrowdfundingUseCase) Execute(input *SettleCrowdfundingInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*SettleCrowdfundingOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type: %T", deposit)
	}

	stablecoin, err := uc.ContractRepository.FindContractBySymbol("STABLECOIN")
	if err != nil {
		return nil, fmt.Errorf("error finding stablecoin contract: %w", err)
	}

	if erc20Deposit.Token != stablecoin.Address {
		return nil, fmt.Errorf("token deposit is not the same as the stablecoin, cannot settle crowdfunding")
	}

	crowdfunding, err := uc.CrowdfundingRepository.FindCrowdfundingById(input.CrowdfundingId)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaign: %w", err)
	}

	switch crowdfunding.State {
	case entity.CrowdfundingStateSettled:
		return nil, fmt.Errorf("crowdfunding campaign already settled")
	case entity.CrowdfundingStateClosed:
		if erc20Deposit.Amount.Cmp(crowdfunding.TotalObligation.ToBig()) != 0 {
			return nil, fmt.Errorf("cannot settle crowdfunding because the deposit amount is not equal to the Total Obligation (sum of amount + interest of all orders)")
		}

		for _, order := range crowdfunding.Orders {
			if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
				order.State = entity.OrderStateSettled
			}
		}

		crowdfunding.State = entity.CrowdfundingStateSettled
		crowdfunding.UpdatedAt = metadata.BlockTimestamp
		res, err := uc.CrowdfundingRepository.UpdateCrowdfunding(crowdfunding)
		if err != nil {
			return nil, fmt.Errorf("error updating crowdfunding campaign: %w", err)
		}

		creator, err := uc.UserRepository.FindUserByAddress(crowdfunding.Creator)
		if err != nil {
			return nil, fmt.Errorf("error finding creator: %w", err)
		}

		creator.DebtIssuanceLimit = new(uint256.Int).Sub(creator.DebtIssuanceLimit, crowdfunding.DebtIssued)
		_, err = uc.UserRepository.UpdateUser(creator)
		if err != nil {
			return nil, fmt.Errorf("error updating creator debt issuance limit: %w", err)
		}

		return &SettleCrowdfundingOutputDTO{
			Id:              res.Id,
			Creator:         res.Creator,
			DebtIssued:      res.DebtIssued,
			MaxInterestRate: res.MaxInterestRate,
			TotalObligation: res.TotalObligation,
			State:           string(res.State),
			Orders:          res.Orders,
			ExpiresAt:       res.ExpiresAt,
			MaturityAt:      res.MaturityAt,
			CreatedAt:       res.CreatedAt,
			UpdatedAt:       res.UpdatedAt,
		}, nil
	default:
		return nil, fmt.Errorf("crowdfunding campaign not closed")
	}
}
