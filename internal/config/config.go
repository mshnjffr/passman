package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	// Password Generation Defaults
	DefaultLength            int  `json:"default_length"`
	DefaultIncludeLowercase  bool `json:"default_include_lowercase"`
	DefaultIncludeUppercase  bool `json:"default_include_uppercase"`
	DefaultIncludeNumbers    bool `json:"default_include_numbers"`
	DefaultIncludeSymbols    bool `json:"default_include_symbols"`
	DefaultExcludeSimilar    bool `json:"default_exclude_similar"`
	DefaultExcludeAmbiguous  bool `json:"default_exclude_ambiguous"`
	
	// Passphrase Defaults
	DefaultPassphraseWords      int    `json:"default_passphrase_words"`
	DefaultPassphraseSeparator  string `json:"default_passphrase_separator"`
	DefaultPassphraseCapitalize bool   `json:"default_passphrase_capitalize"`
	
	// PIN Defaults
	DefaultPinLength            int    `json:"default_pin_length"`
	
	// Clipboard Settings
	AutoCopyToClipboard    bool `json:"auto_copy_to_clipboard"`
	ClearClipboardAfter    int  `json:"clear_clipboard_after_seconds"` // 0 = never
	ShowClipboardSuccess   bool `json:"show_clipboard_success"`
	
	// Export Settings
	DefaultExportFormat    string `json:"default_export_format"`
	DefaultExportPath      string `json:"default_export_path"`
	IncludeTimestampInName bool   `json:"include_timestamp_in_name"`
	
	// History Settings
	HistoryEnabled         bool   `json:"history_enabled"`
	HistoryMaxEntries      int    `json:"history_max_entries"`
	HistoryEncryptionKey   string `json:"history_encryption_key,omitempty"` // Empty = prompt for passphrase
	
	// UI Settings
	Theme                  string `json:"theme"`
	ShowStrengthMeter      bool   `json:"show_strength_meter"`
	ShowGenerationTime     bool   `json:"show_generation_time"`
	ConfirmBeforeExit      bool   `json:"confirm_before_exit"`
	
	// Advanced Settings
	WordlistUpdateInterval int    `json:"wordlist_update_interval_days"`
	EnableTelemetry        bool   `json:"enable_telemetry"`
	Debug                  bool   `json:"debug"`
}

func Default() Config {
	homeDir, _ := os.UserHomeDir()
	defaultExportPath := filepath.Join(homeDir, "Documents", "passwords")
	
	return Config{
		// Password Generation Defaults
		DefaultLength:            12,
		DefaultIncludeLowercase:  true,
		DefaultIncludeUppercase:  true,
		DefaultIncludeNumbers:    true,
		DefaultIncludeSymbols:    false,
		DefaultExcludeSimilar:    false,
		DefaultExcludeAmbiguous:  false,
		
		// Passphrase Defaults
		DefaultPassphraseWords:      4,
		DefaultPassphraseSeparator:  "-",
		DefaultPassphraseCapitalize: false,
		
		// PIN Defaults
		DefaultPinLength:            4,
		
		// Clipboard Settings
		AutoCopyToClipboard:    true,
		ClearClipboardAfter:    0, // Never clear automatically
		ShowClipboardSuccess:   true,
		
		// Export Settings
		DefaultExportFormat:    "txt",
		DefaultExportPath:      defaultExportPath,
		IncludeTimestampInName: true,
		
		// History Settings
		HistoryEnabled:         true, // Enable by default with encryption
		HistoryMaxEntries:      100,
		HistoryEncryptionKey:   "default-key", // Default encryption key
		
		// UI Settings
		Theme:                  "default",
		ShowStrengthMeter:      true,
		ShowGenerationTime:     false,
		ConfirmBeforeExit:      false,
		
		// Advanced Settings
		WordlistUpdateInterval: 30, // 30 days
		EnableTelemetry:        false,
		Debug:                  false,
	}
}

func Load() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Default(), err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Default(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Default(), err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Default(), err
	}

	// Ensure missing fields have default values
	config = mergeWithDefaults(config)

	return config, nil
}

// mergeWithDefaults ensures missing fields have default values
func mergeWithDefaults(config Config) Config {
	defaults := Default()
	
	// Only set defaults for empty/zero values that should have defaults
	if config.HistoryEncryptionKey == "" {
		config.HistoryEncryptionKey = defaults.HistoryEncryptionKey
	}
	
	// Add other fields that might need default merging in the future
	if config.DefaultExportPath == "" {
		config.DefaultExportPath = defaults.DefaultExportPath
	}
	
	return config
}

func (c Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "passman", "config.json"), nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "passman"), nil
}

// Validate validates the configuration settings
func (c *Config) Validate() error {
	if c.DefaultLength < 1 || c.DefaultLength > 512 {
		c.DefaultLength = 12
	}
	
	if c.DefaultPassphraseWords < 1 || c.DefaultPassphraseWords > 20 {
		c.DefaultPassphraseWords = 4
	}
	
	if c.DefaultPinLength < 1 || c.DefaultPinLength > 50 {
		c.DefaultPinLength = 4
	}
	
	if c.DefaultPassphraseSeparator == "" {
		c.DefaultPassphraseSeparator = "-"
	}
	
	if c.ClearClipboardAfter < 0 {
		c.ClearClipboardAfter = 0
	}
	
	if c.HistoryMaxEntries < 1 {
		c.HistoryMaxEntries = 100
	} else if c.HistoryMaxEntries > 10000 {
		c.HistoryMaxEntries = 10000
	}
	
	validFormats := map[string]bool{"txt": true, "json": true, "csv": true}
	if !validFormats[c.DefaultExportFormat] {
		c.DefaultExportFormat = "txt"
	}
	
	if c.WordlistUpdateInterval < 1 {
		c.WordlistUpdateInterval = 30
	}
	
	return nil
}

// Reset resets the configuration to default values
func (c *Config) Reset() {
	*c = Default()
}

// IsHistoryEnabled returns true if history is enabled and properly configured
func (c *Config) IsHistoryEnabled() bool {
	return c.HistoryEnabled
}

// GetExportPath returns the full export path for a given filename
func (c *Config) GetExportPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(c.DefaultExportPath, filename)
}
