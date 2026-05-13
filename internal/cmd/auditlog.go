package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newAuditLogCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auditlog",
		Short: "Access organization audit logs",
	}

	cmd.AddCommand(newAuditLogListCmd(getClient))
	return cmd
}

func newAuditLogListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org     string
		evType  string
		perPage int
		page    int
		from    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit log events",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			q := url.Values{}
			if evType != "" {
				q.Set("type", evType)
			}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if page > 0 {
				q.Set("page", fmt.Sprintf("%d", page))
			}
			if from != "" {
				q.Set("from", from)
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/audit-log", org), q)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&evType, "type", "", "Event type filter")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Results per page")
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().StringVar(&from, "from", "", "Filter events from this timestamp")
	cmd.MarkFlagRequired("org")

	return cmd
}
