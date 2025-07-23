terraform {
  required_providers {
    yamlflattener = {
      source = "local/yamlflattener"
    }
  }
}

provider "yamlflattener" {}

# Example 1: Using the flatten function with a simple structure
locals {
  simple_yaml = <<-EOT
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
EOT

  # Use the function to flatten the YAML
  simple_flattened = yamlflattener_flatten(local.simple_yaml)
}

# Example 2: Using the function with a complex structure (alertmanager example)
locals {
  complex_yaml = <<-EOT
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
EOT

  # Use the function to flatten the complex YAML
  complex_flattened = yamlflattener_flatten(local.complex_yaml)
}

# Example 3: Using the function with different data types
locals {
  data_types_yaml = <<-EOT
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
EOT

  # Use the function to flatten the data types YAML
  data_types_flattened = yamlflattener_flatten(local.data_types_yaml)
}

# Example 4: Using the function with nested arrays
locals {
  nested_arrays_yaml = <<-EOT
matrix:
  - - 1
    - 2
  - - 3
    - 4
EOT

  # Use the function to flatten the nested arrays YAML
  nested_arrays_flattened = yamlflattener_flatten(local.nested_arrays_yaml)
}

# Example 5: Using the function result in another resource
resource "local_file" "flattened_output" {
  content  = jsonencode({
    simple = local.simple_flattened
    complex = local.complex_flattened
    data_types = local.data_types_flattened
    nested_arrays = local.nested_arrays_flattened
  })
  filename = "${path.module}/flattened_output.json"
}

# Output all flattened values for verification
output "simple_flattened" {
  value = local.simple_flattened
}

output "complex_flattened" {
  value = local.complex_flattened
}

output "data_types_flattened" {
  value = local.data_types_flattened
}

output "nested_arrays_flattened" {
  value = local.nested_arrays_flattened
}

# Output specific values to demonstrate accessing individual flattened keys
output "nested_value" {
  value = local.simple_flattened["key2.nested"]
}

output "array_item" {
  value = local.simple_flattened["array[0].name"]
}

output "slack_webhook" {
  value = local.complex_flattened["alertmanager.config.global.slack_api_url"]
}

output "bool_as_string" {
  value = local.data_types_flattened["bool_val"]
}
