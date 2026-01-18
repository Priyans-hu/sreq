package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/Priyans-hu/sreq/internal/resolver"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and manage sreq configuration.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current sreq configuration including providers and settings.`,
	RunE:  runConfigShow,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	RunE:  runConfigPath,
}

var configTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test provider authentication",
	Long: `Verify that authentication is properly configured for all providers.

This command checks:
  - Consul: connectivity and token validity
  - AWS: credentials and Secrets Manager access

Examples:
  sreq config test           # Test all providers
  sreq config test --verbose # Show detailed output`,
	RunE: runConfigTest,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configTestCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	configDir, _ := config.GetConfigDir()
	fmt.Printf("Configuration: %s\n\n", filepath.Join(configDir, config.DefaultConfigFile))

	// Display providers
	fmt.Println("Providers:")
	if len(cfg.Providers) == 0 {
		fmt.Println("  (none configured)")
	} else {
		for name, provider := range cfg.Providers {
			fmt.Printf("  %s:\n", name)
			if provider.Address != "" {
				fmt.Printf("    address: %s\n", provider.Address)
			}
			if provider.Region != "" {
				fmt.Printf("    region: %s\n", provider.Region)
			}
			if provider.Profile != "" {
				fmt.Printf("    profile: %s\n", provider.Profile)
			}
			if len(provider.Paths) > 0 {
				fmt.Println("    paths:")
				for key, path := range provider.Paths {
					fmt.Printf("      %s: %s\n", key, path)
				}
			}
		}
	}

	// Display environments
	fmt.Println("\nEnvironments:")
	if len(cfg.Environments) == 0 {
		fmt.Println("  (none configured)")
	} else {
		for _, env := range cfg.Environments {
			if env == cfg.DefaultEnv {
				fmt.Printf("  - %s (default)\n", env)
			} else {
				fmt.Printf("  - %s\n", env)
			}
		}
	}

	// Display services count
	fmt.Printf("\nServices: %d configured\n", len(cfg.Services))
	if verbose && len(cfg.Services) > 0 {
		for name := range cfg.Services {
			fmt.Printf("  - %s\n", name)
		}
	}

	return nil
}

func runConfigPath(cmd *cobra.Command, args []string) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, config.DefaultConfigFile)

	// Check if using custom path
	if envPath := os.Getenv("SREQ_CONFIG"); envPath != "" {
		configPath = envPath
		fmt.Printf("%s (from SREQ_CONFIG)\n", configPath)
	} else {
		fmt.Println(configPath)
	}

	return nil
}

func runConfigTest(cmd *cobra.Command, args []string) error {
	fmt.Println("Testing provider authentication...")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create resolver to initialize providers
	res, err := resolver.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize providers: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run health checks
	results := res.HealthCheck(ctx)

	// Track overall status
	allPassed := true
	testedCount := 0

	// Check Consul
	if _, exists := cfg.Providers["consul"]; exists {
		testedCount++
		fmt.Print("Consul: ")
		if err, hasErr := results["consul"]; hasErr && err != nil {
			fmt.Printf("FAILED\n")
			fmt.Printf("  Error: %v\n", err)
			if verbose {
				fmt.Printf("  Address: %s\n", cfg.Providers["consul"].Address)
			}
			allPassed = false
		} else if _, ok := results["consul"]; ok {
			fmt.Printf("OK\n")
			if verbose {
				fmt.Printf("  Address: %s\n", cfg.Providers["consul"].Address)
			}
		} else {
			fmt.Printf("SKIPPED (not initialized)\n")
		}
	}

	// Check AWS
	if _, exists := cfg.Providers["aws_secrets"]; exists {
		testedCount++
		fmt.Print("AWS Secrets Manager: ")
		// Check for aws or aws_secrets key in results
		var awsErr error
		var awsFound bool
		if err, ok := results["aws"]; ok {
			awsErr = err
			awsFound = true
		} else if err, ok := results["aws_secrets"]; ok {
			awsErr = err
			awsFound = true
		}

		if awsFound {
			if awsErr != nil {
				fmt.Printf("FAILED\n")
				fmt.Printf("  Error: %v\n", awsErr)
				if verbose {
					fmt.Printf("  Region: %s\n", cfg.Providers["aws_secrets"].Region)
					if cfg.Providers["aws_secrets"].Profile != "" {
						fmt.Printf("  Profile: %s\n", cfg.Providers["aws_secrets"].Profile)
					}
				}
				allPassed = false
			} else {
				fmt.Printf("OK\n")
				if verbose {
					fmt.Printf("  Region: %s\n", cfg.Providers["aws_secrets"].Region)
				}
			}
		} else {
			fmt.Printf("SKIPPED (not initialized)\n")
		}
	} else if _, exists := cfg.Providers["aws"]; exists {
		testedCount++
		fmt.Print("AWS Secrets Manager: ")
		if err, ok := results["aws"]; ok && err != nil {
			fmt.Printf("FAILED\n")
			fmt.Printf("  Error: %v\n", err)
			allPassed = false
		} else if _, ok := results["aws"]; ok {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("SKIPPED (not initialized)\n")
		}
	}

	fmt.Println()

	if testedCount == 0 {
		fmt.Println("No providers configured. Run 'sreq auth' to set up authentication.")
		return nil
	}

	if allPassed {
		fmt.Println("All provider tests passed!")
	} else {
		fmt.Println("Some provider tests failed. Check your configuration with 'sreq auth'.")
		return fmt.Errorf("authentication test failed")
	}

	return nil
}

// LoadServicesConfig loads the services configuration file
func LoadServicesConfig() (map[string]interface{}, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}

	servicesPath := filepath.Join(configDir, config.DefaultServicesFile)
	data, err := os.ReadFile(servicesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("services file not found: %s (run 'sreq init' to create)", servicesPath)
		}
		return nil, err
	}

	var services map[string]interface{}
	if err := yaml.Unmarshal(data, &services); err != nil {
		return nil, fmt.Errorf("failed to parse services file: %w", err)
	}

	return services, nil
}
