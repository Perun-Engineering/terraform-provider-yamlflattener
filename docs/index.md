# YAML Flattener Provider

The YAML Flattener provider is a Terraform provider that allows you to flatten nested YAML structures into a flat key-value map with dot notation for nested objects and bracket notation for arrays. This is particularly useful when you need to extract values from complex YAML configurations and use them in your Terraform resources.

## Example Usage

```hcl
terraform {
  required_providers {
    yamlflattener = {
      source = "local/yamlflattener"
    }
  }
}

provider "yamlflattener" {}

# Using the data source
data "yamlflattener_flatten" "example" {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
  EOT
}

# Using the provider function
locals {
  yaml_content = <<-EOT
key1: value1
key2:
  nested: value2
  EOT

  flattened = yamlflattener_flatten(local.yaml_content)
}

# Accessing values
output "from_data_source" {
  value = data.yamlflattener_flatten.example.flattened["key2.nested"]
}

output "from_function" {
  value = local.flattened["key2.nested"]
}
```

## Provider Configuration

The provider doesn't require any configuration parameters.

```hcl
provider "yamlflattener" {}
```

## Available Resources

The provider offers:

- A data source: `yamlflattener_flatten`
- A provider function: `yamlflattener_flatten()`

## Documentation

### Getting Started
- [Installation Guide](installation.md) - How to install and configure the provider
- [Usage Guide](usage-guide.md) - Comprehensive usage examples and patterns

### Reference Documentation
- [Data Source: yamlflattener_flatten](data-source.md) - Data source reference
- [Function: yamlflattener_flatten](function.md) - Provider function reference

### Help and Support
- [Troubleshooting Guide](troubleshooting.md) - Common issues and solutions
- [FAQ](faq.md) - Frequently asked questions

### Examples
See the `examples/` directory for complete working examples:
- [Basic Examples](../examples/) - Simple usage examples
- [Integration Examples](../examples/integration-examples/) - Real-world integrations with AWS, Kubernetes, and more
