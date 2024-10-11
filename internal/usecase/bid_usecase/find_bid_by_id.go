package bid_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindBidByIdInputDTO struct {
	Id uint `json:"id"`
}

type FindBidByIdUseCase struct {
	BidRepository entity.BidRepository
}

func NewFindBidByIdUseCase(bidRepository entity.BidRepository) *FindBidByIdUseCase {
	return &FindBidByIdUseCase{
		BidRepository: bidRepository,
	}
}

func (c *FindBidByIdUseCase) Execute(input *FindBidByIdInputDTO) (*FindBidOutputDTO, error) {
	res, err := c.BidRepository.FindBidById(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindBidOutputDTO{
		Id:           res.Id,
		AuctionId:    res.AuctionId,
		Bidder:       res.Bidder,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}
