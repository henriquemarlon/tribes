package entity

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrInvalidOrder  = errors.New("invalid order")
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository interface {
	CreateOrder(order *Order) (*Order, error)
	FindOrdersByState(crowdfundingId uint, state string) ([]*Order, error)
	FindOrderById(id uint) (*Order, error)
	FindOrdersByCrowdfundingId(id uint) ([]*Order, error)
	FindAllOrders() ([]*Order, error)
	UpdateOrder(order *Order) (*Order, error)
	DeleteOrder(id uint) error
}

type OrderState string

const (
	OrderStatePending  OrderState = "pending"
	OrderStateAccepted OrderState = "accepted"
	OrderStateExpired  OrderState = "partially_accepted"
	OrderStateRejected OrderState = "rejected"
	OrderStatePaid     OrderState = "paid"
)

type Order struct {
	Id             uint           `json:"id" gorm:"primaryKey"`
	CrowdfundingId uint           `json:"crowdfunding_id" gorm:"not null;index"`
	Investor       common.Address `json:"investor,omitempty" gorm:"not null"`
	Amount         uint256.Int    `json:"amount,omitempty" gorm:"type:bigint;not null"`
	InterestRate   uint256.Int    `json:"interest_rate,omitempty" gorm:"type:bigint;not null"`
	State          OrderState     `json:"state,omitempty" gorm:"type:text;not null"`
	CreatedAt      int64          `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt      int64          `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewOrder(crowdfundingId uint, investor common.Address, amount uint256.Int, interestRate uint256.Int, createdAt int64) (*Order, error) {
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
	if b.CrowdfundingId == 0 || b.Investor == (common.Address{}) || b.Amount.Sign() == 0 || b.InterestRate.Sign() == 0 {
		return ErrInvalidOrder
	}
	return nil
}
