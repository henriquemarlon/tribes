package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
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
		erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
		if ok {
			findUserByAddress := user_usecase.NewFindUserByAddressUseCase(m.UserRepository)
			user, err := findUserByAddress.Execute(&user_usecase.FindUserByAddressInputDTO{
				Address: erc20Deposit.Sender,
			})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("user not found during RBAC middleware check")
				}
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
				return fmt.Errorf("user with address: %v don't have necessary permission", user.Address)
			}
			return handlerFunc(env, metadata, deposit, payload)
		} else {
			findUserByAddress := user_usecase.NewFindUserByAddressUseCase(m.UserRepository)
			user, err := findUserByAddress.Execute(&user_usecase.FindUserByAddressInputDTO{
				Address: metadata.MsgSender,
			})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("user not found during RBAC middleware check")
				}
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
				return fmt.Errorf("user with address: %v don't have necessary permission", user.Address)
			}
			return handlerFunc(env, metadata, deposit, payload)
		}
	}
}
