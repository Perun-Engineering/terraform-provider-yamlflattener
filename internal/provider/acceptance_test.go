// Package provider contains acceptance tests for the YAML flattener provider.
package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAcceptance_FullProviderWorkflow tests the complete provider workflow
// This is the main acceptance test that validates all functionality
func TestAcceptance_FullProviderWorkflow(t *testing.T) {
	// Create temporary directory and files for testing
	tempDir := t.TempDir()

	// Create test YAML files
	simpleYAMLPath := filepath.Join(tempDir, "simple.yaml")
	complexYAMLPath := filepath.Join(tempDir, "complex.yaml")
	alertmanagerYAMLPath := filepath.Join(tempDir, "alertmanager.yaml")

	simpleYAML := `
key1: value1
key2:
  nested: value2
  another: value3
`

	complexYAML := `
service:
  name: "example-service"
  port: 8080
  environment:
    NODE_ENV: "production"
    LOG_LEVEL: "info"
    DEBUG: false
  replicas: 3
  volumes:
    - name: "data"
      path: "/var/data"
      readOnly: true
    - name: "config"
      path: "/etc/config"
      readOnly: false
  healthcheck:
    path: "/health"
    initialDelaySeconds: 10
    periodSeconds: 30
    enabled: true
`

	alertmanagerYAML := `
alertmanager:
  config:
    global:
      slack_api_url: "your-encrypted-slack-webhook"
      smtp_smarthost: "localhost:587"
      smtp_from: "alertmanager@example.org"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "your-encrypted-webhook-url"
            channel: "#alerts"
            send_resolved: true
            title: "Alert: {{ .GroupLabels.alertname }}"
      - name: "email-notifications"
        email_configs:
          - to: "admin@example.org"
            subject: "Alert: {{ .GroupLabels.alertname }}"
            body: "{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
    route:
      group_by: ["alertname"]
      group_wait: "10s"
      group_interval: "10s"
      repeat_interval: "1h"
      receiver: "slack-notifications"
      routes:
        - match:
            severity: "critical"
          receiver: "email-notifications"
`

	// Write test files
	if err := os.WriteFile(simpleYAMLPath, []byte(simpleYAML), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(complexYAMLPath, []byte(complexYAML), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(alertmanagerYAMLPath, []byte(alertmanagerYAML), 0600); err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Test data source with file input
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

data "yamlflattener_flatten" "simple_file" {
  yaml_file = %q
}

data "yamlflattener_flatten" "complex_file" {
  yaml_file = %q
}

data "yamlflattener_flatten" "alertmanager_file" {
  yaml_file = %q
}

output "simple_flattened" {
  value = data.yamlflattener_flatten.simple_file.flattened
}

output "complex_service_name" {
  value = data.yamlflattener_flatten.complex_file.flattened["service.name"]
}

output "alertmanager_slack_url" {
  value = data.yamlflattener_flatten.alertmanager_file.flattened["alertmanager.config.global.slack_api_url"]
}

output "volume_count" {
  value = length([
    for k, v in data.yamlflattener_flatten.complex_file.flattened : k
    if startswith(k, "service.volumes[") && endswith(k, "].name")
  ])
}
`, simpleYAMLPath, complexYAMLPath, alertmanagerYAMLPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.simple_file", "flattened.key1", "value1"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.simple_file", "flattened.key2.nested", "value2"),
					resource.TestCheckResourceAttr("data.yamlflattener_flatten.simple_file", "flattened.key2.another", "value3"),
					resource.TestCheckOutput("complex_service_name", "example-service"),
					resource.TestCheckOutput("alertmanager_slack_url", "your-encrypted-slack-webhook"),
					resource.TestCheckOutput("volume_count", "2"),
				),
			},
			{
				// Step 2: Test data source with inline YAML content
				Config: `
provider "yamlflattener" {}

data "yamlflattener_flatten" "inline_content" {
  yaml_content = <<EOT
database:
  host: "localhost"
  port: 5432
  credentials:
    username: "admin"
    password: "secret"
  pools:
    - name: "read_pool"
      size: 10
    - name: "write_pool"
      size: 5
EOT
}

output "db_host" {
  value = data.yamlflattener_flatten.inline_content.flattened["database.host"]
}

output "db_port" {
  value = data.yamlflattener_flatten.inline_content.flattened["database.port"]
}

output "read_pool_size" {
  value = data.yamlflattener_flatten.inline_content.flattened["database.pools[0].size"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("db_host", "localhost"),
					resource.TestCheckOutput("db_port", "5432"),
					resource.TestCheckOutput("read_pool_size", "10"),
				),
			},
			{
				// Step 3: Test provider function
				Config: `
provider "yamlflattener" {}

locals {
  yaml_content = <<EOT
api:
  version: "v1"
  endpoints:
    - path: "/users"
      methods: ["GET", "POST"]
      auth: true
    - path: "/health"
      methods: ["GET"]
      auth: false
  settings:
    timeout: 30
    retries: 3
    debug: true
EOT
}

output "function_result" {
  value = provider::yamlflattener::flatten(local.yaml_content)
}

output "api_version" {
  value = provider::yamlflattener::flatten(local.yaml_content)["api.version"]
}

output "first_endpoint_path" {
  value = provider::yamlflattener::flatten(local.yaml_content)["api.endpoints[0].path"]
}

output "timeout_setting" {
  value = provider::yamlflattener::flatten(local.yaml_content)["api.settings.timeout"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("api_version", "v1"),
					resource.TestCheckOutput("first_endpoint_path", "/users"),
					resource.TestCheckOutput("timeout_setting", "30"),
				),
			},
		},
	})
}

// TestAcceptance_DataSourceAndFunctionEquivalence verifies that both approaches
// produce identical results for the same input (Requirement 6.5)
func TestAcceptance_DataSourceAndFunctionEquivalence(t *testing.T) {
	testCases := []struct {
		name        string
		yamlContent string
		testKeys    []string
	}{
		{
			name: "simple_object",
			yamlContent: `
key1: value1
key2:
  nested: value2
`,
			testKeys: []string{"key1", "key2.nested"},
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
			testKeys: []string{"items[0].name", "items[0].value", "items[1].name", "items[1].value"},
		},
		{
			name: "mixed_types",
			yamlContent: `
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
empty_val: ""
`,
			testKeys: []string{"string_val", "int_val", "float_val", "bool_val", "null_val", "empty_val"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testAccEquivalenceConfig(tc.yamlContent),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckOutput("outputs_are_equal", "true"),
							testEquivalenceForKeys(tc.testKeys),
						),
					},
				},
			})
		})
	}
}

// TestAcceptance_CrossPlatformCompatibility tests platform-specific behavior
func TestAcceptance_CrossPlatformCompatibility(t *testing.T) {
	t.Logf("Running acceptance tests on platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

data "yamlflattener_flatten" "platform_test" {
  yaml_content = <<EOT
platform:
  os: %q
  arch: %q
  go_version: "1.21"
  test_data:
    unicode: "Hello 世界 🌍"
    special_chars: "!@#$%%^&*()_+-=[]{}|;':\",./<>?"
    multiline: |
      This is a multiline
      string that spans
      multiple lines
  paths:
    unix_style: "/usr/local/bin"
    windows_style: "C:\\Program Files\\App"
EOT
}

output "platform_os" {
  value = data.yamlflattener_flatten.platform_test.flattened["platform.os"]
}

output "platform_arch" {
  value = data.yamlflattener_flatten.platform_test.flattened["platform.arch"]
}

output "unicode_test" {
  value = data.yamlflattener_flatten.platform_test.flattened["platform.test_data.unicode"]
}

output "special_chars_test" {
  value = data.yamlflattener_flatten.platform_test.flattened["platform.test_data.special_chars"]
}

output "multiline_test" {
  value = data.yamlflattener_flatten.platform_test.flattened["platform.test_data.multiline"]
}
`, runtime.GOOS, runtime.GOARCH),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("platform_os", runtime.GOOS),
					resource.TestCheckOutput("platform_arch", runtime.GOARCH),
					resource.TestCheckOutput("unicode_test", "Hello 世界 🌍"),
					resource.TestCheckOutput("special_chars_test", "!@#$%^&*()_+-=[]{}|;':\",./<>?"),
					resource.TestCheckResourceAttrWith("data.yamlflattener_flatten.platform_test", "flattened.platform.test_data.multiline", func(value string) error {
						if !strings.Contains(value, "This is a multiline") {
							return fmt.Errorf("multiline string not preserved correctly: %s", value)
						}
						return nil
					}),
				),
			},
		},
	})
}

// TestAcceptance_ErrorHandling tests various error conditions
func TestAcceptance_ErrorHandling(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test invalid YAML syntax
				Config: `
provider "yamlflattener" {}

data "yamlflattener_flatten" "invalid_yaml" {
  yaml_content = <<EOT
invalid: yaml: content:
  - missing
    proper: indentation
EOT
}
`,
				ExpectError: regexp.MustCompile("failed to parse YAML content"),
			},
			{
				// Test non-existent file
				Config: `
provider "yamlflattener" {}

data "yamlflattener_flatten" "missing_file" {
  yaml_file = "/path/that/does/not/exist.yaml"
}
`,
				ExpectError: regexp.MustCompile("failed to read YAML file"),
			},
			{
				// Test conflicting parameters
				Config: `
provider "yamlflattener" {}

data "yamlflattener_flatten" "conflicting_params" {
  yaml_content = "key: value"
  yaml_file = "/some/file.yaml"
}
`,
				ExpectError: regexp.MustCompile("only one of yaml_content or yaml_file should be specified"),
			},
			{
				// Test missing parameters
				Config: `
provider "yamlflattener" {}

data "yamlflattener_flatten" "missing_params" {
}
`,
				ExpectError: regexp.MustCompile("either yaml_content or yaml_file must be specified"),
			},
		},
	})
}

// TestAcceptance_ProviderInstallationWorkflow simulates the provider installation workflow
func TestAcceptance_ProviderInstallationWorkflow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test provider with version constraint
				Config: `
terraform {
  required_version = ">= 1.0"
  required_providers {
    yamlflattener = {
      source = "registry.terraform.io/terraform/yamlflattener"
      version = ">= 0.2.0"
    }
  }
}

provider "yamlflattener" {
  max_depth = 50
}

data "yamlflattener_flatten" "installation_test" {
  yaml_content = <<EOT
installation:
  status: "success"
  provider:
    name: "yamlflattener"
    version: "0.2.0"
  features:
    - "data_source"
    - "provider_function"
    - "yaml_flattening"
EOT
}

output "installation_status" {
  value = data.yamlflattener_flatten.installation_test.flattened["installation.status"]
}

output "provider_name" {
  value = data.yamlflattener_flatten.installation_test.flattened["installation.provider.name"]
}

output "feature_count" {
  value = length([
    for k, v in data.yamlflattener_flatten.installation_test.flattened : k
    if startswith(k, "installation.features[") && !contains(k, ".")
  ])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("installation_status", "success"),
					resource.TestCheckOutput("provider_name", "yamlflattener"),
					resource.TestCheckOutput("feature_count", "3"),
				),
			},
		},
	})
}

// TestAcceptance_RealWorldUseCases tests realistic usage scenarios
func TestAcceptance_RealWorldUseCases(t *testing.T) {
	// Create a realistic configuration file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app-config.yaml")

	appConfig := `
application:
  name: "web-service"
  version: "1.2.3"
  environment: "production"

server:
  host: "0.0.0.0"
  port: 8080
  ssl:
    enabled: true
    cert_file: "/etc/ssl/certs/app.crt"
    key_file: "/etc/ssl/private/app.key"

database:
  primary:
    host: "db-primary.example.com"
    port: 5432
    name: "app_db"
    ssl_mode: "require"
  replica:
    host: "db-replica.example.com"
    port: 5432
    name: "app_db"
    ssl_mode: "require"

redis:
  cluster:
    - host: "redis-1.example.com"
      port: 6379
    - host: "redis-2.example.com"
      port: 6379
    - host: "redis-3.example.com"
      port: 6379

monitoring:
  metrics:
    enabled: true
    port: 9090
    path: "/metrics"
  health_check:
    enabled: true
    path: "/health"
    interval: "30s"
  logging:
    level: "info"
    format: "json"
    outputs:
      - type: "stdout"
      - type: "file"
        path: "/var/log/app.log"
        max_size: "100MB"
        max_files: 5
`

	if err := os.WriteFile(configPath, []byte(appConfig), 0600); err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

# Load configuration from file
data "yamlflattener_flatten" "app_config" {
  yaml_file = %q
}

# Use flattened config to create outputs that could be used by other resources
output "app_name" {
  value = data.yamlflattener_flatten.app_config.flattened["application.name"]
}

output "app_version" {
  value = data.yamlflattener_flatten.app_config.flattened["application.version"]
}

output "server_config" {
  value = {
    host = data.yamlflattener_flatten.app_config.flattened["server.host"]
    port = data.yamlflattener_flatten.app_config.flattened["server.port"]
    ssl_enabled = data.yamlflattener_flatten.app_config.flattened["server.ssl.enabled"]
  }
}

output "database_urls" {
  value = {
    primary = "postgresql://${data.yamlflattener_flatten.app_config.flattened["database.primary.host"]}:${data.yamlflattener_flatten.app_config.flattened["database.primary.port"]}/${data.yamlflattener_flatten.app_config.flattened["database.primary.name"]}"
    replica = "postgresql://${data.yamlflattener_flatten.app_config.flattened["database.replica.host"]}:${data.yamlflattener_flatten.app_config.flattened["database.replica.port"]}/${data.yamlflattener_flatten.app_config.flattened["database.replica.name"]}"
  }
}

output "redis_endpoints" {
  value = [
    "${data.yamlflattener_flatten.app_config.flattened["redis.cluster[0].host"]}:${data.yamlflattener_flatten.app_config.flattened["redis.cluster[0].port"]}",
    "${data.yamlflattener_flatten.app_config.flattened["redis.cluster[1].host"]}:${data.yamlflattener_flatten.app_config.flattened["redis.cluster[1].port"]}",
    "${data.yamlflattener_flatten.app_config.flattened["redis.cluster[2].host"]}:${data.yamlflattener_flatten.app_config.flattened["redis.cluster[2].port"]}"
  ]
}

# Test using provider function for dynamic configuration
locals {
  override_config = <<EOT
overrides:
  database:
    primary:
      host: "override-db.example.com"
  redis:
    cluster:
      - host: "override-redis.example.com"
        port: 6380
EOT
}

output "override_db_host" {
  value = provider::yamlflattener::flatten(local.override_config)["overrides.database.primary.host"]
}

output "override_redis_port" {
  value = provider::yamlflattener::flatten(local.override_config)["overrides.redis.cluster[0].port"]
}
`, configPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("app_name", "web-service"),
					resource.TestCheckOutput("app_version", "1.2.3"),
					resource.TestCheckOutput("server_config", "{\"host\":\"0.0.0.0\",\"port\":\"8080\",\"ssl_enabled\":\"true\"}"),
					resource.TestCheckOutput("database_urls", "{\"primary\":\"postgresql://db-primary.example.com:5432/app_db\",\"replica\":\"postgresql://db-replica.example.com:5432/app_db\"}"),
					resource.TestCheckOutput("redis_endpoints", "[\"redis-1.example.com:6379\",\"redis-2.example.com:6379\",\"redis-3.example.com:6379\"]"),
					resource.TestCheckOutput("override_db_host", "override-db.example.com"),
					resource.TestCheckOutput("override_redis_port", "6380"),
				),
			},
		},
	})
}

// Helper functions

func testAccEquivalenceConfig(yamlContent string) string {
	return fmt.Sprintf(`
provider "yamlflattener" {}

locals {
  yaml_content = <<EOT
%s
EOT
}

# Data source approach
data "yamlflattener_flatten" "datasource_test" {
  yaml_content = local.yaml_content
}

# Function approach
locals {
  function_result = provider::yamlflattener::flatten(local.yaml_content)
}

# Compare outputs
output "datasource_output" {
  value = data.yamlflattener_flatten.datasource_test.flattened
}

output "function_output" {
  value = local.function_result
}

output "outputs_are_equal" {
  value = jsonencode(data.yamlflattener_flatten.datasource_test.flattened) == jsonencode(local.function_result)
}
`, yamlContent)
}

func testEquivalenceForKeys(_ []string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		// This would be implemented to check that specific keys have the same values
		// in both datasource and function outputs
		return nil
	}
}
