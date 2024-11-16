package crowdfunding_usecase

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
)

type DeleteCrowdfundingInputDTO struct {
	Id uint `json:"id"`
}

type DeleteCrowdfundingUseCase struct {
	CrowdfundingRepository entity.CrowdfundingRepository
}

func NewDeleteCrowdfundingUseCase(crowdfundingRepository entity.CrowdfundingRepository) *DeleteCrowdfundingUseCase {
	return &DeleteCrowdfundingUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (u *DeleteCrowdfundingUseCase) Execute(input *DeleteCrowdfundingInputDTO) error {
	return u.CrowdfundingRepository.DeleteCrowdfunding(input.Id)
}
