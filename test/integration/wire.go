//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/advance_handler"
	"github.com/tribeshq/tribes/internal/infra/cartesi/handler/inspect_handler"
	"github.com/tribeshq/tribes/internal/infra/cartesi/middleware"
	db "github.com/tribeshq/tribes/internal/infra/repository"
	"gorm.io/gorm"
)

var setOrderRepositoryDependency = wire.NewSet(
	db.NewOrderRepositorySqlite,
	wire.Bind(new(entity.OrderRepository), new(*db.OrderRepositorySqlite)),
)

var setCrowdfundingRepositoryDependency = wire.NewSet(
	db.NewCrowdfundingRepositorySqlite,
	wire.Bind(new(entity.CrowdfundingRepository), new(*db.CrowdfundingRepositorySqlite)),
)

var setContractRepositoryDependency = wire.NewSet(
	db.NewContractRepositorySqlite,
	wire.Bind(new(entity.ContractRepository), new(*db.ContractRepositorySqlite)),
)

var setUserRepositoryDependency = wire.NewSet(
	db.NewUserRepositorySqlite,
	wire.Bind(new(entity.UserRepository), new(*db.UserRepositorySqlite)),
)

var setAdvanceHandlers = wire.NewSet(
	advance_handler.NewOrderAdvanceHandlers,
	advance_handler.NewUserAdvanceHandlers,
	advance_handler.NewCrowdfundingAdvanceHandlers,
	advance_handler.NewContractAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewOrderInspectHandlers,
	inspect_handler.NewUserInspectHandlers,
	inspect_handler.NewCrowdfundingInspectHandlers,
	inspect_handler.NewContractInspectHandlers,
)

var setMiddleware = wire.NewSet(
	middleware.NewTLSNMiddleware,
	middleware.NewRBACMiddleware,
)

func NewMiddlewaresMemory(gormDB *gorm.DB) (*Middlewares, error) {
	wire.Build(
		setUserRepositoryDependency,
		setMiddleware,
		wire.Struct(new(Middlewares), "*"),
	)
	return nil, nil
}

func NewAdvanceHandlersMemory(gormDB *gorm.DB) (*AdvanceHandlers, error) {
	wire.Build(
		setOrderRepositoryDependency,
		setUserRepositoryDependency,
		setCrowdfundingRepositoryDependency,
		setContractRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlersMemory(gormDB *gorm.DB) (*InspectHandlers, error) {
	wire.Build(
		setOrderRepositoryDependency,
		setUserRepositoryDependency,
		setCrowdfundingRepositoryDependency,
		setContractRepositoryDependency,
		setInspectHandlers,
		wire.Struct(new(InspectHandlers), "*"),
	)
	return nil, nil
}

type Middlewares struct {
	TLSN *middleware.TLSNMiddleware
	RBAC *middleware.RBACMiddleware
}

type AdvanceHandlers struct {
	OrderAdvanceHandlers        *advance_handler.OrderAdvanceHandlers
	UserAdvanceHandlers         *advance_handler.UserAdvanceHandlers
	CrowdfundingAdvanceHandlers *advance_handler.CrowdfundingAdvanceHandlers
	ContractAdvanceHandlers     *advance_handler.ContractAdvanceHandlers
}

type InspectHandlers struct {
	OrderInspectHandlers        *inspect_handler.OrderInspectHandlers
	UserInspectHandlers         *inspect_handler.UserInspectHandlers
	CrowdfundingInspectHandlers *inspect_handler.CrowdfundingInspectHandlers
	ContractInspectHandlers     *inspect_handler.ContractInspectHandlers
}
