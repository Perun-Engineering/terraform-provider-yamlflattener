terraform {
  required_providers {
    yamlflattener = {
      source = "local/yamlflattener"
    }
  }
}

provider "yamlflattener" {}

# Example 1: Using yaml_content with a simple structure
data "yamlflattener_flatten" "simple" {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
EOT
}

# Example 2: Using yaml_file with a simple structure
data "yamlflattener_flatten" "from_file" {
  yaml_file = "${path.module}/simple.yaml"
}

# Example 3: Using yaml_content with a complex structure (alertmanager example)
data "yamlflattener_flatten" "complex" {
  yaml_content = <<-EOT
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
}

# Example 4: Using yaml_content with different data types
data "yamlflattener_flatten" "data_types" {
  yaml_content = <<-EOT
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
EOT
}

# Example 5: Using yaml_content with nested arrays
data "yamlflattener_flatten" "nested_arrays" {
  yaml_content = <<-EOT
matrix:
  - - 1
    - 2
  - - 3
    - 4
EOT
}

# Output all flattened values for verification
output "simple_flattened" {
  value = data.yamlflattener_flatten.simple.flattened
}

output "from_file_flattened" {
  value = data.yamlflattener_flatten.from_file.flattened
}

output "complex_flattened" {
  value = data.yamlflattener_flatten.complex.flattened
}

output "data_types_flattened" {
  value = data.yamlflattener_flatten.data_types.flattened
}

output "nested_arrays_flattened" {
  value = data.yamlflattener_flatten.nested_arrays.flattened
}

# Output specific values to demonstrate accessing individual flattened keys
output "nested_value" {
  value = data.yamlflattener_flatten.simple.flattened["key2.nested"]
}

output "array_item" {
  value = data.yamlflattener_flatten.simple.flattened["array[0].name"]
}

output "slack_webhook" {
  value = data.yamlflattener_flatten.complex.flattened["alertmanager.config.global.slack_api_url"]
}

output "bool_as_string" {
  value = data.yamlflattener_flatten.data_types.flattened["bool_val"]
}
