package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestOrgsListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/organizations" {
				t.Errorf("path = %q, want /organizations", path)
			}
			return []byte(`{"organizations":{"1":{"name":"test-org"}}}`), nil
		},
	}

	out, err := executeCommand(t, mock, "orgs", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-org") {
		t.Errorf("output = %q, want to contain 'test-org'", out)
	}
}

func TestOrgsListCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 401: unauthorized")
		},
	}

	_, err := executeCommand(t, mock, "orgs", "list")
	if err == nil {
		t.Fatal("expected error")
	}
}
