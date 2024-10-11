package auction_usecase

import (
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindAuctionOutputDTO struct {
	Id           uint                       `json:"id"`
	Creator      custom_type.Address        `json:"creator"`
	DebtIssued   custom_type.BigInt         `json:"debt_issued"`
	InterestRate custom_type.BigInt         `json:"interest_rate"`
	State        string                     `json:"state"`
	Bids         []*FindAuctionOutputSubDTO `json:"bids"`
	ExpiresAt    int64                      `json:"expires_at"`
	CreatedAt    int64                      `json:"created_at"`
	UpdatedAt    int64                      `json:"updated_at"`
}

type FindAuctionOutputSubDTO struct {
	Id           uint                `json:"id"`
	AuctionId    uint                `json:"auction_id"`
	Bidder       custom_type.Address `json:"bidder"`
	Amount       custom_type.BigInt  `json:"amount"`
	InterestRate custom_type.BigInt  `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
	UpdatedAt    int64               `json:"updated_at"`
}
