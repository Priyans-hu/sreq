package env

import (
	"context"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Prefix: "TEST_",
		Paths:  map[string]string{"api_key": "{SERVICE}_API_KEY"},
	}

	provider, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if provider.prefix != "TEST_" {
		t.Errorf("prefix = %q, want %q", provider.prefix, "TEST_")
	}

	if provider.Name() != "env" {
		t.Errorf("Name() = %q, want %q", provider.Name(), "env")
	}
}

func TestProvider_Get(t *testing.T) {
	// Set up test environment variables
	_ = os.Setenv("TEST_API_KEY", "secret123")
	_ = os.Setenv("MY_SERVICE_TOKEN", "token456")
	defer func() { _ = os.Unsetenv("TEST_API_KEY") }()
	defer func() { _ = os.Unsetenv("MY_SERVICE_TOKEN") }()

	tests := []struct {
		name      string
		prefix    string
		key       string
		expected  string
		expectErr bool
	}{
		{
			name:     "simple key without prefix",
			prefix:   "",
			key:      "MY_SERVICE_TOKEN",
			expected: "token456",
		},
		{
			name:     "key with prefix",
			prefix:   "TEST_",
			key:      "API_KEY",
			expected: "secret123",
		},
		{
			name:     "lowercase key gets uppercased",
			prefix:   "TEST_",
			key:      "api_key",
			expected: "secret123",
		},
		{
			name:      "missing key returns error",
			prefix:    "",
			key:       "NONEXISTENT_KEY",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, _ := New(Config{Prefix: tt.prefix})
			ctx := context.Background()

			result, err := provider.Get(ctx, tt.key)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("Get() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestProvider_GetMultiple(t *testing.T) {
	_ = os.Setenv("KEY_ONE", "value1")
	_ = os.Setenv("KEY_TWO", "value2")
	defer func() { _ = os.Unsetenv("KEY_ONE") }()
	defer func() { _ = os.Unsetenv("KEY_TWO") }()

	provider, _ := New(Config{})
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

func TestProvider_GetMultiple_MissingKey(t *testing.T) {
	_ = os.Setenv("KEY_ONE", "value1")
	defer func() { _ = os.Unsetenv("KEY_ONE") }()

	provider, _ := New(Config{})
	ctx := context.Background()

	_, err := provider.GetMultiple(ctx, []string{"KEY_ONE", "MISSING_KEY"})
	if err == nil {
		t.Errorf("expected error for missing key, got nil")
	}
}

func TestProvider_Health(t *testing.T) {
	provider, _ := New(Config{})
	ctx := context.Background()

	err := provider.Health(ctx)
	if err != nil {
		t.Errorf("Health() error = %v, expected nil", err)
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
			name:     "single placeholder",
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
			name:     "hyphen to underscore conversion",
			template: "{SERVICE}_API_KEY",
			vars:     map[string]string{"service": "auth-service"},
			expected: "AUTH_SERVICE_API_KEY",
		},
		{
			name:     "dot to underscore conversion",
			template: "{SERVICE}_API_KEY",
			vars:     map[string]string{"service": "auth.service"},
			expected: "AUTH_SERVICE_API_KEY",
		},
		{
			name:     "no placeholders",
			template: "STATIC_API_KEY",
			vars:     map[string]string{"service": "auth"},
			expected: "STATIC_API_KEY",
		},
		{
			name:     "empty vars",
			template: "{SERVICE}_KEY",
			vars:     map[string]string{},
			expected: "{SERVICE}_KEY",
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

func TestProvider_GetWithTemplate(t *testing.T) {
	_ = os.Setenv("AUTH_PROD_API_KEY", "secret-key")
	defer func() { _ = os.Unsetenv("AUTH_PROD_API_KEY") }()

	provider, _ := New(Config{})
	ctx := context.Background()

	result, err := provider.GetWithTemplate(ctx, "{SERVICE}_{ENV}_API_KEY", map[string]string{
		"service": "auth",
		"env":     "prod",
	})

	if err != nil {
		t.Fatalf("GetWithTemplate() error = %v", err)
	}

	if result != "secret-key" {
		t.Errorf("GetWithTemplate() = %q, want %q", result, "secret-key")
	}
}
