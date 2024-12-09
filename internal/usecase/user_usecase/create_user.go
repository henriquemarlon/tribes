package user_usecase

import (
	"context"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateUserInputDTO struct {
	Role    string              `json:"role"`
	Address custom_type.Address `json:"address"`
}

type CreateUserOutputDTO struct {
	Id                uint                    `json:"id"`
	Role              string                  `json:"role"`
	Address           custom_type.Address     `json:"address"`
	SocialAccounts    []*entity.SocialAccount `json:"social_accounts"`
	InvestmentLimit   *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	DebtIssuanceLimit *uint256.Int            `json:"debt_issuance_limit,omitempty" gorm:"type:bigint"`
	CreatedAt         int64                   `json:"created_at"`
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
	var investmentLimit, debtIssuanceLimit *uint256.Int

	switch input.Role {
	case string(entity.UserRoleQualifiedInvestor):
		investmentLimit = new(uint256.Int).SetAllOne()
	case string(entity.UserRoleNonQualifiedInvestor):
		investmentLimit = uint256.NewInt(20000)
	default:
		investmentLimit = uint256.NewInt(0)
	}

	switch input.Role {
	case string(entity.UserRoleCreator):
		debtIssuanceLimit = uint256.NewInt(15000000)
	default:
		debtIssuanceLimit = uint256.NewInt(0)
	}

	user, err := entity.NewUser(input.Role, investmentLimit, debtIssuanceLimit, input.Address, metadata.BlockTimestamp)
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
