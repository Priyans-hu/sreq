package dotenv

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		cfg           Config
		expectedFiles int
	}{
		{
			name:          "default to .env",
			cfg:           Config{},
			expectedFiles: 1,
		},
		{
			name: "single file",
			cfg: Config{
				File: ".env.local",
			},
			expectedFiles: 1,
		},
		{
			name: "multiple files",
			cfg: Config{
				Files: []string{".env", ".env.local"},
			},
			expectedFiles: 2,
		},
		{
			name: "file and files combined",
			cfg: Config{
				File:  ".env",
				Files: []string{".env.local"},
			},
			expectedFiles: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.cfg)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if len(provider.files) != tt.expectedFiles {
				t.Errorf("files count = %d, want %d", len(provider.files), tt.expectedFiles)
			}

			if provider.Name() != "dotenv" {
				t.Errorf("Name() = %q, want %q", provider.Name(), "dotenv")
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedKey   string
		expectedValue string
		expectErr     bool
	}{
		{
			name:          "simple key=value",
			line:          "API_KEY=secret123",
			expectedKey:   "API_KEY",
			expectedValue: "secret123",
		},
		{
			name:          "double quoted value",
			line:          `DATABASE_URL="postgres://localhost:5432/db"`,
			expectedKey:   "DATABASE_URL",
			expectedValue: "postgres://localhost:5432/db",
		},
		{
			name:          "single quoted value",
			line:          `SECRET='my secret value'`,
			expectedKey:   "SECRET",
			expectedValue: "my secret value",
		},
		{
			name:          "export prefix",
			line:          "export API_KEY=secret123",
			expectedKey:   "API_KEY",
			expectedValue: "secret123",
		},
		{
			name:          "spaces around equals",
			line:          "API_KEY = secret123",
			expectedKey:   "API_KEY",
			expectedValue: "secret123",
		},
		{
			name:          "escape sequence newline",
			line:          "MULTILINE=\"line1\\nline2\"",
			expectedKey:   "MULTILINE",
			expectedValue: "line1\nline2",
		},
		{
			name:          "escape sequence tab",
			line:          "TABBED=\"col1\\tcol2\"",
			expectedKey:   "TABBED",
			expectedValue: "col1\tcol2",
		},
		{
			name:      "no equals sign",
			line:      "INVALID_LINE",
			expectErr: true,
		},
		{
			name:      "empty key",
			line:      "=value",
			expectErr: true,
		},
		{
			name:          "empty value",
			line:          "EMPTY_KEY=",
			expectedKey:   "EMPTY_KEY",
			expectedValue: "",
		},
		{
			name:          "value with equals sign",
			line:          "CONNECTION=host=localhost;port=5432",
			expectedKey:   "CONNECTION",
			expectedValue: "host=localhost;port=5432",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, value, err := parseLine(tt.line)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseLine() error = %v", err)
			}

			if key != tt.expectedKey {
				t.Errorf("key = %q, want %q", key, tt.expectedKey)
			}
			if value != tt.expectedValue {
				t.Errorf("value = %q, want %q", value, tt.expectedValue)
			}
		})
	}
}

func TestProvider_LoadAndGet(t *testing.T) {
	// Create a temporary directory and .env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	content := `# Comment line
API_KEY=test-api-key
DATABASE_URL="postgres://localhost/db"
export SECRET_TOKEN=exported-secret
EMPTY_VALUE=
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test .env file: %v", err)
	}

	provider, err := New(Config{
		File: envFile,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		key       string
		expected  string
		expectErr bool
	}{
		{"API_KEY", "test-api-key", false},
		{"api_key", "test-api-key", false}, // lowercase should work
		{"DATABASE_URL", "postgres://localhost/db", false},
		{"SECRET_TOKEN", "exported-secret", false},
		{"EMPTY_VALUE", "", false},
		{"NONEXISTENT", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result, err := provider.Get(ctx, tt.key)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for key %q, got nil", tt.key)
				}
				return
			}

			if err != nil {
				t.Fatalf("Get(%q) error = %v", tt.key, err)
			}

			if result != tt.expected {
				t.Errorf("Get(%q) = %q, want %q", tt.key, result, tt.expected)
			}
		})
	}
}

func TestProvider_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// First .env file
	envFile1 := filepath.Join(tmpDir, ".env")
	content1 := `API_KEY=original
DATABASE_URL=postgres://localhost/db1
`
	if err := os.WriteFile(envFile1, []byte(content1), 0644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	// Second .env file (should override)
	envFile2 := filepath.Join(tmpDir, ".env.local")
	content2 := `API_KEY=overridden
NEW_KEY=new_value
`
	if err := os.WriteFile(envFile2, []byte(content2), 0644); err != nil {
		t.Fatalf("failed to write .env.local file: %v", err)
	}

	provider, err := New(Config{
		Files: []string{envFile1, envFile2},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()

	// API_KEY should be overridden by .env.local
	apiKey, err := provider.Get(ctx, "API_KEY")
	if err != nil {
		t.Fatalf("Get(API_KEY) error = %v", err)
	}
	if apiKey != "overridden" {
		t.Errorf("API_KEY = %q, want %q", apiKey, "overridden")
	}

	// DATABASE_URL should be from first file
	dbUrl, err := provider.Get(ctx, "DATABASE_URL")
	if err != nil {
		t.Fatalf("Get(DATABASE_URL) error = %v", err)
	}
	if dbUrl != "postgres://localhost/db1" {
		t.Errorf("DATABASE_URL = %q, want %q", dbUrl, "postgres://localhost/db1")
	}

	// NEW_KEY should be from second file
	newKey, err := provider.Get(ctx, "NEW_KEY")
	if err != nil {
		t.Fatalf("Get(NEW_KEY) error = %v", err)
	}
	if newKey != "new_value" {
		t.Errorf("NEW_KEY = %q, want %q", newKey, "new_value")
	}
}

func TestProvider_GetMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	content := `KEY_ONE=value1
KEY_TWO=value2
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test .env file: %v", err)
	}

	provider, _ := New(Config{File: envFile})
	ctx := context.Background()

	results, err := provider.GetMultiple(ctx, []string{"KEY_ONE", "KEY_TWO"})
	if err != nil {
		t.Fatalf("GetMultiple() error = %v", err)
	}

	if results["KEY_ONE"] != "value1" {
		t.Errorf("KEY_ONE = %q, want %q", results["KEY_ONE"], "value1")
	}
	if results["KEY_TWO"] != "value2" {
		t.Errorf("KEY_TWO = %q, want %q", results["KEY_TWO"], "value2")
	}
}

func TestProvider_GetAll(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	content := `KEY_ONE=value1
KEY_TWO=value2
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test .env file: %v", err)
	}

	provider, _ := New(Config{File: envFile})

	results, err := provider.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("GetAll() returned %d items, want 2", len(results))
	}
	if results["KEY_ONE"] != "value1" {
		t.Errorf("KEY_ONE = %q, want %q", results["KEY_ONE"], "value1")
	}
}

func TestProvider_Health(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	// File exists
	if err := os.WriteFile(envFile, []byte("KEY=value"), 0644); err != nil {
		t.Fatalf("failed to write test .env file: %v", err)
	}

	provider, _ := New(Config{File: envFile})
	ctx := context.Background()

	err := provider.Health(ctx)
	if err != nil {
		t.Errorf("Health() error = %v, expected nil", err)
	}

	// File doesn't exist
	provider2, _ := New(Config{File: "/nonexistent/.env"})
	err = provider2.Health(ctx)
	if err == nil {
		t.Errorf("Health() expected error for missing file, got nil")
	}
}

func TestProvider_Reload(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	// Initial content
	if err := os.WriteFile(envFile, []byte("API_KEY=original"), 0644); err != nil {
		t.Fatalf("failed to write test .env file: %v", err)
	}

	provider, _ := New(Config{File: envFile})
	ctx := context.Background()

	// Load initial
	val, _ := provider.Get(ctx, "API_KEY")
	if val != "original" {
		t.Errorf("initial API_KEY = %q, want %q", val, "original")
	}

	// Update file
	if err := os.WriteFile(envFile, []byte("API_KEY=updated"), 0644); err != nil {
		t.Fatalf("failed to update test .env file: %v", err)
	}

	// Should still return cached value
	val, _ = provider.Get(ctx, "API_KEY")
	if val != "original" {
		t.Errorf("cached API_KEY = %q, want %q", val, "original")
	}

	// Reload
	if err := provider.Reload(); err != nil {
		t.Fatalf("Reload() error = %v", err)
	}

	// Should now return updated value
	val, _ = provider.Get(ctx, "API_KEY")
	if val != "updated" {
		t.Errorf("reloaded API_KEY = %q, want %q", val, "updated")
	}
}

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		expected string
	}{
		{
			name:     "single placeholder uppercase",
			template: "{SERVICE}_API_KEY",
			vars:     map[string]string{"service": "auth"},
			expected: "AUTH_API_KEY",
		},
		{
			name:     "multiple placeholders",
			template: "{PROJECT}_{ENV}_{SERVICE}_KEY",
			vars: map[string]string{
				"project": "myapp",
				"env":     "prod",
				"service": "auth",
			},
			expected: "MYAPP_PROD_AUTH_KEY",
		},
		{
			name:     "hyphen to underscore",
			template: "{SERVICE}_KEY",
			vars:     map[string]string{"service": "auth-service"},
			expected: "AUTH_SERVICE_KEY",
		},
		{
			name:     "dot to underscore",
			template: "{SERVICE}_KEY",
			vars:     map[string]string{"service": "auth.service"},
			expected: "AUTH_SERVICE_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePath(tt.template, tt.vars)
			if result != tt.expected {
				t.Errorf("ResolvePath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestProvider_MissingFileSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, ".env")

	// Only create one file
	if err := os.WriteFile(existingFile, []byte("KEY=value"), 0644); err != nil {
		t.Fatalf("failed to write test .env file: %v", err)
	}

	// Include both existing and non-existing files
	provider, err := New(Config{
		Files: []string{existingFile, filepath.Join(tmpDir, ".env.nonexistent")},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()

	// Should still be able to load from existing file
	val, err := provider.Get(ctx, "KEY")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if val != "value" {
		t.Errorf("KEY = %q, want %q", val, "value")
	}
}
