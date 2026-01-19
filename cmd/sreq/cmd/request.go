package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Priyans-hu/sreq/internal/cache"
	"github.com/Priyans-hu/sreq/internal/client"
	"github.com/Priyans-hu/sreq/internal/config"
	sreerrors "github.com/Priyans-hu/sreq/internal/errors"
	"github.com/Priyans-hu/sreq/internal/history"
	"github.com/Priyans-hu/sreq/internal/resolver"
	"github.com/Priyans-hu/sreq/pkg/types"
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
	RunE: runRun,
}

var (
	requestData    string
	requestHeaders []string
	outputFormat   string
	timeout        time.Duration
	offlineMode    bool
	noCache        bool
)

func init() {
	rootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringVarP(&requestData, "data", "d", "", "Request body (or @filename for file)")
	requestCmd.Flags().StringArrayVarP(&requestHeaders, "header", "H", nil, "Add header (repeatable)")
	requestCmd.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json/raw/headers)")
	requestCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")
	requestCmd.Flags().BoolVar(&offlineMode, "offline", false, "Use cached credentials only (no provider calls)")
	requestCmd.Flags().BoolVar(&noCache, "no-cache", false, "Skip cache and fetch fresh credentials")
}

func runRun(cmd *cobra.Command, args []string) error {
	method := strings.ToUpper(args[0])
	path := args[1]
	startTime := time.Now()

	// Validate method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}
	if !validMethods[method] {
		return sreerrors.InvalidMethod(method)
	}

	if serviceName == "" {
		return sreerrors.MissingRequiredFlag("service")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Load context values (from -c flag or default_context)
	ctxName := contextName
	if ctxName == "" {
		ctxName = cfg.DefaultContext
	}

	// Apply context values as base, then override with explicit CLI flags
	if ctxName != "" {
		if ctx, exists := cfg.Contexts[ctxName]; exists {
			if environment == "" {
				environment = ctx.Env
			}
			if region == "" {
				region = ctx.Region
			}
			if project == "" {
				project = ctx.Project
			}
			if app == "" {
				app = ctx.App
			}
		} else if contextName != "" {
			// Only error if user explicitly specified a context that doesn't exist
			return sreerrors.ContextNotFound(contextName)
		}
	}

	// Use default environment if still not specified
	if environment == "" {
		environment = cfg.DefaultEnv
		if environment == "" {
			environment = "dev"
		}
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
		if ctxName != "" {
			fmt.Printf("  Context:     %s\n", ctxName)
		}
		fmt.Printf("  Service:     %s\n", serviceName)
		fmt.Printf("  Environment: %s\n", environment)
		if region != "" {
			fmt.Printf("  Region:      %s\n", region)
		}
		if project != "" {
			fmt.Printf("  Project:     %s\n", project)
		}
		if app != "" {
			fmt.Printf("  App:         %s\n", app)
		}
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
		fmt.Printf("  %s <base-url>%s\n", method, path)
		fmt.Println()
		fmt.Println("Credentials would be resolved from configured providers.")
		return nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Get config directory for cache
	configDir, _ := config.GetConfigDir()

	// Try to get credentials from cache first (unless --no-cache)
	var creds *types.ResolvedCredentials
	var credCache *cache.Cache
	useCache := !noCache && !cache.IsDisabled() && cache.KeyExists(configDir)

	if useCache {
		credCache, _ = cache.New(cache.Config{ConfigDir: configDir})
		if credCache != nil {
			creds, _ = credCache.Get(serviceName, environment)
			if creds != nil && verbose {
				fmt.Println("Using cached credentials")
			}
		}
	}

	// If offline mode, we must have cached credentials
	if offlineMode {
		if creds == nil {
			return fmt.Errorf("no cached credentials found for %s/%s (run 'sreq sync %s' first)",
				serviceName, environment, environment)
		}
	}

	// If no cached credentials, resolve from providers
	if creds == nil {
		// Create resolver
		res, err := resolver.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to create resolver: %w", err)
		}

		// Resolve credentials
		if verbose {
			fmt.Println("Resolving credentials from providers...")
		}

		creds, err = res.Resolve(ctx, resolver.ResolveOptions{
			Service: serviceName,
			Env:     environment,
			Region:  region,
			Project: project,
			App:     app,
		})
		if err != nil {
			return sreerrors.CredentialResolutionFailed(serviceName, environment, err)
		}

		// Cache the credentials for next time (unless --no-cache)
		if useCache && credCache != nil {
			if err := credCache.Set(serviceName, environment, creds); err != nil && verbose {
				fmt.Printf("Warning: failed to cache credentials: %v\n", err)
			}
		}
	}

	if verbose {
		fmt.Printf("  Base URL:  %s\n", creds.BaseURL)
		if creds.Username != "" {
			fmt.Printf("  Username:  %s\n", creds.Username)
			fmt.Printf("  Password:  %s\n", maskPassword(creds.Password))
		}
		if creds.APIKey != "" {
			fmt.Printf("  API Key:   %s\n", maskPassword(creds.APIKey))
		}
		fmt.Println()
	}

	if creds.BaseURL == "" {
		return sreerrors.BaseURLMissing(serviceName, environment)
	}

	// Create HTTP client
	httpClient := client.New(
		client.WithTimeout(timeout),
		client.WithVerbose(verbose),
	)

	// Build request
	req := &types.Request{
		Method:      method,
		Path:        path,
		Service:     serviceName,
		Environment: environment,
		Body:        body,
		Headers:     headers,
	}

	// Execute request
	resp, err := httpClient.Do(ctx, req, creds)
	duration := time.Since(startTime).Milliseconds()

	// Save to history (unless disabled or dry run)
	if os.Getenv("SREQ_NO_HISTORY") != "1" && !dryRun {
		saveHistory(method, path, serviceName, environment, creds.BaseURL, headers, body, resp, duration)
	}

	if err != nil {
		return sreerrors.RequestFailed(creds.BaseURL+path, err)
	}

	// Output response
	return outputResponse(resp, outputFormat)
}

// saveHistory saves the request to history
func saveHistory(method, path, service, env, baseURL string, headers map[string]string, body string, resp *types.Response, durationMs int64) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return // Silently fail - history is optional
	}

	h, err := history.New(configDir)
	if err != nil {
		return
	}

	entry := history.Entry{
		Timestamp: time.Now(),
		Service:   service,
		Env:       env,
		Method:    method,
		Path:      path,
		BaseURL:   baseURL,
		Duration:  durationMs,
		Request: &history.Request{
			Headers: headers,
			Body:    body,
		},
	}

	// Add response info if available
	if resp != nil {
		entry.Status = resp.StatusCode
		entry.Response = &history.Response{
			Status:    resp.Status,
			SizeBytes: len(resp.Body),
		}
	}

	h.Add(entry)
	_ = h.Save() // Ignore save errors - history is optional
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func maskPassword(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

func outputResponse(resp *types.Response, format string) error {
	switch format {
	case "json":
		// Try to pretty-print if it's JSON
		var jsonData interface{}
		if err := json.Unmarshal(resp.Body, &jsonData); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
				return nil
			}
		}
		// Fall back to raw if not valid JSON
		fmt.Println(string(resp.Body))

	case "raw":
		fmt.Println(string(resp.Body))

	case "headers":
		fmt.Printf("HTTP %s\n", resp.Status)
		for key, values := range resp.Headers {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		fmt.Println()
		fmt.Println(string(resp.Body))

	default:
		return fmt.Errorf("unknown output format: %s", format)
	}

	return nil
}
