package db

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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
		"role":       input.Role,
		"address":    input.Address.String(),
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var result map[string]interface{}
	err = r.Db.Raw("SELECT id, role, address, created_at, updated_at FROM users WHERE address = ?", input.Address.String()).
		Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}

	user := &entity.User{
		Id:        uint(result["id"].(int64)),
		Role:      result["role"].(string),
		Address:   common.HexToAddress(result["address"].(string)),
		CreatedAt: result["created_at"].(int64),
		UpdatedAt: result["updated_at"].(int64),
	}
	return user, nil
}

func (r *UserRepositorySqlite) FindUserByAddress(address common.Address) (*entity.User, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, role, address, created_at, updated_at FROM users WHERE address = ? LIMIT 1", address.String()).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}
	return &entity.User{
		Id:        uint(result["id"].(int64)),
		Role:      result["role"].(string),
		Address:   common.HexToAddress(result["address"].(string)),
		CreatedAt: result["created_at"].(int64),
		UpdatedAt: result["updated_at"].(int64),
	}, nil
}

func (r *UserRepositorySqlite) FindUserByRole(role string) (*entity.User, error) {
	var result map[string]interface{}
	err := r.Db.Raw("SELECT id, role, address, created_at, updated_at FROM users WHERE role = ? LIMIT 1", role).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}
	return &entity.User{
		Id:        uint(result["id"].(int64)),
		Role:      result["role"].(string),
		Address:   common.HexToAddress(result["address"].(string)),
		CreatedAt: result["created_at"].(int64),
		UpdatedAt: result["updated_at"].(int64),
	}, nil
}

func (r *UserRepositorySqlite) FindAllUsers() ([]*entity.User, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, role, address, created_at, updated_at FROM users").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var users []*entity.User
	for _, result := range results {
		users = append(users, &entity.User{
			Id:        uint(result["id"].(int64)),
			Role:      result["role"].(string),
			Address:   common.HexToAddress(result["address"].(string)),
			CreatedAt: result["created_at"].(int64),
			UpdatedAt: result["updated_at"].(int64),
		})
	}
	return users, nil
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
