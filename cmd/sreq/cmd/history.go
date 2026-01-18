package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/Priyans-hu/sreq/internal/history"
	"github.com/spf13/cobra"
)

var (
	historyService string
	historyEnv     string
	historyMethod  string
	historyAll     bool
	historyClear   bool
	historyBefore  string
	historyCurl    bool
	historyHTTPie  bool
	historyReplay  bool
)

var historyCmd = &cobra.Command{
	Use:   "history [id]",
	Short: "View and manage request history",
	Long: `View, replay, and manage request history.

Without arguments, shows the last 20 requests. With an ID, shows details
of that specific request.

Examples:
  sreq history                    # List recent requests (last 20)
  sreq history --all              # List all requests
  sreq history --service auth     # Filter by service
  sreq history --env prod         # Filter by env
  sreq history 5                  # Show details of request #5
  sreq history 5 --curl           # Export request #5 as curl command
  sreq history 5 --httpie         # Export request #5 as HTTPie command
  sreq history 5 --replay         # Replay request #5
  sreq history --clear            # Clear all history
  sreq history --clear --before 7d  # Clear entries older than 7 days`,
	RunE: runHistory,
}

func init() {
	rootCmd.AddCommand(historyCmd)

	historyCmd.Flags().StringVar(&historyService, "service", "", "Filter by service name")
	historyCmd.Flags().StringVar(&historyEnv, "env", "", "Filter by environment")
	historyCmd.Flags().StringVar(&historyMethod, "method", "", "Filter by HTTP method")
	historyCmd.Flags().BoolVar(&historyAll, "all", false, "Show all history entries")
	historyCmd.Flags().BoolVar(&historyClear, "clear", false, "Clear history")
	historyCmd.Flags().StringVar(&historyBefore, "before", "", "Clear entries older than duration (e.g., 7d, 24h)")
	historyCmd.Flags().BoolVar(&historyCurl, "curl", false, "Export as curl command")
	historyCmd.Flags().BoolVar(&historyHTTPie, "httpie", false, "Export as HTTPie command")
	historyCmd.Flags().BoolVar(&historyReplay, "replay", false, "Replay the request")
}

func runHistory(cmd *cobra.Command, args []string) error {
	// Check if history is disabled
	if os.Getenv("SREQ_NO_HISTORY") == "1" {
		return fmt.Errorf("history is disabled (SREQ_NO_HISTORY=1)")
	}

	// Get config directory
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	// Load history
	h, err := history.New(configDir)
	if err != nil {
		return err
	}

	// Handle clear
	if historyClear {
		return handleClear(h)
	}

	// Handle specific entry by ID
	if len(args) > 0 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid history ID: %s", args[0])
		}
		return handleEntry(h, id)
	}

	// List entries
	return handleList(h)
}

func handleClear(h *history.History) error {
	if historyBefore != "" {
		duration, err := parseDuration(historyBefore)
		if err != nil {
			return fmt.Errorf("invalid duration: %s (use format like 7d, 24h, 30m)", historyBefore)
		}
		removed := h.ClearBefore(duration)
		if err := h.Save(); err != nil {
			return err
		}
		fmt.Printf("Cleared %d entries older than %s\n", removed, historyBefore)
		return nil
	}

	h.Clear()
	if err := h.Save(); err != nil {
		return err
	}
	fmt.Println("History cleared")
	return nil
}

func handleEntry(h *history.History, id int) error {
	entry, err := h.Get(id)
	if err != nil {
		return err
	}

	// Export as curl
	if historyCurl {
		fmt.Println(entry.ToCurl())
		return nil
	}

	// Export as HTTPie
	if historyHTTPie {
		fmt.Println(entry.ToHTTPie())
		return nil
	}

	// Replay
	if historyReplay {
		return replayEntry(entry)
	}

	// Show details
	printEntryDetails(entry)
	return nil
}

func handleList(h *history.History) error {
	limit := 20
	if historyAll {
		limit = 0
	}

	entries := h.List(history.ListOptions{
		Service: historyService,
		Env:     historyEnv,
		Method:  historyMethod,
		Limit:   limit,
	})

	if len(entries) == 0 {
		fmt.Println("No history entries found")
		if historyService != "" || historyEnv != "" || historyMethod != "" {
			fmt.Println("Try removing filters or run some requests first")
		}
		return nil
	}

	// Print header
	fmt.Printf("%-4s  %-6s  %-30s  %-15s  %-6s  %-8s  %s\n",
		"ID", "METHOD", "PATH", "SERVICE", "ENV", "STATUS", "TIME")
	fmt.Println(strings.Repeat("-", 90))

	// Print entries
	for _, e := range entries {
		path := e.Path
		if len(path) > 28 {
			path = path[:25] + "..."
		}

		service := e.Service
		if len(service) > 13 {
			service = service[:10] + "..."
		}

		statusStr := "-"
		if e.Status > 0 {
			statusStr = fmt.Sprintf("%s%d%s", e.StatusColor(), e.Status, history.ResetColor())
		}

		fmt.Printf("%-4d  %-6s  %-30s  %-15s  %-6s  %-8s  %s\n",
			e.ID,
			e.Method,
			path,
			service,
			e.Env,
			statusStr,
			e.FormatDuration(),
		)
	}

	// Show count
	total := h.Count()
	shown := len(entries)
	if shown < total && !historyAll {
		fmt.Printf("\nShowing %d of %d entries. Use --all to see all.\n", shown, total)
	}

	return nil
}

func printEntryDetails(e *history.Entry) {
	fmt.Printf("Request #%d\n", e.ID)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Time:     %s\n", e.Timestamp.Format(time.RFC3339))
	fmt.Printf("Service:  %s\n", e.Service)
	fmt.Printf("Env:      %s\n", e.Env)
	fmt.Printf("Method:   %s\n", e.Method)
	fmt.Printf("Path:     %s\n", e.Path)
	if e.BaseURL != "" {
		fmt.Printf("Base URL: %s\n", e.BaseURL)
	}
	if e.Status > 0 {
		fmt.Printf("Status:   %s%d%s\n", e.StatusColor(), e.Status, history.ResetColor())
	}
	if e.Duration > 0 {
		fmt.Printf("Duration: %s\n", e.FormatDuration())
	}

	// Request details
	if e.Request != nil {
		if len(e.Request.Headers) > 0 {
			fmt.Println("\nRequest Headers:")
			for k, v := range e.Request.Headers {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
		if e.Request.Body != "" {
			fmt.Println("\nRequest Body:")
			fmt.Printf("  %s\n", e.Request.Body)
		}
	}

	// Response details
	if e.Response != nil {
		if e.Response.Status != "" {
			fmt.Printf("\nResponse: %s", e.Response.Status)
			if e.Response.SizeBytes > 0 {
				fmt.Printf(" (%d bytes)", e.Response.SizeBytes)
			}
			fmt.Println()
		}
	}

	// Export hints
	fmt.Println()
	fmt.Printf("Export: sreq history %d --curl\n", e.ID)
	fmt.Printf("Replay: sreq history %d --replay\n", e.ID)
}

func replayEntry(e *history.Entry) error {
	fmt.Printf("Replaying request #%d: %s %s\n", e.ID, e.Method, e.Path)
	fmt.Println()

	// Build args for run command
	args := []string{e.Method, e.Path}

	// Set global flags (from root.go)
	serviceName = e.Service
	environment = e.Env

	// Run the request
	return runRun(nil, args)
}

// parseDuration parses a duration string like "7d", "24h", "30m"
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Handle days specially
	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	// Handle weeks
	if strings.HasSuffix(s, "w") {
		weeks, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(weeks) * 7 * 24 * time.Hour, nil
	}

	// Standard Go duration parsing
	return time.ParseDuration(s)
}
