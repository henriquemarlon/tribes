package entity

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

func TestNewBid_Success(t *testing.T) {
	auctionId := uint(1)
	bidder := custom_type.Address{Address: common.HexToAddress("0x123")}
	amount := custom_type.BigInt{Int: big.NewInt(100)}
	interestRate := custom_type.BigInt{Int: big.NewInt(50)}
	createdAt := time.Now().Unix()

	bid, err := NewBid(auctionId, bidder, amount, interestRate, createdAt)
	assert.NoError(t, err)
	assert.NotNil(t, bid)
	assert.Equal(t, auctionId, bid.AuctionId)
	assert.Equal(t, bidder, bid.Bidder)
	assert.Equal(t, amount, bid.Amount)
	assert.Equal(t, interestRate, bid.InterestRate)
	assert.Equal(t, BidStatePending, bid.State)
	assert.Equal(t, createdAt, bid.CreatedAt)
}

func TestNewBid_Fail_InvalidBid(t *testing.T) {
	auctionId := uint(0)
	bidder := custom_type.Address{Address: common.HexToAddress("0x123")}
	amount := custom_type.BigInt{Int: big.NewInt(100)}
	interestRate := custom_type.BigInt{Int: big.NewInt(50)}
	createdAt := time.Now().Unix()

	bid, err := NewBid(auctionId, bidder, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, bid)
	assert.Equal(t, ErrInvalidBid, err)

	auctionId = uint(1)
	bidder = custom_type.Address{Address: common.Address{}}

	bid, err = NewBid(auctionId, bidder, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, bid)
	assert.Equal(t, ErrInvalidBid, err)

	bidder = custom_type.Address{Address: common.HexToAddress("0x123")}
	amount = custom_type.BigInt{Int: big.NewInt(0)}

	bid, err = NewBid(auctionId, bidder, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, bid)
	assert.Equal(t, ErrInvalidBid, err)

	amount = custom_type.BigInt{Int: big.NewInt(100)}
	interestRate = custom_type.BigInt{Int: big.NewInt(0)}

	bid, err = NewBid(auctionId, bidder, amount, interestRate, createdAt)
	assert.Error(t, err)
	assert.Nil(t, bid)
	assert.Equal(t, ErrInvalidBid, err)
}
