package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestReposListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/repos" {
				t.Errorf("path = %q, want /orgs/my-org/repos", path)
			}
			return []byte(`{"results":[{"name":"repo1"}]}`), nil
		},
	}

	out, err := executeCommand(t, mock, "repos", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "repo1") {
		t.Errorf("output = %q, want to contain 'repo1'", out)
	}
}

func TestReposListCmd_WithPagination(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("per_page") != "5" {
				t.Errorf("per_page = %q, want 5", query.Get("per_page"))
			}
			if query.Get("page") != "2" {
				t.Errorf("page = %q, want 2", query.Get("page"))
			}
			return []byte(`{"results":[]}`), nil
		},
	}

	_, err := executeCommand(t, mock, "repos", "list", "--org", "x", "--per-page", "5", "--page", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReposListCmd_MissingOrg(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "repos", "list")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestReposGetCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/repos/repo1" {
				t.Errorf("path = %q, want /orgs/org1/repos/repo1", path)
			}
			return []byte(`{"name":"repo1"}`), nil
		},
	}

	out, err := executeCommand(t, mock, "repos", "get", "--org", "org1", "--repo", "repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "repo1") {
		t.Errorf("output = %q, want to contain 'repo1'", out)
	}
}

func TestReposDeleteCmd_Success(t *testing.T) {
	mock := &mockClient{
		deleteFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/repos/repo1" {
				t.Errorf("path = %q, want /orgs/org1/repos/repo1", path)
			}
			return []byte(`{"deleted":true}`), nil
		},
	}

	_, err := executeCommand(t, mock, "repos", "delete", "--org", "org1", "--repo", "repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReposDeleteCmd_APIError(t *testing.T) {
	mock := &mockClient{
		deleteFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 403: forbidden")
		},
	}

	_, err := executeCommand(t, mock, "repos", "delete", "--org", "org1", "--repo", "repo1")
	if err == nil {
		t.Fatal("expected error")
	}
}
