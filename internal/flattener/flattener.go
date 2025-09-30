// Package flattener provides functionality to flatten YAML structures into key-value pairs.
package flattener

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// MaxYAMLSize defines the maximum size of YAML content in bytes (10MB)
	MaxYAMLSize = 10 * 1024 * 1024

	// MaxNestingDepth defines the maximum allowed nesting depth to prevent stack overflow
	MaxNestingDepth = 100

	// MaxResultSize defines the maximum number of key-value pairs in the result
	MaxResultSize = 100000
)

// Flattener provides functionality to flatten nested YAML structures
type Flattener struct {
	// Configuration options
	MaxNestingDepth int
	MaxResultSize   int
	MaxYAMLSize     int
}

// Default creates a Flattener instance with default settings
func Default() *Flattener {
	return &Flattener{
		MaxNestingDepth: MaxNestingDepth,
		MaxResultSize:   MaxResultSize,
		MaxYAMLSize:     MaxYAMLSize,
	}
}

// FlattenYAML takes a parsed YAML structure and flattens it into a map with dot notation
func (f *Flattener) FlattenYAML(yamlData interface{}) (map[string]string, error) {
	if yamlData == nil {
		return nil, fmt.Errorf("cannot flatten nil YAML data")
	}

	result := make(map[string]string)
	err := f.flattenValueWithDepth(yamlData, "", result, 0)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// flattenValueWithDepth recursively flattens a YAML value with the given prefix and tracks depth
func (f *Flattener) flattenValueWithDepth(value interface{}, prefix string, result map[string]string, depth int) error {
	// Check for max nesting depth to prevent stack overflow
	if depth > f.MaxNestingDepth {
		return fmt.Errorf("maximum nesting depth of %d exceeded", f.MaxNestingDepth)
	}

	// Check for max result size to prevent memory exhaustion
	if len(result) >= f.MaxResultSize {
		return fmt.Errorf("maximum result size of %d key-value pairs exceeded", f.MaxResultSize)
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return f.flattenMapWithDepth(v, prefix, result, depth+1)
	case map[interface{}]interface{}:
		return f.flattenInterfaceMapWithDepth(v, prefix, result, depth+1)
	case []interface{}:
		return f.flattenArrayWithDepth(v, prefix, result, depth+1)
	case string:
		result[prefix] = v
	case int:
		result[prefix] = strconv.Itoa(v)
	case int64:
		result[prefix] = strconv.FormatInt(v, 10)
	case float64:
		result[prefix] = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		result[prefix] = strconv.FormatBool(v)
	case nil:
		result[prefix] = ""
	default:
		// For any other type, convert to string
		result[prefix] = fmt.Sprintf("%v", v)
	}

	return nil
}

// flattenMapWithDepth flattens a map[string]interface{} with the given prefix and tracks depth
func (f *Flattener) flattenMapWithDepth(m map[string]interface{}, prefix string, result map[string]string, depth int) error {
	for k, v := range m {
		// Sanitize key to prevent injection attacks
		k = sanitizeKey(k)

		newPrefix := k
		if prefix != "" {
			newPrefix = prefix + "." + k
		}

		if err := f.flattenValueWithDepth(v, newPrefix, result, depth); err != nil {
			return err
		}
	}
	return nil
}

// flattenInterfaceMapWithDepth flattens a map[interface{}]interface{} with the given prefix and tracks depth
func (f *Flattener) flattenInterfaceMapWithDepth(m map[interface{}]interface{}, prefix string, result map[string]string, depth int) error {
	for k, v := range m {
		// Convert key to string
		strKey, ok := k.(string)
		if !ok {
			return fmt.Errorf("non-string key %v in YAML map", k)
		}

		// Sanitize key to prevent injection attacks
		strKey = sanitizeKey(strKey)

		newPrefix := strKey
		if prefix != "" {
			newPrefix = prefix + "." + strKey
		}

		if err := f.flattenValueWithDepth(v, newPrefix, result, depth); err != nil {
			return err
		}
	}
	return nil
}

// flattenArrayWithDepth flattens an array with the given prefix and tracks depth
func (f *Flattener) flattenArrayWithDepth(a []interface{}, prefix string, result map[string]string, depth int) error {
	for i, v := range a {
		newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
		if err := f.flattenValueWithDepth(v, newPrefix, result, depth); err != nil {
			return err
		}
	}
	return nil
}

// FlattenYAMLString takes a YAML string and flattens it into a map with dot notation
func (f *Flattener) FlattenYAMLString(yamlContent string) (map[string]string, error) {
	if yamlContent == "" {
		return nil, fmt.Errorf("YAML content cannot be empty")
	}

	// Check YAML content size
	if len(yamlContent) > f.MaxYAMLSize {
		return nil, fmt.Errorf("YAML content size exceeds maximum allowed size of %d bytes", f.MaxYAMLSize)
	}

	// Sanitize YAML content
	yamlContent = sanitizeYAMLContent(yamlContent)

	var yamlData interface{}

	// Set a timeout for YAML parsing to prevent DoS attacks
	done := make(chan struct{})
	var err error

	go func() {
		err = yaml.Unmarshal([]byte(yamlContent), &yamlData)
		close(done)
	}()

	select {
	case <-done:
		// Parsing completed
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("YAML parsing timed out, content may be too complex")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML content: %w", err)
	}

	return f.FlattenYAML(yamlData)
}

// FlattenYAMLFile reads a YAML file and flattens its content into a map with dot notation
func (f *Flattener) FlattenYAMLFile(filePath string) (map[string]string, error) {
	// Validate and sanitize file path
	cleanPath, err := validateFilePath(filePath)
	if err != nil {
		return nil, err
	}

	// Check file size before reading
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access YAML file: %w", err)
	}

	if fileInfo.Size() > int64(f.MaxYAMLSize) {
		return nil, fmt.Errorf("YAML file size exceeds maximum allowed size of %d bytes", f.MaxYAMLSize)
	}

	content, err := os.ReadFile(cleanPath) // #nosec G304 - cleanPath is validated
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	return f.FlattenYAMLString(string(content))
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

// sanitizeKey sanitizes a map key to prevent injection attacks
func sanitizeKey(key string) string {
	// Remove null bytes and other control characters
	key = strings.Map(func(r rune) rune {
		// Allow printable ASCII characters, common Unicode characters, and common punctuation
		// Remove control characters (0x00-0x1F, 0x7F-0x9F)
		if r < 0x20 || (r >= 0x7F && r <= 0x9F) {
			return -1 // Remove character
		}
		return r
	}, key)

	// Trim whitespace
	key = strings.TrimSpace(key)

	// Limit key length to prevent DoS
	const maxKeyLength = 1000
	if len(key) > maxKeyLength {
		return key[:maxKeyLength]
	}

	return key
}

// sanitizeYAMLContent performs basic sanitization of YAML content
func sanitizeYAMLContent(content string) string {
	// Remove null bytes which could be used in some injection attacks
	return strings.ReplaceAll(content, "\x00", "")
}
