package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestAnalyticsGetCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/analytics/org/30" {
				t.Errorf("path = %q, want /analytics/org/30", path)
			}
			if query.Get("org") != "my-org" {
				t.Errorf("org = %q, want my-org", query.Get("org"))
			}
			return []byte(`{"alerts_count":42}`), nil
		},
	}

	out, err := executeCommand(t, mock, "analytics", "get", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "alerts_count") {
		t.Errorf("output = %q, want to contain 'alerts_count'", out)
	}
}

func TestAnalyticsGetCmd_WithDays(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/analytics/org/7" {
				t.Errorf("path = %q, want /analytics/org/7", path)
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "analytics", "get", "--org", "my-org", "--days", "7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnalyticsGetCmd_WithRepo(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("repo") != "my-repo" {
				t.Errorf("repo = %q, want my-repo", query.Get("repo"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "analytics", "get", "--org", "x", "--repo", "my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnalyticsGetCmd_MissingOrg(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "analytics", "get")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestAnalyticsGetCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 404: not found")
		},
	}

	_, err := executeCommand(t, mock, "analytics", "get", "--org", "x")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAnalyticsGetCmd_MarkdownFormat(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return []byte(`{"alerts_count":42,"org":"my-org"}`), nil
		},
	}

	out, err := executeCommand(t, mock, "analytics", "get", "--org", "my-org", "--format", "markdown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "**alerts_count**") {
		t.Errorf("expected markdown bold key, got: %q", out)
	}
	if strings.Contains(out, "{") {
		t.Errorf("expected non-JSON output for markdown format, got: %q", out)
	}
}
