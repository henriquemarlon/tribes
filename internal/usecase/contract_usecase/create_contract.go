package contract_usecase

import (
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type CreateContractInputDTO struct {
	Symbol  string              `json:"symbol"`
	Address custom_type.Address `json:"address"`
}

type CreateContractOutputDTO struct {
	Id        uint                `json:"id"`
	Symbol    string              `json:"symbol"`
	Address   custom_type.Address `json:"address"`
	CreatedAt int64               `json:"created_at"`
}

type CreateContractUseCase struct {
	ContractRepository entity.ContractRepository
}

func NewCreateContractUseCase(contractRepository entity.ContractRepository) *CreateContractUseCase {
	return &CreateContractUseCase{
		ContractRepository: contractRepository,
	}
}

func (s *CreateContractUseCase) Execute(input *CreateContractInputDTO, metadata rollmelette.Metadata) (*CreateContractOutputDTO, error) {
	contract, err := entity.NewContract(input.Symbol, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := s.ContractRepository.CreateContract(contract)
	if err != nil {
		return nil, err
	}
	output := &CreateContractOutputDTO{
		Id:        res.Id,
		Symbol:    res.Symbol,
		Address:   res.Address,
		CreatedAt: res.CreatedAt,
	}
	return output, nil
}
