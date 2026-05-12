package cmd

import (
	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newOrgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Manage organizations",
	}

	cmd.AddCommand(newOrgsListCmd())
	return cmd
}

func newOrgsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List organizations linked to your API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient()
			data, err := c.Get("/organizations", nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(data)
		},
	}
}
