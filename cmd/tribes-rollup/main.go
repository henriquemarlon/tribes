package main

import (
	"github.com/tribeshq/tribes/cmd/tribes-rollup/root"
	"os"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
