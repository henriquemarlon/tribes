package entity

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	CreateUser(User *User) (*User, error)
	FindUserByRole(role string) (*User, error)
	FindUserByAddress(address custom_type.Address) (*User, error)
	FindAllUsers() ([]*User, error)
	DeleteUserByAddress(address custom_type.Address) error
}

type User struct {
	Id        uint                `json:"id" gorm:"primaryKey"`
	Role      string              `json:"role,omitempty" gorm:"not null"`
	Username  string              `json:"username,omitempty" gorm:"uniqueIndex;not null"`
	Address   custom_type.Address `json:"address,omitempty" gorm:"type:text;uniqueIndex;not null"`
	CreatedAt int64               `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt int64               `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewUser(role string, username string, address custom_type.Address, created_at int64) (*User, error) {
	user := &User{
		Role:      role,
		Username:  username,
		Address:   address,
		CreatedAt: created_at,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {
	if u.Role == "" || u.Username == "" || u.Address.Address == (common.Address{}) {
		return ErrInvalidUser
	}
	return nil
}
