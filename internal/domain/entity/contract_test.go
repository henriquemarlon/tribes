package entity

import (
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
	assert.Nil(t, err)
	assert.NotNil(t, contract)
	assert.Equal(t, symbol, contract.Symbol)
	assert.Equal(t, address, contract.Address)
	assert.NotZero(t, contract.CreatedAt)
}

func TestContract_Validate(t *testing.T) {
	createdAt := time.Now().Unix()

	// Invalid symbol
	contract := &Contract{
		Symbol:    "",
		Address:   common.HexToAddress("0x123"),
		CreatedAt: createdAt,
	}
	err := contract.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidContract, err)

	// Invalid address
	contract.Symbol = "ETH"
	contract.Address = common.Address{}
	err = contract.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidContract, err)

	// Valid contract
	contract.Address = common.HexToAddress("0x123")
	err = contract.Validate()
	assert.Nil(t, err)
}
