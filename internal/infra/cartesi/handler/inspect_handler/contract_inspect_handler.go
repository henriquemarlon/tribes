package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/internal/usecase/contract_usecase"
	"github.com/Mugen-Builders/devolt/pkg/router"
	"github.com/rollmelette/rollmelette"
)

type ContractInspectHandlers struct {
	ContractRepository entity.ContractRepository
}

func NewContractInspectHandlers(contractRepository entity.ContractRepository) *ContractInspectHandlers {
	return &ContractInspectHandlers{
		ContractRepository: contractRepository,
	}
}

func (h *ContractInspectHandlers) FindAllContractsHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	findAllContracts := contract_usecase.NewFindAllContractsUseCase(h.ContractRepository)
	contracts, err := findAllContracts.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all contracts: %w", err)
	}
	contractsBytes, err := json.Marshal(contracts)
	if err != nil {
		return fmt.Errorf("failed to marshal contracts: %w", err)
	}
	env.Report(contractsBytes)
	return nil
}

func (h *ContractInspectHandlers) FindContractBySymbolHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	symbol := strings.ToUpper(router.PathValue(ctx, "symbol"))
	findOrderBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	contract, err := findOrderBySymbol.Execute(&contract_usecase.FindContractBySymbolInputDTO{
		Symbol: symbol,
	})
	if err != nil {
		return err
	}
	contractBytes, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}
	env.Report(contractBytes)
	return nil
}
