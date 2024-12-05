package crowdfunding_usecase

/*
#cgo LDFLAGS: -L./ -lverifier
#cgo CFLAGS: -I./include

#include <stdint.h>

int32_t add_numbers(int32_t a, int32_t b);
*/
import "C"
import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateCrowdfundingInputDTO struct {
	DebitIssued         *uint256.Int `json:"debt_issued"`
	MaxInterestRate     *uint256.Int `json:"max_interest_rate"`
	FundraisingDuration int64        `json:"fundraising_duration"`
	ClosesAt            int64        `json:"closes_at"`
	MaturityAt          int64        `json:"maturity_at"`
	Proof               string       `json:"proof"`
}

type CreateCrowdfundingOutputDTO struct {
	Id                  uint            `json:"id"`
	Token               common.Address  `json:"token,omitempty"`
	Amount              *uint256.Int    `json:"amount,omitempty"`
	Creator             common.Address  `json:"creator,omitempty"`
	DebtIssued          *uint256.Int    `json:"debt_issued"`
	MaxInterestRate     *uint256.Int    `json:"max_interest_rate"`
	Orders              []*entity.Order `json:"orders"`
	State               string          `json:"state"`
	FundraisingDuration int64           `json:"fundraising_duration"`
	ClosesAt            int64           `json:"closes_at"`
	MaturityAt          int64           `json:"maturity_at"`
	CreatedAt           int64           `json:"created_at"`
}

type CreateCrowdfundingUseCase struct {
	UserRepository          entity.UserRepository
	ContractRepository      entity.ContractRepository
	SocialAccountRepository entity.SocialAccountRepository
	CrowdfundingRepository  entity.CrowdfundingRepository
}

func NewCreateCrowdfundingUseCase(userRepository entity.UserRepository, contractRepository entity.ContractRepository, socialAccountRepository entity.SocialAccountRepository, crowdfundingRepository entity.CrowdfundingRepository) *CreateCrowdfundingUseCase {
	return &CreateCrowdfundingUseCase{
		UserRepository:          userRepository,
		ContractRepository:      contractRepository,
		SocialAccountRepository: socialAccountRepository,
		CrowdfundingRepository:  crowdfundingRepository,
	}
}

func (c *CreateCrowdfundingUseCase) Execute(ctx context.Context, input *CreateCrowdfundingInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateCrowdfundingOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type: %T", deposit)
	}

	creator, err := c.UserRepository.FindUserByAddress(ctx, erc20Deposit.Sender)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	if _, err = c.ContractRepository.FindContractByAddress(ctx, erc20Deposit.Token); err != nil {
		return nil, fmt.Errorf("token unknown, cannot create crowdfunding: %w", err)
	}

	// TODO: Replace the logic bellow to validate a notary signature and return the verified values to create a social account.
	a := C.int32_t(3)
	b := C.int32_t(4)
	result := C.add_numbers(a, b)
	slog.Info("TLSN verifier result", "result", result)
	mock, err := entity.NewSocialAccount(creator.Id, "vitalik", 1000, "twitter", metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	_, err = c.SocialAccountRepository.CreateSocialAccount(ctx, mock)
	if err != nil {
		return nil, err
	}

	// Validate debt issuance limit
	if creator.DebtIssuanceLimit.Cmp(input.DebitIssued) < 0 {
		return nil, fmt.Errorf("creator debt issuance limit exceeded")
	}

	crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, creator.Address)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaigns: %w", err)
	}

	// Check for active crowdfunding campaigns within the last 120 days
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State != entity.CrowdfundingStateSettled && metadata.BlockTimestamp-crowdfunding.CreatedAt < 120*24*60*60 {
			return nil, fmt.Errorf("creator already has an active crowdfunding within the last 120 days")
		}
	}

	crowdfunding, err := entity.NewCrowdfunding(erc20Deposit.Token, uint256.MustFromBig(erc20Deposit.Amount), creator.Address, input.DebitIssued, input.MaxInterestRate, input.FundraisingDuration, input.ClosesAt, input.MaturityAt, metadata.BlockTimestamp)
	if err != nil {
		return nil, fmt.Errorf("error creating crowdfunding: %w", err)
	}
	res, err := c.CrowdfundingRepository.CreateCrowdfunding(ctx, crowdfunding)
	if err != nil {
		return nil, fmt.Errorf("error creating crowdfunding: %w", err)
	}

	// Decrease creator's debt issuance limit
	creator.DebtIssuanceLimit.Sub(creator.DebtIssuanceLimit, input.DebitIssued)
	if _, err = c.UserRepository.UpdateUser(ctx, creator); err != nil {
		return nil, fmt.Errorf("error updating creator debt issuance limit: %w", err)
	}

	return &CreateCrowdfundingOutputDTO{
		Id:                  res.Id,
		Token:               res.Token,
		Amount:              res.Amount,
		Creator:             res.Creator,
		DebtIssued:          res.DebtIssued,
		MaxInterestRate:     res.MaxInterestRate,
		Orders:              res.Orders,
		State:               string(res.State),
		FundraisingDuration: res.FundraisingDuration,
		ClosesAt:            res.ClosesAt,
		MaturityAt:          res.MaturityAt,
		CreatedAt:           res.CreatedAt,
	}, nil
}
