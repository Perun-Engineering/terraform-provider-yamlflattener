terraform {
  required_version = ">= 1.0"
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = ">= 0.2.0"
    }
  }
}

# Example 1: Using data source with inline YAML content
locals {
  app_config_yaml = <<EOT
application:
  name: "example-app"
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
    name: "app_db"
  replicas:
    - host: "db-replica-1.example.com"
      port: 5432
      name: "app_db"
    - host: "db-replica-2.example.com"
      port: 5432
      name: "app_db"

cache:
  redis:
    cluster:
      - host: "redis-1.example.com"
        port: 6379
      - host: "redis-2.example.com"
        port: 6379

monitoring:
  metrics:
    enabled: true
    port: 9090
  logging:
    level: "info"
    outputs:
      - type: "stdout"
      - type: "file"
        path: "/var/log/app.log"
EOT
}

data "yamlflattener_flatten" "app_config" {
  yaml_content = local.app_config_yaml
}

# Example 2: Using data source with external YAML file
data "yamlflattener_flatten" "external_config" {
  yaml_file = "${path.module}/config.yaml"
}

# Example 3: Using provider function for inline processing
locals {
  feature_flags = <<EOT
features:
  authentication:
    enabled: true
    providers:
      - name: "oauth2"
        config:
          client_id: "example-client"
          scopes: ["read", "write"]
      - name: "saml"
        config:
          entity_id: "example-entity"

  api_limits:
    rate_limit: 1000
    burst_limit: 100
    per_user_limit: 50

  experimental:
    new_ui: false
    beta_features: true
EOT

  flattened_features = provider::yamlflattener::flatten(local.feature_flags)

}

# Outputs demonstrating various use cases

# Basic configuration values
output "app_info" {
  description = "Application information from flattened YAML"
  value = {
    name        = data.yamlflattener_flatten.app_config.flattened["application.name"]
    version     = data.yamlflattener_flatten.app_config.flattened["application.version"]
    environment = data.yamlflattener_flatten.app_config.flattened["application.environment"]
  }
}

# Server configuration
output "server_config" {
  description = "Server configuration from flattened YAML"
  value = {
    host        = data.yamlflattener_flatten.app_config.flattened["server.host"]
    port        = data.yamlflattener_flatten.app_config.flattened["server.port"]
    ssl_enabled = data.yamlflattener_flatten.app_config.flattened["server.ssl.enabled"]
    ssl_cert    = data.yamlflattener_flatten.app_config.flattened["server.ssl.cert_path"]
  }
}

# Database configuration with primary and replicas
output "database_config" {
  description = "Database configuration including replicas"
  value = {
    primary = {
      host = data.yamlflattener_flatten.app_config.flattened["database.primary.host"]
      port = data.yamlflattener_flatten.app_config.flattened["database.primary.port"]
      name = data.yamlflattener_flatten.app_config.flattened["database.primary.name"]
    }
    replicas = [
      {
        host = data.yamlflattener_flatten.app_config.flattened["database.replicas[0].host"]
        port = data.yamlflattener_flatten.app_config.flattened["database.replicas[0].port"]
        name = data.yamlflattener_flatten.app_config.flattened["database.replicas[0].name"]
      },
      {
        host = data.yamlflattener_flatten.app_config.flattened["database.replicas[1].host"]
        port = data.yamlflattener_flatten.app_config.flattened["database.replicas[1].port"]
        name = data.yamlflattener_flatten.app_config.flattened["database.replicas[1].name"]
      }
    ]
  }
}

# Redis cluster endpoints
output "redis_endpoints" {
  description = "Redis cluster endpoints"
  value = [
    "${data.yamlflattener_flatten.app_config.flattened["cache.redis.cluster[0].host"]}:${data.yamlflattener_flatten.app_config.flattened["cache.redis.cluster[0].port"]}",
    "${data.yamlflattener_flatten.app_config.flattened["cache.redis.cluster[1].host"]}:${data.yamlflattener_flatten.app_config.flattened["cache.redis.cluster[1].port"]}"
  ]
}

# Feature flags using provider function
output "authentication_enabled" {
  description = "Authentication feature flag"
  value       = local.flattened_features["features.authentication.enabled"]
}

output "oauth_client_id" {
  description = "OAuth client ID from feature flags"
  value       = local.flattened_features["features.authentication.providers[0].config.client_id"]
}

output "api_rate_limit" {
  description = "API rate limit configuration"
  value = {
    rate_limit     = local.flattened_features["features.api_limits.rate_limit"]
    burst_limit    = local.flattened_features["features.api_limits.burst_limit"]
    per_user_limit = local.flattened_features["features.api_limits.per_user_limit"]
  }
}

# Demonstrate equivalence between data source and function
output "equivalence_test" {
  description = "Test that data source and function produce identical results"
  value = {
    datasource_result = data.yamlflattener_flatten.app_config.flattened["application.name"]
    function_result   = provider::yamlflattener::flatten(local.app_config_yaml)["application.name"]
    are_equal        = data.yamlflattener_flatten.app_config.flattened["application.name"] == provider::yamlflattener::flatten(local.app_config_yaml)["application.name"]
  }
}

# Show all flattened keys for debugging
output "all_flattened_keys" {
  description = "All flattened keys from the configuration"
  value       = keys(data.yamlflattener_flatten.app_config.flattened)
}

# Complex nested structure demonstration
output "monitoring_config" {
  description = "Monitoring configuration with nested arrays"
  value = {
    metrics_enabled = data.yamlflattener_flatten.app_config.flattened["monitoring.metrics.enabled"]
    metrics_port    = data.yamlflattener_flatten.app_config.flattened["monitoring.metrics.port"]
    log_level       = data.yamlflattener_flatten.app_config.flattened["monitoring.logging.level"]
    log_outputs = [
      {
        type = data.yamlflattener_flatten.app_config.flattened["monitoring.logging.outputs[0].type"]
      },
      {
        type = data.yamlflattener_flatten.app_config.flattened["monitoring.logging.outputs[1].type"]
        path = data.yamlflattener_flatten.app_config.flattened["monitoring.logging.outputs[1].path"]
      }
    ]
  }
}
