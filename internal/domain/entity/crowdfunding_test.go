package entity

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestNewCrowdfunding(t *testing.T) {
	creator := common.HexToAddress("0x123")
	debt_issued := uint256.NewInt(100)
	interestRate := uint256.NewInt(50)
	createdAt := time.Now().Unix()
	expiresAt := createdAt + 3600

	crowdfunding, err := NewCrowdfunding(creator, *debt_issued, *interestRate, expiresAt, createdAt)
	assert.NoError(t, err)
	assert.NotNil(t, crowdfunding)
	assert.Equal(t, *debt_issued, crowdfunding.DebtIssued)
	assert.Equal(t, *interestRate, crowdfunding.MaxInterestRate)
	assert.Equal(t, CrowdfundingStateOngoing, crowdfunding.State)
	assert.Equal(t, expiresAt, crowdfunding.ExpiresAt)
	assert.Equal(t, createdAt, crowdfunding.CreatedAt)
}

func TestNewCrowdfunding_Fail_InvalidCrowdfunding(t *testing.T) {
	creator := common.Address{}
	debt_issued := uint256.NewInt(100)
	interestRate := uint256.NewInt(50)
	createdAt := time.Now().Unix()
	expiresAt := createdAt + 3600

	crowdfunding, err := NewCrowdfunding(creator, *debt_issued, *interestRate, expiresAt, createdAt)
	assert.Error(t, err)
	assert.Nil(t, crowdfunding)
	assert.Equal(t, ErrInvalidCrowdfunding, err)

	debt_issued = uint256.NewInt(100)
	interestRate = uint256.NewInt(0)

	crowdfunding, err = NewCrowdfunding(creator, *debt_issued, *interestRate, expiresAt, createdAt)
	assert.Error(t, err)
	assert.Nil(t, crowdfunding)
	assert.Equal(t, ErrInvalidCrowdfunding, err)

	debt_issued = uint256.NewInt(100)
	interestRate = uint256.NewInt(50)
	expiresAt = createdAt

	crowdfunding, err = NewCrowdfunding(creator, *debt_issued, *interestRate, expiresAt, createdAt)
	assert.Error(t, err)
	assert.Nil(t, crowdfunding)
	assert.Equal(t, ErrInvalidCrowdfunding, err)
}
