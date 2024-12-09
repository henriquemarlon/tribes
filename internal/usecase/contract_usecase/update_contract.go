package contract_usecase

import (
	"context"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type UpdateContractInputDTO struct {
	Id      uint                `json:"id"`
	Address custom_type.Address `json:"address"`
	Symbol  string              `json:"symbol"`
}

type UpdateContractOutputDTO struct {
	Id        uint                `json:"id"`
	Symbol    string              `json:"symbol"`
	Address   custom_type.Address `json:"address"`
	CreatedAt int64               `json:"created_at"`
	UpdatedAt int64               `json:"updated_at"`
}

type UpdateContractUseCase struct {
	ContractReposiotry entity.ContractRepository
}

func NewUpdateContractUseCase(contractRepository entity.ContractRepository) *UpdateContractUseCase {
	return &UpdateContractUseCase{
		ContractReposiotry: contractRepository,
	}
}

func (s *UpdateContractUseCase) Execute(ctx context.Context, input *UpdateContractInputDTO, metadata rollmelette.Metadata) (*UpdateContractOutputDTO, error) {
	contract, err := s.ContractReposiotry.UpdateContract(ctx, &entity.Contract{
		Id:        input.Id,
		Address:   input.Address,
		Symbol:    input.Symbol,
		UpdatedAt: metadata.BlockTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateContractOutputDTO{
		Id:        contract.Id,
		Symbol:    contract.Symbol,
		Address:   contract.Address,
		CreatedAt: contract.CreatedAt,
		UpdatedAt: contract.UpdatedAt,
	}, nil
}
