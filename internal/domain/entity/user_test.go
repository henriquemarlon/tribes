package entity

import (
	"testing"
	"time"

	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewUser_Success(t *testing.T) {
	role := "admin"
	address := custom_type.Address{Address: common.HexToAddress("0x123")}
	createdAt := time.Now().Unix()

	user, err := NewUser(role, address, createdAt)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, role, user.Role)
	assert.Equal(t, address, user.Address)
	assert.Equal(t, createdAt, user.CreatedAt)
}

func TestNewUser_Fail_InvalidUser(t *testing.T) {
	role := ""
	address := custom_type.Address{Address: common.HexToAddress("0x123")}
	createdAt := time.Now().Unix()

	user, err := NewUser(role, address, createdAt)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrInvalidUser, err)

	role = "admin"
	address = custom_type.Address{Address: common.Address{}}

	user, err = NewUser(role, address, createdAt)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrInvalidUser, err)
}
