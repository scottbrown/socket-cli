package cmd

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestQuotaCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/quota" {
				t.Errorf("path = %q, want /quota", path)
			}
			return []byte(`{"maxQuota":1000000,"quota":999990}`), nil
		},
	}

	out, err := executeCommand(t, mock, "quota")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "999990") {
		t.Errorf("output = %q, want to contain '999990'", out)
	}
}

func TestQuotaCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 401: unauthorized")
		},
	}

	_, err := executeCommand(t, mock, "quota")
	if err == nil {
		t.Fatal("expected error")
	}
}
