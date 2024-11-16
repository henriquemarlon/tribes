package advance_handler

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	"github.com/tribeshq/tribes/internal/usecase/crowdfunding_usecase"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
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
	res, err := createCrowdfunding.Execute(input, metadata)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding created - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) FinishCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *crowdfunding_usecase.FinishCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	finishCrowdfunding := crowdfunding_usecase.NewFinishCrowdfundingUseCase(h.CrowdfundingRepository, h.UserRepository, h.OrderRepository)
	res, err := finishCrowdfunding.Execute(input, metadata)
	if err != nil {
		return err
	}

	application, isDefined := env.AppAddress()
	if !isDefined {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}

	findUserByRole := user_usecase.NewFindUserByRoleUseCase(h.UserRepository)
	crowdfundingeer, err := findUserByRole.Execute(&user_usecase.FindUserByRoleInputDTO{Role: "crowdfundingeer"})
	if err != nil {
		return err
	}

	findCrowdfundingBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findCrowdfundingBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{Symbol: "STABLECOIN"})
	if err != nil {
		return err
	}

	var amountRaised *big.Int = big.NewInt(0)

	for _, order := range res.Orders {
		switch order.State {
		case "accepted", "partially_accepted":
			if err := env.ERC20Transfer(stablecoin.Address, crowdfundingeer.Address, res.Creator, order.Amount.ToBig()); err != nil {
				env.Report([]byte(err.Error()))
			}
			amountRaised.Add(amountRaised, order.Amount.ToBig())
		case "rejected":
			if err := env.ERC20Transfer(stablecoin.Address, crowdfundingeer.Address, order.Investor, order.Amount.ToBig()); err != nil {
				env.Report([]byte(err.Error()))
			}
		}
	}

	tribesProfit := new(big.Int).Div(new(big.Int).Mul(amountRaised, big.NewInt(5)), big.NewInt(100))
	if err := env.ERC20Transfer(stablecoin.Address, res.Creator, application, tribesProfit); err != nil {
		env.Report([]byte(err.Error()))
	}

	finishedCrowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}

	env.Notice(append([]byte("crowdfunding finished - "), finishedCrowdfunding...))
	return nil
}
