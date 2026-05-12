package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/scottbrown/socket-cli/internal/config"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    config.GetBaseURL(),
		token:      token,
	}
}

func (c *Client) do(method, path string, query url.Values, body io.Reader, contentType string) (*http.Response, error) {
	u := c.baseURL + path
	if query != nil {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(b))
	}

	return resp, nil
}

func (c *Client) Get(path string, query url.Values) ([]byte, error) {
	resp, err := c.do(http.MethodGet, path, query, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) Post(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
	resp, err := c.do(http.MethodPost, path, query, body, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) Delete(path string, query url.Values) ([]byte, error) {
	resp, err := c.do(http.MethodDelete, path, query, nil, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) Put(path string, query url.Values, body io.Reader, contentType string) ([]byte, error) {
	resp, err := c.do(http.MethodPut, path, query, body, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func PrintJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(string(data))
		return nil
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

