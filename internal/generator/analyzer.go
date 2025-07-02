package generator

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// SecurityAnalyzer analyzes password security and provides detailed metrics
type SecurityAnalyzer struct {
	commonPasswords []string
	commonWords     []string
}

// NewSecurityAnalyzer creates a new security analyzer
func NewSecurityAnalyzer() *SecurityAnalyzer {
	return &SecurityAnalyzer{
		commonPasswords: getCommonPasswords(),
		commonWords:     getCommonWords(),
	}
}

// Analyze performs comprehensive security analysis of a password
func (s *SecurityAnalyzer) Analyze(password string) SecurityAnalysis {
	analysis := SecurityAnalysis{
		Entropy:      s.calculateEntropy(password),
		CharsetSize:  s.calculateCharsetSize(password),
		HasLowercase: s.hasLowercase(password),
		HasUppercase: s.hasUppercase(password),
		HasNumbers:   s.hasNumbers(password),
		HasSymbols:   s.hasSymbols(password),
		HasAmbiguous: s.hasAmbiguous(password),
		CommonWords:  s.findCommonWords(password),
		Feedback:     []string{},
	}
	
	analysis.Level = s.calculateSecurityLevel(analysis.Entropy, len(password), password)
	analysis.CrackTime = s.estimateCrackTime(analysis.Entropy)
	analysis.IsCompromised = s.isCommonPassword(password)
	analysis.Feedback = s.generateFeedback(password, analysis)
	
	return analysis
}

// calculateEntropy estimates password entropy using multiple methods
func (s *SecurityAnalyzer) calculateEntropy(password string) float64 {
	if len(password) == 0 {
		return 0
	}
	
	charsetSize := s.calculateCharsetSize(password)
	basicEntropy := float64(len(password)) * logBase2(float64(charsetSize))
	
	// Apply entropy reduction factors
	repetitionPenalty := s.calculateRepetitionPenalty(password)
	patternPenalty := s.calculatePatternPenalty(password)
	
	adjustedEntropy := basicEntropy * repetitionPenalty * patternPenalty
	
	// Minimum entropy should not be less than 10% of basic entropy
	minEntropy := basicEntropy * 0.1
	if adjustedEntropy < minEntropy {
		adjustedEntropy = minEntropy
	}
	
	return adjustedEntropy
}

// calculateCharsetSize determines the effective charset size
func (s *SecurityAnalyzer) calculateCharsetSize(password string) int {
	hasLower := s.hasLowercase(password)
	hasUpper := s.hasUppercase(password)
	hasNumber := s.hasNumbers(password)
	hasSymbol := s.hasSymbols(password)
	
	size := 0
	if hasLower {
		size += 26
	}
	if hasUpper {
		size += 26
	}
	if hasNumber {
		size += 10
	}
	if hasSymbol {
		size += 32 // Common symbols
	}
	
	// If only unique characters are used, charset size is the unique count
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	
	if len(uniqueChars) < size {
		return len(uniqueChars)
	}
	
	return size
}

// calculateRepetitionPenalty reduces entropy for repeated characters
func (s *SecurityAnalyzer) calculateRepetitionPenalty(password string) float64 {
	if len(password) == 0 {
		return 1.0
	}
	
	charCount := make(map[rune]int)
	for _, char := range password {
		charCount[char]++
	}
	
	totalRepeats := 0
	for _, count := range charCount {
		if count > 1 {
			totalRepeats += count - 1
		}
	}
	
	penalty := 1.0 - (float64(totalRepeats) / float64(len(password)) * 0.5)
	if penalty < 0.3 {
		penalty = 0.3 // Minimum penalty
	}
	
	return penalty
}

// calculatePatternPenalty reduces entropy for common patterns
func (s *SecurityAnalyzer) calculatePatternPenalty(password string) float64 {
	penalty := 1.0
	lower := strings.ToLower(password)
	
	// Sequential characters (abc, 123, etc.)
	if s.hasSequentialChars(lower) {
		penalty *= 0.7
	}
	
	// Keyboard patterns (qwerty, asdf, etc.)
	if s.hasKeyboardPattern(lower) {
		penalty *= 0.6
	}
	
	// Common substitutions (@ for a, 3 for e, etc.)
	if s.hasCommonSubstitutions(password) {
		penalty *= 0.8
	}
	
	// Repeated patterns (abcabc, 123123, etc.)
	if s.hasRepeatedPatterns(lower) {
		penalty *= 0.5
	}
	
	return penalty
}

// hasSequentialChars checks for sequential character patterns
func (s *SecurityAnalyzer) hasSequentialChars(password string) bool {
	if len(password) < 3 {
		return false
	}
	
	sequential := 0
	for i := 0; i < len(password)-1; i++ {
		if password[i]+1 == password[i+1] {
			sequential++
			if sequential >= 2 {
				return true
			}
		} else {
			sequential = 0
		}
	}
	
	return false
}

// hasKeyboardPattern checks for common keyboard patterns
func (s *SecurityAnalyzer) hasKeyboardPattern(password string) bool {
	patterns := []string{
		"qwerty", "asdf", "zxcv", "qwer", "asdf", "zxcv",
		"123456", "abcdef", "qazwsx", "wsxedc",
	}
	
	for _, pattern := range patterns {
		if strings.Contains(password, pattern) {
			return true
		}
	}
	
	return false
}

// hasCommonSubstitutions checks for common character substitutions
func (s *SecurityAnalyzer) hasCommonSubstitutions(password string) bool {
	substitutions := map[string]string{
		"@": "a", "3": "e", "1": "i", "0": "o", "5": "s",
		"7": "t", "4": "a", "8": "b", "6": "g", "2": "z",
	}
	
	subCount := 0
	for sub := range substitutions {
		if strings.Contains(password, sub) {
			subCount++
		}
	}
	
	return subCount >= 2
}

// hasRepeatedPatterns checks for repeated substrings
func (s *SecurityAnalyzer) hasRepeatedPatterns(password string) bool {
	if len(password) < 4 {
		return false
	}
	
	for length := 2; length <= len(password)/2; length++ {
		for start := 0; start <= len(password)-length*2; start++ {
			pattern := password[start : start+length]
			if strings.Contains(password[start+length:], pattern) {
				return true
			}
		}
	}
	
	return false
}

// calculateSecurityLevel determines overall security level
func (s *SecurityAnalyzer) calculateSecurityLevel(entropy float64, length int, password string) SecurityLevel {
	// Base level on entropy
	level := VeryWeak
	
	switch {
	case entropy >= 80:
		level = VeryStrong
	case entropy >= 60:
		level = Strong
	case entropy >= 45:
		level = Good
	case entropy >= 30:
		level = Fair
	case entropy >= 20:
		level = Weak
	default:
		level = VeryWeak
	}
	
	// Adjust for length
	if length < 8 {
		if level > Weak {
			level--
		}
	}
	
	// Check for compromised password
	if s.isCommonPassword(password) {
		level = VeryWeak
	}
	
	return level
}

// estimateCrackTime provides human-readable crack time estimates
func (s *SecurityAnalyzer) estimateCrackTime(entropy float64) string {
	if entropy <= 0 {
		return "Instantly"
	}
	
	// Assume 1 billion guesses per second
	guessesPerSecond := 1e9
	combinations := math.Pow(2, entropy)
	seconds := combinations / (2 * guessesPerSecond) // Average case
	
	switch {
	case seconds < 1:
		return "Instantly"
	case seconds < 60:
		return "Less than a minute"
	case seconds < 3600:
		return fmt.Sprintf("%.0f minutes", seconds/60)
	case seconds < 86400:
		return fmt.Sprintf("%.0f hours", seconds/3600)
	case seconds < 2592000: // 30 days
		return fmt.Sprintf("%.0f days", seconds/86400)
	case seconds < 31536000: // 1 year
		return fmt.Sprintf("%.0f months", seconds/2592000)
	case seconds < 31536000000: // 1000 years
		return fmt.Sprintf("%.0f years", seconds/31536000)
	default:
		return "Centuries"
	}
}

// Character type checking functions
func (s *SecurityAnalyzer) hasLowercase(password string) bool {
	for _, char := range password {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func (s *SecurityAnalyzer) hasUppercase(password string) bool {
	for _, char := range password {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func (s *SecurityAnalyzer) hasNumbers(password string) bool {
	for _, char := range password {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func (s *SecurityAnalyzer) hasSymbols(password string) bool {
	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			return true
		}
	}
	return false
}

func (s *SecurityAnalyzer) hasAmbiguous(password string) bool {
	ambiguous := "0O1lI"
	for _, char := range password {
		if strings.ContainsRune(ambiguous, char) {
			return true
		}
	}
	return false
}

// findCommonWords identifies dictionary words in the password
func (s *SecurityAnalyzer) findCommonWords(password string) []string {
	var found []string
	lower := strings.ToLower(password)
	
	for _, word := range s.commonWords {
		if len(word) >= 3 && strings.Contains(lower, word) {
			found = append(found, word)
		}
	}
	
	return found
}

// isCommonPassword checks if password is in common password lists
func (s *SecurityAnalyzer) isCommonPassword(password string) bool {
	lower := strings.ToLower(password)
	for _, common := range s.commonPasswords {
		if lower == common {
			return true
		}
	}
	return false
}

// generateFeedback provides actionable improvement suggestions
func (s *SecurityAnalyzer) generateFeedback(password string, analysis SecurityAnalysis) []string {
	var feedback []string
	
	if len(password) < 12 {
		feedback = append(feedback, "Use at least 12 characters for better security")
	}
	
	if !analysis.HasLowercase {
		feedback = append(feedback, "Add lowercase letters")
	}
	
	if !analysis.HasUppercase {
		feedback = append(feedback, "Add uppercase letters")
	}
	
	if !analysis.HasNumbers {
		feedback = append(feedback, "Add numbers")
	}
	
	if !analysis.HasSymbols {
		feedback = append(feedback, "Add symbols (!@#$%^&*)")
	}
	
	if len(analysis.CommonWords) > 0 {
		feedback = append(feedback, "Avoid dictionary words")
	}
	
	if analysis.IsCompromised {
		feedback = append(feedback, "This password has been found in data breaches")
	}
	
	if s.hasSequentialChars(strings.ToLower(password)) {
		feedback = append(feedback, "Avoid sequential characters (abc, 123)")
	}
	
	if s.hasKeyboardPattern(strings.ToLower(password)) {
		feedback = append(feedback, "Avoid keyboard patterns (qwerty, asdf)")
	}
	
	if analysis.Level <= Fair {
		feedback = append(feedback, "Consider using a passphrase with multiple words")
	}
	
	return feedback
}



// Helper functions for common passwords and words
func getCommonPasswords() []string {
	return []string{
		"password", "123456", "123456789", "guest", "qwerty", "12345678", "111111", "12345",
		"col123456", "123123", "1234567", "1234", "1234567890", "000000", "555555", "666666",
		"123321", "654321", "7777777", "123", "D1lakiss", "777777", "110110jp", "1111", "987654321",
		"121212", "Gizli", "abc123", "112233", "azerty", "159753", "1q2w3e4r", "54321", "pass@123",
		"222222", "qwertyui", "1234554321", "123qwe", "qwerty123", "password1", "administrator",
		"1111111", "123456a", "qwerty1", "password123", "Passwd", "welcome", "admin", "master",
		"hello", "dragon", "monkey", "letmein", "login", "princess", "qwertyuiop", "solo",
		"passw0rd", "starwars", "shadow", "sunshine", "12345678910", "football", "iloveyou",
		"superman", "trustno1", "jesus", "mustang", "ninja", "michael", "charlie",
	}
}

func getCommonWords() []string {
	return []string{
		"password", "admin", "user", "login", "welcome", "hello", "world", "test", "home",
		"love", "life", "work", "time", "year", "good", "great", "best", "free", "new",
		"first", "last", "long", "little", "right", "big", "high", "different", "small",
		"large", "next", "early", "young", "important", "few", "public", "bad", "same",
		"able", "house", "service", "party", "company", "system", "program", "question",
		"government", "place", "case", "part", "group", "problem", "fact", "hand", "right",
		"thing", "person", "woman", "man", "child", "people", "family", "community", "name",
	}
}

// AnalyzePassword is a convenience function to analyze a password using default analyzer
func AnalyzePassword(password string) *SecurityAnalysis {
	analyzer := NewSecurityAnalyzer()
	analysis := analyzer.Analyze(password)
	return &analysis
}
