package middleware

import "C"

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
	"github.com/tribeshq/tribes/pkg/router"
)

type TLSNMiddleware struct {
	UserRepository entity.UserRepository
}

func NewTLSNMiddleware(userRepository entity.UserRepository) *TLSNMiddleware {
	return &TLSNMiddleware{
		UserRepository: userRepository,
	}
}

func (m *TLSNMiddleware) Middleware(handlerFunc router.AdvanceHandlerFunc) router.AdvanceHandlerFunc {
	return func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
		erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
		if !ok {
			return fmt.Errorf("invalid deposit type: %T", deposit)
		}
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
		if user.Role != "creator" {
			return fmt.Errorf("user with address: %v don't have necessary permission", user.Address)
		}

		// TODO: call tlsn verifier here
		// a := C.int32_t(3)
		// b := C.int32_t(4)
		// result := C.add_numbers(a, b)
		// fmt.Printf("result: %d\n", result)

		return handlerFunc(env, metadata, deposit, payload)
	}
}
