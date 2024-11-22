package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestOrder_Validate(t *testing.T) {
	crowdfundingId := uint(1)
	investor := common.HexToAddress("0x123")
	amount := uint256.NewInt(100000)
	interestRate := uint256.NewInt(500)
	createdAt := time.Now().Unix()

	t.Run("Invalid Crowdfunding ID", func(t *testing.T) {
		order := &Order{
			CrowdfundingId: 0,
			Investor:       investor,
			Amount:         amount,
			InterestRate:   interestRate,
			CreatedAt:      createdAt,
		}
		err := order.Validate()
		assert.NotNil(t, err, "An error should be returned for a zero crowdfunding ID")
		assert.True(t, errors.Is(err, ErrInvalidOrder), "The error should be ErrInvalidOrder")
		assert.Contains(t, err.Error(), "crowdfunding ID cannot be zero", "The error message should mention the invalid crowdfunding ID")
	})

	t.Run("Invalid Investor Address", func(t *testing.T) {
		order := &Order{
			CrowdfundingId: crowdfundingId,
			Investor:       common.Address{},
			Amount:         amount,
			InterestRate:   interestRate,
			CreatedAt:      createdAt,
		}
		err := order.Validate()
		assert.NotNil(t, err, "An error should be returned for an empty investor address")
		assert.True(t, errors.Is(err, ErrInvalidOrder), "The error should be ErrInvalidOrder")
		assert.Contains(t, err.Error(), "investor address cannot be empty", "The error message should mention the empty investor address")
	})

	t.Run("Invalid Amount - Zero", func(t *testing.T) {
		order := &Order{
			CrowdfundingId: crowdfundingId,
			Investor:       investor,
			Amount:         uint256.NewInt(0),
			InterestRate:   interestRate,
			CreatedAt:      createdAt,
		}
		err := order.Validate()
		assert.NotNil(t, err, "An error should be returned for a zero amount")
		assert.True(t, errors.Is(err, ErrInvalidOrder), "The error should be ErrInvalidOrder")
		assert.Contains(t, err.Error(), "amount cannot be zero or negative", "The error message should mention the invalid amount")
	})

	t.Run("Invalid Interest Rate - Zero", func(t *testing.T) {
		order := &Order{
			CrowdfundingId: crowdfundingId,
			Investor:       investor,
			Amount:         amount,
			InterestRate:   uint256.NewInt(0),
			CreatedAt:      createdAt,
		}
		err := order.Validate()
		assert.NotNil(t, err, "An error should be returned for a zero interest rate")
		assert.True(t, errors.Is(err, ErrInvalidOrder), "The error should be ErrInvalidOrder")
		assert.Contains(t, err.Error(), "interest rate cannot be zero or negative", "The error message should mention the invalid interest rate")
	})

	t.Run("Invalid Creation Date", func(t *testing.T) {
		order := &Order{
			CrowdfundingId: crowdfundingId,
			Investor:       investor,
			Amount:         amount,
			InterestRate:   interestRate,
			CreatedAt:      0,
		}
		err := order.Validate()
		assert.NotNil(t, err, "An error should be returned for a missing creation date")
		assert.True(t, errors.Is(err, ErrInvalidOrder), "The error should be ErrInvalidOrder")
		assert.Contains(t, err.Error(), "creation date is missing", "The error message should mention the missing creation date")
	})

	t.Run("Valid Order", func(t *testing.T) {
		order := &Order{
			CrowdfundingId: crowdfundingId,
			Investor:       investor,
			Amount:         amount,
			InterestRate:   interestRate,
			CreatedAt:      createdAt,
		}
		err := order.Validate()
		assert.Nil(t, err, "No error should be returned for a valid order")
	})
}
