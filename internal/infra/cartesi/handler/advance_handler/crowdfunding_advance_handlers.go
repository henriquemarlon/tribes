package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/crowdfunding_usecase"
)

type CrowdfundingAdvanceHandlers struct {
	OrderRepository        entity.OrderRepository
	UserRepository         entity.UserRepository
	CrowdfundingRepository entity.CrowdfundingRepository
	ContractRepository     entity.ContractRepository
}

func NewCrowdfundingAdvanceHandlers(
	orderRepository entity.OrderRepository,
	userRepository entity.UserRepository,
	crowdfundingRepository entity.CrowdfundingRepository,
	contractRepository entity.ContractRepository,
) *CrowdfundingAdvanceHandlers {
	return &CrowdfundingAdvanceHandlers{
		OrderRepository:        orderRepository,
		UserRepository:         userRepository,
		CrowdfundingRepository: crowdfundingRepository,
		ContractRepository:     contractRepository,
	}
}

func (h *CrowdfundingAdvanceHandlers) CreateCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *crowdfunding_usecase.CreateCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createCrowdfunding := crowdfunding_usecase.NewCreateCrowdfundingUseCase(h.UserRepository, h.CrowdfundingRepository)
	res, err := createCrowdfunding.Execute(input, deposit, metadata)
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
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding created - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) CloseCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	return nil
}

func (h *CrowdfundingAdvanceHandlers) SettleCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	return nil
}

func (h *CrowdfundingAdvanceHandlers) UpdateCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.UpdateCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	updateCrowdfunding := crowdfunding_usecase.NewUpdateCrowdfundingUseCase(h.CrowdfundingRepository)
	res, err := updateCrowdfunding.Execute(input, metadata)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding updated - "), crowdfunding...))
	return nil
}
