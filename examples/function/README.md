# YAML Flattener Function Example

This example demonstrates how to use the YAML Flattener provider's function to flatten nested YAML structures into a map with dot notation for nested objects and bracket notation for arrays.

## Usage

```hcl
locals {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
  EOT

  # Use the function to flatten the YAML
  flattened = yamlflattener_flatten(local.yaml_content)
}

# Accessing flattened values
output "example_value" {
  value = local.flattened["key2.nested"]
}
```

## Examples in this Directory

This directory contains several examples demonstrating different aspects of the YAML Flattener function:

1. **Simple Structure**: Basic nested objects and arrays
2. **Complex Structure**: The alertmanager example from the requirements
3. **Data Types**: Handling different data types (strings, numbers, booleans, nulls)
4. **Nested Arrays**: Handling multi-dimensional arrays
5. **Resource Integration**: Using the function result in another resource

## Running the Example

To run this example, execute:

```bash
# Initialize Terraform
terraform init

# Apply the configuration
terraform apply
```

The outputs will show the flattened YAML structures and demonstrate how to access specific values from the flattened map.

## Function vs Data Source

The function approach offers more flexibility for inline usage within expressions and resource configurations, while the data source approach is more declarative and allows the flattened data to be referenced as data source attributes. Both implementations produce identical flattened output for the same input.
