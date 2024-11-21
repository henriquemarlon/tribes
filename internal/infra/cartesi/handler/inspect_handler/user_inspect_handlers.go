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

func (h *UserInspectHandlers) FindUserByAddressHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	address := strings.ToLower(router.PathValue(ctx, "address"))
	findUserByAddress := user_usecase.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(&user_usecase.FindUserByAddressInputDTO{
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

func (h *UserInspectHandlers) FindAllUsersHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	findAllUsers := user_usecase.NewFindAllUsersUseCase(h.UserRepository)
	res, err := findAllUsers.Execute()
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

func (h *UserInspectHandlers) BalanceHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	findAllContracts := contract_usecase.NewFindAllContractsUseCase(h.ContractRepository)
	contracts, err := findAllContracts.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all contracts: %w", err)
	}
	balances := make(map[string]string)
	for _, contract := range contracts {
		balances[contract.Symbol] = env.ERC20BalanceOf(contract.Address, contract.Address).String()
	}
	balanceBytes, err := json.Marshal(balances)
	if err != nil {
		return fmt.Errorf("failed to marshal balances: %w", err)
	}
	env.Report(balanceBytes)
	return nil
}
