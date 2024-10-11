package db

import (
	"fmt"
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/pkg/custom_type"
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
	err := r.Db.Create(input).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return input, nil
}

func (r *UserRepositorySqlite) FindUserByRole(role string) (*entity.User, error) {
	var user entity.User
	err := r.Db.Where("role = ?", role).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user by role: %w", err)
	}
	return &user, nil
}

func (r *UserRepositorySqlite) FindUserByAddress(address custom_type.Address) (*entity.User, error) {
	var user entity.User
	err := r.Db.Where("address = ?", address).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user by address %v: %w", address.Address, err)
	}
	return &user, nil
}

func (r *UserRepositorySqlite) FindAllUsers() ([]*entity.User, error) {
	var users []*entity.User
	err := r.Db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}
	return users, nil
}

func (r *UserRepositorySqlite) DeleteUserByAddress(address custom_type.Address) error {
	res := r.Db.Delete(&entity.User{}, "address = ?", address)
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}
