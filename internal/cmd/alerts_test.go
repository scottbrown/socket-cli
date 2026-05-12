package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestAlertsListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/alerts" {
				t.Errorf("path = %q, want /orgs/my-org/alerts", path)
			}
			return []byte(`{"items":[{"id":"ALERT-1","severity":"high"}]}`), nil
		},
	}

	out, err := executeCommand(t, mock, "alerts", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ALERT-1") {
		t.Errorf("output = %q, want to contain 'ALERT-1'", out)
	}
}

func TestAlertsListCmd_Pagination(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("per_page") != "3" {
				t.Errorf("per_page = %q, want 3", query.Get("per_page"))
			}
			if query.Get("page") != "2" {
				t.Errorf("page = %q, want 2", query.Get("page"))
			}
			return []byte(`{"items":[]}`), nil
		},
	}

	_, err := executeCommand(t, mock, "alerts", "list", "--org", "x", "--per-page", "3", "--page", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlertsListCmd_MissingOrg(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "alerts", "list")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestAlertsTriageCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/triage/alerts" {
				t.Errorf("path = %q, want /orgs/my-org/triage/alerts", path)
			}
			return []byte(`{"triaged":[]}`), nil
		},
	}

	_, err := executeCommand(t, mock, "alerts", "triage", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAlertsTriageCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 500: internal")
		},
	}

	_, err := executeCommand(t, mock, "alerts", "triage", "--org", "x")
	if err == nil {
		t.Fatal("expected error")
	}
}
