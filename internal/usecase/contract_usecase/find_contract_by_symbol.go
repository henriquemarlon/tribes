package contract_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
)

type FindContractBySymbolInputDTO struct {
	Symbol string
}

type FindContractBySymbolUseCase struct {
	ContractRepository entity.ContractRepository
}

func NewFindContractBySymbolUseCase(contractRepository entity.ContractRepository) *FindContractBySymbolUseCase {
	return &FindContractBySymbolUseCase{
		ContractRepository: contractRepository,
	}
}

func (s *FindContractBySymbolUseCase) Execute(input *FindContractBySymbolInputDTO) (*FindContractOutputDTO, error) {
	contract, err := s.ContractRepository.FindContractBySymbol(input.Symbol)
	if err != nil {
		return nil, err
	}
	return &FindContractOutputDTO{
		Id:        contract.Id,
		Symbol:    contract.Symbol,
		Address:   contract.Address,
		CreatedAt: contract.CreatedAt,
		UpdatedAt: contract.UpdatedAt,
	}, nil
}
