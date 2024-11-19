package advance_handler

import (
	// "encoding/json"
	// "fmt"

	"encoding/json"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/order_usecase"
	// "github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	// "github.com/tribeshq/tribes/internal/usecase/order_usecase"
	// "github.com/tribeshq/tribes/internal/usecase/user_usecase"
)

type OrderAdvanceHandlers struct {
	OrderRepository        entity.OrderRepository
	UserRepository         entity.UserRepository
	CrowdfundingRepository entity.CrowdfundingRepository
	ContractRepository     entity.ContractRepository
}

func NewOrderAdvanceHandlers(orderRepository entity.OrderRepository, userRepository entity.UserRepository, contractRepository entity.ContractRepository, crowdfundingRepository entity.CrowdfundingRepository) *OrderAdvanceHandlers {
	return &OrderAdvanceHandlers{
		OrderRepository:        orderRepository,
		UserRepository:         userRepository,
		CrowdfundingRepository: crowdfundingRepository,
		ContractRepository:     contractRepository,
	}
}

func (h *OrderAdvanceHandlers) CreateOrderHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order_usecase.CreateOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createOrder := order_usecase.NewCreateOrderUseCase(h.OrderRepository, h.ContractRepository, h.CrowdfundingRepository)
	res, err := createOrder.Execute(&input, deposit, metadata)
	if err != nil {
		return err
	}
	order, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("order created - "), order...))
	return nil
}