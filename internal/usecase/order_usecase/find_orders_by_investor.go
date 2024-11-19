package order_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type FinsOrdersByInvestorInputDTO struct {
	Investor common.Address `json:"investor"`
}

type FindOrdersByInvestorOutputDTO []*FindOrderOutputDTO

type FindOrderByInvestorUsecase struct {
	OrderRepository entity.OrderRepository
}

func NewFindOrderByInvestorUsecase(orderRepository entity.OrderRepository) *FindOrderByInvestorUsecase {
	return &FindOrderByInvestorUsecase{
		OrderRepository: orderRepository,
	}
}

func (o *FindOrderByInvestorUsecase) Execute(input *FinsOrdersByInvestorInputDTO) (FindOrdersByInvestorOutputDTO, error) {
	res, err := o.OrderRepository.FindOrdersByInvestor(input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByInvestorOutputDTO, len(res))
	for i, order := range res {
		output[i] = &FindOrderOutputDTO{
			Id:             order.Id,
			CrowdfundingId: order.CrowdfundingId,
			Investor:       order.Investor,
			Amount:         order.Amount,
			InterestRate:   order.InterestRate,
			State:          string(order.State),
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		}
	}
	return output, nil
}
