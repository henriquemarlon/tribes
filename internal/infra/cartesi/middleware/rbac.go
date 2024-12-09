package middleware

import (
	"context"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
	"github.com/tribeshq/tribes/pkg/custom_type"
	"github.com/tribeshq/tribes/pkg/router"
)

type RBACMiddleware struct {
	UserRepository entity.UserRepository
}

func NewRBACMiddleware(userRepository entity.UserRepository) *RBACMiddleware {
	return &RBACMiddleware{
		UserRepository: userRepository,
	}
}

func (m *RBACMiddleware) Middleware(handlerFunc router.AdvanceHandlerFunc, roles []string) router.AdvanceHandlerFunc {
	return func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
		var address custom_type.Address
		ctx := context.Background()
		erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
		if ok {
			address = custom_type.Address(erc20Deposit.Sender)
		} else {
			address = custom_type.Address(metadata.MsgSender)
		}
		findUserByAddress := user_usecase.NewFindUserByAddressUseCase(m.UserRepository)
		user, err := findUserByAddress.Execute(ctx, &user_usecase.FindUserByAddressInputDTO{
			Address: address,
		})
		if err != nil {
			return err
		}
		var hasRole bool
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}
		if !hasRole {
			return fmt.Errorf("user with address: %v does not have necessary permissions: %v", user.Address, roles)
		}
		return handlerFunc(env, metadata, deposit, payload)
	}
}
