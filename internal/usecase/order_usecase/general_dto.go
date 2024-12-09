package order_usecase

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindOrderOutputDTO struct {
	Id             uint                `json:"id"`
	CrowdfundingId uint                `json:"crowdfunding_id"`
	Investor       custom_type.Address `json:"investor"`
	Amount         *uint256.Int        `json:"amount"`
	InterestRate   *uint256.Int        `json:"interest_rate"`
	State          string              `json:"state"`
	CreatedAt      int64               `json:"created_at"`
	UpdatedAt      int64               `json:"updated_at"`
}
