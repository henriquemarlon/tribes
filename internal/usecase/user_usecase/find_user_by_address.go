package user_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindUserByAddressInputDTO struct {
	Address common.Address `json:"address"`
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
		Id:                res.Id,
		Role:              string(res.Role),
		Address:           res.Address,
		InvestmentLimit:   res.InvestmentLimit,
		DebtIssuanceLimit: res.DebtIssuanceLimit,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
