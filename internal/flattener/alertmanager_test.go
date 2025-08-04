package flattener

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFlattenYAML_AlertmanagerConfig(t *testing.T) {
	flattener := NewFlattener()

	yamlStr := `alertmanager:
    config:
        receivers:
            - name: "null"
            - name: discord_prometheus
              webhook_configs:
                - url: https://example.com/webhook/prometheus
                  send_resolved: true
                  http_config:
                    headers:
                        Content-Type: application/json
                  body: |-
                    {
                      "content": "**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
                    }
            - name: discord_alerts
              webhook_configs:
                - url: https://example.com/webhook/alerts
                  send_resolved: true
                  http_config:
                    headers:
                        Content-Type: application/json
                  body: |-
                    {
                      "content": "**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
                    }`

	expected := map[string]string{
		"alertmanager.config.receivers[0].name":                                                "null",
		"alertmanager.config.receivers[1].name":                                                "discord_prometheus",
		"alertmanager.config.receivers[1].webhook_configs[0].url":                              "https://example.com/webhook/prometheus",
		"alertmanager.config.receivers[1].webhook_configs[0].send_resolved":                    "true",
		"alertmanager.config.receivers[1].webhook_configs[0].http_config.headers.Content-Type": "application/json",
		"alertmanager.config.receivers[1].webhook_configs[0].body":                             "{\n  \"content\": \"**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}\"\n}",
		"alertmanager.config.receivers[2].name":                                                "discord_alerts",
		"alertmanager.config.receivers[2].webhook_configs[0].url":                              "https://example.com/webhook/alerts",
		"alertmanager.config.receivers[2].webhook_configs[0].send_resolved":                    "true",
		"alertmanager.config.receivers[2].webhook_configs[0].http_config.headers.Content-Type": "application/json",
		"alertmanager.config.receivers[2].webhook_configs[0].body":                             "{\n  \"content\": \"**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}\"\n}",
	}

	var yamlData interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &yamlData)
	if err != nil {
		t.Fatalf("Failed to parse test YAML: %v", err)
	}

	result, err := flattener.FlattenYAML(yamlData)
	if err != nil {
		t.Errorf("FlattenYAML() error = %v", err)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenYAML() result mismatch")
		t.Logf("Expected:")
		for k, v := range expected {
			t.Logf("  %q: %q", k, v)
		}
		t.Logf("Got:")
		for k, v := range result {
			t.Logf("  %q: %q", k, v)
		}
	}
}

func TestFlattenYAML_AlertmanagerConfigWithEscapedNewlines(t *testing.T) {
	flattener := NewFlattenerWithOptions(true) // Enable escape newlines

	yamlStr := `alertmanager:
    config:
        receivers:
            - name: "null"
            - name: discord_prometheus
              webhook_configs:
                - url: https://example.com/webhook/prometheus
                  send_resolved: true
                  http_config:
                    headers:
                        Content-Type: application/json
                  body: |-
                    {
                      "content": "**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
                    }
            - name: discord_alerts
              webhook_configs:
                - url: https://example.com/webhook/alerts
                  send_resolved: true
                  http_config:
                    headers:
                        Content-Type: application/json
                  body: |-
                    {
                      "content": "**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
                    }`

	expected := map[string]string{
		"alertmanager.config.receivers[0].name":                                                "null",
		"alertmanager.config.receivers[1].name":                                                "discord_prometheus",
		"alertmanager.config.receivers[1].webhook_configs[0].url":                              "https://example.com/webhook/prometheus",
		"alertmanager.config.receivers[1].webhook_configs[0].send_resolved":                    "true",
		"alertmanager.config.receivers[1].webhook_configs[0].http_config.headers.Content-Type": "application/json",
		"alertmanager.config.receivers[1].webhook_configs[0].body":                             "{\\n  \"content\": \"**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}\"\\n}",
		"alertmanager.config.receivers[2].name":                                                "discord_alerts",
		"alertmanager.config.receivers[2].webhook_configs[0].url":                              "https://example.com/webhook/alerts",
		"alertmanager.config.receivers[2].webhook_configs[0].send_resolved":                    "true",
		"alertmanager.config.receivers[2].webhook_configs[0].http_config.headers.Content-Type": "application/json",
		"alertmanager.config.receivers[2].webhook_configs[0].body":                             "{\\n  \"content\": \"**{{ .Status | title }}**: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}\"\\n}",
	}

	var yamlData interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &yamlData)
	if err != nil {
		t.Fatalf("Failed to parse test YAML: %v", err)
	}

	result, err := flattener.FlattenYAML(yamlData)
	if err != nil {
		t.Errorf("FlattenYAML() error = %v", err)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlattenYAML() result mismatch")
		t.Logf("Expected:")
		for k, v := range expected {
			t.Logf("  %q: %q", k, v)
		}
		t.Logf("Got:")
		for k, v := range result {
			t.Logf("  %q: %q", k, v)
		}
	}
}
