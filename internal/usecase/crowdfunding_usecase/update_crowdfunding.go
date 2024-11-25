package crowdfunding_usecase

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type UpdateCrowdfundingInputDTO struct {
	Id              uint         `json:"id"`
	DebtIssued      *uint256.Int `json:"debt_issued"`
	MaxInterestRate *uint256.Int `json:"max_interest_rate"`
	TotalObligation *uint256.Int `json:"total_obligation"`
	State           string       `json:"state"`
	ClosesAt        int64        `json:"closes_at"`
	MaturityAt      int64        `json:"maturity_at"`
}

type UpdateCrowdfundingOutputDTO struct {
	Id              uint            `json:"id"`
	Creator         common.Address  `json:"creator"`
	DebtIssued      *uint256.Int    `json:"debt_issued"`
	MaxInterestRate *uint256.Int    `json:"max_interest_rate"`
	TotalObligation *uint256.Int    `json:"total_obligation"`
	State           string          `json:"state"`
	Orders          []*entity.Order `json:"orders"`
	ClosesAt        int64           `json:"closes_at"`
	MaturityAt      int64           `json:"maturity_at"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`
}

type UpdateCrowdfundingUsecase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewUpdateCrowdfundingUseCase(crowdfundingRepository entity.CrowdfundingRepository) *UpdateCrowdfundingUsecase {
	return &UpdateCrowdfundingUsecase{
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (uc *UpdateCrowdfundingUsecase) Execute(ctx context.Context, input UpdateCrowdfundingInputDTO, metadata rollmelette.Metadata) (*UpdateCrowdfundingOutputDTO, error) {
	crowdfunding, err := uc.CrowdfundingRepository.UpdateCrowdfunding(ctx, &entity.Crowdfunding{
		Id:              input.Id,
		DebtIssued:      input.DebtIssued,
		MaxInterestRate: input.MaxInterestRate,
		TotalObligation: input.TotalObligation,
		State:           entity.CrowdfundingState(input.State),
		ClosesAt:        input.ClosesAt,
		MaturityAt:      input.ClosesAt,
		UpdatedAt:       metadata.BlockTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateCrowdfundingOutputDTO{
		Id:              crowdfunding.Id,
		Creator:         crowdfunding.Creator,
		DebtIssued:      crowdfunding.DebtIssued,
		MaxInterestRate: crowdfunding.MaxInterestRate,
		TotalObligation: crowdfunding.TotalObligation,
		State:           string(crowdfunding.State),
		Orders:          crowdfunding.Orders,
		ClosesAt:        crowdfunding.ClosesAt,
		MaturityAt:      crowdfunding.MaturityAt,
		CreatedAt:       crowdfunding.CreatedAt,
		UpdatedAt:       crowdfunding.UpdatedAt,
	}, nil
}
