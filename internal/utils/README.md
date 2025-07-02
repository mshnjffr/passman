# Password Generator TUI - Utility Systems

This package provides comprehensive utility systems for the Password Generator TUI application, including configuration management, clipboard operations, file export, wordlist management, and encrypted history.

## Components

### 1. Configuration Management (`config.go`)

Enhanced JSON-based configuration system with comprehensive settings for all password generation and application features.

**Features:**
- Password generation defaults (length, character sets, exclusions)
- Passphrase generation settings (word count, separator, capitalization)
- Clipboard integration settings (auto-copy, clear timer, notifications)
- Export preferences (format, path, filename patterns)
- History management (encryption, retention limits)
- UI customization (themes, meters, confirmations)
- Advanced settings (wordlist updates, telemetry, debug mode)

**Usage:**
```go
// Load configuration
cfg, err := config.Load()
if err != nil {
    cfg = config.Default()
}

// Validate and save
cfg.Validate()
cfg.Save()
```

### 2. Clipboard Management (`clipboard.go`)

Cross-platform clipboard operations with error handling and availability detection.

**Features:**
- Copy/paste text to/from system clipboard
- Clipboard availability detection
- Clear clipboard functionality
- Cross-platform compatibility (Windows, macOS, Linux)

**Usage:**
```go
clipboard := NewClipboardManager()

// Check availability
if clipboard.IsAvailable() {
    // Copy password
    err := clipboard.Copy("my-secure-password")
    
    // Retrieve from clipboard
    text, err := clipboard.Paste()
    
    // Clear clipboard
    err := clipboard.Clear()
}
```

### 3. File Export (`export.go`)

Multi-format password export with structured data and metadata.

**Supported Formats:**
- **Text (.txt)**: Human-readable format with metadata
- **JSON (.json)**: Structured data with export timestamp
- **CSV (.csv)**: Spreadsheet-compatible format

**Features:**
- Single password or batch export
- Automatic filename generation with timestamps
- Directory creation and path validation
- Metadata inclusion (creation time, type, description)

**Usage:**
```go
export := NewExportManager()

// Export single password
err := export.ExportSingle(
    "my-password", 
    "High security password", 
    FormatJSON, 
    "/path/to/export.json"
)

// Export multiple entries
entries := []PasswordEntry{...}
err := export.Export(entries, FormatCSV, "/path/to/passwords.csv")
```

### 4. EFF Wordlist Management (`wordlist.go`)

EFF Large Wordlist integration for memorable passphrase generation.

**Features:**
- Embedded wordlist (7,776 words) for offline use
- Automatic download and caching from EFF servers
- Configurable passphrase generation (word count, separators, capitalization)
- Wordlist validation and integrity checking

**Usage:**
```go
wordlist := NewWordlistManager()

// Load wordlist (embedded or cached)
err := wordlist.LoadWordlist()

// Generate passphrase
passphrase, err := wordlist.GeneratePassphrase(
    4,        // number of words
    "-",      // separator
    false,    // capitalize
)
// Result: "correct-horse-battery-staple"
```

### 5. Encrypted History (`history.go`)

Optional encrypted password generation history with AES-256-GCM encryption.

**Security Features:**
- AES-256-GCM encryption with PBKDF2 key derivation
- User-provided passphrase for encryption key
- Secure file permissions (0600)
- Configurable retention limits
- Optional functionality (disabled by default)

**Features:**
- Encrypted storage of password generation history
- Search functionality across entries
- Configurable maximum entries
- Secure deletion and cleanup

**Usage:**
```go
// Create history manager (disabled by default for security)
history := NewHistoryManager(true, "encryption-passphrase", 100)

// Add entry
entry := HistoryEntry{
    Password:    "generated-password",
    Type:        "standard",
    Settings:    "length:16,symbols:true",
    Description: "High security password",
}
err := history.AddEntry(entry)

// Load history
entries, err := history.LoadHistory()

// Search entries
matches, err := history.SearchEntries("high security")
```

### 6. Utilities Manager (`manager.go`)

Centralized management of all utility systems with configuration integration.

**Features:**
- Single point of access for all utilities
- Configuration-driven initialization
- System health checks and testing
- Coordinated cleanup operations

**Usage:**
```go
// Load configuration
cfg, err := config.Load()

// Create manager
manager, err := NewManager(&cfg)

// Access utilities
manager.Clipboard.Copy("password")
manager.Wordlist.GeneratePassphrase(4, "-", false)
manager.Export.ExportSingle("password", "desc", FormatJSON, "file.json")

// System information
info := manager.GetSystemInfo()

// Test all systems
results := manager.TestSystems()
```

## Configuration File Structure

The configuration file is stored at `~/.config/passman/config.json`:

```json
{
  "default_length": 12,
  "default_include_lowercase": true,
  "default_include_uppercase": true,
  "default_include_numbers": true,
  "default_include_symbols": false,
  "default_exclude_similar": false,
  "default_exclude_ambiguous": false,
  "default_passphrase_words": 4,
  "default_passphrase_separator": "-",
  "default_passphrase_capitalize": false,
  "auto_copy_to_clipboard": true,
  "clear_clipboard_after_seconds": 0,
  "show_clipboard_success": true,
  "default_export_format": "txt",
  "default_export_path": "~/Documents/passwords",
  "include_timestamp_in_name": true,
  "history_enabled": false,
  "history_max_entries": 100,
  "history_encryption_key": "",
  "theme": "default",
  "show_strength_meter": true,
  "show_generation_time": false,
  "confirm_before_exit": false,
  "wordlist_update_interval_days": 30,
  "enable_telemetry": false,
  "debug": false
}
```

## Security Considerations

### History Encryption
- **AES-256-GCM**: Industry-standard encryption with authenticated encryption
- **PBKDF2**: 100,000 iterations with SHA-256 for key derivation
- **Random salt and nonce**: Unique for each encryption operation
- **Secure file permissions**: 0600 (owner read/write only)

### Clipboard Security
- **Auto-clear**: Optional automatic clipboard clearing after specified time
- **Secure clearing**: Overwrites clipboard with empty string
- **Availability checking**: Prevents errors on headless systems

### Export Security
- **Directory permissions**: Creates directories with 0755 permissions
- **File validation**: Validates paths and formats before writing
- **No password obfuscation**: Exports contain plaintext passwords (by design)

## Dependencies

- `github.com/atotto/clipboard` - Cross-platform clipboard operations
- `golang.org/x/crypto/pbkdf2` - PBKDF2 key derivation
- Standard library: `crypto/aes`, `crypto/cipher`, `crypto/rand`, `crypto/sha256`

## Error Handling

All utility functions return descriptive errors with context:
- Configuration validation errors
- Clipboard operation failures
- File system errors (permissions, disk space)
- Encryption/decryption failures
- Network errors (wordlist download)

## Testing

Each utility component includes comprehensive test functions:
- `TestSystems()` - Tests all systems for basic functionality
- Individual component tests for specific features
- Configuration validation tests
- Error condition testing

## Examples

See `examples.go` for comprehensive usage examples of all utility systems.

## Cross-Platform Compatibility

The utility systems are designed for cross-platform operation:
- **Windows**: Clipboard, file operations, configuration paths
- **macOS**: Native clipboard integration, standard paths
- **Linux**: X11/Wayland clipboard support, XDG compliance

## Performance Considerations

- **Lazy loading**: Wordlist loaded only when needed
- **Caching**: EFF wordlist cached locally after first download
- **Memory efficiency**: Streaming operations for large exports
- **Encryption overhead**: Minimal impact on history operations
