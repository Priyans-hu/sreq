package resolver

import (
	"context"
	"fmt"
	"testing"

	"github.com/Priyans-hu/sreq/internal/providers"
	"github.com/Priyans-hu/sreq/pkg/types"
)

// mockProvider implements providers.Provider for testing
type mockProvider struct {
	name   string
	values map[string]string
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Get(ctx context.Context, key string) (string, error) {
	if val, ok := m.values[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key '%s' not found", key)
}

func (m *mockProvider) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, key := range keys {
		val, err := m.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		results[key] = val
	}
	return results, nil
}

func (m *mockProvider) Health(ctx context.Context) error {
	return nil
}

var _ providers.Provider = (*mockProvider)(nil)

func TestNew(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services:  map[string]types.ServiceConfig{},
	}

	r, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if r == nil {
		t.Fatal("New() returned nil")
	}
}

func TestResolver_GetProvider(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services:  map[string]types.ServiceConfig{},
	}

	r, _ := New(cfg)

	// Add mock provider directly
	mock := &mockProvider{name: "test", values: map[string]string{"key": "value"}}
	r.providers["test"] = mock

	// Test GetProvider
	p, ok := r.GetProvider("test")
	if !ok {
		t.Error("GetProvider() should find 'test' provider")
	}
	if p.Name() != "test" {
		t.Errorf("Provider name = %q, want %q", p.Name(), "test")
	}

	// Test non-existent provider
	_, ok = r.GetProvider("nonexistent")
	if ok {
		t.Error("GetProvider() should not find 'nonexistent' provider")
	}
}

func TestResolver_HealthCheck(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services:  map[string]types.ServiceConfig{},
	}

	r, _ := New(cfg)
	r.providers["mock1"] = &mockProvider{name: "mock1", values: map[string]string{}}
	r.providers["mock2"] = &mockProvider{name: "mock2", values: map[string]string{}}

	ctx := context.Background()
	results := r.HealthCheck(ctx)

	if len(results) != 2 {
		t.Errorf("HealthCheck() returned %d results, want 2", len(results))
	}

	// Both mock providers should return nil (healthy)
	for name, err := range results {
		if err != nil {
			t.Errorf("HealthCheck() %s returned error: %v", name, err)
		}
	}
}

func TestResolver_Resolve_ServiceNotFound(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services:  map[string]types.ServiceConfig{},
	}

	r, _ := New(cfg)
	ctx := context.Background()

	_, err := r.Resolve(ctx, ResolveOptions{
		Service: "nonexistent",
		Env:     "dev",
	})

	if err == nil {
		t.Error("Resolve() should return error for nonexistent service")
	}
}

func TestResolver_Resolve_AdvancedMode(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services: map[string]types.ServiceConfig{
			"test-service": {
				Paths: map[string]string{
					"base_url": "mock:services/test/url",
					"username": "mock:services/test/username",
					"password": "mock:services/test/password",
					"api_key":  "mock:services/test/api_key",
					"custom":   "mock:services/test/custom",
				},
			},
		},
	}

	r, _ := New(cfg)
	r.providers["mock"] = &mockProvider{
		name: "mock",
		values: map[string]string{
			"services/test/url":      "https://api.test.com",
			"services/test/username": "admin",
			"services/test/password": "secret",
			"services/test/api_key":  "key123",
			"services/test/custom":   "custom-value",
		},
	}

	ctx := context.Background()
	creds, err := r.Resolve(ctx, ResolveOptions{
		Service: "test-service",
		Env:     "dev",
	})

	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if creds.BaseURL != "https://api.test.com" {
		t.Errorf("BaseURL = %q, want %q", creds.BaseURL, "https://api.test.com")
	}
	if creds.Username != "admin" {
		t.Errorf("Username = %q, want %q", creds.Username, "admin")
	}
	if creds.Password != "secret" {
		t.Errorf("Password = %q, want %q", creds.Password, "secret")
	}
	if creds.APIKey != "key123" {
		t.Errorf("APIKey = %q, want %q", creds.APIKey, "key123")
	}
	if creds.Custom["custom"] != "custom-value" {
		t.Errorf("Custom[custom] = %q, want %q", creds.Custom["custom"], "custom-value")
	}
}

func TestResolver_Resolve_ProviderNotConfigured(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services: map[string]types.ServiceConfig{
			"test-service": {
				Paths: map[string]string{
					"base_url": "nonexistent:services/test/url",
				},
			},
		},
	}

	r, _ := New(cfg)
	ctx := context.Background()

	_, err := r.Resolve(ctx, ResolveOptions{
		Service: "test-service",
		Env:     "dev",
	})

	if err == nil {
		t.Error("Resolve() should return error for unconfigured provider")
	}
}

func TestResolver_Resolve_WithPlaceholders(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services: map[string]types.ServiceConfig{
			"auth": {
				Paths: map[string]string{
					"base_url": "mock:{service}/{env}/url",
				},
			},
		},
	}

	r, _ := New(cfg)
	r.providers["mock"] = &mockProvider{
		name: "mock",
		values: map[string]string{
			"auth/prod/url": "https://auth.prod.example.com",
		},
	}

	ctx := context.Background()
	creds, err := r.Resolve(ctx, ResolveOptions{
		Service: "auth",
		Env:     "prod",
	})

	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if creds.BaseURL != "https://auth.prod.example.com" {
		t.Errorf("BaseURL = %q, want %q", creds.BaseURL, "https://auth.prod.example.com")
	}
}

func TestResolver_Resolve_WithJSONKey(t *testing.T) {
	cfg := &types.Config{
		Providers: map[string]types.ProviderConfig{},
		Services: map[string]types.ServiceConfig{
			"db-service": {
				Paths: map[string]string{
					"password": "mock:secrets/db#password",
					"username": "mock:secrets/db#username",
				},
			},
		},
	}

	r, _ := New(cfg)
	r.providers["mock"] = &mockProvider{
		name: "mock",
		values: map[string]string{
			"secrets/db": `{"username": "dbadmin", "password": "dbsecret123"}`,
		},
	}

	ctx := context.Background()
	creds, err := r.Resolve(ctx, ResolveOptions{
		Service: "db-service",
		Env:     "dev",
	})

	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if creds.Password != "dbsecret123" {
		t.Errorf("Password = %q, want %q", creds.Password, "dbsecret123")
	}
	if creds.Username != "dbadmin" {
		t.Errorf("Username = %q, want %q", creds.Username, "dbadmin")
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected PathSpec
	}{
		{
			name:  "simple path",
			input: "services/auth/url",
			expected: PathSpec{
				Provider: "",
				Path:     "services/auth/url",
				JSONKey:  "",
			},
		},
		{
			name:  "path with consul provider",
			input: "consul:services/auth/url",
			expected: PathSpec{
				Provider: "consul",
				Path:     "services/auth/url",
				JSONKey:  "",
			},
		},
		{
			name:  "path with aws provider",
			input: "aws:secrets/prod/db",
			expected: PathSpec{
				Provider: "aws",
				Path:     "secrets/prod/db",
				JSONKey:  "",
			},
		},
		{
			name:  "path with json key",
			input: "secrets/prod/db#password",
			expected: PathSpec{
				Provider: "",
				Path:     "secrets/prod/db",
				JSONKey:  "password",
			},
		},
		{
			name:  "full path with provider and json key",
			input: "aws:secrets/prod/db#password",
			expected: PathSpec{
				Provider: "aws",
				Path:     "secrets/prod/db",
				JSONKey:  "password",
			},
		},
		{
			name:  "path with multiple slashes",
			input: "consul:project/env/service/region/config/key",
			expected: PathSpec{
				Provider: "consul",
				Path:     "project/env/service/region/config/key",
				JSONKey:  "",
			},
		},
		{
			name:  "path with nested json key",
			input: "aws:myapp/prod/credentials#db.password",
			expected: PathSpec{
				Provider: "aws",
				Path:     "myapp/prod/credentials",
				JSONKey:  "db.password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePath(tt.input)

			if result.Provider != tt.expected.Provider {
				t.Errorf("Provider: got %q, want %q", result.Provider, tt.expected.Provider)
			}
			if result.Path != tt.expected.Path {
				t.Errorf("Path: got %q, want %q", result.Path, tt.expected.Path)
			}
			if result.JSONKey != tt.expected.JSONKey {
				t.Errorf("JSONKey: got %q, want %q", result.JSONKey, tt.expected.JSONKey)
			}
		})
	}
}

func TestExtractJSONKey(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		key       string
		expected  string
		expectErr bool
	}{
		{
			name:     "simple string value",
			json:     `{"password": "secret123"}`,
			key:      "password",
			expected: "secret123",
		},
		{
			name:     "multiple keys",
			json:     `{"username": "admin", "password": "secret123"}`,
			key:      "password",
			expected: "secret123",
		},
		{
			name:     "first key",
			json:     `{"username": "admin", "password": "secret123"}`,
			key:      "username",
			expected: "admin",
		},
		{
			name:     "numeric value",
			json:     `{"port": 5432, "host": "localhost"}`,
			key:      "port",
			expected: "5432",
		},
		{
			name:     "boolean value",
			json:     `{"enabled": true, "debug": false}`,
			key:      "enabled",
			expected: "true",
		},
		{
			name:     "null value",
			json:     `{"value": null}`,
			key:      "value",
			expected: "null",
		},
		{
			name:     "no spaces around colon",
			json:     `{"password":"secret123"}`,
			key:      "password",
			expected: "secret123",
		},
		{
			name:      "key not found",
			json:      `{"username": "admin"}`,
			key:       "password",
			expectErr: true,
		},
		{
			name:      "empty json",
			json:      `{}`,
			key:       "password",
			expectErr: true,
		},
		{
			name:     "value with special characters",
			json:     `{"password": "p@ss!word#123"}`,
			key:      "password",
			expected: "p@ss!word#123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractJSONKey(tt.json, tt.key)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
