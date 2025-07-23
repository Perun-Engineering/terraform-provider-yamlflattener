package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestFinalIntegration_ComprehensiveEquivalence performs a comprehensive test
// to verify that both data source and function produce identical results for
// complex YAML structures (Requirement 6.5)
func TestFinalIntegration_ComprehensiveEquivalence(t *testing.T) {
	// Create a complex YAML structure with various data types and nesting levels
	complexYAML := `
# Complex YAML with various data types and structures
application:
  name: "terraform-yaml-flattener"
  version: "1.0.0"
  description: "A Terraform provider for flattening YAML structures"
  metadata:
    created_at: "2023-07-22T10:30:00Z"
    updated_at: "2023-07-22T15:45:00Z"
    tags:
      - "terraform"
      - "yaml"
      - "provider"
    maintainers:
      - name: "Developer One"
        email: "dev1@example.com"
        roles:
          - "lead"
          - "maintainer"
      - name: "Developer Two"
        email: "dev2@example.com"
        roles:
          - "contributor"
  settings:
    debug: true
    log_level: "info"
    timeout: 30
    retry:
      enabled: true
      max_attempts: 3
      backoff:
        initial: 1.0
        multiplier: 2.0
        max: 60.0
  features:
    - name: "yaml_parsing"
      enabled: true
      config:
        strict_mode: false
    - name: "flattening"
      enabled: true
      config:
        preserve_types: false
        max_depth: 100
    - name: "terraform_integration"
      enabled: true
      config:
        provider_function: true
        data_source: true
  environments:
    development:
      url: "http://localhost:8080"
      debug: true
    staging:
      url: "https://staging.example.com"
      debug: true
    production:
      url: "https://production.example.com"
      debug: false
  empty_array: []
  null_value: null
  empty_string: ""
  zero_value: 0
  false_value: false
`

	// Define specific paths to check for equivalence
	pathsToCheck := []string{
		"application.name",
		"application.version",
		"application.metadata.tags[0]",
		"application.metadata.tags[1]",
		"application.metadata.maintainers[0].name",
		"application.metadata.maintainers[0].roles[0]",
		"application.metadata.maintainers[1].email",
		"application.settings.retry.backoff.initial",
		"application.features[0].config.strict_mode",
		"application.features[1].name",
		"application.environments.production.url",
		"application.environments.production.debug",
		"application.null_value",
		"application.empty_string",
		"application.zero_value",
		"application.false_value",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testFinalIntegrationConfig_Equivalence(complexYAML),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check overall equivalence
					resource.TestCheckOutput("outputs_are_equal", "true"),

					// Check specific paths for equivalence
					testFinalIntegration_CheckPathEquivalence(pathsToCheck),

					// Verify specific values from both approaches
					resource.TestCheckOutput("ds_app_name", "terraform-yaml-flattener"),
					resource.TestCheckOutput("fn_app_name", "terraform-yaml-flattener"),
					resource.TestCheckOutput("ds_first_tag", "terraform"),
					resource.TestCheckOutput("fn_first_tag", "terraform"),
					resource.TestCheckOutput("ds_prod_debug", "false"),
					resource.TestCheckOutput("fn_prod_debug", "false"),
				),
			},
		},
	})
}

// TestFinalIntegration_CrossPlatformFeatures tests platform-specific features
// and ensures the provider works correctly across different platforms
func TestFinalIntegration_CrossPlatformFeatures(t *testing.T) {
	// Create a temporary directory for platform-specific file paths
	tempDir := t.TempDir()

	// Create platform-specific paths
	var platformPath string
	if runtime.GOOS == "windows" {
		platformPath = filepath.Join(tempDir, "windows\\path\\file.yaml")
	} else {
		platformPath = filepath.Join(tempDir, "unix/path/file.yaml")
	}

	// Ensure directory exists
	err := os.MkdirAll(filepath.Dir(platformPath), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create test YAML file
	yamlContent := fmt.Sprintf(`
platform:
  os: %q
  arch: %q
  path_separator: %q
  line_endings: %q
  unicode_test: "Hello 世界 🌍"
  special_chars: "!@#$%%^&*()_+-=[]{}|;':\",./<>?"
`,
		runtime.GOOS,
		runtime.GOARCH,
		string(os.PathSeparator),
		"\r\n",
	)

	err = os.WriteFile(platformPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

# Test with file path specific to this platform
data "yamlflattener_flatten" "platform_file" {
  yaml_file = %q
}

# Test with inline content
data "yamlflattener_flatten" "platform_content" {
  yaml_content = <<EOT
%s
EOT
}

# Test function with the same content
locals {
  yaml_content = <<EOT
%s
EOT
  function_result = provider::yamlflattener::flatten(local.yaml_content)
}

output "file_os" {
  value = data.yamlflattener_flatten.platform_file.flattened["platform.os"]
}

output "content_os" {
  value = data.yamlflattener_flatten.platform_content.flattened["platform.os"]
}

output "function_os" {
  value = local.function_result["platform.os"]
}

output "unicode_test" {
  value = data.yamlflattener_flatten.platform_file.flattened["platform.unicode_test"]
}

output "special_chars" {
  value = data.yamlflattener_flatten.platform_file.flattened["platform.special_chars"]
}

output "all_equal" {
  value = (
    data.yamlflattener_flatten.platform_file.flattened["platform.os"] ==
    data.yamlflattener_flatten.platform_content.flattened["platform.os"] &&
    data.yamlflattener_flatten.platform_content.flattened["platform.os"] ==
    local.function_result["platform.os"]
  )
}
`, platformPath, yamlContent, yamlContent),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that OS is correctly detected
					resource.TestCheckOutput("file_os", runtime.GOOS),
					resource.TestCheckOutput("content_os", runtime.GOOS),
					resource.TestCheckOutput("function_os", runtime.GOOS),

					// Check that Unicode and special characters are preserved
					resource.TestCheckOutput("unicode_test", "Hello 世界 🌍"),
					resource.TestCheckOutput("special_chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?"),

					// Check that all methods produce the same result
					resource.TestCheckOutput("all_equal", "true"),
				),
			},
		},
	})
}

// TestFinalIntegration_ProviderInstallationWorkflow tests the complete provider
// installation and usage workflow
func TestFinalIntegration_ProviderInstallationWorkflow(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a test configuration file
	configPath := filepath.Join(tempDir, "terraform.tf")
	configContent := `
terraform {
  required_providers {
    yamlflattener = {
      source = "registry.terraform.io/terraform/yamlflattener"
      version = ">= 0.1.0"
    }
  }
}

provider "yamlflattener" {
  # Provider configuration options
  max_depth = 100
}

# Test data source
data "yamlflattener_flatten" "test" {
  yaml_content = <<EOT
installation:
  test: "success"
  provider: "yamlflattener"
  version: "0.1.0"
EOT
}

# Test provider function
output "function_test" {
  value = provider::yamlflattener::flatten(<<EOT
function:
  test: "success"
  provider: "yamlflattener"
  version: "0.1.0"
EOT
)
}

# Output test results
output "datasource_result" {
  value = data.yamlflattener_flatten.test.flattened["installation.test"]
}

output "function_result" {
  value = provider::yamlflattener::flatten(<<EOT
function:
  test: "success"
EOT
)["function.test"]
}
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// We can't actually test the installation process in a unit test,
	// but we can verify the provider configuration is valid
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "yamlflattener" {
  max_depth = 100
}

data "yamlflattener_flatten" "installation_test" {
  yaml_content = <<EOT
installation:
  test: "success"
  provider: "yamlflattener"
  version: "0.1.0"
EOT
}

output "installation_test" {
  value = data.yamlflattener_flatten.installation_test.flattened["installation.test"]
}

output "provider_name" {
  value = data.yamlflattener_flatten.installation_test.flattened["installation.provider"]
}

output "provider_version" {
  value = data.yamlflattener_flatten.installation_test.flattened["installation.version"]
}

# Test function as well
locals {
  function_yaml = <<EOT
function:
  test: "success"
EOT
}

output "function_test" {
  value = provider::yamlflattener::flatten(local.function_yaml)["function.test"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("installation_test", "success"),
					resource.TestCheckOutput("provider_name", "yamlflattener"),
					resource.TestCheckOutput("provider_version", "0.1.0"),
					resource.TestCheckOutput("function_test", "success"),
				),
			},
		},
	})
}

// TestFinalIntegration_RealWorldScenarios tests the provider in realistic scenarios
// that might be encountered in production environments
func TestFinalIntegration_RealWorldScenarios(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Create multiple YAML files that might be used in a real project
	files := map[string]string{
		"app-config.yaml": `
application:
  name: "my-service"
  version: "2.1.0"
  environment: "production"

server:
  host: "0.0.0.0"
  port: 8080
  ssl:
    enabled: true
    cert_file: "/etc/ssl/certs/app.crt"
    key_file: "/etc/ssl/private/app.key"
`,
		"database.yaml": `
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
`,
		"monitoring.yaml": `
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
`,
	}

	// Write all test files
	for filename, content := range files {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "yamlflattener" {}

# Load configuration from multiple files
data "yamlflattener_flatten" "app_config" {
  yaml_file = %q
}

data "yamlflattener_flatten" "db_config" {
  yaml_file = %q
}

data "yamlflattener_flatten" "monitoring_config" {
  yaml_file = %q
}

# Use provider function for dynamic configuration
locals {
  override_yaml = <<EOT
overrides:
  server:
    port: 9000
  database:
    primary:
      host: "custom-db.example.com"
  monitoring:
    logging:
      level: "debug"
EOT

  overrides = provider::yamlflattener::flatten(local.override_yaml)
}

# Combine configurations using both data sources and functions
locals {
  # Get base values from data sources
  base_server_port = data.yamlflattener_flatten.app_config.flattened["server.port"]
  base_db_host = data.yamlflattener_flatten.db_config.flattened["database.primary.host"]
  base_log_level = data.yamlflattener_flatten.monitoring_config.flattened["monitoring.logging.level"]

  # Get override values with defaults
  override_server_port = lookup(local.overrides, "overrides.server.port", local.base_server_port)
  override_db_host = lookup(local.overrides, "overrides.database.primary.host", local.base_db_host)
  override_log_level = lookup(local.overrides, "overrides.monitoring.logging.level", local.base_log_level)
}

# Output the final configuration
output "final_config" {
  value = {
    application = {
      name = data.yamlflattener_flatten.app_config.flattened["application.name"]
      version = data.yamlflattener_flatten.app_config.flattened["application.version"]
      environment = data.yamlflattener_flatten.app_config.flattened["application.environment"]
    }
    server = {
      host = data.yamlflattener_flatten.app_config.flattened["server.host"]
      port = local.override_server_port
      ssl_enabled = data.yamlflattener_flatten.app_config.flattened["server.ssl.enabled"]
    }
    database = {
      host = local.override_db_host
      port = data.yamlflattener_flatten.db_config.flattened["database.primary.port"]
      name = data.yamlflattener_flatten.db_config.flattened["database.primary.name"]
    }
    monitoring = {
      metrics_enabled = data.yamlflattener_flatten.monitoring_config.flattened["monitoring.metrics.enabled"]
      log_level = local.override_log_level
    }
  }
}

# Test that both data source and function produce the same results for the same input
locals {
  test_yaml = <<EOT
test:
  value: "same result"
EOT
}

data "yamlflattener_flatten" "test_ds" {
  yaml_content = local.test_yaml
}

output "datasource_result" {
  value = data.yamlflattener_flatten.test_ds.flattened["test.value"]
}

output "function_result" {
  value = provider::yamlflattener::flatten(local.test_yaml)["test.value"]
}

output "results_match" {
  value = data.yamlflattener_flatten.test_ds.flattened["test.value"] == provider::yamlflattener::flatten(local.test_yaml)["test.value"]
}
`,
					filepath.Join(tempDir, "app-config.yaml"),
					filepath.Join(tempDir, "database.yaml"),
					filepath.Join(tempDir, "monitoring.yaml")),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that data source and function produce the same results
					resource.TestCheckOutput("datasource_result", "same result"),
					resource.TestCheckOutput("function_result", "same result"),
					resource.TestCheckOutput("results_match", "true"),
				),
			},
		},
	})
}

// Helper functions

// testFinalIntegrationConfig_Equivalence creates a Terraform configuration that tests
// equivalence between data source and function outputs
func testFinalIntegrationConfig_Equivalence(yamlContent string) string {
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

# Test specific values from both approaches
output "ds_app_name" {
  value = data.yamlflattener_flatten.datasource_test.flattened["application.name"]
}

output "fn_app_name" {
  value = local.function_result["application.name"]
}

output "ds_first_tag" {
  value = data.yamlflattener_flatten.datasource_test.flattened["application.metadata.tags[0]"]
}

output "fn_first_tag" {
  value = local.function_result["application.metadata.tags[0]"]
}

output "ds_prod_debug" {
  value = data.yamlflattener_flatten.datasource_test.flattened["application.environments.production.debug"]
}

output "fn_prod_debug" {
  value = local.function_result["application.environments.production.debug"]
}
`, yamlContent)
}

// testFinalIntegration_CheckPathEquivalence checks that specific paths have the same values
// in both datasource and function outputs
func testFinalIntegration_CheckPathEquivalence(paths []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get the outputs from the state
		dsOutput, ok := s.RootModule().Outputs["datasource_output"]
		if !ok {
			return fmt.Errorf("datasource_output not found in state")
		}

		fnOutput, ok := s.RootModule().Outputs["function_output"]
		if !ok {
			return fmt.Errorf("function_output not found in state")
		}

		// Convert outputs to maps
		dsMap, ok := dsOutput.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("datasource_output is not a map")
		}

		fnMap, ok := fnOutput.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("function_output is not a map")
		}

		// Check each path
		for _, path := range paths {
			dsValue, dsOk := dsMap[path]
			fnValue, fnOk := fnMap[path]

			if !dsOk && !fnOk {
				// If both don't have the path, that's fine
				continue
			}

			if dsOk != fnOk {
				return fmt.Errorf("path %s exists in one output but not the other", path)
			}

			if fmt.Sprintf("%v", dsValue) != fmt.Sprintf("%v", fnValue) {
				return fmt.Errorf("values for path %s don't match: datasource=%v, function=%v",
					path, dsValue, fnValue)
			}
		}

		return nil
	}
}
