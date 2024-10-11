package auction_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindActiveAuctionUseCase struct {
	AuctionRepository entity.AuctionRepository
}

func NewFindActiveAuctionUseCase(auctionRepository entity.AuctionRepository) *FindActiveAuctionUseCase {
	return &FindActiveAuctionUseCase{
		AuctionRepository: auctionRepository,
	}
}

func (f *FindActiveAuctionUseCase) Execute() (*FindAuctionOutputDTO, error) {
	res, err := f.AuctionRepository.FindActiveAuction()
	if err != nil {
		return nil, err
	}
	var bids []*FindAuctionOutputSubDTO
	for _, bid := range res.Bids {
		bids = append(bids, &FindAuctionOutputSubDTO{
			Id: bid.Id,

			AuctionId:    bid.AuctionId,
			Bidder:       bid.Bidder,
			Amount:       bid.Amount,
			InterestRate: bid.InterestRate,
			State:        string(bid.State),
			CreatedAt:    bid.CreatedAt,
			UpdatedAt:    bid.UpdatedAt,
		})
	}
	return &FindAuctionOutputDTO{
		Id:           res.Id,
		DebtIssued:   res.DebtIssued,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		Bids:         bids,
		ExpiresAt:    res.ExpiresAt,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}
