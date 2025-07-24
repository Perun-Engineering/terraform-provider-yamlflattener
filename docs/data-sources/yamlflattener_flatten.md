---
page_title: "yamlflattener_flatten Data Source - yamlflattener"
subcategory: ""
description: |-
  Flattens a nested YAML structure into a flat key-value map using dot notation for objects and bracket notation for arrays.
---

# yamlflattener_flatten (Data Source)

Flattens a nested YAML structure into a flat key-value map using dot notation for objects and bracket notation for arrays.

This data source is particularly useful when working with the Helm provider's `set_sensitive` blocks, allowing you to dynamically flatten complex YAML configurations.

## Example Usage

```terraform
# Flatten YAML content directly
data "yamlflattener_flatten" "config" {
  yaml_content = <<EOF
database:
  host: "localhost"
  port: 5432
  replicas:
    - host: "replica1"
      port: 5433
    - host: "replica2"
      port: 5434
EOF
}

# Flatten YAML from file
data "yamlflattener_flatten" "from_file" {
  yaml_file = "${path.module}/config.yaml"
}

# Access flattened values
output "database_host" {
  value = data.yamlflattener_flatten.config.flattened["database.host"]
}

output "first_replica" {
  value = data.yamlflattener_flatten.config.flattened["database.replicas[0].host"]
}
```

## Schema

### Required

One of the following must be specified:

- `yaml_content` (String) - The YAML content to flatten as a string
- `yaml_file` (String) - Path to a YAML file to read and flatten

### Read-Only

- `flattened` (Map of String) - The flattened key-value map
- `id` (String) - The ID of this resource

## Flattening Rules

- **Objects**: Flattened using dot notation (e.g., `key.subkey`)
- **Arrays**: Flattened using bracket notation (e.g., `key[0]`, `key[1]`)
- **Mixed**: Combinations use both notations (e.g., `key.array[0].subkey`)
- **Values**: All values are converted to strings
- **Null values**: Represented as empty strings
