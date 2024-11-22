package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
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
	var input *crowdfunding_usecase.CloseCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	closeCrowdfunding := crowdfunding_usecase.NewCloseCrowdfundingUseCase(h.CrowdfundingRepository, h.UserRepository, h.OrderRepository)
	res, err := closeCrowdfunding.Execute(input, metadata)
	if err != nil {
		return err
	}

	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{
		Symbol: "STABLECOIN",
	})
	if err != nil {
		return err
	}
	// TODO: remove this check when update to V2
	appAddress, isSet := env.AppAddress()
	if !isSet {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}

	// Return the funds to the investors who had their orders rejected
	for _, order := range res.Orders {
		if order.State == entity.OrderStateRejected {
			if err = env.ERC20Transfer(
				stablecoin.Address,
				appAddress,
				order.Investor,
				order.Amount.ToBig(),
			); err != nil {
				return err
			}
		}
	}

	// Transfer the raised funds to the creator
	if err = env.ERC20Transfer(
		stablecoin.Address,
		appAddress,
		res.Creator,
		res.DebtIssued.ToBig(),
	); err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding closed - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) SettleCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *crowdfunding_usecase.SettleCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	settleCrowdfunding := crowdfunding_usecase.NewSettleCrowdfundingUseCase(h.UserRepository, h.CrowdfundingRepository, h.ContractRepository)
	res, err := settleCrowdfunding.Execute(input, deposit, metadata)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	contract, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{
		Symbol: "STABLECOIN",
	})
	if err != nil {
		return err
	}
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettled {
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			if err := env.ERC20Transfer(
				contract.Address,
				res.Creator,
				order.Investor,
				new(uint256.Int).Add(order.Amount, interest).ToBig(),
			); err != nil {
				return err
			}
		}
	}
	env.Notice(append([]byte("crowdfunding settled - "), crowdfunding...))
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

func (h *CrowdfundingAdvanceHandlers) DeleteCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *crowdfunding_usecase.DeleteCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	deleteCrowdfunding := crowdfunding_usecase.NewDeleteCrowdfundingUseCase(h.CrowdfundingRepository)
	err := deleteCrowdfunding.Execute(input)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding deleted - "), crowdfunding...))
	return nil
}