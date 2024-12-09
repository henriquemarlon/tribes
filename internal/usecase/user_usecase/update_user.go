package user_usecase

import (
	"context"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type UpdateUserInputDTO struct {
	Role              string              `json:"role"`
	Address           custom_type.Address `json:"address"`
	InvestmentLimit   *uint256.Int        `json:"investment_limit,omitempty" gorm:"type:bigint"`
	DebtIssuanceLimit *uint256.Int        `json:"debt_issuance_limit,omitempty" gorm:"type:bigint"`
}

type UpdateUserOutputDTO struct {
	Id                uint                    `json:"id"`
	Role              string                  `json:"role"`
	Address           custom_type.Address     `json:"address"`
	SocialAccounts    []*entity.SocialAccount `json:"social_accounts"`
	InvestmentLimit   *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	DebtIssuanceLimit *uint256.Int            `json:"debt_issuance_limit,omitempty" gorm:"type:bigint"`
	CreatedAt         int64                   `json:"created_at"`
	UpdatedAt         int64                   `json:"updated_at"`
}

type UpdateUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewUpdateUserUseCase(userRepository entity.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *UpdateUserUseCase) Execute(ctx context.Context, input *UpdateUserInputDTO, metadata rollmelette.Metadata) (*UpdateUserOutputDTO, error) {
	user, err := u.UserRepository.UpdateUser(ctx, &entity.User{
		Role:              entity.UserRole(input.Role),
		Address:           input.Address,
		InvestmentLimit:   input.InvestmentLimit,
		DebtIssuanceLimit: input.DebtIssuanceLimit,
		UpdatedAt:         metadata.BlockTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateUserOutputDTO{
		Id:                user.Id,
		Role:              string(user.Role),
		Address:           user.Address,
		SocialAccounts:    user.SocialAccounts,
		InvestmentLimit:   user.InvestmentLimit,
		DebtIssuanceLimit: user.DebtIssuanceLimit,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}, nil
}
