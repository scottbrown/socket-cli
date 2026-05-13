package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newFixesCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fixes",
		Short: "Get fix recommendations for vulnerabilities",
	}

	cmd.AddCommand(newFixesListCmd(getClient))
	return cmd
}

func newFixesListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org              string
		repo             string
		scanID           string
		vulnIDs          string
		allowMajor       bool
		minReleaseAge    string
		includeDetails   bool
		includeResponsible bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Fetch fix recommendations for vulnerabilities",
		Long: `Fetch recommended package upgrades to resolve vulnerabilities.

Specify a target using --repo (latest scan on default branch), --scan-id
(a specific scan), or neither (requires vulnerability IDs).

Use --vuln-ids "*" to get fixes for all detected vulnerabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			q := url.Values{}
			q.Set("vulnerability_ids", vulnIDs)
			q.Set("allow_major_updates", fmt.Sprintf("%t", allowMajor))
			if repo != "" {
				q.Set("repo_slug", repo)
			}
			if scanID != "" {
				q.Set("full_scan_id", scanID)
			}
			if minReleaseAge != "" {
				q.Set("minimum_release_age", minReleaseAge)
			}
			if includeDetails {
				q.Set("include_details", "true")
			}
			if includeResponsible {
				q.Set("include_responsible_direct_dependencies", "true")
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/fixes", org), q)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository slug (uses latest scan on default branch)")
	cmd.Flags().StringVar(&scanID, "scan-id", "", "Full scan ID to get fixes for")
	cmd.Flags().StringVar(&vulnIDs, "vuln-ids", "*", "Comma-separated GHSA/CVE IDs, or \"*\" for all")
	cmd.Flags().BoolVar(&allowMajor, "allow-major", false, "Allow major version updates in fixes")
	cmd.Flags().StringVar(&minReleaseAge, "min-age", "", "Minimum release age (e.g. \"1h\", \"2d\", \"1w\")")
	cmd.Flags().BoolVar(&includeDetails, "details", false, "Include advisory details in response")
	cmd.Flags().BoolVar(&includeResponsible, "responsible", false, "Include direct dependencies responsible for the vulnerability")
	cmd.MarkFlagRequired("org")

	return cmd
}
