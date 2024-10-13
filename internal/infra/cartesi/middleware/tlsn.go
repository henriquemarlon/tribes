package middleware

/*
#cgo LDFLAGS: -L../../../../cmd/dapp/lib/target/release -lverifier -lpthread -ldl -lm -lstdc++
#include <stdint.h>

int32_t add_numbers(int32_t a, int32_t b);
*/
import "C"

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
	"github.com/tribeshq/tribes/pkg/custom_type"
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
		findUserByAddress := user_usecase.NewFindUserByAddressUseCase(m.UserRepository)
		user, err := findUserByAddress.Execute(&user_usecase.FindUserByAddressInputDTO{
			Address: custom_type.NewAddress(metadata.MsgSender),
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
