package entity

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder_Success(t *testing.T) {
	crowdfundingId := uint(1)
	investor := common.HexToAddress("0x123")
	amount := uint256.NewInt(100)
	interestRate := uint256.NewInt(50)
	createdAt := time.Now().Unix()

	order, err := NewOrder(crowdfundingId, investor, amount, interestRate, createdAt)
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, crowdfundingId, order.CrowdfundingId)
	assert.Equal(t, investor, order.Investor)
	assert.Equal(t, *amount, order.Amount)
	assert.Equal(t, *interestRate, order.InterestRate)
	assert.Equal(t, OrderStatePending, order.State)
	assert.Equal(t, createdAt, order.CreatedAt)
}

func TestNewOrder_Fail_InvalidOrder(t *testing.T) {
	crowdfundingId := uint(0)
	investor := common.HexToAddress("0x123")
	amount := uint256.NewInt(100)
	interestRate := uint256.NewInt(50)
	createdAt := time.Now().Unix()

	order, err := NewOrder(crowdfundingId, investor, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, ErrInvalidOrder, err)

	crowdfundingId = uint(1)
	investor = common.Address{}

	order, err = NewOrder(crowdfundingId, investor, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, ErrInvalidOrder, err)

	investor = common.HexToAddress("0x123")
	amount = uint256.NewInt(0)

	order, err = NewOrder(crowdfundingId, investor, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, ErrInvalidOrder, err)

	amount = uint256.NewInt(100)
	interestRate = uint256.NewInt(0)

	order, err = NewOrder(crowdfundingId, investor, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, ErrInvalidOrder, err)
}
