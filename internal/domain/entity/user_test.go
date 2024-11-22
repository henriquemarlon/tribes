package entity

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	address := common.HexToAddress("0x123")
	createdAt := time.Now().Unix()

	t.Run("Valid User - Qualified Investor", func(t *testing.T) {
		role := string(UserRoleQualifiedInvestor)
		user, err := NewUser(role, address, createdAt)

		assert.Nil(t, err, "Unexpected error when creating a valid qualified investor user")
		assert.NotNil(t, user, "The created user should not be nil")
		assert.Equal(t, UserRole(role), user.Role, "Role is incorrect")
		assert.Equal(t, address, user.Address, "Address is incorrect")
		assert.Equal(t, createdAt, user.CreatedAt, "Creation date is incorrect")
		assert.Equal(t, new(uint256.Int).SetAllOne(), user.InvestmentLimit, "Investment limit is incorrect for a qualified investor")
		assert.Equal(t, uint256.NewInt(0), user.DebtIssuanceLimit, "Debt issuance limit is incorrect for a qualified investor")
	})

	t.Run("Valid User - Non-Qualified Investor", func(t *testing.T) {
		role := string(UserRoleNonQualifiedInvestor)
		user, err := NewUser(role, address, createdAt)

		assert.Nil(t, err, "Unexpected error when creating a valid non-qualified investor user")
		assert.NotNil(t, user, "The created user should not be nil")
		assert.Equal(t, UserRole(role), user.Role, "Role is incorrect")
		assert.Equal(t, uint256.NewInt(20000), user.InvestmentLimit, "Investment limit is incorrect for a non-qualified investor")
		assert.Equal(t, uint256.NewInt(0), user.DebtIssuanceLimit, "Debt issuance limit is incorrect for a non-qualified investor")
	})

	t.Run("Valid User - Creator", func(t *testing.T) {
		role := string(UserRoleCreator)
		user, err := NewUser(role, address, createdAt)

		assert.Nil(t, err, "Unexpected error when creating a valid creator user")
		assert.NotNil(t, user, "The created user should not be nil")
		assert.Equal(t, UserRole(role), user.Role, "Role is incorrect")
		assert.Equal(t, uint256.NewInt(0), user.InvestmentLimit, "Investment limit should be zero for a creator")
		assert.Equal(t, uint256.NewInt(15000000), user.DebtIssuanceLimit, "Debt issuance limit is incorrect for a creator")
	})

	t.Run("Invalid User - Empty Role", func(t *testing.T) {
		_, err := NewUser("", address, createdAt)

		assert.NotNil(t, err, "An error should be returned for an empty role")
		assert.True(t, errors.Is(err, ErrInvalidUser), "The error should be ErrInvalidUser")
		assert.Contains(t, err.Error(), "role cannot be empty", "The error message should mention the empty role")
	})

	t.Run("Invalid User - Empty Address", func(t *testing.T) {
		_, err := NewUser(string(UserRoleAdmin), common.Address{}, createdAt)

		assert.NotNil(t, err, "An error should be returned for an empty address")
		assert.True(t, errors.Is(err, ErrInvalidUser), "The error should be ErrInvalidUser")
		assert.Contains(t, err.Error(), "address cannot be empty", "The error message should mention the empty address")
	})

	t.Run("Invalid User - Missing Creation Date", func(t *testing.T) {
		_, err := NewUser(string(UserRoleAdmin), address, 0)

		assert.NotNil(t, err, "An error should be returned for a missing creation date")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "creation date is missing", "The error message should mention the missing creation date")
	})
}

func TestUser_Validate(t *testing.T) {
	t.Run("Valid User", func(t *testing.T) {
		user := &User{
			Role:              UserRoleAdmin,
			Address:           common.HexToAddress("0x123"),
			CreatedAt:         time.Now().Unix(),
			InvestmentLimit:   uint256.NewInt(0),
			DebtIssuanceLimit: uint256.NewInt(0),
		}
		err := user.Validate()
		assert.Nil(t, err, "No error should be returned for a valid user")
	})

	t.Run("Invalid Role", func(t *testing.T) {
		user := &User{
			Role:              "",
			Address:           common.HexToAddress("0x123"),
			CreatedAt:         time.Now().Unix(),
			InvestmentLimit:   uint256.NewInt(0),
			DebtIssuanceLimit: uint256.NewInt(0),
		}
		err := user.Validate()
		assert.NotNil(t, err, "An error should be returned for an empty role")
		assert.True(t, errors.Is(err, ErrInvalidUser), "The error should be ErrInvalidUser")
		assert.Contains(t, err.Error(), "role cannot be empty", "The error message should mention the empty role")
	})

	t.Run("Invalid Address", func(t *testing.T) {
		user := &User{
			Role:              UserRoleAdmin,
			Address:           common.Address{},
			CreatedAt:         time.Now().Unix(),
			InvestmentLimit:   uint256.NewInt(0),
			DebtIssuanceLimit: uint256.NewInt(0),
		}
		err := user.Validate()
		assert.NotNil(t, err, "An error should be returned for an empty address")
		assert.True(t, errors.Is(err, ErrInvalidUser), "The error should be ErrInvalidUser")
		assert.Contains(t, err.Error(), "address cannot be empty", "The error message should mention the empty address")
	})

	t.Run("Invalid Creation Date", func(t *testing.T) {
		user := &User{
			Role:              UserRoleAdmin,
			Address:           common.HexToAddress("0x123"),
			CreatedAt:         0,
			InvestmentLimit:   uint256.NewInt(0),
			DebtIssuanceLimit: uint256.NewInt(0),
		}
		err := user.Validate()
		assert.NotNil(t, err, "An error should be returned for a missing creation date")
		assert.True(t, errors.Is(err, ErrInvalidCrowdfunding), "The error should be ErrInvalidCrowdfunding")
		assert.Contains(t, err.Error(), "creation date is missing", "The error message should mention the missing creation date")
	})
}
