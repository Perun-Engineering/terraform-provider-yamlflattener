# Data Source: yamlflattener_flatten

The `yamlflattener_flatten` data source flattens nested YAML structures into a flat key-value map with dot notation for nested objects and bracket notation for arrays.

## Example Usage

```hcl
# Using yaml_content directly
data "yamlflattener_flatten" "example" {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
  EOT
}

# Using a YAML file
data "yamlflattener_flatten" "from_file" {
  yaml_file = "${path.module}/config.yaml"
}

# Accessing flattened values
output "nested_value" {
  value = data.yamlflattener_flatten.example.flattened["key2.nested"]
}

output "array_item" {
  value = data.yamlflattener_flatten.example.flattened["array[0].name"]
}
```

## Argument Reference

The following arguments are supported:

* `yaml_content` - (Optional) The YAML content to flatten as a string. Either `yaml_content` or `yaml_file` must be provided. Maximum size is 10MB.
* `yaml_file` - (Optional) The path to a YAML file to read and flatten. Either `yaml_content` or `yaml_file` must be provided. File paths with directory traversal patterns (`..`) are rejected for security reasons. Maximum file size is 10MB.

## Attribute Reference

The following attributes are exported:

* `flattened` - A map of string key-value pairs representing the flattened YAML structure.

## Performance and Security Considerations

The data source implements several performance and security measures:

* **Size Limits**: YAML content is limited to 10MB to prevent memory exhaustion
* **Nesting Depth**: Maximum nesting depth is limited to 100 levels to prevent stack overflow
* **Result Size**: Maximum number of flattened key-value pairs is limited to 100,000
* **Path Validation**: File paths are validated to prevent directory traversal attacks
* **Input Sanitization**: YAML content is sanitized to remove potentially dangerous characters

## Flattening Rules

The provider follows these rules when flattening YAML:

1. **Objects**: Nested objects are flattened using dot notation
   ```yaml
   parent:
     child:
       key: value
   ```
   Becomes: `parent.child.key = "value"`

2. **Arrays**: Arrays are flattened using bracket notation with indices
   ```yaml
   array:
     - name: item1
     - name: item2
   ```
   Becomes: `array[0].name = "item1"` and `array[1].name = "item2"`

3. **Mixed Structures**: Both notations are combined for complex structures
   ```yaml
   parent:
     children:
       - name: child1
         age: 10
       - name: child2
         age: 12
   ```
   Becomes:
   - `parent.children[0].name = "child1"`
   - `parent.children[0].age = "10"`
   - `parent.children[1].name = "child2"`
   - `parent.children[1].age = "12"`

4. **Data Types**: All values are converted to strings in the flattened output
   - Strings: Preserved as-is
   - Numbers: Converted to string representation
   - Booleans: Converted to "true" or "false"
   - Null: Represented as empty string

## Complete Example: Alertmanager Configuration

```hcl
data "yamlflattener_flatten" "alertmanager" {
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

# Access specific configuration values
output "slack_webhook" {
  value = data.yamlflattener_flatten.alertmanager.flattened["alertmanager.config.global.slack_api_url"]
}

output "channel" {
  value = data.yamlflattener_flatten.alertmanager.flattened["alertmanager.config.receivers[0].slack_configs[0].channel"]
}

output "send_resolved" {
  value = data.yamlflattener_flatten.alertmanager.flattened["alertmanager.config.receivers[0].slack_configs[0].send_resolved"]
}
```
