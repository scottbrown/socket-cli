package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newPackagesCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages",
		Short: "Query package information",
	}

	cmd.AddCommand(
		newPackagesLookupCmd(getClient),
		newPackagesShallowCmd(getClient),
	)
	return cmd
}

func newPackagesLookupCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org string

	cmd := &cobra.Command{
		Use:   "lookup [purl...]",
		Short: "Look up packages by Package URL (purl)",
		Long:  "Look up package information using Package URLs. Example purls: pkg:npm/express@4.18.2, pkg:pypi/requests@2.31.0",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()

			body := map[string][]string{"purls": args}
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}

			path := "/purl"
			if org != "" {
				path = fmt.Sprintf("/orgs/%s/purl", org)
			}

			data, err := c.Post(path, nil, bytes.NewReader(b), "application/json")
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (uses org-specific scoring if set)")

	return cmd
}

func newPackagesShallowCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org string

	cmd := &cobra.Command{
		Use:   "shallow [purl...]",
		Short: "Look up packages without transitive dependencies",
		Long:  "Look up package information using Package URLs without resolving transitive dependencies. Faster but less complete than lookup.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()

			body := map[string]interface{}{
				"purls":            args,
				"include_transitives": false,
			}
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}

			path := "/purl"
			if org != "" {
				path = fmt.Sprintf("/orgs/%s/purl", org)
			}

			data, err := c.Post(path, nil, bytes.NewReader(b), "application/json")
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (uses org-specific scoring if set)")

	return cmd
}
