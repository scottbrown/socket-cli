package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestFixesListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/fixes" {
				t.Errorf("path = %q, want /orgs/my-org/fixes", path)
			}
			if query.Get("vulnerability_ids") != "*" {
				t.Errorf("vulnerability_ids = %q, want *", query.Get("vulnerability_ids"))
			}
			if query.Get("allow_major_updates") != "false" {
				t.Errorf("allow_major_updates = %q, want false", query.Get("allow_major_updates"))
			}
			return []byte(`{"fixDetails":{}}`), nil
		},
	}

	out, err := executeCommand(t, mock, "fixes", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "fixDetails") {
		t.Errorf("output = %q, want to contain 'fixDetails'", out)
	}
}

func TestFixesListCmd_WithRepo(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("repo_slug") != "my-repo" {
				t.Errorf("repo_slug = %q, want my-repo", query.Get("repo_slug"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "fixes", "list", "--org", "x", "--repo", "my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFixesListCmd_WithScanID(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("full_scan_id") != "scan-123" {
				t.Errorf("full_scan_id = %q, want scan-123", query.Get("full_scan_id"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "fixes", "list", "--org", "x", "--scan-id", "scan-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFixesListCmd_WithAllFlags(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("vulnerability_ids") != "GHSA-1234,CVE-2024-5678" {
				t.Errorf("vulnerability_ids = %q", query.Get("vulnerability_ids"))
			}
			if query.Get("allow_major_updates") != "true" {
				t.Errorf("allow_major_updates = %q, want true", query.Get("allow_major_updates"))
			}
			if query.Get("minimum_release_age") != "1w" {
				t.Errorf("minimum_release_age = %q, want 1w", query.Get("minimum_release_age"))
			}
			if query.Get("include_details") != "true" {
				t.Errorf("include_details = %q, want true", query.Get("include_details"))
			}
			if query.Get("include_responsible_direct_dependencies") != "true" {
				t.Errorf("include_responsible = %q, want true", query.Get("include_responsible_direct_dependencies"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "fixes", "list", "--org", "x",
		"--vuln-ids", "GHSA-1234,CVE-2024-5678",
		"--allow-major",
		"--min-age", "1w",
		"--details",
		"--responsible",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFixesListCmd_MissingOrg(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "fixes", "list")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestFixesListCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 404: not found")
		},
	}

	_, err := executeCommand(t, mock, "fixes", "list", "--org", "x", "--repo", "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}
