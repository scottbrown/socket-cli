package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClient_Get_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/test-path" {
			t.Errorf("path = %s, want /test-path", r.URL.Path)
		}
		if r.URL.Query().Get("key") != "val" {
			t.Errorf("query key = %q, want %q", r.URL.Query().Get("key"), "val")
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer my-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer my-token")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "my-token")
	q := url.Values{"key": {"val"}}
	data, err := c.Get("/test-path", q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("body = %q, want %q", string(data), `{"ok":true}`)
	}
}

func TestClient_Get_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token")
	_, err := c.Get("/missing", nil)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %q, want to contain '404'", err.Error())
	}
}

func TestClient_Post(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"test":1}` {
			t.Errorf("body = %q, want %q", string(body), `{"test":1}`)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"abc"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token")
	data, err := c.Post("/create", nil, strings.NewReader(`{"test":1}`), "application/json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"id":"abc"}` {
		t.Errorf("body = %q", string(data))
	}
}

func TestClient_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token")
	data, err := c.Delete("/resource/123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"deleted":true}` {
		t.Errorf("body = %q", string(data))
	}
}

func TestClient_Put(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated":true}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token")
	data, err := c.Put("/resource/123", nil, strings.NewReader("body"), "text/plain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"updated":true}` {
		t.Errorf("body = %q", string(data))
	}
}

func TestClient_NetworkError(t *testing.T) {
	c := NewClient("http://127.0.0.1:1", "token")
	_, err := c.Get("/anything", nil)
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestPrintJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := PrintJSON(&buf, []byte(`{"name":"test","value":42}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"name\": \"test\"") {
		t.Errorf("output not pretty-printed: %q", out)
	}
	if !strings.Contains(out, "\"value\": 42") {
		t.Errorf("output missing value: %q", out)
	}
}

func TestPrintJSON_InvalidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := PrintJSON(&buf, []byte("not json at all"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "not json at all") {
		t.Errorf("output should contain raw input: %q", out)
	}
}
