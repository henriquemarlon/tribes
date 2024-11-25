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
	closesAt := time.Now().Add(24 * time.Hour).Unix()
	maturityAt := time.Now().Add(48 * time.Hour).Unix() // Adiciona MaturityAt válido
	createdAt := time.Now().Unix()

	t.Run("Valid Crowdfunding", func(t *testing.T) {
		crowdfunding, err := NewCrowdfunding(creator, debtIssued, maxInterestRate, closesAt, maturityAt, createdAt)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.NoError(t, err, "Unexpected error for valid crowdfunding")
		assert.NotNil(t, crowdfunding, "Crowdfunding should not be nil for valid input")
		assert.Equal(t, creator, crowdfunding.Creator, "Creator mismatch")
		assert.Equal(t, debtIssued, crowdfunding.DebtIssued, "DebtIssued mismatch")
		assert.Equal(t, maxInterestRate, crowdfunding.MaxInterestRate, "MaxInterestRate mismatch")
		assert.Equal(t, CrowdfundingStateUnderReview, crowdfunding.State, "State mismatch")
		assert.Equal(t, closesAt, crowdfunding.ClosesAt, "ClosesAt mismatch")
		assert.Equal(t, maturityAt, crowdfunding.MaturityAt, "MaturityAt mismatch")
		assert.Equal(t, createdAt, crowdfunding.CreatedAt, "CreatedAt mismatch")
	})

	t.Run("Invalid Creator", func(t *testing.T) {
		_, err := NewCrowdfunding(common.Address{}, debtIssued, maxInterestRate, closesAt, maturityAt, createdAt)
		assert.Error(t, err, "Error expected for invalid creator address")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "Error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "invalid creator address", "Error message should mention invalid creator")
	})

	t.Run("Invalid Debt Issued", func(t *testing.T) {
		_, err := NewCrowdfunding(creator, uint256.NewInt(0), maxInterestRate, closesAt, maturityAt, createdAt)
		assert.Error(t, err, "Error expected for zero DebtIssued")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "Error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "debt issued cannot be zero", "Error message should mention DebtIssued")
	})

	t.Run("Debt Issued Exceeds Maximum", func(t *testing.T) {
		_, err := NewCrowdfunding(creator, uint256.NewInt(15000000001), maxInterestRate, closesAt, maturityAt, createdAt)
		assert.Error(t, err, "Error expected for exceeding max DebtIssued")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "Error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "debt issued exceeds the maximum allowed value", "Error message should mention maximum allowed value")
	})

	t.Run("Invalid Max Interest Rate", func(t *testing.T) {
		_, err := NewCrowdfunding(creator, debtIssued, uint256.NewInt(0), closesAt, maturityAt, createdAt)
		assert.Error(t, err, "Error expected for zero MaxInterestRate")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "Error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "max interest rate cannot be zero", "Error message should mention MaxInterestRate")
	})

	t.Run("Invalid Expiration Date", func(t *testing.T) {
		_, err := NewCrowdfunding(creator, debtIssued, maxInterestRate, 0, maturityAt, createdAt)
		assert.Error(t, err, "Error expected for missing Expiration Date")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "Error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "expiration date is missing", "Error message should mention Expiration Date")
	})

	t.Run("Invalid Creation Date >= Expiration Date", func(t *testing.T) {
		_, err := NewCrowdfunding(creator, debtIssued, maxInterestRate, createdAt, maturityAt, createdAt)
		assert.Error(t, err, "Error expected for CreatedAt >= Expiration Date")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "Error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "creation date cannot be greater than or equal to expiration date", "Error message should mention date issue")
	})
}

func TestCrowdfunding_Validate(t *testing.T) {
	creator := common.HexToAddress("0x123")
	debtIssued := uint256.NewInt(1000000000)
	maxInterestRate := uint256.NewInt(50000000)
	closesAt := time.Now().Add(24 * time.Hour).Unix()
	maturityAt := time.Now().Add(48 * time.Hour).Unix() // Adiciona MaturityAt válido
	createdAt := time.Now().Unix()

	t.Run("Valid Crowdfunding", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ClosesAt:        closesAt,
			MaturityAt:      maturityAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.NoError(t, err, "No error expected for valid crowdfunding")
	})

	t.Run("Invalid Creator Address", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         common.Address{},
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ClosesAt:        closesAt,
			MaturityAt:      maturityAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.Error(t, err, "Error expected for invalid creator address")
		assert.Contains(t, err.Error(), "invalid creator address", "Error message should mention invalid creator")
	})

	t.Run("Invalid Debt Issued Zero", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      uint256.NewInt(0),
			MaxInterestRate: maxInterestRate,
			ClosesAt:        closesAt,
			MaturityAt:      maturityAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.Error(t, err, "Error expected for zero DebtIssued")
		assert.Contains(t, err.Error(), "debt issued cannot be zero", "Error message should mention DebtIssued")
	})

	t.Run("Invalid Expiration Date", func(t *testing.T) {
		crowdfunding := &Crowdfunding{
			Creator:         creator,
			DebtIssued:      debtIssued,
			MaxInterestRate: maxInterestRate,
			ClosesAt:        0,
			MaturityAt:      maturityAt,
			CreatedAt:       createdAt,
		}
		err := crowdfunding.Validate()
		assert.Error(t, err, "Error expected for missing Expiration Date")
		assert.Contains(t, err.Error(), "expiration date is missing", "Error message should mention Expiration Date")
	})
}
