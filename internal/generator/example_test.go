package generator

import (
	"context"
	"fmt"
	"log"
)

// Example demonstrates how to use all the password generators
func ExampleGenerators() {
	ctx := context.Background()
	analyzer := NewSecurityAnalyzer()

	// Random Password Generator
	fmt.Println("=== Random Password Generator ===")
	randomGen := NewRandomGenerator(16, Lowercase, Uppercase, Numbers, Symbols)
	randomPassword, err := randomGen.Generate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Password: %s\n", randomPassword)
	fmt.Printf("Entropy: %.2f bits\n", randomGen.EstimateEntropy())
	
	analysis := analyzer.Analyze(randomPassword)
	fmt.Printf("Security Level: %s\n", SecurityLevelToString(analysis.Level))
	fmt.Printf("Crack Time: %s\n", analysis.CrackTime)
	fmt.Println()

	// Memorable Passphrase Generator
	fmt.Println("=== Memorable Passphrase Generator ===")
	wordlist := GetEFFWordlist()
	memorableGen := NewMemorableGenerator(4, "-", wordlist)
	passphrase, err := memorableGen.Generate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Passphrase: %s\n", passphrase)
	fmt.Printf("Entropy: %.2f bits\n", memorableGen.EstimateEntropy())
	
	analysis = analyzer.Analyze(passphrase)
	fmt.Printf("Security Level: %s\n", SecurityLevelToString(analysis.Level))
	fmt.Printf("Crack Time: %s\n", analysis.CrackTime)
	fmt.Println()

	// PIN Generator
	fmt.Println("=== PIN Generator ===")
	pinGen := NewPINGenerator(6)
	pin, err := pinGen.Generate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PIN: %s\n", pin)
	fmt.Printf("Entropy: %.2f bits\n", pinGen.EstimateEntropy())
	
	// Formatted PIN
	formattedPIN, err := pinGen.GenerateFormatted(ctx, "-", 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Formatted PIN: %s\n", formattedPIN)
	fmt.Println()

	// Security Analysis Example
	fmt.Println("=== Security Analysis Example ===")
	testPassword := "MyTestP@ssw0rd123"
	analysis = analyzer.Analyze(testPassword)
	
	fmt.Printf("Password: %s\n", testPassword)
	fmt.Printf("Entropy: %.2f bits\n", analysis.Entropy)
	fmt.Printf("Security Level: %s\n", SecurityLevelToString(analysis.Level))
	fmt.Printf("Crack Time: %s\n", analysis.CrackTime)
	fmt.Printf("Character Types: L:%v U:%v N:%v S:%v A:%v\n", 
		analysis.HasLowercase, analysis.HasUppercase, 
		analysis.HasNumbers, analysis.HasSymbols, analysis.HasAmbiguous)
	
	if len(analysis.Feedback) > 0 {
		fmt.Println("Feedback:")
		for _, feedback := range analysis.Feedback {
			fmt.Printf("  - %s\n", feedback)
		}
	}
	
	if len(analysis.CommonWords) > 0 {
		fmt.Printf("Common words found: %v\n", analysis.CommonWords)
	}
	
	fmt.Printf("Is compromised: %v\n", analysis.IsCompromised)
}

// Example shows how to use custom character sets
func ExampleCustomCharSets() {
	ctx := context.Background()
	
	// Password with only alphanumeric (no symbols)
	gen1 := NewRandomGenerator(12, Lowercase, Uppercase, Numbers)
	gen1.SetExcludeChars("0O1lI") // Exclude ambiguous characters
	
	password1, _ := gen1.Generate(ctx)
	fmt.Printf("Alphanumeric (no ambiguous): %s\n", password1)
	
	// Password with all character types
	gen2 := NewRandomGenerator(16, Lowercase, Uppercase, Numbers, Symbols)
	password2, _ := gen2.Generate(ctx)
	fmt.Printf("Full character set: %s\n", password2)
	
	// Symbols only password (for special use cases)
	gen3 := NewRandomGenerator(8, Symbols)
	password3, _ := gen3.Generate(ctx)
	fmt.Printf("Symbols only: %s\n", password3)
}

// Example shows different memorable passphrase configurations
func ExampleMemorableVariations() {
	ctx := context.Background()
	wordlist := GetEFFWordlist()
	
	// Short passphrase with spaces
	gen1 := NewMemorableGenerator(3, " ", wordlist)
	phrase1, _ := gen1.Generate(ctx)
	fmt.Printf("3 words with spaces: %s\n", phrase1)
	
	// Longer passphrase with dashes
	gen2 := NewMemorableGenerator(6, "-", wordlist)
	phrase2, _ := gen2.Generate(ctx)
	fmt.Printf("6 words with dashes: %s\n", phrase2)
	
	// Camel case style (no separator)
	gen3 := NewMemorableGenerator(4, "", wordlist)
	phrase3, _ := gen3.Generate(ctx)
	fmt.Printf("4 words joined: %s\n", phrase3)
}
