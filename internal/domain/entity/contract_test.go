package entity

import (
	"testing"
	"time"

	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewContract(t *testing.T) {
	symbol := "ETH"
	address := custom_type.NewAddress(common.HexToAddress("0x123"))
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
		Address:   custom_type.NewAddress(common.HexToAddress("0x123")),
		CreatedAt: createdAt,
	}
	err := contract.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidContract, err)

	// Invalid address
	contract.Symbol = "ETH"
	contract.Address = custom_type.NewAddress(common.Address{})
	err = contract.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidContract, err)

	// Valid contract
	contract.Address = custom_type.NewAddress(common.HexToAddress("0x123"))
	err = contract.Validate()
	assert.Nil(t, err)
}
