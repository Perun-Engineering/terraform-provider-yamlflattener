terraform {
  required_providers {
    yamlflattener = {
      source = "Perun-Engineering/yamlflattener"
    }
  }
}

# Create a YAML file to read from
resource "local_file" "config" {
  content = <<-EOT
    database:
      host: localhost
      port: 5432
      credentials:
        username: admin
        password: secret
    services:
      - name: web
        port: 8080
      - name: api
        port: 3000
  EOT
  filename = "${path.module}/config.yaml"
}

# Example using data source with file input
data "yamlflattener_flatten" "config" {
  yaml_file = local_file.config.filename
}

# Output the flattened result
output "flattened_config" {
  value = data.yamlflattener_flatten.config.flattened
}
