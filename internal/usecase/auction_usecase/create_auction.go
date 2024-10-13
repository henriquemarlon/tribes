package auction_usecase

import (
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateAuctionInputDTO struct {
	DebtIssued      custom_type.BigInt  `json:"debt_issued"`
	MaxInterestRate custom_type.BigInt  `json:"max_interest_rate"`
	ExpiresAt       int64               `json:"expires_at"`
	CreatedAt       int64               `json:"created_at"`
}

type CreateAuctionOutputDTO struct {
	Id              uint                `json:"id"`
	Creator         custom_type.Address `json:"creator,omitempty"`
	DebtIssued      custom_type.BigInt  `json:"debt_issued"`
	MaxInterestRate custom_type.BigInt  `json:"max_interest_rate"`
	State           string              `json:"state"`
	ExpiresAt       int64               `json:"expires_at"`
	CreatedAt       int64               `json:"created_at"`
}

type CreateAuctionUseCase struct {
	DeviceRepository entity.AuctionRepository
}

func NewCreateAuctionUseCase(deviceRepository entity.AuctionRepository) *CreateAuctionUseCase {
	return &CreateAuctionUseCase{
		DeviceRepository: deviceRepository,
	}
}

func (c *CreateAuctionUseCase) Execute(input *CreateAuctionInputDTO, metadata rollmelette.Metadata) (*CreateAuctionOutputDTO, error) {
	auction, err := entity.NewAuction(custom_type.NewAddress(metadata.MsgSender), input.DebtIssued, input.MaxInterestRate, input.ExpiresAt, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.DeviceRepository.CreateAuction(auction)
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
