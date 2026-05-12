package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/scottbrown/socket-cli/internal/api"
)

type mockClient struct {
	getFunc    func(path string, query url.Values) ([]byte, error)
	postFunc   func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error)
	deleteFunc func(path string, query url.Values) ([]byte, error)
	putFunc    func(path string, query url.Values, body io.Reader, contentType string) ([]byte, error)
}

func (m *mockClient) Get(path string, query url.Values) ([]byte, error) {
	if m.getFunc != nil {
		return m.getFunc(path, query)
	}
	return nil, fmt.Errorf("Get not mocked")
}

func (m *mockClient) Post(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
	if m.postFunc != nil {
		return m.postFunc(path, query, body, contentType)
	}
	return nil, fmt.Errorf("Post not mocked")
}

func (m *mockClient) Delete(path string, query url.Values) ([]byte, error) {
	if m.deleteFunc != nil {
		return m.deleteFunc(path, query)
	}
	return nil, fmt.Errorf("Delete not mocked")
}

func (m *mockClient) Put(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
	if m.putFunc != nil {
		return m.putFunc(path, query, body, contentType)
	}
	return nil, fmt.Errorf("Put not mocked")
}

var _ api.SocketAPI = (*mockClient)(nil)

func executeCommand(t *testing.T, client api.SocketAPI, args ...string) (string, error) {
	t.Helper()
	root := newRootCmd(client)
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}
