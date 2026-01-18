package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHistory_AddAndGet(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "sreq-history-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	h, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Add entry
	entry := Entry{
		Service: "auth-service",
		Env:     "dev",
		Method:  "GET",
		Path:    "/api/v1/users",
		BaseURL: "https://auth.dev.example.com",
		Status:  200,
	}
	h.Add(entry)

	// Verify entry was added with ID 1
	got, err := h.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Service != "auth-service" {
		t.Errorf("got service %q, want %q", got.Service, "auth-service")
	}
	if got.ID != 1 {
		t.Errorf("got ID %d, want 1", got.ID)
	}
}

func TestHistory_List(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-history-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	h, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Add multiple entries
	h.Add(Entry{Service: "auth", Env: "dev", Method: "GET", Path: "/users"})
	h.Add(Entry{Service: "billing", Env: "prod", Method: "POST", Path: "/invoice"})
	h.Add(Entry{Service: "auth", Env: "prod", Method: "GET", Path: "/health"})

	// List all
	entries := h.List(ListOptions{})
	if len(entries) != 3 {
		t.Errorf("got %d entries, want 3", len(entries))
	}

	// Filter by service
	entries = h.List(ListOptions{Service: "auth"})
	if len(entries) != 2 {
		t.Errorf("got %d entries for auth, want 2", len(entries))
	}

	// Filter by env
	entries = h.List(ListOptions{Env: "prod"})
	if len(entries) != 2 {
		t.Errorf("got %d entries for prod, want 2", len(entries))
	}

	// Filter by method
	entries = h.List(ListOptions{Method: "POST"})
	if len(entries) != 1 {
		t.Errorf("got %d entries for POST, want 1", len(entries))
	}

	// Limit
	entries = h.List(ListOptions{Limit: 2})
	if len(entries) != 2 {
		t.Errorf("got %d entries with limit 2, want 2", len(entries))
	}
}

func TestHistory_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-history-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	h, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Add entries
	h.Add(Entry{Service: "auth", Env: "dev"})
	h.Add(Entry{Service: "billing", Env: "prod"})

	if h.Count() != 2 {
		t.Errorf("got %d entries, want 2", h.Count())
	}

	// Clear
	h.Clear()
	if h.Count() != 0 {
		t.Errorf("got %d entries after clear, want 0", h.Count())
	}
}

func TestHistory_ClearBefore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-history-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	h, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Add old entry
	oldEntry := Entry{
		Service:   "auth",
		Env:       "dev",
		Timestamp: time.Now().Add(-48 * time.Hour),
	}
	h.Add(oldEntry)

	// Add recent entry
	recentEntry := Entry{
		Service: "billing",
		Env:     "prod",
	}
	h.Add(recentEntry)

	// Clear entries older than 24h
	removed := h.ClearBefore(24 * time.Hour)
	if removed != 1 {
		t.Errorf("removed %d entries, want 1", removed)
	}
	if h.Count() != 1 {
		t.Errorf("got %d entries after clear, want 1", h.Count())
	}
}

func TestHistory_SaveAndLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-history-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create and save
	h1, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	h1.Add(Entry{Service: "auth", Env: "dev", Method: "GET", Path: "/users"})
	if err := h1.Save(); err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tmpDir, DefaultHistoryFile)); err != nil {
		t.Errorf("history file not created: %v", err)
	}

	// Load fresh
	h2, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if h2.Count() != 1 {
		t.Errorf("got %d entries after load, want 1", h2.Count())
	}

	entry, _ := h2.Get(1)
	if entry.Service != "auth" {
		t.Errorf("got service %q, want %q", entry.Service, "auth")
	}
}

func TestHistory_MaxEntries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sreq-history-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	h, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	h.MaxEntries = 3

	// Add more than max
	for i := 0; i < 5; i++ {
		h.Add(Entry{Service: "svc", Path: "/test"})
	}

	if h.Count() != 3 {
		t.Errorf("got %d entries, want 3 (max)", h.Count())
	}
}

func TestEntry_ToCurl(t *testing.T) {
	tests := []struct {
		name     string
		entry    Entry
		contains []string
	}{
		{
			name: "GET request",
			entry: Entry{
				Method:  "GET",
				Path:    "/api/v1/users",
				BaseURL: "https://api.example.com",
			},
			contains: []string{"curl", "'https://api.example.com/api/v1/users'"},
		},
		{
			name: "POST with body",
			entry: Entry{
				Method:  "POST",
				Path:    "/api/v1/users",
				BaseURL: "https://api.example.com",
				Request: &Request{
					Body: `{"name":"test"}`,
				},
			},
			contains: []string{"curl", "-X", "POST", "-d"},
		},
		{
			name: "with headers",
			entry: Entry{
				Method:  "GET",
				Path:    "/api",
				BaseURL: "https://api.example.com",
				Request: &Request{
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			contains: []string{"-H", "Content-Type: application/json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curl := tt.entry.ToCurl()
			for _, want := range tt.contains {
				if !contains(curl, want) {
					t.Errorf("curl output %q missing %q", curl, want)
				}
			}
		})
	}
}

func TestRedactHeaders(t *testing.T) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer secret-token",
		"X-Api-Key":     "my-api-key",
		"Accept":        "application/json",
	}

	redacted := redactHeaders(headers)

	// Should not be redacted
	if redacted["Content-Type"] != "application/json" {
		t.Error("Content-Type should not be redacted")
	}
	if redacted["Accept"] != "application/json" {
		t.Error("Accept should not be redacted")
	}

	// Should be redacted
	if redacted["Authorization"] != "***REDACTED***" {
		t.Errorf("Authorization should be redacted, got %q", redacted["Authorization"])
	}
	if redacted["X-Api-Key"] != "***REDACTED***" {
		t.Errorf("X-Api-Key should be redacted, got %q", redacted["X-Api-Key"])
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
