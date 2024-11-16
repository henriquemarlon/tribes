package entity

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrExpired              = errors.New("crowdfunding expired")
	ErrCrowdfundingNotFound = errors.New("crowdfunding not found")
	ErrInvalidCrowdfunding  = errors.New("invalid crowdfunding")
)

type CrowdfundingRepository interface {
	CreateCrowdfunding(crowdfunding *Crowdfunding) (*Crowdfunding, error)
	FindCrowdfundingsByCreator(creator common.Address) ([]*Crowdfunding, error)
	FindCrowdfundingById(id uint) (*Crowdfunding, error)
	FindAllCrowdfundings() ([]*Crowdfunding, error)
	CloseCrowdfunding(id uint) ([]*Order, error)
	SettleCrowdfunding(id uint) ([]*Order, error)
	UpdateCrowdfunding(crowdfunding *Crowdfunding) (*Crowdfunding, error)
	DeleteCrowdfunding(id uint) error
}

type CrowdfundingState string

const (
	CrowdfundingStateOngoing CrowdfundingState = "ongoing"
	CrowdfundingStateClosed  CrowdfundingState = "closed"
	CrowdfundingStateSettled CrowdfundingState = "settled"
)

type Crowdfunding struct {
	Id              uint              `json:"id" gorm:"primaryKey"`
	Creator         common.Address    `json:"creator,omitempty" gorm:"type:text;not null"`
	DebtIssued      uint256.Int       `json:"debt_issued,omitempty" gorm:"type:bigint;not null"`
	MaxInterestRate uint256.Int       `json:"max_interest_rate,omitempty" gorm:"type:bigint;not null"`
	State           CrowdfundingState `json:"state,omitempty" gorm:"type:text;not null"`
	Orders          []*Order          `json:"orders,omitempty" gorm:"foreignKey:CrowdfundingId;constraint:OnDelete:CASCADE"`
	ExpiresAt       int64             `json:"expires_at,omitempty" gorm:"not null"`
	CreatedAt       int64             `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt       int64             `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewCrowdfunding(creator common.Address, debt_issued uint256.Int, maxInterestRate uint256.Int, expiresAt int64, createdAt int64) (*Crowdfunding, error) {
	crowdfunding := &Crowdfunding{
		Creator:         creator,
		DebtIssued:      debt_issued,
		MaxInterestRate: maxInterestRate,
		State:           CrowdfundingStateOngoing,
		ExpiresAt:       expiresAt,
		CreatedAt:       createdAt,
	}
	if err := crowdfunding.Validate(); err != nil {
		return nil, err
	}
	return crowdfunding, nil
}

func (a *Crowdfunding) Validate() error {
	if a.Creator == (common.Address{}) || a.DebtIssued.Sign() == 0 || a.MaxInterestRate.Sign() == 0 || a.ExpiresAt == 0 || a.CreatedAt == 0 || a.CreatedAt >= a.ExpiresAt {
		return ErrInvalidCrowdfunding
	}
	return nil
}
