package cmd

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/scottbrown/socket-cli/internal/api"
)

func TestFullScansListCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/my-org/full-scans" {
				t.Errorf("path = %q, want /orgs/my-org/full-scans", path)
			}
			return []byte(`{"results":[]}`), nil
		},
	}

	_, err := executeCommand(t, mock, "fullscans", "list", "--org", "my-org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFullScansListCmd_WithFilters(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if query.Get("repo") != "my-repo" {
				t.Errorf("repo = %q, want my-repo", query.Get("repo"))
			}
			if query.Get("branch") != "main" {
				t.Errorf("branch = %q, want main", query.Get("branch"))
			}
			if query.Get("sort") != "created_at" {
				t.Errorf("sort = %q, want created_at", query.Get("sort"))
			}
			if query.Get("direction") != "asc" {
				t.Errorf("direction = %q, want asc", query.Get("direction"))
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "fullscans", "list", "--org", "x",
		"--repo", "my-repo", "--branch", "main", "--sort", "created_at", "--direction", "asc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFullScansGetCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/full-scans/scan-abc" {
				t.Errorf("path = %q", path)
			}
			return []byte(`{"id":"scan-abc"}`), nil
		},
	}

	out, err := executeCommand(t, mock, "fullscans", "get", "--org", "org1", "--id", "scan-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "scan-abc") {
		t.Errorf("output = %q", out)
	}
}

func TestFullScansDeleteCmd_Success(t *testing.T) {
	mock := &mockClient{
		deleteFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/full-scans/scan-xyz" {
				t.Errorf("path = %q", path)
			}
			return []byte(`{"deleted":true}`), nil
		},
	}

	_, err := executeCommand(t, mock, "fullscans", "delete", "--org", "org1", "--id", "scan-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFullScansMetadataCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/full-scans/scan-1/metadata" {
				t.Errorf("path = %q", path)
			}
			return []byte(`{"scan_type":"default"}`), nil
		},
	}

	out, err := executeCommand(t, mock, "fullscans", "metadata", "--org", "org1", "--id", "scan-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "scan_type") {
		t.Errorf("output = %q", out)
	}
}

func TestFullScansCreateCmd_Success(t *testing.T) {
	var capturedPath string
	var capturedQuery url.Values
	var capturedContentType string
	var capturedBody []byte

	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			capturedPath = path
			capturedQuery = query
			capturedContentType = contentType
			capturedBody, _ = io.ReadAll(body)
			return []byte(`{"id":"new-scan"}`), nil
		},
	}

	opener := func(name string) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(`{"name":"test-package"}`)), nil
	}

	getClient := func() api.SocketAPI { return mock }
	cmd := newFullScansCreateCmd(getClient, opener)

	root := newRootCmd(mock)
	root.RemoveCommand(root.Commands()[2]) // remove default fullscans
	fsCmd := newFullScansCmd(getClient)
	fsCmd.RemoveCommand(fsCmd.Commands()[2]) // remove default create
	fsCmd.AddCommand(cmd)
	root.AddCommand(fsCmd)

	// Simpler approach: test the create command directly
	cmd2 := newFullScansCreateCmd(getClient, opener)
	cmd2.SetArgs([]string{"--org", "org1", "--repo", "my-repo", "--branch", "main", "--file", "package.json"})
	cmd2.SetOut(io.Discard)
	err := cmd2.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedPath != "/orgs/org1/full-scans" {
		t.Errorf("path = %q, want /orgs/org1/full-scans", capturedPath)
	}
	if capturedQuery.Get("repo") != "my-repo" {
		t.Errorf("repo param = %q", capturedQuery.Get("repo"))
	}
	if capturedQuery.Get("branch") != "main" {
		t.Errorf("branch param = %q", capturedQuery.Get("branch"))
	}
	if !strings.Contains(capturedContentType, "multipart/form-data") {
		t.Errorf("content-type = %q, want multipart/form-data", capturedContentType)
	}
	if !strings.Contains(string(capturedBody), "test-package") {
		t.Errorf("body doesn't contain file content")
	}
}

func TestFullScansReportCmd_Success(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			if path != "/orgs/org1/full-scans/scan-1/stream" {
				t.Errorf("path = %q, want /orgs/org1/full-scans/scan-1/stream", path)
			}
			return []byte(`{"policy_result":"pass","alerts":[]}`), nil
		},
	}

	out, err := executeCommand(t, mock, "fullscans", "report", "--org", "org1", "--id", "scan-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "policy_result") {
		t.Errorf("output = %q, want to contain 'policy_result'", out)
	}
}

func TestFullScansReportCmd_MissingFlags(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "fullscans", "report", "--org", "org1")
	if err == nil {
		t.Fatal("expected error for missing --id flag")
	}
}

func TestFullScansReportCmd_APIError(t *testing.T) {
	mock := &mockClient{
		getFunc: func(path string, query url.Values) ([]byte, error) {
			return nil, fmt.Errorf("API error 404: scan not found")
		},
	}

	_, err := executeCommand(t, mock, "fullscans", "report", "--org", "org1", "--id", "bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFullScansCreateCmd_FileOpenError(t *testing.T) {
	mock := &mockClient{}

	opener := func(name string) (io.ReadCloser, error) {
		return nil, fmt.Errorf("file not found: %s", name)
	}

	getClient := func() api.SocketAPI { return mock }
	cmd := newFullScansCreateCmd(getClient, opener)
	cmd.SetArgs([]string{"--org", "org1", "--repo", "my-repo", "--file", "missing.json"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("error = %q, want to contain 'file not found'", err.Error())
	}
}
