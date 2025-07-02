package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// HistoryEntry represents a password generation history entry
type HistoryEntry struct {
	ID          string    `json:"id"`
	Password    string    `json:"password"`
	Length      int       `json:"length"`
	Type        string    `json:"type"`
	Settings    string    `json:"settings"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description,omitempty"`
}

// HistoryManager handles encrypted password history
type HistoryManager struct {
	enabled    bool
	passphrase string
	maxEntries int
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(enabled bool, passphrase string, maxEntries int) *HistoryManager {
	if maxEntries <= 0 {
		maxEntries = 100
	}

	return &HistoryManager{
		enabled:    enabled,
		passphrase: passphrase,
		maxEntries: maxEntries,
	}
}

// AddEntry adds a new entry to the history
func (h *HistoryManager) AddEntry(entry HistoryEntry) error {
	if !h.enabled {
		return fmt.Errorf("history is disabled")
	}

	if h.passphrase == "" {
		return fmt.Errorf("history passphrase not set")
	}

	entries, err := h.LoadHistory()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load existing history: %w", err)
	}

	// Generate ID if not provided
	if entry.ID == "" {
		entry.ID = h.generateID()
	}

	// Set creation time if not provided
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	// Add new entry at the beginning
	entries = append([]HistoryEntry{entry}, entries...)

	// Trim to max entries
	if len(entries) > h.maxEntries {
		entries = entries[:h.maxEntries]
	}

	return h.saveHistory(entries)
}

// LoadHistory loads and decrypts the history
func (h *HistoryManager) LoadHistory() ([]HistoryEntry, error) {
	if !h.enabled {
		return nil, fmt.Errorf("history is disabled")
	}

	if h.passphrase == "" {
		return nil, fmt.Errorf("history passphrase not set")
	}

	historyPath, err := h.getHistoryPath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return []HistoryEntry{}, nil
	}

	// Read encrypted data
	encryptedData, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	// Decrypt data
	decryptedData, err := h.decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt history: %w", err)
	}

	// Parse JSON
	var entries []HistoryEntry
	if err := json.Unmarshal(decryptedData, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse history data: %w", err)
	}

	return entries, nil
}

// saveHistory encrypts and saves the history
func (h *HistoryManager) saveHistory(entries []HistoryEntry) error {
	historyPath, err := h.getHistoryPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(historyPath), 0700); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history data: %w", err)
	}

	// Encrypt data
	encryptedData, err := h.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt history data: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(historyPath, encryptedData, 0600); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// ClearHistory removes all history entries
func (h *HistoryManager) ClearHistory() error {
	if !h.enabled {
		return fmt.Errorf("history is disabled")
	}

	historyPath, err := h.getHistoryPath()
	if err != nil {
		return err
	}

	// Remove the file
	if err := os.Remove(historyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove history file: %w", err)
	}

	return nil
}

// GetRecentEntries returns the most recent entries
func (h *HistoryManager) GetRecentEntries(limit int) ([]HistoryEntry, error) {
	entries, err := h.LoadHistory()
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > len(entries) {
		limit = len(entries)
	}

	return entries[:limit], nil
}

// SearchEntries searches for entries matching criteria
func (h *HistoryManager) SearchEntries(query string) ([]HistoryEntry, error) {
	entries, err := h.LoadHistory()
	if err != nil {
		return nil, err
	}

	var matches []HistoryEntry
	for _, entry := range entries {
		if h.matchesQuery(entry, query) {
			matches = append(matches, entry)
		}
	}

	return matches, nil
}

// matchesQuery checks if an entry matches the search query
func (h *HistoryManager) matchesQuery(entry HistoryEntry, query string) bool {
	query = strings.ToLower(query)
	
	return strings.Contains(strings.ToLower(entry.Type), query) ||
		   strings.Contains(strings.ToLower(entry.Description), query) ||
		   strings.Contains(strings.ToLower(entry.Settings), query)
}

// getHistoryPath returns the path to the history file
func (h *HistoryManager) getHistoryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "passman", "history.enc"), nil
}

// generateID generates a unique ID for history entries
func (h *HistoryManager) generateID() string {
	randNum, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), randNum.Int64())
}

// encrypt encrypts data using AES-GCM with a key derived from passphrase
func (h *HistoryManager) encrypt(data []byte) ([]byte, error) {
	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive key from passphrase
	key := pbkdf2.Key([]byte(h.passphrase), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Combine salt + nonce + ciphertext
	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// decrypt decrypts data using AES-GCM
func (h *HistoryManager) decrypt(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < 16+12 { // salt + nonce minimum
		return nil, fmt.Errorf("encrypted data too short")
	}

	// Extract salt, nonce, and ciphertext
	salt := encryptedData[:16]
	nonce := encryptedData[16:28]
	ciphertext := encryptedData[28:]

	// Derive key from passphrase
	key := pbkdf2.Key([]byte(h.passphrase), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// IsEnabled returns whether history is enabled
func (h *HistoryManager) IsEnabled() bool {
	return h.enabled
}

// SetEnabled enables or disables history
func (h *HistoryManager) SetEnabled(enabled bool) {
	h.enabled = enabled
}

// SetPassphrase sets the encryption passphrase
func (h *HistoryManager) SetPassphrase(passphrase string) {
	h.passphrase = passphrase
}

// GetEntryCount returns the number of entries in history
func (h *HistoryManager) GetEntryCount() (int, error) {
	entries, err := h.LoadHistory()
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}
