package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
	"github.com/tribeshq/tribes/pkg/router"
)

type UserInspectHandlers struct {
	UserRepository     entity.UserRepository
	ContractRepository entity.ContractRepository
}

func NewUserInspectHandlers(userRepository entity.UserRepository, crowdfundingRepository entity.ContractRepository) *UserInspectHandlers {
	return &UserInspectHandlers{
		UserRepository:     userRepository,
		ContractRepository: crowdfundingRepository,
	}
}

func (h *UserInspectHandlers) FindUserByAddressHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	address := strings.ToLower(router.PathValue(ctx, "address"))
	findUserByAddress := user_usecase.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(ctx, &user_usecase.FindUserByAddressInputDTO{
		Address: common.HexToAddress(address),
	})
	if err != nil {
		return fmt.Errorf("failed to find User: %w", err)
	}
	User, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal User: %w", err)
	}
	env.Report(User)
	return nil
}

func (h *UserInspectHandlers) FindAllUsersHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	findAllUsers := user_usecase.NewFindAllUsersUseCase(h.UserRepository)
	res, err := findAllUsers.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all Users: %w", err)
	}
	allUsers, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all Users: %w", err)
	}
	env.Report(allUsers)
	return nil
}

func (h *UserInspectHandlers) BalanceHandler(ctx context.Context, env rollmelette.EnvInspector) error {
	findAllContracts := contract_usecase.NewFindAllContractsUseCase(h.ContractRepository)
	contracts, err := findAllContracts.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all contracts: %w", err)
	}

	findUserbyAddress := user_usecase.NewFindUserByAddressUseCase(h.UserRepository)
	user, err := findUserbyAddress.Execute(ctx, &user_usecase.FindUserByAddressInputDTO{
		Address: common.HexToAddress(router.PathValue(ctx, "address")),
	})
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	switch user.Role {
	case string(entity.UserRoleAdmin):
		appAddress, isSet := env.AppAddress()
		if !isSet {
			return fmt.Errorf("no application address defined yet, contact the Tribes support")
		}
		balances := make(map[string]string)
		for _, contract := range contracts {
			balances[contract.Symbol] = env.ERC20BalanceOf(contract.Address, appAddress).String()
		}
		balanceBytes, err := json.Marshal(balances)
		if err != nil {
			return fmt.Errorf("failed to marshal balances: %w", err)
		}
		env.Report(balanceBytes)
		return nil
	default:
		balances := make(map[string]string)
		for _, contract := range contracts {
			balances[contract.Symbol] = env.ERC20BalanceOf(contract.Address, common.HexToAddress(router.PathValue(ctx, "address"))).String()
		}
		balanceBytes, err := json.Marshal(balances)
		if err != nil {
			return fmt.Errorf("failed to marshal balances: %w", err)
		}
		env.Report(balanceBytes)
		return nil
	}
}
