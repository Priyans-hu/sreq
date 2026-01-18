package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Long:  `Add, list, or remove service configurations.`,
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured services",
	RunE:  runServiceList,
}

var serviceAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new service",
	Long: `Add a new service configuration.

Example:
  sreq service add auth-service --consul-key auth --aws-prefix auth-svc`,
	Args: cobra.ExactArgs(1),
	RunE: runServiceAdd,
}

var serviceRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a service",
	Args:  cobra.ExactArgs(1),
	RunE:  runServiceRemove,
}

// Flags for service add
var (
	consulKey string
	awsPrefix string
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceAddCmd)
	serviceCmd.AddCommand(serviceRemoveCmd)

	serviceAddCmd.Flags().StringVar(&consulKey, "consul-key", "", "Consul key prefix for this service")
	serviceAddCmd.Flags().StringVar(&awsPrefix, "aws-prefix", "", "AWS Secrets Manager prefix for this service")
}

func runServiceList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Also try to load from services.yaml
	servicesData, _ := loadServicesFile()

	fmt.Println("Configured services:")
	fmt.Println()

	count := 0

	// Services from main config
	for name, svc := range cfg.Services {
		count++
		fmt.Printf("  %s\n", name)
		if svc.ConsulKey != "" {
			fmt.Printf("    consul_key: %s\n", svc.ConsulKey)
		}
		if svc.AWSPrefix != "" {
			fmt.Printf("    aws_prefix: %s\n", svc.AWSPrefix)
		}
		fmt.Println()
	}

	// Services from services.yaml
	if services, ok := servicesData["services"].(map[string]interface{}); ok {
		for name, svcData := range services {
			// Skip if already in main config
			if _, exists := cfg.Services[name]; exists {
				continue
			}
			count++
			fmt.Printf("  %s\n", name)
			if svc, ok := svcData.(map[string]interface{}); ok {
				if consulKey, ok := svc["consul_key"].(string); ok {
					fmt.Printf("    consul_key: %s\n", consulKey)
				}
				if awsPrefix, ok := svc["aws_prefix"].(string); ok {
					fmt.Printf("    aws_prefix: %s\n", awsPrefix)
				}
			}
			fmt.Println()
		}
	}

	if count == 0 {
		fmt.Println("  (no services configured)")
		fmt.Println()
		fmt.Println("Add a service with:")
		fmt.Println("  sreq service add <name> --consul-key <key> --aws-prefix <prefix>")
	}

	return nil
}

func runServiceAdd(cmd *cobra.Command, args []string) error {
	name := args[0]

	if consulKey == "" && awsPrefix == "" {
		return fmt.Errorf("at least one of --consul-key or --aws-prefix is required")
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	servicesPath := filepath.Join(configDir, config.DefaultServicesFile)

	// Load existing services
	data, err := loadServicesFile()
	if err != nil {
		// Start with empty structure if file doesn't exist
		data = map[string]interface{}{
			"services": map[string]interface{}{},
		}
	}

	services, ok := data["services"].(map[string]interface{})
	if !ok {
		services = map[string]interface{}{}
		data["services"] = services
	}

	// Check if service already exists
	if _, exists := services[name]; exists {
		return fmt.Errorf("service '%s' already exists", name)
	}

	// Add new service
	svcConfig := map[string]string{}
	if consulKey != "" {
		svcConfig["consul_key"] = consulKey
	}
	if awsPrefix != "" {
		svcConfig["aws_prefix"] = awsPrefix
	}
	services[name] = svcConfig

	// Write back
	output, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal services: %w", err)
	}

	if err := os.WriteFile(servicesPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write services file: %w", err)
	}

	fmt.Printf("Added service: %s\n", name)
	if consulKey != "" {
		fmt.Printf("  consul_key: %s\n", consulKey)
	}
	if awsPrefix != "" {
		fmt.Printf("  aws_prefix: %s\n", awsPrefix)
	}

	return nil
}

func runServiceRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	servicesPath := filepath.Join(configDir, config.DefaultServicesFile)

	data, err := loadServicesFile()
	if err != nil {
		return fmt.Errorf("failed to load services: %w", err)
	}

	services, ok := data["services"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("service '%s' not found", name)
	}

	if _, exists := services[name]; !exists {
		return fmt.Errorf("service '%s' not found", name)
	}

	delete(services, name)

	output, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal services: %w", err)
	}

	if err := os.WriteFile(servicesPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write services file: %w", err)
	}

	fmt.Printf("Removed service: %s\n", name)
	return nil
}

func loadServicesFile() (map[string]interface{}, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}

	servicesPath := filepath.Join(configDir, config.DefaultServicesFile)
	fileData, err := os.ReadFile(servicesPath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(fileData, &data); err != nil {
		return nil, err
	}

	return data, nil
}
