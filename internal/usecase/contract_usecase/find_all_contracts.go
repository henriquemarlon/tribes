package contract_usecase

import "github.com/Mugen-Builders/devolt/internal/domain/entity"

type FindAllContractsOutputDTO []*FindContractOutputDTO

type FindAllContractsUsecase struct {
	ContractRepository entity.ContractRepository
}

func NewFindAllContractsUseCase(contractRepository entity.ContractRepository) *FindAllContractsUsecase {
	return &FindAllContractsUsecase{
		ContractRepository: contractRepository,
	}
}

func (s *FindAllContractsUsecase) Execute() (FindAllContractsOutputDTO, error) {
	res, err := s.ContractRepository.FindAllContracts()
	if err != nil {
		return nil, err
	}
	var output FindAllContractsOutputDTO
	for _, contract := range res {
		dto := &FindContractOutputDTO{
			Id:        contract.Id,
			Symbol:    contract.Symbol,
			Address:   contract.Address,
			CreatedAt: contract.CreatedAt,
			UpdatedAt: contract.UpdatedAt,
		}
		output = append(output, dto)
	}
	return output, nil
}
