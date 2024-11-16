package crowdfunding_usecase

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateCrowdfundingInputDTO struct {
	DebtIssued      uint256.Int `json:"debt_issued"`
	MaxInterestRate uint256.Int `json:"max_interest_rate"`
	ExpiresAt       int64       `json:"expires_at"`
	CreatedAt       int64       `json:"created_at"`
}

type CreateCrowdfundingOutputDTO struct {
	Id              uint           `json:"id"`
	Creator         common.Address `json:"creator,omitempty"`
	DebtIssued      uint256.Int    `json:"debt_issued"`
	MaxInterestRate uint256.Int    `json:"max_interest_rate"`
	State           string         `json:"state"`
	ExpiresAt       int64          `json:"expires_at"`
	CreatedAt       int64          `json:"created_at"`
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

func (c *CreateCrowdfundingUseCase) Execute(input *CreateCrowdfundingInputDTO, metadata rollmelette.Metadata) (*CreateCrowdfundingOutputDTO, error) {
	creator, err := c.UserRepository.FindUserByAddress(metadata.MsgSender)
	if err != nil {
		return nil, err
	}
	crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(creator.Address)
	if err != nil {
		return nil, err
	}
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State == entity.CrowdfundingOngoing || crowdfunding.State == entity.CrowdfundingFinished {
			return nil, fmt.Errorf("creator already has an non paid crowdfunding")
		}
	}
	crowdfunding, err := entity.NewCrowdfunding(creator.Address, input.DebtIssued, input.MaxInterestRate, input.ExpiresAt, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.CrowdfundingRepository.CreateCrowdfunding(crowdfunding)
	if err != nil {
		return nil, err
	}
	return &CreateCrowdfundingOutputDTO{
		Id:              res.Id,
		Creator:         res.Creator,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		State:           string(res.State),
		ExpiresAt:       res.ExpiresAt,
		CreatedAt:       res.CreatedAt,
	}, nil
}
