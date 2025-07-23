package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// MaxYAMLSize defines the maximum size of YAML content in bytes (10MB)
	MaxYAMLSize = 10 * 1024 * 1024

	// MaxParsingTime defines the maximum time allowed for parsing YAML in seconds
	MaxParsingTime = 5 * time.Second
)

// ParseYAML parses YAML content from a string and returns the parsed structure
func ParseYAML(content string) (interface{}, error) {
	if content == "" {
		return nil, fmt.Errorf("YAML content cannot be empty")
	}

	// Check content size
	if len(content) > MaxYAMLSize {
		return nil, fmt.Errorf("YAML content size exceeds maximum allowed size of %d bytes", MaxYAMLSize)
	}

	// Sanitize content
	content = sanitizeYAMLContent(content)

	var parsedYAML interface{}

	// Set a timeout for YAML parsing to prevent DoS attacks
	done := make(chan struct{})
	var err error

	go func() {
		err = yaml.Unmarshal([]byte(content), &parsedYAML)
		close(done)
	}()

	select {
	case <-done:
		// Parsing completed
	case <-time.After(MaxParsingTime):
		return nil, fmt.Errorf("YAML parsing timed out, content may be too complex")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML content: %w", err)
	}

	return parsedYAML, nil
}

// ReadYAMLFile reads a YAML file from the given path and returns its content as a string
func ReadYAMLFile(filePath string) (string, error) {
	// Validate and sanitize file path
	cleanPath, err := validateFilePath(filePath)
	if err != nil {
		return "", err
	}

	// Check file size before reading
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("YAML file does not exist: %s", cleanPath)
		}
		return "", fmt.Errorf("error accessing YAML file: %w", err)
	}

	// Check if it's a regular file
	if fileInfo.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file: %s", cleanPath)
	}

	// Check file size
	if fileInfo.Size() > MaxYAMLSize {
		return "", fmt.Errorf("YAML file size exceeds maximum allowed size of %d bytes", MaxYAMLSize)
	}

	// Read file content with timeout
	content, err := readFileWithTimeout(cleanPath, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Sanitize content
	sanitizedContent := sanitizeYAMLContent(string(content))

	return sanitizedContent, nil
}

// ValidateYAML checks if the provided content is valid YAML
func ValidateYAML(content string) error {
	if content == "" {
		return fmt.Errorf("YAML content cannot be empty")
	}

	// Check content size
	if len(content) > MaxYAMLSize {
		return fmt.Errorf("YAML content size exceeds maximum allowed size of %d bytes", MaxYAMLSize)
	}

	// Sanitize content
	content = sanitizeYAMLContent(content)

	// Set a timeout for YAML validation to prevent DoS attacks
	done := make(chan error)

	go func() {
		var parsedYAML interface{}
		err := yaml.Unmarshal([]byte(content), &parsedYAML)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("invalid YAML: %w", err)
		}
	case <-time.After(MaxParsingTime):
		return fmt.Errorf("YAML validation timed out, content may be too complex")
	}

	return nil
}

// ReadAndParseYAMLFile combines reading and parsing a YAML file
func ReadAndParseYAMLFile(filePath string) (interface{}, error) {
	content, err := ReadYAMLFile(filePath)
	if err != nil {
		return nil, err
	}

	return ParseYAML(content)
}

// validateFilePath validates and sanitizes a file path to prevent directory traversal
func validateFilePath(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("file path cannot be empty")
	}

	// Clean and validate the file path to prevent directory traversal
	cleanPath := filepath.Clean(filePath)

	// Check for directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("file path contains invalid directory traversal patterns")
	}

	// Check if path is absolute and within allowed directories
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	// Additional security checks could be added here
	// For example, restricting to specific directories

	return absPath, nil
}

// readFileWithTimeout reads a file with a timeout to prevent hanging on large files
func readFileWithTimeout(filePath string, timeout time.Duration) ([]byte, error) {
	done := make(chan struct{})
	var content []byte
	var err error

	go func() {
		content, err = os.ReadFile(filePath)
		close(done)
	}()

	select {
	case <-done:
		return content, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("file reading timed out, file may be too large or system too busy")
	}
}

// sanitizeYAMLContent performs basic sanitization of YAML content
func sanitizeYAMLContent(content string) string {
	// Remove null bytes which could be used in some injection attacks
	return strings.ReplaceAll(content, "\x00", "")
}
