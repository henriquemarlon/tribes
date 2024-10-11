package bid_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindBidsByStateInputDTO struct {
	AuctionId uint   `json:"auction_id"`
	State     string `json:"state"`
}

type FindBidsByStateOutputDTO []*FindBidOutputDTO

type FindBidsByStateUseCase struct {
	BidRepository entity.BidRepository
}

func NewFindBidsByStateUseCase(bidRepository entity.BidRepository) *FindBidsByStateUseCase {
	return &FindBidsByStateUseCase{
		BidRepository: bidRepository,
	}
}

func (f *FindBidsByStateUseCase) Execute(input *FindBidsByStateInputDTO) (FindBidsByStateOutputDTO, error) {
	res, err := f.BidRepository.FindBidsByState(input.AuctionId, input.State)
	if err != nil {
		return nil, err
	}
	output := make(FindBidsByStateOutputDTO, len(res))
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
	return output, nil
}
