package cmd

import (
	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newOrgsCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Manage organizations",
	}

	cmd.AddCommand(newOrgsListCmd(getClient))
	return cmd
}

func newOrgsListCmd(getClient func() api.SocketAPI) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List organizations linked to your API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Get("/organizations", nil)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}
}
