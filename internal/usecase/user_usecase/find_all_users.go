package user_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FindAllUsersOutputDTO []*FindUserOutputDTO

type FindAllUsersUseCase struct {
	UserRepository entity.UserRepository
}

func NewFindAllUsersUseCase(userRepository entity.UserRepository) *FindAllUsersUseCase {
	return &FindAllUsersUseCase{
		UserRepository: userRepository,
	}
}

func (u *FindAllUsersUseCase) Execute(ctx context.Context) (*FindAllUsersOutputDTO, error) {
	res, err := u.UserRepository.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllUsersOutputDTO, len(res))
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
	return &output, nil
}
