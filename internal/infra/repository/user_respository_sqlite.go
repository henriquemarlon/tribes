package repository

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepositorySqlite struct {
	Db *gorm.DB
}

func NewUserRepositorySqlite(db *gorm.DB) *UserRepositorySqlite {
	return &UserRepositorySqlite{
		Db: db,
	}
}

func (r *UserRepositorySqlite) CreateUser(input *entity.User) (*entity.User, error) {
	err := r.Db.Model(&entity.User{}).Create(map[string]interface{}{
		"role":                input.Role,
		"address":             input.Address.String(),
		"investment_limit":    input.InvestmentLimit.String(),
		"debt_issuance_limit": input.DebtIssuanceLimit.String(),
		"created_at":          input.CreatedAt,
		"updated_at":          input.UpdatedAt,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return r.FindUserByAddress(input.Address)
}

func (r *UserRepositorySqlite) FindUserByAddress(address common.Address) (*entity.User, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at FROM users WHERE address = ? LIMIT 1", address.String()).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}
	return &entity.User{
		Id:                uint(result["id"].(int64)),
		Role:              result["role"].(string),
		Address:           common.HexToAddress(result["address"].(string)),
		InvestmentLimit:   uint256.MustFromHex(result["investment_limit"].(string)),
		DebtIssuanceLimit: uint256.MustFromHex(result["debt_issuance_limit"].(string)),
		CreatedAt:         result["created_at"].(int64),
		UpdatedAt:         result["updated_at"].(int64),
	}, nil
}

func (r *UserRepositorySqlite) FindUsersByRole(role string) ([]*entity.User, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at FROM users WHERE role = ?", role).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}
	var users []*entity.User
	for _, result := range results {
		users = append(users, &entity.User{
			Id:                uint(result["id"].(int64)),
			Role:              result["role"].(string),
			Address:           common.HexToAddress(result["address"].(string)),
			InvestmentLimit:   uint256.MustFromHex(result["investment_limit"].(string)),
			DebtIssuanceLimit: uint256.MustFromHex(result["debt_issuance_limit"].(string)),
			CreatedAt:         result["created_at"].(int64),
			UpdatedAt:         result["updated_at"].(int64),
		})
	}
	return users, nil
}

func (r *UserRepositorySqlite) FindAllUsers() ([]*entity.User, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at FROM users").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	for _, result := range results {
		users = append(users, &entity.User{
			Id:                uint(result["id"].(int64)),
			Role:              result["role"].(string),
			Address:           common.HexToAddress(result["address"].(string)),
			InvestmentLimit:   uint256.MustFromHex(result["investment_limit"].(string)),
			DebtIssuanceLimit: uint256.MustFromHex(result["debt_issuance_limit"].(string)),
			CreatedAt:         result["created_at"].(int64),
			UpdatedAt:         result["updated_at"].(int64),
		})
	}
	return users, nil
}

func (r *UserRepositorySqlite) UpdateUser(input *entity.User) (*entity.User, error) {
	var userJSON map[string]interface{}
	err := r.Db.Where("address = ?", input.Address.String()).First(&userJSON).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	if input.Role != "" {
		userJSON["role"] = input.Role
	}
	if input.Address != (common.Address{}) {
		userJSON["address"] = input.Address.String()
	}
	if input.InvestmentLimit != nil {
		userJSON["investment_limit"] = input.InvestmentLimit.String()
	}
	if input.InvestmentLimit != nil {
		userJSON["debt_issuance_limit"] = input.DebtIssuanceLimit.String()
	}
	userJSON["updated_at"] = input.UpdatedAt

	userBytes, err := json.Marshal(userJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	var user entity.User
	if err = json.Unmarshal(userBytes, &user); err != nil {
		return nil, err
	}
	if err = user.Validate(); err != nil {
		return nil, err
	}

	res := r.Db.Save(&user)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", res.Error)
	}
	return &user, nil
}


func (r *UserRepositorySqlite) DeleteUser(address common.Address) error {
	res := r.Db.Delete(&entity.User{}, "address = ?", address.String())
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}
