package entity

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)

type UserRole string

const (
	UserRoleAdmin                UserRole = "admin"
	UserRoleCreator              UserRole = "creator"
	UserRoleNonQualifiedInvestor UserRole = "non_qualified_investor"
	UserRoleQualifiedInvestor    UserRole = "qualified_investor"
)

type UserRepository interface {
	CreateUser(User *User) (*User, error)
	FindUserByRole(role string) (*User, error)
	FindUserByAddress(address common.Address) (*User, error)
	FindAllUsers() ([]*User, error)
	DeleteUser(address common.Address) error
}

type User struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	Role      string         `json:"role,omitempty" gorm:"not null"`
	Address   common.Address `json:"address,omitempty" gorm:"type:text;uniqueIndex;not null"`
	CreatedAt int64          `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt int64          `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewUser(role string, address common.Address, created_at int64) (*User, error) {
	user := &User{
		Role:      role,
		Address:   address,
		CreatedAt: created_at,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {
	if u.Role == "" || u.Address == (common.Address{}) {
		return ErrInvalidUser
	}
	return nil
}
