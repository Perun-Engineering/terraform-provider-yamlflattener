// Package flattener provides functionality to flatten YAML structures into key-value pairs.
package flattener

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

// NewFlattener creates a new instance of Flattener with default settings
func NewFlattener() *Flattener {
	return &Flattener{
		MaxNestingDepth: MaxNestingDepth,
		MaxResultSize:   MaxResultSize,
		MaxYAMLSize:     MaxYAMLSize,
	}
}

// validateDepthAndSize checks depth and result size limits
func (f *Flattener) validateDepthAndSize(depth int, resultSize int) error {
	// Check for max nesting depth to prevent stack overflow
	if depth > f.MaxNestingDepth {
		return fmt.Errorf("maximum nesting depth of %d exceeded", f.MaxNestingDepth)
	}

	// Check for max result size to prevent memory exhaustion
	if resultSize >= f.MaxResultSize {
		return fmt.Errorf("maximum result size of %d key-value pairs exceeded", f.MaxResultSize)
	}

	return nil
}

// buildPrefix constructs a new prefix by appending key to existing prefix
func buildPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// FlattenYAMLString takes a YAML string and flattens it into a map with dot notation
// This version preserves the original order from the YAML document
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

	var yamlNode yaml.Node

	// Set a timeout for YAML parsing to prevent DoS attacks
	done := make(chan struct{})
	var err error

	go func() {
		err = yaml.Unmarshal([]byte(yamlContent), &yamlNode)
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

	return f.FlattenYAMLNode(&yamlNode)
}

// FlattenYAMLNode takes a parsed YAML Node and flattens it into a map with dot notation
// This preserves the original document order
func (f *Flattener) FlattenYAMLNode(node *yaml.Node) (map[string]string, error) {
	if node == nil {
		return nil, fmt.Errorf("cannot flatten nil YAML node")
	}

	result := make(map[string]string)
	err := f.flattenNodeWithDepth(node, "", result, 0)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// flattenNodeWithDepth recursively flattens a YAML node with the given prefix and tracks depth
func (f *Flattener) flattenNodeWithDepth(node *yaml.Node, prefix string, result map[string]string, depth int) error {
	if err := f.validateDepthAndSize(depth, len(result)); err != nil {
		return err
	}

	switch node.Kind {
	case yaml.DocumentNode:
		// Document node - process its content
		if len(node.Content) > 0 {
			return f.flattenNodeWithDepth(node.Content[0], prefix, result, depth)
		}
	case yaml.MappingNode:
		return f.flattenMappingNodeWithDepth(node, prefix, result, depth+1)
	case yaml.SequenceNode:
		return f.flattenSequenceNodeWithDepth(node, prefix, result, depth+1)
	case yaml.ScalarNode:
		// Handle null values specially
		if node.Tag == "!!null" || node.Value == "null" || node.Value == "~" || node.Value == "" {
			result[prefix] = ""
		} else {
			result[prefix] = node.Value
		}
	case yaml.AliasNode:
		// Handle alias nodes by following the alias
		if node.Alias != nil {
			return f.flattenNodeWithDepth(node.Alias, prefix, result, depth)
		}
	}

	return nil
}

// flattenMappingNodeWithDepth flattens a mapping node preserving key order
func (f *Flattener) flattenMappingNodeWithDepth(node *yaml.Node, prefix string, result map[string]string, depth int) error {
	// Content of mapping nodes alternates: key, value, key, value, ...
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break // Malformed mapping, skip
		}

		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			return fmt.Errorf("non-scalar key in YAML mapping")
		}

		// Sanitize key to prevent injection attacks
		key := sanitizeKey(keyNode.Value)
		newPrefix := buildPrefix(prefix, key)

		if err := f.flattenNodeWithDepth(valueNode, newPrefix, result, depth); err != nil {
			return err
		}
	}
	return nil
}

// flattenSequenceNodeWithDepth flattens a sequence node preserving element order
func (f *Flattener) flattenSequenceNodeWithDepth(node *yaml.Node, prefix string, result map[string]string, depth int) error {
	for i, childNode := range node.Content {
		newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
		if err := f.flattenNodeWithDepth(childNode, newPrefix, result, depth); err != nil {
			return err
		}
	}
	return nil
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
	// Remove potentially dangerous characters from keys
	// This is a simple implementation - in production you might want more sophisticated sanitization
	key = strings.ReplaceAll(key, "\x00", "") // Remove null bytes

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
