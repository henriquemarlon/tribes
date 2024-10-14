package auction_usecase

import (
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateAuctionInputDTO struct {
	DebtIssued      custom_type.BigInt `json:"debt_issued"`
	MaxInterestRate custom_type.BigInt `json:"max_interest_rate"`
	ExpiresAt       int64              `json:"expires_at"`
	CreatedAt       int64              `json:"created_at"`
}

type CreateAuctionOutputDTO struct {
	Id              uint               `json:"id"`
	Creator         string             `json:"creator,omitempty"`
	DebtIssued      custom_type.BigInt `json:"debt_issued"`
	MaxInterestRate custom_type.BigInt `json:"max_interest_rate"`
	State           string             `json:"state"`
	ExpiresAt       int64              `json:"expires_at"`
	CreatedAt       int64              `json:"created_at"`
}

type CreateAuctionUseCase struct {
	UserRepository    entity.UserRepository
	AuctionRepository entity.AuctionRepository
}

func NewCreateAuctionUseCase(userRepository entity.UserRepository, auctionRepository entity.AuctionRepository) *CreateAuctionUseCase {
	return &CreateAuctionUseCase{
		UserRepository:    userRepository,
		AuctionRepository: auctionRepository,
	}
}

func (c *CreateAuctionUseCase) Execute(input *CreateAuctionInputDTO, metadata rollmelette.Metadata) (*CreateAuctionOutputDTO, error) {
	creator, err := c.UserRepository.FindUserByAddress(custom_type.NewAddress(metadata.MsgSender))
	if err != nil {
		return nil, err
	}
	auctions, err := c.AuctionRepository.FindAuctionsByCreator(creator.Username)
	if err != nil {
		return nil, err
	}
	for _, auction := range auctions {
		if auction.State == entity.AuctionOngoing || auction.State == entity.AuctionFinished {
			return nil, fmt.Errorf("creator already has an non paid auction")
		}
	}
	auction, err := entity.NewAuction(creator.Username, input.DebtIssued, input.MaxInterestRate, input.ExpiresAt, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.AuctionRepository.CreateAuction(auction)
	if err != nil {
		return nil, err
	}
	return &CreateAuctionOutputDTO{
		Id:              res.Id,
		Creator:         res.Creator,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		State:           string(res.State),
		ExpiresAt:       res.ExpiresAt,
		CreatedAt:       res.CreatedAt,
	}, nil
}
