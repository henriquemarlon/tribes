package advance_handler

import (
	"encoding/json"
	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/internal/usecase/contract_usecase"
	"github.com/rollmelette/rollmelette"
)

type ContractAdvanceHandlers struct {
	ContractRepository entity.ContractRepository
}

func NewContractAdvanceHandlers(contractRepository entity.ContractRepository) *ContractAdvanceHandlers {
	return &ContractAdvanceHandlers{
		ContractRepository: contractRepository,
	}
}

func (h *ContractAdvanceHandlers) CreateContractHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.CreateContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createContract := contract_usecase.NewCreateContractUseCase(h.ContractRepository)
	res, err := createContract.Execute(&input, metadata)
	if err != nil {
		return err
	}
	contract, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("created contract - "), contract...))
	return nil
}

func (h *ContractAdvanceHandlers) UpdateContractHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.UpdateContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	updateContract := contract_usecase.NewUpdateContractUseCase(h.ContractRepository)
	res, err := updateContract.Execute(&input, metadata)
	if err != nil {
		return err
	}
	contract, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("updated contract - "), contract...))
	return nil
}

func (h *ContractAdvanceHandlers) DeleteContractHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.DeleteContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	deleteContract := contract_usecase.NewDeleteContractUseCase(h.ContractRepository)
	err := deleteContract.Execute(&input)
	if err != nil {
		return err
	}
	contract, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("deleted contract with - "), contract...))
	return nil
}
