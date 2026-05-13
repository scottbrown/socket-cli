package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newAnalyticsCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "View organization analytics",
	}

	cmd.AddCommand(newAnalyticsGetCmd(getClient))
	return cmd
}

func newAnalyticsGetCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org    string
		repo   string
		days   int
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get analytics for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			q := url.Values{}
			if repo != "" {
				q.Set("repo", repo)
			}
			q.Set("org", org)

			path := fmt.Sprintf("/analytics/org/%d", days)

			data, err := c.Get(path, q)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Filter by repository")
	cmd.Flags().IntVar(&days, "days", 30, "Time range in days (7, 30, or 90)")
	cmd.MarkFlagRequired("org")

	return cmd
}
