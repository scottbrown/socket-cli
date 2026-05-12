package cmd

import (
	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newQuotaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "quota",
		Short: "Show API quota information",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient()
			data, err := c.Get("/quota", nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(data)
		},
	}
}
