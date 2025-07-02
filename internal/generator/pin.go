package generator

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// PINGenerator generates numeric PIN codes
type PINGenerator struct {
	config Config
}

// NewPINGenerator creates a new PIN generator
func NewPINGenerator(length int) *PINGenerator {
	return &PINGenerator{
		config: Config{
			Length: length,
		},
	}
}

// Generate creates a cryptographically secure numeric PIN
func (p *PINGenerator) Generate(ctx context.Context) (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}

	pin := make([]byte, p.config.Length)
	ten := big.NewInt(10)

	for i := 0; i < p.config.Length; i++ {
		select {
		case <-ctx.Done():
			// Clear sensitive data before returning
			clearBytes(pin[:i])
			return "", ctx.Err()
		default:
		}

		randomDigit, err := rand.Int(rand.Reader, ten)
		if err != nil {
			clearBytes(pin[:i])
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		
		pin[i] = byte('0' + randomDigit.Int64())
	}

	result := string(pin)
	clearBytes(pin) // Clear sensitive data from memory
	
	return result, nil
}

// GenerateFormatted creates a PIN with optional formatting (e.g., "1234-5678")
func (p *PINGenerator) GenerateFormatted(ctx context.Context, separator string, groupSize int) (string, error) {
	pin, err := p.Generate(ctx)
	if err != nil {
		return "", err
	}
	
	if separator == "" || groupSize <= 0 || groupSize >= len(pin) {
		return pin, nil
	}
	
	var formatted strings.Builder
	for i, digit := range pin {
		if i > 0 && i%groupSize == 0 {
			formatted.WriteString(separator)
		}
		formatted.WriteRune(digit)
	}
	
	return formatted.String(), nil
}

// EstimateEntropy calculates the theoretical entropy for numeric PINs
func (p *PINGenerator) EstimateEntropy() float64 {
	// Each digit has 10 possible values (0-9)
	return float64(p.config.Length) * logBase2(10.0)
}

// GetName returns the generator name
func (p *PINGenerator) GetName() string {
	return "Numeric PIN"
}

// Validate checks if the configuration is valid
func (p *PINGenerator) Validate() error {
	if p.config.Length <= 0 {
		return errors.New("PIN length must be positive")
	}
	
	if p.config.Length > 50 {
		return errors.New("PIN length too long (max 50)")
	}
	
	return nil
}

// SetLength sets the PIN length
func (p *PINGenerator) SetLength(length int) {
	p.config.Length = length
}
