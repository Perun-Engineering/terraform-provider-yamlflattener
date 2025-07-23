# yamlflattener_flatten Data Source

Flattens YAML content into a key-value map with dot notation for objects and bracket notation for arrays.

## Example Usage

### Inline YAML Content

```hcl
data "yamlflattener_flatten" "example" {
  yaml_content = <<-EOT
    database:
      host: localhost
      port: 5432
    services:
      - name: web
        port: 8080
      - name: api
        port: 3000
  EOT
}

output "database_host" {
  value = data.yamlflattener_flatten.example.flattened["database.host"]
}

output "first_service" {
  value = data.yamlflattener_flatten.example.flattened["services[0].name"]
}
```

### File Input

```hcl
data "yamlflattener_flatten" "config" {
  yaml_file = "config.yaml"
}
```

## Argument Reference

- `yaml_content` - (Optional) The YAML content to flatten. Mutually exclusive with `yaml_file`.
- `yaml_file` - (Optional) Path to YAML file to flatten. Mutually exclusive with `yaml_content`.

## Attribute Reference

- `flattened` - A map containing the flattened YAML structure.

## Flattening Rules

- **Objects**: Nested objects use dot notation (e.g., `parent.child.key`)
- **Arrays**: Array elements use bracket notation (e.g., `parent.array[0].key`)
- **Primitive values**: Converted to strings
- **Null values**: Represented as empty strings
- **Boolean values**: Converted to "true"/"false"

## Security

- File paths are validated to prevent directory traversal attacks
- Absolute paths are not allowed
- Input sanitization is performed on YAML content
