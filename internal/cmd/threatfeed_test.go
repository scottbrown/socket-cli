package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestThreatFeedListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/threat-feed" {
				t.Errorf("path = %q, want /orgs/my-org/threat-feed", path)
			}
			return []byte(`{"items":[{"name":"malware-pkg"}]}`), nil
		},
	}

	out, err := executeCommand(t, mock, "threatfeed", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "malware-pkg") {
		t.Errorf("output = %q, want to contain 'malware-pkg'", out)
	}
}

func TestThreatFeedListCmd_WithAllFlags(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("per_page") != "10" {
				t.Errorf("per_page = %q, want 10", query.Get("per_page"))
			}
			if query.Get("cursor") != "abc123" {
				t.Errorf("cursor = %q, want abc123", query.Get("cursor"))
			}
			if query.Get("sort") != "updated_at" {
				t.Errorf("sort = %q, want updated_at", query.Get("sort"))
			}
			if query.Get("direction") != "desc" {
				t.Errorf("direction = %q, want desc", query.Get("direction"))
			}
			if query.Get("filter") != "mal" {
				t.Errorf("filter = %q, want mal", query.Get("filter"))
			}
			if query.Get("ecosystem") != "npm" {
				t.Errorf("ecosystem = %q, want npm", query.Get("ecosystem"))
			}
			if query.Get("name") != "evil-pkg" {
				t.Errorf("name = %q, want evil-pkg", query.Get("name"))
			}
			if query.Get("version") != "1.0" {
				t.Errorf("version = %q, want 1.0", query.Get("version"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "threatfeed", "list", "--org", "x",
		"--per-page", "10",
		"--cursor", "abc123",
		"--sort", "updated_at",
		"--direction", "desc",
		"--filter", "mal",
		"--ecosystem", "npm",
		"--name", "evil-pkg",
		"--version", "1.0",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestThreatFeedListCmd_MissingOrg(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "threatfeed", "list")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestThreatFeedListCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 500: internal error")
		},
	}

	_, err := executeCommand(t, mock, "threatfeed", "list", "--org", "x")
	if err == nil {
		t.Fatal("expected error")
	}
}
