package main

import (
	"os"

	"github.com/tribeshq/tribes/cmd/dapp/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}