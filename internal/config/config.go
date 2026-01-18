package config

import (
	"fmt"
	"os"
	"path/filepath"

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
			return nil, fmt.Errorf("config file not found: %s (run 'sreq init' to create)", path)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
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
		return fmt.Errorf("failed to parse services file: %w", err)
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

providers:
  consul:
    address: localhost:8500
    # token: ${CONSUL_TOKEN}
    paths:
      base_url: "services/{service}/config/base_url"
      username: "services/{service}/config/username"

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
