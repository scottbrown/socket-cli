package cmd

import (
	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/scottbrown/socket-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return newRootCmd(nil)
}

func newRootCmd(injectedClient api.SocketAPI) *cobra.Command {
	var client api.SocketAPI

	root := &cobra.Command{
		Use:   "socket",
		Short: "CLI for the Socket.dev API",
		Long:  "A command-line interface for interacting with the Socket.dev supply chain security platform.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if injectedClient != nil {
				client = injectedClient
				return nil
			}
			token, err := config.GetAPIToken()
			if err != nil {
				return err
			}
			client = api.NewClient(config.GetBaseURL(), token)
			return nil
		},
	}

	getClient := func() api.SocketAPI { return client }

	root.AddCommand(
		newOrgsCmd(getClient),
		newReposCmd(getClient),
		newFullScansCmd(getClient),
		newDiffScansCmd(getClient),
		newAlertsCmd(getClient),
		newFixesCmd(getClient),
		newQuotaCmd(getClient),
		newPackagesCmd(getClient),
	)

	return root
}
