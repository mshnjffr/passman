package generator

import (
	"context"
	"strings"
	"testing"
	"unicode"
)

func TestPINGenerator(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Valid 4-digit PIN",
			length:  4,
			wantErr: false,
		},
		{
			name:    "Valid 6-digit PIN",
			length:  6,
			wantErr: false,
		},
		{
			name:    "Valid long PIN",
			length:  16,
			wantErr: false,
		},
		{
			name:    "Invalid zero length",
			length:  0,
			wantErr: true,
		},
		{
			name:    "Invalid negative length",
			length:  -1,
			wantErr: true,
		},
		{
			name:    "Invalid too long",
			length:  100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewPINGenerator(tt.length)
			ctx := context.Background()
			
			pin, err := gen.Generate(ctx)
			
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
			
			if len(pin) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(pin))
			}
			
			// Verify all characters are digits
			for _, r := range pin {
				if !unicode.IsDigit(r) {
					t.Errorf("PIN contains non-digit character: %c", r)
				}
			}
		})
	}
}

func TestPINGeneratorEntropy(t *testing.T) {
	gen := NewPINGenerator(6)
	entropy := gen.EstimateEntropy()
	
	// With 10 digits and 6 length, entropy should be around 19.93 bits
	expectedEntropy := 6 * logBase2(10)
	if entropy < expectedEntropy*0.9 || entropy > expectedEntropy*1.1 {
		t.Errorf("Expected entropy around %.2f, got %.2f", expectedEntropy, entropy)
	}
}

func TestPINGeneratorFormatted(t *testing.T) {
	gen := NewPINGenerator(8)
	ctx := context.Background()
	
	tests := []struct {
		name      string
		separator string
		groupSize int
		want      func(string) bool
	}{
		{
			name:      "Hyphen separated groups of 4",
			separator: "-",
			groupSize: 4,
			want: func(pin string) bool {
				return strings.Contains(pin, "-") && len(strings.Split(pin, "-")) == 2
			},
		},
		{
			name:      "Space separated groups of 2",
			separator: " ",
			groupSize: 2,
			want: func(pin string) bool {
				parts := strings.Split(pin, " ")
				return len(parts) == 4 && len(parts[0]) == 2
			},
		},
		{
			name:      "No formatting (empty separator)",
			separator: "",
			groupSize: 4,
			want: func(pin string) bool {
				return !strings.Contains(pin, "-") && !strings.Contains(pin, " ")
			},
		},
		{
			name:      "No formatting (zero group size)",
			separator: "-",
			groupSize: 0,
			want: func(pin string) bool {
				return !strings.Contains(pin, "-")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin, err := gen.GenerateFormatted(ctx, tt.separator, tt.groupSize)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if !tt.want(pin) {
				t.Errorf("PIN format doesn't match expectation: %s", pin)
			}
			
			// Remove separators and verify it's still all digits
			cleanPin := strings.ReplaceAll(pin, tt.separator, "")
			if len(cleanPin) != 8 {
				t.Errorf("Expected 8 digits after removing separators, got %d", len(cleanPin))
			}
			
			for _, r := range cleanPin {
				if !unicode.IsDigit(r) {
					t.Errorf("Cleaned PIN contains non-digit: %c", r)
				}
			}
		})
	}
}

func TestPINGeneratorUniqueness(t *testing.T) {
	gen := NewPINGenerator(6)
	ctx := context.Background()
	
	pins := make(map[string]bool)
	iterations := 100
	
	for i := 0; i < iterations; i++ {
		pin, err := gen.Generate(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if pins[pin] {
			t.Errorf("Generated duplicate PIN: %s", pin)
		}
		pins[pin] = true
	}
}

func TestPINGeneratorSetLength(t *testing.T) {
	gen := NewPINGenerator(4)
	
	// Change length
	gen.SetLength(8)
	
	ctx := context.Background()
	pin, err := gen.Generate(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	if len(pin) != 8 {
		t.Errorf("Expected length 8 after SetLength, got %d", len(pin))
	}
}

func TestPINGeneratorDistribution(t *testing.T) {
	gen := NewPINGenerator(50) // Maximum allowed PIN length to test distribution
	ctx := context.Background()
	
	pin, err := gen.Generate(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	// Count digit frequency
	digitCount := make(map[rune]int)
	for _, r := range pin {
		digitCount[r]++
	}
	
	// Should have all 10 digits represented (with high probability)
	if len(digitCount) < 8 { // Allow some variance
		t.Errorf("Expected at least 8 different digits, got %d", len(digitCount))
	}
	
	// Check for reasonable distribution (no digit should be more than 20% of total)
	maxCount := len(pin) / 5
	for digit, count := range digitCount {
		if count > maxCount {
			t.Errorf("Digit %c appears too frequently: %d times (max expected: %d)", digit, count, maxCount)
		}
	}
}

func TestPINGeneratorContextCancellation(t *testing.T) {
	gen := NewPINGenerator(6)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	_, err := gen.Generate(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}
