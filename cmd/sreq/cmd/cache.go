package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Priyans-hu/sreq/internal/cache"
	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/Priyans-hu/sreq/internal/resolver"
	"github.com/spf13/cobra"
)

var (
	syncAll   bool
	syncForce bool
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage credential cache",
	Long: `Manage the local credential cache.

The cache stores encrypted credentials locally for faster access
and offline use.

Examples:
  sreq cache status     # Show cache status
  sreq cache clear      # Clear all cached credentials
  sreq cache clear dev  # Clear cache for dev environment`,
}

var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache status",
	RunE:  runCacheStatus,
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear [env]",
	Short: "Clear cached credentials",
	Long: `Clear cached credentials.

Without arguments, clears all cached credentials.
With an environment argument, clears only that environment's cache.

Examples:
  sreq cache clear        # Clear all
  sreq cache clear dev    # Clear only dev environment`,
	RunE: runCacheClear,
}

var syncCmd = &cobra.Command{
	Use:   "sync [env]",
	Short: "Sync credentials to local cache",
	Long: `Sync credentials from providers to local cache.

This fetches credentials for all configured services and caches
them locally for faster access and offline use.

Examples:
  sreq sync dev         # Sync credentials for dev environment
  sreq sync --all       # Sync all environments
  sreq sync dev --force # Force refresh even if cache is valid`,
	RunE: runSync,
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(syncCmd)

	cacheCmd.AddCommand(cacheStatusCmd)
	cacheCmd.AddCommand(cacheClearCmd)

	syncCmd.Flags().BoolVar(&syncAll, "all", false, "Sync all environments")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "Force refresh even if cache is valid")
}

func runCacheStatus(cmd *cobra.Command, args []string) error {
	if cache.IsDisabled() {
		fmt.Println("Cache is disabled (SREQ_NO_CACHE=1 or CI environment detected)")
		return nil
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	// Check if key exists
	if !cache.KeyExists(configDir) {
		fmt.Println("Cache not initialized (run 'sreq init' first)")
		return nil
	}

	c, err := cache.New(cache.Config{
		ConfigDir: configDir,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	status, err := c.Status()
	if err != nil {
		return fmt.Errorf("failed to get cache status: %w", err)
	}

	fmt.Println("Cache Status")
	fmt.Println("============")
	fmt.Printf("Enabled:    %v\n", status.Enabled)
	fmt.Printf("Directory:  %s\n", status.CacheDir)
	fmt.Printf("TTL:        %v\n", status.TTL)
	fmt.Printf("Entries:    %d\n", status.EntryCount)
	fmt.Printf("Total Size: %d bytes\n", status.TotalSize)

	if len(status.Entries) > 0 {
		fmt.Println("\nCached Entries:")
		for _, e := range status.Entries {
			expiredStr := ""
			if e.Expired {
				expiredStr = " (EXPIRED)"
			}
			fmt.Printf("  %s/%s - cached %s, expires %s%s\n",
				e.Service,
				e.Env,
				e.CachedAt.Format("15:04:05"),
				e.ExpiresAt.Format("15:04:05"),
				expiredStr,
			)
		}
	}

	return nil
}

func runCacheClear(cmd *cobra.Command, args []string) error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	if !cache.KeyExists(configDir) {
		fmt.Println("Cache not initialized")
		return nil
	}

	c, err := cache.New(cache.Config{
		ConfigDir: configDir,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	if len(args) > 0 {
		// Clear specific environment
		env := args[0]
		if err := c.ClearEnv(env); err != nil {
			return err
		}
		fmt.Printf("Cleared cache for environment: %s\n", env)
	} else {
		// Clear all
		if err := c.Clear(); err != nil {
			return err
		}
		fmt.Println("Cleared all cached credentials")
	}

	return nil
}

func runSync(cmd *cobra.Command, args []string) error {
	if cache.IsDisabled() {
		return fmt.Errorf("cache is disabled (SREQ_NO_CACHE=1 or CI environment detected)")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	// Initialize cache
	c, err := cache.New(cache.Config{
		ConfigDir: configDir,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Create resolver
	res, err := resolver.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create resolver: %w", err)
	}

	// Determine environments to sync
	var envs []string
	if syncAll {
		envs = cfg.Environments
	} else if len(args) > 0 {
		envs = []string{args[0]}
	} else {
		// Default to current env
		envs = []string{cfg.DefaultEnv}
		if envs[0] == "" {
			envs = []string{"dev"}
		}
	}

	// Sync each environment
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	total := 0
	errors := 0

	for _, env := range envs {
		fmt.Printf("Syncing %s environment...\n", env)

		for serviceName := range cfg.Services {
			// Check if already cached (unless force)
			if !syncForce {
				cached, _ := c.Get(serviceName, env)
				if cached != nil {
					fmt.Printf("  %s: cached (use --force to refresh)\n", serviceName)
					continue
				}
			}

			// Resolve credentials
			creds, err := res.Resolve(ctx, resolver.ResolveOptions{
				Service: serviceName,
				Env:     env,
			})

			if err != nil {
				fmt.Printf("  %s: FAILED - %v\n", serviceName, err)
				errors++
				continue
			}

			// Cache credentials
			if err := c.Set(serviceName, env, creds); err != nil {
				fmt.Printf("  %s: FAILED to cache - %v\n", serviceName, err)
				errors++
				continue
			}

			fmt.Printf("  %s: synced\n", serviceName)
			total++
		}
	}

	fmt.Println()
	if errors > 0 {
		fmt.Printf("Synced %d credentials with %d errors\n", total, errors)
	} else {
		fmt.Printf("Synced %d credentials successfully\n", total)
	}

	return nil
}
