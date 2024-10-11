package bid_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindAllBidsOutputDTO []*FindBidOutputDTO

type FindAllBidsUseCase struct {
	BidRepository entity.BidRepository
}

func NewFindAllBidsUseCase(bidRepository entity.BidRepository) *FindAllBidsUseCase {
	return &FindAllBidsUseCase{
		BidRepository: bidRepository,
	}
}

func (f *FindAllBidsUseCase) Execute() (*FindAllBidsOutputDTO, error) {
	res, err := f.BidRepository.FindAllBids()
	if err != nil {
		return nil, err
	}
	output := make(FindAllBidsOutputDTO, len(res))
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
