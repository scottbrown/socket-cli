package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newThreatFeedCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threatfeed",
		Short: "Access the threat feed",
	}

	cmd.AddCommand(newThreatFeedListCmd(getClient))
	return cmd
}

func newThreatFeedListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org       string
		perPage   int
		cursor    string
		sort      string
		direction string
		filter    string
		ecosystem string
		name      string
		version   string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List threat feed items",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			q := url.Values{}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			if sort != "" {
				q.Set("sort", sort)
			}
			if direction != "" {
				q.Set("direction", direction)
			}
			if filter != "" {
				q.Set("filter", filter)
			}
			if ecosystem != "" {
				q.Set("ecosystem", ecosystem)
			}
			if name != "" {
				q.Set("name", name)
			}
			if version != "" {
				q.Set("version", version)
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/threat-feed", org), q)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Results per page")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort by: updated_at or created_at")
	cmd.Flags().StringVar(&direction, "direction", "", "Sort direction: asc or desc")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter type: mal, typo, joke, spy, vuln")
	cmd.Flags().StringVar(&ecosystem, "ecosystem", "", "Ecosystem: npm, pypi, golang, maven, gem, nuget")
	cmd.Flags().StringVar(&name, "name", "", "Package name filter")
	cmd.Flags().StringVar(&version, "version", "", "Version prefix filter")
	cmd.MarkFlagRequired("org")

	return cmd
}
