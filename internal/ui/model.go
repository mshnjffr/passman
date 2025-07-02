package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mshnjffr/passman/internal/utils"
)

// NewModel creates and returns the initial menu model
func NewModel() tea.Model {
	return NewMenuModel(nil)
}

// NewModelWithManager creates and returns the initial menu model with manager
func NewModelWithManager(manager *utils.Manager) tea.Model {
	return NewMenuModel(manager)
}
