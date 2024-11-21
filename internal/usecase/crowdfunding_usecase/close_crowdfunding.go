package crowdfunding_usecase

import (
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
	State           string          `json:"state,omitempty"`
	Orders          []*entity.Order `json:"orders,omitempty"`
	ExpiresAt       int64           `json:"expires_at,omitempty"`
	CreatedAt       int64           `json:"created_at,omitempty"`
	UpdatedAt       int64           `json:"updated_at,omitempty"`
}

type CloseCrowdfundingUseCase struct {
	OrderRepository        entity.OrderRepository
	UserRepository         entity.UserRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewCloseCrowdfundingUseCase(crowdfundingRepository entity.CrowdfundingRepository, userRepository entity.UserRepository, orderRepository entity.OrderRepository) *CloseCrowdfundingUseCase {
	return &CloseCrowdfundingUseCase{
		OrderRepository:        orderRepository,
		UserRepository:         userRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (u *CloseCrowdfundingUseCase) Execute(input *CloseCrowdfundingInputDTO, metadata rollmelette.Metadata) (*CloseCrowdfundingOutputDTO, error) {
	crowdfundings, err := u.CrowdfundingRepository.FindCrowdfundingsByCreator(input.Creator)
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

	if metadata.BlockTimestamp < ongoingCrowdfunding.ExpiresAt {
		return nil, fmt.Errorf("active crowdfunding not expired, you can't finish it yet")
	}

	orders, err := u.OrderRepository.FindOrdersByCrowdfundingId(ongoingCrowdfunding.Id)
	if err != nil {
		return nil, err
	}

	totalCollected := uint256.NewInt(0)
	for _, order := range orders {
		totalCollected.Add(totalCollected, order.Amount)
	}

	// According to CVM Resolution 88
	twoThirdsTarget := new(uint256.Int).Mul(ongoingCrowdfunding.DebtIssued, uint256.NewInt(2)).Div(ongoingCrowdfunding.DebtIssued, uint256.NewInt(3))

	if totalCollected.Lt(twoThirdsTarget) {
		ongoingCrowdfunding.State = entity.CrowdfundingState("canceled")
		ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
		res, err := u.CrowdfundingRepository.UpdateCrowdfunding(ongoingCrowdfunding)
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
			CreatedAt:       res.CreatedAt,
			UpdatedAt:       res.UpdatedAt,
		}, nil
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].InterestRate.Cmp(orders[j].InterestRate) < 0
	})

	debtIssuedRemaining := new(uint256.Int).Set(ongoingCrowdfunding.DebtIssued)

	for _, order := range orders {
		if debtIssuedRemaining.IsZero() {
			order.State = "rejected"
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(order)
			if err != nil {
				return nil, err
			}
			continue
		}

		if debtIssuedRemaining.Gt(order.Amount) {
			order.State = "accepted"
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(order)
			if err != nil {
				return nil, err
			}
			debtIssuedRemaining.Sub(debtIssuedRemaining, order.Amount)
		} else {
			partiallyAcceptedAmount := new(uint256.Int).Set(debtIssuedRemaining)
			_, err := u.OrderRepository.CreateOrder(&entity.Order{
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         partiallyAcceptedAmount,
				InterestRate:   order.InterestRate,
				State:          "partially_accepted",
				CreatedAt:      metadata.BlockTimestamp,
				UpdatedAt:      metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}

			rejectedAmount := new(uint256.Int).Sub(order.Amount, partiallyAcceptedAmount)
			_, err = u.OrderRepository.CreateOrder(&entity.Order{
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

			err = u.OrderRepository.DeleteOrder(order.Id)
			if err != nil {
				return nil, err
			}
			debtIssuedRemaining.Clear()
		}
	}

	totalObligation := new(uint256.Int)
	for _, order := range orders {
		interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
		interest.Div(interest, uint256.NewInt(100))
		obligation := new(uint256.Int).Add(order.Amount, interest)
		totalObligation.Add(totalObligation, obligation)
	}

	ongoingCrowdfunding.TotalObligation = totalObligation
	ongoingCrowdfunding.State = entity.CrowdfundingState("closed")
	ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
	res, err := u.CrowdfundingRepository.UpdateCrowdfunding(ongoingCrowdfunding)
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
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}, nil
}
