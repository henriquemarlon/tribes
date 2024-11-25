package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
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
	ctx := context.Background()
	createUser := user_usecase.NewCreateUserUseCase(h.UserRepository)
	res, err := createUser.Execute(ctx, &input, metadata)
	if err != nil {
		return err
	}
	user, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("user created - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) UpdateUserHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.UpdateUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	updateUser := user_usecase.NewUpdateUserUseCase(h.UserRepository)
	res, err := updateUser.Execute(ctx, &input, metadata)
	if err != nil {
		return err
	}
	user, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("user updated - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) DeleteUserHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.DeleteUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	deleteUserByAddress := user_usecase.NewDeleteUserUseCase(h.UserRepository)
	err := deleteUserByAddress.Execute(ctx, &input)
	if err != nil {
		return err
	}
	user, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("user deleted - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) WithdrawHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	// TODO: remove this check when update to V2
	appAddress, isSet := env.AppAddress()
	if !isSet {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}
	var input user_usecase.WithdrawInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	findUserByAddress := user_usecase.NewCreateUserUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(ctx, &user_usecase.CreateUserInputDTO{
		Address: metadata.MsgSender,
	}, metadata)
	if err != nil {
		return err
	}

	switch entity.UserRole(res.Role) {
	case entity.UserRoleAdmin:
		// The Admin role can withdraw the entire Application Balance if wanted
		err := env.ERC20Transfer(
			input.Token,
			appAddress,
			metadata.MsgSender,
			input.Amount.ToBig(),
		)
		if err != nil {
			return err
		}
		_, err = env.ERC20Withdraw(
			input.Token,
			metadata.MsgSender,
			input.Amount.ToBig(),
		)
		if err != nil {
			return err
		}
	default:
		_, err := env.ERC20Withdraw(
			input.Token,
			metadata.MsgSender,
			input.Amount.ToBig(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
