package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrCrowdfundingNotFound = errors.New("crowdfunding not found")
	ErrInvalidCrowdfunding  = errors.New("invalid crowdfunding")
)

type CrowdfundingState string

const (
	CrowdfundingStateUnderReview CrowdfundingState = "under_review"
	CrowdfundingStateClosed      CrowdfundingState = "closed"
	CrowdfundingStateOngoing     CrowdfundingState = "ongoing"
	CrowdfundingStateCanceled    CrowdfundingState = "canceled"
	CrowdfundingStateSettled     CrowdfundingState = "settled"
)

type CrowdfundingRepository interface {
	CreateCrowdfunding(ctx context.Context, crowdfunding *Crowdfunding) (*Crowdfunding, error)
	FindCrowdfundingsByCreator(ctx context.Context, creator common.Address) ([]*Crowdfunding, error)
	FindCrowdfundingsByInvestor(ctx context.Context, investor common.Address) ([]*Crowdfunding, error)
	FindCrowdfundingById(ctx context.Context, id uint) (*Crowdfunding, error)
	FindAllCrowdfundings(ctx context.Context) ([]*Crowdfunding, error)
	UpdateCrowdfunding(ctx context.Context, crowdfunding *Crowdfunding) (*Crowdfunding, error)
	DeleteCrowdfunding(ctx context.Context, id uint) error
}

type Crowdfunding struct {
	Id                  uint              `json:"id" gorm:"primaryKey"`
	Token               common.Address    `json:"token,omitempty" gorm:"type:text;not null"`
	Amount              *uint256.Int      `json:"amount,omitempty" gorm:"type:text;not null"`
	Creator             common.Address    `json:"creator,omitempty" gorm:"type:text;not null"`
	DebtIssued          *uint256.Int      `json:"debt_issued,omitempty" gorm:"type:text;not null"`
	MaxInterestRate     *uint256.Int      `json:"max_interest_rate,omitempty" gorm:"type:text;not null"`
	TotalObligation     *uint256.Int      `json:"total_obligation,omitempty" gorm:"type:text;not null;default:0"`
	State               CrowdfundingState `json:"state,omitempty" gorm:"type:text;not null"`
	Orders              []*Order          `json:"orders,omitempty" gorm:"foreignKey:CrowdfundingId;constraint:OnDelete:CASCADE"`
	FundraisingDuration int64             `json:"fundraising_duration,omitempty" gorm:"not null"`
	ClosesAt            int64             `json:"closes_at,omitempty" gorm:"not null"`
	MaturityAt          int64             `json:"maturity_at,omitempty" gorm:"not null"`
	CreatedAt           int64             `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt           int64             `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewCrowdfunding(token common.Address, amount *uint256.Int, creator common.Address, debt_issued *uint256.Int, maxInterestRate *uint256.Int, fundraisingDuration int64, closesAt int64, maturityAt int64, createdAt int64) (*Crowdfunding, error) {
	crowdfunding := &Crowdfunding{
		Token:               token,
		Amount:              amount,
		Creator:             creator,
		DebtIssued:          debt_issued,
		MaxInterestRate:     maxInterestRate,
		State:               CrowdfundingStateUnderReview,
		FundraisingDuration: fundraisingDuration,
		ClosesAt:            closesAt,
		MaturityAt:          maturityAt,
		CreatedAt:           createdAt,
	}
	if err := crowdfunding.Validate(); err != nil {
		return nil, err
	}
	return crowdfunding, nil
}

func (a *Crowdfunding) Validate() error {
	if a.Token == (common.Address{}) {
		return fmt.Errorf("%w: invalid token address", ErrInvalidCrowdfunding)
	}
	if a.Amount.Sign() == 0 {
		return fmt.Errorf("%w: amount cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.Creator == (common.Address{}) {
		return fmt.Errorf("%w: invalid creator address", ErrInvalidCrowdfunding)
	}
	if a.DebtIssued.Sign() == 0 {
		return fmt.Errorf("%w: debt issued cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.MaxInterestRate.Sign() == 0 {
		return fmt.Errorf("%w: max interest rate cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidCrowdfunding)
	}
	if a.ClosesAt == 0 {
		return fmt.Errorf("%w: close date is missing", ErrInvalidCrowdfunding)
	}
	if a.MaturityAt == 0 {
		return fmt.Errorf("%w: maturity date is missing", ErrInvalidCrowdfunding)
	}
	return nil
}
