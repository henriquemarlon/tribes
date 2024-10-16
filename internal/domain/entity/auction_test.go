package entity

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

func TestNewAuction(t *testing.T) {
	creator := "matt"
	debt_issued := custom_type.BigInt{Int: big.NewInt(100)}
	interestRate := custom_type.BigInt{Int: big.NewInt(50)}
	createdAt := time.Now().Unix()
	expiresAt := createdAt + 3600

	auction, err := NewAuction(creator, debt_issued, interestRate, expiresAt, createdAt)
	assert.NoError(t, err)
	assert.NotNil(t, auction)
	assert.Equal(t, debt_issued, auction.DebtIssued)
	assert.Equal(t, interestRate, auction.MaxInterestRate)
	assert.Equal(t, AuctionOngoing, auction.State)
	assert.Equal(t, expiresAt, auction.ExpiresAt)
	assert.Equal(t, createdAt, auction.CreatedAt)
}

func TestNewAuction_Fail_InvalidAuction(t *testing.T) {
	creator := "matt"
	debt_issued := custom_type.BigInt{Int: big.NewInt(0)}
	interestRate := custom_type.BigInt{Int: big.NewInt(50)}
	createdAt := time.Now().Unix()
	expiresAt := createdAt + 3600

	auction, err := NewAuction(creator, debt_issued, interestRate, expiresAt, createdAt)
	assert.Error(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, ErrInvalidAuction, err)

	debt_issued = custom_type.BigInt{Int: big.NewInt(100)}
	interestRate = custom_type.BigInt{Int: big.NewInt(0)}

	auction, err = NewAuction(creator, debt_issued, interestRate, expiresAt, createdAt)
	assert.Error(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, ErrInvalidAuction, err)

	debt_issued = custom_type.BigInt{Int: big.NewInt(100)}
	interestRate = custom_type.BigInt{Int: big.NewInt(50)}
	expiresAt = createdAt

	auction, err = NewAuction(creator, debt_issued, interestRate, expiresAt, createdAt)
	assert.Error(t, err)
	assert.Nil(t, auction)
	assert.Equal(t, ErrInvalidAuction, err)
}
