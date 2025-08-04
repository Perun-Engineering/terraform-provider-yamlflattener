---
page_title: "Provider: yamlflattener"
description: |-
  The yamlflattener provider flattens nested YAML structures into flat key-value maps with dot notation for objects and bracket notation for arrays.
---

# yamlflattener Provider

The yamlflattener provider flattens nested YAML structures into flat key-value maps with dot notation for objects and bracket notation for arrays.

## Why This Provider?

When using the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs), you can define sensitive variables using `set_sensitive` blocks. However, Terraform lacks a built-in function to flatten nested YAML structures, making it difficult to dynamically create `set_sensitive` blocks for complex configurations with nested objects and arrays.

This provider solves that problem by converting complex YAML structures into flat key-value pairs that can be easily used with Helm's `set_sensitive` blocks or any other Terraform resource that requires flattened configuration data.

## Example Usage

```terraform
terraform {
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 0.2"
    }
  }
}

# Flatten YAML content
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

# Use with Helm provider
resource "helm_release" "app" {
  name  = "my-app"
  chart = "my-chart"

  dynamic "set_sensitive" {
    for_each = data.yamlflattener_flatten.example.flattened
    content {
      name  = set_sensitive.key
      value = set_sensitive.value
    }
  }
}
```

## Schema

This provider does not require any configuration.
