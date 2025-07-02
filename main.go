package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"passman/internal/config"
	"passman/internal/generator"
	"passman/internal/ui"
	"passman/internal/utils"
)

const (
	appName    = "passman"
	appVersion = "1.0.0"
)

func main() {
	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help", "-h", "help":
			showHelp()
			return
		case "--version", "-v", "version":
			fmt.Printf("%s %s\n", appName, appVersion)
			return
		case "--test", "test":
			runComponentTests()
			return
		case "--reset", "reset":
			resetConfiguration()
			return
		}
	}

	// Initialize logging
	initLogging()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		fmt.Fprintf(os.Stderr, "Error: Failed to load configuration: %v\n", err)
		return
	}

	// Initialize the utilities manager
	manager, err := utils.NewManager(&cfg)
	if err != nil {
		log.Printf("Failed to initialize manager: %v", err)
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize utilities: %v\n", err)
		return
	}

	// Initialize the UI with manager
	model := ui.NewModelWithManager(manager)

	// Create and run the Bubble Tea program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	log.Println("Application shutdown gracefully")
}

func showHelp() {
	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.json")

	fmt.Printf(`%s %s
A beautiful, secure password generator with a stunning terminal UI

USAGE:
  ./%s [options]

OPTIONS:
  -help, -h        Show this help information
  -version, -v     Show version information
  -test            Test system components and exit
  -reset           Reset configuration to defaults

FEATURES:
  🔐 Cryptographically secure password generation
  🎨 Beautiful neon-themed terminal interface
  📊 Real-time strength visualization
  🔄 Animated generation with spinners
  📋 Instant clipboard integration
  💾 Export to multiple formats
  📈 Advanced security analysis

KEYBOARD SHORTCUTS:
  Tab/Shift+Tab    Navigate between components
  g                Generate password
  c                Copy to clipboard
  s                Save/Export
  q, Ctrl+C        Quit

CONFIGURATION:
  Config directory: %s
  Config file: %s

EXAMPLES:
  ./%s              Start the beautiful TUI
  ./%s --test       Test system components
  ./%s --reset      Reset configuration

For more information, visit: https://github.com/mshnjffr/passman
`, appName, appVersion, appName, configDir, configFile, appName, appName, appName)
}

func runComponentTests() {
	fmt.Println("Testing system components...\n")

	// Test config loading
	fmt.Print("config:      ")
	if _, err := config.Load(); err != nil {
		fmt.Printf("✗ FAIL: %v\n", err)
	} else {
		fmt.Println("✓ PASS")
	}

	// Test generators
	fmt.Print("generators:  ")
	randomGen := generator.NewRandomGenerator(12, generator.Lowercase, generator.Uppercase)
	memorableGen := generator.NewMemorableGenerator(3, "-", generator.GetEFFWordlist())
	pinGen := generator.NewPINGenerator(4)
	
	ctx := context.Background()
	if _, err := randomGen.Generate(ctx); err != nil {
		fmt.Printf("✗ FAIL: random generator: %v\n", err)
	} else if _, err := memorableGen.Generate(ctx); err != nil {
		fmt.Printf("✗ FAIL: memorable generator: %v\n", err)
	} else if _, err := pinGen.Generate(ctx); err != nil {
		fmt.Printf("✗ FAIL: PIN generator: %v\n", err)
	} else {
		fmt.Println("✓ PASS")
	}

	// Test utilities
	fmt.Print("utilities:   ")
	cfg, _ := config.Load()
	if _, err := utils.NewManager(&cfg); err != nil {
		fmt.Printf("✗ FAIL: %v\n", err)
	} else {
		fmt.Println("✓ PASS")
	}

	fmt.Println("\nAll components tested successfully! 🎉")
}

func resetConfiguration() {
	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.json")
	
	if err := os.Remove(configFile); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error removing config file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Configuration reset to defaults.\nConfig file: %s\n", configFile)
}

func initLogging() {
	configDir := getConfigDir()
	logDir := filepath.Join(configDir, "logs")
	
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// Fallback to stderr if we can't create log directory
		return
	}
	
	logFile := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stderr if we can't open log file
		return
	}
	
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Application started - %s %s", appName, appVersion)
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".passman"
	}
	return filepath.Join(homeDir, ".config", appName)
}
