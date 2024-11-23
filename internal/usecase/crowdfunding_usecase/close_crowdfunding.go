package crowdfunding_usecase

import (
	"context"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type CloseCrowdfundingInputDTO struct {
	Creator common.Address `json:"creator"`
}

type CloseCrowdfundingOutputDTO struct {
	Id              uint            `json:"id"`
	Creator         common.Address  `json:"creator,omitempty"`
	DebtIssued      *uint256.Int    `json:"debt_issued,omitempty"`
	MaxInterestRate *uint256.Int    `json:"max_interest_rate,omitempty"`
	TotalObligation *uint256.Int    `json:"total_obligation,omitempty"`
	State           string          `json:"state,omitempty"`
	Orders          []*entity.Order `json:"orders,omitempty"`
	ExpiresAt       int64           `json:"expires_at,omitempty"`
	MaturityAt      int64           `json:"maturity_at,omitempty"`
	CreatedAt       int64           `json:"created_at,omitempty"`
	UpdatedAt       int64           `json:"updated_at,omitempty"`
}

type CloseCrowdfundingUseCase struct {
	OrderRepository        entity.OrderRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCloseCrowdfundingUseCase(crowdfundingRepository entity.CrowdfundingRepository, orderRepository entity.OrderRepository) *CloseCrowdfundingUseCase {
	return &CloseCrowdfundingUseCase{
		OrderRepository:        orderRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (u *CloseCrowdfundingUseCase) Execute(ctx context.Context, input *CloseCrowdfundingInputDTO, metadata rollmelette.Metadata) (*CloseCrowdfundingOutputDTO, error) {
	// Retrieve ongoing crowdfunding by creator
	crowdfundings, err := u.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}

	var ongoingCrowdfunding *entity.Crowdfunding
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State == entity.CrowdfundingState("ongoing") {
			ongoingCrowdfunding = crowdfunding
			break
		}
	}

	if ongoingCrowdfunding == nil {
		return nil, fmt.Errorf("no ongoing crowdfunding found, cannot finish crowdfunding")
	}

	// Ensure crowdfunding has expired before closing
	if metadata.BlockTimestamp < ongoingCrowdfunding.ExpiresAt {
		return nil, fmt.Errorf("active crowdfunding not expired, you can't finish it yet")
	}

	// Retrieve all orders related to the crowdfunding
	orders, err := u.OrderRepository.FindOrdersByCrowdfundingId(ctx, ongoingCrowdfunding.Id)
	if err != nil {
		return nil, err
	}

	// Sort orders by the best interest-to-amount ratio (ascending)
	sort.Slice(orders, func(i, j int) bool {
		ratioI := new(uint256.Int).Div(orders[i].InterestRate, orders[i].Amount)
		ratioJ := new(uint256.Int).Div(orders[j].InterestRate, orders[j].Amount)
		return ratioI.Cmp(ratioJ) < 0
	})

	debtIssuedRemaining := new(uint256.Int).Set(ongoingCrowdfunding.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0) // Initialize the total obligation

	// Process each order
	for _, order := range orders {
		if debtIssuedRemaining.IsZero() {
			// Mark remaining orders as rejected
			order.State = "rejected"
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}
			continue
		}

		if debtIssuedRemaining.Gt(order.Amount) {
			// Fully accept the order
			order.State = "accepted"
			order.UpdatedAt = metadata.BlockTimestamp
			totalCollected.Add(totalCollected, order.Amount)

			// Calculate interest and add to total obligation
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100)) // Interest = (amount * rate) / 100
			totalObligation.Add(totalObligation, new(uint256.Int).Add(order.Amount, interest))

			// Update the accepted order in the database
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.Sub(debtIssuedRemaining, order.Amount)
		} else {
			// Partially accept the order
			acceptedAmount := new(uint256.Int).Set(debtIssuedRemaining)
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptedAmount)

			// Update the order with the accepted portion
			order.Amount = acceptedAmount
			order.State = "partially_accepted"
			order.UpdatedAt = metadata.BlockTimestamp
			totalCollected.Add(totalCollected, acceptedAmount)

			// Calculate interest for the accepted portion and add to total obligation
			interest := new(uint256.Int).Mul(acceptedAmount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100)) // Interest = (amount * rate) / 100
			totalObligation.Add(totalObligation, new(uint256.Int).Add(acceptedAmount, interest))

			// Update the partially accepted order in the database
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}

			// Create a new order for the rejected portion
			_, err = u.OrderRepository.CreateOrder(ctx, &entity.Order{
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         rejectedAmount,
				InterestRate:   order.InterestRate,
				State:          "rejected",
				CreatedAt:      metadata.BlockTimestamp,
				UpdatedAt:      metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.Clear()
		}
	}

	// Check if total collected meets the minimum threshold (2/3 of DebtIssued)
	twoThirdsTarget := new(uint256.Int).Mul(ongoingCrowdfunding.DebtIssued, uint256.NewInt(2)).Div(ongoingCrowdfunding.DebtIssued, uint256.NewInt(3))
	if totalCollected.Lt(twoThirdsTarget) {
		// Cancel crowdfunding and mark all orders as rejected
		for _, order := range orders {
			order.State = "rejected"
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}
		}

		ongoingCrowdfunding.State = entity.CrowdfundingState("canceled")
		ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
		res, err := u.CrowdfundingRepository.UpdateCrowdfunding(ctx, ongoingCrowdfunding)
		if err != nil {
			return nil, err
		}

		return &CloseCrowdfundingOutputDTO{
			Id:              res.Id,
			Creator:         res.Creator,
			DebtIssued:      res.DebtIssued,
			MaxInterestRate: res.MaxInterestRate,
			State:           string(res.State),
			Orders:          orders,
			ExpiresAt:       res.ExpiresAt,
			MaturityAt:      res.MaturityAt,
			CreatedAt:       res.CreatedAt,
			UpdatedAt:       res.UpdatedAt,
		}, nil
	}

	// Close crowdfunding if threshold is met
	ongoingCrowdfunding.State = entity.CrowdfundingState("closed")
	ongoingCrowdfunding.TotalObligation = totalObligation
	ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
	res, err := u.CrowdfundingRepository.UpdateCrowdfunding(ctx, ongoingCrowdfunding)
	if err != nil {
		return nil, err
	}

	// Return the final state of the crowdfunding
	return &CloseCrowdfundingOutputDTO{
		Id:              res.Id,
		Creator:         res.Creator,
		DebtIssued:      res.DebtIssued,
		MaxInterestRate: res.MaxInterestRate,
		TotalObligation: res.TotalObligation,
		State:           string(res.State),
		Orders:          orders,
		ExpiresAt:       res.ExpiresAt,
		MaturityAt:      res.MaturityAt,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}
