package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewContract(t *testing.T) {
	symbol := "ETH"
	address := common.HexToAddress("0x123")
	createdAt := time.Now().Unix()

	contract, err := NewContract(symbol, address, createdAt)

	assert.Nil(t, err, "Unexpected error when creating a valid contract")
	assert.NotNil(t, contract, "The created contract should not be nil")
	assert.Equal(t, symbol, contract.Symbol, "The contract symbol is incorrect")
	assert.Equal(t, address, contract.Address, "The contract address is incorrect")
	assert.Equal(t, createdAt, contract.CreatedAt, "The contract creation date is incorrect")
}

func TestContract_Validate(t *testing.T) {
	t.Run("Invalid symbol", func(t *testing.T) {
		contract := &Contract{
			Symbol:    "",
			Address:   common.HexToAddress("0x123"),
			CreatedAt: time.Now().Unix(),
		}
		err := contract.Validate()
		assert.NotNil(t, err, "An error should be returned for an invalid symbol")
		assert.True(t, errors.Is(err, ErrInvalidContract), "The returned error should be ErrInvalidContract")
		assert.Contains(t, err.Error(), "symbol cannot be empty", "The error message should mention that the symbol is empty")
	})

	t.Run("Invalid address", func(t *testing.T) {
		contract := &Contract{
			Symbol:    "ETH",
			Address:   common.Address{},
			CreatedAt: time.Now().Unix(),
		}
		err := contract.Validate()
		assert.NotNil(t, err, "An error should be returned for an invalid address")
		assert.True(t, errors.Is(err, ErrInvalidContract), "The returned error should be ErrInvalidContract")
		assert.Contains(t, err.Error(), "address cannot be empty", "The error message should mention that the address is empty")
	})

	t.Run("Valid contract", func(t *testing.T) {
		contract := &Contract{
			Symbol:    "ETH",
			Address:   common.HexToAddress("0x123"),
			CreatedAt: time.Now().Unix(),
		}
		err := contract.Validate()
		assert.Nil(t, err, "No error should be returned for a valid contract")
	})
}
