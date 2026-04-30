package flattener

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFlattenYAML(t *testing.T) {
	flattener := New()

	tests := []struct {
		name     string
		yamlStr  string
		expected map[string]string
		wantErr  bool
	}{
		{
			name:    "Simple key-value",
			yamlStr: "key: value",
			expected: map[string]string{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "Nested object",
			yamlStr: `
key1: value1
key2:
  nested: value2
`,
			expected: map[string]string{
				"key1":        "value1",
				"key2.nested": "value2",
			},
			wantErr: false,
		},
		{
			name: "Array values",
			yamlStr: `
items:
  - item1
  - item2
`,
			expected: map[string]string{
				"items[0]": "item1",
				"items[1]": "item2",
			},
			wantErr: false,
		},
		{
			name: "Mixed nested structure",
			yamlStr: `
alertmanager:
  config:
    global:
      slack_api_url: "your-encrypted-slack-webhook"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "your-encrypted-webhook-url"
`,
			expected: map[string]string{
				"alertmanager.config.global.slack_api_url":                  "your-encrypted-slack-webhook",
				"alertmanager.config.receivers[0].name":                     "slack-notifications",
				"alertmanager.config.receivers[0].slack_configs[0].api_url": "your-encrypted-webhook-url",
			},
			wantErr: false,
		},
		{
			name: "Different data types",
			yamlStr: `
string_value: "string"
int_value: 42
float_value: 3.14
bool_value: true
null_value: null
`,
			expected: map[string]string{
				"string_value": "string",
				"int_value":    "42",
				"float_value":  "3.14",
				"bool_value":   "true",
				"null_value":   "",
			},
			wantErr: false,
		},
		{
			name: "Empty array and object",
			yamlStr: `
empty_array: []
empty_object: {}
`,
			expected: map[string]string{},
			wantErr:  false,
		},
		{
			name:     "Nil input",
			yamlStr:  "",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var yamlData interface{}
			var err error

			if tt.yamlStr != "" {
				err = yaml.Unmarshal([]byte(tt.yamlStr), &yamlData)
				if err != nil {
					t.Fatalf("Failed to parse test YAML: %v", err)
				}
			}

			result, err := flattener.FlattenYAML(yamlData)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("FlattenYAML() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestFlattenEdgeCases(t *testing.T) {
	flattener := New()

	tests := []struct {
		name     string
		yamlStr  string
		expected map[string]string
		wantErr  bool
	}{
		{
			name: "Deeply nested structure",
			yamlStr: `
level1:
  level2:
    level3:
      level4:
        level5:
          level6:
            level7:
              level8:
                level9:
                  level10: "deep value"
`,
			expected: map[string]string{
				"level1.level2.level3.level4.level5.level6.level7.level8.level9.level10": "deep value",
			},
			wantErr: false,
		},
		{
			name: "Array with mixed types",
			yamlStr: `
mixed_array:
  - "string"
  - 42
  - true
  - null
  - key: value
  - [1, 2, 3]
`,
			expected: map[string]string{
				"mixed_array[0]":     "string",
				"mixed_array[1]":     "42",
				"mixed_array[2]":     "true",
				"mixed_array[3]":     "",
				"mixed_array[4].key": "value",
				"mixed_array[5][0]":  "1",
				"mixed_array[5][1]":  "2",
				"mixed_array[5][2]":  "3",
			},
			wantErr: false,
		},
		{
			name: "Special characters in keys",
			yamlStr: `
"key-with-dash": "dash-value"
"key.with.dots": "dots-value"
"key with spaces": "spaces-value"
`,
			expected: map[string]string{
				"key-with-dash":   "dash-value",
				"key.with.dots":   "dots-value",
				"key with spaces": "spaces-value",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var yamlData interface{}
			var err error

			if tt.yamlStr != "" {
				err = yaml.Unmarshal([]byte(tt.yamlStr), &yamlData)
				if err != nil {
					t.Fatalf("Failed to parse test YAML: %v", err)
				}
			}

			result, err := flattener.FlattenYAML(yamlData)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("FlattenYAML() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestFlattenYAMLString(t *testing.T) {
	flattener := New()

	tests := []struct {
		name       string
		yamlString string
		expected   map[string]string
		wantErr    bool
	}{
		{
			name:       "Valid YAML string",
			yamlString: "key: value\nkey2:\n  nested: value2",
			expected: map[string]string{
				"key":         "value",
				"key2.nested": "value2",
			},
			wantErr: false,
		},
		{
			name:       "Empty YAML string",
			yamlString: "",
			expected:   nil,
			wantErr:    true,
		},
		{
			name:       "Invalid YAML string",
			yamlString: "key: : value",
			expected:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := flattener.FlattenYAMLString(tt.yamlString)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenYAMLString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FlattenYAMLString() = %v, want %v", result, tt.expected)
			}
		})
	}
}
