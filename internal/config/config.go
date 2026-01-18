package config

import (
	"fmt"
	"os"
	"path/filepath"

	sreerrors "github.com/Priyans-hu/sreq/internal/errors"
	"github.com/Priyans-hu/sreq/pkg/types"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigDir is the default configuration directory
	DefaultConfigDir = ".sreq"

	// DefaultConfigFile is the default configuration file name
	DefaultConfigFile = "config.yaml"

	// DefaultServicesFile is the default services file name
	DefaultServicesFile = "services.yaml"
)

// Load loads the configuration from the default or specified path
func Load() (*types.Config, error) {
	configPath := os.Getenv("SREQ_CONFIG")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, DefaultConfigDir, DefaultConfigFile)
	}

	return LoadFromFile(configPath)
}

// LoadFromFile loads configuration from a specific file
func LoadFromFile(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, sreerrors.ConfigNotFound(path)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, sreerrors.ConfigParseError(path, err)
	}

	// Initialize services map if nil
	if config.Services == nil {
		config.Services = make(map[string]types.ServiceConfig)
	}

	// Load additional services from services.yaml
	configDir := filepath.Dir(path)
	servicesPath := filepath.Join(configDir, DefaultServicesFile)
	if err := loadServicesInto(&config, servicesPath); err != nil {
		// Ignore if services file doesn't exist
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return &config, nil
}

// loadServicesInto loads services from services.yaml into the config
func loadServicesInto(config *types.Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var servicesFile struct {
		Services map[string]types.ServiceConfig `yaml:"services"`
	}

	if err := yaml.Unmarshal(data, &servicesFile); err != nil {
		return sreerrors.ConfigParseError(path, err)
	}

	// Merge services (services.yaml takes precedence for conflicts)
	for name, svc := range servicesFile.Services {
		svc.Name = name
		config.Services[name] = svc
	}

	return nil
}

// Init creates the default configuration files
func Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, DefaultConfigDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default config
	configPath := filepath.Join(configDir, DefaultConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	// Create default services file
	servicesPath := filepath.Join(configDir, DefaultServicesFile)
	if _, err := os.Stat(servicesPath); os.IsNotExist(err) {
		if err := os.WriteFile(servicesPath, []byte(defaultServices), 0644); err != nil {
			return fmt.Errorf("failed to create services file: %w", err)
		}
	}

	return nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, DefaultConfigDir), nil
}

var defaultConfig = `# sreq configuration
# Documentation: https://github.com/Priyans-hu/sreq

# Available placeholders in paths:
#   {service} - Service name (from -s flag or consul_key)
#   {env}     - Environment (from -e flag)
#   {region}  - Region (from -r flag)
#   {project} - Project name (from -p flag)
#   {app}     - App name (from -a flag)

providers:
  consul:
    address: localhost:8500
    # token: ${CONSUL_TOKEN}
    # datacenter: us-east-1
    paths:
      base_url: "{project}/{env}/{app}/{region}/config/{service}/base_url"
      username: "{project}/{env}/{app}/{region}/config/{service}/username"

  aws_secrets:
    region: us-east-1
    # profile: default
    paths:
      password: "{service}/{env}/credentials#password"
      api_key: "{service}/{env}/credentials#api_key"

environments:
  - dev
  - staging
  - prod

default_env: dev

# Contexts are presets for quick switching between environments
# Use with: sreq run GET /api -s myservice -c dev-us
contexts:
  # Example contexts - customize for your setup
  # dev-us:
  #   project: myproject
  #   env: dev
  #   region: us-east-1
  #   app: myapp
  #
  # prod-eu:
  #   project: myproject
  #   env: prod
  #   region: eu-west-1
  #   app: myapp

# default_context: dev-us
`

var defaultServices = `# sreq services configuration
# Documentation: https://github.com/Priyans-hu/sreq

services:
  # ===================
  # SIMPLE MODE
  # ===================
  # Use when your Consul/AWS structure follows a standard pattern.
  # sreq uses the path templates from config.yaml to resolve credentials.
  #
  # example-service:
  #   consul_key: example           # Used in: services/{consul_key}/config/*
  #   aws_prefix: example-service   # Used in: {aws_prefix}/{env}/credentials

  # ===================
  # ADVANCED MODE
  # ===================
  # Use when you need explicit control over credential paths.
  # Supports mixed providers and complex structures.
  #
  # invoice:
  #   paths:
  #     base_url: "billing_service/invoice_svc_url"           # Consul (default)
  #     username: "billing_service/invoice_svc_username"      # Consul
  #     password: "aws:billing/{env}/invoice#password"        # AWS with JSON key
  #
  # Path format: [provider:]path[#jsonkey]
  #   - provider: consul (default), aws, vault, env
  #   - path: the key path (supports {service}, {env} placeholders)
  #   - #jsonkey: extract a key from JSON value (for AWS secrets)
`
