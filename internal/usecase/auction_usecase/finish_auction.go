package auction_usecase

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FinishAuctionInputDTO struct {
	Creator string `json:"creator"`
}

type FinishAuctionOutputDTO struct {
	Id              uint               `json:"id"`
	Creator         string             `json:"creator,omitempty"`
	DebtIssued      custom_type.BigInt `json:"debt_issued,omitempty"`
	MaxInterestRate custom_type.BigInt `json:"max_interest_rate,omitempty"`
	State           string             `json:"state,omitempty"`
	Bids            []*entity.Bid      `json:"bids,omitempty"`
	ExpiresAt       int64              `json:"expires_at,omitempty"`
	CreatedAt       int64              `json:"created_at,omitempty"`
	UpdatedAt       int64              `json:"updated_at,omitempty"`
}

type FinishAuctionUseCase struct {
	BidRepository     entity.BidRepository
	UserRepository    entity.UserRepository
	AuctionRepository entity.AuctionRepository
}

func NewFinishAuctionUseCase(auctionRepository entity.AuctionRepository, userRepository entity.UserRepository, bidRepository entity.BidRepository) *FinishAuctionUseCase {
	return &FinishAuctionUseCase{
		BidRepository:     bidRepository,
		UserRepository:    userRepository,
		AuctionRepository: auctionRepository,
	}
}

func (u *FinishAuctionUseCase) Execute(input *FinishAuctionInputDTO, metadata rollmelette.Metadata) (*FinishAuctionOutputDTO, error) {
	auctions, err := u.AuctionRepository.FindAuctionByStateFromCreator(input.Creator, string(entity.AuctionState("ongoing")))
	if err != nil {
		return nil, err
	}
	if len(auctions) == 0 {
		return nil, fmt.Errorf("no active auction found, cannot finish auction")
	}
	activeAuction := auctions[0]

	if metadata.BlockTimestamp < activeAuction.ExpiresAt {
		return nil, fmt.Errorf("active auction not expired, you can't finish it yet")
	}

	bids, err := u.BidRepository.FindBidsByAuctionId(activeAuction.Id)
	if err != nil {
		return nil, err
	}

	if len(bids) == 0 {
		activeAuction.State = entity.AuctionState("canceled")
		activeAuction.UpdatedAt = metadata.BlockTimestamp
		res, err := u.AuctionRepository.UpdateAuction(activeAuction)
		if err != nil {
			return nil, err
		}
		return &FinishAuctionOutputDTO{
			Id:              res.Id,
			Creator:         res.Creator,
			DebtIssued:      res.DebtIssued,
			MaxInterestRate: res.MaxInterestRate,
			State:           string(res.State),
			Bids:            bids,
			ExpiresAt:       res.ExpiresAt,
			CreatedAt:       res.CreatedAt,
			UpdatedAt:       res.UpdatedAt,
		}, nil
	}

	sort.Slice(bids, func(i, j int) bool {
		return bids[i].InterestRate.Cmp(bids[j].InterestRate.Int) < 0
	})

	debtIssuedRemaining := new(big.Int).Set(activeAuction.DebtIssued.Int)

	for _, bid := range bids {
		if debtIssuedRemaining.Sign() == 0 {
			bid.State = "rejected"
			bid.UpdatedAt = metadata.BlockTimestamp
			_, err := u.BidRepository.UpdateBid(bid)
			if err != nil {
				return nil, err
			}
			continue
		}

		if debtIssuedRemaining.Cmp(bid.Amount.Int) >= 0 {
			bid.State = "accepted"
			bid.UpdatedAt = metadata.BlockTimestamp
			_, err := u.BidRepository.UpdateBid(bid)
			if err != nil {
				return nil, err
			}
			debtIssuedRemaining.Sub(debtIssuedRemaining, bid.Amount.Int)
		} else {
			// Partially accept the bid
			partiallyAcceptedAmount := new(big.Int).Set(debtIssuedRemaining)
			_, err := u.BidRepository.CreateBid(&entity.Bid{
				AuctionId:    bid.AuctionId,
				Bidder:       bid.Bidder,
				Amount:       custom_type.NewBigInt(partiallyAcceptedAmount),
				InterestRate: bid.InterestRate,
				State:        "partially_accepted",
				CreatedAt:    metadata.BlockTimestamp,
				UpdatedAt:    metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}

			// Reject the remaining amount
			rejectedAmount := new(big.Int).Sub(bid.Amount.Int, partiallyAcceptedAmount)
			_, err = u.BidRepository.CreateBid(&entity.Bid{
				AuctionId:    bid.AuctionId,
				Bidder:       bid.Bidder,
				Amount:       custom_type.NewBigInt(rejectedAmount),
				InterestRate: bid.InterestRate,
				State:        "rejected",
				CreatedAt:    metadata.BlockTimestamp,
				UpdatedAt:    metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}

			// Delete original bid
			err = u.BidRepository.DeleteBid(bid.Id)
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.SetInt64(0)
		}
	}

	activeAuction.State = entity.AuctionState("finished")
	activeAuction.UpdatedAt = metadata.BlockTimestamp
	res, err := u.AuctionRepository.UpdateAuction(activeAuction)
	if err != nil {
		return nil, err
	}

	return &FinishAuctionOutputDTO{
		Id:              res.Id,
		Creator:         res.Creator,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		State:           string(res.State),
		Bids:            bids,
		ExpiresAt:       res.ExpiresAt,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}
