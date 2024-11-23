package user_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindUserByRoleInputDTO struct {
	Role string `json:"role"`
}

type FindUserByRoleOutputDTO []*FindUserOutputDTO

type FindUserByRoleUseCase struct {
	userRepository entity.UserRepository
}

func NewFindUserByRoleUseCase(userRepository entity.UserRepository) *FindUserByRoleUseCase {
	return &FindUserByRoleUseCase{
		userRepository: userRepository,
	}
}

func (u *FindUserByRoleUseCase) Execute(ctx context.Context, input *FindUserByRoleInputDTO) ([]*FindUserOutputDTO, error) {
	res, err := u.userRepository.FindUsersByRole(ctx, input.Role)
	if err != nil {
		return nil, err
	}
	output := make(FindUserByRoleOutputDTO, len(res))
	for i, user := range res {
		output[i] = &FindUserOutputDTO{
			Id:                user.Id,
			Role:              string(user.Role),
			Address:           user.Address,
			InvestmentLimit:   user.InvestmentLimit,
			DebtIssuanceLimit: user.DebtIssuanceLimit,
			CreatedAt:         user.CreatedAt,
			UpdatedAt:         user.UpdatedAt,
		}
	}
	return output, nil
}
