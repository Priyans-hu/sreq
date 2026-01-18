package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Priyans-hu/sreq/internal/providers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Provider implements the providers.Provider interface for AWS Secrets Manager
type Provider struct {
	client  *secretsmanager.Client
	region  string
	profile string
	paths   map[string]string
}

// Config holds AWS Secrets Manager provider configuration
type Config struct {
	Region  string
	Profile string
	Paths   map[string]string
}

// New creates a new AWS Secrets Manager provider
func New(cfg Config) (*Provider, error) {
	region := cfg.Region
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			region = "us-east-1"
		}
	}

	// Build AWS config options
	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(region))

	// Use profile if specified
	if cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(cfg.Profile))
	}

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Secrets Manager client
	client := secretsmanager.NewFromConfig(awsCfg)

	return &Provider{
		client:  client,
		region:  region,
		profile: cfg.Profile,
		paths:   cfg.Paths,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "aws_secrets"
}

// Get retrieves a secret from AWS Secrets Manager
// The key format is: secret-name or secret-name#json-key
func (p *Provider) Get(ctx context.Context, key string) (string, error) {
	// Parse the key for JSON key extraction (secret-name#jsonkey)
	secretName := key
	var jsonKey string
	if idx := strings.LastIndex(key, "#"); idx != -1 {
		secretName = key[:idx]
		jsonKey = key[idx+1:]
	}

	// Get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := p.client.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret '%s': %w", secretName, err)
	}

	var secretValue string
	if result.SecretString != nil {
		secretValue = *result.SecretString
	} else {
		return "", fmt.Errorf("secret '%s' has no string value (binary secrets not supported)", secretName)
	}

	// If no JSON key specified, return the whole value
	if jsonKey == "" {
		return secretValue, nil
	}

	// Extract JSON key
	extracted, err := extractJSONKey(secretValue, jsonKey)
	if err != nil {
		return "", fmt.Errorf("failed to extract key '%s' from secret '%s': %w", jsonKey, secretName, err)
	}

	return extracted, nil
}

// GetMultiple retrieves multiple secrets from AWS Secrets Manager
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

// Health checks if AWS Secrets Manager is reachable
func (p *Provider) Health(ctx context.Context) error {
	// List secrets with max 1 result just to verify connectivity
	input := &secretsmanager.ListSecretsInput{
		MaxResults: aws.Int32(1),
	}

	_, err := p.client.ListSecrets(ctx, input)
	if err != nil {
		return fmt.Errorf("AWS Secrets Manager health check failed: %w", err)
	}
	return nil
}

// extractJSONKey extracts a value from a JSON string
// Handles simple cases like {"password": "secret", "username": "admin"}
func extractJSONKey(jsonStr, key string) (string, error) {
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

// ResolvePath replaces placeholders in a path template
// Placeholders: {service}, {env}, {region}, {project}, {app}
func ResolvePath(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

// Ensure Provider implements the interface
var _ providers.Provider = (*Provider)(nil)
