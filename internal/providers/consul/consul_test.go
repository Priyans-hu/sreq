package consul

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/consul/api"
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

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		expectErr bool
	}{
		{
			name: "valid config with address",
			cfg: Config{
				Address: "localhost:8500",
			},
			expectErr: false,
		},
		{
			name: "valid config with env addresses only",
			cfg: Config{
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with both",
			cfg: Config{
				Address: "consul-default:8500",
				EnvAddresses: map[string]string{
					"prod": "consul-prod:8500",
				},
				Token:      "test-token",
				Datacenter: "dc1",
			},
			expectErr: false,
		},
		{
			name:      "empty config - no addresses",
			cfg:       Config{},
			expectErr: true,
		},
		{
			name: "config with empty strings",
			cfg: Config{
				Address:      "",
				EnvAddresses: map[string]string{},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.cfg)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if p == nil {
				t.Error("expected provider, got nil")
			}
		})
	}
}

func TestNew_TokenFromEnv(t *testing.T) {
	// Set up test environment variable
	testToken := "test-env-token"
	_ = os.Setenv("TEST_CONSUL_TOKEN", testToken)
	defer func() { _ = os.Unsetenv("TEST_CONSUL_TOKEN") }()

	cfg := Config{
		Address: "localhost:8500",
		Token:   "${TEST_CONSUL_TOKEN}",
	}

	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if p.token != testToken {
		t.Errorf("token = %q, want %q", p.token, testToken)
	}
}

func TestNew_TokenFromEnv_NotSet(t *testing.T) {
	// Ensure env var is not set
	_ = os.Unsetenv("NONEXISTENT_TOKEN")

	cfg := Config{
		Address: "localhost:8500",
		Token:   "${NONEXISTENT_TOKEN}",
	}

	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Token should be empty when env var is not set
	if p.token != "" {
		t.Errorf("token = %q, want empty string", p.token)
	}
}

func TestProvider_Name(t *testing.T) {
	p := &Provider{}
	if p.Name() != "consul" {
		t.Errorf("Name() = %q, want %q", p.Name(), "consul")
	}
}

func TestProvider_queryOptions(t *testing.T) {
	tests := []struct {
		name       string
		datacenter string
	}{
		{
			name:       "with datacenter",
			datacenter: "dc1",
		},
		{
			name:       "without datacenter",
			datacenter: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{datacenter: tt.datacenter}
			ctx := context.Background()
			opts := p.queryOptions(ctx)

			if opts == nil {
				t.Error("queryOptions() returned nil")
				return
			}

			if tt.datacenter != "" && opts.Datacenter != tt.datacenter {
				t.Errorf("Datacenter = %q, want %q", opts.Datacenter, tt.datacenter)
			}
		})
	}
}

func TestProvider_getClientForEnv_NoAddress(t *testing.T) {
	p := &Provider{
		defaultAddress: "",
		envAddresses:   nil,
		clients:        make(map[string]*api.Client),
	}

	_, err := p.getClientForEnv("prod")
	if err == nil {
		t.Error("expected error for missing address, got nil")
	}
}
