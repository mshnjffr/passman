package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExportFormat represents different export formats
type ExportFormat string

const (
	FormatText ExportFormat = "txt"
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// PasswordEntry represents a password entry for export
type PasswordEntry struct {
	Password    string    `json:"password"`
	Length      int       `json:"length"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description,omitempty"`
}

// ExportManager handles password export operations
type ExportManager struct{}

// NewExportManager creates a new export manager instance
func NewExportManager() *ExportManager {
	return &ExportManager{}
}

// ExportSingle exports a single password to a file
func (e *ExportManager) ExportSingle(password, description string, format ExportFormat, filePath string) error {
	entry := PasswordEntry{
		Password:    password,
		Length:      len(password),
		Type:        "generated",
		CreatedAt:   time.Now(),
		Description: description,
	}

	return e.Export([]PasswordEntry{entry}, format, filePath)
}

// Export exports multiple password entries to a file
func (e *ExportManager) Export(entries []PasswordEntry, format ExportFormat, filePath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	switch format {
	case FormatText:
		return e.exportText(entries, filePath)
	case FormatJSON:
		return e.exportJSON(entries, filePath)
	case FormatCSV:
		return e.exportCSV(entries, filePath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportText exports entries as plain text
func (e *ExportManager) exportText(entries []PasswordEntry, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for i, entry := range entries {
		if i > 0 {
			fmt.Fprintln(file, "---")
		}
		
		fmt.Fprintf(file, "Password: %s\n", entry.Password)
		fmt.Fprintf(file, "Length: %d\n", entry.Length)
		fmt.Fprintf(file, "Type: %s\n", entry.Type)
		fmt.Fprintf(file, "Created: %s\n", entry.CreatedAt.Format(time.RFC3339))
		
		if entry.Description != "" {
			fmt.Fprintf(file, "Description: %s\n", entry.Description)
		}
		fmt.Fprintln(file)
	}

	return nil
}

// exportJSON exports entries as JSON
func (e *ExportManager) exportJSON(entries []PasswordEntry, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	exportData := struct {
		ExportedAt time.Time       `json:"exported_at"`
		Count      int             `json:"count"`
		Entries    []PasswordEntry `json:"entries"`
	}{
		ExportedAt: time.Now(),
		Count:      len(entries),
		Entries:    entries,
	}

	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// exportCSV exports entries as CSV
func (e *ExportManager) exportCSV(entries []PasswordEntry, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{
		"Password", "Length", "Type", "Created At", "Description",
	}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write entries
	for _, entry := range entries {
		record := []string{
			entry.Password,
			fmt.Sprintf("%d", entry.Length),
			entry.Type,
			entry.CreatedAt.Format(time.RFC3339),
			entry.Description,
		}
		
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// GetSuggestedFilename generates a suggested filename for export
func (e *ExportManager) GetSuggestedFilename(format ExportFormat, baseName string) string {
	if baseName == "" {
		baseName = "passwords"
	}
	
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s", baseName, timestamp, string(format))
	
	// Sanitize filename
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, ":", "-")
	
	return filename
}

// ValidateExportPath validates the export path and format
func (e *ExportManager) ValidateExportPath(filePath string, format ExportFormat) error {
	if filePath == "" {
		return fmt.Errorf("export path cannot be empty")
	}

	// Check if we can write to the directory
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
	}

	// Validate format matches extension
	ext := strings.ToLower(filepath.Ext(filePath))
	expectedExt := "." + string(format)
	
	if ext != expectedExt {
		return fmt.Errorf("file extension %s does not match format %s", ext, format)
	}

	return nil
}
