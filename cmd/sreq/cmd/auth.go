package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure authentication for providers",
	Long: `Interactively configure authentication for Consul and AWS providers.

This command helps you set up credentials for:
  - Consul: address and token
  - AWS Secrets Manager: region and profile

Examples:
  sreq auth              # Interactive setup for all providers
  sreq auth consul       # Configure Consul only
  sreq auth aws          # Configure AWS only`,
	RunE: runAuthInteractive,
}

var authConsulCmd = &cobra.Command{
	Use:   "consul",
	Short: "Configure Consul authentication",
	RunE:  runAuthConsul,
}

var authAWSCmd = &cobra.Command{
	Use:   "aws",
	Short: "Configure AWS authentication",
	RunE:  runAuthAWS,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authConsulCmd)
	authCmd.AddCommand(authAWSCmd)
}

func runAuthInteractive(cmd *cobra.Command, args []string) error {
	fmt.Println("sreq Authentication Setup")
	fmt.Println("==========================")
	fmt.Println()

	// Configure Consul
	fmt.Println("1. Consul Configuration")
	fmt.Println("------------------------")
	if err := configureConsul(); err != nil {
		return err
	}
	fmt.Println()

	// Configure AWS
	fmt.Println("2. AWS Configuration")
	fmt.Println("---------------------")
	if err := configureAWS(); err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("Authentication setup complete!")
	fmt.Println("Run 'sreq config test' to verify your configuration.")
	return nil
}

func runAuthConsul(cmd *cobra.Command, args []string) error {
	fmt.Println("Consul Authentication Setup")
	fmt.Println("============================")
	fmt.Println()
	return configureConsul()
}

func runAuthAWS(cmd *cobra.Command, args []string) error {
	fmt.Println("AWS Authentication Setup")
	fmt.Println("=========================")
	fmt.Println()
	return configureAWS()
}

func configureConsul() error {
	reader := bufio.NewReader(os.Stdin)

	// Get current config
	cfg, configPath, err := loadCurrentConfig()
	if err != nil {
		return err
	}

	// Initialize providers map if nil
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]providerConfig)
	}

	consulCfg := cfg.Providers["consul"]

	// Consul address
	defaultAddr := consulCfg.Address
	if defaultAddr == "" {
		defaultAddr = "localhost:8500"
	}
	fmt.Printf("Consul address [%s]: ", defaultAddr)
	addr, _ := reader.ReadString('\n')
	addr = strings.TrimSpace(addr)
	if addr == "" {
		addr = defaultAddr
	}
	consulCfg.Address = addr

	// Consul token
	fmt.Print("Consul token (leave empty to use CONSUL_HTTP_TOKEN env): ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)
	if token != "" {
		consulCfg.Token = token
	} else if consulCfg.Token == "" {
		consulCfg.Token = "${CONSUL_HTTP_TOKEN}"
	}

	// Datacenter (optional)
	defaultDC := consulCfg.Datacenter
	fmt.Printf("Datacenter (optional) [%s]: ", defaultDC)
	dc, _ := reader.ReadString('\n')
	dc = strings.TrimSpace(dc)
	if dc != "" {
		consulCfg.Datacenter = dc
	}

	// Set default paths if not configured
	if consulCfg.Paths == nil {
		consulCfg.Paths = map[string]string{
			"base_url": "{project}/{env}/{app}/{region}/config/{service}/base_url",
			"username": "{project}/{env}/{app}/{region}/config/{service}/username",
		}
	}

	cfg.Providers["consul"] = consulCfg

	// Save config
	if err := saveConfig(cfg, configPath); err != nil {
		return err
	}

	fmt.Println("Consul configuration saved.")
	return nil
}

func configureAWS() error {
	reader := bufio.NewReader(os.Stdin)

	// Get current config
	cfg, configPath, err := loadCurrentConfig()
	if err != nil {
		return err
	}

	// Initialize providers map if nil
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]providerConfig)
	}

	awsCfg := cfg.Providers["aws_secrets"]

	// AWS region
	defaultRegion := awsCfg.Region
	if defaultRegion == "" {
		defaultRegion = os.Getenv("AWS_REGION")
		if defaultRegion == "" {
			defaultRegion = "us-east-1"
		}
	}
	fmt.Printf("AWS region [%s]: ", defaultRegion)
	region, _ := reader.ReadString('\n')
	region = strings.TrimSpace(region)
	if region == "" {
		region = defaultRegion
	}
	awsCfg.Region = region

	// AWS profile
	defaultProfile := awsCfg.Profile
	if defaultProfile == "" {
		defaultProfile = "default"
	}
	fmt.Printf("AWS profile [%s]: ", defaultProfile)
	profile, _ := reader.ReadString('\n')
	profile = strings.TrimSpace(profile)
	if profile == "" {
		profile = defaultProfile
	}
	if profile != "default" {
		awsCfg.Profile = profile
	}

	// Set default paths if not configured
	if awsCfg.Paths == nil {
		awsCfg.Paths = map[string]string{
			"password": "{service}/{env}/credentials#password",
			"api_key":  "{service}/{env}/credentials#api_key",
		}
	}

	cfg.Providers["aws_secrets"] = awsCfg

	// Save config
	if err := saveConfig(cfg, configPath); err != nil {
		return err
	}

	fmt.Println("AWS configuration saved.")
	fmt.Println()
	fmt.Println("Note: AWS credentials are read from:")
	fmt.Println("  1. Environment: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY")
	fmt.Println("  2. Credentials file: ~/.aws/credentials")
	fmt.Println("  3. IAM role (EC2/ECS/Lambda)")
	return nil
}

// providerConfig for yaml marshaling
type providerConfig struct {
	Address    string            `yaml:"address,omitempty"`
	Token      string            `yaml:"token,omitempty"`
	Region     string            `yaml:"region,omitempty"`
	Profile    string            `yaml:"profile,omitempty"`
	Datacenter string            `yaml:"datacenter,omitempty"`
	Paths      map[string]string `yaml:"paths,omitempty"`
}

// configFile represents the config file structure
type configFile struct {
	Providers      map[string]providerConfig `yaml:"providers"`
	Environments   []string                  `yaml:"environments,omitempty"`
	DefaultEnv     string                    `yaml:"default_env,omitempty"`
	Services       map[string]interface{}    `yaml:"services,omitempty"`
	Contexts       map[string]interface{}    `yaml:"contexts,omitempty"`
	DefaultContext string                    `yaml:"default_context,omitempty"`
}

func loadCurrentConfig() (*configFile, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".sreq", "config.yaml")

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return empty config
		return &configFile{
			Providers:    make(map[string]providerConfig),
			Environments: []string{"dev", "staging", "prod"},
			DefaultEnv:   "dev",
		}, configPath, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read config: %w", err)
	}

	var cfg configFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, "", fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Providers == nil {
		cfg.Providers = make(map[string]providerConfig)
	}

	return &cfg, configPath, nil
}

func saveConfig(cfg *configFile, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
