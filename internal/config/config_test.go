package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAPIToken_FromEnv(t *testing.T) {
	t.Setenv(EnvAPIToken, "test-token-123")

	token, err := GetAPIToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "test-token-123" {
		t.Errorf("got %q, want %q", token, "test-token-123")
	}
}

func TestGetAPIToken_FromFile(t *testing.T) {
	t.Setenv(EnvAPIToken, "")

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	tokenDir := filepath.Join(tmpDir, ".config", "socket")
	if err := os.MkdirAll(tokenDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tokenDir, "token"), []byte("file-token-456\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	token, err := GetAPIToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "file-token-456" {
		t.Errorf("got %q, want %q", token, "file-token-456")
	}
}

func TestGetAPIToken_MissingBoth(t *testing.T) {
	t.Setenv(EnvAPIToken, "")

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	_, err := GetAPIToken()
	if err == nil {
		t.Fatal("expected error when no token available")
	}
}

func TestGetAPIToken_EmptyFile(t *testing.T) {
	t.Setenv(EnvAPIToken, "")

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	tokenDir := filepath.Join(tmpDir, ".config", "socket")
	if err := os.MkdirAll(tokenDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tokenDir, "token"), []byte("  \n"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := GetAPIToken()
	if err == nil {
		t.Fatal("expected error for empty token file")
	}
}

func TestGetBaseURL_FromEnv(t *testing.T) {
	t.Setenv(EnvBaseURL, "https://custom.api.dev/v1/")

	u := GetBaseURL()
	if u != "https://custom.api.dev/v1" {
		t.Errorf("got %q, want trailing slash trimmed", u)
	}
}

func TestGetBaseURL_Default(t *testing.T) {
	t.Setenv(EnvBaseURL, "")

	u := GetBaseURL()
	if u != DefaultBase {
		t.Errorf("got %q, want %q", u, DefaultBase)
	}
}
