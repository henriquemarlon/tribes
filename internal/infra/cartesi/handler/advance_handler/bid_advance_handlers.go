package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/bid_usecase"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
)

type BidAdvanceHandlers struct {
	BidRepository      entity.BidRepository
	UserRepository     entity.UserRepository
	AuctionRepository  entity.AuctionRepository
	ContractRepository entity.ContractRepository
}

func NewBidAdvanceHandlers(bidRepository entity.BidRepository, userRepository entity.UserRepository, contractRepository entity.ContractRepository, auctionRepository entity.AuctionRepository) *BidAdvanceHandlers {
	return &BidAdvanceHandlers{
		BidRepository:      bidRepository,
		UserRepository:     userRepository,
		AuctionRepository:  auctionRepository,
		ContractRepository: contractRepository,
	}
}

func (h *BidAdvanceHandlers) CreateBidHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	switch deposit := deposit.(type) {
	case *rollmelette.ERC20Deposit:
		var input bid_usecase.CreateBidInputDTO
		if err := json.Unmarshal(payload, &input); err != nil {
			return err
		}
		createBid := bid_usecase.NewCreateBidUseCase(h.BidRepository, h.ContractRepository, h.AuctionRepository)
		res, err := createBid.Execute(&input, deposit, metadata)
		if err != nil {
			return err
		}

		findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
		stablecoin, err := findContractBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{Symbol: "STABLECOIN"})
		if err != nil {
			return err
		}
		findUserByRole := user_usecase.NewFindUserByRoleUseCase(h.UserRepository)
		auctioneer, err := findUserByRole.Execute(&user_usecase.FindUserByRoleInputDTO{Role: "auctioneer"})
		if err != nil {
			return err
		}

		if err := env.ERC20Transfer(stablecoin.Address.Address, res.Bidder.Address, auctioneer.Address.Address, res.Amount.Int); err != nil {
			return err
		}
		bid, err := json.Marshal(res)
		if err != nil {
			return err
		}
		env.Notice(append([]byte("created bid - "), bid...))
		return nil
	default:
		return fmt.Errorf("unsupported deposit type")
	}
}
