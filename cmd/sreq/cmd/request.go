package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:   "run [METHOD] [path]",
	Short: "Make an HTTP request",
	Long: `Make an HTTP request to a service endpoint.

The run command resolves credentials from configured providers
and makes the HTTP request.

Examples:
  sreq run GET /api/v1/users -s auth-service -e dev
  sreq run POST /api/v1/users -s auth-service -d '{"name":"test"}'
  sreq run GET /health -s billing-service --verbose`,
	Args: cobra.MinimumNArgs(2),
	RunE: runRequest,
}

var (
	requestData    string
	requestHeaders []string
	outputFormat   string
)

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringVarP(&requestData, "data", "d", "", "Request body (or @filename for file)")
	requestCmd.Flags().StringArrayVarP(&requestHeaders, "header", "H", nil, "Add header (repeatable)")
	requestCmd.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json/raw/headers)")
}

func runRequest(cmd *cobra.Command, args []string) error {
	method := strings.ToUpper(args[0])
	path := args[1]

	// Validate method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}
	if !validMethods[method] {
		return fmt.Errorf("invalid HTTP method: %s", method)
	}

	if serviceName == "" {
		return fmt.Errorf("--service (-s) is required\n\nUsage: sreq run %s %s -s <service-name>", method, path)
	}

	// Use default environment if not specified
	if environment == "" {
		environment = "dev"
	}

	// Parse headers
	headers := make(map[string]string)
	for _, h := range requestHeaders {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s (expected 'Key: Value')", h)
		}
		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	// Handle @filename for body
	body := requestData
	if strings.HasPrefix(body, "@") {
		filename := body[1:]
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read body file: %w", err)
		}
		body = string(data)
	}

	if verbose {
		fmt.Println("Request Details:")
		fmt.Printf("  Service:     %s\n", serviceName)
		fmt.Printf("  Environment: %s\n", environment)
		fmt.Printf("  Method:      %s\n", method)
		fmt.Printf("  Path:        %s\n", path)
		if len(headers) > 0 {
			fmt.Println("  Headers:")
			for k, v := range headers {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}
		if body != "" {
			fmt.Printf("  Body:        %s\n", truncate(body, 100))
		}
		fmt.Println()
	}

	if dryRun {
		fmt.Println("[DRY RUN] Would execute:")
		fmt.Printf("  %s https://<base-url-from-%s>%s\n", method, serviceName, path)
		fmt.Println()
		fmt.Println("Credentials would be fetched from:")
		fmt.Println("  - Consul: services/<service>/config/*")
		fmt.Println("  - AWS Secrets Manager: <service>/<env>/credentials")
		return nil
	}

	// TODO: Implement actual request logic
	// 1. Load config
	// 2. Resolve credentials from providers
	// 3. Build full URL
	// 4. Make HTTP request
	// 5. Display response

	fmt.Printf("Making %s request to %s...\n", method, path)
	fmt.Printf("Service: %s | Environment: %s\n", serviceName, environment)
	fmt.Println()
	fmt.Println("Note: Provider integration not yet implemented.")
	fmt.Println("This will fetch credentials from Consul/AWS and make the actual request.")

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
