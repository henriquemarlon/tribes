package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/rollmelette/rollmelette"
)

func main() {
	//////////////////////// Setup DApp /////////////////////////
	app := NewDApp()

	ctx := context.Background()
	opts := rollmelette.NewRunOpts()
	if rollupUrl, isSet := os.LookupEnv("ROLLUP_HTTP_SERVER_URL"); isSet {
		opts.RollupURL = rollupUrl
	}
	if err := rollmelette.Run(ctx, opts, app); err != nil {
		slog.Error("application error", "error", err)
	}
}
