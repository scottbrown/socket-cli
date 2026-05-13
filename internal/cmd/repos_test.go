package cmd

import (
	"encoding/json"
	"fmt"
	"io"
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

func TestReposCreateCmd_Success(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			if path != "/orgs/org1/repos" {
				t.Errorf("path = %q, want /orgs/org1/repos", path)
			}
			if contentType != "application/json" {
				t.Errorf("contentType = %q, want application/json", contentType)
			}
			var payload map[string]string
			if err := json.NewDecoder(body).Decode(&payload); err != nil {
				t.Fatalf("failed to decode body: %v", err)
			}
			if payload["name"] != "new-repo" {
				t.Errorf("name = %q, want new-repo", payload["name"])
			}
			if payload["description"] != "A test repo" {
				t.Errorf("description = %q, want 'A test repo'", payload["description"])
			}
			return []byte(`{"name":"new-repo"}`), nil
		},
	}

	out, err := executeCommand(t, mock, "repos", "create", "--org", "org1", "--name", "new-repo", "--description", "A test repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "new-repo") {
		t.Errorf("output = %q, want to contain 'new-repo'", out)
	}
}

func TestReposCreateCmd_MissingFlags(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "repos", "create", "--org", "org1")
	if err == nil {
		t.Fatal("expected error for missing --name flag")
	}

	_, err = executeCommand(t, mock, "repos", "create", "--name", "x")
	if err == nil {
		t.Fatal("expected error for missing --org flag")
	}
}

func TestReposCreateCmd_APIError(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			return nil, fmt.Errorf("API error 409: conflict")
		},
	}

	_, err := executeCommand(t, mock, "repos", "create", "--org", "org1", "--name", "dup")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReposUpdateCmd_Success(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			if path != "/orgs/org1/repos/repo1" {
				t.Errorf("path = %q, want /orgs/org1/repos/repo1", path)
			}
			var payload map[string]interface{}
			if err := json.NewDecoder(body).Decode(&payload); err != nil {
				t.Fatalf("failed to decode body: %v", err)
			}
			if payload["description"] != "updated desc" {
				t.Errorf("description = %q, want 'updated desc'", payload["description"])
			}
			if payload["archived"] != true {
				t.Errorf("archived = %v, want true", payload["archived"])
			}
			return []byte(`{"name":"repo1","archived":true}`), nil
		},
	}

	out, err := executeCommand(t, mock, "repos", "update", "--org", "org1", "--repo", "repo1",
		"--description", "updated desc", "--archived")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "repo1") {
		t.Errorf("output = %q, want to contain 'repo1'", out)
	}
}

func TestReposUpdateCmd_MissingFlags(t *testing.T) {
	mock := &mockClient{}
	_, err := executeCommand(t, mock, "repos", "update", "--org", "org1")
	if err == nil {
		t.Fatal("expected error for missing --repo flag")
	}
}

func TestReposUpdateCmd_OnlyChangedFields(t *testing.T) {
	mock := &mockClient{
		postFunc: func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
			var payload map[string]interface{}
			if err := json.NewDecoder(body).Decode(&payload); err != nil {
				t.Fatalf("failed to decode body: %v", err)
			}
			if _, ok := payload["name"]; ok {
				t.Error("name should not be in body when not passed as flag")
			}
			if payload["visibility"] != "private" {
				t.Errorf("visibility = %q, want private", payload["visibility"])
			}
			return []byte(`{}`), nil
		},
	}

	_, err := executeCommand(t, mock, "repos", "update", "--org", "org1", "--repo", "repo1",
		"--visibility", "private")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
