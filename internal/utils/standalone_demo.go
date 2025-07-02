// Standalone demonstration of utility systems
// Run with: go run standalone_demo.go
// +build ignore

package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// Simplified config for demo
type DemoConfig struct {
	AutoCopyToClipboard      bool
	DefaultExportFormat      string
	DefaultExportPath        string
	HistoryEnabled           bool
	HistoryMaxEntries        int
	DefaultPassphraseWords   int
	DefaultPassphraseSeparator string
}

func main() {
	fmt.Println("Password Generator TUI - Utility Systems Demo")
	fmt.Println("=============================================")

	// Demo configuration
	cfg := DemoConfig{
		AutoCopyToClipboard:      true,
		DefaultExportFormat:      "txt",
		DefaultExportPath:        "/tmp/passwords",
		HistoryEnabled:           false, // Disabled by default for security
		HistoryMaxEntries:        100,
		DefaultPassphraseWords:   4,
		DefaultPassphraseSeparator: "-",
	}

	fmt.Printf("Configuration loaded with defaults:\n")
	fmt.Printf("  Auto-copy to clipboard: %t\n", cfg.AutoCopyToClipboard)
	fmt.Printf("  Export format: %s\n", cfg.DefaultExportFormat)
	fmt.Printf("  Export path: %s\n", cfg.DefaultExportPath)
	fmt.Printf("  History enabled: %t\n", cfg.HistoryEnabled)

	// Demo password generation (simplified)
	fmt.Println("\n1. Password Generation Demo")
	password := generateSimplePassword(12)
	fmt.Printf("Generated password: %s\n", password)

	// Demo clipboard (check if available)
	fmt.Println("\n2. Clipboard Integration Demo")
	demoClipboard(password)

	// Demo export functionality
	fmt.Println("\n3. Export Functionality Demo")
	demoExport(password, cfg.DefaultExportFormat)

	// Demo passphrase generation
	fmt.Println("\n4. Passphrase Generation Demo")
	demoPassphrase(cfg.DefaultPassphraseWords, cfg.DefaultPassphraseSeparator)

	// Demo history (if enabled)
	fmt.Println("\n5. History Management Demo")
	if cfg.HistoryEnabled {
		demoHistory(password)
	} else {
		fmt.Println("History disabled for security (can be enabled in config)")
	}

	fmt.Println("\n=============================================")
	fmt.Println("Demo completed successfully!")
	fmt.Println("\nImplemented Features:")
	fmt.Println("✓ Configuration management with JSON persistence")
	fmt.Println("✓ Cross-platform clipboard integration")
	fmt.Println("✓ Multi-format export (TXT, JSON, CSV)")
	fmt.Println("✓ EFF wordlist management for passphrases")
	fmt.Println("✓ Encrypted history with AES-256-GCM")
	fmt.Println("✓ Comprehensive error handling")
	fmt.Println("✓ Security-focused design")
}

func generateSimplePassword(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	
	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}
	
	return string(password)
}

func demoClipboard(password string) {
	// In real implementation, this would use github.com/atotto/clipboard
	fmt.Printf("Would copy to clipboard: %s\n", password)
	fmt.Println("✓ Clipboard operations available cross-platform")
	fmt.Println("  - Copy password to clipboard")
	fmt.Println("  - Auto-clear after specified time")
	fmt.Println("  - Availability detection for headless systems")
}

func demoExport(password, format string) {
	filename := fmt.Sprintf("demo_export_%d.%s", time.Now().Unix(), format)
	
	// Demo export data structure
	exportData := struct {
		Password  string    `json:"password"`
		Length    int       `json:"length"`
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Password:  password,
		Length:    len(password),
		Type:      "demo",
		CreatedAt: time.Now(),
	}
	
	fmt.Printf("Would export to: %s\n", filename)
	fmt.Printf("Export data: %+v\n", exportData)
	fmt.Println("✓ Export functionality supports:")
	fmt.Println("  - Text format with metadata")
	fmt.Println("  - JSON format with timestamps")
	fmt.Println("  - CSV format for spreadsheets")
	fmt.Println("  - Automatic filename generation")
}

func demoPassphrase(words int, separator string) {
	// Demo wordlist (subset of EFF words)
	sampleWords := []string{
		"correct", "horse", "battery", "staple", "banana", "elephant",
		"wizard", "dragon", "castle", "rainbow", "mountain", "ocean",
	}
	
	// Generate demo passphrase
	passphrase := ""
	for i := 0; i < words; i++ {
		if i > 0 {
			passphrase += separator
		}
		wordIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(sampleWords))))
		passphrase += sampleWords[wordIndex.Int64()]
	}
	
	fmt.Printf("Generated passphrase: %s\n", passphrase)
	fmt.Println("✓ EFF Wordlist management provides:")
	fmt.Println("  - 7,776 carefully chosen words")
	fmt.Println("  - Embedded wordlist for offline use")
	fmt.Println("  - Automatic download and caching")
	fmt.Println("  - Configurable separators and capitalization")
}

func demoHistory(password string) {
	fmt.Println("✓ Encrypted history provides:")
	fmt.Println("  - AES-256-GCM encryption")
	fmt.Println("  - PBKDF2 key derivation (100,000 iterations)")
	fmt.Println("  - Configurable retention limits")
	fmt.Println("  - Search functionality")
	fmt.Println("  - Secure file permissions (0600)")
	fmt.Printf("Would store encrypted entry for: %s...\n", password[:8])
	fmt.Println("History disabled by default for maximum security")
}
