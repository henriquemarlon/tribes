package main

import (
	"log"

	"github.com/tribeshq/tribes/pkg/router"
)

func NewDApp() *router.Router {
	//////////////////////// Setup Handlers //////////////////////////
	ah, err := NewAdvanceHandlers()
	if err != nil {
		log.Fatalf("Failed to initialize advance handlers from wire: %v", err)
	}

	ih, err := NewInspectHandlers()
	if err != nil {
		log.Fatalf("Failed to initialize inspect handlers from wire: %v", err)
	}

	ms, err := NewMiddlewares()
	if err != nil {
		log.Fatalf("Failed to initialize middlewares from wire: %v", err)
	}

	//////////////////////// Router //////////////////////////
	app := router.NewRouter()

	//////////////////////// Advance //////////////////////////
	app.HandleAdvance("createContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.CreateContractHandler, "admin"))
	app.HandleAdvance("updateContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.UpdateContractHandler, "admin"))
	app.HandleAdvance("deleteContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.DeleteContractHandler, "admin"))

	app.HandleAdvance("createBid", ah.BidAdvanceHandlers.CreateBidHandler)

	app.HandleAdvance("createAuction", ms.RBAC.Middleware(ah.AuctionAdvanceHandlers.CreateAuctionHandler, "admin"))
	app.HandleAdvance("finishAuction", ms.RBAC.Middleware(ah.AuctionAdvanceHandlers.FinishAuctionHandler, "admin"))

	// app.HandleAdvance("withdrawApp", ms.RBAC.Middleware(ah.UserAdvanceHandlers.WithdrawAppHandler, "admin"))
	app.HandleAdvance("withdraw", ah.UserAdvanceHandlers.WithdrawHandler)

	app.HandleAdvance("createUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.CreateUserHandler, "admin"))
	app.HandleAdvance("deleteUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.DeleteUserByAddressHandler, "admin"))

	//////////////////////// Inspect //////////////////////////

	app.HandleInspect("auction", ih.AuctionInspectHandlers.FindAllAuctionsHandler)
	app.HandleInspect("auction/{id}", ih.AuctionInspectHandlers.FindAuctionByIdHandler)
	app.HandleInspect("auction/active", ih.AuctionInspectHandlers.FindActiveAuctionHandler)

	app.HandleInspect("bid", ih.BidInspectHandlers.FindAllBidsHandler)
	app.HandleInspect("bid/{id}", ih.BidInspectHandlers.FindBidByIdHandler)
	app.HandleInspect("bid/auction/{id}", ih.BidInspectHandlers.FindBisdByAuctionIdHandler)

	app.HandleInspect("contract", ih.ContractInspectHandlers.FindAllContractsHandler)
	app.HandleInspect("contract/{symbol}", ih.ContractInspectHandlers.FindContractBySymbolHandler)

	app.HandleInspect("user", ih.UserInspectHandlers.FindAllUsersHandler)
	app.HandleInspect("user/{address}", ih.UserInspectHandlers.FindUserByAddressHandler)
	app.HandleInspect("balance/{symbol}/{address}", ih.UserInspectHandlers.BalanceHandler)

	return app
}

func NewDAppMemory() *router.Router {
	//////////////////////// Setup Handlers //////////////////////////
	ah, err := NewAdvanceHandlersMemory()
	if err != nil {
		log.Fatalf("Failed to initialize advance handlers from wire: %v", err)
	}

	ih, err := NewInspectHandlersMemory()
	if err != nil {
		log.Fatalf("Failed to initialize inspect handlers from wire: %v", err)
	}

	ms, err := NewMiddlewaresMemory()
	if err != nil {
		log.Fatalf("Failed to initialize middlewares from wire: %v", err)
	}

	//////////////////////// Router //////////////////////////
	app := router.NewRouter()

	//////////////////////// Advance //////////////////////////
	app.HandleAdvance("createContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.CreateContractHandler, "admin"))
	app.HandleAdvance("updateContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.UpdateContractHandler, "admin"))
	app.HandleAdvance("deleteContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.DeleteContractHandler, "admin"))

	app.HandleAdvance("createBid", ah.BidAdvanceHandlers.CreateBidHandler)

	app.HandleAdvance("createAuction", ms.RBAC.Middleware(ah.AuctionAdvanceHandlers.CreateAuctionHandler, "admin"))
	app.HandleAdvance("finishAuction", ms.RBAC.Middleware(ah.AuctionAdvanceHandlers.FinishAuctionHandler, "admin"))

	// app.HandleAdvance("withdrawApp", ms.RBAC.Middleware(ah.UserAdvanceHandlers.WithdrawAppHandler, "admin"))
	app.HandleAdvance("withdrawStablecoin", ah.UserAdvanceHandlers.WithdrawHandler)

	app.HandleAdvance("createUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.CreateUserHandler, "admin"))
	app.HandleAdvance("deleteUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.DeleteUserByAddressHandler, "admin"))

	//////////////////////// Inspect //////////////////////////
	app.HandleInspect("auction", ih.AuctionInspectHandlers.FindAllAuctionsHandler)
	app.HandleInspect("auction/{id}", ih.AuctionInspectHandlers.FindAuctionByIdHandler)
	app.HandleInspect("auction/active", ih.AuctionInspectHandlers.FindActiveAuctionHandler)

	app.HandleInspect("bid", ih.BidInspectHandlers.FindAllBidsHandler)
	app.HandleInspect("bid/{id}", ih.BidInspectHandlers.FindBidByIdHandler)
	app.HandleInspect("bid/auction/{id}", ih.BidInspectHandlers.FindBisdByAuctionIdHandler)

	app.HandleInspect("contract", ih.ContractInspectHandlers.FindAllContractsHandler)
	app.HandleInspect("contract/{symbol}", ih.ContractInspectHandlers.FindContractBySymbolHandler)

	app.HandleInspect("user", ih.UserInspectHandlers.FindAllUsersHandler)
	app.HandleInspect("user/{address}", ih.UserInspectHandlers.FindUserByAddressHandler)
	app.HandleInspect("balance/{symbol}/{address}", ih.UserInspectHandlers.BalanceHandler)

	return app
}
