package bid_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindBidsByAuctionIdInputDTO struct {
	AuctionId uint `json:"auction_id"`
}

type FindBidsByAuctionIdOutputDTO []*FindBidOutputDTO

type FindBidsByAuctionIdUseCase struct {
	BidRepository entity.BidRepository
}

func NewFindBidsByAuctionIdUseCase(bidRepository entity.BidRepository) *FindBidsByAuctionIdUseCase {
	return &FindBidsByAuctionIdUseCase{
		BidRepository: bidRepository,
	}
}

func (c *FindBidsByAuctionIdUseCase) Execute(input *FindBidsByAuctionIdInputDTO) (*FindBidsByAuctionIdOutputDTO, error) {
	res, err := c.BidRepository.FindBidsByAuctionId(input.AuctionId)
	if err != nil {
		return nil, err
	}
	output := make(FindBidsByAuctionIdOutputDTO, len(res))
	for i, bid := range res {
		output[i] = &FindBidOutputDTO{
			Id:           bid.Id,
			AuctionId:    bid.AuctionId,
			Bidder:       bid.Bidder,
			Amount:       bid.Amount,
			InterestRate: bid.InterestRate,
			State:        string(bid.State),
			CreatedAt:    bid.CreatedAt,
			UpdatedAt:    bid.UpdatedAt,
		}
	}
	return &output, nil
}
