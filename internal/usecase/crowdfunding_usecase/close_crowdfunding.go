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
	crowdfundings, err := u.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}

	var ongoingCrowdfunding *entity.Crowdfunding
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State == entity.CrowdfundingStateOngoing {
			ongoingCrowdfunding = crowdfunding
			break
		}
	}

	if ongoingCrowdfunding == nil {
		return nil, fmt.Errorf("no ongoing crowdfunding found, cannot close it")
	}

	// Ensure crowdfunding has expired before closing
	if metadata.BlockTimestamp < ongoingCrowdfunding.ExpiresAt {
		return nil, fmt.Errorf("crowdfunding not expired yet, you can't close it")
	}

	// Retrieve all orders related to the crowdfunding
	orders, err := u.OrderRepository.FindOrdersByCrowdfundingId(ctx, ongoingCrowdfunding.Id)
	if err != nil {
		return nil, err
	}

	// Sort orders by InterestRate ascending, Amount descending
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].InterestRate.Cmp(orders[j].InterestRate) == 0 {
			return orders[i].Amount.Cmp(orders[j].Amount) > 0 // larger amounts first
		}
		return orders[i].InterestRate.Cmp(orders[j].InterestRate) < 0
	})

	debtIssuedRemaining := new(uint256.Int).Set(ongoingCrowdfunding.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)

	for _, order := range orders {
		if debtIssuedRemaining.IsZero() {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}
			continue
		}

		if debtIssuedRemaining.Gt(order.Amount) || debtIssuedRemaining.Eq(order.Amount) {
			order.State = entity.OrderStateAccepted
			order.UpdatedAt = metadata.BlockTimestamp
			totalCollected.Add(totalCollected, order.Amount)

			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100)) // Interest = (amount * rate) / 100
			orderObligation := new(uint256.Int).Add(order.Amount, interest)
			totalObligation.Add(totalObligation, orderObligation)

			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.Sub(debtIssuedRemaining, order.Amount)
		} else {
			acceptedAmount := new(uint256.Int).Set(debtIssuedRemaining)
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptedAmount)

			order.Amount = acceptedAmount
			order.State = entity.OrderStatePartiallyAccepted
			order.UpdatedAt = metadata.BlockTimestamp
			totalCollected.Add(totalCollected, acceptedAmount)

			interest := new(uint256.Int).Mul(acceptedAmount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100)) // Interest = (amount * rate) / 100
			orderObligation := new(uint256.Int).Add(acceptedAmount, interest)
			totalObligation.Add(totalObligation, orderObligation)

			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}

			_, err = u.OrderRepository.CreateOrder(ctx, &entity.Order{
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         rejectedAmount,
				InterestRate:   order.InterestRate,
				State:          entity.OrderStateRejected,
				CreatedAt:      metadata.BlockTimestamp,
				UpdatedAt:      metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.Clear()
		}
	}

	twoThirdsTarget := new(uint256.Int).Mul(ongoingCrowdfunding.DebtIssued, uint256.NewInt(2))
	twoThirdsTarget.Div(twoThirdsTarget, uint256.NewInt(3))
	if totalCollected.Lt(twoThirdsTarget) {
		// Cancel crowdfunding and mark all orders as rejected
		for _, order := range orders {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}
		}

		ongoingCrowdfunding.State = entity.CrowdfundingStateCanceled
		ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
		_, err := u.CrowdfundingRepository.UpdateCrowdfunding(ctx, ongoingCrowdfunding)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("crowdfunding canceled due to insufficient funds collected")
	}

	ongoingCrowdfunding.State = entity.CrowdfundingStateClosed
	ongoingCrowdfunding.TotalObligation = totalObligation
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
		TotalObligation: res.TotalObligation,
		State:           string(res.State),
		Orders:          res.Orders,
		ExpiresAt:       res.ExpiresAt,
		MaturityAt:      res.MaturityAt,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}
