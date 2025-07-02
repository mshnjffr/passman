package generator

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// RandomGenerator generates cryptographically secure random passwords
type RandomGenerator struct {
	config Config
}

// NewRandomGenerator creates a new random password generator
func NewRandomGenerator(length int, charSets ...CharSet) *RandomGenerator {
	if len(charSets) == 0 {
		charSets = []CharSet{Lowercase, Uppercase, Numbers}
	}
	
	return &RandomGenerator{
		config: Config{
			Length:   length,
			CharSets: charSets,
		},
	}
}

// Generate creates a cryptographically secure random password
func (r *RandomGenerator) Generate(ctx context.Context) (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	// Build individual charsets for each enabled character type
	charsets := r.buildIndividualCharsets()
	if len(charsets) == 0 {
		return "", errors.New("no valid character sets")
	}

	// If password length is less than number of character sets, 
	// we can't guarantee all types are included
	if r.config.Length < len(charsets) {
		return "", errors.New("password length must be at least equal to number of enabled character types")
	}

	password := make([]byte, r.config.Length)

	// First, ensure at least one character from each enabled character set
	for i, charset := range charsets {
		select {
		case <-ctx.Done():
			clearBytes(password[:i])
			return "", ctx.Err()
		default:
		}

		if len(charset) == 0 {
			continue
		}

		charsetSize := big.NewInt(int64(len(charset)))
		randomIndex, err := rand.Int(rand.Reader, charsetSize)
		if err != nil {
			clearBytes(password[:i])
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		
		password[i] = charset[randomIndex.Int64()]
	}

	// Fill the remaining positions with random characters from all charsets
	fullCharset := r.buildCharset()
	if len(fullCharset) == 0 {
		clearBytes(password)
		return "", errors.New("no valid characters in charset")
	}

	fullCharsetSize := big.NewInt(int64(len(fullCharset)))
	
	for i := len(charsets); i < r.config.Length; i++ {
		select {
		case <-ctx.Done():
			clearBytes(password[:i])
			return "", ctx.Err()
		default:
		}

		randomIndex, err := rand.Int(rand.Reader, fullCharsetSize)
		if err != nil {
			clearBytes(password[:i])
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		
		password[i] = fullCharset[randomIndex.Int64()]
	}

	// Shuffle the password to randomize the positions
	err := r.shufflePassword(password)
	if err != nil {
		clearBytes(password)
		return "", fmt.Errorf("failed to shuffle password: %w", err)
	}

	result := string(password)
	clearBytes(password) // Clear sensitive data from memory
	
	return result, nil
}

// EstimateEntropy calculates the theoretical entropy for random passwords
func (r *RandomGenerator) EstimateEntropy() float64 {
	charset := r.buildCharset()
	if len(charset) == 0 {
		return 0
	}
	
	return float64(r.config.Length) * logBase2(float64(len(charset)))
}

// GetName returns the generator name
func (r *RandomGenerator) GetName() string {
	return "Random Password"
}

// Validate checks if the configuration is valid
func (r *RandomGenerator) Validate() error {
	if r.config.Length <= 0 {
		return errors.New("password length must be positive")
	}
	
	if r.config.Length > 1024 {
		return errors.New("password length too long (max 1024)")
	}
	
	if len(r.config.CharSets) == 0 {
		return errors.New("at least one character set must be specified")
	}
	
	return nil
}

// SetExcludeChars sets characters to exclude from generation
func (r *RandomGenerator) SetExcludeChars(chars string) {
	r.config.ExcludeChar = chars
}

// buildIndividualCharsets builds separate charsets for each enabled character type
func (r *RandomGenerator) buildIndividualCharsets() []string {
	var charsets []string
	
	for _, cs := range r.config.CharSets {
		var charset string
		switch cs {
		case Lowercase:
			charset = "abcdefghijklmnopqrstuvwxyz"
		case Uppercase:
			charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		case Numbers:
			charset = "0123456789"
		case Symbols:
			charset = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		case Ambiguous:
			charset = "0O1lI"
		}
		
		// Remove excluded characters
		if r.config.ExcludeChar != "" {
			charset = removeChars(charset, r.config.ExcludeChar)
		}
		
		if len(charset) > 0 {
			charsets = append(charsets, charset)
		}
	}
	
	return charsets
}

// buildCharset constructs the character set based on configuration
func (r *RandomGenerator) buildCharset() string {
	var charset strings.Builder
	
	for _, cs := range r.config.CharSets {
		switch cs {
		case Lowercase:
			charset.WriteString("abcdefghijklmnopqrstuvwxyz")
		case Uppercase:
			charset.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		case Numbers:
			charset.WriteString("0123456789")
		case Symbols:
			charset.WriteString("!@#$%^&*()_+-=[]{}|;:,.<>?")
		case Ambiguous:
			charset.WriteString("0O1lI")
		}
	}
	
	result := charset.String()
	
	// Remove excluded characters
	if r.config.ExcludeChar != "" {
		result = removeChars(result, r.config.ExcludeChar)
	}
	
	return result
}

// shufflePassword securely shuffles the password bytes using Fisher-Yates algorithm
func (r *RandomGenerator) shufflePassword(password []byte) error {
	n := len(password)
	for i := n - 1; i > 0; i-- {
		// Generate a random index from 0 to i
		maxIndex := big.NewInt(int64(i + 1))
		randomIndex, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return fmt.Errorf("failed to generate random index for shuffle: %w", err)
		}
		
		j := randomIndex.Int64()
		// Swap elements at positions i and j
		password[i], password[j] = password[j], password[i]
	}
	return nil
}

// removeChars removes specified characters from the charset
func removeChars(charset, exclude string) string {
	excludeMap := make(map[rune]bool)
	for _, char := range exclude {
		excludeMap[char] = true
	}
	
	var result strings.Builder
	for _, char := range charset {
		if !excludeMap[char] {
			result.WriteRune(char)
		}
	}
	
	return result.String()
}

// clearBytes securely clears sensitive data from memory
func clearBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
