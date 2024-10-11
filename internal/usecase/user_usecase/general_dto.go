package user_usecase

import "github.com/Mugen-Builders/devolt/pkg/custom_type"

type FindUserOutputDTO struct {
	Id        uint                `json:"id"`
	Role      string              `json:"role"`
	Address   custom_type.Address `json:"address"`
	CreatedAt int64               `json:"created_at"`
	UpdatedAt int64               `json:"updated_at"`
}
