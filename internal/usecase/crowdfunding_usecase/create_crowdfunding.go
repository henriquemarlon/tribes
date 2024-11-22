package crowdfunding_usecase

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateCrowdfundingInputDTO struct {
	DebitIssued     *uint256.Int `json:"debt_issued"`
	MaxInterestRate *uint256.Int `json:"max_interest_rate"`
	ExpiresAt       int64        `json:"expires_at"`
	MaturityAt      int64        `json:"maturity_at"`
}

type CreateCrowdfundingOutputDTO struct {
	Id              uint            `json:"id"`
	Creator         common.Address  `json:"creator,omitempty"`
	DebtIssued      *uint256.Int    `json:"debt_issued"`
	MaxInterestRate *uint256.Int    `json:"max_interest_rate"`
	State           string          `json:"state"`
	Orders          []*entity.Order `json:"orders"`
	ExpiresAt       int64           `json:"expires_at"`
	MaturityAt      int64           `json:"maturity_at"`
	CreatedAt       int64           `json:"created_at"`
}

type CreateCrowdfundingUseCase struct {
	UserRepository         entity.UserRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCreateCrowdfundingUseCase(userRepository entity.UserRepository, crowdfundingRepository entity.CrowdfundingRepository) *CreateCrowdfundingUseCase {
	return &CreateCrowdfundingUseCase{
		UserRepository:         userRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (c *CreateCrowdfundingUseCase) Execute(input *CreateCrowdfundingInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateCrowdfundingOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type: %T", deposit)
	}

	creator, err := c.UserRepository.FindUserByAddress(erc20Deposit.Sender)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	// Validate debt issuance limit
	if creator.DebtIssuanceLimit.Cmp(input.DebitIssued) < 0 {
		return nil, fmt.Errorf("creator debt issuance limit exceeded")
	}

	crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(creator.Address)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaigns: %w", err)
	}

	// Check for active crowdfunding campaigns within the last 120 days
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State != entity.CrowdfundingStateSettled && metadata.BlockTimestamp-crowdfunding.CreatedAt < 120*24*60*60 {
			return nil, fmt.Errorf("creator already has an active crowdfunding within the last 120 days")
		}
	}

	crowdfunding, err := entity.NewCrowdfunding(creator.Address, input.DebitIssued, input.MaxInterestRate, input.ExpiresAt, input.MaturityAt, metadata.BlockTimestamp)
	if err != nil {
		return nil, fmt.Errorf("error creating crowdfunding: %w", err)
	}
	res, err := c.CrowdfundingRepository.CreateCrowdfunding(crowdfunding)
	if err != nil {
		return nil, fmt.Errorf("error creating crowdfunding: %w", err)
	}

	// Decrease creator's debt issuance limit
	creator.DebtIssuanceLimit.Sub(creator.DebtIssuanceLimit, input.DebitIssued)
	if _, err = c.UserRepository.UpdateUser(creator); err != nil {
		return nil, fmt.Errorf("error updating creator debt issuance limit: %w", err)
	}

	return &CreateCrowdfundingOutputDTO{
		Id:              res.Id,
		Creator:         res.Creator,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		State:           string(res.State),
		Orders:          res.Orders,
		ExpiresAt:       res.ExpiresAt,
		MaturityAt:      res.MaturityAt,
		CreatedAt:       res.CreatedAt,
	}, nil
}
