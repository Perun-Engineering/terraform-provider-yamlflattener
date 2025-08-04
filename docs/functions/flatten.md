---
page_title: "flatten Function - yamlflattener"
subcategory: ""
description: |-
  Flattens a nested YAML structure into an ordered list of key-value pairs using dot notation for objects and bracket notation for arrays.
---

# flatten Function

Flattens a nested YAML structure into an ordered list of key-value pairs using dot notation for objects and bracket notation for arrays.

This function provides the same functionality as the data source but returns an ordered list that preserves the exact key order from the original YAML structure.

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

  # Get ordered list of key-value pairs
  flattened_pairs = provider::yamlflattener::flatten(local.yaml_config)

  # Convert to map when needed (loses order)
  flattened_map = { for pair in local.flattened_pairs : pair[0] => pair[1] }
}

# Use with Helm provider - preserves order!
resource "helm_release" "app" {
  name  = "my-app"
  chart = "my-chart"

  dynamic "set_sensitive" {
    for_each = local.flattened_pairs
    content {
      name  = set_sensitive.value[0]  # key
      value = set_sensitive.value[1]  # value
    }
  }
}

# Access specific values by position (order preserved)
output "database_host" {
  value = local.flattened_map["database.host"]
}

# Show all keys in original YAML order
output "ordered_keys" {
  value = [for pair in local.flattened_pairs : pair[0]]
}
```

## Signature

```
flatten(yaml_content string) list(tuple([string, string]))
```

## Arguments

1. `yaml_content` (String) - The YAML content to flatten as a string

## Return Type

The function returns a list of tuples where:
- Each tuple contains `[key, value]`
- Keys use dot and bracket notation for nested structures
- Values are string representations of the original YAML values
- **Order is preserved** exactly as it appears in the YAML

## Order Preservation

Unlike maps, this list format preserves the exact order from your YAML:

```terraform
# Input YAML:
# receivers:
#   - name: blackhole
#   - name: discord_prometheus
#     discord_config:
#       - webhook_url: "https://example.com/webhook"

# Output (order preserved):
# [
#   ["receivers[0].name", "blackhole"],
#   ["receivers[1].name", "discord_prometheus"],
#   ["receivers[1].discord_config[0].webhook_url", "https://example.com/webhook"]
# ]
```
