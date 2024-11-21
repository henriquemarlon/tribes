package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/order_usecase"
)

type OrderAdvanceHandlers struct {
	UserRepository         entity.UserRepository
	OrderRepository        entity.OrderRepository
	CrowdfundingRepository entity.CrowdfundingRepository
	ContractRepository     entity.ContractRepository
}

func NewOrderAdvanceHandlers(userRepository entity.UserRepository, orderRepository entity.OrderRepository, contractRepository entity.ContractRepository, crowdfundingRepository entity.CrowdfundingRepository) *OrderAdvanceHandlers {
	return &OrderAdvanceHandlers{
		UserRepository:         userRepository,
		OrderRepository:        orderRepository,
		CrowdfundingRepository: crowdfundingRepository,
		ContractRepository:     contractRepository,
	}
}

func (h *OrderAdvanceHandlers) CreateOrderHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order_usecase.CreateOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createOrder := order_usecase.NewCreateOrderUseCase(h.UserRepository, h.OrderRepository, h.ContractRepository, h.CrowdfundingRepository)
	res, err := createOrder.Execute(&input, deposit, metadata)
	if err != nil {
		return err
	}
	// TODO: remove this check when update to V2
	appAddress, isSet := env.AppAddress()
	if !isSet {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}
	if err := env.ERC20Transfer(
		deposit.(*rollmelette.ERC20Deposit).Token,
		metadata.MsgSender,
		appAddress,
		deposit.(*rollmelette.ERC20Deposit).Amount,
	); err != nil {
		return err
	}
	order, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("order created - "), order...))
	return nil
}
