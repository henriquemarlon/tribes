package user_usecase

import (
	"context"

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

func (u *DeleteUserUseCase) Execute(ctx context.Context, input *DeleteUserInputDTO) error {
	return u.UserRepository.DeleteUser(ctx, input.Address)
}
