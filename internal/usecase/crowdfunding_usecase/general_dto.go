package crowdfunding_usecase

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCrowdfundingOutputDTO struct {
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
