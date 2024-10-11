package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/internal/usecase/contract_usecase"
	"github.com/Mugen-Builders/devolt/internal/usecase/user_usecase"
	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/Mugen-Builders/devolt/pkg/router"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
)

type UserInspectHandlers struct {
	UserRepository     entity.UserRepository
	ContractRepository entity.ContractRepository
}

func NewUserInspectHandlers(userRepository entity.UserRepository, contractRepository entity.ContractRepository) *UserInspectHandlers {
	return &UserInspectHandlers{
		UserRepository:     userRepository,
		ContractRepository: contractRepository,
	}
}

func (h *UserInspectHandlers) FindUserByAddressHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	address := strings.ToLower(router.PathValue(ctx, "address"))
	findUserByAddress := user_usecase.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(&user_usecase.FindUserByAddressInputDTO{
		Address: custom_type.NewAddress(common.HexToAddress(address)),
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
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	contract, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{
		Symbol: strings.ToUpper(router.PathValue(ctx, "symbol")),
	})
	if err != nil {
		return fmt.Errorf("failed to find contract: %w", err)
	}
	balanceBytes, err := json.Marshal(env.ERC20BalanceOf(contract.Address.Address, common.HexToAddress(router.PathValue(ctx, "address"))))
	if err != nil {
		return fmt.Errorf("failed to marshal balance: %w", err)
	}
	env.Report(balanceBytes)
	return nil
}
