package generator

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRandomGenerator(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		charSets []CharSet
		wantErr  bool
	}{
		{
			name:     "Valid lowercase only",
			length:   10,
			charSets: []CharSet{Lowercase},
			wantErr:  false,
		},
		{
			name:     "Valid mixed case",
			length:   12,
			charSets: []CharSet{Lowercase, Uppercase},
			wantErr:  false,
		},
		{
			name:     "Valid full character set",
			length:   16,
			charSets: []CharSet{Lowercase, Uppercase, Numbers, Symbols},
			wantErr:  false,
		},
		{
			name:     "Invalid zero length",
			length:   0,
			charSets: []CharSet{Lowercase},
			wantErr:  true,
		},
		{
			name:     "Invalid negative length",
			length:   -1,
			charSets: []CharSet{Lowercase},
			wantErr:  true,
		},
		{
			name:     "Invalid too long",
			length:   2000,
			charSets: []CharSet{Lowercase},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewRandomGenerator(tt.length, tt.charSets...)
			ctx := context.Background()
			
			password, err := gen.Generate(ctx)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(password) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(password))
			}
			
			// Verify character sets are respected
			if contains(tt.charSets, Lowercase) {
				if !hasLowercase(password) {
					t.Error("Password should contain lowercase characters")
				}
			}
			
			if contains(tt.charSets, Uppercase) {
				if !hasUppercase(password) {
					t.Error("Password should contain uppercase characters")
				}
			}
			
			if contains(tt.charSets, Numbers) {
				if !hasNumbers(password) {
					t.Error("Password should contain numbers")
				}
			}
			
			if contains(tt.charSets, Symbols) {
				if !hasSymbols(password) {
					t.Error("Password should contain symbols")
				}
			}
		})
	}
}

func TestRandomGeneratorEntropy(t *testing.T) {
	gen := NewRandomGenerator(12, Lowercase, Uppercase, Numbers, Symbols)
	entropy := gen.EstimateEntropy()
	
	// With 94 characters (26+26+10+32) and 12 length, entropy should be around 79 bits
	expectedEntropy := 12 * logBase2(94)
	if entropy < expectedEntropy*0.9 || entropy > expectedEntropy*1.1 {
		t.Errorf("Expected entropy around %.2f, got %.2f", expectedEntropy, entropy)
	}
}

func TestRandomGeneratorExcludeChars(t *testing.T) {
	gen := NewRandomGenerator(100, Lowercase, Numbers)
	gen.SetExcludeChars("lo1")
	
	ctx := context.Background()
	password, err := gen.Generate(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if strings.ContainsAny(password, "lo1") {
		t.Error("Password contains excluded characters")
	}
}

func TestRandomGeneratorCancelation(t *testing.T) {
	gen := NewRandomGenerator(10, Lowercase)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	_, err := gen.Generate(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestRandomGeneratorTimeout(t *testing.T) {
	gen := NewRandomGenerator(10, Lowercase)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	time.Sleep(1 * time.Millisecond) // Ensure timeout
	
	_, err := gen.Generate(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestRandomGeneratorUniqueness(t *testing.T) {
	gen := NewRandomGenerator(12, Lowercase, Uppercase, Numbers, Symbols)
	ctx := context.Background()
	
	passwords := make(map[string]bool)
	iterations := 100
	
	for i := 0; i < iterations; i++ {
		password, err := gen.Generate(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if passwords[password] {
			t.Errorf("Generated duplicate password: %s", password)
		}
		passwords[password] = true
	}
}

// Helper functions for testing
func hasLowercase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func hasUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func hasNumbers(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func hasSymbols(s string) bool {
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		if strings.ContainsRune(symbols, r) {
			return true
		}
	}
	return false
}
