package repository

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
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

func (r *ContractRepositorySqlite) FindContractByAddress(ctx context.Context, address custom_type.Address) (*entity.Contract, error) {
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
	if err := r.Db.WithContext(ctx).Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}
	contract, err := r.FindContractBySymbol(ctx, input.Symbol)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

func (r *ContractRepositorySqlite) DeleteContract(ctx context.Context, symbol string) error {
	if err := r.Db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&entity.Contract{}).Error; err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}
	return nil
}
