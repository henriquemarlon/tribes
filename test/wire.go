//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/advance_handler"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/inspect_handler"
	"github.com/tribeshq/tribes/internal/infra/cartesi/middleware"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"gorm.io/gorm"
)

var setOrderRepositoryDependency = wire.NewSet(
	repository.NewOrderRepositorySqlite,
	wire.Bind(new(entity.OrderRepository), new(*repository.OrderRepositorySqlite)),
)

var setCrowdfundingRepositoryDependency = wire.NewSet(
	repository.NewCrowdfundingRepositorySqlite,
	wire.Bind(new(entity.CrowdfundingRepository), new(*repository.CrowdfundingRepositorySqlite)),
)

var setContractRepositoryDependency = wire.NewSet(
	repository.NewContractRepositorySqlite,
	wire.Bind(new(entity.ContractRepository), new(*repository.ContractRepositorySqlite)),
)

var setUserRepositoryDependency = wire.NewSet(
	repository.NewUserRepositorySqlite,
	wire.Bind(new(entity.UserRepository), new(*repository.UserRepositorySqlite)),
)

var setSocialAccountDependecy = wire.NewSet(
	repository.NewSocialAccountRepositorySqlite,
	wire.Bind(new(entity.SocialAccountRepository), new(*repository.SocialAccountRepositorySqlite)),
)

var setAdvanceHandlers = wire.NewSet(
	advance_handler.NewOrderAdvanceHandlers,
	advance_handler.NewUserAdvanceHandlers,
	advance_handler.NewSocialAccountAdvanceHandlers,
	advance_handler.NewCrowdfundingAdvanceHandlers,
	advance_handler.NewContractAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewOrderInspectHandlers,
	inspect_handler.NewUserInspectHandlers,
	inspect_handler.NewSocialAccountInspectHandlers,
	inspect_handler.NewCrowdfundingInspectHandlers,
	inspect_handler.NewContractInspectHandlers,
)

var setMiddleware = wire.NewSet(
	middleware.NewRBACMiddleware,
)

func NewMiddlewares(gormDB *gorm.DB) (*Middlewares, error) {
	wire.Build(
		setUserRepositoryDependency,
		setMiddleware,
		wire.Struct(new(Middlewares), "*"),
	)
	return nil, nil
}

func NewAdvanceHandlers(gormDB *gorm.DB) (*AdvanceHandlers, error) {
	wire.Build(
		setOrderRepositoryDependency,
		setUserRepositoryDependency,
		setSocialAccountDependecy,
		setCrowdfundingRepositoryDependency,
		setContractRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlers(gormDB *gorm.DB) (*InspectHandlers, error) {
	wire.Build(
		setOrderRepositoryDependency,
		setUserRepositoryDependency,
		setSocialAccountDependecy,
		setCrowdfundingRepositoryDependency,
		setContractRepositoryDependency,
		setInspectHandlers,
		wire.Struct(new(InspectHandlers), "*"),
	)
	return nil, nil
}

type Middlewares struct {
	RBAC *middleware.RBACMiddleware
}

type AdvanceHandlers struct {
	OrderAdvanceHandlers        *advance_handler.OrderAdvanceHandlers
	UserAdvanceHandlers         *advance_handler.UserAdvanceHandlers
	SocialAccountsHandlers *advance_handler.SocialAccountAdvanceHandlers
	CrowdfundingAdvanceHandlers *advance_handler.CrowdfundingAdvanceHandlers
	ContractAdvanceHandlers     *advance_handler.ContractAdvanceHandlers
}

type InspectHandlers struct {
	OrderInspectHandlers        *inspect_handler.OrderInspectHandlers
	UserInspectHandlers         *inspect_handler.UserInspectHandlers
	SocialAccountHandlers			 	*inspect_handler.SocialAccountInspectHandlers
	CrowdfundingInspectHandlers *inspect_handler.CrowdfundingInspectHandlers
	ContractInspectHandlers     *inspect_handler.ContractInspectHandlers
}
