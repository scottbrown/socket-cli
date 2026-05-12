package cmd

import (
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newReposCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "Manage repositories",
	}

	cmd.AddCommand(
		newReposListCmd(),
		newReposGetCmd(),
		newReposDeleteCmd(),
	)
	return cmd
}

func newReposListCmd() *cobra.Command {
	var (
		org     string
		sort    string
		dir     string
		perPage int
		page    int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repositories for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient()
			q := url.Values{}
			if sort != "" {
				q.Set("sort", sort)
			}
			if dir != "" {
				q.Set("direction", dir)
			}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if page > 0 {
				q.Set("page", fmt.Sprintf("%d", page))
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/repos", org), q)
			if err != nil {
				return err
			}
			return api.PrintJSON(data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort by: name or created_at")
	cmd.Flags().StringVar(&dir, "direction", "", "Sort direction: asc or desc")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Results per page (1-100)")
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.MarkFlagRequired("org")

	return cmd
}

func newReposGetCmd() *cobra.Command {
	var org, repo string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get repository details",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient()
			data, err := c.Get(fmt.Sprintf("/orgs/%s/repos/%s", org, repo), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository slug (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("repo")

	return cmd
}

func newReposDeleteCmd() *cobra.Command {
	var org, repo string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient()
			data, err := c.Delete(fmt.Sprintf("/orgs/%s/repos/%s", org, repo), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository slug (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("repo")

	return cmd
}
