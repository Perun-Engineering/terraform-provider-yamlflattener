# Function: yamlflattener_flatten

The `yamlflattener_flatten` function flattens nested YAML structures into a flat key-value map with dot notation for nested objects and bracket notation for arrays.

## Example Usage

```hcl
locals {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
array:
  - name: item1
  - name: item2
  EOT

  # Use the function to flatten the YAML
  flattened = yamlflattener_flatten(local.yaml_content)
}

# Accessing flattened values
output "nested_value" {
  value = local.flattened["key2.nested"]
}

output "array_item" {
  value = local.flattened["array[0].name"]
}
```

## Argument Reference

The function accepts a single argument:

* `yaml_content` - (Required) The YAML content to flatten as a string. Maximum size is 10MB.

## Return Value

The function returns a map of string key-value pairs representing the flattened YAML structure.

## Performance and Security Considerations

The function implements several performance and security measures:

* **Size Limits**: YAML content is limited to 10MB to prevent memory exhaustion
* **Nesting Depth**: Maximum nesting depth is limited to 100 levels to prevent stack overflow
* **Result Size**: Maximum number of flattened key-value pairs is limited to 100,000
* **Input Sanitization**: YAML content is sanitized to remove potentially dangerous characters
* **Parsing Timeouts**: YAML parsing operations have timeouts to prevent hanging on malicious input

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
locals {
  alertmanager_yaml = <<-EOT
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

  flattened = yamlflattener_flatten(local.alertmanager_yaml)
}

# Access specific configuration values
output "slack_webhook" {
  value = local.flattened["alertmanager.config.global.slack_api_url"]
}

output "channel" {
  value = local.flattened["alertmanager.config.receivers[0].slack_configs[0].channel"]
}

output "send_resolved" {
  value = local.flattened["alertmanager.config.receivers[0].slack_configs[0].send_resolved"]
}
```

## Function vs Data Source

The function approach offers more flexibility for inline usage within expressions and resource configurations, while the data source approach is more declarative and allows the flattened data to be referenced as data source attributes. Both implementations produce identical flattened output for the same input.

Key differences:

1. **Inline Usage**: The function can be used directly in expressions and resource attributes
2. **Input Methods**: The function only accepts YAML content as a string, while the data source can read from a file
3. **State Management**: The data source is tracked in Terraform state, while function results are computed during each plan/apply
