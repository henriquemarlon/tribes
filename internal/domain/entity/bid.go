package entity

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

var (
	ErrInvalidBid  = errors.New("invalid bid")
	ErrBidNotFound = errors.New("bid not found")
)

type BidRepository interface {
	CreateBid(bid *Bid) (*Bid, error)
	FindBidsByState(auctionId uint, state string) ([]*Bid, error)
	FindBidById(id uint) (*Bid, error)
	FindBidsByAuctionId(id uint) ([]*Bid, error)
	FindAllBids() ([]*Bid, error)
	UpdateBid(bid *Bid) (*Bid, error)
	DeleteBid(id uint) error
}

type BidState string

const (
	BidStatePending  BidState = "pending"
	BidStateAccepted BidState = "accepted"
	BidStateExpired  BidState = "partially_accepted"
	BidStateRejected BidState = "rejected"
	BidStatePaid     BidState = "paid"
)

type Bid struct {
	Id           uint                `json:"id" gorm:"primaryKey"`
	AuctionId    uint                `json:"auction_id" gorm:"not null;index"`
	Bidder       custom_type.Address `json:"bidder,omitempty" gorm:"not null"`
	Amount       custom_type.BigInt  `json:"amount,omitempty" gorm:"type:bigint;not null"`
	InterestRate custom_type.BigInt  `json:"interest_rate,omitempty" gorm:"type:bigint;not null"`
	State        BidState            `json:"state,omitempty" gorm:"type:text;not null"`
	CreatedAt    int64               `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt    int64               `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewBid(auctionId uint, bidder custom_type.Address, amount custom_type.BigInt, interestRate custom_type.BigInt, createdAt int64) (*Bid, error) {
	bid := &Bid{
		AuctionId:    auctionId,
		Bidder:       bidder,
		Amount:       amount,
		InterestRate: interestRate,
		State:        BidStatePending,
		CreatedAt:    createdAt,
	}
	if err := bid.Validate(); err != nil {
		return nil, err
	}
	return bid, nil
}

func (b *Bid) Validate() error {
	if b.AuctionId == 0 || b.Bidder.Address == (common.Address{}) || b.Amount.Int.Sign() == 0 || b.InterestRate.Int.Sign() == 0 {
		return ErrInvalidBid
	}
	return nil
}
