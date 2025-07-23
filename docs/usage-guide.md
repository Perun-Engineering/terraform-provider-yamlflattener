# Usage Guide

This comprehensive guide covers all aspects of using the YAML Flattener provider, from basic usage to advanced integration patterns.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Basic Usage](#basic-usage)
3. [Advanced Patterns](#advanced-patterns)
4. [Integration Examples](#integration-examples)
5. [Best Practices](#best-practices)
6. [Performance Optimization](#performance-optimization)
7. [Error Handling](#error-handling)

## Quick Start

### 1. Provider Configuration

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 1.0"
    }
  }
}

provider "yamlflattener" {}
```

### 2. Basic Flattening

```hcl
# Using data source with inline YAML
data "yamlflattener_flatten" "config" {
  yaml_content = <<-EOT
    app:
      name: "my-app"
      database:
        host: "localhost"
        port: 5432
  EOT
}

# Access flattened values
output "app_name" {
  value = data.yamlflattener_flatten.config.flattened["app.name"]
}

output "db_host" {
  value = data.yamlflattener_flatten.config.flattened["app.database.host"]
}
```

### 3. Using with Files

```hcl
# Using data source with external file
data "yamlflattener_flatten" "file_config" {
  yaml_file = "${path.module}/config.yaml"
}

# Using provider function
locals {
  yaml_content = file("${path.module}/config.yaml")
  flattened = provider::yamlflattener::flatten(local.yaml_content)
}
```

## Basic Usage

### Data Source Approach

The data source approach is ideal for declarative configurations:

```hcl
data "yamlflattener_flatten" "app_config" {
  yaml_content = <<-EOT
    application:
      name: "web-service"
      version: "1.0.0"
      environment: "production"

    server:
      host: "0.0.0.0"
      port: 8080
      ssl:
        enabled: true
        cert_path: "/etc/ssl/certs/app.crt"

    database:
      primary:
        host: "db-primary.example.com"
        port: 5432
      replicas:
        - host: "db-replica-1.example.com"
          port: 5432
        - host: "db-replica-2.example.com"
          port: 5432
  EOT
}

# Extract values
locals {
  app_name = data.yamlflattener_flatten.app_config.flattened["application.name"]
  db_host  = data.yamlflattener_flatten.app_config.flattened["database.primary.host"]
  replica1 = data.yamlflattener_flatten.app_config.flattened["database.replicas[0].host"]
}
```

### Provider Function Approach

The function approach offers more flexibility for inline processing:

```hcl
locals {
  config_template = <<-EOT
    api:
      version: "v1"
      endpoints:
        - path: "/users"
          methods: ["GET", "POST"]
        - path: "/orders"
          methods: ["GET", "POST", "PUT"]

    cache:
      redis:
        host: "${var.redis_host}"
        port: ${var.redis_port}
  EOT

  # Flatten the templated YAML
  flattened_config = provider::yamlflattener::flatten(local.config_template)

  # Use flattened values
  api_version = local.flattened_config["api.version"]
  first_endpoint = local.flattened_config["api.endpoints[0].path"]
  redis_host = local.flattened_config["cache.redis.host"]
}
```

## Advanced Patterns

### 1. Multi-Environment Configuration

Manage different environments using the same YAML structure:

```hcl
# environments/dev.yaml, environments/prod.yaml, etc.
data "yamlflattener_flatten" "env_config" {
  yaml_file = "${path.module}/environments/${var.environment}.yaml"
}

locals {
  config = data.yamlflattener_flatten.env_config.flattened

  # Environment-specific values
  instance_type = local.config["infrastructure.instance_type"]
  database_size = local.config["database.instance_class"]
  replica_count = tonumber(local.config["application.replicas"])
}

# Use in resources
resource "aws_instance" "app" {
  instance_type = local.instance_type
  count         = local.replica_count

  tags = {
    Environment = var.environment
    Application = local.config["application.name"]
  }
}
```

### 2. Dynamic Array Processing

Handle arrays of unknown length using Terraform's `for` expressions:

```hcl
data "yamlflattener_flatten" "services" {
  yaml_content = <<-EOT
    services:
      - name: "web"
        port: 80
        replicas: 3
      - name: "api"
        port: 8080
        replicas: 2
      - name: "worker"
        port: 9090
        replicas: 1
  EOT
}

locals {
  flattened = data.yamlflattener_flatten.services.flattened

  # Extract service configurations dynamically
  service_names = [
    for key in keys(local.flattened) :
    regex("services\\[(\\d+)\\]\\.name", key)[0]
    if can(regex("services\\[\\d+\\]\\.name$", key))
  ]

  services = {
    for idx in local.service_names :
    local.flattened["services[${idx}].name"] => {
      name     = local.flattened["services[${idx}].name"]
      port     = tonumber(local.flattened["services[${idx}].port"])
      replicas = tonumber(local.flattened["services[${idx}].replicas"])
    }
  }
}

# Create resources for each service
resource "kubernetes_deployment" "services" {
  for_each = local.services

  metadata {
    name = each.value.name
  }

  spec {
    replicas = each.value.replicas

    # ... rest of deployment configuration
  }
}
```

### 3. Configuration Validation

Validate configuration values using Terraform's validation features:

```hcl
data "yamlflattener_flatten" "config" {
  yaml_file = "${path.module}/config.yaml"
}

locals {
  config = data.yamlflattener_flatten.config.flattened

  # Validate required keys exist
  required_keys = [
    "application.name",
    "application.version",
    "database.host",
    "database.port"
  ]

  missing_keys = [
    for key in local.required_keys : key
    if !contains(keys(local.config), key)
  ]
}

# Validation check
resource "null_resource" "config_validation" {
  count = length(local.missing_keys) > 0 ? 1 : 0

  provisioner "local-exec" {
    command = "echo 'Missing required configuration keys: ${join(", ", local.missing_keys)}' && exit 1"
  }
}
```

### 4. Configuration Merging

Merge multiple YAML configurations:

```hcl
# Base configuration
data "yamlflattener_flatten" "base_config" {
  yaml_file = "${path.module}/base-config.yaml"
}

# Environment-specific overrides
data "yamlflattener_flatten" "env_config" {
  yaml_file = "${path.module}/env-${var.environment}.yaml"
}

locals {
  # Merge configurations (environment overrides base)
  merged_config = merge(
    data.yamlflattener_flatten.base_config.flattened,
    data.yamlflattener_flatten.env_config.flattened
  )

  # Use merged configuration
  app_name = local.merged_config["application.name"]
  db_host  = local.merged_config["database.host"]
}
```

## Integration Examples

### AWS Integration

```hcl
data "yamlflattener_flatten" "aws_config" {
  yaml_file = "${path.module}/aws-infrastructure.yaml"
}

locals {
  config = data.yamlflattener_flatten.aws_config.flattened
}

# VPC
resource "aws_vpc" "main" {
  cidr_block = local.config["vpc.cidr"]

  tags = {
    Name        = local.config["vpc.name"]
    Environment = local.config["environment"]
  }
}

# Security Groups
resource "aws_security_group" "web" {
  name_prefix = "${local.config["application.name"]}-web-"
  vpc_id      = aws_vpc.main.id

  dynamic "ingress" {
    for_each = {
      for key, value in local.config :
      key => value
      if can(regex("^security_groups\\.web\\.ingress\\[\\d+\\]", key))
    }

    content {
      from_port   = tonumber(local.config["${ingress.key}.from_port"])
      to_port     = tonumber(local.config["${ingress.key}.to_port"])
      protocol    = local.config["${ingress.key}.protocol"]
      cidr_blocks = [local.config["${ingress.key}.cidr_blocks[0]"]]
    }
  }
}
```

### Kubernetes Integration

```hcl
data "yamlflattener_flatten" "k8s_config" {
  yaml_file = "${path.module}/kubernetes-app.yaml"
}

locals {
  config = data.yamlflattener_flatten.k8s_config.flattened
}

# ConfigMap from flattened YAML
resource "kubernetes_config_map" "app_config" {
  metadata {
    name      = local.config["application.name"]
    namespace = local.config["application.namespace"]
  }

  data = {
    for key, value in local.config :
    replace(replace(key, "application.config.", ""), ".", "_") => value
    if startswith(key, "application.config.")
  }
}

# Deployment
resource "kubernetes_deployment" "app" {
  metadata {
    name      = local.config["application.name"]
    namespace = local.config["application.namespace"]
  }

  spec {
    replicas = tonumber(local.config["application.replicas"])

    template {
      spec {
        container {
          name  = local.config["application.name"]
          image = "${local.config["application.image.repository"]}:${local.config["application.image.tag"]}"

          env_from {
            config_map_ref {
              name = kubernetes_config_map.app_config.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = local.config["application.resources.requests.cpu"]
              memory = local.config["application.resources.requests.memory"]
            }
            limits = {
              cpu    = local.config["application.resources.limits.cpu"]
              memory = local.config["application.resources.limits.memory"]
            }
          }
        }
      }
    }
  }
}
```

## Best Practices

### 1. Configuration Organization

```hcl
# Organize configuration extraction in locals
locals {
  config = data.yamlflattener_flatten.app_config.flattened

  # Group related configurations
  app = {
    name        = local.config["application.name"]
    version     = local.config["application.version"]
    environment = local.config["application.environment"]
  }

  database = {
    host     = local.config["database.host"]
    port     = tonumber(local.config["database.port"])
    name     = local.config["database.name"]
    ssl_mode = local.config["database.ssl_mode"]
  }

  cache = {
    redis_host = local.config["cache.redis.host"]
    redis_port = tonumber(local.config["cache.redis.port"])
  }
}
```

### 2. Error Handling

```hcl
locals {
  config = data.yamlflattener_flatten.app_config.flattened

  # Safe access with defaults
  database_port = try(tonumber(local.config["database.port"]), 5432)
  ssl_enabled   = try(tobool(local.config["ssl.enabled"]), false)

  # Conditional configuration
  has_redis = contains(keys(local.config), "cache.redis.host")
  redis_config = local.has_redis ? {
    host = local.config["cache.redis.host"]
    port = tonumber(local.config["cache.redis.port"])
  } : null
}
```

### 3. Type Conversion

```hcl
locals {
  config = data.yamlflattener_flatten.app_config.flattened

  # Convert types as needed
  numeric_values = {
    port         = tonumber(local.config["server.port"])
    timeout      = tonumber(local.config["server.timeout"])
    max_connections = tonumber(local.config["server.max_connections"])
  }

  boolean_values = {
    ssl_enabled    = tobool(local.config["server.ssl.enabled"])
    debug_mode     = tobool(local.config["application.debug"])
    metrics_enabled = tobool(local.config["monitoring.enabled"])
  }

  list_values = {
    allowed_origins = split(",", local.config["cors.allowed_origins"])
    log_levels     = split(",", local.config["logging.levels"])
  }
}
```

## Performance Optimization

### 1. Caching Results

```hcl
# Cache flattened results in locals
locals {
  config = data.yamlflattener_flatten.large_config.flattened

  # Use cached results throughout
  app_name = local.config["application.name"]
  db_host  = local.config["database.host"]
}

# Don't repeat data source references
output "app_info" {
  value = {
    name = local.app_name  # ✅ Good
    host = local.db_host   # ✅ Good
    # name = data.yamlflattener_flatten.large_config.flattened["application.name"]  # ❌ Avoid
  }
}
```

### 2. Selective Processing

```hcl
# Process only what you need
data "yamlflattener_flatten" "database_config" {
  yaml_content = yamlencode({
    database = yamldecode(file("${path.module}/full-config.yaml")).database
  })
}

# Or split large configurations
data "yamlflattener_flatten" "app_config" {
  yaml_file = "${path.module}/app-config.yaml"
}

data "yamlflattener_flatten" "db_config" {
  yaml_file = "${path.module}/database-config.yaml"
}
```

### 3. Provider Configuration

```hcl
# Optimize provider settings for your use case
provider "yamlflattener" {
  max_depth     = 50     # Reduce if you don't need deep nesting
  max_keys      = 10000  # Adjust based on your configuration size
  max_file_size = 5242880  # 5MB - adjust based on your files
}
```

## Error Handling

### Common Error Patterns

```hcl
locals {
  config = data.yamlflattener_flatten.app_config.flattened

  # Handle missing keys
  database_config = {
    host = try(local.config["database.host"], "localhost")
    port = try(tonumber(local.config["database.port"]), 5432)
    ssl  = try(tobool(local.config["database.ssl"]), false)
  }

  # Validate required keys
  required_keys = ["application.name", "database.host"]
  missing_keys = [
    for key in local.required_keys : key
    if !contains(keys(local.config), key)
  ]
}

# Fail fast on missing required configuration
resource "null_resource" "validate_config" {
  count = length(local.missing_keys) > 0 ? 1 : 0

  provisioner "local-exec" {
    command = "echo 'ERROR: Missing required keys: ${join(", ", local.missing_keys)}' && exit 1"
  }
}
```

### Debugging Configuration

```hcl
# Debug outputs for troubleshooting
output "debug_all_keys" {
  value = keys(data.yamlflattener_flatten.config.flattened)
}

output "debug_database_keys" {
  value = [
    for key in keys(data.yamlflattener_flatten.config.flattened) : key
    if startswith(key, "database.")
  ]
}

output "debug_config_sample" {
  value = {
    for key, value in data.yamlflattener_flatten.config.flattened :
    key => value
    if length(regexall("^(application|database)\\.", key)) > 0
  }
}
```

This usage guide provides comprehensive coverage of the YAML Flattener provider's capabilities and should help users implement it effectively in their Terraform configurations.
