package user_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
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

func (u *FindAllUsersUseCase) Execute() (*FindAllUsersOutputDTO, error) {
	res, err := u.UserRepository.FindAllUsers()
	if err != nil {
		return nil, err
	}
	output := make(FindAllUsersOutputDTO, len(res))
	for i, user := range res {
		output[i] = &FindUserOutputDTO{
			Id:        user.Id,
			Role:      user.Role,
			Address:   user.Address,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return &output, nil
}
