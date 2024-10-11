package auction_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type CreateAuctionInputDTO struct {
	Creator      custom_type.Address `json:"creator,omitempty"`
	DebtIssued   custom_type.BigInt `json:"debt_issued"`
	InterestRate custom_type.BigInt `json:"interest_rate"`
	ExpiresAt    int64              `json:"expires_at"`
	CreatedAt    int64              `json:"created_at"`
}

type CreateAuctionOutputDTO struct {
	Id           uint               `json:"id"`
	Creator      custom_type.Address `json:"creator,omitempty"`
	DebtIssued   custom_type.BigInt `json:"debt_issued"`
	InterestRate custom_type.BigInt `json:"interest_rate"`
	State        string             `json:"state"`
	ExpiresAt    int64              `json:"expires_at"`
	CreatedAt    int64              `json:"created_at"`
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
	auction, err := entity.NewAuction(input.Creator, input.DebtIssued, input.InterestRate, input.ExpiresAt, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := c.DeviceRepository.CreateAuction(auction)
	if err != nil {
		return nil, err
	}
	return &CreateAuctionOutputDTO{
		Id:           res.Id,
		Creator:      res.Creator,
		DebtIssued:   res.DebtIssued,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		ExpiresAt:    res.ExpiresAt,
		CreatedAt:    res.CreatedAt,
	}, nil
}
