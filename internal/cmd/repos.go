package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

func newReposCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "Manage repositories",
	}

	cmd.AddCommand(
		newReposListCmd(getClient),
		newReposGetCmd(getClient),
		newReposCreateCmd(getClient),
		newReposUpdateCmd(getClient),
		newReposDeleteCmd(getClient),
	)
	return cmd
}

func newReposListCmd(getClient func() api.SocketAPI) *cobra.Command {
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
			c := getClient()
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
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
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

func newReposGetCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, repo string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get repository details",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Get(fmt.Sprintf("/orgs/%s/repos/%s", org, repo), nil)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository slug (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("repo")

	return cmd
}

func newReposCreateCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org           string
		name          string
		description   string
		homepage      string
		visibility    string
		defaultBranch string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			body := map[string]string{"name": name}
			if description != "" {
				body["description"] = description
			}
			if homepage != "" {
				body["homepage"] = homepage
			}
			if visibility != "" {
				body["visibility"] = visibility
			}
			if defaultBranch != "" {
				body["default_branch"] = defaultBranch
			}
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}
			data, err := c.Post(fmt.Sprintf("/orgs/%s/repos", org), nil, bytes.NewReader(b), "application/json")
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&name, "name", "", "Repository name (required)")
	cmd.Flags().StringVar(&description, "description", "", "Repository description")
	cmd.Flags().StringVar(&homepage, "homepage", "", "Repository homepage URL")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Repository visibility")
	cmd.Flags().StringVar(&defaultBranch, "default-branch", "", "Default branch name")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newReposUpdateCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org           string
		repo          string
		name          string
		description   string
		homepage      string
		visibility    string
		defaultBranch string
		archived      bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			body := map[string]interface{}{}
			if cmd.Flags().Changed("name") {
				body["name"] = name
			}
			if cmd.Flags().Changed("description") {
				body["description"] = description
			}
			if cmd.Flags().Changed("homepage") {
				body["homepage"] = homepage
			}
			if cmd.Flags().Changed("visibility") {
				body["visibility"] = visibility
			}
			if cmd.Flags().Changed("default-branch") {
				body["default_branch"] = defaultBranch
			}
			if cmd.Flags().Changed("archived") {
				body["archived"] = archived
			}
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}
			data, err := c.Post(fmt.Sprintf("/orgs/%s/repos/%s", org, repo), nil, bytes.NewReader(b), "application/json")
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository slug (required)")
	cmd.Flags().StringVar(&name, "name", "", "New repository name")
	cmd.Flags().StringVar(&description, "description", "", "Repository description")
	cmd.Flags().StringVar(&homepage, "homepage", "", "Repository homepage URL")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Repository visibility")
	cmd.Flags().StringVar(&defaultBranch, "default-branch", "", "Default branch name")
	cmd.Flags().BoolVar(&archived, "archived", false, "Archive the repository")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("repo")

	return cmd
}

func newReposDeleteCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, repo string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Delete(fmt.Sprintf("/orgs/%s/repos/%s", org, repo), nil)
			if err != nil {
				return err
			}
			return api.PrintOutput(cmd.OutOrStdout(), data, getFormat(cmd))
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository slug (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("repo")

	return cmd
}
