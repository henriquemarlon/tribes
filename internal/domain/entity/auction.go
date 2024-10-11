package entity

import (
	"errors"

	"github.com/Mugen-Builders/devolt/pkg/custom_type"
)

var (
	ErrExpired         = errors.New("auction expired")
	ErrAuctionNotFound = errors.New("auction not found")
	ErrInvalidAuction  = errors.New("invalid auction")
)

type AuctionRepository interface {
	DeleteAuction(id uint) error
	FindActiveAuction() (*Auction, error)
	FindAllAuctions() ([]*Auction, error)
	FindAuctionById(id uint) (*Auction, error)
	CreateAuction(auction *Auction) (*Auction, error)
	UpdateAuction(auction *Auction) (*Auction, error)
}

type AuctionState string

const (
	AuctionOngoing   AuctionState = "ongoing"
	AuctionFinished  AuctionState = "finished"
	AuctionCancelled AuctionState = "cancelled"
)

type Auction struct {
	Id           uint               `json:"id" gorm:"primaryKey"`
	Creator      custom_type.Address `json:"creator,omitempty" gorm:"type:text;not null"`
	DebtIssued   custom_type.BigInt `json:"debt_issued,omitempty" gorm:"type:bigint;not null"`
	InterestRate custom_type.BigInt `json:"interest_rate,omitempty" gorm:"type:bigint;not null"`
	State        AuctionState       `json:"state,omitempty" gorm:"type:text;not null"`
	Bids         []*Bid             `json:"bids,omitempty" gorm:"foreignKey:AuctionId;constraint:OnDelete:CASCADE"`
	ExpiresAt    int64              `json:"expires_at,omitempty" gorm:"not null"`
	CreatedAt    int64              `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt    int64              `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewAuction(creator custom_type.Address, debt_issued custom_type.BigInt, interestRate custom_type.BigInt, expiresAt int64, createdAt int64) (*Auction, error) {
	auction := &Auction{
		Creator:      creator,
		DebtIssued:   debt_issued,
		InterestRate: interestRate,
		State:        AuctionOngoing,
		ExpiresAt:    expiresAt,
		CreatedAt:    createdAt,
	}
	if err := auction.Validate(); err != nil {
		return nil, err
	}
	return auction, nil
}

func (a *Auction) Validate() error {
	if a.DebtIssued.Int.Sign() == 0 || a.InterestRate.Int.Sign() == 0 || a.ExpiresAt == 0 || a.CreatedAt == 0 || a.CreatedAt >= a.ExpiresAt {
		return ErrInvalidAuction
	}
	return nil
}
