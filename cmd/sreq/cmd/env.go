package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
	Long:  `List environments or switch the default environment.`,
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available environments",
	RunE:  runEnvList,
}

var envSwitchCmd = &cobra.Command{
	Use:   "switch [env]",
	Short: "Switch default environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnvSwitch,
}

var envCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current default environment",
	RunE:  runEnvCurrent,
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envSwitchCmd)
	envCmd.AddCommand(envCurrentCmd)
}

func runEnvList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("Available environments:")
	fmt.Println()

	if len(cfg.Environments) == 0 {
		fmt.Println("  (none configured)")
		fmt.Println()
		fmt.Println("Add environments to your config.yaml:")
		fmt.Println("  environments:")
		fmt.Println("    - dev")
		fmt.Println("    - staging")
		fmt.Println("    - prod")
		return nil
	}

	for _, env := range cfg.Environments {
		if env == cfg.DefaultEnv {
			fmt.Printf("  * %s (default)\n", env)
		} else {
			fmt.Printf("    %s\n", env)
		}
	}

	return nil
}

func runEnvSwitch(cmd *cobra.Command, args []string) error {
	newEnv := args[0]

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Validate environment exists
	valid := false
	for _, env := range cfg.Environments {
		if env == newEnv {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("environment '%s' not found in config\n\nAvailable: %v", newEnv, cfg.Environments)
	}

	// Update config file
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, config.DefaultConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var configData map[string]interface{}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	configData["default_env"] = newEnv

	output, err := yaml.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Switched default environment to: %s\n", newEnv)
	return nil
}

func runEnvCurrent(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.DefaultEnv == "" {
		fmt.Println("No default environment set")
		return nil
	}

	fmt.Println(cfg.DefaultEnv)
	return nil
}
