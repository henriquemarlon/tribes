package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidContract  = errors.New("invalid contract")
	ErrContractNotFound = errors.New("contract not found")
)

type ContractRepository interface {
	CreateContract(ctx context.Context, contract *Contract) (*Contract, error)
	FindAllContracts(ctx context.Context) ([]*Contract, error)
	FindContractBySymbol(ctx context.Context, symbol string) (*Contract, error)
	FindContractByAddress(ctx context.Context, address common.Address) (*Contract, error)
	UpdateContract(ctx context.Context, contract *Contract) (*Contract, error)
	DeleteContract(ctx context.Context, symbol string) error
}

type Contract struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	Symbol    string         `json:"symbol,omitempty" gorm:"uniqueIndex;not null"`
	Address   common.Address `json:"address,omitempty" gorm:"type:text;not null"`
	CreatedAt int64          `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt int64          `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewContract(symbol string, address common.Address, createdAt int64) (*Contract, error) {
	contract := &Contract{
		Symbol:    symbol,
		Address:   address,
		CreatedAt: createdAt,
	}
	if err := contract.Validate(); err != nil {
		return nil, err
	}
	return contract, nil
}

func (c *Contract) Validate() error {
	if c.Symbol == "" {
		return fmt.Errorf("%w: symbol cannot be empty", ErrInvalidContract)
	}
	if c.Address == (common.Address{}) {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidContract)
	}
	return nil
}
