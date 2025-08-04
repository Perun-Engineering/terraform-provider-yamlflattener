# Terraform YAML Flattener Provider

A Terraform provider that flattens nested YAML structures into flat key-value maps with dot notation for objects and bracket notation for arrays. **Preserves the original YAML key order** for consistent and predictable results.

## Why This Provider?

When using the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs), you can define sensitive variables using `set_sensitive` blocks. However, Terraform lacks a built-in function to flatten nested YAML structures, making it difficult to dynamically create `set_sensitive` blocks for complex configurations with nested objects and arrays.

This provider solves that problem by converting complex YAML structures into flat key-value pairs that can be easily used with Helm's `set_sensitive` blocks or any other Terraform resource that requires flattened configuration data.

## Key Features

- **Order Preservation**: Maintains the exact key order from your original YAML files
- **Dot Notation**: Nested objects become `parent.child` keys
- **Bracket Notation**: Arrays become `parent[0]`, `parent[1]` keys
- **Mixed Structures**: Handles complex combinations like `parent.array[0].child`
- **Type Conversion**: All values converted to strings for Terraform compatibility
- **Null Handling**: Null values become empty strings
- **Performance**: Single-pass parsing with order preservation

## Installation

```hcl
terraform {
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 0.2"
    }
  }
}
```

## Usage

### Data Source

```hcl
data "yamlflattener_flatten" "example" {
  yaml_content = <<EOF
database:
  host: "localhost"
  port: 5432
  credentials:
    username: "admin"
    password: "secret"
EOF
}

output "flattened_config" {
  value = data.yamlflattener_flatten.example.flattened
}
```

### Provider Function

```hcl
locals {
  # Get ordered list of key-value pairs
  flattened_pairs = provider::yamlflattener::flatten(file("config.yaml"))

  # Convert to map when needed (loses order)
  flattened_map = { for pair in local.flattened_pairs : pair[0] => pair[1] }
}

# Use with Helm - preserves order!
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
```

## Example Output

**Input:**
```yaml
database:
  host: "localhost"
  replicas:
    - host: "replica1"
    - host: "replica2"
```

**Output:**
```json
{
  "database.host": "localhost",
  "database.replicas[0].host": "replica1",
  "database.replicas[1].host": "replica2"
}
```

## Order Preservation Example

This provider maintains the exact order from your YAML files, which is crucial for consistent reconstruction:

**Input:**
```yaml
receivers:
  - name: blackhole
  - name: discord_prometheus
    discord_config:
      - webhook_url: "https://example.com/webhook"
```

**Output (order preserved):**
```json
{
  "receivers[0].name": "blackhole",
  "receivers[1].name": "discord_prometheus",
  "receivers[1].discord_config[0].webhook_url": "https://example.com/webhook"
}
```

When reconstructed, this maintains the original structure:
```yaml
receivers:
- name: blackhole
- name: discord_prometheus          # ✅ name comes first (as in original)
  discord_config:                   # ✅ discord_config comes second (as in original)
  - webhook_url: https://example.com/webhook
```

## License

MIT License - see [LICENSE](LICENSE) file.
