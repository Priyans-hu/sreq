package cmd

import (
	"fmt"

	"github.com/Priyans-hu/sreq/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open interactive terminal UI",
	Long: `Open an interactive terminal UI for sreq.

The TUI provides a visual interface for:
  - Browsing configured services
  - Viewing request history
  - Viewing request details and exporting as curl

Navigation:
  ↑/k, ↓/j  Navigate lists
  Enter     Select/view details
  Esc       Go back
  s         Services view
  h         History view
  c         Copy as curl (in detail view)
  q         Quit`,
	RunE: runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) error {
	model, err := tui.New()
	if err != nil {
		return fmt.Errorf("failed to initialize TUI: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
