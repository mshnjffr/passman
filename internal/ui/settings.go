package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"passman/internal/utils"
)

// SettingsModel represents the settings screen
type SettingsModel struct {
	width    int
	height   int
	manager  *utils.Manager
	cursor   int
	settings []SettingItem
}

// SettingItem represents a configurable setting
type SettingItem struct {
	Name        string
	Description string
	Type        string // "toggle", "number", "text"
	Value       interface{}
	Key         string // Config key
}

// NewSettingsModel creates a new settings model
func NewSettingsModel(manager *utils.Manager) *SettingsModel {
	// Load current values from manager/config
	historyEnabled := false
	autoCopy := true
	defaultLength := 16
	showStrength := true
	
	if manager != nil {
		if manager.History != nil {
			historyEnabled = manager.History.IsEnabled()
		}
		if manager.Config != nil {
			autoCopy = manager.Config.AutoCopyToClipboard
			defaultLength = manager.Config.DefaultLength
			showStrength = manager.Config.ShowStrengthMeter
		}
	}
	
	settings := []SettingItem{
		{
			Name:        "Password History",
			Description: "Save generated passwords to encrypted history",
			Type:        "toggle",
			Value:       historyEnabled,
			Key:         "history_enabled",
		},
		{
			Name:        "Auto Copy to Clipboard",
			Description: "Automatically copy generated passwords",
			Type:        "toggle",
			Value:       autoCopy,
			Key:         "auto_copy_to_clipboard",
		},
		{
			Name:        "Default Password Length",
			Description: "Default length for random passwords",
			Type:        "number",
			Value:       defaultLength,
			Key:         "default_length",
		},
		{
			Name:        "Show Strength Meter",
			Description: "Display password strength analysis",
			Type:        "toggle",
			Value:       showStrength,
			Key:         "show_strength_meter",
		},
	}
	
	return &SettingsModel{
		manager:  manager,
		cursor:   0,
		settings: settings,
	}
}

func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return NewMenuModel(m.manager), nil
		case "esc":
			return NewMenuModel(m.manager), nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.settings)-1 {
				m.cursor++
			}
		case "enter", " ":
			// Toggle or modify the selected setting
			m.toggleSetting(m.cursor)
		}
	}

	return m, nil
}

func (m *SettingsModel) View() string {
	// Title with white text like main menu
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Render("Settings")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Render("Use ↑/↓ to navigate, Enter to change settings")

	// Build the settings list like main menu
	var settingsItems []string
	for i, setting := range m.settings {
		var valueStr string
		switch setting.Type {
		case "toggle":
			if val, ok := setting.Value.(bool); ok && val {
				valueStr = "Enabled"
			} else {
				valueStr = "Disabled"
			}
		case "number":
			valueStr = fmt.Sprintf("%v", setting.Value)
		default:
			valueStr = fmt.Sprintf("%v", setting.Value)
		}
		
		line := fmt.Sprintf("%s: %s", setting.Name, valueStr)
		settingsItems = append(settingsItems, checkbox(line, m.cursor == i))
	}

	settingsList := strings.Join(settingsItems, "\n")

	// Helper commands like main menu
	help := subtleStyle.Render("↑/↓: navigate") + dotStyle +
		subtleStyle.Render("enter: change") + dotStyle +
		subtleStyle.Render("esc: back") + dotStyle +
		subtleStyle.Render("q: quit")

	// Combine everything like main menu
	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
		title,
		subtitle,
		settingsList,
		help,
	)

	// Apply main style like the main menu
	return mainStyle.Render("\n" + content + "\n\n")
}

// toggleSetting handles toggling or modifying settings values
func (m *SettingsModel) toggleSetting(index int) {
	if index < 0 || index >= len(m.settings) {
		return
	}
	
	setting := &m.settings[index]
	var newValue interface{}
	
	switch setting.Type {
	case "toggle":
		if val, ok := setting.Value.(bool); ok {
			newValue = !val
			setting.Value = newValue
		}
	case "number":
		// For now, cycle through common values for password length
		if setting.Key == "default_length" {
			lengths := []int{8, 12, 16, 20, 24, 32}
			if val, ok := setting.Value.(int); ok {
				for i, length := range lengths {
					if length == val {
						newValue = lengths[(i+1)%len(lengths)]
						setting.Value = newValue
						break
					}
				}
			}
		}
	}
	
	// Apply the setting change to the manager/config
	m.applySetting(setting.Key, newValue)
}

// applySetting applies a setting change to the manager and config
func (m *SettingsModel) applySetting(key string, value interface{}) {
	if m.manager == nil || m.manager.Config == nil {
		return
	}
	
	// Update the config with the new value
	switch key {
	case "history_enabled":
		if val, ok := value.(bool); ok {
			m.manager.Config.HistoryEnabled = val
			// Also update the history manager
			if m.manager.History != nil {
				m.manager.History.SetEnabled(val)
				// Ensure passphrase is set when enabling history
				if val && m.manager.Config.HistoryEncryptionKey != "" {
					m.manager.History.SetPassphrase(m.manager.Config.HistoryEncryptionKey)
				} else if val {
					// Set a default passphrase if none exists
					m.manager.Config.HistoryEncryptionKey = "default-encryption-key"
					m.manager.History.SetPassphrase(m.manager.Config.HistoryEncryptionKey)
				}
			}
		}
	case "auto_copy_to_clipboard":
		if val, ok := value.(bool); ok {
			m.manager.Config.AutoCopyToClipboard = val
		}
	case "default_length":
		if val, ok := value.(int); ok {
			m.manager.Config.DefaultLength = val
		}
	case "show_strength_meter":
		if val, ok := value.(bool); ok {
			m.manager.Config.ShowStrengthMeter = val
		}
	}
	
	// Save the updated config to file
	if err := m.manager.Config.Save(); err != nil {
		// If save fails, we could show an error message in the future
		// For now, changes are still applied in memory
	}
}
