package bid_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type DeleteBidInputDTO struct {
	Id uint `json:"id"`
}

type DeleteBidUseCase struct {
	BidRepository entity.BidRepository
}

func NewDeleteBidUseCase(bidRepository entity.BidRepository) *DeleteBidUseCase {
	return &DeleteBidUseCase{
		BidRepository: bidRepository,
	}
}

func (c *DeleteBidUseCase) Execute(input *DeleteBidInputDTO) error {
	return c.BidRepository.DeleteBid(input.Id)
}
