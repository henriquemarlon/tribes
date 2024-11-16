package user_usecase

import "github.com/ethereum/go-ethereum/common"

type FindUserOutputDTO struct {
	Id        uint           `json:"id"`
	Role      string         `json:"role"`
	Address   common.Address `json:"address"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
}
