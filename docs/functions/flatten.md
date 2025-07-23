# flatten Function

Flattens YAML content into a key-value map with dot notation for objects and bracket notation for arrays.

## Example Usage

```hcl
locals {
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

  flattened = provider::yamlflattener::flatten(local.yaml_content)
}

output "database_host" {
  value = local.flattened["database.host"]
}

output "first_service_port" {
  value = local.flattened["services[0].port"]
}
```

## Signature

```
flatten(yaml_content string) map(string)
```

## Arguments

- `yaml_content` - (Required) The YAML content to flatten as a string.

## Return Type

Returns a map where keys are flattened paths and values are strings.

## Flattening Rules

- **Objects**: Nested objects use dot notation (e.g., `parent.child.key`)
- **Arrays**: Array elements use bracket notation (e.g., `parent.array[0].key`)
- **Primitive values**: Converted to strings
- **Null values**: Represented as empty strings
- **Boolean values**: Converted to "true"/"false"

## Error Handling

The function will fail if:
- YAML content is empty or invalid
- YAML syntax is malformed
- Maximum nesting depth is exceeded (protection against stack overflow)
