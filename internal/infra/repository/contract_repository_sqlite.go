package repository

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/datatype"
	"gorm.io/gorm"
)

type ContractRepositorySqlite struct {
	Db *gorm.DB
}

func NewContractRepositorySqlite(db *gorm.DB) *ContractRepositorySqlite {
	return &ContractRepositorySqlite{Db: db}
}

func (r *ContractRepositorySqlite) CreateContract(ctx context.Context, input *entity.Contract) (*entity.Contract, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}
	return input, nil
}

func (r *ContractRepositorySqlite) FindAllContracts(ctx context.Context) ([]*entity.Contract, error) {
	var contracts []*entity.Contract
	if err := r.Db.WithContext(ctx).Find(&contracts).Error; err != nil {
		return nil, fmt.Errorf("failed to find all contracts: %w", err)
	}
	return contracts, nil
}

func (r *ContractRepositorySqlite) FindContractBySymbol(ctx context.Context, symbol string) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.Db.WithContext(ctx).Where("symbol = ?", symbol).First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract by symbol: %w", err)
	}
	return &contract, nil
}

func (r *ContractRepositorySqlite) FindContractByAddress(ctx context.Context, address datatype.Address) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.Db.WithContext(ctx).Where("address = ?", address).First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract by address: %w", err)
	}
	return &contract, nil
}

func (r *ContractRepositorySqlite) UpdateContract(ctx context.Context, input *entity.Contract) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.Db.WithContext(ctx).Where("symbol = ?", input.Symbol).First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract for update: %w", err)
	}

	if input.Address != (datatype.Address{}) {
		contract.Address = input.Address
	}
	contract.UpdatedAt = input.UpdatedAt

	if err := r.Db.WithContext(ctx).Save(&contract).Error; err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}
	return &contract, nil
}

func (r *ContractRepositorySqlite) DeleteContract(ctx context.Context, symbol string) error {
	if err := r.Db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&entity.Contract{}).Error; err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}
	return nil
}
