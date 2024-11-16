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
	CMD_NAME = "tribes-dapp"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   CMD_NAME,
		Short: "Runs Tribes DApp",
		Long:  `Runs Tribes DApp`,
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

	r.HandleAdvance("createContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.CreateContractHandler, "admin"))
	r.HandleAdvance("updateContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.UpdateContractHandler, "admin"))
	r.HandleAdvance("deleteContract", ms.RBAC.Middleware(ah.ContractAdvanceHandlers.DeleteContractHandler, "admin"))

	r.HandleAdvance("createOrder", ah.OrderAdvanceHandlers.CreateOrderHandler)

	r.HandleAdvance("createCrowdfunding", ms.TLSN.Middleware(ah.CrowdfundingAdvanceHandlers.CreateCrowdfundingHandler))
	r.HandleAdvance("closeCrowdfunding", ah.CrowdfundingAdvanceHandlers.CloseCrowdfundingHandler)
	r.HandleAdvance("settleCrowdfunding", ms.RBAC.Middleware(ah.CrowdfundingAdvanceHandlers.SettleCrowdfundingHandler, "creator"))

	r.HandleAdvance("withdraw", ah.UserAdvanceHandlers.WithdrawHandler)
	r.HandleAdvance("withdrawApp", ms.RBAC.Middleware(ah.UserAdvanceHandlers.WithdrawAppHandler, "admin"))

	r.HandleAdvance("createUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.CreateUserHandler, "admin"))
	r.HandleAdvance("deleteUser", ms.RBAC.Middleware(ah.UserAdvanceHandlers.DeleteUserHandler, "admin"))

	r.HandleInspect("crowdfunding", ih.CrowdfundingInspectHandlers.FindAllCrowdfundingsHandler)
	r.HandleInspect("crowdfunding/{id}", ih.CrowdfundingInspectHandlers.FindCrowdfundingByIdHandler)

	r.HandleInspect("order", ih.OrderInspectHandlers.FindAllOrdersHandler)
	r.HandleInspect("order/{id}", ih.OrderInspectHandlers.FindOrderByIdHandler)
	r.HandleInspect("order/crowdfunding/{id}", ih.OrderInspectHandlers.FindBisdByCrowdfundingIdHandler)

	r.HandleInspect("contract", ih.ContractInspectHandlers.FindAllContractsHandler)
	r.HandleInspect("contract/{symbol}", ih.ContractInspectHandlers.FindContractBySymbolHandler)

	r.HandleInspect("user", ih.UserInspectHandlers.FindAllUsersHandler)
	r.HandleInspect("user/{address}", ih.UserInspectHandlers.FindUserByAddressHandler)
	r.HandleInspect("balance/{symbol}/{address}", ih.UserInspectHandlers.BalanceHandler)

	return r
}
