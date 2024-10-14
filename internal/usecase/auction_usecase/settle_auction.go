package auction_usecase

// import (
// 	"fmt"

// 	"github.com/rollmelette/rollmelette"
// 	"github.com/tribeshq/tribes/internal/domain/entity"
// 	"github.com/tribeshq/tribes/pkg/custom_type"
// )

// type SettleAuctionOutputDTO struct {
// 	Id                  uint               `json:"id"`
// 	Creator             string             `json:"creator,omitempty"`
// 	DebtIssued          custom_type.BigInt `json:"debt_issued,omitempty"`
// 	MaxInterestRate     custom_type.BigInt `json:"max_interest_rate,omitempty"`
// 	AverageInterestRate custom_type.BigInt `json:"average_interest_rate,omitempty"`
// 	State               string             `json:"state,omitempty"`
// 	Bids                []*entity.Bid      `json:"bids,omitempty"`
// 	ExpiresAt           int64              `json:"expires_at,omitempty"`
// 	CreatedAt           int64              `json:"created_at,omitempty"`
// 	UpdatedAt           int64              `json:"updated_at,omitempty"`
// }

// type SettleAuctionUseCase struct {
// 	BidRepository     entity.BidRepository
// 	UserRepository    entity.UserRepository
// 	AuctionRepository entity.AuctionRepository
// }

// func NewSettleAuctionUseCase(auctionRepository entity.AuctionRepository, userRepository entity.UserRepository, bidRepository entity.BidRepository) *SettleAuctionUseCase {
// 	return &SettleAuctionUseCase{
// 		BidRepository:     bidRepository,
// 		UserRepository:    userRepository,
// 		AuctionRepository: auctionRepository,
// 	}
// }

// func (s *SettleAuctionUseCase) Execute(metadata rollmelette.Metadata) (*SettleAuctionOutputDTO, error) {
// 	creator, err := s.UserRepository.FindUserByAddress(custom_type.NewAddress(metadata.MsgSender))
// 	if err != nil {
// 		return nil, err
// 	}

// 	auctions, err := s.AuctionRepository.FindAuctionByStateFromCreator(creator.Username, string(entity.AuctionFinished))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(auctions) == 0 {
// 		return nil, fmt.Errorf("no finished auction found")
// 	}
// 	auction := auctions[0]

// 	bids, err := s.BidRepository.FindBidsByAuctionId(auction.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var bidsToBePaid []*entity.Bid
// 	var totalInterestRate *big.Float = big.NewFloat(0)
// 	var bidCount int

// 	for _, bid := range bids {
// 		if bid.State == entity.BidState("accepted") || bid.State == entity.BidState("partially_accepted") {
// 			bid.State = entity.BidState("paid")
// 			bidsToBePaid = append(bidsToBePaid, bid)

// 			// Convert bid interest rate to big.Float and accumulate
// 			bidInterestRate := new(big.Float).SetInt(bid.InterestRate.Int()) // Assuming InterestRate is a *big.Int
// 			totalInterestRate = new(big.Float).Add(totalInterestRate, bidInterestRate)
// 			bidCount++
// 		}
// 	}

// 	var averageInterestRate *big.Float
// 	if bidCount > 0 {
// 		averageInterestRate = new(big.Float).Quo(totalInterestRate, big.NewFloat(float64(bidCount)))
// 	} else {
// 		averageInterestRate = big.NewFloat(0)
// 	}

// 	return &SettleAuctionOutputDTO{
// 		Id:              auction.Id,
// 		Creator:         auction.Creator,
// 		DebtIssued:      auction.DebtIssued,
// 		MaxInterestRate: auction.MaxInterestRate,
// 		State:           string(auction.State),
// 		Bids:            bidsToBePaid,
// 		ExpiresAt:       auction.ExpiresAt,
// 		CreatedAt:       auction.CreatedAt,
// 		UpdatedAt:       auction.UpdatedAt,
// 	}, nil
// }
