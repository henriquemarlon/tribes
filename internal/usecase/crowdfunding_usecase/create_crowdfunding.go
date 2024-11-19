package crowdfunding_usecase

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateCrowdfundingInputDTO struct {
	MaxInterestRate *uint256.Int `json:"max_interest_rate"`
	ExpiresAt       int64        `json:"expires_at"`
}

type CreateCrowdfundingOutputDTO struct {
	Id              uint            `json:"id"`
	Creator         common.Address  `json:"creator,omitempty"`
	DebtIssued      *uint256.Int    `json:"debt_issued"`
	MaxInterestRate *uint256.Int    `json:"max_interest_rate"`
	State           string          `json:"state"`
	Orders          []*entity.Order `json:"orders"`
	ExpiresAt       int64           `json:"expires_at"`
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
	switch deposit := deposit.(type) {
	case *rollmelette.ERC20Deposit:
		creator, err := c.UserRepository.FindUserByAddress(metadata.MsgSender)
		if err != nil {
			return nil, fmt.Errorf("error finding creator: %w", err)
		}

		crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(creator.Address)
		if err != nil {
			return nil, fmt.Errorf("error finding crowdfunding campaigns: %w", err)
		}

		debtIssued := uint256.MustFromBig(deposit.Amount)

		// According with resolution 88 ( each company has a limit to be offered in a year )
		if creator.DebtIssuanceLimit.Cmp(debtIssued) == -1 {
			return nil, fmt.Errorf("creator debt issuance limit exceeded")
		}

		// According with resolution 88 ( the interval to create a crowdfunding campaign is 120 days )
		for _, crowdfunding := range crowdfundings {
			if crowdfunding.State != entity.CrowdfundingStateSettled && metadata.BlockTimestamp-crowdfunding.CreatedAt < 120*24*60*60 {
				return nil, fmt.Errorf("creator already has an active crowdfunding within the last 120 days")
			}
		}

		creator.DebtIssuanceLimit.Sub(creator.DebtIssuanceLimit, debtIssued)
		_, err = c.UserRepository.UpdateUser(creator)
		if err != nil {
			return nil, fmt.Errorf("error decreasing creator debt issuance limit: %w", err)
		}

		res, err := c.CrowdfundingRepository.CreateCrowdfunding(&entity.Crowdfunding{
			Creator:         creator.Address,
			DebtIssued:      debtIssued,
			MaxInterestRate: input.MaxInterestRate,
			State:           entity.CrowdfundingStateUnderReview,
			ExpiresAt:       input.ExpiresAt,
			CreatedAt:       metadata.BlockTimestamp,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating crowdfunding: %w", err)
		}

		return &CreateCrowdfundingOutputDTO{
			Id:              res.Id,
			Creator:         res.Creator,
			DebtIssued:      res.DebtIssued,
			MaxInterestRate: res.MaxInterestRate,
			State:           string(res.State),
			Orders:          res.Orders,
			ExpiresAt:       res.ExpiresAt,
			CreatedAt:       res.CreatedAt,
		}, nil
	default:
		return nil, fmt.Errorf("invalid deposit type: %T", deposit)
	}
}
