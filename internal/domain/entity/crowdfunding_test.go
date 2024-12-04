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
	debtIssued := uint256.NewInt(10000000)
	maxInterestRate := uint256.NewInt(50000000)
	fundraisingDuration := int64(604800)
	createdAt := time.Now().Unix()
	closesAt := createdAt + 24*3600
	maturityAt := createdAt + 48*3600

	tests := []struct {
		name             string
		creator          common.Address
		debtIssued       *uint256.Int
		maxInterestRate  *uint256.Int
		fundraisingDur   int64
		closesAt         int64
		maturityAt       int64
		createdAt        int64
		expectedError    error
		errorMessagePart string
	}{
		{
			name:            "Valid Crowdfunding",
			creator:         creator,
			debtIssued:      debtIssued,
			maxInterestRate: maxInterestRate,
			fundraisingDur:  fundraisingDuration,
			closesAt:        closesAt,
			maturityAt:      maturityAt,
			createdAt:       createdAt,
			expectedError:   nil,
		},
		{
			name:             "Invalid Creator",
			creator:          common.Address{},
			debtIssued:       debtIssued,
			maxInterestRate:  maxInterestRate,
			fundraisingDur:   fundraisingDuration,
			closesAt:         closesAt,
			maturityAt:       maturityAt,
			createdAt:        createdAt,
			expectedError:    ErrInvalidCrowdfunding,
			errorMessagePart: "invalid creator address",
		},
		{
			name:             "Invalid Debt Issued",
			creator:          creator,
			debtIssued:       uint256.NewInt(0),
			maxInterestRate:  maxInterestRate,
			fundraisingDur:   fundraisingDuration,
			closesAt:         closesAt,
			maturityAt:       maturityAt,
			createdAt:        createdAt,
			expectedError:    ErrInvalidCrowdfunding,
			errorMessagePart: "debt issued cannot be zero",
		},
		{
			name:             "Debt Issued Exceeds Maximum",
			creator:          creator,
			debtIssued:       uint256.NewInt(15000000001),
			maxInterestRate:  maxInterestRate,
			fundraisingDur:   fundraisingDuration,
			closesAt:         closesAt,
			maturityAt:       maturityAt,
			createdAt:        createdAt,
			expectedError:    ErrInvalidCrowdfunding,
			errorMessagePart: "debt issued exceeds the maximum allowed value",
		},
		{
			name:             "Invalid Close Date",
			creator:          creator,
			debtIssued:       debtIssued,
			maxInterestRate:  maxInterestRate,
			fundraisingDur:   fundraisingDuration,
			closesAt:         0,
			maturityAt:       maturityAt,
			createdAt:        createdAt,
			expectedError:    ErrInvalidCrowdfunding,
			errorMessagePart: "close date is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crowdfunding, err := NewCrowdfunding(tt.creator, tt.debtIssued, tt.maxInterestRate, tt.fundraisingDur, tt.closesAt, tt.maturityAt, tt.createdAt)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Contains(t, err.Error(), tt.errorMessagePart)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, crowdfunding)
			}
		})
	}
}

func TestCrowdfunding_Validate(t *testing.T) {
	creator := common.HexToAddress("0x123")
	debtIssued := uint256.NewInt(10000000)
	maxInterestRate := uint256.NewInt(50000000)
	createdAt := time.Now().Unix()
	closesAt := createdAt + 24*3600
	maturityAt := createdAt + 48*3600

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
		assert.NoError(t, err)
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
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid creator address")
	})
}
