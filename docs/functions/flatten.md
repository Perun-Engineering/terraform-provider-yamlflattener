---
page_title: "flatten Function - yamlflattener"
subcategory: ""
description: |-
  Flattens a nested YAML structure into a flat key-value map using dot notation for objects and bracket notation for arrays.
---

# flatten Function

Flattens a nested YAML structure into a flat key-value map using dot notation for objects and bracket notation for arrays.

This function provides the same functionality as the data source but can be used inline within expressions.

## Example Usage

```terraform
locals {
  yaml_config = <<EOF
database:
  host: "localhost"
  port: 5432
  replicas:
    - host: "replica1"
    - host: "replica2"
EOF

  flattened = provider::yamlflattener::flatten(local.yaml_config)
}

# Use with Helm provider
resource "helm_release" "app" {
  name  = "my-app"
  chart = "my-chart"

  dynamic "set_sensitive" {
    for_each = local.flattened
    content {
      name  = set_sensitive.key
      value = set_sensitive.value
    }
  }
}

# Access specific values
output "database_host" {
  value = local.flattened["database.host"]
}
```

## Signature

```
flatten(yaml_content string) map(string)
```

## Arguments

1. `yaml_content` (String) - The YAML content to flatten as a string

## Return Type

The function returns a map of strings where:
- Keys are the flattened paths using dot and bracket notation
- Values are string representations of the original YAML values
