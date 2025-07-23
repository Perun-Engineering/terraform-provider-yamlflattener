# YAML Flattener Data Source Example

This example demonstrates how to use the YAML Flattener provider's data source to flatten nested YAML structures into a map with dot notation for nested objects and bracket notation for arrays.

## Usage

```hcl
# Using yaml_content directly
data "yamlflattener_flatten" "example" {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
  EOT
}

# Using a YAML file
data "yamlflattener_flatten" "from_file" {
  yaml_file = "/path/to/file.yaml"
}

# Accessing flattened values
output "example_value" {
  value = data.yamlflattener_flatten.example.flattened["key2.nested"]
}
```

## Examples in this Directory

This directory contains several examples demonstrating different aspects of the YAML Flattener data source:

1. **Simple Structure**: Basic nested objects and arrays
2. **File Input**: Reading from a YAML file
3. **Complex Structure**: The alertmanager example from the requirements
4. **Data Types**: Handling different data types (strings, numbers, booleans, nulls)
5. **Nested Arrays**: Handling multi-dimensional arrays

## Running the Example

To run this example, execute:

```bash
# Initialize Terraform
terraform init

# Apply the configuration
terraform apply
```

The outputs will show the flattened YAML structures and demonstrate how to access specific values from the flattened map.
