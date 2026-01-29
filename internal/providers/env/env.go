package env

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Priyans-hu/sreq/internal/providers"
)

// Provider implements the providers.Provider interface for environment variables
type Provider struct {
	paths  map[string]string
	prefix string
}

// Config holds environment variable provider configuration
type Config struct {
	Paths  map[string]string
	Prefix string // Optional prefix for all env vars (e.g., "SREQ_")
}

// New creates a new environment variable provider
func New(cfg Config) (*Provider, error) {
	return &Provider{
		paths:  cfg.Paths,
		prefix: cfg.Prefix,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "env"
}

// Get retrieves a value from environment variables
// The key can be:
// - A direct env var name: "API_KEY"
// - A template with placeholders: "{SERVICE}_API_KEY"
func (p *Provider) Get(ctx context.Context, key string) (string, error) {
	// Apply prefix if set
	envKey := key
	if p.prefix != "" && !strings.HasPrefix(key, p.prefix) {
		envKey = p.prefix + key
	}

	// Convert to uppercase (env vars are typically uppercase)
	envKey = strings.ToUpper(envKey)

	value := os.Getenv(envKey)
	if value == "" {
		return "", fmt.Errorf("environment variable '%s' not set", envKey)
	}

	return value, nil
}

// GetMultiple retrieves multiple values from environment variables
func (p *Provider) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)

	for _, key := range keys {
		value, err := p.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		results[key] = value
	}

	return results, nil
}

// GetWithTemplate retrieves a value using a path template
// Supports placeholders: {service}, {env}, {region}, {project}
func (p *Provider) GetWithTemplate(ctx context.Context, template string, vars map[string]string) (string, error) {
	key := ResolvePath(template, vars)
	return p.Get(ctx, key)
}

// Health checks if the provider is available (always returns nil for env)
func (p *Provider) Health(ctx context.Context) error {
	// Environment variables are always available
	return nil
}

// ResolvePath replaces placeholders in a path template
// Placeholders: {service}, {env}, {region}, {project} (case-insensitive)
// Also converts to uppercase and replaces hyphens/dots with underscores
func ResolvePath(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		// Convert value to env-friendly format (uppercase, underscores)
		envValue := strings.ToUpper(value)
		envValue = strings.ReplaceAll(envValue, "-", "_")
		envValue = strings.ReplaceAll(envValue, ".", "_")
		// Replace both lowercase and uppercase placeholders
		result = strings.ReplaceAll(result, "{"+key+"}", envValue)
		result = strings.ReplaceAll(result, "{"+strings.ToUpper(key)+"}", envValue)
	}
	return result
}

// Ensure Provider implements the interface
var _ providers.Provider = (*Provider)(nil)
