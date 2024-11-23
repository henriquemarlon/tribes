package repository

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/gorm"
)

type ContractRepositorySqlite struct {
	Db *gorm.DB
}

func NewContractRepositorySqlite(db *gorm.DB) *ContractRepositorySqlite {
	return &ContractRepositorySqlite{
		Db: db,
	}
}

func (r *ContractRepositorySqlite) CreateContract(ctx context.Context, input *entity.Contract) (*entity.Contract, error) {
	err := r.Db.Raw(`
		INSERT INTO contracts (symbol, address, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		RETURNING id
	`, input.Symbol, input.Address.String(), input.CreatedAt, input.UpdatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}
	return input, nil
}

func (r *ContractRepositorySqlite) FindAllContracts(ctx context.Context) ([]*entity.Contract, error) {
	return r.findContractsByQuery("SELECT id, symbol, address, created_at, updated_at FROM contracts")
}

func (r *ContractRepositorySqlite) FindContractBySymbol(ctx context.Context, symbol string) (*entity.Contract, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, symbol, address, created_at, updated_at FROM contracts WHERE symbol = ? LIMIT 1", symbol).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to find contract by symbol: %w", err)
	}
	return r.mapToContractEntity(result), nil
}

func (r *ContractRepositorySqlite) UpdateContract(ctx context.Context, input *entity.Contract) (*entity.Contract, error) {
	existingContract, err := r.FindContractBySymbol(ctx, input.Symbol)
	if err != nil {
		return nil, err
	}

	if input.Address != (common.Address{}) {
		existingContract.Address = input.Address
	}
	existingContract.UpdatedAt = input.UpdatedAt

	res := r.Db.Model(&entity.Contract{}).Where("symbol = ?", input.Symbol).Updates(map[string]interface{}{
		"address":    existingContract.Address.String(),
		"updated_at": existingContract.UpdatedAt,
	})
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update contract: %w", res.Error)
	}
	return existingContract, nil
}

func (r *ContractRepositorySqlite) DeleteContract(ctx context.Context, symbol string) error {
	res := r.Db.Delete(&entity.Contract{}, "symbol = ?", symbol)
	if res.Error != nil {
		return fmt.Errorf("failed to delete contract: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrContractNotFound
	}
	return nil
}

func (r *ContractRepositorySqlite) findContractsByQuery(query string, args ...interface{}) ([]*entity.Contract, error) {
	var results []map[string]interface{}
	err := r.Db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find contracts: %w", err)
	}

	var contracts []*entity.Contract
	for _, result := range results {
		contracts = append(contracts, r.mapToContractEntity(result))
	}
	return contracts, nil
}

func (r *ContractRepositorySqlite) mapToContractEntity(data map[string]interface{}) *entity.Contract {
	return &entity.Contract{
		Id:        uint(data["id"].(int64)),
		Symbol:    data["symbol"].(string),
		Address:   common.HexToAddress(data["address"].(string)),
		CreatedAt: data["created_at"].(int64),
		UpdatedAt: data["updated_at"].(int64),
	}
}
