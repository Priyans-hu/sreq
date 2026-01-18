package resolver

import (
	"context"
	"fmt"
	"strings"

	sreerrors "github.com/Priyans-hu/sreq/internal/errors"
	"github.com/Priyans-hu/sreq/internal/providers"
	"github.com/Priyans-hu/sreq/internal/providers/aws"
	"github.com/Priyans-hu/sreq/internal/providers/consul"
	"github.com/Priyans-hu/sreq/pkg/types"
)

// Resolver resolves credentials for services using configured providers
type Resolver struct {
	config    *types.Config
	providers map[string]providers.Provider
}

// New creates a new resolver with the given configuration
func New(cfg *types.Config) (*Resolver, error) {
	r := &Resolver{
		config:    cfg,
		providers: make(map[string]providers.Provider),
	}

	// Initialize providers based on config
	if err := r.initProviders(); err != nil {
		return nil, err
	}

	return r, nil
}

// initProviders initializes all configured providers
func (r *Resolver) initProviders() error {
	for name, providerCfg := range r.config.Providers {
		switch name {
		case "consul":
			provider, err := consul.New(consul.Config{
				Address:    providerCfg.Address,
				Token:      providerCfg.Token,
				Datacenter: providerCfg.Datacenter,
				Paths:      providerCfg.Paths,
			})
			if err != nil {
				return sreerrors.ProviderInitFailed("Consul", err)
			}
			r.providers["consul"] = provider

		case "aws_secrets", "aws":
			provider, err := aws.New(aws.Config{
				Region:  providerCfg.Region,
				Profile: providerCfg.Profile,
				Paths:   providerCfg.Paths,
			})
			if err != nil {
				return sreerrors.ProviderInitFailed("AWS Secrets Manager", err)
			}
			r.providers["aws"] = provider
			r.providers["aws_secrets"] = provider // alias

		default:
			// Unknown provider, skip
		}
	}

	return nil
}

// ResolveOptions contains options for credential resolution
type ResolveOptions struct {
	Service string
	Env     string
	Region  string
	Project string
	App     string
}

// Resolve resolves credentials for a service with given options
func (r *Resolver) Resolve(ctx context.Context, opts ResolveOptions) (*types.ResolvedCredentials, error) {
	// Find service config
	svcCfg, exists := r.config.Services[opts.Service]
	if !exists {
		return nil, sreerrors.ServiceNotFound(opts.Service)
	}

	creds := &types.ResolvedCredentials{
		Headers: make(map[string]string),
		Custom:  make(map[string]string),
	}

	// Build variables map for path resolution
	vars := map[string]string{
		"service": opts.Service,
		"env":     opts.Env,
	}
	if opts.Region != "" {
		vars["region"] = opts.Region
	}
	if opts.Project != "" {
		vars["project"] = opts.Project
	}
	if opts.App != "" {
		vars["app"] = opts.App
	}

	if svcCfg.IsAdvancedMode() {
		// Advanced mode: use explicit path mappings
		return r.resolveAdvanced(ctx, &svcCfg, vars, creds)
	}

	// Simple mode: use path templates from provider config
	return r.resolveSimple(ctx, &svcCfg, vars, creds)
}

// resolveSimple resolves credentials using simple mode (consul_key, aws_prefix)
func (r *Resolver) resolveSimple(ctx context.Context, svc *types.ServiceConfig, vars map[string]string, creds *types.ResolvedCredentials) (*types.ResolvedCredentials, error) {
	// Get Consul provider
	consulProvider, hasConsul := r.providers["consul"]

	if hasConsul && svc.ConsulKey != "" {
		// Get path templates from provider config
		consulCfg := r.config.Providers["consul"]

		// Add consul_key to vars for template resolution
		vars["service"] = svc.ConsulKey

		for key, template := range consulCfg.Paths {
			// Replace placeholders
			path := consul.ResolvePath(template, vars)

			value, err := consulProvider.Get(ctx, path)
			if err != nil {
				// Log warning but continue - not all keys may exist
				continue
			}

			// Map to credential fields
			switch key {
			case "base_url":
				creds.BaseURL = value
			case "username":
				creds.Username = value
			case "password":
				creds.Password = value
			case "api_key":
				creds.APIKey = value
			default:
				creds.Custom[key] = value
			}
		}
	}

	// Get AWS provider
	awsProvider, hasAWS := r.providers["aws"]

	if hasAWS && svc.AWSPrefix != "" {
		// Get path templates from provider config
		awsCfg := r.config.Providers["aws_secrets"]
		if awsCfg.Paths == nil {
			awsCfg = r.config.Providers["aws"]
		}

		// Add aws_prefix to vars for template resolution
		vars["service"] = svc.AWSPrefix

		for key, template := range awsCfg.Paths {
			// Replace placeholders
			path := aws.ResolvePath(template, vars)

			value, err := awsProvider.Get(ctx, path)
			if err != nil {
				// Log warning but continue - not all keys may exist
				continue
			}

			// Map to credential fields (AWS typically provides password/api_key)
			switch key {
			case "base_url":
				if creds.BaseURL == "" {
					creds.BaseURL = value
				}
			case "username":
				if creds.Username == "" {
					creds.Username = value
				}
			case "password":
				creds.Password = value
			case "api_key":
				creds.APIKey = value
			default:
				creds.Custom[key] = value
			}
		}
	}

	return creds, nil
}

// resolveAdvanced resolves credentials using advanced mode (explicit paths)
func (r *Resolver) resolveAdvanced(ctx context.Context, svc *types.ServiceConfig, vars map[string]string, creds *types.ResolvedCredentials) (*types.ResolvedCredentials, error) {
	for key, pathSpec := range svc.Paths {
		value, err := r.resolvePath(ctx, pathSpec, vars)
		if err != nil {
			return nil, sreerrors.PathResolutionFailed(key, err)
		}

		// Map to credential fields
		switch key {
		case "base_url":
			creds.BaseURL = value
		case "username":
			creds.Username = value
		case "password":
			creds.Password = value
		case "api_key":
			creds.APIKey = value
		default:
			creds.Custom[key] = value
		}
	}

	return creds, nil
}

// resolvePath resolves a single path specification
// Format: [provider:]path[#jsonkey]
func (r *Resolver) resolvePath(ctx context.Context, pathSpec string, vars map[string]string) (string, error) {
	parsed := parsePath(pathSpec)

	// Replace placeholders in path
	path := consul.ResolvePath(parsed.Path, vars)

	// Get provider (default to consul)
	providerName := parsed.Provider
	if providerName == "" {
		providerName = "consul"
	}

	provider, exists := r.providers[providerName]
	if !exists {
		return "", sreerrors.ProviderNotConfigured(providerName)
	}

	// Get value
	value, err := provider.Get(ctx, path)
	if err != nil {
		return "", err
	}

	// Extract JSON key if specified
	if parsed.JSONKey != "" {
		value, err = extractJSONKey(value, parsed.JSONKey)
		if err != nil {
			return "", sreerrors.JSONKeyNotFound(parsed.JSONKey, path)
		}
	}

	return value, nil
}

// PathSpec represents a parsed path specification
type PathSpec struct {
	Provider string // consul, aws, vault, env
	Path     string // The actual path
	JSONKey  string // Optional JSON key (after #)
}

// parsePath parses a path specification
// Format: [provider:]path[#jsonkey]
// Examples:
//   - "billing_service/invoice_url" -> consul:billing_service/invoice_url
//   - "consul:services/auth/url" -> consul:services/auth/url
//   - "aws:secrets/prod/db#password" -> aws:secrets/prod/db, key=password
func parsePath(spec string) PathSpec {
	result := PathSpec{}

	// Check for JSON key
	if idx := strings.LastIndex(spec, "#"); idx != -1 {
		result.JSONKey = spec[idx+1:]
		spec = spec[:idx]
	}

	// Check for provider prefix
	if idx := strings.Index(spec, ":"); idx != -1 {
		// Make sure it's not a Windows path (C:\...)
		if idx > 1 || (idx == 1 && len(spec) > 2 && spec[2] != '\\') {
			result.Provider = spec[:idx]
			spec = spec[idx+1:]
		}
	}

	result.Path = spec
	return result
}

// extractJSONKey extracts a value from a JSON string
func extractJSONKey(jsonStr, key string) (string, error) {
	// Simple JSON extraction - for complex cases, use encoding/json
	// This handles simple cases like {"password": "secret"}

	// Look for "key": "value" or "key":"value"
	searchKey := fmt.Sprintf(`"%s"`, key)
	idx := strings.Index(jsonStr, searchKey)
	if idx == -1 {
		return "", fmt.Errorf("key '%s' not found in JSON", key)
	}

	// Find the colon after the key
	rest := jsonStr[idx+len(searchKey):]
	colonIdx := strings.Index(rest, ":")
	if colonIdx == -1 {
		return "", fmt.Errorf("invalid JSON format")
	}

	// Find the value
	rest = strings.TrimSpace(rest[colonIdx+1:])

	if len(rest) == 0 {
		return "", fmt.Errorf("empty value for key '%s'", key)
	}

	// Handle string value
	if rest[0] == '"' {
		endIdx := strings.Index(rest[1:], `"`)
		if endIdx == -1 {
			return "", fmt.Errorf("unterminated string value")
		}
		return rest[1 : endIdx+1], nil
	}

	// Handle non-string value (number, bool, null)
	endIdx := strings.IndexAny(rest, ",}")
	if endIdx == -1 {
		endIdx = len(rest)
	}
	return strings.TrimSpace(rest[:endIdx]), nil
}

// GetProvider returns a provider by name
func (r *Resolver) GetProvider(name string) (providers.Provider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

// HealthCheck checks all providers
func (r *Resolver) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)
	for name, provider := range r.providers {
		results[name] = provider.Health(ctx)
	}
	return results
}
