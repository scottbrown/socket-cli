package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/scottbrown/socket-cli/internal/api"
	"github.com/spf13/cobra"
)

type FileOpener func(name string) (io.ReadCloser, error)

func defaultFileOpener(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func newFullScansCmd(getClient func() api.SocketAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fullscans",
		Short: "Manage full scans",
	}

	cmd.AddCommand(
		newFullScansListCmd(getClient),
		newFullScansGetCmd(getClient),
		newFullScansCreateCmd(getClient, defaultFileOpener),
		newFullScansDeleteCmd(getClient),
		newFullScansMetadataCmd(getClient),
	)
	return cmd
}

func newFullScansListCmd(getClient func() api.SocketAPI) *cobra.Command {
	var (
		org      string
		sort     string
		dir      string
		perPage  int
		page     int
		repo     string
		branch   string
		from     string
		scanType string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List full scans for an organization",
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
			if repo != "" {
				q.Set("repo", repo)
			}
			if branch != "" {
				q.Set("branch", branch)
			}
			if from != "" {
				q.Set("from", from)
			}
			if scanType != "" {
				q.Set("scan_type", scanType)
			}
			data, err := c.Get(fmt.Sprintf("/orgs/%s/full-scans", org), q)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort by: name or created_at")
	cmd.Flags().StringVar(&dir, "direction", "", "Sort direction: asc or desc")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "Results per page (1-100)")
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().StringVar(&repo, "repo", "", "Filter by repository")
	cmd.Flags().StringVar(&branch, "branch", "", "Filter by branch")
	cmd.Flags().StringVar(&from, "from", "", "Filter by Unix timestamp")
	cmd.Flags().StringVar(&scanType, "scan-type", "", "Filter by scan type")
	cmd.MarkFlagRequired("org")

	return cmd
}

func newFullScansGetCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, scanID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a full scan by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Get(fmt.Sprintf("/orgs/%s/full-scans/%s", org, scanID), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&scanID, "id", "", "Full scan ID (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newFullScansCreateCmd(getClient func() api.SocketAPI, openFile FileOpener) *cobra.Command {
	var (
		org           string
		repo          string
		branch        string
		commitHash    string
		commitMessage string
		pullRequest   int
		files         []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new full scan by uploading manifest files",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()

			q := url.Values{}
			q.Set("repo", repo)
			if branch != "" {
				q.Set("branch", branch)
			}
			if commitHash != "" {
				q.Set("commit_hash", commitHash)
			}
			if commitMessage != "" {
				q.Set("commit_message", commitMessage)
			}
			if pullRequest > 0 {
				q.Set("pull_request", fmt.Sprintf("%d", pullRequest))
			}

			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			for _, f := range files {
				file, err := openFile(f)
				if err != nil {
					return fmt.Errorf("opening %s: %w", f, err)
				}
				part, err := writer.CreateFormFile(filepath.Base(f), filepath.Base(f))
				if err != nil {
					file.Close()
					return err
				}
				if _, err := io.Copy(part, file); err != nil {
					file.Close()
					return err
				}
				file.Close()
			}
			writer.Close()

			data, err := c.Post(
				fmt.Sprintf("/orgs/%s/full-scans", org),
				q, &body, writer.FormDataContentType(),
			)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "Repository name (required)")
	cmd.Flags().StringVar(&branch, "branch", "", "Git branch")
	cmd.Flags().StringVar(&commitHash, "commit", "", "Commit hash")
	cmd.Flags().StringVar(&commitMessage, "message", "", "Commit message")
	cmd.Flags().IntVar(&pullRequest, "pr", 0, "Pull request number")
	cmd.Flags().StringSliceVar(&files, "file", nil, "Manifest files to upload (required, repeatable)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("repo")
	cmd.MarkFlagRequired("file")

	return cmd
}

func newFullScansDeleteCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, scanID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a full scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Delete(fmt.Sprintf("/orgs/%s/full-scans/%s", org, scanID), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&scanID, "id", "", "Full scan ID (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newFullScansMetadataCmd(getClient func() api.SocketAPI) *cobra.Command {
	var org, scanID string

	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Get metadata for a full scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getClient()
			data, err := c.Get(fmt.Sprintf("/orgs/%s/full-scans/%s/metadata", org, scanID), nil)
			if err != nil {
				return err
			}
			return api.PrintJSON(cmd.OutOrStdout(), data)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "Organization slug (required)")
	cmd.Flags().StringVar(&scanID, "id", "", "Full scan ID (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("id")

	return cmd
}
