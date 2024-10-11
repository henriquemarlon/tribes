package user_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindUserByRoleInputDTO struct {
	Role string `json:"role"`
}

type FindUserByRoleUseCase struct {
	userRepository entity.UserRepository
}

func NewFindUserByRoleUseCase(userRepository entity.UserRepository) *FindUserByRoleUseCase {
	return &FindUserByRoleUseCase{
		userRepository: userRepository,
	}
}

func (u *FindUserByRoleUseCase) Execute(input *FindUserByRoleInputDTO) (*FindUserOutputDTO, error) {
	res, err := u.userRepository.FindUserByRole(input.Role)
	if err != nil {
		return nil, err
	}
	return &FindUserOutputDTO{
		Id:        res.Id,
		Role:      res.Role,
		Address:   res.Address,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
