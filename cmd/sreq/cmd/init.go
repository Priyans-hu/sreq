package cmd

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sreq configuration",
	Long: `Creates the default configuration files in ~/.sreq/

This command will:
  - Create ~/.sreq/ directory
  - Generate config.yaml with provider templates
  - Generate services.yaml for service definitions
  - Generate encryption key for credential caching`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Check if already initialized
	configPath := filepath.Join(configDir, config.DefaultConfigFile)
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration already exists at %s\n", configDir)
		fmt.Println("Use --force to reinitialize (this will not overwrite existing files)")
		return nil
	}

	fmt.Println("Initializing sreq configuration...")

	// Create config directory and files
	if err := config.Init(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Generate encryption key
	if err := generateEncryptionKey(configDir); err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Create cache directory
	cacheDir := filepath.Join(configDir, "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	fmt.Println()
	fmt.Printf("Created %s/\n", configDir)
	fmt.Printf("  ├── config.yaml      # Provider configuration\n")
	fmt.Printf("  ├── services.yaml    # Service definitions\n")
	fmt.Printf("  ├── .key             # Encryption key (do not share!)\n")
	fmt.Printf("  └── cache/           # Credential cache directory\n")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit config.yaml to configure your providers (Consul, AWS, etc.)")
	fmt.Println("  2. Add services to services.yaml")
	fmt.Println("  3. Run: sreq GET /api/health -s <service-name> -e dev")

	return nil
}

// generateEncryptionKey creates a 32-byte random key for AES-256 encryption
func generateEncryptionKey(configDir string) error {
	keyPath := filepath.Join(configDir, ".key")

	// Don't overwrite existing key
	if _, err := os.Stat(keyPath); err == nil {
		return nil
	}

	key := make([]byte, 32) // 256 bits for AES-256
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate random key: %w", err)
	}

	// Write with restrictive permissions (owner read/write only)
	if err := os.WriteFile(keyPath, key, 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}
