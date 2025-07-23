package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestIntegration_DataSourceAndFunctionEquivalence verifies that both the data source and function
// produce identical results for the same input
func TestIntegration_DataSourceAndFunctionEquivalence(t *testing.T) {
	// Test with various YAML structures to ensure consistent behavior
	testCases := []struct {
		name        string
		yamlContent string
	}{
		{
			name: "simple_object",
			yamlContent: `
key1: value1
key2:
  nested: value2
`,
		},
		{
			name: "array_structure",
			yamlContent: `
items:
  - name: item1
    value: val1
  - name: item2
    value: val2
`,
		},
		{
			name: "complex_nested",
			yamlContent: `
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
`,
		},
		{
			name: "mixed_types",
			yamlContent: `
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
array:
  - 1
  - "two"
  - true
  - null
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testAccIntegrationConfig_Equivalence(tc.yamlContent),
						Check: resource.ComposeAggregateTestCheckFunc(
							// Check that data source and function outputs are identical
							resource.TestCheckOutput("are_equal", "true"),
						),
					},
				},
			})
		})
	}
}

// TestIntegration_CrossPlatformCompatibility tests the provider on the current platform
// and verifies basic functionality works as expected
func TestIntegration_CrossPlatformCompatibility(t *testing.T) {
	t.Logf("Running on platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	// This test will run on whatever platform the tests are executed on
	// In a CI environment, this should be run on multiple platforms
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

# Platform-specific test
data "yamlflattener_flatten" "platform_test" {
  yaml_content = <<EOT
platform:
  os: %q
  arch: %q
  test: "Platform compatibility test"
EOT
}

output "platform_info" {
  value = data.yamlflattener_flatten.platform_test.flattened
}
`, runtime.GOOS, runtime.GOARCH),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.platform_test", "flattened.platform.os", runtime.GOOS),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.platform_test", "flattened.platform.arch", runtime.GOARCH),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.platform_test", "flattened.platform.test", "Platform compatibility test"),
				),
			},
		},
	})
}

// TestIntegration_CompleteWorkflow tests a complete workflow that simulates real usage
func TestIntegration_CompleteWorkflow(t *testing.T) {
	// Create a temporary YAML file for testing
	tempDir := t.TempDir()
	yamlFilePath := filepath.Join(tempDir, "config.yaml")
	yamlContent := `
service:
  name: "example-service"
  port: 8080
  environment:
    NODE_ENV: "production"
    LOG_LEVEL: "info"
  replicas: 3
  volumes:
    - name: "data"
      path: "/var/data"
    - name: "config"
      path: "/etc/config"
  healthcheck:
    path: "/health"
    initialDelaySeconds: 10
    periodSeconds: 30
`
	err := os.WriteFile(yamlFilePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Use data source with file input
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

data "yamlflattener_flatten" "config_file" {
  yaml_file = %q
}

output "service_name" {
  value = data.yamlflattener_flatten.config_file.flattened["service.name"]
}

output "service_port" {
  value = data.yamlflattener_flatten.config_file.flattened["service.port"]
}

output "environment_vars" {
  value = {
    node_env = data.yamlflattener_flatten.config_file.flattened["service.environment.NODE_ENV"]
    log_level = data.yamlflattener_flatten.config_file.flattened["service.environment.LOG_LEVEL"]
  }
}

output "volume_paths" {
  value = [
    data.yamlflattener_flatten.config_file.flattened["service.volumes[0].path"],
    data.yamlflattener_flatten.config_file.flattened["service.volumes[1].path"]
  ]
}
`, yamlFilePath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("service_name", "example-service"),
					resource.TestCheckOutput("service_port", "8080"),
					resource.TestCheckOutput("environment_vars", "{\"log_level\":\"info\",\"node_env\":\"production\"}"),
					resource.TestCheckOutput("volume_paths", "[\""+"/var/data"+"\",\""+"/etc/config"+"\"]"),
				),
			},
			{
				// Step 2: Use provider function with inline YAML
				Config: `
provider "yamlflattener" {}

locals {
  yaml_content = <<EOT
service:
  name: "example-service"
  port: 8080
  environment:
    NODE_ENV: "production"
    LOG_LEVEL: "info"
EOT
}

output "flattened_map" {
  value = provider::yamlflattener::flatten(local.yaml_content)
}

output "service_name_from_function" {
  value = provider::yamlflattener::flatten(local.yaml_content)["service.name"]
}

output "combined_output" {
  value = {
    name = provider::yamlflattener::flatten(local.yaml_content)["service.name"]
    port = provider::yamlflattener::flatten(local.yaml_content)["service.port"]
    env = {
      node_env = provider::yamlflattener::flatten(local.yaml_content)["service.environment.NODE_ENV"]
      log_level = provider::yamlflattener::flatten(local.yaml_content)["service.environment.LOG_LEVEL"]
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("service_name_from_function", "example-service"),
					resource.TestCheckOutput("combined_output", "{\"env\":{\"log_level\":\"info\",\"node_env\":\"production\"},\"name\":\"example-service\",\"port\":\"8080\"}"),
				),
			},
			{
				// Step 3: Test provider configuration options
				Config: `
provider "yamlflattener" {
  max_depth = 20
}

locals {
  deep_yaml = <<EOT
level1:
  level2:
    level3:
      level4:
        level5:
          value: "deep value"
EOT
}

output "deep_value" {
  value = provider::yamlflattener::flatten(local.deep_yaml)["level1.level2.level3.level4.level5.value"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("deep_value", "deep value"),
				),
			},
		},
	})
}

// TestIntegration_InstallationWorkflow simulates the provider installation and usage workflow
func TestIntegration_InstallationWorkflow(t *testing.T) {
	// This test simulates what happens when a user installs and uses the provider
	// We can't actually test the installation process in a unit test, but we can
	// verify the provider works as expected once "installed"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Basic provider configuration
				Config: `
terraform {
  required_providers {
    yamlflattener = {
      source = "example/yamlflattener"
      version = "0.1.0"
    }
  }
}

provider "yamlflattener" {}

data "yamlflattener_flatten" "test" {
  yaml_content = "key: value"
}

output "result" {
  value = data.yamlflattener_flatten.test.flattened["key"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("result", "value"),
				),
			},
		},
	})
}

// Helper function to create a Terraform configuration that tests equivalence
// between data source and function outputs
func testAccIntegrationConfig_Equivalence(yamlContent string) string {
	return fmt.Sprintf(`
provider "yamlflattener" {}

# Use the same YAML content for both data source and function
locals {
  yaml_content = <<EOT
%s
EOT
}

# Data source approach
data "yamlflattener_flatten" "test_datasource" {
  yaml_content = local.yaml_content
}

# Function approach
locals {
  function_result = provider::yamlflattener::flatten(local.yaml_content)
}

# Compare the outputs
output "datasource_output" {
  value = data.yamlflattener_flatten.test_datasource.flattened
}

output "function_output" {
  value = local.function_result
}

# Check if they're equal
output "are_equal" {
  value = jsonencode(data.yamlflattener_flatten.test_datasource.flattened) == jsonencode(local.function_result)
}
`, yamlContent)
}
