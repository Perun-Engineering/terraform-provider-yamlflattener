package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Use the testAccProtoV6ProviderFactories from provider_test.go

func TestAccFlattenDataSource_YAMLContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test basic YAML content
			{
				Config: testAccFlattenDataSourceConfig_YAMLContent(`
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key1", "value1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key2.nested", "value2"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[0].name", "item1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[1].name", "item2"),
				),
			},
			// Test complex YAML content with mixed types
			{
				Config: testAccFlattenDataSourceConfig_YAMLContent(`
alertmanager:
  config:
    global:
      slack_api_url: "your-encrypted-slack-webhook"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "your-encrypted-webhook-url"
            channel: "#alerts"
            send_resolved: true
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.alertmanager.config.global.slack_api_url", "your-encrypted-slack-webhook"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.alertmanager.config.receivers[0].name", "slack-notifications"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.alertmanager.config.receivers[0].slack_configs[0].api_url", "your-encrypted-webhook-url"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.alertmanager.config.receivers[0].slack_configs[0].channel", "#alerts"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.alertmanager.config.receivers[0].slack_configs[0].send_resolved", "true"),
				),
			},
		},
	})
}

func TestAccFlattenDataSource_YAMLFile(t *testing.T) {
	// Create a temporary YAML file for testing
	tempDir := t.TempDir()
	yamlFilePath := filepath.Join(tempDir, "test.yaml")
	yamlContent := `
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
`
	err := os.WriteFile(yamlFilePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFlattenDataSourceConfig_YAMLFile(yamlFilePath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key1", "value1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.key2.nested", "value2"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[0].name", "item1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.test", "flattened.array[1].name", "item2"),
				),
			},
		},
	})
}

func TestAccFlattenDataSource_ErrorHandling(t *testing.T) {
	// Test with invalid YAML content
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfig_YAMLContent(`invalid: yaml: : content`),
				ExpectError: regexp.MustCompile(`Failed to Flatten YAML Content`),
			},
		},
	})

	// Test with non-existent file
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfig_YAMLFile("/path/to/nonexistent/file.yaml"),
				ExpectError: regexp.MustCompile(`Failed to Flatten YAML File`),
			},
		},
	})

	// Test with both yaml_content and yaml_file provided
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfig_BothInputs(`key: value`, "/path/to/file.yaml"),
				ExpectError: regexp.MustCompile(`Conflicting Inputs`),
			},
		},
	})

	// Test with neither yaml_content nor yaml_file provided
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFlattenDataSourceConfig_NoInputs(),
				ExpectError: regexp.MustCompile(`Missing Required Input`),
			},
		},
	})
}

func testAccFlattenDataSourceConfig_YAMLContent(yamlContent string) string {
	return fmt.Sprintf(`
data "yamlflattener_flatten" "test" {
  yaml_content = <<EOT
%s
EOT
}
`, yamlContent)
}

func testAccFlattenDataSourceConfig_YAMLFile(filePath string) string {
	return fmt.Sprintf(`
data "yamlflattener_flatten" "test" {
  yaml_file = %q
}
`, filePath)
}

func testAccFlattenDataSourceConfig_BothInputs(yamlContent, filePath string) string {
	return fmt.Sprintf(`
data "yamlflattener_flatten" "test" {
  yaml_content = <<EOT
%s
EOT
  yaml_file = %q
}
`, yamlContent, filePath)
}

func testAccFlattenDataSourceConfig_NoInputs() string {
	return `
data "yamlflattener_flatten" "test" {
}
`
}
