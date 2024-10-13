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

var setBidRepositoryDependency = wire.NewSet(
	db.NewBidRepositorySqlite,
	wire.Bind(new(entity.BidRepository), new(*db.BidRepositorySqlite)),
)

var setAuctionRepositoryDependency = wire.NewSet(
	db.NewAuctionRepositorySqlite,
	wire.Bind(new(entity.AuctionRepository), new(*db.AuctionRepositorySqlite)),
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
	advance_handler.NewBidAdvanceHandlers,
	advance_handler.NewUserAdvanceHandlers,
	advance_handler.NewAuctionAdvanceHandlers,
	advance_handler.NewContractAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewBidInspectHandlers,
	inspect_handler.NewUserInspectHandlers,
	inspect_handler.NewAuctionInspectHandlers,
	inspect_handler.NewContractInspectHandlers,
)

var setMiddleware = wire.NewSet(
	middleware.NewTLSNMiddleware,
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

func NewMiddlewaresMemory(gormDB *gorm.DB) (*Middlewares, error) {
	wire.Build(
		setUserRepositoryDependency,
		setMiddleware,
		wire.Struct(new(Middlewares), "*"),
	)
	return nil, nil
}

func NewAdvanceHandlers(gormDB *gorm.DB) (*AdvanceHandlers, error) {
	wire.Build(
		setBidRepositoryDependency,
		setUserRepositoryDependency,
		setAuctionRepositoryDependency,
		setContractRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewAdvanceHandlersMemory(gormDB *gorm.DB) (*AdvanceHandlers, error) {
	wire.Build(
		setBidRepositoryDependency,
		setUserRepositoryDependency,
		setAuctionRepositoryDependency,
		setContractRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlers(gormDB *gorm.DB) (*InspectHandlers, error) {
	wire.Build(
		setBidRepositoryDependency,
		setUserRepositoryDependency,
		setAuctionRepositoryDependency,
		setContractRepositoryDependency,
		setInspectHandlers,
		wire.Struct(new(InspectHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlersMemory(gormDB *gorm.DB) (*InspectHandlers, error) {
	wire.Build(
		setBidRepositoryDependency,
		setUserRepositoryDependency,
		setAuctionRepositoryDependency,
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
	BidAdvanceHandlers      *advance_handler.BidAdvanceHandlers
	UserAdvanceHandlers     *advance_handler.UserAdvanceHandlers
	AuctionAdvanceHandlers  *advance_handler.AuctionAdvanceHandlers
	ContractAdvanceHandlers *advance_handler.ContractAdvanceHandlers
}

type InspectHandlers struct {
	BidInspectHandlers      *inspect_handler.BidInspectHandlers
	UserInspectHandlers     *inspect_handler.UserInspectHandlers
	AuctionInspectHandlers  *inspect_handler.AuctionInspectHandlers
	ContractInspectHandlers *inspect_handler.ContractInspectHandlers
}
