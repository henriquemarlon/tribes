package crowdfunding_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindCrowdfundingOutputDTO struct {
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
