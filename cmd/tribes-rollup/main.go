package main

import (
	"os"
	"github.com/tribeshq/tribes/cmd/tribes-rollup/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
