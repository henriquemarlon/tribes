package user_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
	// "github.com/ethereum/go-ethereum/common"
)

type DeleteUserByAddressInputDTO struct {
	Address custom_type.Address `json:"address"`
}

type DeleteUserByAddressUseCase struct {
	UserRepository entity.UserRepository
}

func NewDeleteUserByAddressUseCase(userRepository entity.UserRepository) *DeleteUserByAddressUseCase {
	return &DeleteUserByAddressUseCase{
		UserRepository: userRepository,
	}
}

func (u *DeleteUserByAddressUseCase) Execute(input *DeleteUserByAddressInputDTO) error {
	return u.UserRepository.DeleteUserByAddress(input.Address)
}
