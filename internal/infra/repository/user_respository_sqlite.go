package repository

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/datatype"
	"gorm.io/gorm"
)

type UserRepositorySqlite struct {
	Db *gorm.DB
}

func NewUserRepositorySqlite(db *gorm.DB) *UserRepositorySqlite {
	return &UserRepositorySqlite{Db: db}
}

func (r *UserRepositorySqlite) CreateUser(ctx context.Context, input *entity.User) (*entity.User, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return input, nil
}

func (r *UserRepositorySqlite) FindUserByAddress(ctx context.Context, address datatype.Address) (*entity.User, error) {
	var user entity.User
	if err := r.Db.WithContext(ctx).Preload("SocialAccounts").Where("address = ?", address).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by address: %w", err)
	}
	return &user, nil
}

func (r *UserRepositorySqlite) FindUsersByRole(ctx context.Context, role string) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.Db.WithContext(ctx).Preload("SocialAccounts").Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}
	return users, nil
}

func (r *UserRepositorySqlite) FindAllUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.Db.WithContext(ctx).Preload("SocialAccounts").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
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

	if err := r.Db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (r *UserRepositorySqlite) DeleteUser(ctx context.Context, address datatype.Address) error {
	res := r.Db.WithContext(ctx).Where("address = ?", address).Delete(&entity.User{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}
