package user_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type DeleteUserInputDTO struct {
	Address common.Address `json:"address"`
}

type DeleteUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewDeleteUserUseCase(userRepository entity.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *DeleteUserUseCase) Execute(input *DeleteUserInputDTO) error {
	return u.UserRepository.DeleteUser(input.Address)
}
