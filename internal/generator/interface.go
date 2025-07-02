package generator

import "context"

// Generator defines the common interface for all password generators
type Generator interface {
	// Generate creates a new password/passphrase based on the generator's configuration
	Generate(ctx context.Context) (string, error)
	
	// EstimateEntropy calculates the theoretical entropy for passwords generated
	// with the current configuration
	EstimateEntropy() float64
	
	// GetName returns a human-readable name for this generator
	GetName() string
	
	// Validate checks if the current configuration is valid
	Validate() error
}

// Config holds common configuration options
type Config struct {
	Length      int
	CharSets    []CharSet
	WordCount   int
	Separator   string
	ExcludeChar string
}

// CharSet represents different character types
type CharSet int

const (
	Lowercase CharSet = 1 << iota
	Uppercase
	Numbers
	Symbols
	Ambiguous // Characters that can be confused (0, O, l, 1, etc.)
)

// SecurityLevel represents password strength levels
type SecurityLevel int

const (
	VeryWeak SecurityLevel = iota
	Weak
	Fair
	Good
	Strong
	VeryStrong
)

// SecurityAnalysis contains detailed password security metrics
type SecurityAnalysis struct {
	Entropy       float64
	Level         SecurityLevel
	CrackTime     string
	Feedback      []string
	CharsetSize   int
	HasLowercase  bool
	HasUppercase  bool
	HasNumbers    bool
	HasSymbols    bool
	HasAmbiguous  bool
	CommonWords   []string
	IsCompromised bool
}
