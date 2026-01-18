package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Priyans-hu/sreq/internal/config"
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

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
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
