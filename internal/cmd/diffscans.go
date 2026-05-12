package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newDiffScansCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diffscans",
		Short: "Manage diff scans",
	}

	cmd.AddCommand(
		newDiffScansListCmd(getClient),
		newDiffScansGetCmd(getClient),
		newDiffScansDeleteCmd(getClient),
	)
	return cmd
}

func newDiffScansListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org     string
		perPage int
		page    int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List diff scans for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			q := url.Values{}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if page > 0 {
				q.Set("page", fmt.Sprintf("%d", page))
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/diff-scans", org), q)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Results per page (1-100)")
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.MarkFlagRequired("org")

	return cmd
}

func newDiffScansGetCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, scanID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a diff scan by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Get(fmt.Sprintf("/orgs/%s/diff-scans/%s", org, scanID), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&scanID, "id", "", "Diff scan ID (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newDiffScansDeleteCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, scanID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a diff scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Delete(fmt.Sprintf("/orgs/%s/diff-scans/%s", org, scanID), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&scanID, "id", "", "Diff scan ID (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("id")

	return cmd
}
