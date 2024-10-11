package bid_usecase

import "github.com/Mugen-Builders/devolt/pkg/custom_type"

type FindBidOutputDTO struct {
	Id           uint                `json:"id"`
	AuctionId    uint                `json:"auction_id"`
	Bidder       custom_type.Address `json:"bidder"`
	Amount       custom_type.BigInt  `json:"amount"`
	InterestRate custom_type.BigInt  `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
	UpdatedAt    int64               `json:"updated_at"`
}
