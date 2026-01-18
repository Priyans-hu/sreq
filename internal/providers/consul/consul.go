package consul

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Priyans-hu/sreq/internal/providers"
	"github.com/hashicorp/consul/api"
)

// Provider implements the providers.Provider interface for Consul KV
type Provider struct {
	client     *api.Client
	address    string
	token      string
	datacenter string
	paths      map[string]string
}

// Config holds Consul provider configuration
type Config struct {
	Address    string
	Token      string
	Datacenter string
	Paths      map[string]string
}

// New creates a new Consul provider
func New(cfg Config) (*Provider, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("consul address is required")
	}

	// Resolve token from environment variable if needed
	token := cfg.Token
	if strings.HasPrefix(token, "${") && strings.HasSuffix(token, "}") {
		envVar := token[2 : len(token)-1]
		token = os.Getenv(envVar)
	}

	// Create Consul client config
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Address

	if token != "" {
		consulConfig.Token = token
	}

	if cfg.Datacenter != "" {
		consulConfig.Datacenter = cfg.Datacenter
	}

	// Create the client
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &Provider{
		client:     client,
		address:    cfg.Address,
		token:      token,
		datacenter: cfg.Datacenter,
		paths:      cfg.Paths,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "consul"
}

// Get retrieves a value from Consul KV
func (p *Provider) Get(ctx context.Context, key string) (string, error) {
	kv := p.client.KV()

	// Get the value
	pair, _, err := kv.Get(key, p.queryOptions(ctx))
	if err != nil {
		return "", fmt.Errorf("failed to get key '%s' from consul: %w", key, err)
	}

	if pair == nil {
		return "", fmt.Errorf("key '%s' not found in consul", key)
	}

	return string(pair.Value), nil
}

// GetMultiple retrieves multiple values from Consul KV
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

// Health checks if Consul is reachable
func (p *Provider) Health(ctx context.Context) error {
	_, err := p.client.Status().Leader()
	if err != nil {
		return fmt.Errorf("consul health check failed: %w", err)
	}
	return nil
}

// ListKeys lists all keys under a prefix
func (p *Provider) ListKeys(ctx context.Context, prefix string) ([]string, error) {
	kv := p.client.KV()

	keys, _, err := kv.Keys(prefix, "", p.queryOptions(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list keys with prefix '%s': %w", prefix, err)
	}

	return keys, nil
}

// queryOptions creates query options with context
func (p *Provider) queryOptions(ctx context.Context) *api.QueryOptions {
	opts := &api.QueryOptions{}
	if p.datacenter != "" {
		opts.Datacenter = p.datacenter
	}
	return opts.WithContext(ctx)
}

// ResolvePath replaces placeholders in a path template
// Placeholders: {service}, {env}, {region}, {project}
func ResolvePath(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

// ResolvePathSimple is a convenience function for basic service/env resolution
func ResolvePathSimple(template, service, env string) string {
	return ResolvePath(template, map[string]string{
		"service": service,
		"env":     env,
	})
}

// Ensure Provider implements the interface
var _ providers.Provider = (*Provider)(nil)
