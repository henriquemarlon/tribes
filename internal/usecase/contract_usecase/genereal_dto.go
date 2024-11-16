package contract_usecase

import "github.com/ethereum/go-ethereum/common"

type FindContractOutputDTO struct {
	Id        uint           `json:"id"`
	Symbol    string         `json:"symbol"`
	Address   common.Address `json:"address"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
}
