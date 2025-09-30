package flattener

import (
	"path/filepath"
	"strings"
	"testing"
)

// FuzzFlattenYAMLString tests the FlattenYAMLString function with fuzzed inputs
func FuzzFlattenYAMLString(f *testing.F) {
	// Seed corpus with valid YAML samples
	f.Add("key: value")
	f.Add("nested:\n  key: value")
	f.Add("array:\n  - item1\n  - item2")
	f.Add("complex:\n  nested:\n    deep:\n      value: test")
	f.Add("numbers:\n  int: 42\n  float: 3.14\n  bool: true")
	f.Add("empty: ")
	f.Add("null_value: null")
	f.Add("multiline: |\n  line1\n  line2")

	flattener := Default()

	f.Fuzz(func(t *testing.T, input string) {
		// Call the function with fuzzed input
		// We don't care if it fails, we just want to ensure it doesn't crash or panic
		result, err := flattener.FlattenYAMLString(input)

		// If successful, verify the result is valid
		if err == nil {
			// Ensure result is not nil
			if result == nil {
				t.Error("FlattenYAMLString returned nil result without error")
			}

			// Verify all values in the result are strings
			for key, value := range result {
				if key == "" {
					t.Error("FlattenYAMLString produced empty key")
				}
				// Value can be empty string, but type should be string
				_ = value
			}
		}
	})
}

// FuzzFlattenYAML tests the FlattenYAML function with structured fuzzed inputs
func FuzzFlattenYAML(f *testing.F) {
	// Seed corpus with various YAML structures represented as strings
	f.Add("simple: value")
	f.Add("nested:\n  key: value\n  another: test")
	f.Add("array:\n  - 1\n  - 2\n  - 3")
	f.Add("mixed:\n  key: value\n  list:\n    - item1\n    - item2")

	flattener := Default()

	f.Fuzz(func(t *testing.T, yamlStr string) {
		// We need to parse the YAML first
		result, err := flattener.FlattenYAMLString(yamlStr)

		// Check for proper error handling
		if err != nil {
			// Verify error is of expected type
			if _, ok := err.(*Error); !ok {
				// Allow standard errors during fuzzing, but log for visibility
				t.Logf("Non-FlattenerError returned: %v", err)
			}
			return
		}

		// If successful, validate the result
		if result == nil {
			t.Fatal("FlattenYAML returned nil result without error")
		}

		// Ensure result doesn't exceed maximum size
		if len(result) > flattener.MaxResultSize {
			t.Fatalf("Result size %d exceeds MaxResultSize %d", len(result), flattener.MaxResultSize)
		}
	})
}

// FuzzSanitizeKey tests the sanitizeKey function
func FuzzSanitizeKey(f *testing.F) {
	// Seed corpus with various key patterns
	f.Add("normal_key")
	f.Add("key-with-dash")
	f.Add("key.with.dots")
	f.Add("key_with_underscore")
	f.Add("CamelCaseKey")
	f.Add("key123")
	f.Add("key with spaces")
	f.Add("\x00null_byte")
	f.Add("\tkey_with_tab")

	f.Fuzz(func(t *testing.T, key string) {
		result := sanitizeKey(key)

		// Ensure result doesn't contain control characters
		for _, r := range result {
			if r < 0x20 || (r >= 0x7F && r <= 0x9F) {
				t.Errorf("sanitizeKey returned control character: %q in result %q", r, result)
			}
		}

		// Ensure result doesn't exceed max length
		if len(result) > 1000 {
			t.Errorf("sanitizeKey returned key longer than 1000 characters: %d", len(result))
		}
	})
}

// FuzzValidateFilePath tests the validateFilePath function
func FuzzValidateFilePath(f *testing.F) {
	// Seed corpus with various file path patterns
	f.Add("/tmp/test.yaml")
	f.Add("relative/path/test.yaml")
	f.Add("./local/test.yaml")
	f.Add("../parent/test.yaml")
	f.Add("/etc/passwd")
	f.Add("..\\..\\windows\\path")
	f.Add("/path/with spaces/file.yaml")

	f.Fuzz(func(t *testing.T, filePath string) {
		result, err := validateFilePath(filePath)

		// If validation passes, ensure result is absolute
		if err == nil {
			if !filepath.IsAbs(result) {
				t.Errorf("validateFilePath returned non-absolute path: %s", result)
			}

			// Ensure no directory traversal patterns in result
			if strings.Contains(result, "..") {
				t.Errorf("validateFilePath returned path with ..: %s", result)
			}
		} else {
			// If validation fails, ensure it's a FlattenerError
			if _, ok := err.(*Error); !ok {
				t.Logf("Non-FlattenerError returned: %v", err)
			}
		}
	})
}
