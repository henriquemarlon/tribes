package user_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type CreateUserInputDTO struct {
	Role    string              `json:"role"`
	Address custom_type.Address `json:"address"`
}

type CreateUserOutputDTO struct {
	Id        uint                `json:"id"`
	Role      string              `json:"role"`
	Address   custom_type.Address `json:"address"`
	CreatedAt int64               `json:"created_at"`
}

type CreateUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewCreateUserUseCase(userRepository entity.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *CreateUserUseCase) Execute(input *CreateUserInputDTO, metadata rollmelette.Metadata) (*CreateUserOutputDTO, error) {
	user, err := entity.NewUser(input.Role, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := u.UserRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &CreateUserOutputDTO{
		Id:        res.Id,
		Role:      res.Role,
		Address:   res.Address,
		CreatedAt: res.CreatedAt,
	}, nil
}
