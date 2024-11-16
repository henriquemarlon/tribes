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

func (m *RBACMiddleware) Middleware(handlerFunc router.AdvanceHandlerFunc, role string) router.AdvanceHandlerFunc {
	return func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
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
		if user.Role != role && user.Role != "admin" {
			return fmt.Errorf("user with address: %v don't have necessary permission: %v", user.Address, role)
		}
		return handlerFunc(env, metadata, deposit, payload)
	}
}
