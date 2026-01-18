package types

import (
	"testing"
)

func TestServiceConfig_IsAdvancedMode(t *testing.T) {
	tests := []struct {
		name     string
		config   ServiceConfig
		expected bool
	}{
		{
			name: "simple mode - consul_key only",
			config: ServiceConfig{
				ConsulKey: "auth",
			},
			expected: false,
		},
		{
			name: "simple mode - aws_prefix only",
			config: ServiceConfig{
				AWSPrefix: "auth-svc",
			},
			expected: false,
		},
		{
			name: "simple mode - both consul and aws",
			config: ServiceConfig{
				ConsulKey: "auth",
				AWSPrefix: "auth-svc",
			},
			expected: false,
		},
		{
			name: "advanced mode - with paths",
			config: ServiceConfig{
				Paths: map[string]string{
					"base_url": "consul:services/auth/url",
					"password": "aws:secrets/auth#password",
				},
			},
			expected: true,
		},
		{
			name: "advanced mode - single path",
			config: ServiceConfig{
				Paths: map[string]string{
					"base_url": "services/auth/url",
				},
			},
			expected: true,
		},
		{
			name:     "empty config",
			config:   ServiceConfig{},
			expected: false,
		},
		{
			name: "mixed - has both consul_key and paths (paths wins)",
			config: ServiceConfig{
				ConsulKey: "auth",
				Paths: map[string]string{
					"base_url": "custom/path",
				},
			},
			expected: true,
		},
		{
			name: "empty paths map",
			config: ServiceConfig{
				Paths: map[string]string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsAdvancedMode()
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
