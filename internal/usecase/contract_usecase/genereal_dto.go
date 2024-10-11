package contract_usecase

import "github.com/Mugen-Builders/devolt/pkg/custom_type"

type FindContractOutputDTO struct {
	Id        uint                `json:"id"`
	Symbol    string              `json:"symbol"`
	Address   custom_type.Address `json:"address"`
	CreatedAt int64               `json:"created_at"`
	UpdatedAt int64               `json:"updated_at"`
}
