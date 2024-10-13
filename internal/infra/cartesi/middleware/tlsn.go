package middleware

import (
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
		_, err := findUserByAddress.Execute(&user_usecase.FindUserByAddressInputDTO{
			Address: custom_type.NewAddress(metadata.MsgSender),
		})
		if err != nil {
			// if errors.Is(err, sql.ErrNoRows) {
			// 	return fmt.Errorf("user not found during RBAC middleware check")
			// }
			return err
		}
		// if user.Role != "creator" {
		// 	return fmt.Errorf("user with address: %v don't have necessary permission", user.Address)
		// }
		// TODO: call tlsn verifier here
		return handlerFunc(env, metadata, deposit, payload)
	}
}
