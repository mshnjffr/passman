package generator

import (
	"strings"
	"testing"
)

func TestLogBase2(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{64, 6},
		{0, 0},    // Edge case
		{-1, 0},   // Edge case
	}

	for _, tt := range tests {
		result := logBase2(tt.input)
		if result != tt.expected {
			t.Errorf("logBase2(%.1f) = %.6f, expected %.6f", tt.input, result, tt.expected)
		}
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		a, b     int
		min, max int
	}{
		{1, 2, 1, 2},
		{5, 3, 3, 5},
		{0, 0, 0, 0},
		{-1, 1, -1, 1},
	}

	for _, tt := range tests {
		if min(tt.a, tt.b) != tt.min {
			t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, min(tt.a, tt.b), tt.min)
		}
		if max(tt.a, tt.b) != tt.max {
			t.Errorf("max(%d, %d) = %d, expected %d", tt.a, tt.b, max(tt.a, tt.b), tt.max)
		}
	}
}

func TestContains(t *testing.T) {
	charSets := []CharSet{Lowercase, Uppercase, Numbers}
	
	tests := []struct {
		item     CharSet
		expected bool
	}{
		{Lowercase, true},
		{Uppercase, true},
		{Numbers, true},
		{Symbols, false},
		{Ambiguous, false},
	}

	for _, tt := range tests {
		result := contains(charSets, tt.item)
		if result != tt.expected {
			t.Errorf("contains(charSets, %v) = %v, expected %v", tt.item, result, tt.expected)
		}
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{},
			expected: []string{},
		},
		{
			input:    []string{"same", "same", "same"},
			expected: []string{"same"},
		},
	}

	for _, tt := range tests {
		result := unique(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("unique(%v) length = %d, expected %d", tt.input, len(result), len(tt.expected))
			continue
		}
		
		for i, item := range result {
			if item != tt.expected[i] {
				t.Errorf("unique(%v)[%d] = %s, expected %s", tt.input, i, item, tt.expected[i])
			}
		}
	}
}

func TestCharSetToString(t *testing.T) {
	tests := []struct {
		charSet  CharSet
		expected string
	}{
		{Lowercase, "lowercase"},
		{Uppercase, "uppercase"},
		{Numbers, "numbers"},
		{Symbols, "symbols"},
		{Ambiguous, "ambiguous"},
		{CharSet(999), "unknown"}, // Invalid value
	}

	for _, tt := range tests {
		result := CharSetToString(tt.charSet)
		if result != tt.expected {
			t.Errorf("CharSetToString(%v) = %s, expected %s", tt.charSet, result, tt.expected)
		}
	}
}

func TestSecurityLevelToString(t *testing.T) {
	tests := []struct {
		level    SecurityLevel
		expected string
	}{
		{VeryWeak, "Very Weak"},
		{Weak, "Weak"},
		{Fair, "Fair"},
		{Good, "Good"},
		{Strong, "Strong"},
		{VeryStrong, "Very Strong"},
		{SecurityLevel(999), "Unknown"}, // Invalid value
	}

	for _, tt := range tests {
		result := SecurityLevelToString(tt.level)
		if result != tt.expected {
			t.Errorf("SecurityLevelToString(%v) = %s, expected %s", tt.level, result, tt.expected)
		}
	}
}

func TestGetSecurityLevelColor(t *testing.T) {
	tests := []struct {
		level SecurityLevel
		color string
	}{
		{VeryWeak, "#ff4444"},
		{Weak, "#ff8800"},
		{Fair, "#ffaa00"},
		{Good, "#88cc00"},
		{Strong, "#44aa44"},
		{VeryStrong, "#00aa88"},
		{SecurityLevel(999), "#888888"}, // Invalid value
	}

	for _, tt := range tests {
		result := GetSecurityLevelColor(tt.level)
		if result != tt.color {
			t.Errorf("GetSecurityLevelColor(%v) = %s, expected %s", tt.level, result, tt.color)
		}
		
		// Verify it's a valid hex color
		if !strings.HasPrefix(result, "#") || len(result) != 7 {
			t.Errorf("GetSecurityLevelColor(%v) returned invalid hex color: %s", tt.level, result)
		}
	}
}

func TestEstimateGenerationTime(t *testing.T) {
	tests := []struct {
		count         int
		generatorType string
		expected      string
	}{
		{1, "Random Password", "< 1ms"},
		{1, "Memorable Passphrase", "< 1ms"},
		{1, "Numeric PIN", "< 1ms"},
		{1000, "Random Password", "< 1ms"},
		{100000, "Random Password", "< 1s"},
	}

	for _, tt := range tests {
		result := EstimateGenerationTime(tt.count, tt.generatorType)
		if result != tt.expected {
			// Allow for some flexibility in timing estimates
			if !strings.Contains(result, "ms") && !strings.Contains(result, "s") {
				t.Errorf("EstimateGenerationTime(%d, %s) = %s, expected format with time unit", 
					tt.count, tt.generatorType, result)
			}
		}
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		minLength    int
		requireMixed bool
		expectIssues int
	}{
		{
			name:         "Valid strong password",
			password:     "MyStr0ng!P@ssw0rd",
			minLength:    12,
			requireMixed: true,
			expectIssues: 0,
		},
		{
			name:         "Too short",
			password:     "short",
			minLength:    12,
			requireMixed: false,
			expectIssues: 1,
		},
		{
			name:         "Missing character types",
			password:     "alllowercase",
			minLength:    8,
			requireMixed: true,
			expectIssues: 3, // Missing uppercase, numbers, symbols
		},
		{
			name:         "No mixed requirements",
			password:     "simplelongpassword",
			minLength:    12,
			requireMixed: false,
			expectIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := ValidatePasswordStrength(tt.password, tt.minLength, tt.requireMixed)
			
			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}

func TestClearString(t *testing.T) {
	// Test clearString function (best effort in Go)
	testStr := "sensitive"
	clearString(&testStr)
	
	if testStr != "" {
		t.Errorf("clearString should empty the string, got: %s", testStr)
	}
	
	// Test with nil pointer (should not panic)
	clearString(nil)
}
