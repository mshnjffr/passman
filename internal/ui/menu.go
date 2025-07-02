package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mshnjffr/passman/internal/utils"
)

// Styling constants following the views example
var (
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotChar       = " • "
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

// Screen represents different app screens
type Screen int

const (
	MenuScreen Screen = iota
	GenerateScreen
	HistoryScreen
	SettingsScreen
)

// MenuModel represents the main menu state
type MenuModel struct {
	choices  []string
	actions  []string
	cursor   int
	quitting bool
	width    int
	height   int
	manager  *utils.Manager
}

// NewMenuModel creates a new menu model
func NewMenuModel(manager *utils.Manager) *MenuModel {
	choices := []string{
		"Generate Random Password",
		"Generate Memorable Passphrase",
		"Generate PIN Code",
		"View Password History",
		"Settings",
		"Quit",
	}

	actions := []string{
		"random",
		"memorable", 
		"pin",
		"history",
		"settings",
		"quit",
	}

	return &MenuModel{
		choices: choices,
		actions: actions,
		cursor:  0,
		manager: manager,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			action := m.actions[m.cursor]
			switch action {
			case "quit":
				m.quitting = true
				return m, tea.Quit
			case "random":
				return NewGeneratorModel("random", m.manager), nil
			case "memorable":
				return NewGeneratorModel("memorable", m.manager), nil
			case "pin":
				return NewGeneratorModel("pin", m.manager), nil
			case "history":
				return NewHistoryModel(m.manager), nil
			case "settings":
				return NewSettingsModel(m.manager), nil
			}
		}
	}

	return m, nil
}

func (m *MenuModel) View() string {
	if m.quitting {
		return "\n  Thanks for using Password Generator TUI! 👋\n\n"
	}

	// Title with white color
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Render("Password Generator TUI")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Render("What would you like to do today?")

	// Build the checkbox-style menu exactly like the views example
	var menuItems []string
	for i, choice := range m.choices {
		menuItems = append(menuItems, checkbox(choice, m.cursor == i))
	}

	menu := strings.Join(menuItems, "\n")

	// Footer with arrows and cleaner formatting like the help example
	help := subtleStyle.Render("↑/↓: navigate") + dotStyle +
		subtleStyle.Render("enter: select") + dotStyle +
		subtleStyle.Render("q: quit")

	// Combine everything
	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
		title,
		subtitle,
		menu,
		help,
	)

	// Apply main style (margin left) like the example
	return mainStyle.Render("\n" + content + "\n\n")
}

// checkbox renders a checkbox with label, exactly like the views example
func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}


