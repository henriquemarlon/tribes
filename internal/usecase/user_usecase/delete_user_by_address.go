package user_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type DeleteUserInputDTO struct {
	Address custom_type.Address `json:"address"`
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
