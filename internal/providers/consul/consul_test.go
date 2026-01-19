package consul

import (
	"testing"
)

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		expected string
	}{
		{
			name:     "single placeholder",
			template: "services/{service}/config",
			vars:     map[string]string{"service": "auth"},
			expected: "services/auth/config",
		},
		{
			name:     "multiple placeholders",
			template: "{project}/{env}/{service}/config",
			vars: map[string]string{
				"project": "myapp",
				"env":     "prod",
				"service": "auth",
			},
			expected: "myapp/prod/auth/config",
		},
		{
			name:     "all placeholders",
			template: "{project}/{env}/{app}/{region}/config/{service}/url",
			vars: map[string]string{
				"project": "contacto",
				"env":     "dev",
				"app":     "hodor",
				"region":  "us-east-1",
				"service": "auth",
			},
			expected: "contacto/dev/hodor/us-east-1/config/auth/url",
		},
		{
			name:     "no placeholders",
			template: "static/path/to/key",
			vars:     map[string]string{"service": "auth"},
			expected: "static/path/to/key",
		},
		{
			name:     "missing variable",
			template: "{project}/{env}/{service}/config",
			vars: map[string]string{
				"project": "myapp",
				"service": "auth",
			},
			expected: "myapp/{env}/auth/config",
		},
		{
			name:     "empty vars",
			template: "{service}/config",
			vars:     map[string]string{},
			expected: "{service}/config",
		},
		{
			name:     "repeated placeholder",
			template: "{env}/{service}/{env}/config",
			vars: map[string]string{
				"env":     "prod",
				"service": "auth",
			},
			expected: "prod/auth/prod/config",
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

func TestProvider_getAddressForEnv(t *testing.T) {
	tests := []struct {
		name           string
		defaultAddress string
		envAddresses   map[string]string
		env            string
		expected       string
	}{
		{
			name:           "default address only",
			defaultAddress: "consul.default:8500",
			envAddresses:   nil,
			env:            "prod",
			expected:       "consul.default:8500",
		},
		{
			name:           "env override exists",
			defaultAddress: "consul-nonprod:8500",
			envAddresses: map[string]string{
				"prod": "consul-prod:8500",
			},
			env:      "prod",
			expected: "consul-prod:8500",
		},
		{
			name:           "env override not found - use default",
			defaultAddress: "consul-nonprod:8500",
			envAddresses: map[string]string{
				"prod": "consul-prod:8500",
			},
			env:      "dev",
			expected: "consul-nonprod:8500",
		},
		{
			name:           "multiple env addresses",
			defaultAddress: "consul-default:8500",
			envAddresses: map[string]string{
				"prod":    "consul-prod:8500",
				"staging": "consul-staging:8500",
			},
			env:      "staging",
			expected: "consul-staging:8500",
		},
		{
			name:           "empty env uses default",
			defaultAddress: "consul.default:8500",
			envAddresses: map[string]string{
				"prod": "consul-prod:8500",
			},
			env:      "",
			expected: "consul.default:8500",
		},
		{
			name:           "nonprod default, prod separate",
			defaultAddress: "consul-nonprod.internal:8500",
			envAddresses: map[string]string{
				"prod": "consul-prod.internal:8500",
			},
			env:      "qa",
			expected: "consul-nonprod.internal:8500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				defaultAddress: tt.defaultAddress,
				envAddresses:   tt.envAddresses,
			}
			result := p.getAddressForEnv(tt.env)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestProvider_GetAddresses(t *testing.T) {
	tests := []struct {
		name           string
		defaultAddress string
		envAddresses   map[string]string
		expected       map[string]string
	}{
		{
			name:           "default only",
			defaultAddress: "consul.default:8500",
			envAddresses:   nil,
			expected: map[string]string{
				"default": "consul.default:8500",
			},
		},
		{
			name:           "default and env addresses",
			defaultAddress: "consul-nonprod:8500",
			envAddresses: map[string]string{
				"prod": "consul-prod:8500",
			},
			expected: map[string]string{
				"default": "consul-nonprod:8500",
				"prod":    "consul-prod:8500",
			},
		},
		{
			name:           "only env addresses",
			defaultAddress: "",
			envAddresses: map[string]string{
				"prod": "consul-prod:8500",
				"dev":  "consul-dev:8500",
			},
			expected: map[string]string{
				"prod": "consul-prod:8500",
				"dev":  "consul-dev:8500",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				defaultAddress: tt.defaultAddress,
				envAddresses:   tt.envAddresses,
			}
			result := p.GetAddresses()
			if len(result) != len(tt.expected) {
				t.Errorf("got %d addresses, want %d", len(result), len(tt.expected))
			}
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("for key %q, got %q, want %q", key, result[key], expectedVal)
				}
			}
		})
	}
}

func TestResolvePathSimple(t *testing.T) {
	tests := []struct {
		name     string
		template string
		service  string
		env      string
		expected string
	}{
		{
			name:     "basic resolution",
			template: "services/{service}/{env}/config",
			service:  "auth",
			env:      "prod",
			expected: "services/auth/prod/config",
		},
		{
			name:     "service only",
			template: "services/{service}/url",
			service:  "billing",
			env:      "dev",
			expected: "services/billing/url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePathSimple(tt.template, tt.service, tt.env)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
