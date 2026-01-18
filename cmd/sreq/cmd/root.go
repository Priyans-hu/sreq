package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Global flags
	serviceName string
	environment string
	verbose     bool
	dryRun      bool
)

var rootCmd = &cobra.Command{
	Use:   "sreq",
	Short: "Service-aware API client with automatic credential resolution",
	Long: `sreq eliminates the overhead of manually fetching credentials from
multiple sources when testing APIs. Just specify the service name
and environment — sreq handles the rest.

Quick Start:
  sreq init                                    # Initialize configuration
  sreq service add auth-service --consul-key auth --aws-prefix auth
  sreq run GET /api/v1/users -s auth-service -e dev

Examples:
  sreq run GET /api/v1/users -s auth-service -e dev
  sreq run POST /api/v1/users -s auth-service -d '{"name":"test"}'
  sreq run GET /health -s billing-service --verbose --dry-run`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&serviceName, "service", "s", "", "Service name")
	rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "", "Environment (dev/staging/prod)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be sent without executing")
}
