package user_usecase

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CreateUserInputDTO struct {
	Role    string         `json:"role"`
	Address common.Address `json:"address"`
}

type CreateUserOutputDTO struct {
	Id                uint           `json:"id"`
	Role              string         `json:"role"`
	Address           common.Address `json:"address"`
	InvestmentLimit   *uint256.Int   `json:"investment_limit,omitempty" gorm:"type:bigint"`
	DebtIssuanceLimit *uint256.Int   `json:"debt_issuance_limit,omitempty" gorm:"type:bigint"`
	CreatedAt         int64          `json:"created_at"`
}

type CreateUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewCreateUserUseCase(userRepository entity.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, input *CreateUserInputDTO, metadata rollmelette.Metadata) (*CreateUserOutputDTO, error) {
	user, err := entity.NewUser(input.Role, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := u.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &CreateUserOutputDTO{
		Id:                res.Id,
		Role:              string(res.Role),
		Address:           res.Address,
		InvestmentLimit:   res.InvestmentLimit,
		DebtIssuanceLimit: res.DebtIssuanceLimit,
		CreatedAt:         res.CreatedAt,
	}, nil
}
