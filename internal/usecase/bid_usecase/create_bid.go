package bid_usecase

import (
	"fmt"

	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type CreateBidInputDTO struct {
	Price custom_type.BigInt `json:"interest_rate"`
}

type CreateBidOutputDTO struct {
	Id           uint                `json:"id"`
	AuctionId    uint                `json:"auction_id"`
	Bidder       custom_type.Address `json:"bidder"`
	Amount       custom_type.BigInt  `json:"amount"`
	InterestRate custom_type.BigInt  `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
}

type CreateBidUseCase struct {
	BidRepository      entity.BidRepository
	ContractRepository entity.ContractRepository
	AuctionRepository  entity.AuctionRepository
}

func NewCreateBidUseCase(bidRepository entity.BidRepository, contractRepository entity.ContractRepository, auctionRepository entity.AuctionRepository) *CreateBidUseCase {
	return &CreateBidUseCase{
		BidRepository:      bidRepository,
		ContractRepository: contractRepository,
		AuctionRepository:  auctionRepository,
	}
}

func (c *CreateBidUseCase) Execute(input *CreateBidInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateBidOutputDTO, error) {
	activeAuction, err := c.AuctionRepository.FindActiveAuction()
	if err != nil {
		return nil, err
	}
	if activeAuction == nil {
		return nil, fmt.Errorf("no active auction found, cannot create bid")
	}

	if metadata.BlockTimestamp > activeAuction.ExpiresAt {
		return nil, fmt.Errorf("active auction expired, cannot create bid")
	}
	volt, err := c.ContractRepository.FindContractBySymbol("VOLT")
	if err != nil {
		return nil, err
	}
	if deposit.(*rollmelette.ERC20Deposit).Token != volt.Address.Address {
		return nil, fmt.Errorf("invalid contract address provided for bid creation: %v", deposit.(*rollmelette.ERC20Deposit).Token)
	}

	if input.Price.Cmp(activeAuction.InterestRate.Int) == 1 {
		return nil, fmt.Errorf("bid price exceeds active auction price limit")
	}

	bid, err := entity.NewBid(activeAuction.Id, custom_type.NewAddress(deposit.(*rollmelette.ERC20Deposit).Sender), custom_type.NewBigInt(deposit.(*rollmelette.ERC20Deposit).Amount), input.Price, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.BidRepository.CreateBid(bid)
	if err != nil {
		return nil, err
	}

	return &CreateBidOutputDTO{
		Id:           res.Id,
		AuctionId:    res.AuctionId,
		Bidder:       res.Bidder,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
	}, nil
}
