package entity

import (
	"errors"

	"github.com/tribeshq/tribes/pkg/custom_type"
)

var (
	ErrExpired         = errors.New("auction expired")
	ErrAuctionNotFound = errors.New("auction not found")
	ErrInvalidAuction  = errors.New("invalid auction")
)

type AuctionRepository interface {
	DeleteAuction(id uint) error
	FindAuctionsByCreator(creator string) ([]*Auction, error)
	FindAuctionByStateFromCreator(creator string, state string) ([]*Auction, error)
	FindAllAuctions() ([]*Auction, error)
	FindAuctionsByState(state string) ([]*Auction, error)
	FindAuctionById(id uint) (*Auction, error)
	CreateAuction(auction *Auction) (*Auction, error)
	UpdateAuction(auction *Auction) (*Auction, error)
}

type AuctionState string

const (
	AuctionOngoing   AuctionState = "ongoing"
	AuctionFinished  AuctionState = "finished"
	AuctionCancelled AuctionState = "cancelled"
	AuctionPaid      AuctionState = "paid"
)

type Auction struct {
	Id              uint               `json:"id" gorm:"primaryKey"`
	Creator         string             `json:"creator,omitempty" gorm:"type:text;not null"`
	DebtIssued      custom_type.BigInt `json:"debt_issued,omitempty" gorm:"type:bigint;not null"`
	MaxInterestRate custom_type.BigInt `json:"max_interest_rate,omitempty" gorm:"type:bigint;not null"`
	State           AuctionState       `json:"state,omitempty" gorm:"type:text;not null"`
	Bids            []*Bid             `json:"bids,omitempty" gorm:"foreignKey:AuctionId;constraint:OnDelete:CASCADE"`
	ExpiresAt       int64              `json:"expires_at,omitempty" gorm:"not null"`
	CreatedAt       int64              `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt       int64              `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewAuction(creator string, debt_issued custom_type.BigInt, maxInterestRate custom_type.BigInt, expiresAt int64, createdAt int64) (*Auction, error) {
	auction := &Auction{
		Creator:         creator,
		DebtIssued:      debt_issued,
		MaxInterestRate: maxInterestRate,
		State:           AuctionOngoing,
		ExpiresAt:       expiresAt,
		CreatedAt:       createdAt,
	}
	if err := auction.Validate(); err != nil {
		return nil, err
	}
	return auction, nil
}

func (a *Auction) Validate() error {
	if a.Creator == "" || a.DebtIssued.Int.Sign() == 0 || a.MaxInterestRate.Int.Sign() == 0 || a.ExpiresAt == 0 || a.CreatedAt == 0 || a.CreatedAt >= a.ExpiresAt {
		return ErrInvalidAuction
	}
	return nil
}
