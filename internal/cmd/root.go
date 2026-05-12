package cmd

import (
	"fmt"
	"os"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/scottbrown/socket-cli/internal/config"
	"github.com/spf13/cobra"
)

var client *api.Client

func newClient() *api.Client {
	if client != nil {
		return client
	}
	token, err := config.GetAPIToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	client = api.NewClient(token)
	return client
}

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "socket",
		Short: "CLI for the Socket.dev API",
		Long:  "A command-line interface for interacting with the Socket.dev supply chain security platform.",
	}

	root.AddCommand(
		newOrgsCmd(),
		newReposCmd(),
		newFullScansCmd(),
		newDiffScansCmd(),
		newAlertsCmd(),
		newQuotaCmd(),
		newPackagesCmd(),
	)

	return root
}
