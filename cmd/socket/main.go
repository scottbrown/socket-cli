package main

import (
	"os"

	"github.com/scottbrown/socket-cli/internal/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
