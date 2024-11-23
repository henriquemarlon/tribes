package repository

import (
	"context"
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

func (r *UserRepositorySqlite) CreateUser(ctx context.Context, input *entity.User) (*entity.User, error) {
	err := r.Db.WithContext(ctx).Raw(`
		INSERT INTO users (role, address, investment_limit, debt_issuance_limit, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`, input.Role, input.Address.String(), input.InvestmentLimit.Hex(), input.DebtIssuanceLimit.Hex(), input.CreatedAt, input.UpdatedAt).Scan(&input.Id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return input, nil
}

func (r *UserRepositorySqlite) FindUserByAddress(ctx context.Context, address common.Address) (*entity.User, error) {
	var result map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at 
		FROM users WHERE address = ? LIMIT 1
	`, address.String()).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by address: %w", err)
	}

	return &entity.User{
		Id:                uint(result["id"].(int64)),
		Role:              entity.UserRole(result["role"].(string)),
		Address:           common.HexToAddress(result["address"].(string)),
		InvestmentLimit:   uint256.MustFromHex(result["investment_limit"].(string)),
		DebtIssuanceLimit: uint256.MustFromHex(result["debt_issuance_limit"].(string)),
		CreatedAt:         result["created_at"].(int64),
		UpdatedAt:         result["updated_at"].(int64),
	}, nil
}

func (r *UserRepositorySqlite) FindUsersByRole(ctx context.Context, role string) ([]*entity.User, error) {
	var results []map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at 
		FROM users WHERE role = ?
	`, role).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	var users []*entity.User
	for _, data := range results {
		users = append(users, &entity.User{
			Id:                uint(data["id"].(int64)),
			Role:              entity.UserRole(data["role"].(string)),
			Address:           common.HexToAddress(data["address"].(string)),
			InvestmentLimit:   uint256.MustFromHex(data["investment_limit"].(string)),
			DebtIssuanceLimit: uint256.MustFromHex(data["debt_issuance_limit"].(string)),
			CreatedAt:         data["created_at"].(int64),
			UpdatedAt:         data["updated_at"].(int64),
		})
	}
	return users, nil
}

func (r *UserRepositorySqlite) FindAllUsers(ctx context.Context) ([]*entity.User, error) {
	var results []map[string]interface{}
	err := r.Db.WithContext(ctx).Raw(`
		SELECT id, role, address, investment_limit, debt_issuance_limit, created_at, updated_at 
		FROM users
	`).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	var users []*entity.User
	for _, data := range results {
		users = append(users, &entity.User{
			Id:                uint(data["id"].(int64)),
			Role:              entity.UserRole(data["role"].(string)),
			Address:           common.HexToAddress(data["address"].(string)),
			InvestmentLimit:   uint256.MustFromHex(data["investment_limit"].(string)),
			DebtIssuanceLimit: uint256.MustFromHex(data["debt_issuance_limit"].(string)),
			CreatedAt:         data["created_at"].(int64),
			UpdatedAt:         data["updated_at"].(int64),
		})
	}
	return users, nil
}

func (r *UserRepositorySqlite) UpdateUser(ctx context.Context, input *entity.User) (*entity.User, error) {
	user, err := r.FindUserByAddress(ctx, input.Address)
	if err != nil {
		return nil, err
	}

	if input.Role != "" {
		user.Role = input.Role
	}
	if input.InvestmentLimit != nil && input.InvestmentLimit.Sign() > 0 {
		user.InvestmentLimit = input.InvestmentLimit
	}
	if input.DebtIssuanceLimit != nil && input.DebtIssuanceLimit.Sign() > 0 {
		user.DebtIssuanceLimit = input.DebtIssuanceLimit
	}
	user.UpdatedAt = input.UpdatedAt

	res := r.Db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.Id).Updates(map[string]interface{}{
		"role":                user.Role,
		"investment_limit":    user.InvestmentLimit.Hex(),
		"debt_issuance_limit": user.DebtIssuanceLimit.Hex(),
		"updated_at":          user.UpdatedAt,
	})
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", res.Error)
	}
	return user, nil
}

func (r *UserRepositorySqlite) DeleteUser(ctx context.Context, address common.Address) error {
	res := r.Db.Delete(&entity.User{}, "address = ?", address.String())
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}
