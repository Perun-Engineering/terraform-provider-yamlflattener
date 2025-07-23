# Frequently Asked Questions (FAQ)

## General Questions

### Q: What is the YAML Flattener provider?

**A:** The YAML Flattener provider is a Terraform provider that converts nested YAML structures into flat key-value maps using dot notation for objects and bracket notation for arrays. This makes it easier to extract specific values from complex YAML configurations for use in Terraform resources.

### Q: When should I use this provider?

**A:** Use this provider when you need to:
- Extract specific values from complex YAML configuration files
- Convert nested YAML structures for use with other Terraform resources
- Process configuration files like Kubernetes manifests, application configs, or CI/CD configurations
- Transform YAML data for integration with cloud resources that expect flat key-value pairs

### Q: What's the difference between the data source and function approaches?

**A:** Both produce identical results but have different use cases:

- **Data Source (`yamlflattener_flatten`):**
  - More declarative approach
  - Can read from files or inline content
  - Tracked in Terraform state
  - Better for static configurations

- **Function (`provider::yamlflattener::flatten()`):**
  - More flexible for inline usage
  - Only accepts string content (no file reading)
  - Computed during each plan/apply
  - Better for dynamic processing

## Technical Questions

### Q: What YAML features are supported?

**A:** The provider supports:
- ✅ Nested objects and arrays
- ✅ All basic data types (strings, numbers, booleans, null)
- ✅ Multi-line strings
- ✅ Unicode characters
- ✅ Complex nested structures

**Not supported:**
- ❌ YAML anchors (`&`) and aliases (`*`)
- ❌ YAML tags
- ❌ Multi-document YAML files
- ❌ Custom data types

### Q: How are different data types handled in the flattened output?

**A:** All values are converted to strings:
- **Strings:** Preserved as-is
- **Numbers:** `123` → `"123"`, `3.14` → `"3.14"`
- **Booleans:** `true` → `"true"`, `false` → `"false"`
- **Null:** `null` → `""` (empty string)

Use Terraform's type conversion functions if you need other types:
```hcl
locals {
  port = tonumber(local.flattened["server.port"])
  enabled = tobool(local.flattened["feature.enabled"])
}
```

### Q: What are the size and performance limits?

**A:** Default limits (configurable):
- **File size:** 10MB maximum
- **Nesting depth:** 100 levels maximum
- **Number of keys:** 100,000 maximum
- **Memory usage:** Optimized for typical configurations

Configure limits in the provider:
```hcl
provider "yamlflattener" {
  max_file_size = 20971520  # 20MB
  max_depth     = 50
  max_keys      = 50000
}
```

### Q: How does array flattening work?

**A:** Arrays use bracket notation with zero-based indices:

```yaml
# Input YAML
items:
  - name: "first"
    value: 1
  - name: "second"
    value: 2
```

```hcl
# Flattened output
{
  "items[0].name"  = "first"
  "items[0].value" = "1"
  "items[1].name"  = "second"
  "items[1].value" = "2"
}
```

### Q: Can I use the provider with dynamic YAML content?

**A:** Yes, both approaches support dynamic content:

```hcl
# Dynamic content with data source
data "yamlflattener_flatten" "dynamic" {
  yaml_content = templatefile("${path.module}/config.yaml.tpl", {
    environment = var.environment
    region      = var.region
  })
}

# Dynamic content with function
locals {
  dynamic_yaml = templatefile("${path.module}/config.yaml.tpl", var.config)
  flattened = provider::yamlflattener::flatten(local.dynamic_yaml)
}
```

## Usage Questions

### Q: How do I handle missing keys gracefully?

**A:** Use Terraform's `try()` function or conditional expressions:

```hcl
# Using try() function
locals {
  optional_value = try(data.yamlflattener_flatten.config.flattened["optional.key"], "default")
}

# Using conditional expression
locals {
  has_ssl = contains(keys(data.yamlflattener_flatten.config.flattened), "server.ssl.enabled")
  ssl_enabled = local.has_ssl ? data.yamlflattener_flatten.config.flattened["server.ssl.enabled"] : "false"
}
```

### Q: How do I work with arrays of unknown length?

**A:** Use Terraform's `for` expressions to process arrays dynamically:

```hcl
# Get all database replica configurations
locals {
  flattened = data.yamlflattener_flatten.config.flattened

  # Find all replica keys
  replica_keys = [for k in keys(local.flattened) : k if can(regex("^database\\.replicas\\[\\d+\\]\\.host$", k))]

  # Extract replica configurations
  replicas = [for key in local.replica_keys : {
    host = local.flattened[key]
    port = local.flattened[replace(key, ".host", ".port")]
    name = local.flattened[replace(key, ".host", ".name")]
  }]
}
```

### Q: Can I use this provider with Kubernetes YAML manifests?

**A:** Yes, but with considerations:

```hcl
# Flatten a Kubernetes deployment
data "yamlflattener_flatten" "k8s_deployment" {
  yaml_file = "${path.module}/deployment.yaml"
}

# Extract specific values
locals {
  image = data.yamlflattener_flatten.k8s_deployment.flattened["spec.template.spec.containers[0].image"]
  replicas = data.yamlflattener_flatten.k8s_deployment.flattened["spec.replicas"]
}
```

**Note:** Multi-document YAML files (separated by `---`) are not supported. Split them into separate files first.

### Q: How do I integrate with other Terraform resources?

**A:** Use flattened values directly in resource configurations:

```hcl
# Flatten application configuration
data "yamlflattener_flatten" "app_config" {
  yaml_file = "${path.module}/app-config.yaml"
}

# Use in AWS resources
resource "aws_instance" "app_server" {
  ami           = var.ami_id
  instance_type = data.yamlflattener_flatten.app_config.flattened["server.instance_type"]

  tags = {
    Name        = data.yamlflattener_flatten.app_config.flattened["application.name"]
    Environment = data.yamlflattener_flatten.app_config.flattened["application.environment"]
  }
}

# Use in Kubernetes resources
resource "kubernetes_config_map" "app_config" {
  metadata {
    name = data.yamlflattener_flatten.app_config.flattened["application.name"]
  }

  data = {
    database_host = data.yamlflattener_flatten.app_config.flattened["database.host"]
    database_port = data.yamlflattener_flatten.app_config.flattened["database.port"]
    redis_url     = "${data.yamlflattener_flatten.app_config.flattened["cache.redis.host"]}:${data.yamlflattener_flatten.app_config.flattened["cache.redis.port"]}"
  }
}
```

## Troubleshooting Questions

### Q: Why am I getting "key not found" errors?

**A:** Common causes:
1. **Typo in key name** - Check exact spelling and case
2. **Wrong array index** - Arrays are zero-based
3. **Missing nested structure** - Verify YAML structure

Debug by showing all keys:
```hcl
output "debug_keys" {
  value = keys(data.yamlflattener_flatten.config.flattened)
}
```

### Q: Why is my YAML not parsing correctly?

**A:** Common issues:
1. **Indentation problems** - Use spaces, not tabs
2. **Special characters** - Quote strings with special characters
3. **YAML anchors/aliases** - Not supported, expand manually

Validate YAML syntax separately before using with Terraform.

### Q: How do I handle large YAML files efficiently?

**A:** Optimization strategies:
1. **Split large files** into smaller, focused configurations
2. **Adjust provider limits** based on your needs
3. **Use caching** by storing results in locals
4. **Process selectively** - only flatten what you need

```hcl
# Split approach
data "yamlflattener_flatten" "database_config" {
  yaml_file = "${path.module}/database.yaml"
}

data "yamlflattener_flatten" "application_config" {
  yaml_file = "${path.module}/application.yaml"
}

# Caching approach
locals {
  config = data.yamlflattener_flatten.large_config.flattened
  # Use local.config throughout instead of repeating the data source reference
}
```

## Best Practices

### Q: What are the recommended best practices?

**A:**
1. **Validate YAML** before using with Terraform
2. **Use meaningful variable names** for flattened results
3. **Handle missing keys gracefully** with `try()` or conditionals
4. **Cache results** in locals for repeated access
5. **Split large configurations** into focused files
6. **Document key patterns** in your code comments
7. **Use type conversion** functions when needed

### Q: Should I use the data source or function approach?

**A:** Choose based on your use case:

**Use Data Source when:**
- Reading from files
- Configuration is relatively static
- You want declarative approach
- You need state tracking

**Use Function when:**
- Processing dynamic content
- Using in expressions or calculations
- You need inline processing
- Working with templated YAML

### Q: How do I version control YAML configurations?

**A:** Best practices:
1. **Store YAML files** in version control alongside Terraform
2. **Use relative paths** with `${path.module}`
3. **Document key structures** in README files
4. **Validate YAML** in CI/CD pipelines
5. **Use semantic versioning** for configuration changes

## Migration and Compatibility

### Q: How do I migrate from other YAML processing methods?

**A:** Common migration patterns:

```hcl
# From yamldecode() function
# Old approach
locals {
  config = yamldecode(file("${path.module}/config.yaml"))
  database_host = local.config.database.host
}

# New approach with flattener
data "yamlflattener_flatten" "config" {
  yaml_file = "${path.module}/config.yaml"
}

locals {
  database_host = data.yamlflattener_flatten.config.flattened["database.host"]
}
```

### Q: Is the provider compatible with all Terraform versions?

**A:** Requirements:
- **Terraform >= 1.0** (required for provider functions)
- **Provider Framework v1.x** compatibility
- **Cross-platform support** (Linux, macOS, Windows)

Check the [installation guide](installation.md) for specific version requirements.

## Support and Community

### Q: Where can I get help?

**A:** Resources:
1. **Documentation** - This docs folder
2. **GitHub Issues** - Bug reports and feature requests
3. **Examples** - See the `examples/` directory
4. **Community** - Terraform community forums

### Q: How do I report bugs or request features?

**A:**
1. **Search existing issues** first
2. **Create detailed bug reports** with reproduction steps
3. **Include environment details** (OS, Terraform version, provider version)
4. **Provide minimal examples** that demonstrate the issue

Submit at: https://github.com/Perun-Engineering/terraform-provider-yamlflattener/issues

### Q: Can I contribute to the provider?

**A:** Yes! See the contribution guidelines in the repository. Common contributions:
- Bug fixes
- Documentation improvements
- Example configurations
- Feature enhancements
- Test cases
