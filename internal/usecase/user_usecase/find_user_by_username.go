package user_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindUserByUsernameInputDTO struct {
	Username string `json:"username"`
}

type FindUserByUsernameUseCase struct {
	userRepository entity.UserRepository
}

func NewFindUserByUsernameUseCase(userRepository entity.UserRepository) *FindUserByUsernameUseCase {
	return &FindUserByUsernameUseCase{
		userRepository: userRepository,
	}
}

func (u *FindUserByUsernameUseCase) Execute(input *FindUserByUsernameInputDTO) (*FindUserOutputDTO, error) {
	res, err := u.userRepository.FindUserByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	return &FindUserOutputDTO{
		Id:        res.Id,
		Role:      res.Role,
		Username:  res.Username,
		Address:   res.Address,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
