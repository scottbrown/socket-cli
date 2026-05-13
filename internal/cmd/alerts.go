package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newAlertsCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "View and manage alerts",
	}

	cmd.AddCommand(
		newAlertsListCmd(getClient),
		newAlertsTriageListCmd(getClient),
	)
	return cmd
}

func newAlertsListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org     string
		perPage int
		page    int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alerts for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			q := url.Values{}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if page > 0 {
				q.Set("page", fmt.Sprintf("%d", page))
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/alerts", org), q)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Results per page")
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.MarkFlagRequired("org")

	return cmd
}

func newAlertsTriageListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org string

	cmd := &cobra.Command{
		Use:   "triage",
		Short: "List triaged alerts",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Get(fmt.Sprintf("/orgs/%s/triage/alerts", org), nil)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.MarkFlagRequired("org")

	return cmd
}
