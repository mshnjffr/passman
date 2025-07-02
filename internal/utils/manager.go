package utils

import (
	"fmt"
	"os"
	"passman/internal/config"
)

// Manager centralizes access to all utility systems
type Manager struct {
	Config    *config.Config
	Clipboard *ClipboardManager
	Export    *ExportManager
	Wordlist  *WordlistManager
	History   *HistoryManager
}

// NewManager creates a new utilities manager with initialized components
func NewManager(cfg *config.Config) (*Manager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize components
	clipboard := NewClipboardManager()
	export := NewExportManager()
	wordlist := NewWordlistManager()
	
	// Initialize history manager with encryption if enabled
	var history *HistoryManager
	if cfg.HistoryEnabled {
		history = NewHistoryManager(
			cfg.HistoryEnabled,
			cfg.HistoryEncryptionKey,
			cfg.HistoryMaxEntries,
		)
	} else {
		history = NewHistoryManager(false, "", 0)
	}

	manager := &Manager{
		Config:    cfg,
		Clipboard: clipboard,
		Export:    export,
		Wordlist:  wordlist,
		History:   history,
	}

	// Load wordlist if needed
	if err := manager.initializeWordlist(); err != nil {
		// Don't fail initialization if wordlist loading fails
		// This allows the app to work even if wordlist is unavailable
		if cfg.Debug {
			fmt.Printf("Warning: Failed to load wordlist: %v\n", err)
		}
	}

	return manager, nil
}

// initializeWordlist loads the wordlist for passphrase generation
func (m *Manager) initializeWordlist() error {
	return m.Wordlist.LoadWordlist()
}

// UpdateConfig updates the manager's configuration and reinitializes components if needed
func (m *Manager) UpdateConfig(newConfig *config.Config) error {
	if newConfig == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	oldConfig := m.Config
	m.Config = newConfig

	// Reinitialize history if settings changed
	if oldConfig.HistoryEnabled != newConfig.HistoryEnabled ||
		oldConfig.HistoryMaxEntries != newConfig.HistoryMaxEntries ||
		oldConfig.HistoryEncryptionKey != newConfig.HistoryEncryptionKey {
		
		m.History = NewHistoryManager(
			newConfig.HistoryEnabled,
			newConfig.HistoryEncryptionKey,
			newConfig.HistoryMaxEntries,
		)
	}

	return nil
}

// GetSystemInfo returns information about the utility systems
func (m *Manager) GetSystemInfo() map[string]interface{} {
	info := map[string]interface{}{
		"clipboard_available": m.Clipboard.IsAvailable(),
		"wordlist_loaded":     m.Wordlist.IsLoaded(),
		"wordlist_source":     m.Wordlist.GetLoadedFrom(),
		"wordlist_word_count": m.Wordlist.GetWordCount(),
		"history_enabled":     m.History.IsEnabled(),
		"config_valid":        m.Config != nil,
	}

	if m.History.IsEnabled() {
		if count, err := m.History.GetEntryCount(); err == nil {
			info["history_entry_count"] = count
		}
	}

	return info
}

// TestSystems performs basic tests on all utility systems
func (m *Manager) TestSystems() map[string]error {
	results := make(map[string]error)

	// Test clipboard
	if m.Clipboard.IsAvailable() {
		if err := m.Clipboard.Copy("test"); err != nil {
			results["clipboard"] = err
		} else {
			results["clipboard"] = nil
		}
	} else {
		results["clipboard"] = fmt.Errorf("clipboard not available")
	}

	// Test wordlist
	if m.Wordlist.IsLoaded() {
		if _, err := m.Wordlist.GeneratePassphrase(2, "-", false); err != nil {
			results["wordlist"] = err
		} else {
			results["wordlist"] = nil
		}
	} else {
		results["wordlist"] = fmt.Errorf("wordlist not loaded")
	}

	// Test export
	tempPath := "/tmp/test_export.txt"
	if err := m.Export.ExportSingle("test-password", "test", FormatText, tempPath); err != nil {
		results["export"] = err
	} else {
		results["export"] = nil
		// Clean up test file
		if err := os.Remove(tempPath); err != nil && m.Config.Debug {
			fmt.Printf("Warning: Failed to clean up test file: %v\n", err)
		}
	}

	// Test history (if enabled)
	if m.History.IsEnabled() {
		testEntry := HistoryEntry{
			Password:    "test-password",
			Length:      13,
			Type:        "test",
			Settings:    "test-settings",
			Description: "test entry",
		}
		
		if err := m.History.AddEntry(testEntry); err != nil {
			results["history"] = err
		} else {
			results["history"] = nil
		}
	} else {
		results["history"] = fmt.Errorf("history disabled")
	}

	return results
}

// Cleanup performs cleanup operations for all utility systems
func (m *Manager) Cleanup() error {
	var errors []error

	// Clear clipboard if auto-clear is enabled
	if m.Config.ClearClipboardAfter > 0 {
		if err := m.Clipboard.Clear(); err != nil {
			errors = append(errors, fmt.Errorf("failed to clear clipboard: %w", err))
		}
	}

	// Save current configuration
	if err := m.Config.Save(); err != nil {
		errors = append(errors, fmt.Errorf("failed to save configuration: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}
