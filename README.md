# Passman 🔐

A **beautiful**, secure password manager with an elegant terminal UI built with Go and Bubble Tea. Generate, store, and manage your passwords with style through an intuitive menu-driven interface.

![Passman](logo.png)

## ✨ **Features**

### 🎯 **Clean Menu-Driven Design** 
- **Simple main menu** - easy to navigate options
- **Focused screens** - one task at a time for clarity
- **Intuitive navigation** - arrow keys and enter to select
- **Clean user experience** - no overwhelming interfaces

### 🌈 **Stunning Visual Components**
- **🔄 Animated Spinners** - Neon pink loading animations during generation
- **📊 Progress Bars** - Real-time strength visualization with gradient colors  
- **📝 Modern Text Inputs** - Sleek input fields with neon focus states
- **📋 Beautiful Tables** - Organized password history with neon borders
- **📜 Scrollable Viewports** - Detailed history with full information
- **🎯 Selection Lists** - Generator type selection with beautiful highlighting

### 🎨 **Neon Theme Design**
- **Vibrant color palette** - Pink, Blue, Green, Yellow, Purple, Cyan
- **Rounded borders** and modern typography
- **Gradient effects** on titles and important elements
- **Adaptive theming** for light/dark terminals
- **Smooth focus transitions** with bright color indicators

### 🔐 **Advanced Security Features**
- **Cryptographically secure** random generation using `crypto/rand`
- **High-quality passwords** - no patterns or repetition (e.g., no "iiiiiiiiiiqqqqq")
- **Dynamic configuration** - settings instantly applied to generation
- **Memory safe** with automatic cleanup of sensitive data
- **Real-time entropy calculation** and strength scoring
- **Pattern detection** (sequences, repetition, keyboard patterns)
- **Crack time estimation** based on current hardware
- **No data collection** - everything stays local

### 🚀 **Password Generation Modes**
- **🔐 Random Passwords**: Strong random passwords with customizable character sets
  - ✅ **Fixed Quality Issue**: No more repeated characters or patterns
  - ✅ **Respects Settings**: Length, numbers, symbols, letters all properly applied
  - ✅ **High Entropy**: 16-char passwords achieve ~103 bits of entropy
- **🧠 Memorable Passphrases**: EFF wordlist-based for easy recall (~46 bits entropy)
- **🔢 Numeric PINs**: Secure PIN codes with customizable length 
- **⚡ Live Configuration**: Settings instantly applied to password generation

### 💎 **Enhanced User Experience**
- **Instant clipboard integration** with visual confirmation
- **Real-time strength meters** using animated progress bars
- **Tabbed navigation** - seamlessly move between all components
- **Keyboard shortcuts** for power users
- **Visual feedback** for every action and state change

## Installation

### 🍺 Homebrew (macOS/Linux)

```bash
# Add the tap
brew tap mshnjffr/passman

# Install passman
brew install passman

# Run from anywhere
passman
```

```bash
# Update to latest version
brew upgrade passman
```

> **Advantages**: No Go installation required, automatic dependency management, easy updates

### 🔧 Go Install (Recommended for Go users)

```bash
# Install the latest version directly from GitHub
go install github.com/mshnjffr/passman@latest

# Run from anywhere
passman
```

```bash
# Update to latest version (same command)
go install github.com/mshnjffr/passman@latest
```

> **Advantages**: Always latest version, development builds available, faster installation

> **Note**: Ensure `$GOPATH/bin` (usually `~/go/bin`) is in your PATH. If `passman` command is not found after installation, add Go's bin directory to your PATH:
> 
> **For Zsh (most common on macOS):**
> ```bash
> # Open your zsh configuration file
> nano ~/.zshrc
> 
> # Add this line at the end of the file
> export PATH="$PATH:$(go env GOPATH)/bin"
> 
> # Save and reload your shell configuration
> source ~/.zshrc
> ```
>
> **For Bash:**
> ```bash
> # Open your bash configuration file
> nano ~/.bashrc   # or ~/.bash_profile on macOS
> 
> # Add this line at the end of the file
> export PATH="$PATH:$(go env GOPATH)/bin"
> 
> # Save and reload your shell configuration
> source ~/.bashrc   # or source ~/.bash_profile
> ```
>
> **Quick one-liner for Zsh:**
> ```bash
> echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc
> ```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mshnjffr/passman.git
cd passman

# Build the application
go build -o passman

# Run the application
./passman
```

### Alternative Installation Methods

```bash
# Install specific version
go install github.com/mshnjffr/passman@v1.0.1

# Install from main branch (development)
go install github.com/mshnjffr/passman@main
```

### Dependencies

- Go 1.21 or later (required for installation)
- Terminal with Unicode support
- Git (for building from source)
- Optional: `xclip` (Linux) or `pbcopy` (macOS) for clipboard support

## Usage

### Basic Usage

```bash
# Start the interactive TUI
passman

# Show version and help
passman --help
passman --version

# Test system components
passman --test

# Reset configuration to defaults
passman --reset

# Enable debug logging
passman --debug
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate menu options |
| `Enter` | Select menu item |
| `g` | Generate password (in generator screens) |
| `c` | Copy to clipboard |
| `l/u/n/s` | Toggle character types (lowercase/uppercase/numbers/symbols) |
| `Esc` | Back to main menu |
| `q`, `Ctrl+C` | Quit |

### Generation Modes

#### Random Passwords
- Customizable length (1-128 characters)
- Character set selection (lowercase, uppercase, numbers, symbols)
- Exclude similar characters (0/O, 1/l/I)
- Custom character sets

#### Memorable Passphrases
- EFF wordlist-based generation
- Customizable word count (2-12 words)
- Multiple separator options
- Capitalization control
- Word filtering

#### Numeric PINs
- Customizable length (4-20 digits)
- Optional formatting with separators
- Pattern exclusion (no repeating sequences)

## Configuration

Configuration is stored at `~/.config/passman/config.json`:

```json
{
  "defaults": {
    "random": {
      "length": 16,
      "include_lowercase": true,
      "include_uppercase": true,
      "include_numbers": true,
      "include_symbols": true,
      "exclude_similar": false
    },
    "memorable": {
      "word_count": 4,
      "separator": "-",
      "capitalize": false,
      "include_numbers": false
    },
    "pin": {
      "length": 6,
      "include_separators": false,
      "separator": "-",
      "group_size": 3
    }
  },
  "ui": {
    "theme": "default",
    "show_help": true,
    "confirm_quit": true
  },
  "clipboard": {
    "enabled": true,
    "auto_clear": false,
    "clear_delay": 30
  },
  "history": {
    "enabled": false,
    "max_entries": 100,
    "encrypt": true
  }
}
```

## Architecture

The application features a clean, component-focused architecture:

```
├── main.go                    # Clean application entry point
├── internal/
│   ├── generator/            # Password generation engines
│   │   ├── interface.go      # Common generator interface
│   │   ├── random.go         # Random password generator
│   │   ├── memorable.go      # Memorable passphrase generator
│   │   ├── pin.go           # PIN generator
│   │   ├── analyzer.go      # Security analysis
│   │   └── utils.go         # Helper functions
│   ├── ui/                  # Beautiful single-screen UI
│   │   ├── model.go         # Main UI model with all components
│   │   ├── view.go          # View rendering
│   │   ├── styles.go        # Neon theme styling
│   │   └── generators.go    # Generator configurations
│   ├── config/              # Configuration management
│   │   └── config.go        # Config loading/saving
│   └── utils/               # Utilities and helpers
│       ├── clipboard.go     # Clipboard operations
│       ├── export.go        # File export
│       ├── wordlist.go      # EFF wordlist management
│       └── history.go       # Password history
├── go.mod
└── README.md
```

### Key Architecture Benefits:
- **Single-screen design** - all functionality visible at once
- **Component integration** - all Bubble Tea components work together seamlessly
- **Clean state management** - simplified, maintainable code
- **Beautiful visual design** - focus on user experience and security

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Development build
go build -o passman

# Release build
go build -ldflags="-s -w" -o passman
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Security

This application takes security seriously:

- **No network access** - everything runs locally
- **Cryptographically secure random generation** using OS entropy
- **Memory safety** with automatic cleanup of sensitive data
- **Optional encryption** for stored data
- **No telemetry or data collection**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [motus](https://github.com/oleiade/motus) by oleiade
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- Uses [EFF Wordlists](https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases) for memorable passwords
- Security recommendations based on [NIST SP 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html)

## Support

If you encounter any issues or have feature requests, please open an issue on GitHub.

---

**Happy password generating! 🔐**
