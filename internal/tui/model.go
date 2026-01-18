package tui

import (
	"fmt"
	"strings"

	"github.com/Priyans-hu/sreq/internal/config"
	"github.com/Priyans-hu/sreq/internal/history"
	"github.com/Priyans-hu/sreq/pkg/types"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View represents the current view in the TUI
type View int

const (
	ViewDashboard View = iota
	ViewServices
	ViewHistory
	ViewHistoryDetail
)

// Model represents the TUI state
type Model struct {
	config       *types.Config
	history      *history.History
	services     []string
	historyItems []history.Entry

	// UI state
	view            View
	serviceList     list.Model
	historyList     list.Model
	selectedHistory *history.Entry
	currentEnv      string

	// Dimensions
	width  int
	height int

	// Error state
	err error
}

// serviceItem implements list.Item for services
type serviceItem struct {
	name string
}

func (i serviceItem) Title() string       { return i.name }
func (i serviceItem) Description() string { return "" }
func (i serviceItem) FilterValue() string { return i.name }

// historyItem implements list.Item for history entries
type historyItem struct {
	entry history.Entry
}

func (i historyItem) Title() string {
	return fmt.Sprintf("%s %s", i.entry.Method, i.entry.Path)
}

func (i historyItem) Description() string {
	status := "-"
	if i.entry.Status > 0 {
		status = fmt.Sprintf("%d", i.entry.Status)
	}
	return fmt.Sprintf("%s | %s | %s | %s",
		i.entry.Service,
		i.entry.Env,
		status,
		i.entry.FormatDuration(),
	)
}

func (i historyItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s", i.entry.Service, i.entry.Method, i.entry.Path)
}

// keyMap defines keyboard shortcuts
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Back     key.Binding
	Services key.Binding
	History  key.Binding
	Env      key.Binding
	Curl     key.Binding
	Quit     key.Binding
	Help     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Services: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "services"),
	),
	History: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "history"),
	),
	Env: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "switch env"),
	),
	Curl: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy curl"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// New creates a new TUI model
func New() (*Model, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Load history
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %w", err)
	}

	h, err := history.New(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load history: %w", err)
	}

	// Get services
	var services []string
	for name := range cfg.Services {
		services = append(services, name)
	}

	// Get history entries
	historyEntries := h.List(history.ListOptions{Limit: 50})

	// Create service list
	serviceItems := make([]list.Item, len(services))
	for i, s := range services {
		serviceItems[i] = serviceItem{name: s}
	}
	serviceList := list.New(serviceItems, list.NewDefaultDelegate(), 0, 0)
	serviceList.Title = "Services"
	serviceList.SetShowHelp(false)

	// Create history list
	historyListItems := make([]list.Item, len(historyEntries))
	for i, e := range historyEntries {
		historyListItems[i] = historyItem{entry: e}
	}
	historyList := list.New(historyListItems, list.NewDefaultDelegate(), 0, 0)
	historyList.Title = "Recent Requests"
	historyList.SetShowHelp(false)

	// Determine current env
	currentEnv := cfg.DefaultEnv
	if currentEnv == "" {
		currentEnv = "dev"
	}

	return &Model{
		config:       cfg,
		history:      h,
		services:     services,
		historyItems: historyEntries,
		view:         ViewDashboard,
		serviceList:  serviceList,
		historyList:  historyList,
		currentEnv:   currentEnv,
	}, nil
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update list sizes
		listHeight := m.height - 8
		listWidth := m.width - 4

		m.serviceList.SetSize(listWidth/2, listHeight)
		m.historyList.SetSize(listWidth/2, listHeight)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Back):
			if m.view != ViewDashboard {
				m.view = ViewDashboard
				m.selectedHistory = nil
			}
			return m, nil

		case key.Matches(msg, keys.Services):
			m.view = ViewServices
			return m, nil

		case key.Matches(msg, keys.History):
			m.view = ViewHistory
			return m, nil

		case key.Matches(msg, keys.Enter):
			return m.handleEnter()

		case key.Matches(msg, keys.Curl):
			if m.selectedHistory != nil {
				// In a real app, we'd copy to clipboard
				// For now, just show a message
			}
			return m, nil
		}
	}

	// Update the appropriate list based on current view
	var cmd tea.Cmd
	switch m.view {
	case ViewServices:
		m.serviceList, cmd = m.serviceList.Update(msg)
	case ViewHistory:
		m.historyList, cmd = m.historyList.Update(msg)
	case ViewDashboard:
		// In dashboard, history list is active
		m.historyList, cmd = m.historyList.Update(msg)
	}

	return m, cmd
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.view {
	case ViewHistory, ViewDashboard:
		if item, ok := m.historyList.SelectedItem().(historyItem); ok {
			m.selectedHistory = &item.entry
			m.view = ViewHistoryDetail
		}
	}
	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	switch m.view {
	case ViewDashboard:
		content = m.viewDashboard()
	case ViewServices:
		content = m.viewServices()
	case ViewHistory:
		content = m.viewHistory()
	case ViewHistoryDetail:
		content = m.viewHistoryDetail()
	}

	return content
}

func (m Model) viewDashboard() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("sreq - Service Request CLI")
	envBadge := fmt.Sprintf("env: %s", m.currentEnv)
	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		header,
		strings.Repeat(" ", max(0, m.width-lipgloss.Width(header)-lipgloss.Width(envBadge)-4)),
		mutedStyle.Render(envBadge),
	)
	b.WriteString(headerLine)
	b.WriteString("\n\n")

	// Stats
	stats := fmt.Sprintf("Services: %d | History: %d entries",
		len(m.services),
		m.history.Count(),
	)
	b.WriteString(mutedStyle.Render(stats))
	b.WriteString("\n\n")

	// Recent requests
	b.WriteString(subtitleStyle.Render("Recent Requests"))
	b.WriteString("\n")
	b.WriteString(m.historyList.View())

	// Help
	b.WriteString("\n")
	b.WriteString(m.helpView())

	return b.String()
}

func (m Model) viewServices() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Services"))
	b.WriteString("\n\n")
	b.WriteString(m.serviceList.View())
	b.WriteString("\n")
	b.WriteString(m.helpView())

	return b.String()
}

func (m Model) viewHistory() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Request History"))
	b.WriteString("\n\n")
	b.WriteString(m.historyList.View())
	b.WriteString("\n")
	b.WriteString(m.helpView())

	return b.String()
}

func (m Model) viewHistoryDetail() string {
	if m.selectedHistory == nil {
		return "No entry selected"
	}

	e := m.selectedHistory
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("Request #%d", e.ID)))
	b.WriteString("\n\n")

	// Details
	details := []string{
		fmt.Sprintf("Time:     %s", e.Timestamp.Format("2006-01-02 15:04:05")),
		fmt.Sprintf("Service:  %s", e.Service),
		fmt.Sprintf("Env:      %s", e.Env),
		fmt.Sprintf("Method:   %s", e.Method),
		fmt.Sprintf("Path:     %s", e.Path),
	}

	if e.BaseURL != "" {
		details = append(details, fmt.Sprintf("Base URL: %s", e.BaseURL))
	}

	if e.Status > 0 {
		statusLine := fmt.Sprintf("Status:   %d", e.Status)
		details = append(details, StatusStyle(e.Status).Render(statusLine))
	}

	if e.Duration > 0 {
		details = append(details, fmt.Sprintf("Duration: %s", e.FormatDuration()))
	}

	b.WriteString(strings.Join(details, "\n"))
	b.WriteString("\n\n")

	// Curl command
	b.WriteString(subtitleStyle.Render("curl command:"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render(e.ToCurl()))
	b.WriteString("\n\n")

	// Help
	b.WriteString(helpStyle.Render("[esc] back  [c] copy curl  [q] quit"))

	return b.String()
}

func (m Model) helpView() string {
	var help string
	switch m.view {
	case ViewDashboard:
		help = "[s] services  [h] history  [enter] view  [q] quit"
	case ViewServices:
		help = "[enter] select  [esc] back  [q] quit"
	case ViewHistory:
		help = "[enter] view  [esc] back  [q] quit"
	case ViewHistoryDetail:
		help = "[c] copy curl  [esc] back  [q] quit"
	}
	return helpStyle.Render(help)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
