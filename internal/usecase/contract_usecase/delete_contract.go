package contract_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type DeleteContractInputDTO struct {
	Symbol string
}

type DeleteContractUseCase struct {
	ContractReposiotry entity.ContractRepository
}

func NewDeleteContractUseCase(contractRepository entity.ContractRepository) *DeleteContractUseCase {
	return &DeleteContractUseCase{
		ContractReposiotry: contractRepository,
	}
}

func (s *DeleteContractUseCase) Execute(ctx context.Context, input *DeleteContractInputDTO) error {
	return s.ContractReposiotry.DeleteContract(ctx, input.Symbol)
}
