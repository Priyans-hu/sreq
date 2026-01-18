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
