package flattener

import (
	"reflect"
	"testing"
)

// TestOrderPreservation tests that the flattener preserves the original YAML key order
func TestOrderPreservation(t *testing.T) {
	flattener := NewFlattener()

	tests := []struct {
		name     string
		yamlStr  string
		expected map[string]string
	}{
		{
			name: "Alertmanager receivers preserve original order",
			yamlStr: `alertmanager:
    config:
        receivers:
            - name: blackhole
            - name: discord_prometheus
              discord_config:
                - webhook_url: https://example.com/webhook1
            - name: discord_alerts
              discord_config:
                - webhook_url: https://example.com/webhook2`,
			expected: map[string]string{
				"alertmanager.config.receivers[0].name":                          "blackhole",
				"alertmanager.config.receivers[1].name":                          "discord_prometheus",
				"alertmanager.config.receivers[1].discord_config[0].webhook_url": "https://example.com/webhook1",
				"alertmanager.config.receivers[2].name":                          "discord_alerts",
				"alertmanager.config.receivers[2].discord_config[0].webhook_url": "https://example.com/webhook2",
			},
		},
		{
			name: "Simple object preserves key order",
			yamlStr: `name: discord_prometheus
discord_config:
  - webhook_url: test_url`,
			expected: map[string]string{
				"name":                          "discord_prometheus",
				"discord_config[0].webhook_url": "test_url",
			},
		},
		{
			name: "Mixed order keys preserved",
			yamlStr: `zebra: last
alpha: first
beta: second`,
			expected: map[string]string{
				"zebra": "last",
				"alpha": "first",
				"beta":  "second",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := flattener.FlattenYAMLString(tt.yamlStr)
			if err != nil {
				t.Fatalf("FlattenYAMLString() error = %v", err)
			}

			if !reflect.DeepEqual(result.ToMap(), tt.expected) {
				t.Errorf("FlattenYAMLString() = %v, want %v", result.ToMap(), tt.expected)
			}

			// Run multiple times to ensure consistency
			for i := 0; i < 3; i++ {
				result2, err := flattener.FlattenYAMLString(tt.yamlStr)
				if err != nil {
					t.Fatalf("Run %d: FlattenYAMLString() error = %v", i+1, err)
				}

				if !reflect.DeepEqual(result2.Keys(), result.Keys()) || !reflect.DeepEqual(result2.ToMap(), result.ToMap()) {
					t.Errorf("Run %d: Results are not consistent", i+1)
				}
			}
		})
	}
}

// TestOrderPreservationConsistency ensures that the order is preserved consistently across runs
func TestOrderPreservationConsistency(t *testing.T) {
	flattener := NewFlattener()

	yamlStr := `name: discord_prometheus
discord_config:
  - webhook_url: test_url
other_field: value`

	// Run multiple times and ensure the same result
	var firstResult *OrderedMap
	for i := 0; i < 5; i++ {
		result, err := flattener.FlattenYAMLString(yamlStr)
		if err != nil {
			t.Fatalf("Run %d: FlattenYAMLString() error = %v", i+1, err)
		}

		if i == 0 {
			firstResult = result
		} else {
			// Compare the keys and values
			if !reflect.DeepEqual(result.Keys(), firstResult.Keys()) {
				t.Errorf("Run %d: Key order is not consistent with first run", i+1)
				t.Errorf("First result keys: %v", firstResult.Keys())
				t.Errorf("Current result keys: %v", result.Keys())
			}
			if !reflect.DeepEqual(result.ToMap(), firstResult.ToMap()) {
				t.Errorf("Run %d: Values are not consistent with first run", i+1)
				t.Errorf("First result: %v", firstResult.ToMap())
				t.Errorf("Current result: %v", result.ToMap())
			}
		}
	}
}
