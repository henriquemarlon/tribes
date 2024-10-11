package contract_usecase

import "github.com/Mugen-Builders/devolt/internal/domain/entity"

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

func (s *DeleteContractUseCase) Execute(input *DeleteContractInputDTO) error {
	return s.ContractReposiotry.DeleteContract(input.Symbol)
}
