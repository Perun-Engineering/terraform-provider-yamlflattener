# Troubleshooting Guide

This guide helps you resolve common issues when using the YAML Flattener provider.

## Common Issues

### 1. YAML Parsing Errors

#### Problem: "Invalid YAML syntax" error

**Symptoms:**
```
Error: failed to parse YAML content: yaml: line 5: found character that cannot start any token
```

**Solutions:**
- Check YAML indentation (use spaces, not tabs)
- Validate YAML syntax using an online validator
- Ensure proper quoting of special characters
- Check for trailing spaces or invisible characters

**Example Fix:**
```yaml
# ❌ Incorrect (mixed tabs and spaces)
key1: value1
	key2:
  nested: value2

# ✅ Correct (consistent spaces)
key1: value1
key2:
  nested: value2
```

#### Problem: "YAML contains unsupported features" error

**Symptoms:**
```
Error: YAML contains unsupported features: anchors and aliases are not supported
```

**Solutions:**
- Remove YAML anchors (`&`) and aliases (`*`)
- Expand repeated sections manually
- Use Terraform variables for repeated values instead

### 2. File Reading Errors

#### Problem: "File not found" error

**Symptoms:**
```
Error: failed to read YAML file: open /path/to/file.yaml: no such file or directory
```

**Solutions:**
- Verify the file path is correct
- Use absolute paths or `${path.module}` for relative paths
- Check file permissions
- Ensure the file exists in the expected location

**Example Fix:**
```hcl
# ❌ Incorrect (relative path without module reference)
data "yamlflattener_flatten" "config" {
  yaml_file = "config.yaml"
}

# ✅ Correct (using path.module)
data "yamlflattener_flatten" "config" {
  yaml_file = "${path.module}/config.yaml"
}
```

#### Problem: "Permission denied" error

**Symptoms:**
```
Error: failed to read YAML file: open /path/to/file.yaml: permission denied
```

**Solutions:**
- Check file permissions: `ls -la /path/to/file.yaml`
- Ensure Terraform has read access to the file
- On Unix systems, use `chmod 644 file.yaml`

### 3. Size and Performance Issues

#### Problem: "File too large" error

**Symptoms:**
```
Error: YAML file size exceeds maximum limit of 10MB
```

**Solutions:**
- Split large YAML files into smaller chunks
- Increase the provider's `max_file_size` setting
- Use external processing for very large files

**Example Fix:**
```hcl
provider "yamlflattener" {
  max_file_size = 20971520  # 20MB
}
```

#### Problem: "Too many keys" error

**Symptoms:**
```
Error: flattened result exceeds maximum key limit of 100,000
```

**Solutions:**
- Reduce YAML complexity
- Increase the provider's `max_keys` setting
- Process YAML in smaller sections

#### Problem: "Maximum nesting depth exceeded" error

**Symptoms:**
```
Error: YAML nesting depth exceeds maximum limit of 100 levels
```

**Solutions:**
- Reduce YAML nesting depth
- Increase the provider's `max_depth` setting
- Restructure YAML to be less deeply nested

### 4. Configuration Issues

#### Problem: "Both yaml_content and yaml_file provided" error

**Symptoms:**
```
Error: either yaml_content or yaml_file must be provided, not both
```

**Solutions:**
- Use only one input method per data source
- Create separate data sources for different inputs

**Example Fix:**
```hcl
# ❌ Incorrect (both inputs provided)
data "yamlflattener_flatten" "config" {
  yaml_content = "key: value"
  yaml_file    = "config.yaml"
}

# ✅ Correct (single input)
data "yamlflattener_flatten" "config" {
  yaml_content = "key: value"
}
```

#### Problem: "No input provided" error

**Symptoms:**
```
Error: either yaml_content or yaml_file must be provided
```

**Solutions:**
- Provide either `yaml_content` or `yaml_file`
- Check for empty variables or conditionals

### 5. Output and Access Issues

#### Problem: "Key not found in flattened map" error

**Symptoms:**
```
Error: This object does not have an attribute named "nonexistent.key"
```

**Solutions:**
- Check the exact key name in the flattened output
- Use `keys()` function to see all available keys
- Verify YAML structure matches expected format

**Example Debug:**
```hcl
# Debug: Show all available keys
output "debug_keys" {
  value = keys(data.yamlflattener_flatten.config.flattened)
}

# Then access the correct key
output "value" {
  value = data.yamlflattener_flatten.config.flattened["correct.key.name"]
}
```

#### Problem: "Unexpected data type" error

**Symptoms:**
Values are strings when numbers or booleans are expected.

**Solutions:**
- Remember all flattened values are strings
- Use Terraform type conversion functions
- Handle type conversion in your configuration

**Example Fix:**
```hcl
# Convert string to number
locals {
  port_number = tonumber(data.yamlflattener_flatten.config.flattened["server.port"])
  ssl_enabled = tobool(data.yamlflattener_flatten.config.flattened["server.ssl.enabled"])
}
```

### 6. Provider Function Issues

#### Problem: "Function not available" error

**Symptoms:**
```
Error: Call to unknown function "yamlflattener_flatten"
```

**Solutions:**
- Ensure provider is properly configured
- Use correct function syntax: `provider::yamlflattener::flatten()`
- Check Terraform version supports provider functions

**Example Fix:**
```hcl
# ❌ Incorrect function call
locals {
  result = yamlflattener_flatten(local.yaml_content)
}

# ✅ Correct function call
locals {
  result = provider::yamlflattener::flatten(local.yaml_content)
}
```

## Performance Optimization

### Large YAML Files

For better performance with large YAML files:

1. **Use appropriate limits:**
   ```hcl
   provider "yamlflattener" {
     max_depth     = 50    # Reduce if not needed
     max_keys      = 10000 # Adjust based on needs
     max_file_size = 5242880  # 5MB
   }
   ```

2. **Process in chunks:**
   ```hcl
   # Split large configurations
   data "yamlflattener_flatten" "app_config" {
     yaml_file = "${path.module}/app-config.yaml"
   }

   data "yamlflattener_flatten" "db_config" {
     yaml_file = "${path.module}/db-config.yaml"
   }
   ```

3. **Use caching:**
   ```hcl
   # Store results in locals for reuse
   locals {
     config = data.yamlflattener_flatten.app_config.flattened
     db_host = local.config["database.host"]
     db_port = local.config["database.port"]
   }
   ```

## Debugging Tips

### 1. Enable Debug Logging

Set environment variables for detailed logging:

```bash
export TF_LOG=DEBUG
export TF_LOG_PROVIDER=DEBUG
terraform plan
```

### 2. Validate YAML Separately

Test YAML syntax before using with Terraform:

```bash
# Using yq (if installed)
yq eval . config.yaml

# Using Python
python -c "import yaml; yaml.safe_load(open('config.yaml'))"

# Online validators
# https://yamlchecker.com/
```

### 3. Inspect Flattened Output

Use outputs to debug flattening results:

```hcl
# Show all flattened keys
output "all_keys" {
  value = keys(data.yamlflattener_flatten.config.flattened)
}

# Show specific key patterns
output "database_keys" {
  value = [for k in keys(data.yamlflattener_flatten.config.flattened) : k if can(regex("^database\\.", k))]
}

# Show full flattened map
output "full_config" {
  value = data.yamlflattener_flatten.config.flattened
}
```

## Getting Help

If you're still experiencing issues:

1. **Check the FAQ** in this documentation
2. **Search existing issues** on GitHub
3. **Create a new issue** with:
   - Terraform version
   - Provider version
   - Minimal reproduction case
   - Error messages
   - YAML content (sanitized)

## Reporting Bugs

When reporting bugs, please include:

- **Environment details:**
  - Operating system
  - Terraform version
  - Provider version

- **Reproduction steps:**
  - Minimal Terraform configuration
  - YAML content that causes the issue
  - Commands run

- **Expected vs actual behavior**
- **Error messages and logs**

Submit issues at: https://github.com/Perun-Engineering/terraform-provider-yamlflattener/issues
