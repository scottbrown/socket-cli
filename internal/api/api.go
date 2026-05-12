package api

import (
	"io"
	"net/url"
)

type SocketAPI interface {
	Get(path string, query url.Values) ([]byte, error)
	Post(path string, query url.Values, body io.Reader, contentType string) ([]byte, error)
	Delete(path string, query url.Values) ([]byte, error)
	Put(path string, query url.Values, body io.Reader, contentType string) ([]byte, error)
}
