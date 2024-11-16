package repository

import (
	"encoding/json"
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

func (r *ContractRepositorySqlite) CreateContract(input *entity.Contract) (*entity.Contract, error) {
	err := r.Db.Model(&entity.Contract{}).Create(map[string]interface{}{
		"symbol":     input.Symbol,
		"address":    input.Address.String(),
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	var result map[string]interface{}
	err = r.Db.Raw("SELECT id, symbol, address, created_at, updated_at FROM contracts WHERE symbol = ? LIMIT 1", input.Symbol).
		Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, err
	}

	contract := &entity.Contract{
		Id:        uint(result["id"].(int64)),
		Symbol:    result["symbol"].(string),
		Address:   common.HexToAddress(result["address"].(string)),
		CreatedAt: result["created_at"].(int64),
		UpdatedAt: result["updated_at"].(int64),
	}
	return contract, nil
}

func (r *ContractRepositorySqlite) FindAllContracts() ([]*entity.Contract, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, symbol, address, created_at, updated_at FROM contracts").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var contracts []*entity.Contract
	for _, result := range results {
		contracts = append(contracts, &entity.Contract{
			Id:        uint(result["id"].(int64)),
			Symbol:    result["symbol"].(string),
			Address:   common.HexToAddress(result["address"].(string)),
			CreatedAt: result["created_at"].(int64),
			UpdatedAt: result["updated_at"].(int64),
		})
	}
	return contracts, nil
}

func (r *ContractRepositorySqlite) FindContractBySymbol(symbol string) (*entity.Contract, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, symbol, address, created_at, updated_at FROM contracts WHERE symbol = ? LIMIT 1", symbol).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, err
	}
	return &entity.Contract{
		Id:        uint(result["id"].(int64)),
		Symbol:    result["symbol"].(string),
		Address:   common.HexToAddress(result["address"].(string)),
		CreatedAt: result["created_at"].(int64),
		UpdatedAt: result["updated_at"].(int64),
	}, nil
}

func (r *ContractRepositorySqlite) UpdateContract(input *entity.Contract) (*entity.Contract, error) {
	var contractJSON map[string]interface{}
	err := r.Db.Where("symbol = ?", input.Symbol).First(&contractJSON).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrContractNotFound
		}
		return nil, err
	}

	contractJSON["address"] = input.Address.String()
	contractJSON["updated_at"] = input.UpdatedAt

	contractBytes, err := json.Marshal(contractJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract: %w", err)
	}

	var contract entity.Contract
	if err = json.Unmarshal(contractBytes, &contract); err != nil {
		return nil, err
	}
	if err = contract.Validate(); err != nil {
		return nil, err
	}

	res := r.Db.Save(contractJSON)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update contract: %w", res.Error)
	}
	return &contract, nil
}

func (r *ContractRepositorySqlite) DeleteContract(symbol string) error {
	res := r.Db.Delete(&entity.Contract{}, "symbol = ?", symbol)
	if res.Error != nil {
		return fmt.Errorf("failed to delete contract: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrContractNotFound
	}
	return nil
}
