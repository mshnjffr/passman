package generator

import (
	"fmt"
	"math"
)

// logBase2 calculates logarithm base 2
func logBase2(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Log2(x)
}

// clearString securely clears a string from memory (best effort)
func clearString(s *string) {
	if s == nil {
		return
	}
	// Note: In Go, strings are immutable, so we can't actually clear them
	// This is a placeholder for the pattern - in practice, use []byte for sensitive data
	*s = ""
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// contains checks if a slice contains a specific item
func contains(slice []CharSet, item CharSet) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// unique removes duplicate strings from a slice
func unique(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// CharSetToString converts a CharSet to its string representation
func CharSetToString(cs CharSet) string {
	switch cs {
	case Lowercase:
		return "lowercase"
	case Uppercase:
		return "uppercase"
	case Numbers:
		return "numbers"
	case Symbols:
		return "symbols"
	case Ambiguous:
		return "ambiguous"
	default:
		return "unknown"
	}
}

// SecurityLevelToString converts a SecurityLevel to its string representation
func SecurityLevelToString(level SecurityLevel) string {
	switch level {
	case VeryWeak:
		return "Very Weak"
	case Weak:
		return "Weak"
	case Fair:
		return "Fair"
	case Good:
		return "Good"
	case Strong:
		return "Strong"
	case VeryStrong:
		return "Very Strong"
	default:
		return "Unknown"
	}
}

// GetSecurityLevelColor returns a color code for the security level (for UI)
func GetSecurityLevelColor(level SecurityLevel) string {
	switch level {
	case VeryWeak:
		return "#ff4444" // Red
	case Weak:
		return "#ff8800" // Orange
	case Fair:
		return "#ffaa00" // Yellow-Orange
	case Good:
		return "#88cc00" // Yellow-Green
	case Strong:
		return "#44aa44" // Green
	case VeryStrong:
		return "#00aa88" // Teal
	default:
		return "#888888" // Gray
	}
}

// EstimateGenerationTime estimates how long it would take to generate passwords
func EstimateGenerationTime(count int, generatorType string) string {
	baseTimeNs := 1000 // Base time in nanoseconds
	
	switch generatorType {
	case "Random Password":
		baseTimeNs *= 10 // Random generation is relatively fast
	case "Memorable Passphrase":
		baseTimeNs *= 50 // Wordlist lookup takes more time
	case "Numeric PIN":
		baseTimeNs *= 5  // Numeric generation is fastest
	default:
		baseTimeNs *= 20
	}
	
	totalNs := int64(count * baseTimeNs)
	
	if totalNs < 1000000 { // Less than 1ms
		return "< 1ms"
	} else if totalNs < 1000000000 { // Less than 1s
		return "< 1s"
	} else {
		seconds := totalNs / 1000000000
		return fmt.Sprintf("%ds", seconds)
	}
}

// ValidatePasswordStrength validates if a password meets minimum requirements
func ValidatePasswordStrength(password string, minLength int, requireMixed bool) []string {
	var issues []string
	
	if len(password) < minLength {
		issues = append(issues, fmt.Sprintf("Password must be at least %d characters long", minLength))
	}
	
	if requireMixed {
		analyzer := NewSecurityAnalyzer()
		analysis := analyzer.Analyze(password)
		
		if !analysis.HasLowercase {
			issues = append(issues, "Password must contain lowercase letters")
		}
		if !analysis.HasUppercase {
			issues = append(issues, "Password must contain uppercase letters")
		}
		if !analysis.HasNumbers {
			issues = append(issues, "Password must contain numbers")
		}
		if !analysis.HasSymbols {
			issues = append(issues, "Password must contain symbols")
		}
	}
	
	return issues
}
