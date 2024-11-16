package main

import (
	"log"

	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/pkg/router"
)

func NewDAppMemory() *router.Router {
	//////////////////////// Setup Database //////////////////////////
	db, err := configs.SetupSQliteMemory()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	//////////////////////// Setup Handlers //////////////////////////
	ah, err := NewAdvanceHandlersMemory(db)
	if err != nil {
		log.Fatalf("Failed to initialize advance handlers from wire: %v", err)
	}

	ih, err := NewInspectHandlersMemory(db)
	if err != nil {
		log.Fatalf("Failed to initialize inspect handlers from wire: %v", err)
	}

	ms, err := NewMiddlewaresMemory(db)
	if err != nil {
		log.Fatalf("Failed to initialize middlewares from wire: %v", err)
	}

	//////////////////////// Router //////////////////////////
	app := router.NewRouter()

	//////////////////////// Advance //////////////////////////
	app.HandleAdvance("createContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.CreateContractHandler, "admin"))
	app.HandleAdvance("updateContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.UpdateContractHandler, "admin"))
	app.HandleAdvance("deleteContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.DeleteContractHandler, "admin"))

	app.HandleAdvance("createOrder", ah.OrderAdvanceHandlers.CreateOrderHandler)

	app.HandleAdvance("createCrowdfunding", ms.TLSN.Middleware(ah.CrowdfundingAdvanceHandlers.CreateCrowdfundingHandler))
	app.HandleAdvance("closeCrowdfunding", ah.CrowdfundingAdvanceHandlers.CloseCrowdfundingHandler)
	app.HandleAdvance("settleCrowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.SettleCrowdfundingHandler, "creator"))

	app.HandleAdvance("withdraw", ah.UserAdvanceHandlers.WithdrawHandler)
	app.HandleAdvance("withdrawApp", ms.RBAC.Middleware(ah.UserAdvanceHandlers.WithdrawAppHandler, "admin"))

	app.HandleAdvance("createUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.CreateUserHandler, "admin"))
	app.HandleAdvance("deleteUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.DeleteUserHandler, "admin"))

	//////////////////////// Inspect //////////////////////////
	app.HandleInspect("crowdfunding", ih.CrowdfundingInspectHandlers.FindAllCrowdfundingsHandler)
	app.HandleInspect("crowdfunding/{id}", ih.CrowdfundingInspectHandlers.FindCrowdfundingByIdHandler)

	app.HandleInspect("order", ih.OrderInspectHandlers.FindAllOrdersHandler)
	app.HandleInspect("order/{id}", ih.OrderInspectHandlers.FindOrderByIdHandler)
	app.HandleInspect("order/crowdfunding/{id}", ih.OrderInspectHandlers.FindBisdByCrowdfundingIdHandler)

	app.HandleInspect("contract", ih.ContractInspectHandlers.FindAllContractsHandler)
	app.HandleInspect("contract/{symbol}", ih.ContractInspectHandlers.FindContractBySymbolHandler)

	app.HandleInspect("user", ih.UserInspectHandlers.FindAllUsersHandler)
	app.HandleInspect("user/{address}", ih.UserInspectHandlers.FindUserByAddressHandler)
	app.HandleInspect("balance/{symbol}/{address}", ih.UserInspectHandlers.BalanceHandler)

	return app
}
