package utils

import (
	"bufio"
	"crypto/rand"
	"embed"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed data/eff_large_wordlist.txt
var embeddedWordlist embed.FS

// WordlistManager handles EFF wordlist operations
type WordlistManager struct {
	wordlist       []string
	loadedFromFile bool
}

// NewWordlistManager creates a new wordlist manager instance
func NewWordlistManager() *WordlistManager {
	return &WordlistManager{}
}

// LoadWordlist loads the EFF wordlist (embedded or from file)
func (w *WordlistManager) LoadWordlist() error {
	// Try to load from embedded first
	if err := w.loadEmbeddedWordlist(); err == nil {
		return nil
	}

	// Try to load from config directory
	configPath, err := w.getWordlistPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); err == nil {
		return w.loadFromFile(configPath)
	}

	// Download and cache the wordlist
	return w.downloadAndCacheWordlist()
}

// loadEmbeddedWordlist loads the wordlist from embedded files
func (w *WordlistManager) loadEmbeddedWordlist() error {
	file, err := embeddedWordlist.Open("data/eff_large_wordlist.txt")
	if err != nil {
		return fmt.Errorf("failed to open embedded wordlist: %w", err)
	}
	defer file.Close()

	return w.parseWordlist(file)
}

// loadFromFile loads the wordlist from a file
func (w *WordlistManager) loadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open wordlist file: %w", err)
	}
	defer file.Close()

	w.loadedFromFile = true
	return w.parseWordlist(file)
}

// parseWordlist parses the wordlist from a reader
func (w *WordlistManager) parseWordlist(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	words := make([]string, 0, 7776) // EFF large wordlist has 7776 words

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// EFF wordlist format: "11111	abacus"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			word := parts[1]
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read wordlist: %w", err)
	}

	if len(words) == 0 {
		return fmt.Errorf("wordlist is empty or invalid")
	}

	w.wordlist = words
	return nil
}

// downloadAndCacheWordlist downloads the EFF wordlist and caches it
func (w *WordlistManager) downloadAndCacheWordlist() error {
	const effWordlistURL = "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt"

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(effWordlistURL)
	if err != nil {
		return fmt.Errorf("failed to download wordlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download wordlist: HTTP %d", resp.StatusCode)
	}

	// Parse wordlist from response
	if err := w.parseWordlist(resp.Body); err != nil {
		return err
	}

	// Cache to file
	if err := w.cacheWordlist(); err != nil {
		// Don't fail if we can't cache, just log it
		// In a real application, you'd use proper logging
		return nil
	}

	return nil
}

// cacheWordlist saves the current wordlist to cache
func (w *WordlistManager) cacheWordlist() error {
	cachePath, err := w.getWordlistPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for i, word := range w.wordlist {
		// Write in EFF format
		fmt.Fprintf(file, "%05d\t%s\n", i+1, word)
	}

	return nil
}

// getWordlistPath returns the path for cached wordlist
func (w *WordlistManager) getWordlistPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "passman", "eff_wordlist.txt"), nil
}

// GeneratePassphrase generates a memorable passphrase using EFF wordlist
func (w *WordlistManager) GeneratePassphrase(numWords int, separator string, capitalize bool) (string, error) {
	if len(w.wordlist) == 0 {
		if err := w.LoadWordlist(); err != nil {
			return "", fmt.Errorf("failed to load wordlist: %w", err)
		}
	}

	if numWords <= 0 {
		numWords = 4
	}

	if separator == "" {
		separator = "-"
	}

	words := make([]string, numWords)
	for i := 0; i < numWords; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(w.wordlist))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}

		word := w.wordlist[index.Int64()]
		if capitalize {
			word = strings.Title(word)
		}
		words[i] = word
	}

	return strings.Join(words, separator), nil
}

// GetWordCount returns the number of words in the loaded wordlist
func (w *WordlistManager) GetWordCount() int {
	return len(w.wordlist)
}

// IsLoaded returns true if wordlist is loaded
func (w *WordlistManager) IsLoaded() bool {
	return len(w.wordlist) > 0
}

// GetLoadedFrom returns information about wordlist source
func (w *WordlistManager) GetLoadedFrom() string {
	if !w.IsLoaded() {
		return "not loaded"
	}
	if w.loadedFromFile {
		return "cached file"
	}
	return "embedded"
}
