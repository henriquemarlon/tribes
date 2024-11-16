package crowdfunding_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type FindCrowdfundingOutputDTO struct {
	Id              uint                            `json:"id"`
	Creator         common.Address                  `json:"creator"`
	DebtIssued      *uint256.Int                    `json:"debt_issued"`
	MaxInterestRate *uint256.Int                    `json:"max_interest_rate"`
	State           string                          `json:"state"`
	Orders          []*FindCrowdfundingOutputSubDTO `json:"orders"`
	ExpiresAt       int64                           `json:"expires_at"`
	CreatedAt       int64                           `json:"created_at"`
	UpdatedAt       int64                           `json:"updated_at"`
}

type FindCrowdfundingOutputSubDTO struct {
	Id             uint           `json:"id"`
	CrowdfundingId uint           `json:"crowdfunding_id"`
	Investor       common.Address `json:"investor"`
	Amount         *uint256.Int   `json:"amount"`
	InterestRate   *uint256.Int   `json:"interest_rate"`
	State          string         `json:"state"`
	CreatedAt      int64          `json:"created_at"`
	UpdatedAt      int64          `json:"updated_at"`
}
