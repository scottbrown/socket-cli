package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestDiffScansListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/diff-scans" {
				t.Errorf("path = %q, want /orgs/my-org/diff-scans", path)
			}
			return []byte(`{"results":[]}`), nil
		},
	}

	_, err := executeCommand(t, mock, "diffscans", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffScansListCmd_Pagination(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("per_page") != "10" {
				t.Errorf("per_page = %q, want 10", query.Get("per_page"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "diffscans", "list", "--org", "x", "--per-page", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffScansGetCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/diff-scans/scan-123" {
				t.Errorf("path = %q, want /orgs/org1/diff-scans/scan-123", path)
			}
			return []byte(`{"id":"scan-123"}`), nil
		},
	}

	out, err := executeCommand(t, mock, "diffscans", "get", "--org", "org1", "--id", "scan-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "scan-123") {
		t.Errorf("output = %q, want to contain 'scan-123'", out)
	}
}

func TestDiffScansDeleteCmd_Success(t *testing.T) {
	mock := &mockClient{
		deleteFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/diff-scans/scan-456" {
				t.Errorf("path = %q", path)
			}
			return []byte(`{"deleted":true}`), nil
		},
	}

	_, err := executeCommand(t, mock, "diffscans", "delete", "--org", "org1", "--id", "scan-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffScansDeleteCmd_APIError(t *testing.T) {
	mock := &mockClient{
		deleteFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 404: not found")
		},
	}

	_, err := executeCommand(t, mock, "diffscans", "delete", "--org", "org1", "--id", "bad")
	if err == nil {
		t.Fatal("expected error")
	}
}
