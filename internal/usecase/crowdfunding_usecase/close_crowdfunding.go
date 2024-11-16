package crowdfunding_usecase

import (
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FinishCrowdfundingInputDTO struct {
	Creator common.Address `json:"creator"`
}

type FinishCrowdfundingOutputDTO struct {
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

type FinishCrowdfundingUseCase struct {
	OrderRepository        entity.OrderRepository
	UserRepository         entity.UserRepository
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewFinishCrowdfundingUseCase(crowdfundingRepository entity.CrowdfundingRepository, userRepository entity.UserRepository, orderRepository entity.OrderRepository) *FinishCrowdfundingUseCase {
	return &FinishCrowdfundingUseCase{
		OrderRepository:        orderRepository,
		UserRepository:         userRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (u *FinishCrowdfundingUseCase) Execute(input *FinishCrowdfundingInputDTO, metadata rollmelette.Metadata) (*FinishCrowdfundingOutputDTO, error) {
	crowdfundings, err := u.CrowdfundingRepository.FindCrowdfundingsByCreator(input.Creator)
	if err != nil {
		return nil, err
	}
	if len(crowdfundings) == 0 {
		return nil, fmt.Errorf("no active crowdfunding found, cannot finish crowdfunding")
	}
	activeCrowdfunding := crowdfundings[0]

	if metadata.BlockTimestamp < activeCrowdfunding.ExpiresAt {
		return nil, fmt.Errorf("active crowdfunding not expired, you can't finish it yet")
	}

	orders, err := u.OrderRepository.FindOrdersByCrowdfundingId(activeCrowdfunding.Id)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		activeCrowdfunding.State = entity.CrowdfundingState("canceled")
		activeCrowdfunding.UpdatedAt = metadata.BlockTimestamp
		res, err := u.CrowdfundingRepository.UpdateCrowdfunding(activeCrowdfunding)
		if err != nil {
			return nil, err
		}
		return &FinishCrowdfundingOutputDTO{
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

	debtIssuedRemaining := new(uint256.Int).Set(activeCrowdfunding.DebtIssued)

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

	activeCrowdfunding.State = entity.CrowdfundingState("finished")
	activeCrowdfunding.UpdatedAt = metadata.BlockTimestamp
	res, err := u.CrowdfundingRepository.UpdateCrowdfunding(activeCrowdfunding)
	if err != nil {
		return nil, err
	}

	return &FinishCrowdfundingOutputDTO{
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
