package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

var (
	ErrInvalidOrder  = errors.New("invalid order")
	ErrOrderNotFound = errors.New("order not found")
)

type OrderState string

const (
	OrderStatePending           OrderState = "pending"
	OrderStateAccepted          OrderState = "accepted"
	OrderCancelled              OrderState = "cancelled"
	OrderStatePartiallyAccepted OrderState = "partially_accepted"
	OrderStateRejected          OrderState = "rejected"
	OrderStateSettled           OrderState = "settled"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	FindOrderById(ctx context.Context, id uint) (*Order, error)
	FindOrdersByCrowdfundingId(ctx context.Context, id uint) ([]*Order, error)
	FindOrdersByState(ctx context.Context, crowdfundingId uint, state string) ([]*Order, error)
	FindOrdersByInvestor(ctx context.Context, investor custom_type.Address) ([]*Order, error)
	FindAllOrders(ctx context.Context) ([]*Order, error)
	UpdateOrder(ctx context.Context, order *Order) (*Order, error)
	DeleteOrder(ctx context.Context, id uint) error
}

type Order struct {
	Id             uint                `json:"id" gorm:"primaryKey"`
	CrowdfundingId uint                `json:"crowdfunding_id" gorm:"not null;index"`
	Investor       custom_type.Address `json:"investor,omitempty" gorm:"not null"`
	Amount         *uint256.Int        `json:"amount,omitempty" gorm:"type:text;not null"`
	InterestRate   *uint256.Int        `json:"interest_rate,omitempty" gorm:"type:text;not null"`
	State          OrderState          `json:"state,omitempty" gorm:"type:text;not null"`
	CreatedAt      int64               `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt      int64               `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewOrder(crowdfundingId uint, investor custom_type.Address, amount *uint256.Int, interestRate *uint256.Int, createdAt int64) (*Order, error) {
	order := &Order{
		CrowdfundingId: crowdfundingId,
		Investor:       investor,
		Amount:         amount,
		InterestRate:   interestRate,
		State:          OrderStatePending,
		CreatedAt:      createdAt,
	}
	if err := order.Validate(); err != nil {
		return nil, err
	}
	return order, nil
}

func (b *Order) Validate() error {
	if b.CrowdfundingId == 0 {
		return fmt.Errorf("%w: crowdfunding ID cannot be zero", ErrInvalidOrder)
	}
	if b.Investor == (custom_type.Address{}) {
		return fmt.Errorf("%w: investor address cannot be empty", ErrInvalidOrder)
	}
	if b.Amount.Sign() <= 0 {
		return fmt.Errorf("%w: amount cannot be zero or negative", ErrInvalidOrder)
	}
	if b.InterestRate.Sign() <= 0 {
		return fmt.Errorf("%w: interest rate cannot be zero or negative", ErrInvalidOrder)
	}
	if b.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidOrder)
	}
	return nil
}
