package aws

import (
	"testing"
)

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
			name:     "multiple keys - get first",
			json:     `{"username": "admin", "password": "secret123"}`,
			key:      "username",
			expected: "admin",
		},
		{
			name:     "multiple keys - get second",
			json:     `{"username": "admin", "password": "secret123"}`,
			key:      "password",
			expected: "secret123",
		},
		{
			name:     "numeric value",
			json:     `{"port": 5432}`,
			key:      "port",
			expected: "5432",
		},
		{
			name:     "boolean true",
			json:     `{"enabled": true}`,
			key:      "enabled",
			expected: "true",
		},
		{
			name:     "boolean false",
			json:     `{"debug": false}`,
			key:      "debug",
			expected: "false",
		},
		{
			name:     "no spaces",
			json:     `{"key":"value"}`,
			key:      "key",
			expected: "value",
		},
		{
			name:     "extra whitespace",
			json:     `{  "key"  :  "value"  }`,
			key:      "key",
			expected: "value",
		},
		{
			name:      "key not found",
			json:      `{"other": "value"}`,
			key:       "password",
			expectErr: true,
		},
		{
			name:      "empty object",
			json:      `{}`,
			key:       "key",
			expectErr: true,
		},
		{
			name:     "complex password",
			json:     `{"password": "P@$$w0rd!#%^&*()"}`,
			key:      "password",
			expected: "P@$$w0rd!#%^&*()",
		},
		{
			name:     "url value",
			json:     `{"url": "https://api.example.com/v1"}`,
			key:      "url",
			expected: "https://api.example.com/v1",
		},
		{
			name:     "value with colon",
			json:     `{"connection": "host:port"}`,
			key:      "connection",
			expected: "host:port",
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

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		expected string
	}{
		{
			name:     "service and env",
			template: "{service}/{env}/credentials",
			vars: map[string]string{
				"service": "auth-svc",
				"env":     "prod",
			},
			expected: "auth-svc/prod/credentials",
		},
		{
			name:     "all placeholders",
			template: "{project}/{app}/{env}/{service}",
			vars: map[string]string{
				"project": "myproject",
				"app":     "backend",
				"env":     "staging",
				"service": "api",
			},
			expected: "myproject/backend/staging/api",
		},
		{
			name:     "no placeholders",
			template: "static/secret/path",
			vars:     map[string]string{"service": "test"},
			expected: "static/secret/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePath(tt.template, tt.vars)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
