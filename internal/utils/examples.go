package utils

import (
	"context"
	"fmt"
	"time"
	"github.com/mshnjffr/passman/internal/config"
	"github.com/mshnjffr/passman/internal/generator"
)

// Example demonstrates how to use the utility systems
func Example() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create utilities manager
	manager, err := NewManager(&cfg)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	// Generate a password using the new generator interface
	var charSets []generator.CharSet
	if cfg.DefaultIncludeLowercase {
		charSets = append(charSets, generator.Lowercase)
	}
	if cfg.DefaultIncludeUppercase {
		charSets = append(charSets, generator.Uppercase)
	}
	if cfg.DefaultIncludeNumbers {
		charSets = append(charSets, generator.Numbers)
	}
	if cfg.DefaultIncludeSymbols {
		charSets = append(charSets, generator.Symbols)
	}
	
	randomGen := generator.NewRandomGenerator(cfg.DefaultLength, charSets...)
	password, err := randomGen.Generate(context.Background())
	if err != nil {
		return fmt.Errorf("failed to generate password: %w", err)
	}

	fmt.Printf("Generated password: %s\n", password)

	// Copy to clipboard if enabled
	if cfg.AutoCopyToClipboard && manager.Clipboard.IsAvailable() {
		if err := manager.Clipboard.Copy(password); err != nil {
			fmt.Printf("Warning: Failed to copy to clipboard: %v\n", err)
		} else if cfg.ShowClipboardSuccess {
			fmt.Println("Password copied to clipboard!")
		}
	}

	// Add to history if enabled
	if cfg.HistoryEnabled {
		entry := HistoryEntry{
			Password:    password,
			Length:      len(password),
			Type:        "standard",
			Settings:    fmt.Sprintf("L:%d,U:%t,N:%t,S:%t", cfg.DefaultLength, cfg.DefaultIncludeUppercase, cfg.DefaultIncludeNumbers, cfg.DefaultIncludeSymbols),
			Description: "Generated via example",
		}

		if err := manager.History.AddEntry(entry); err != nil {
			fmt.Printf("Warning: Failed to add to history: %v\n", err)
		}
	}

	// Generate a passphrase
	passphrase, err := manager.Wordlist.GeneratePassphrase(
		cfg.DefaultPassphraseWords,
		cfg.DefaultPassphraseSeparator,
		cfg.DefaultPassphraseCapitalize,
	)
	if err != nil {
		fmt.Printf("Warning: Failed to generate passphrase: %v\n", err)
	} else {
		fmt.Printf("Generated passphrase: %s\n", passphrase)
	}

	// Export examples
	entries := []PasswordEntry{
		{
			Password:    password,
			Length:      len(password),
			Type:        "standard",
			CreatedAt:   time.Now(),
			Description: "Example standard password",
		},
	}

	if passphrase != "" {
		entries = append(entries, PasswordEntry{
			Password:    passphrase,
			Length:      len(passphrase),
			Type:        "passphrase",
			CreatedAt:   time.Now(),
			Description: "Example passphrase",
		})
	}

	// Export in different formats
	formats := []ExportFormat{FormatText, FormatJSON, FormatCSV}
	for _, format := range formats {
		filename := manager.Export.GetSuggestedFilename(format, "example")
		fullPath := cfg.GetExportPath(filename)
		
		if err := manager.Export.Export(entries, format, fullPath); err != nil {
			fmt.Printf("Warning: Failed to export as %s: %v\n", format, err)
		} else {
			fmt.Printf("Exported passwords to %s\n", fullPath)
		}
	}

	// Display system information
	fmt.Println("\nSystem Information:")
	info := manager.GetSystemInfo()
	for key, value := range info {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// Test all systems
	fmt.Println("\nSystem Tests:")
	tests := manager.TestSystems()
	for system, err := range tests {
		status := "OK"
		if err != nil {
			status = fmt.Sprintf("ERROR: %v", err)
		}
		fmt.Printf("  %s: %s\n", system, status)
	}

	// Cleanup
	if err := manager.Cleanup(); err != nil {
		fmt.Printf("Warning: Cleanup failed: %v\n", err)
	}

	return nil
}

// ExamplePassphraseGeneration demonstrates passphrase generation
func ExamplePassphraseGeneration() error {
	wordlist := NewWordlistManager()
	
	if err := wordlist.LoadWordlist(); err != nil {
		return fmt.Errorf("failed to load wordlist: %w", err)
	}

	fmt.Printf("Wordlist loaded: %d words from %s\n", 
		wordlist.GetWordCount(), wordlist.GetLoadedFrom())

	// Generate various passphrases
	examples := []struct {
		words      int
		separator  string
		capitalize bool
		description string
	}{
		{4, "-", false, "Standard passphrase"},
		{6, " ", true, "Long capitalized passphrase"},
		{3, ".", false, "Short dot-separated passphrase"},
		{5, "", false, "Concatenated passphrase"},
	}

	for _, example := range examples {
		passphrase, err := wordlist.GeneratePassphrase(
			example.words, 
			example.separator, 
			example.capitalize,
		)
		if err != nil {
			fmt.Printf("Failed to generate %s: %v\n", example.description, err)
			continue
		}
		fmt.Printf("%s: %s\n", example.description, passphrase)
	}

	return nil
}

// ExampleHistoryManagement demonstrates history functionality
func ExampleHistoryManagement() error {
	// Create history manager with encryption
	history := NewHistoryManager(true, "example-passphrase", 10)

	// Add some example entries
	entries := []HistoryEntry{
		{
			Password:    "ExamplePassword123!",
			Length:      19,
			Type:        "standard",
			Settings:    "Length:19,Upper:true,Numbers:true,Symbols:true",
			Description: "High security password",
		},
		{
			Password:    "correct-horse-battery-staple",
			Length:      28,
			Type:        "passphrase",
			Settings:    "Words:4,Separator:-,Capitalize:false",
			Description: "XKCD-style passphrase",
		},
	}

	for _, entry := range entries {
		if err := history.AddEntry(entry); err != nil {
			return fmt.Errorf("failed to add entry: %w", err)
		}
	}

	// Load and display history
	loadedEntries, err := history.LoadHistory()
	if err != nil {
		return fmt.Errorf("failed to load history: %w", err)
	}

	fmt.Printf("History contains %d entries:\n", len(loadedEntries))
	for i, entry := range loadedEntries {
		fmt.Printf("  %d. %s (%s) - %s\n", 
			i+1, entry.Type, entry.CreatedAt.Format("2006-01-02 15:04"), entry.Description)
	}

	// Search example
	matches, err := history.SearchEntries("passphrase")
	if err != nil {
		return fmt.Errorf("failed to search history: %w", err)
	}

	fmt.Printf("\nFound %d matches for 'passphrase':\n", len(matches))
	for _, match := range matches {
		fmt.Printf("  - %s: %s\n", match.Type, match.Description)
	}

	return nil
}

// ExampleClipboardOperations demonstrates clipboard functionality
func ExampleClipboardOperations() error {
	clipboard := NewClipboardManager()

	if !clipboard.IsAvailable() {
		return fmt.Errorf("clipboard not available")
	}

	// Test copy and paste
	testText := "example-password-123"
	if err := clipboard.Copy(testText); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	pastedText, err := clipboard.Paste()
	if err != nil {
		return fmt.Errorf("failed to paste: %w", err)
	}

	if pastedText != testText {
		return fmt.Errorf("clipboard mismatch: expected %q, got %q", testText, pastedText)
	}

	fmt.Printf("Clipboard test successful: copied and retrieved %q\n", testText)

	// Clear clipboard
	if err := clipboard.Clear(); err != nil {
		fmt.Printf("Warning: Failed to clear clipboard: %v\n", err)
	}

	return nil
}
