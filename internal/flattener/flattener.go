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
	MaxNestingDepth int
	MaxResultSize   int
	MaxYAMLSize     int
}

// New creates a Flattener instance with default settings
func New() *Flattener {
	return &Flattener{
		MaxNestingDepth: MaxNestingDepth,
		MaxResultSize:   MaxResultSize,
		MaxYAMLSize:     MaxYAMLSize,
	}
}

// FlattenYAML takes a parsed YAML structure and flattens it into a map with dot notation
func (f *Flattener) FlattenYAML(yamlData interface{}) (map[string]string, error) {
	if yamlData == nil {
		return nil, ValidationError("cannot flatten nil YAML data", nil)
	}

	result := make(map[string]string)
	if err := f.flattenValueWithDepth(yamlData, "", result, 0); err != nil {
		return nil, err
	}

	return result, nil
}

// flattenValueWithDepth recursively flattens a YAML value with the given prefix and tracks depth
func (f *Flattener) flattenValueWithDepth(value interface{}, prefix string, result map[string]string, depth int) error {
	if depth > f.MaxNestingDepth {
		return DepthLimitError(f.MaxNestingDepth)
	}

	if len(result) >= f.MaxResultSize {
		return SizeLimitError(f.MaxResultSize, "result")
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
		result[prefix] = fmt.Sprintf("%v", v)
	}

	return nil
}

// flattenMapWithDepth flattens a map[string]interface{} with the given prefix and tracks depth
func (f *Flattener) flattenMapWithDepth(m map[string]interface{}, prefix string, result map[string]string, depth int) error {
	for k, v := range m {
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
		strKey, ok := k.(string)
		if !ok {
			return ParsingError(fmt.Sprintf("non-string key %v in YAML map", k), nil)
		}
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
		return nil, ValidationError("YAML content cannot be empty", nil)
	}

	if len(yamlContent) > f.MaxYAMLSize {
		return nil, SizeLimitError(f.MaxYAMLSize, "YAML content")
	}

	yamlContent = sanitizeYAMLContent(yamlContent)

	var yamlData interface{}

	done := make(chan struct{})
	var err error

	go func() {
		err = yaml.Unmarshal([]byte(yamlContent), &yamlData)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		return nil, TimeoutError("YAML parsing")
	}

	if err != nil {
		return nil, ParsingError("failed to parse YAML content", err)
	}

	return f.FlattenYAML(yamlData)
}

// FlattenYAMLFile reads a YAML file and flattens it into a map with dot notation.
// It validates the path for security (rejects directory traversal), checks file size
// against MaxYAMLSize, and delegates to FlattenYAMLString for parsing and flattening.
func (f *Flattener) FlattenYAMLFile(path string) (map[string]string, error) {
	if path == "" {
		return nil, ValidationError("file path cannot be empty", nil)
	}

	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return nil, PathSecurityError("file path contains invalid directory traversal patterns")
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, FileAccessError(fmt.Sprintf("invalid file path: %s", err), err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, FileAccessError(fmt.Sprintf("failed to access YAML file: %s", err), err)
	}

	if fileInfo.Size() > int64(f.MaxYAMLSize) {
		return nil, SizeLimitError(f.MaxYAMLSize, "YAML file")
	}

	content, err := os.ReadFile(absPath) // #nosec G304 - absPath is validated
	if err != nil {
		return nil, FileAccessError(fmt.Sprintf("failed to read YAML file: %s", err), err)
	}

	return f.FlattenYAMLString(string(content))
}

// sanitizeKey sanitizes a map key to prevent injection attacks
func sanitizeKey(key string) string {
	key = strings.Map(func(r rune) rune {
		if r < 0x20 || (r >= 0x7F && r <= 0x9F) {
			return -1
		}
		return r
	}, key)
	key = strings.TrimSpace(key)
	const maxKeyLength = 1000
	if len(key) > maxKeyLength {
		return key[:maxKeyLength]
	}
	return key
}

// sanitizeYAMLContent performs basic sanitization of YAML content
func sanitizeYAMLContent(content string) string {
	return strings.ReplaceAll(content, "\x00", "")
}
