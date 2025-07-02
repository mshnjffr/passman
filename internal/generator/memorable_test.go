package generator

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestMemorableGenerator(t *testing.T) {
	// Create a larger wordlist to pass validation (min 100 words)
	wordlist := make([]string, 120)
	baseWords := []string{"apple", "banana", "cherry", "dog", "elephant", "fox", "grape", "house", "ice", "jungle"}
	for i := 0; i < 120; i++ {
		wordlist[i] = baseWords[i%len(baseWords)] + fmt.Sprintf("%d", i)
	}
	
	tests := []struct {
		name      string
		wordCount int
		separator string
		wordlist  []string
		wantErr   bool
	}{
		{
			name:      "Valid 4 words",
			wordCount: 4,
			separator: "-",
			wordlist:  wordlist,
			wantErr:   false,
		},
		{
			name:      "Valid 6 words with space separator",
			wordCount: 6,
			separator: " ",
			wordlist:  wordlist,
			wantErr:   false,
		},
		{
			name:      "Valid single word",
			wordCount: 1,
			separator: "-",
			wordlist:  wordlist,
			wantErr:   false,
		},
		{
			name:      "Invalid zero words",
			wordCount: 0,
			separator: "-",
			wordlist:  wordlist,
			wantErr:   true,
		},
		{
			name:      "Invalid too many words",
			wordCount: 25,
			separator: "-",
			wordlist:  wordlist,
			wantErr:   true,
		},
		{
			name:      "Invalid empty wordlist",
			wordCount: 4,
			separator: "-",
			wordlist:  []string{},
			wantErr:   true,
		},
		{
			name:      "Invalid small wordlist",
			wordCount: 4,
			separator: "-",
			wordlist:  []string{"word1", "word2"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewMemorableGenerator(tt.wordCount, tt.separator, tt.wordlist)
			ctx := context.Background()
			
			passphrase, err := gen.Generate(ctx)
			
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
			
			// Verify structure
			words := strings.Split(passphrase, tt.separator)
			if len(words) != tt.wordCount {
				t.Errorf("Expected %d words, got %d", tt.wordCount, len(words))
			}
			
			// Verify all words are from wordlist
			for _, word := range words {
				if !containsString(tt.wordlist, word) {
					t.Errorf("Word '%s' not found in wordlist", word)
				}
			}
		})
	}
}

func TestMemorableGeneratorEntropy(t *testing.T) {
	wordlist := make([]string, 7776) // EFF wordlist size
	for i := range wordlist {
		wordlist[i] = "word" + string(rune(i))
	}
	
	gen := NewMemorableGenerator(5, "-", wordlist)
	entropy := gen.EstimateEntropy()
	
	// With 7776 words and 5 words, entropy should be around 64.6 bits
	expectedEntropy := 5 * logBase2(7776)
	if entropy < expectedEntropy*0.9 || entropy > expectedEntropy*1.1 {
		t.Errorf("Expected entropy around %.2f, got %.2f", expectedEntropy, entropy)
	}
}

func TestMemorableGeneratorEFFWordlist(t *testing.T) {
	wordlist := GetEFFWordlist()
	
	if len(wordlist) == 0 {
		t.Error("EFF wordlist should not be empty")
	}
	
	// Check that words are reasonable length
	for _, word := range wordlist[:10] { // Check first 10
		if len(word) < 3 || len(word) > 15 {
			t.Errorf("Word '%s' has unusual length %d", word, len(word))
		}
	}
}

func TestMemorableGeneratorSeparator(t *testing.T) {
	// Create a larger wordlist to pass validation
	wordlist := make([]string, 120)
	baseWords := []string{"apple", "banana", "cherry", "dog"}
	for i := 0; i < 120; i++ {
		wordlist[i] = baseWords[i%len(baseWords)] + fmt.Sprintf("%d", i)
	}
	gen := NewMemorableGenerator(3, "-", wordlist)
	
	// Test changing separator
	gen.SetSeparator("_")
	ctx := context.Background()
	
	passphrase, err := gen.Generate(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !strings.Contains(passphrase, "_") {
		t.Error("Passphrase should contain underscore separator")
	}
	
	if strings.Contains(passphrase, "-") {
		t.Error("Passphrase should not contain old separator")
	}
}

func TestMemorableGeneratorUniqueness(t *testing.T) {
	wordlist := make([]string, 1000)
	for i := range wordlist {
		wordlist[i] = "word" + string(rune(i))
	}
	
	gen := NewMemorableGenerator(4, "-", wordlist)
	ctx := context.Background()
	
	passphrases := make(map[string]bool)
	iterations := 50
	
	for i := 0; i < iterations; i++ {
		passphrase, err := gen.Generate(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if passphrases[passphrase] {
			t.Errorf("Generated duplicate passphrase: %s", passphrase)
		}
		passphrases[passphrase] = true
	}
}

func TestMemorableGeneratorWordlistOperations(t *testing.T) {
	wordlist1 := []string{"apple", "banana", "cherry"}
	wordlist2 := []string{"dog", "elephant", "fox"}
	
	gen := NewMemorableGenerator(2, "-", wordlist1)
	
	// Test getting wordlist
	retrieved := gen.GetWordlist()
	if len(retrieved) != len(wordlist1) {
		t.Errorf("Expected %d words, got %d", len(wordlist1), len(retrieved))
	}
	
	// Test setting new wordlist
	gen.SetWordlist(wordlist2)
	newRetrieved := gen.GetWordlist()
	if len(newRetrieved) != len(wordlist2) {
		t.Errorf("Expected %d words after setting, got %d", len(wordlist2), len(newRetrieved))
	}
}

// Helper function for testing
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
