package resolver

import (
	"testing"
)

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
