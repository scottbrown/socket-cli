package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"
)

func TestPackagesLookupCmd_Success(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			if path != "/purl" {
				t.Errorf("path = %q, want /purl", path)
			}
			if contentType != "application/json" {
				t.Errorf("contentType = %q, want application/json", contentType)
			}
			b, _ := io.ReadAll(body)
			var payload map[string][]string
			if err := json.Unmarshal(b, &payload); err != nil {
				t.Fatalf("invalid body JSON: %v", err)
			}
			if len(payload["purls"]) != 1 || payload["purls"][0] != "pkg:npm/express@4.18.2" {
				t.Errorf("purls = %v", payload["purls"])
			}
			return []byte(`{"packages":[{"name":"express"}]}`), nil
		},
	}

	out, err := executeCommand(t, mock, "packages", "lookup", "pkg:npm/express@4.18.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "express") {
		t.Errorf("output = %q, want to contain 'express'", out)
	}
}

func TestPackagesLookupCmd_WithOrg(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			if path != "/orgs/my-org/purl" {
				t.Errorf("path = %q, want /orgs/my-org/purl", path)
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "packages", "lookup", "--org", "my-org", "pkg:npm/lodash@4.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPackagesLookupCmd_MultiplePurls(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			b, _ := io.ReadAll(body)
			var payload map[string][]string
			json.Unmarshal(b, &payload)
			if len(payload["purls"]) != 3 {
				t.Errorf("expected 3 purls, got %d", len(payload["purls"]))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "packages", "lookup", "pkg:npm/a@1", "pkg:npm/b@2", "pkg:npm/c@3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPackagesLookupCmd_NoPurls(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "packages", "lookup")
	if err == nil {
		t.Fatal("expected error for missing purl args")
	}
}

func TestPackagesLookupCmd_APIError(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			return nil, fmt.Errorf("API error 429: rate limited")
		},
	}

	_, err := executeCommand(t, mock, "packages", "lookup", "pkg:npm/x@1")
	if err == nil {
		t.Fatal("expected error")
	}
}
