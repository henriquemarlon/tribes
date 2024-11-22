package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestNewCrowdfunding(t *testing.T) {
	creator := common.HexToAddress("0x123")
	debtIssued := uint256.NewInt(1000000000)
	maxInterestRate := uint256.NewInt(50000000)
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	createdAt := time.Now().Unix()

	t.Run("Valid Crowdfunding", func(t *testing.T) {
		crowdfunding, err := NewCrowdfunding(creator, debtIssued, maxInterestRate, expiresAt, createdAt)
		assert.Nil(t, err, "Unexpected error when creating a valid crowdfunding")
		assert.NotNil(t, crowdfunding, "The created crowdfunding should not be nil")
		assert.Equal(t, creator, crowdfunding.Creator, "Creator address is incorrect")
		assert.Equal(t, debtIssued, crowdfunding.DebtIssued, "Debt issued is incorrect")
		assert.Equal(t, maxInterestRate, crowdfunding.MaxInterestRate, "Max interest rate is incorrect")
		assert.Equal(t, CrowdfundingStateOngoing, crowdfunding.State, "Initial state should be 'ongoing'")
		assert.Equal(t, expiresAt, crowdfunding.ExpiresAt, "Expiration date is incorrect")
		assert.Equal(t, createdAt, crowdfunding.CreatedAt, "Creation date is incorrect")
	})

	t.Run("Invalid Crowdfunding - Validation Error", func(t *testing.T) {
		_, err := NewCrowdfunding(common.Address{}, debtIssued, maxInterestRate, expiresAt, createdAt)
		assert.NotNil(t, err, "An error should be returned for an invalid creator address")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
	})
}

func TestCrowdfunding_Validate(t *testing.T) {
	creator := common.HexToAddress("0x123")
	debtIssued := uint256.NewInt(1000000000)
	maxInterestRate := uint256.NewInt(50000000)
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	createdAt := time.Now().Unix()

	t.Run("Invalid Creator Address", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         common.Address{},
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ExpiresAt:       expiresAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NotNil(t, err, "An error should be returned for an invalid creator address")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "invalid creator address", "The error message should mention invalid address")
	})

	t.Run("Invalid Debt Issued - Zero", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      uint256.NewInt(0),
			MaxInterestRate: maxInterestRate,
			ExpiresAt:       expiresAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NotNil(t, err, "An error should be returned for zero debt issued")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "debt issued cannot be zero", "The error message should mention invalid debt issued")
	})

	t.Run("Invalid Debt Issued - Exceeds Maximum", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      uint256.NewInt(15000000001),
			MaxInterestRate: maxInterestRate,
			ExpiresAt:       expiresAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NotNil(t, err, "An error should be returned for debt issued exceeding the maximum allowed value")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "debt issued exceeds the maximum allowed value", "The error message should mention exceeded limit")
	})

	t.Run("Invalid Max Interest Rate", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      debtIssued,
			MaxInterestRate: uint256.NewInt(0),
			ExpiresAt:       expiresAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NotNil(t, err, "An error should be returned for zero max interest rate")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "max interest rate cannot be zero", "The error message should mention invalid max interest rate")
	})

	t.Run("Invalid Expiration Date", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ExpiresAt:       0,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NotNil(t, err, "An error should be returned for missing expiration date")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "expiration date is missing", "The error message should mention missing expiration date")
	})

	t.Run("Invalid Creation Date - Greater or Equal to Expiration Date", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ExpiresAt:       createdAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NotNil(t, err, "An error should be returned for creation date greater than or equal to expiration date")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "creation date cannot be greater than or equal to expiration date", "The error message should mention invalid creation date")
	})

	t.Run("Valid Crowdfunding", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ExpiresAt:       expiresAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.Nil(t, err, "No error should be returned for a valid crowdfunding")
	})
}
