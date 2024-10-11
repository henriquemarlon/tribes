package auction_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type DeleteAuctionInputDTO struct {
	Id uint `json:"id"`
}

type DeleteAuctionUseCase struct {
	AuctionRepository entity.AuctionRepository
}

func NewDeleteAuctionUseCase(auctionRepository entity.AuctionRepository) *DeleteAuctionUseCase {
	return &DeleteAuctionUseCase{AuctionRepository: auctionRepository}
}

func (u *DeleteAuctionUseCase) Execute(input *DeleteAuctionInputDTO) error {
	return u.AuctionRepository.DeleteAuction(input.Id)
}
