package types

import (
	"testing"
)

func TestProviderConfig_GetAddressForEnv(t *testing.T) {
	tests := []struct {
		name     string
		config   ProviderConfig
		env      string
		expected string
	}{
		{
			name: "default address only - returns default",
			config: ProviderConfig{
				Address: "consul.default:8500",
			},
			env:      "prod",
			expected: "consul.default:8500",
		},
		{
			name: "env_addresses only - returns matching env",
			config: ProviderConfig{
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
					"dev":  "consul-dev:8500",
				},
			},
			env:      "prod",
			expected: "consul-prod:8500",
		},
		{
			name: "env_addresses only - no match returns empty",
			config: ProviderConfig{
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
				},
			},
			env:      "staging",
			expected: "",
		},
		{
			name: "both - env match overrides default",
			config: ProviderConfig{
				Address: "consul-nonprod:8500",
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
				},
			},
			env:      "prod",
			expected: "consul-prod:8500",
		},
		{
			name: "both - no env match returns default",
			config: ProviderConfig{
				Address: "consul-nonprod:8500",
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
				},
			},
			env:      "dev",
			expected: "consul-nonprod:8500",
		},
		{
			name: "nonprod default, prod override",
			config: ProviderConfig{
				Address: "consul-nonprod.internal:8500",
				EnvAddresses: map[string]string{
					"prod": "consul-prod.internal:8500",
				},
			},
			env:      "staging",
			expected: "consul-nonprod.internal:8500",
		},
		{
			name: "empty env string uses default",
			config: ProviderConfig{
				Address: "consul.default:8500",
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
				},
			},
			env:      "",
			expected: "consul.default:8500",
		},
		{
			name:     "empty config returns empty",
			config:   ProviderConfig{},
			env:      "prod",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetAddressForEnv(tt.env)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

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
