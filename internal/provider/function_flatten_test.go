package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenFunction_Metadata(t *testing.T) {
	f := NewFlattenFunction()

	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}

	f.Metadata(context.Background(), req, resp)

	if resp.Name != "flatten" {
		t.Errorf("Expected function name 'flatten', got %s", resp.Name)
	}
}

func TestFlattenFunction_Definition(t *testing.T) {
	f := NewFlattenFunction()

	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}

	f.Definition(context.Background(), req, resp)

	if len(resp.Definition.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(resp.Definition.Parameters))
	}

	if resp.Definition.Parameters[0].GetName() != "yaml_content" {
		t.Errorf("Expected parameter name 'yaml_content', got %s", resp.Definition.Parameters[0].GetName())
	}

	if resp.Definition.Parameters[1].GetName() != "escape_newlines" {
		t.Errorf("Expected parameter name 'escape_newlines', got %s", resp.Definition.Parameters[1].GetName())
	}
}

func TestFlattenFunction_Run_SimpleObject(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `
key1: value1
key2:
  nested: value2
`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	elements := result.Elements()

	// Check expected keys
	expectedKeys := map[string]string{
		"key1":        "value1",
		"key2.nested": "value2",
	}

	if len(elements) != len(expectedKeys) {
		t.Errorf("Expected %d elements, got %d", len(expectedKeys), len(elements))
	}

	for expectedKey, expectedValue := range expectedKeys {
		if element, exists := elements[expectedKey]; exists {
			actualValue := element.(types.String).ValueString()
			if actualValue != expectedValue {
				t.Errorf("Expected %s=%s, got %s=%s", expectedKey, expectedValue, expectedKey, actualValue)
			}
		} else {
			t.Errorf("Expected key %s not found in result", expectedKey)
		}
	}
}

func TestFlattenFunction_Run_WithArray(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `
items:
  - name: item1
    value: val1
  - name: item2
    value: val2
`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	elements := result.Elements()

	// Check expected keys
	expectedKeys := map[string]string{
		"items[0].name":  "item1",
		"items[0].value": "val1",
		"items[1].name":  "item2",
		"items[1].value": "val2",
	}

	if len(elements) != len(expectedKeys) {
		t.Errorf("Expected %d elements, got %d", len(expectedKeys), len(elements))
	}

	for expectedKey, expectedValue := range expectedKeys {
		if element, exists := elements[expectedKey]; exists {
			actualValue := element.(types.String).ValueString()
			if actualValue != expectedValue {
				t.Errorf("Expected %s=%s, got %s=%s", expectedKey, expectedValue, expectedKey, actualValue)
			}
		} else {
			t.Errorf("Expected key %s not found in result", expectedKey)
		}
	}
}

func TestFlattenFunction_Run_EmptyContent(t *testing.T) {
	f := NewFlattenFunction()

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(""),
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("Expected error for empty YAML content, got nil")
	}
}

func TestFlattenFunction_Run_InvalidYAML(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `
key1: value1
key2: [
  invalid yaml
`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("Expected error for invalid YAML content, got nil")
	}
}

func TestFlattenFunction_Run_ComplexExample(t *testing.T) {
	f := NewFlattenFunction()

	// Using the alertmanager example from requirements
	yamlContent := `
alertmanager:
  config:
    global:
      slack_api_url: "your-encrypted-slack-webhook"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "your-encrypted-webhook-url"
`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	elements := result.Elements()

	// Check some expected keys
	expectedKeys := map[string]string{
		"alertmanager.config.global.slack_api_url":                  "your-encrypted-slack-webhook",
		"alertmanager.config.receivers[0].name":                     "slack-notifications",
		"alertmanager.config.receivers[0].slack_configs[0].api_url": "your-encrypted-webhook-url",
	}

	for expectedKey, expectedValue := range expectedKeys {
		if element, exists := elements[expectedKey]; exists {
			actualValue := element.(types.String).ValueString()
			if actualValue != expectedValue {
				t.Errorf("Expected %s=%s, got %s=%s", expectedKey, expectedValue, expectedKey, actualValue)
			}
		} else {
			t.Errorf("Expected key %s not found in result", expectedKey)
		}
	}
}

func TestFlattenFunction_Run_DataTypes(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	elements := result.Elements()

	// Check expected keys and their string representations
	expectedKeys := map[string]string{
		"string_val": "hello",
		"int_val":    "42",
		"float_val":  "3.14",
		"bool_val":   "true",
		"null_val":   "",
	}

	for expectedKey, expectedValue := range expectedKeys {
		if element, exists := elements[expectedKey]; exists {
			actualValue := element.(types.String).ValueString()
			if actualValue != expectedValue {
				t.Errorf("Expected %s=%s, got %s=%s", expectedKey, expectedValue, expectedKey, actualValue)
			}
		} else {
			t.Errorf("Expected key %s not found in result", expectedKey)
		}
	}
}

func TestFlattenFunction_Run_EmptyYAML(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `{}`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	elements := result.Elements()

	// Empty YAML should result in empty map
	if len(elements) != 0 {
		t.Errorf("Expected empty result for empty YAML, got %d elements", len(elements))
	}
}

func TestFlattenFunction_Run_NestedArrays(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `
matrix:
  - - 1
    - 2
  - - 3
    - 4
`

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	elements := result.Elements()

	// Check expected keys for nested arrays
	expectedKeys := map[string]string{
		"matrix[0][0]": "1",
		"matrix[0][1]": "2",
		"matrix[1][0]": "3",
		"matrix[1][1]": "4",
	}

	for expectedKey, expectedValue := range expectedKeys {
		if element, exists := elements[expectedKey]; exists {
			actualValue := element.(types.String).ValueString()
			if actualValue != expectedValue {
				t.Errorf("Expected %s=%s, got %s=%s", expectedKey, expectedValue, expectedKey, actualValue)
			}
		} else {
			t.Errorf("Expected key %s not found in result", expectedKey)
		}
	}
}

func TestFlattenFunction_Run_WhitespaceOnly(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := "   \n   \t   "

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(false), // escape_newlines = false
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error == nil {
		t.Error("Expected error for whitespace-only YAML content, got nil")
	}
}

func TestFlattenFunction_Run_EscapeNewlines(t *testing.T) {
	f := NewFlattenFunction()

	yamlContent := `
alertmanager:
  config:
    receivers:
      - name: discord_prometheus
        webhook_configs:
          - url: https://example.com/webhook
            body: |
              {
                "content": "**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
              }
`

	// Test with escape_newlines = true
	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.StringValue(yamlContent),
			types.BoolValue(true), // escape_newlines = true
		}),
	}
	resp := &function.RunResponse{}

	f.Run(context.Background(), req, resp)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %v", resp.Error)
	}

	// Get the result directly as a Map value
	resultValue := resp.Result.Value()
	result, ok := resultValue.(types.Map)
	if !ok {
		t.Fatalf("Expected result to be types.Map, got %T", resultValue)
	}

	// Check that the body contains escaped newlines
	bodyKey := "alertmanager.config.receivers[0].webhook_configs[0].body"
	bodyValue, exists := result.Elements()[bodyKey]
	if !exists {
		t.Fatalf("Expected key '%s' not found in result", bodyKey)
	}

	bodyString := bodyValue.(types.String).ValueString()
	if !strings.Contains(bodyString, "\\n") {
		t.Errorf("Expected body to contain escaped newlines (\\n), got: %s", bodyString)
	}

	// Verify newlines are escaped
	expectedContent := "{\\n  \"content\": \"**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}\"\\n}"
	if bodyString != expectedContent {
		t.Errorf("Expected body content '%s', got '%s'", expectedContent, bodyString)
	}
}
