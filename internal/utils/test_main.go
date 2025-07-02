// +build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mshnjffr/passman/internal/config"
	"github.com/mshnjffr/passman/internal/utils"
)

func main() {
	fmt.Println("Testing Password Generator TUI Utility Systems")
	fmt.Println("==============================================")

	// Test configuration
	fmt.Println("\n1. Testing Configuration System...")
	cfg := config.Default()
	if err := cfg.Validate(); err != nil {
		log.Printf("Config validation error: %v", err)
	} else {
		fmt.Println("✓ Configuration system working")
	}

	// Test clipboard
	fmt.Println("\n2. Testing Clipboard System...")
	clipboard := utils.NewClipboardManager()
	if clipboard.IsAvailable() {
		if err := clipboard.Copy("test-password-123"); err != nil {
			log.Printf("Clipboard copy error: %v", err)
		} else {
			fmt.Println("✓ Clipboard system working")
		}
	} else {
		fmt.Println("⚠ Clipboard not available (expected in headless environments)")
	}

	// Test export
	fmt.Println("\n3. Testing Export System...")
	export := utils.NewExportManager()
	tempFile := "/tmp/test_export.txt"
	
	err := export.ExportSingle(
		"test-password-123",
		"Test password for export",
		utils.FormatText,
		tempFile,
	)
	
	if err != nil {
		log.Printf("Export error: %v", err)
	} else {
		fmt.Println("✓ Export system working")
		// Clean up
		os.Remove(tempFile)
	}

	// Test wordlist
	fmt.Println("\n4. Testing Wordlist System...")
	wordlist := utils.NewWordlistManager()
	if err := wordlist.LoadWordlist(); err != nil {
		log.Printf("Wordlist load error: %v", err)
	} else {
		passphrase, err := wordlist.GeneratePassphrase(3, "-", false)
		if err != nil {
			log.Printf("Passphrase generation error: %v", err)
		} else {
			fmt.Printf("✓ Wordlist system working - generated: %s\n", passphrase)
		}
	}

	// Test history (with temporary passphrase)
	fmt.Println("\n5. Testing History System...")
	history := utils.NewHistoryManager(true, "test-passphrase-123", 10)
	
	testEntry := utils.HistoryEntry{
		Password:    "test-password-456",
		Length:      16,
		Type:        "test",
		Settings:    "test-settings",
		Description: "Test entry for validation",
	}
	
	if err := history.AddEntry(testEntry); err != nil {
		log.Printf("History add error: %v", err)
	} else {
		entries, err := history.LoadHistory()
		if err != nil {
			log.Printf("History load error: %v", err)
		} else {
			fmt.Printf("✓ History system working - %d entries stored\n", len(entries))
		}
	}

	// Test manager integration
	fmt.Println("\n6. Testing Manager Integration...")
	manager, err := utils.NewManager(&cfg)
	if err != nil {
		log.Printf("Manager creation error: %v", err)
	} else {
		info := manager.GetSystemInfo()
		fmt.Println("✓ Manager system working")
		fmt.Println("System Information:")
		for key, value := range info {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println("\n==============================================")
	fmt.Println("Utility Systems Test Complete")
}
