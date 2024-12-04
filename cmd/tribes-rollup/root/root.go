package root

import (
	"log/slog"
	"os"
	"time"

	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/pkg/router"
	"gorm.io/gorm"
)

const (
	CMD_NAME = "tribes-rollup"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   CMD_NAME,
		Short: "Runs Tribes Rollup",
		Long:  `Runs Tribes Rollup`,
		Run:   run,
	}
)

func init() {
	Cmd.PersistentFlags().BoolVar(
		&useMemoryDB,
		"memory-db",
		false,
		"Use in-memory SQLite database instead of persistent",
	)
}

func run(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	ctx := cmd.Context()

	var db *gorm.DB
	var err error
	if useMemoryDB {
		db, err = configs.SetupSQlite(":memory:")
		if err != nil {
			slog.Error("Failed to setup in-memory SQLite database", "error", err)
			os.Exit(1)
		}
		slog.Info("In-memory database initialized")
	} else {
		db, err = configs.SetupSQlite("tribes.db")
		if err != nil {
			slog.Error("Failed to setup SQLite database", "error", err)
			os.Exit(1)
		}
		slog.Info("Persistent database initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get SQL DB from GORM", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	ah, err := NewAdvanceHandlers(db)
	if err != nil {
		slog.Error("Failed to initialize advance handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Advance handlers initialized")

	ih, err := NewInspectHandlers(db)
	if err != nil {
		slog.Error("Failed to initialize inspect handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Inspect handlers initialized")

	ms, err := NewMiddlewares(db)
	if err != nil {
		slog.Error("Failed to initialize middlewares", "error", err)
		os.Exit(1)
	}
	slog.Info("Middlewares initialized")

	r := NewDApp(ah, ih, ms)
	slog.Info("Router initialized")

	opts := rollmelette.NewRunOpts()
	if rollupUrl, isSet := os.LookupEnv("ROLLUP_HTTP_SERVER_URL"); isSet {
		opts.RollupURL = rollupUrl
	}

	ready := make(chan struct{}, 1)
	go func() {
		select {
		case <-ready:
			duration := time.Since(startTime)
			slog.Info("DApp is ready", "after", duration)
		case <-ctx.Done():
		}
	}()

	if err := rollmelette.Run(ctx, opts, r); err != nil {
		slog.Error("Application exited with an error", "error", err)
		os.Exit(1)
	}
}

func NewDApp(ah *AdvanceHandlers, ih *InspectHandlers, ms *Middlewares) *router.Router {
	r := router.NewRouter()

	r.HandleAdvance("createContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.CreateContractHandler, []string{"admin"}))
	r.HandleAdvance("updateContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.UpdateContractHandler, []string{"admin"}))
	r.HandleAdvance("deleteContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.DeleteContractHandler, []string{"admin"}))

	r.HandleAdvance("createOrder", ms.RBAC.Middleware(ah.OrderAdvanceHandlers.CreateOrderHandler, []string{"non_qualified_investor", "qualified_investor"}))
	r.HandleAdvance("cancelOrder", ms.RBAC.Middleware(ah.OrderAdvanceHandlers.CancelOrderHandler, []string{"non_qualified_investor", "qualified_investor"}))

	r.HandleAdvance("createCrowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.CreateCrowdfundingHandler, []string{"creator"}))
	r.HandleAdvance("deleteCrowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.DeleteCrowdfundingHandler, []string{"admin"}))
	r.HandleAdvance("updateCrowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.UpdateCrowdfundingHandler, []string{"admin"}))
	r.HandleAdvance("closeCrowdfunding", ah.CrowdfundingAdvanceHandlers.CloseCrowdfundingHandler)
	r.HandleAdvance("settleCrowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.SettleCrowdfundingHandler, []string{"creator"}))

	r.HandleAdvance("createUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.CreateUserHandler, []string{"admin"}))
	r.HandleAdvance("updateUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.UpdateUserHandler, []string{"admin"}))
	r.HandleAdvance("deleteUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.DeleteUserHandler, []string{"admin"}))
	r.HandleAdvance("withdraw", ah.UserAdvanceHandlers.WithdrawHandler)

	r.HandleAdvance("createSocialAccount", ms.RBAC.Middleware(ah.SocialAccountsHandlers.CreateSocialAccountHandler, []string{"admin"}))
	r.HandleAdvance("deleteSocialAccount", ms.RBAC.Middleware(ah.SocialAccountsHandlers.DeleteSocialAccountHandler, []string{"admin"}))

	r.HandleInspect("crowdfunding", ih.CrowdfundingInspectHandlers.FindAllCrowdfundingsHandler)
	r.HandleInspect("crowdfunding/{id}", ih.CrowdfundingInspectHandlers.FindCrowdfundingByIdHandler)
	r.HandleInspect("crowdfunding/creator/{address}", ih.CrowdfundingInspectHandlers.FindCrowdfundingsByCreatorHandler)
	r.HandleInspect("crowdfunding/investor/{address}", ih.CrowdfundingInspectHandlers.FindCrowdfundingsByInvestorHandler)

	r.HandleInspect("order", ih.OrderInspectHandlers.FindAllOrdersHandler)
	r.HandleInspect("order/{id}", ih.OrderInspectHandlers.FindOrderByIdHandler)
	r.HandleInspect("order/investor/{address}", ih.OrderInspectHandlers.FindOrdersByInvestorHandler)
	r.HandleInspect("order/crowdfunding/{id}", ih.OrderInspectHandlers.FindBisdByCrowdfundingIdHandler)

	r.HandleInspect("contract", ih.ContractInspectHandlers.FindAllContractsHandler)
	r.HandleInspect("contract/{symbol}", ih.ContractInspectHandlers.FindContractBySymbolHandler)

	r.HandleInspect("user", ih.UserInspectHandlers.FindAllUsersHandler)
	r.HandleInspect("user/{address}", ih.UserInspectHandlers.FindUserByAddressHandler)
	r.HandleInspect("balance/{address}", ih.UserInspectHandlers.BalanceHandler)

	r.HandleInspect("social/{id}", ih.SocialAccountHandlers.FindSocialAccountById)
	r.HandleInspect("social/user/{id}", ih.SocialAccountHandlers.FindSocialAccountsByUserId)

	return r
}
