package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/internal/usecase/contract_usecase"
	"github.com/Mugen-Builders/devolt/internal/usecase/user_usecase"
	"github.com/rollmelette/rollmelette"
)

type UserAdvanceHandlers struct {
	UserRepository     entity.UserRepository
	ContractRepository entity.ContractRepository
}

func NewUserAdvanceHandlers(userRepository entity.UserRepository, contractRepository entity.ContractRepository) *UserAdvanceHandlers {
	return &UserAdvanceHandlers{
		UserRepository:     userRepository,
		ContractRepository: contractRepository,
	}
}

func (h *UserAdvanceHandlers) CreateUserHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.CreateUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createUser := user_usecase.NewCreateUserUseCase(h.UserRepository)
	res, err := createUser.Execute(&input, metadata)
	if err != nil {
		return err
	}
	user, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("created user - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) DeleteUserByAddressHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.DeleteUserByAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	deleteUserByAddress := user_usecase.NewDeleteUserByAddressUseCase(h.UserRepository)
	err := deleteUserByAddress.Execute(&input)
	if err != nil {
		return err
	}
	user, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("deleted user with - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) WithdrawHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{Symbol: "STABLECOIN"})
	if err != nil {
		return err
	}
	stablecoinBalance := env.ERC20BalanceOf(stablecoin.Address.Address, metadata.MsgSender)
	if stablecoinBalance.Sign() == 0 {
		return fmt.Errorf("no balance of %v to withdraw", stablecoin.Symbol)
	}
	stablecoinVoucherIndex, err := env.ERC20Withdraw(stablecoin.Address.Address, metadata.MsgSender, stablecoinBalance)
	if err != nil {
		return err
	}
	env.Notice([]byte(fmt.Sprintf("withdrawn %v and %v from %v with voucher index: %v", stablecoin.Symbol, stablecoinBalance, metadata.MsgSender, stablecoinVoucherIndex)))
	return nil
}

func (h *UserAdvanceHandlers) WithdrawApp(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{Symbol: "STABLECOIN"})
	if err != nil {
		return err
	}
	application, isDefined := env.AppAddress()
	if !isDefined {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}
	stablecoinBalance := env.ERC20BalanceOf(stablecoin.Address.Address, application)
	if stablecoinBalance.Sign() == 0 {
		return fmt.Errorf("no balance of %v to withdraw", stablecoin.Symbol)
	}

	if err := env.ERC20Transfer(stablecoin.Address.Address, application, metadata.MsgSender, stablecoinBalance); err != nil {
		return err
	}

	stablecoinVoucherIndex, err := env.ERC20Withdraw(stablecoin.Address.Address, metadata.MsgSender, stablecoinBalance)
	if err != nil {
		return err
	}
	env.Notice([]byte(fmt.Sprintf("withdrawn %v and %v from %v with voucher index: %v", stablecoin.Symbol, stablecoinBalance, metadata.MsgSender, stablecoinVoucherIndex)))
	return nil
}