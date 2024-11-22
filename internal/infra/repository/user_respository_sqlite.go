package repository

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/holiman/uint256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
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
	err := r.Db.Raw(`
		INSERT INTO users (role, address, investment_limit, debt_issuance_limit, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`, input.Role, input.Address.String(), input.InvestmentLimit.Hex(), input.DebtIssuanceLimit.Hex(), input.CreatedAt, input.UpdatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return input, nil
}

func (r *UserRepositorySqlite) FindUserByAddress(address common.Address) (*entity.User, error) {
	return r.findUserByQuery("SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at FROM users WHERE address = ? LIMIT 1", address.String())
}

func (r *UserRepositorySqlite) FindUsersByRole(role string) ([]*entity.User, error) {
	return r.findUsersByQuery("SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at FROM users WHERE role = ?", role)
}

func (r *UserRepositorySqlite) FindAllUsers() ([]*entity.User, error) {
	return r.findUsersByQuery("SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at FROM users")
}

func (r *UserRepositorySqlite) UpdateUser(input *entity.User) (*entity.User, error) {
	existingUser, err := r.FindUserByAddress(input.Address)
	if err != nil {
		return nil, err
	}

	if input.Role != "" {
		existingUser.Role = input.Role
	}
	if input.InvestmentLimit != nil && input.InvestmentLimit.Sign() > 0 {
		existingUser.InvestmentLimit = input.InvestmentLimit
	}
	if input.DebtIssuanceLimit != nil && input.DebtIssuanceLimit.Sign() > 0 {
		existingUser.DebtIssuanceLimit = input.DebtIssuanceLimit
	}
	existingUser.UpdatedAt = input.UpdatedAt

	res := r.Db.Model(&entity.User{}).Where("id = ?", existingUser.Id).Updates(map[string]interface{}{
		"role":                existingUser.Role,
		"investment_limit":    existingUser.InvestmentLimit.Hex(),
		"debt_issuance_limit": existingUser.DebtIssuanceLimit.Hex(),
		"updated_at":          existingUser.UpdatedAt,
	})
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", res.Error)
	}
	return existingUser, nil
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

func (r *UserRepositorySqlite) findUserByQuery(query string, args ...interface{}) (*entity.User, error) {
	var result map[string]interface{}
	err := r.Db.Raw(query, args...).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.mapToUserEntity(result), nil
}

func (r *UserRepositorySqlite) findUsersByQuery(query string, args ...interface{}) ([]*entity.User, error) {
	var results []map[string]interface{}
	err := r.Db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	var users []*entity.User
	for _, result := range results {
		users = append(users, r.mapToUserEntity(result))
	}
	return users, nil
}

func (r *UserRepositorySqlite) mapToUserEntity(data map[string]interface{}) *entity.User {
	return &entity.User{
		Id:                uint(data["id"].(int64)),
		Role:              entity.UserRole(data["role"].(string)),
		Address:           common.HexToAddress(data["address"].(string)),
		InvestmentLimit:   uint256.MustFromHex(data["investment_limit"].(string)),
		DebtIssuanceLimit: uint256.MustFromHex(data["debt_issuance_limit"].(string)),
		CreatedAt:         data["created_at"].(int64),
		UpdatedAt:         data["updated_at"].(int64),
	}
}
