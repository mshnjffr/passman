package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mshnjffr/passman/internal/utils"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type clearStatusMsg struct{}

// HistoryModel represents the password history screen
type HistoryModel struct {
	table       table.Model
	manager     *utils.Manager
	width       int
	height      int
	statusMsg   string
	filterType  string // "all", "random", "memorable", "pin"
	allEntries  []utils.HistoryEntry // Cache all entries
	displayedEntries []utils.HistoryEntry // Currently displayed entries for copying
}

// NewHistoryModel creates a new history model
func NewHistoryModel(manager *utils.Manager) *HistoryModel {
	// Start with minimal default columns
	columns := []table.Column{
		{Title: "Time", Width: 8},
		{Title: "Password", Width: 15},
		{Title: "Length", Width: 4},
		{Title: "Type", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false).
		Foreground(lipgloss.Color("15"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.Foreground(lipgloss.Color("15"))
	t.SetStyles(s)

	model := &HistoryModel{
		table:      t,
		manager:    manager,
		width:      40,  // Conservative default for small terminals
		height:     12,  // Conservative default for small terminals
		filterType: "all", // Show all types by default
	}
	
	return model
}

// NewHistoryModelWithSize creates a new history model with specified dimensions
func NewHistoryModelWithSize(manager *utils.Manager, width, height int) *HistoryModel {
	model := NewHistoryModel(manager)
	model.width = width
	model.height = height
	model.updateTableSize()
	return model
}

// RefreshCache clears the cached entries to force a reload
func (m *HistoryModel) RefreshCache() {
	m.allEntries = nil
}

func (m *HistoryModel) Init() tea.Cmd {
	return nil
}

func (m *HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateTableSize()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return NewMenuModelWithSize(m.manager, m.width, m.height), nil
		case "esc":
			return NewMenuModelWithSize(m.manager, m.width, m.height), nil
		case "enter":
			// Copy selected password to clipboard (full password, not truncated)
			selectedIndex := m.table.Cursor()
			if selectedIndex >= 0 && selectedIndex < len(m.displayedEntries) && m.manager != nil && m.manager.Clipboard != nil {
				fullPassword := m.displayedEntries[selectedIndex].Password
				if err := m.manager.Clipboard.Copy(fullPassword); err == nil {
					m.statusMsg = "Password copied to clipboard!"
					return m, tea.Batch(cmd, m.clearStatusAfter(2*time.Second))
				} else {
					m.statusMsg = "Failed to copy to clipboard"
					return m, tea.Batch(cmd, m.clearStatusAfter(3*time.Second))
				}
			}
		case "a":
			// Show all types
			m.filterType = "all"
			m.statusMsg = "Showing all password types"
			return m, tea.Batch(cmd, m.clearStatusAfter(2*time.Second))
		case "r":
			// Filter by random passwords
			m.filterType = "random"
			m.statusMsg = "Filtering by Random passwords"
			return m, tea.Batch(cmd, m.clearStatusAfter(2*time.Second))
		case "m":
			// Filter by memorable passwords  
			m.filterType = "memorable"
			m.statusMsg = "Filtering by Memorable passwords"
			return m, tea.Batch(cmd, m.clearStatusAfter(2*time.Second))
		case "p":
			// Filter by PIN codes
			m.filterType = "pin"
			m.statusMsg = "Filtering by PIN codes"
			return m, tea.Batch(cmd, m.clearStatusAfter(2*time.Second))
		}
	case clearStatusMsg:
		m.statusMsg = ""
		return m, nil
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *HistoryModel) clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m *HistoryModel) updateTableSize() {
	// Adjust table size based on terminal dimensions
	tableWidth := m.width - 4  // Account for padding
	tableHeight := m.height - 8 // Account for title, help, and padding

	if tableHeight < 5 {
		tableHeight = 5
	}
	if tableHeight > 15 {
		tableHeight = 15
	}

	// Calculate responsive column widths
	var timeWidth, lengthWidth, typeWidth, passwordWidth int

	if m.width < 60 {
		// Very small terminals
		timeWidth = 8
		lengthWidth = 4
		typeWidth = 8
		passwordWidth = tableWidth - timeWidth - lengthWidth - typeWidth - 8
		if passwordWidth < 12 {
			passwordWidth = 12
		}
	} else if m.width < 100 {
		// Medium terminals
		timeWidth = 11
		lengthWidth = 6
		typeWidth = 10
		passwordWidth = tableWidth - timeWidth - lengthWidth - typeWidth - 8
		if passwordWidth < 20 {
			passwordWidth = 20
		}
	} else {
		// Large terminals
		timeWidth = 12
		lengthWidth = 8
		typeWidth = 12
		passwordWidth = tableWidth - timeWidth - lengthWidth - typeWidth - 8
		if passwordWidth < 30 {
			passwordWidth = 30
		}
	}

	columns := []table.Column{
		{Title: "Time", Width: timeWidth},
		{Title: "Password", Width: passwordWidth},
		{Title: "Length", Width: lengthWidth},
		{Title: "Type", Width: typeWidth},
	}

	m.table.SetColumns(columns)
	m.table.SetHeight(tableHeight)
}

func (m *HistoryModel) loadHistoryData() {
	if m.manager == nil || m.manager.History == nil || !m.manager.History.IsEnabled() {
		return
	}

	// Load all entries if not cached or refresh cache
	if len(m.allEntries) == 0 {
		entries, err := m.manager.History.LoadHistory() // Get ALL entries, not just recent
		if err != nil {
			return
		}
		m.allEntries = entries
	}

	// Filter entries based on current filter
	var filteredEntries []utils.HistoryEntry
	for _, entry := range m.allEntries {
		if m.filterType == "all" || strings.ToLower(entry.Type) == m.filterType {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	// Store displayed entries for copying (full passwords)
	m.displayedEntries = filteredEntries

	// Calculate password display width based on current column width
	passwordColumnWidth := 20 // Default fallback
	if len(m.table.Columns()) > 1 {
		passwordColumnWidth = m.table.Columns()[1].Width
	}

	// Convert to table rows
	var rows []table.Row
	for _, entry := range filteredEntries {
		timeStr := entry.CreatedAt.Format("Jan 2 15:04")
		
		// Handle password display based on available width
		password := entry.Password
		if passwordColumnWidth < 15 {
			// Very small width - show just first few chars
			if len(password) > 8 {
				password = password[:5] + "..."
			}
		} else if len(password) > passwordColumnWidth-3 {
			// Normal truncation for medium/large widths
			truncateAt := passwordColumnWidth - 6
			if truncateAt < 5 {
				truncateAt = 5
			}
			password = password[:truncateAt] + "..."
		}
		
		typeStr := strings.Title(entry.Type)
		lengthStr := strconv.Itoa(entry.Length)

		rows = append(rows, table.Row{
			timeStr,
			password,
			lengthStr,
			typeStr,
		})
	}

	m.table.SetRows(rows)
}

func (m *HistoryModel) View() string {
	// Load fresh data each time we render
	m.loadHistoryData()

	// Title with filter indicator
	titleText := "Password History"
	if m.filterType != "all" {
		titleText += " - " + strings.Title(m.filterType) + " Only"
	}
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Render(titleText)

	// Check if history is enabled and has data
	var content string
	if m.manager == nil || m.manager.History == nil || !m.manager.History.IsEnabled() {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Render("History is disabled.\n\nEnable it in settings to track your generated passwords.")
	} else {
		entries, _ := m.manager.History.GetRecentEntries(1)
		if len(entries) == 0 {
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Render("No passwords in history yet.\n\nGenerate some passwords to see them here!")
		} else {
			content = baseStyle.Render(m.table.View())
			
			// Add count information when filtering
			if m.filterType != "all" {
				filteredCount := len(m.table.Rows())
				totalCount := len(m.allEntries)
				countInfo := lipgloss.NewStyle().
					Foreground(lipgloss.Color("241")).
					Render(fmt.Sprintf("Showing %d of %d entries", filteredCount, totalCount))
				content += "\n" + countInfo
			}
		}
	}

	// Help text with filter shortcuts
	help := subtleStyle.Render("↑/↓: navigate") + dotStyle +
		subtleStyle.Render("enter: copy") + dotStyle +
		subtleStyle.Render("a/r/m/p: filter") + dotStyle +
		subtleStyle.Render("esc: back") + dotStyle +
		subtleStyle.Render("q: quit")

	// Status message
	status := ""
	if m.statusMsg != "" {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Render(m.statusMsg)
	}

	// Combine everything
	sections := []string{title, content}
	if status != "" {
		sections = append(sections, status)
	}
	sections = append(sections, help)
	fullContent := strings.Join(sections, "\n\n")

	// Apply main style with responsive spacing
	topSpacing := "\n\n"
	bottomSpacing := "\n"
	
	if m.height < 15 {
		topSpacing = ""
		bottomSpacing = ""
	} else if m.height < 20 {
		topSpacing = "\n"
		bottomSpacing = ""
	}

	return mainStyle.Render(topSpacing + fullContent + bottomSpacing)
}
