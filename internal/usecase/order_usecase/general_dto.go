package order_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type FindOrderOutputDTO struct {
	Id             uint           `json:"id"`
	CrowdfundingId uint           `json:"crowdfunding_id"`
	Investor       common.Address `json:"investor"`
	Amount         *uint256.Int   `json:"amount"`
	InterestRate   *uint256.Int   `json:"interest_rate"`
	State          string         `json:"state"`
	CreatedAt      int64          `json:"created_at"`
	UpdatedAt      int64          `json:"updated_at"`
}
