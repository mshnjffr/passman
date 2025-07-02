package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mshnjffr/passman/internal/generator"
	"github.com/mshnjffr/passman/internal/utils"
)

// GeneratorModel represents the password generation screen
type GeneratorModel struct {
	generatorType   string
	lengthInput     textinput.Model
	wordCountInput  textinput.Model
	spinner         spinner.Model
	generating      bool
	currentPassword string
	strength        string
	statusMsg       string
	width           int
	height          int

	// Settings
	includeLower    bool
	includeUpper    bool
	includeNumbers  bool
	includeSymbols  bool
	
	// Manager for history and other utilities
	manager         *utils.Manager
}

type generateMsg struct {
	password string
	strength string
}

// NewGeneratorModel creates a new generator model
func NewGeneratorModel(genType string, manager *utils.Manager) *GeneratorModel {
	lengthInput := textinput.New()
	if genType == "pin" {
		pinLength := "4"
		if manager != nil {
			pinLength = fmt.Sprintf("%d", manager.Config.DefaultPinLength)
		}
		lengthInput.Placeholder = pinLength
		lengthInput.SetValue(pinLength)
	} else {
		lengthInput.Placeholder = "16"
		lengthInput.SetValue("16")
	}
	lengthInput.CharLimit = 3
	lengthInput.Width = 10
	// Don't focus by default so character toggles work immediately

	wordCountInput := textinput.New()
	wordCountInput.Placeholder = "4"
	wordCountInput.SetValue("4")
	wordCountInput.CharLimit = 2
	wordCountInput.Width = 10

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF10F0"))

	return &GeneratorModel{
		generatorType:   genType,
		lengthInput:     lengthInput,
		wordCountInput:  wordCountInput,
		spinner:         s,
		includeLower:    true,
		includeUpper:    true,
		includeNumbers:  true,
		includeSymbols:  true,
		statusMsg:       "",
		manager:         manager,
	}
}

func (m *GeneratorModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *GeneratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
		case "enter", "g":
			if !m.generating {
				m.generating = true
				m.statusMsg = "Generating password..."
				return m, tea.Batch(m.generatePassword(), m.spinner.Tick)
			}
		case "c":
			if m.currentPassword != "" && !strings.HasPrefix(m.currentPassword, "Error:") {
				// Try to copy to clipboard using the manager
				if m.manager != nil && m.manager.Clipboard != nil {
					if err := m.manager.Clipboard.Copy(m.currentPassword); err != nil {
						m.statusMsg = "Failed to copy to clipboard: " + err.Error()
					} else {
						m.statusMsg = "Password copied to clipboard!"
					}
				} else {
					m.statusMsg = "Clipboard not available"
				}
			} else if m.currentPassword == "" {
				m.statusMsg = "No password to copy. Generate one first!"
			} else {
				m.statusMsg = "Cannot copy error message to clipboard"
			}
		case "tab":
			// Toggle focus between inputs based on generator type
			if m.generatorType == "memorable" {
				// For memorable passphrase, toggle word count input focus
				if m.wordCountInput.Focused() {
					m.wordCountInput.Blur()
				} else {
					m.wordCountInput.Focus()
				}
			} else {
				// For random/pin, toggle length input focus
				if m.lengthInput.Focused() {
					m.lengthInput.Blur()
				} else {
					m.lengthInput.Focus()
				}
			}
		case "n":
			// Only toggle if input is not focused
			if !m.lengthInput.Focused() && !(m.generatorType == "memorable" && m.wordCountInput.Focused()) {
				m.includeNumbers = !m.includeNumbers
			}
		case "s":
			// Only toggle if input is not focused
			if !m.lengthInput.Focused() && !(m.generatorType == "memorable" && m.wordCountInput.Focused()) {
				m.includeSymbols = !m.includeSymbols
			}
		case "l":
			// Only toggle if input is not focused
			if !m.lengthInput.Focused() && !(m.generatorType == "memorable" && m.wordCountInput.Focused()) {
				m.includeLower = !m.includeLower
			}
		case "u":
			// Only toggle if input is not focused
			if !m.lengthInput.Focused() && !(m.generatorType == "memorable" && m.wordCountInput.Focused()) {
				m.includeUpper = !m.includeUpper
			}
		}

	case generateMsg:
		m.generating = false
		m.currentPassword = msg.password
		m.strength = msg.strength
		m.statusMsg = "Password generated successfully!"
		
		// Save to history if manager is available and password is valid
		if m.manager != nil && m.manager.History != nil && m.manager.History.IsEnabled() && msg.password != "" && !strings.HasPrefix(msg.password, "Error:") {
			settings := m.buildSettingsString()
			entry := utils.HistoryEntry{
				Password:    msg.password,
				Length:      len(msg.password),
				Type:        m.generatorType,
				Settings:    settings,
				Description: fmt.Sprintf("%s password", strings.Title(m.generatorType)),
			}
			if err := m.manager.History.AddEntry(entry); err != nil {
				// Don't fail the UI if history fails, just log it
				m.statusMsg = "Password generated successfully! (History save failed)"
			}
		}

	case spinner.TickMsg:
		if m.generating {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update inputs
	var cmd tea.Cmd
	m.lengthInput, cmd = m.lengthInput.Update(msg)
	cmds = append(cmds, cmd)

	if m.generatorType == "memorable" {
		m.wordCountInput, cmd = m.wordCountInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *GeneratorModel) generatePassword() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var gen generator.Generator
		var password string
		var err error

		switch m.generatorType {
		case "random":
			length, _ := strconv.Atoi(m.lengthInput.Value())
			if length <= 0 {
				length = 16
			}

			var charSets []generator.CharSet
			if m.includeLower {
				charSets = append(charSets, generator.Lowercase)
			}
			if m.includeUpper {
				charSets = append(charSets, generator.Uppercase)
			}
			if m.includeNumbers {
				charSets = append(charSets, generator.Numbers)
			}
			if m.includeSymbols {
				charSets = append(charSets, generator.Symbols)
			}

			gen = generator.NewRandomGenerator(length, charSets...)
			password, err = gen.Generate(ctx)

		case "memorable":
			wordCount, _ := strconv.Atoi(m.wordCountInput.Value())
			if wordCount <= 0 {
				wordCount = 4
			}
			gen = generator.NewMemorableGenerator(wordCount, " ", generator.GetEFFWordlist())
			password, err = gen.Generate(ctx)

		case "pin":
			length, _ := strconv.Atoi(m.lengthInput.Value())
			if length <= 0 {
				length = m.manager.Config.DefaultPinLength
			}
			gen = generator.NewPINGenerator(length)
			password, err = gen.Generate(ctx)
		}

		if err != nil {
			return generateMsg{password: "Error: " + err.Error(), strength: "Error"}
		}

		// Calculate strength
		strength := "Strong"
		if len(password) < 8 {
			strength = "Weak"
		} else if len(password) < 12 {
			strength = "Medium"
		}

		return generateMsg{password: password, strength: strength}
	}
}

func (m *GeneratorModel) View() string {
	var title string
	switch m.generatorType {
	case "random":
		title = "🔐 Generate Random Password"
	case "memorable":
		title = "🧠 Generate Memorable Passphrase"
	case "pin":
		title = "🔢 Generate PIN Code"
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Align(lipgloss.Center)

	// Settings section with white text
	var settings string
	if m.generatorType == "random" {
		var focusHint string
		if m.lengthInput.Focused() {
			focusHint = " (Press Tab to toggle character types)"
		} else {
			focusHint = " (Press Tab to edit length)"
		}
		
		settingsContent := fmt.Sprintf(`Settings:
Length: %s%s

Character Types:
%s
%s
%s
%s`,
			m.lengthInput.View(),
			focusHint,
			checkbox("Lowercase (l)", m.includeLower),
			checkbox("Uppercase (u)", m.includeUpper),
			checkbox("Numbers (n)", m.includeNumbers),
			checkbox("Symbols (s)", m.includeSymbols))
		settings = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(settingsContent)
	} else if m.generatorType == "memorable" {
		var focusHint string
		if m.wordCountInput.Focused() {
			focusHint = " (Press Tab to exit editing)"
		} else {
			focusHint = " (Press Tab to edit word count)"
		}
		
		settingsContent := fmt.Sprintf(`Settings:
Word Count: %s%s`, m.wordCountInput.View(), focusHint)
		settings = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(settingsContent)
	} else if m.generatorType == "pin" {
		settingsContent := fmt.Sprintf(`Settings:
PIN Length: %s`, m.lengthInput.View())
		settings = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(settingsContent)
	}

	// Password output
	var passwordDisplay string
	if m.generating {
		passwordDisplay = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Render(fmt.Sprintf("%s Generating...", m.spinner.View()))
	} else if m.currentPassword != "" {
		passwordDisplay = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Render(m.currentPassword)
		// Only show strength if enabled in settings
		if m.strength != "" && m.manager != nil && m.manager.Config != nil && m.manager.Config.ShowStrengthMeter {
			passwordDisplay += "\nStrength: " + lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Render(m.strength)
		}
	} else {
		passwordDisplay = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("Press Enter to generate a password")
	}

	// Helper commands at bottom like main menu
	help := subtleStyle.Render("enter/g: generate") + dotStyle +
		subtleStyle.Render("tab: toggle focus") + dotStyle +
		subtleStyle.Render("l/u/n/s: toggle types") + dotStyle +
		subtleStyle.Render("c: copy") + dotStyle +
		subtleStyle.Render("esc: back")

	// Status
	status := ""
	if m.statusMsg != "" {
		status = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Render(m.statusMsg)
	}

	// Calculate responsive box sizes based on terminal width
	var settingsWidth, passwordWidth int
	
	if m.width < 40 {
		// Very small terminals - minimal padding
		settingsWidth = m.width - 2
		passwordWidth = m.width - 2
	} else if m.width < 60 {
		// Small terminals - compact layout  
		settingsWidth = m.width - 4
		passwordWidth = m.width - 4
	} else if m.width < 80 {
		// Medium sized terminals - vertical layout
		availableWidth := m.width - 6
		settingsWidth = availableWidth - 2
		passwordWidth = availableWidth - 2
	} else {
		// Large terminals - horizontal layout
		availableWidth := m.width - 6
		settingsWidth = int(float64(availableWidth) * 0.45)
		passwordWidth = int(float64(availableWidth) * 0.50)
	}
	
	// Adjust height based on terminal height
	passwordHeight := 3
	if m.height < 20 {
		passwordHeight = 2 // Smaller height for medium terminals
	}
	if m.height < 15 {
		passwordHeight = 1 // Very small height for small terminals
	}
	
	// Adjust styling based on terminal size
	var settingsBoxStyle, passwordBoxStyle lipgloss.Style
	if m.width < 40 {
		// Minimal styling for very small terminals
		settingsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("15")).
			Padding(0, 1).
			Width(settingsWidth)
		passwordBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("15")).
			Padding(0, 1).
			Width(passwordWidth).
			Height(passwordHeight).
			Align(lipgloss.Center, lipgloss.Center)
	} else {
		// Normal styling for larger terminals
		settingsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("15")).
			Padding(1, 2).
			Width(settingsWidth)
		passwordBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("15")).
			Padding(1, 2).
			Width(passwordWidth).
			Height(passwordHeight).
			Align(lipgloss.Center, lipgloss.Center)
	}

	settingsBox := settingsBoxStyle.Render(settings)
	passwordBox := passwordBoxStyle.Render(passwordDisplay)

	// Combine boxes - use vertical layout for small terminals
	var mainContent string
	if m.width < 80 { // Use vertical layout for most terminals
		// Vertical layout for small and medium terminals
		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			settingsBox,
			"",
			passwordBox,
		)
	} else {
		// Horizontal layout for very large terminals
		mainContent = lipgloss.JoinHorizontal(
			lipgloss.Top,
			settingsBox,
			" ",
			passwordBox,
		)
	}

	// Combine everything like main menu - always reserve space for status
	var contentParts []string
	contentParts = append(contentParts, titleStyle.Render(title))
	
	// Responsive spacing between sections
	if m.height < 15 {
		// Very compact spacing for small terminals
		contentParts = append(contentParts, mainContent)
		// Always add status line to prevent layout shift
		if status != "" {
			contentParts = append(contentParts, status)
		} else {
			contentParts = append(contentParts, " ") // Empty space to maintain layout
		}
		contentParts = append(contentParts, help)
	} else if m.height < 20 {
		// Compact spacing for small terminals
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, mainContent)
		// Always add status line to prevent layout shift
		if status != "" {
			contentParts = append(contentParts, status)
		} else {
			contentParts = append(contentParts, " ") // Empty space to maintain layout
		}
		contentParts = append(contentParts, help)
	} else {
		// Reduced spacing for larger terminals
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, mainContent)
		contentParts = append(contentParts, "")
		// Always add status line to prevent layout shift
		if status != "" {
			contentParts = append(contentParts, status)
		} else {
			contentParts = append(contentParts, " ") // Empty space to maintain layout
		}
		contentParts = append(contentParts, help)
	}

	content := strings.Join(contentParts, "\n")

	// Apply main style with responsive spacing
	topSpacing := "\n\n"
	bottomSpacing := "\n"
	
	// Reduce spacing for small terminals
	if m.height < 15 {
		topSpacing = ""
		bottomSpacing = ""
	} else if m.height < 20 {
		topSpacing = "\n"
		bottomSpacing = ""
	}
	
	return mainStyle.Render(topSpacing + content + bottomSpacing)
}

// buildSettingsString creates a string representation of current settings
func (m *GeneratorModel) buildSettingsString() string {
	if m.generatorType == "random" {
		return fmt.Sprintf("Length: %s, Lower: %t, Upper: %t, Numbers: %t, Symbols: %t",
			m.lengthInput.Value(), m.includeLower, m.includeUpper, m.includeNumbers, m.includeSymbols)
	} else if m.generatorType == "memorable" {
		return fmt.Sprintf("Word Count: %s", m.wordCountInput.Value())
	} else if m.generatorType == "pin" {
		return fmt.Sprintf("PIN Length: %s", m.lengthInput.Value())
	}
	return ""
}


