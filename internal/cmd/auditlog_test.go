package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestAuditLogListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/audit-log" {
				t.Errorf("path = %q, want /orgs/my-org/audit-log", path)
			}
			return []byte(`{"events":[{"type":"repo.created"}]}`), nil
		},
	}

	out, err := executeCommand(t, mock, "auditlog", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "repo.created") {
		t.Errorf("output = %q, want to contain 'repo.created'", out)
	}
}

func TestAuditLogListCmd_WithAllFlags(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("type") != "repo.created" {
				t.Errorf("type = %q, want repo.created", query.Get("type"))
			}
			if query.Get("per_page") != "20" {
				t.Errorf("per_page = %q, want 20", query.Get("per_page"))
			}
			if query.Get("page") != "3" {
				t.Errorf("page = %q, want 3", query.Get("page"))
			}
			if query.Get("from") != "2026-01-01T00:00:00Z" {
				t.Errorf("from = %q, want 2026-01-01T00:00:00Z", query.Get("from"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "auditlog", "list", "--org", "x",
		"--type", "repo.created",
		"--per-page", "20",
		"--page", "3",
		"--from", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuditLogListCmd_MissingOrg(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "auditlog", "list")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestAuditLogListCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 403: forbidden")
		},
	}

	_, err := executeCommand(t, mock, "auditlog", "list", "--org", "x")
	if err == nil {
		t.Fatal("expected error")
	}
}
