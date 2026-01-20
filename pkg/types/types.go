package types

// ServiceConfig represents a service's configuration
// Supports two modes:
//
// Simple mode (for straightforward structures):
//
//	services:
//	  auth-service:
//	    consul_key: auth
//	    aws_prefix: auth-service
//
// Advanced mode (for complex structures):
//
//	services:
//	  invoice:
//	    paths:
//	      base_url: "billing_service/invoice_svc_url"
//	      username: "billing_service/invoice_svc_username"
//	      password: "aws:billing/{env}/invoice#password"
type ServiceConfig struct {
	Name      string `yaml:"name,omitempty"`
	ConsulKey string `yaml:"consul_key,omitempty"` // Simple mode: Consul key prefix
	AWSPrefix string `yaml:"aws_prefix,omitempty"` // Simple mode: AWS secret prefix

	// Advanced mode: explicit path mappings
	// Keys: base_url, username, password, api_key, or custom
	// Values: path with optional provider prefix (consul:, aws:)
	//         and JSON key suffix (#key)
	// Examples:
	//   "billing_service/invoice_svc_url"           -> Consul (default)
	//   "consul:billing_service/invoice_svc_url"    -> Consul (explicit)
	//   "aws:billing/dev/creds#password"            -> AWS with JSON key
	Paths map[string]string `yaml:"paths,omitempty"`
}

// IsAdvancedMode returns true if the service uses explicit path mappings
func (s *ServiceConfig) IsAdvancedMode() bool {
	return len(s.Paths) > 0
}

// ProviderConfig represents a secret provider's configuration
type ProviderConfig struct {
	Type       string            `yaml:"type,omitempty"`
	Address    string            `yaml:"address,omitempty"`
	Token      string            `yaml:"token,omitempty"`
	Region     string            `yaml:"region,omitempty"`
	Profile    string            `yaml:"profile,omitempty"`
	Datacenter string            `yaml:"datacenter,omitempty"` // Consul datacenter

	// Environment-specific addresses (overrides Address for specific envs)
	// Example:
	//   address: consul-nonprod.internal:8500      # default
	//   env_addresses:
	//     prod: consul-prod.internal:8500          # override for prod
	EnvAddresses map[string]string `yaml:"env_addresses,omitempty"`

	// Default path templates for simple mode
	// Use {service} and {env} as placeholders
	Paths map[string]string `yaml:"paths,omitempty"`

	// Env provider: prefix for environment variables (e.g., "SREQ_")
	Prefix string `yaml:"prefix,omitempty"`

	// Dotenv provider: single file path (for backward compatibility)
	File string `yaml:"file,omitempty"`

	// Dotenv provider: list of .env files to load (later files override earlier)
	// Supports: ".env", ".env.local", ".env.{env}", etc.
	Files []string `yaml:"files,omitempty"`
}

// GetAddressForEnv returns the appropriate address for the given environment.
// It checks env_addresses first, then falls back to the default address.
func (p *ProviderConfig) GetAddressForEnv(env string) string {
	if p.EnvAddresses != nil {
		if addr, ok := p.EnvAddresses[env]; ok {
			return addr
		}
	}
	return p.Address
}

// Context represents a preset of variables for quick switching
type Context struct {
	Project string `yaml:"project,omitempty"`
	Env     string `yaml:"env,omitempty"`
	Region  string `yaml:"region,omitempty"`
	App     string `yaml:"app,omitempty"`
}

// Config represents the main configuration
type Config struct {
	Providers      map[string]ProviderConfig `yaml:"providers"`
	Environments   []string                  `yaml:"environments"`
	DefaultEnv     string                    `yaml:"default_env"`
	Services       map[string]ServiceConfig  `yaml:"services,omitempty"`
	Contexts       map[string]Context        `yaml:"contexts,omitempty"`
	DefaultContext string                    `yaml:"default_context,omitempty"`
}

// ResolvedCredentials contains the resolved credentials for a service
type ResolvedCredentials struct {
	BaseURL  string            // The base URL for the service
	Username string            // Basic auth username
	Password string            // Basic auth password
	APIKey   string            // API key (if used instead of basic auth)
	Headers  map[string]string // Additional headers to add to requests
	Custom   map[string]string // Custom resolved values from paths
}

// Request represents an HTTP request to be made
type Request struct {
	Method      string
	Path        string
	Service     string
	Environment string
	Body        string
	Headers     map[string]string
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Status     string
	Headers    map[string][]string
	Body       []byte
}

// PathMapping represents a parsed path configuration
type PathMapping struct {
	Provider string // consul, aws, env, vault
	Path     string // The actual path
	JSONKey  string // Optional JSON key (after #)
}
