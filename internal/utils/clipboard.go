package utils

import (
	"errors"
	"fmt"

	"github.com/atotto/clipboard"
)

// ClipboardManager handles cross-platform clipboard operations
type ClipboardManager struct{}

// NewClipboardManager creates a new clipboard manager instance
func NewClipboardManager() *ClipboardManager {
	return &ClipboardManager{}
}

// Copy copies the given text to the system clipboard
func (c *ClipboardManager) Copy(text string) error {
	if text == "" {
		return errors.New("cannot copy empty text to clipboard")
	}

	err := clipboard.WriteAll(text)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return nil
}

// Paste retrieves text from the system clipboard
func (c *ClipboardManager) Paste() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read from clipboard: %w", err)
	}

	return text, nil
}

// IsAvailable checks if clipboard functionality is available
func (c *ClipboardManager) IsAvailable() bool {
	// Try to read from clipboard to test availability
	_, err := clipboard.ReadAll()
	return err == nil
}

// Clear clears the clipboard (platform-dependent)
func (c *ClipboardManager) Clear() error {
	return clipboard.WriteAll("")
}
