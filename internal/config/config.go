package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvAPIToken = "SOCKET_API_TOKEN"
	EnvBaseURL  = "SOCKET_API_URL"
	DefaultBase = "https://api.socket.dev/v0"
)

func GetAPIToken() (string, error) {
	token := os.Getenv(EnvAPIToken)
	if token != "" {
		return token, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("set %s or create ~/.config/socket/token", EnvAPIToken)
	}

	data, err := os.ReadFile(filepath.Join(home, ".config", "socket", "token"))
	if err != nil {
		return "", fmt.Errorf("set %s or create ~/.config/socket/token", EnvAPIToken)
	}

	token = strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("token file is empty")
	}
	return token, nil
}

func GetBaseURL() string {
	if u := os.Getenv(EnvBaseURL); u != "" {
		return strings.TrimRight(u, "/")
	}
	return DefaultBase
}
