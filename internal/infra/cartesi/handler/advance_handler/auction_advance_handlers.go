package advance_handler

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/auction_usecase"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
)

type AuctionAdvanceHandlers struct {
	BidRepository      entity.BidRepository
	UserRepository     entity.UserRepository
	AuctionRepository  entity.AuctionRepository
	ContractRepository entity.ContractRepository
}

func NewAuctionAdvanceHandlers(
	bidRepository entity.BidRepository,
	userRepository entity.UserRepository,
	auctionRepository entity.AuctionRepository,
	contractRepository entity.ContractRepository,
) *AuctionAdvanceHandlers {
	return &AuctionAdvanceHandlers{
		BidRepository:      bidRepository,
		UserRepository:     userRepository,
		AuctionRepository:  auctionRepository,
		ContractRepository: contractRepository,
	}
}

func (h *AuctionAdvanceHandlers) CreateAuctionHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *auction_usecase.CreateAuctionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createAuction := auction_usecase.NewCreateAuctionUseCase(h.UserRepository, h.AuctionRepository)
	res, err := createAuction.Execute(input, metadata)
	if err != nil {
		return err
	}
	auction, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("auction created - "), auction...))
	return nil
}

func (h *AuctionAdvanceHandlers) FinishAuctionHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *auction_usecase.FinishAuctionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	finishAuction := auction_usecase.NewFinishAuctionUseCase(h.AuctionRepository, h.UserRepository, h.BidRepository)
	res, err := finishAuction.Execute(input, metadata)
	if err != nil {
		return err
	}

	application, isDefined := env.AppAddress()
	if !isDefined {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}

	findUserByRole := user_usecase.NewFindUserByRoleUseCase(h.UserRepository)
	auctioneer, err := findUserByRole.Execute(&user_usecase.FindUserByRoleInputDTO{Role: "auctioneer"})
	if err != nil {
		return err
	}

	findUserByUsername := user_usecase.NewFindUserByUsernameUseCase(h.UserRepository)
	user, err := findUserByUsername.Execute(&user_usecase.FindUserByUsernameInputDTO{Username: res.Creator})
	if err != nil {
		return err
	}

	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{Symbol: "STABLECOIN"})
	if err != nil {
		return err
	}

	var amountRaised *big.Int = big.NewInt(0)

	for _, bid := range res.Bids {
		switch bid.State {
		case "accepted", "partially_accepted":
			if err := env.ERC20Transfer(stablecoin.Address.Address, auctioneer.Address.Address, user.Address.Address, bid.Amount.Int); err != nil {
				env.Report([]byte(err.Error()))
			}
			amountRaised.Add(amountRaised, bid.Amount.Int)
		case "rejected":
			if err := env.ERC20Transfer(stablecoin.Address.Address, auctioneer.Address.Address, bid.Bidder.Address, bid.Amount.Int); err != nil {
				env.Report([]byte(err.Error()))
			}
		}
	}

	tribesProfit := new(big.Int).Div(new(big.Int).Mul(amountRaised, big.NewInt(5)), big.NewInt(100))
	if err := env.ERC20Transfer(stablecoin.Address.Address, user.Address.Address, application, tribesProfit); err != nil {
		env.Report([]byte(err.Error()))
	}

	finishedAuction, err := json.Marshal(res)
	if err != nil {
		return err
	}

	env.Notice(append([]byte("auction finished - "), finishedAuction...))
	return nil
}
