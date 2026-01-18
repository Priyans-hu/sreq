package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	// DefaultHistoryFile is the default history file name
	DefaultHistoryFile = "history.json"

	// DefaultMaxEntries is the default maximum number of history entries
	DefaultMaxEntries = 100
)

// Entry represents a single request history entry
type Entry struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Env       string    `json:"env"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	BaseURL   string    `json:"base_url,omitempty"`
	Status    int       `json:"status,omitempty"`
	Duration  int64     `json:"duration_ms,omitempty"`
	Request   *Request  `json:"request,omitempty"`
	Response  *Response `json:"response,omitempty"`
}

// Request contains request details (with sensitive data redacted)
type Request struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// Response contains response details
type Response struct {
	Status    string `json:"status,omitempty"`
	SizeBytes int    `json:"size_bytes,omitempty"`
}

// History manages request history
type History struct {
	Entries    []Entry `json:"entries"`
	MaxEntries int     `json:"max_entries"`
	path       string
}

// sensitiveHeaders are headers that should be redacted
var sensitiveHeaders = []string{
	"authorization",
	"x-api-key",
	"x-auth-token",
	"cookie",
	"set-cookie",
}

// New creates a new history manager
func New(configDir string) (*History, error) {
	path := filepath.Join(configDir, DefaultHistoryFile)

	h := &History{
		Entries:    []Entry{},
		MaxEntries: DefaultMaxEntries,
		path:       path,
	}

	// Load existing history if file exists
	if _, err := os.Stat(path); err == nil {
		if err := h.load(); err != nil {
			return nil, fmt.Errorf("failed to load history: %w", err)
		}
	}

	return h, nil
}

// load reads history from disk
func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, h)
}

// Save writes history to disk
func (h *History) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(h.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.path, data, 0644)
}

// Add adds a new entry to the history
func (h *History) Add(entry Entry) {
	// Assign ID based on highest existing ID + 1
	maxID := 0
	for _, e := range h.Entries {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	entry.ID = maxID + 1

	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Redact sensitive headers
	if entry.Request != nil && entry.Request.Headers != nil {
		entry.Request.Headers = redactHeaders(entry.Request.Headers)
	}

	// Prepend to keep most recent first
	h.Entries = append([]Entry{entry}, h.Entries...)

	// Trim to max entries
	if len(h.Entries) > h.MaxEntries {
		h.Entries = h.Entries[:h.MaxEntries]
	}
}

// Get returns an entry by ID
func (h *History) Get(id int) (*Entry, error) {
	for i := range h.Entries {
		if h.Entries[i].ID == id {
			return &h.Entries[i], nil
		}
	}
	return nil, fmt.Errorf("history entry #%d not found", id)
}

// List returns entries with optional filtering
func (h *History) List(opts ListOptions) []Entry {
	var result []Entry

	for _, e := range h.Entries {
		// Apply filters
		if opts.Service != "" && e.Service != opts.Service {
			continue
		}
		if opts.Env != "" && e.Env != opts.Env {
			continue
		}
		if opts.Method != "" && !strings.EqualFold(e.Method, opts.Method) {
			continue
		}

		result = append(result, e)
	}

	// Apply limit
	if opts.Limit > 0 && len(result) > opts.Limit {
		result = result[:opts.Limit]
	}

	return result
}

// ListOptions contains filtering options for listing history
type ListOptions struct {
	Service string
	Env     string
	Method  string
	Limit   int
}

// Clear removes all history entries
func (h *History) Clear() {
	h.Entries = []Entry{}
}

// ClearBefore removes entries older than the given duration
func (h *History) ClearBefore(d time.Duration) int {
	cutoff := time.Now().Add(-d)
	var kept []Entry
	removed := 0

	for _, e := range h.Entries {
		if e.Timestamp.After(cutoff) {
			kept = append(kept, e)
		} else {
			removed++
		}
	}

	h.Entries = kept
	return removed
}

// Count returns the number of entries
func (h *History) Count() int {
	return len(h.Entries)
}

// redactHeaders redacts sensitive header values
func redactHeaders(headers map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range headers {
		lowerKey := strings.ToLower(k)
		redact := false
		for _, sensitive := range sensitiveHeaders {
			if lowerKey == sensitive {
				redact = true
				break
			}
		}
		if redact {
			result[k] = "***REDACTED***"
		} else {
			result[k] = v
		}
	}
	return result
}

// ToCurl converts an entry to a curl command
func (e *Entry) ToCurl() string {
	var parts []string
	parts = append(parts, "curl")

	// Method
	if e.Method != "GET" {
		parts = append(parts, "-X", e.Method)
	}

	// Headers
	if e.Request != nil && e.Request.Headers != nil {
		// Sort headers for consistent output
		var keys []string
		for k := range e.Request.Headers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := e.Request.Headers[k]
			parts = append(parts, "-H", fmt.Sprintf("'%s: %s'", k, v))
		}
	}

	// Body
	if e.Request != nil && e.Request.Body != "" {
		// Escape single quotes in body
		body := strings.ReplaceAll(e.Request.Body, "'", "'\\''")
		parts = append(parts, "-d", fmt.Sprintf("'%s'", body))
	}

	// URL
	url := e.BaseURL + e.Path
	parts = append(parts, fmt.Sprintf("'%s'", url))

	return strings.Join(parts, " ")
}

// ToHTTPie converts an entry to an HTTPie command
func (e *Entry) ToHTTPie() string {
	var parts []string
	parts = append(parts, "http")

	// Method and URL
	url := e.BaseURL + e.Path
	parts = append(parts, e.Method, url)

	// Headers
	if e.Request != nil && e.Request.Headers != nil {
		var keys []string
		for k := range e.Request.Headers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := e.Request.Headers[k]
			parts = append(parts, fmt.Sprintf("%s:'%s'", k, v))
		}
	}

	return strings.Join(parts, " ")
}

// FormatDuration formats the duration nicely
func (e *Entry) FormatDuration() string {
	if e.Duration == 0 {
		return "-"
	}
	if e.Duration < 1000 {
		return fmt.Sprintf("%dms", e.Duration)
	}
	return fmt.Sprintf("%.2fs", float64(e.Duration)/1000)
}

// StatusColor returns ANSI color code for the status
func (e *Entry) StatusColor() string {
	switch {
	case e.Status >= 200 && e.Status < 300:
		return "\033[32m" // Green
	case e.Status >= 300 && e.Status < 400:
		return "\033[33m" // Yellow
	case e.Status >= 400:
		return "\033[31m" // Red
	default:
		return ""
	}
}

// ResetColor returns ANSI reset code
func ResetColor() string {
	return "\033[0m"
}
