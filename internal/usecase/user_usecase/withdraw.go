package user_usecase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type WithdrawInputDTO struct {
	Token  common.Address `json:"token"`
	Amount *uint256.Int   `json:"amount"`
}
