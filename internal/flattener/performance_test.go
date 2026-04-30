package flattener

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestLargeYAMLPerformance tests the performance of flattening large YAML structures
func TestLargeYAMLPerformance(t *testing.T) {
	flattener := New()

	tests := []struct {
		name      string
		generator func() string
		maxTime   time.Duration // Maximum allowed time for processing
	}{
		{
			name: "Large Nested Object (1000 keys)",
			generator: func() string {
				var builder strings.Builder
				builder.WriteString("root:\n")
				for i := 0; i < 1000; i++ {
					fmt.Fprintf(&builder, "  key%d: value%d\n", i, i)
				}
				return builder.String()
			},
			maxTime: 1 * time.Second,
		},
		{
			name: "Deep Nesting (50 levels)",
			generator: func() string {
				var builder strings.Builder
				indent := ""
				key := "root"
				for i := 0; i < 50; i++ {
					fmt.Fprintf(&builder, "%s%s:\n", indent, key)
					indent += "  "
					key = "nested"
				}
				fmt.Fprintf(&builder, "%svalue: \"deep value\"\n", indent)
				return builder.String()
			},
			maxTime: 1 * time.Second,
		},
		{
			name: "Large Array (1000 items)",
			generator: func() string {
				var builder strings.Builder
				builder.WriteString("items:\n")
				for i := 0; i < 1000; i++ {
					fmt.Fprintf(&builder, "  - item%d\n", i)
				}
				return builder.String()
			},
			maxTime: 1 * time.Second,
		},
		{
			name: "Complex Mixed Structure (500 nested items)",
			generator: func() string {
				var builder strings.Builder
				builder.WriteString("root:\n")
				for i := 0; i < 500; i++ {
					fmt.Fprintf(&builder, "  group%d:\n", i)
					fmt.Fprintf(&builder, "    name: \"Group %d\"\n", i)
					builder.WriteString("    items:\n")
					for j := 0; j < 5; j++ {
						fmt.Fprintf(&builder, "      - id: %d\n", j)
						fmt.Fprintf(&builder, "        name: \"Item %d-%d\"\n", i, j)
						builder.WriteString("        attributes:\n")
						fmt.Fprintf(&builder, "          attr1: \"value-%d-%d-1\"\n", i, j)
						fmt.Fprintf(&builder, "          attr2: \"value-%d-%d-2\"\n", i, j)
					}
				}
				return builder.String()
			},
			maxTime: 2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlContent := tt.generator()

			start := time.Now()
			result, err := flattener.FlattenYAMLString(yamlContent)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("Failed to flatten YAML: %v", err)
			}

			t.Logf("Flattened %s in %v, resulting in %d key-value pairs",
				tt.name, duration, len(result))

			if duration > tt.maxTime {
				t.Errorf("Performance test failed: processing took %v, which exceeds maximum allowed time of %v",
					duration, tt.maxTime)
			}
		})
	}
}

// TestMemoryLimits tests that the memory limits are enforced
func TestMemoryLimits(t *testing.T) {
	tests := []struct {
		name            string
		yamlGenerator   func() string
		maxNestingDepth int
		maxResultSize   int
		shouldError     bool
	}{
		{
			name: "Exceeding max nesting depth",
			yamlGenerator: func() string {
				var builder strings.Builder
				indent := ""
				key := "root"
				for i := 0; i < 150; i++ { // Generate 150 levels of nesting
					fmt.Fprintf(&builder, "%s%s:\n", indent, key)
					indent += "  "
					key = "nested"
				}
				fmt.Fprintf(&builder, "%svalue: \"deep value\"\n", indent)
				return builder.String()
			},
			maxNestingDepth: 100, // Set limit to 100
			maxResultSize:   100000,
			shouldError:     true,
		},
		{
			name: "Exceeding max result size",
			yamlGenerator: func() string {
				var builder strings.Builder
				builder.WriteString("root:\n")
				for i := 0; i < 20000; i++ { // Generate 20000 keys
					fmt.Fprintf(&builder, "  key%d: value%d\n", i, i)
				}
				return builder.String()
			},
			maxNestingDepth: 100,
			maxResultSize:   10000, // Set limit to 10000
			shouldError:     true,
		},
		{
			name: "Within limits",
			yamlGenerator: func() string {
				var builder strings.Builder
				builder.WriteString("root:\n")
				for i := 0; i < 50; i++ {
					fmt.Fprintf(&builder, "  key%d: value%d\n", i, i)
				}
				return builder.String()
			},
			maxNestingDepth: 100,
			maxResultSize:   10000,
			shouldError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flattener := New()
			flattener.MaxNestingDepth = tt.maxNestingDepth
			flattener.MaxResultSize = tt.maxResultSize

			yamlContent := tt.yamlGenerator()

			_, err := flattener.FlattenYAMLString(yamlContent)

			if tt.shouldError && err == nil {
				t.Errorf("Expected error due to exceeding limits, but got none")
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

// TestSecurityMeasures tests the security measures
func TestSecurityMeasures(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		shouldError bool
	}{
		{
			name:        "YAML with null bytes",
			yamlContent: "key: value\x00malicious",
			shouldError: false, // Should sanitize, not error
		},
		{
			name:        "Extremely large content",
			yamlContent: strings.Repeat("a: b\n", 11*1024*1024), // 11MB
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flattener := New()

			_, err := flattener.FlattenYAMLString(tt.yamlContent)

			if tt.shouldError && err == nil {
				t.Errorf("Expected security error, but got none")
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}
