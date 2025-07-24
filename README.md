# Terraform YAML Flattener Provider

A Terraform provider that flattens nested YAML structures into flat key-value maps with dot notation for objects and bracket notation for arrays.

## Why This Provider?

When using the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs), you can define sensitive variables using `set_sensitive` blocks. However, Terraform lacks a built-in function to flatten nested YAML structures, making it difficult to dynamically create `set_sensitive` blocks for complex configurations with nested objects and arrays.

This provider solves that problem by converting complex YAML structures into flat key-value pairs that can be easily used with Helm's `set_sensitive` blocks or any other Terraform resource that requires flattened configuration data.

## Installation

```hcl
terraform {
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 0.1"
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
  flattened = provider::yamlflattener::flatten(file("config.yaml"))
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

## License

MIT License - see [LICENSE](LICENSE) file.
