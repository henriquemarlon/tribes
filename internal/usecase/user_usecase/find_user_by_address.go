package user_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindUserByAddressInputDTO struct {
	Address custom_type.Address `json:"address"`
}

type FindUserByAddressUseCase struct {
	UserRepository entity.UserRepository
}

func NewFindUserByAddressUseCase(userRepository entity.UserRepository) *FindUserByAddressUseCase {
	return &FindUserByAddressUseCase{
		UserRepository: userRepository,
	}
}

func (u *FindUserByAddressUseCase) Execute(input *FindUserByAddressInputDTO) (*FindUserOutputDTO, error) {
	res, err := u.UserRepository.FindUserByAddress(input.Address)
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
