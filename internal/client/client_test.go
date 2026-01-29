package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Priyans-hu/sreq/pkg/types"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() returned nil")
	}
	if c.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestNew_WithTimeout(t *testing.T) {
	c := New(WithTimeout(60 * time.Second))
	if c.httpClient.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want %v", c.httpClient.Timeout, 60*time.Second)
	}
}

func TestNew_WithVerbose(t *testing.T) {
	c := New(WithVerbose(true))
	if !c.verbose {
		t.Error("verbose should be true")
	}
}

func TestClient_Do_GET(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/api/test" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/api/test")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/test",
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if !strings.Contains(string(resp.Body), "ok") {
		t.Errorf("Body = %q, should contain 'ok'", string(resp.Body))
	}
}

func TestClient_Do_POST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want %q", r.Method, http.MethodPost)
		}

		// Check content type was set
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want %q", r.Header.Get("Content-Type"), "application/json")
		}

		// Decode body
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}
		if body["name"] != "test" {
			t.Errorf("body.name = %q, want %q", body["name"], "test")
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id": "123"}`))
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req := &types.Request{
		Method: "POST",
		Path:   "/api/users",
		Body:   `{"name": "test"}`,
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
}

func TestClient_Do_WithBasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Basic auth not provided")
		}
		if username != "testuser" {
			t.Errorf("Username = %q, want %q", username, "testuser")
		}
		if password != "testpass" {
			t.Errorf("Password = %q, want %q", password, "testpass")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/secure",
	}

	creds := &types.ResolvedCredentials{
		BaseURL:  server.URL,
		Username: "testuser",
		Password: "testpass",
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Do_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request headers
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("X-Custom-Header = %q, want %q", r.Header.Get("X-Custom-Header"), "custom-value")
		}
		// Check credential headers
		if r.Header.Get("Authorization") != "Bearer token123" {
			t.Errorf("Authorization = %q, want %q", r.Header.Get("Authorization"), "Bearer token123")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/test",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
		Headers: map[string]string{
			"Authorization": "Bearer token123",
		},
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Do_ResponseHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Header", "response-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/test",
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if len(resp.Headers["X-Response-Header"]) == 0 || resp.Headers["X-Response-Header"][0] != "response-value" {
		t.Errorf("Response header = %v, want %q", resp.Headers["X-Response-Header"], "response-value")
	}
}

func TestClient_Do_Error(t *testing.T) {
	client := New(WithTimeout(100 * time.Millisecond))
	ctx := context.Background()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/test",
	}

	// Use invalid URL to trigger error
	creds := &types.ResolvedCredentials{
		BaseURL: "http://invalid-host-that-does-not-exist.local:12345",
	}

	_, err := client.Do(ctx, req, creds)
	if err == nil {
		t.Error("Do() expected error for invalid host")
	}
}

func TestClient_Do_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Long delay
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/slow",
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
	}

	_, err := client.Do(ctx, req, creds)
	if err == nil {
		t.Error("Do() expected error for cancelled context")
	}
}

func TestClient_Do_AllMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Errorf("Method = %q, want %q", r.Method, method)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := New()
			ctx := context.Background()

			req := &types.Request{
				Method: method,
				Path:   "/api/test",
			}

			creds := &types.ResolvedCredentials{
				BaseURL: server.URL,
			}

			resp, err := client.Do(ctx, req, creds)
			if err != nil {
				t.Fatalf("Do() error = %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
			}
		})
	}
}

func TestClient_Do_CustomContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// When custom content type is set, it should not be overwritten
		if r.Header.Get("Content-Type") != "application/xml" {
			t.Errorf("Content-Type = %q, want %q", r.Header.Get("Content-Type"), "application/xml")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New()
	ctx := context.Background()

	req := &types.Request{
		Method: "POST",
		Path:   "/api/test",
		Body:   "<xml>data</xml>",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Do_Verbose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Verbose mode should not cause errors
	client := New(WithVerbose(true))
	ctx := context.Background()

	req := &types.Request{
		Method: "GET",
		Path:   "/api/test",
		Headers: map[string]string{
			"X-Test": "value",
		},
	}

	creds := &types.ResolvedCredentials{
		BaseURL: server.URL,
	}

	resp, err := client.Do(ctx, req, creds)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Do_StatusCodes(t *testing.T) {
	codes := []int{200, 201, 204, 400, 401, 403, 404, 500, 502, 503}

	for _, code := range codes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			defer server.Close()

			client := New()
			ctx := context.Background()

			req := &types.Request{
				Method: "GET",
				Path:   "/api/test",
			}

			creds := &types.ResolvedCredentials{
				BaseURL: server.URL,
			}

			resp, err := client.Do(ctx, req, creds)
			if err != nil {
				t.Fatalf("Do() error = %v", err)
			}

			if resp.StatusCode != code {
				t.Errorf("StatusCode = %d, want %d", resp.StatusCode, code)
			}
		})
	}
}
