# Multi-Resource Integration Example
# This example demonstrates using the YAML Flattener provider with multiple types of Terraform resources

terraform {
  required_version = ">= 1.0"
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 1.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

provider "yamlflattener" {}

# Comprehensive application configuration
data "yamlflattener_flatten" "app_config" {
  yaml_content = <<-EOT
    application:
      name: "multi-tier-app"
      version: "2.1.0"
      environment: "production"

      components:
        frontend:
          type: "react"
          build_command: "npm run build"
          output_dir: "dist"
          env_vars:
            API_URL: "https://api.example.com"
            CDN_URL: "https://cdn.example.com"

        backend:
          type: "nodejs"
          port: 3000
          health_check: "/health"
          dependencies:
            - name: "express"
              version: "4.18.0"
            - name: "mongoose"
              version: "6.5.0"

        database:
          type: "mongodb"
          host: "mongodb.example.com"
          port: 27017
          name: "app_production"
          auth:
            username: "app_user"
            mechanism: "SCRAM-SHA-256"

        cache:
          type: "redis"
          cluster:
            - host: "redis-1.example.com"
              port: 6379
              role: "master"
            - host: "redis-2.example.com"
              port: 6379
              role: "slave"

    infrastructure:
      monitoring:
        prometheus:
          enabled: true
          port: 9090
          scrape_interval: "30s"
          retention: "15d"
        grafana:
          enabled: true
          port: 3001
          admin_user: "admin"

      logging:
        elasticsearch:
          enabled: true
          cluster_name: "app-logs"
          nodes:
            - host: "es-1.example.com"
              port: 9200
            - host: "es-2.example.com"
              port: 9200
        kibana:
          enabled: true
          port: 5601

      security:
        ssl:
          enabled: true
          cert_path: "/etc/ssl/certs"
          key_path: "/etc/ssl/private"
        firewall:
          rules:
            - name: "allow-http"
              port: 80
              protocol: "tcp"
              source: "0.0.0.0/0"
            - name: "allow-https"
              port: 443
              protocol: "tcp"
              source: "0.0.0.0/0"
            - name: "allow-ssh"
              port: 22
              protocol: "tcp"
              source: "10.0.0.0/8"

    deployment:
      strategy: "blue-green"
      rollback:
        enabled: true
        max_revisions: 5
      health_checks:
        startup_timeout: 300
        readiness_timeout: 30
        liveness_interval: 10

    backup:
      database:
        enabled: true
        schedule: "0 2 * * *"
        retention_days: 30
        storage:
          type: "s3"
          bucket: "app-backups-prod"
          region: "us-west-2"
      files:
        enabled: true
        paths:
          - "/var/app/uploads"
          - "/var/app/logs"
        schedule: "0 3 * * *"
  EOT
}

# Extract configuration into organized locals
locals {
  config = data.yamlflattener_flatten.app_config.flattened

  # Application configuration
  app = {
    name        = local.config["application.name"]
    version     = local.config["application.version"]
    environment = local.config["application.environment"]
  }

  # Frontend configuration
  frontend = {
    type         = local.config["application.components.frontend.type"]
    build_cmd    = local.config["application.components.frontend.build_command"]
    output_dir   = local.config["application.components.frontend.output_dir"]
    api_url      = local.config["application.components.frontend.env_vars.API_URL"]
    cdn_url      = local.config["application.components.frontend.env_vars.CDN_URL"]
  }

  # Backend configuration
  backend = {
    type         = local.config["application.components.backend.type"]
    port         = tonumber(local.config["application.components.backend.port"])
    health_check = local.config["application.components.backend.health_check"]
  }

  # Database configuration
  database = {
    type      = local.config["application.components.database.type"]
    host      = local.config["application.components.database.host"]
    port      = tonumber(local.config["application.components.database.port"])
    name      = local.config["application.components.database.name"]
    username  = local.config["application.components.database.auth.username"]
    mechanism = local.config["application.components.database.auth.mechanism"]
  }

  # Redis cluster configuration
  redis_nodes = [
    {
      host = local.config["application.components.cache.cluster[0].host"]
      port = tonumber(local.config["application.components.cache.cluster[0].port"])
      role = local.config["application.components.cache.cluster[0].role"]
    },
    {
      host = local.config["application.components.cache.cluster[1].host"]
      port = tonumber(local.config["application.components.cache.cluster[1].port"])
      role = local.config["application.components.cache.cluster[1].role"]
    }
  ]

  # Monitoring configuration
  monitoring = {
    prometheus_enabled = tobool(local.config["infrastructure.monitoring.prometheus.enabled"])
    prometheus_port    = tonumber(local.config["infrastructure.monitoring.prometheus.port"])
    grafana_enabled    = tobool(local.config["infrastructure.monitoring.grafana.enabled"])
    grafana_port       = tonumber(local.config["infrastructure.monitoring.grafana.port"])
  }

  # Security configuration
  security = {
    ssl_enabled = tobool(local.config["infrastructure.security.ssl.enabled"])
    cert_path   = local.config["infrastructure.security.ssl.cert_path"]
    key_path    = local.config["infrastructure.security.ssl.key_path"]
  }

  # Backup configuration
  backup = {
    db_enabled      = tobool(local.config["backup.database.enabled"])
    db_schedule     = local.config["backup.database.schedule"]
    db_retention    = tonumber(local.config["backup.database.retention_days"])
    storage_bucket  = local.config["backup.database.storage.bucket"]
    storage_region  = local.config["backup.database.storage.region"]
  }
}

# Generate application configuration files
resource "local_file" "frontend_env" {
  filename = "${path.module}/generated/frontend/.env"
  content = templatefile("${path.module}/templates/frontend.env.tpl", {
    api_url = local.frontend.api_url
    cdn_url = local.frontend.cdn_url
    version = local.app.version
    environment = local.app.environment
  })

  depends_on = [null_resource.create_directories]
}

resource "local_file" "backend_config" {
  filename = "${path.module}/generated/backend/config.json"
  content = jsonencode({
    app = {
      name        = local.app.name
      version     = local.app.version
      environment = local.app.environment
      port        = local.backend.port
      health_check = local.backend.health_check
    }
    database = {
      type     = local.database.type
      host     = local.database.host
      port     = local.database.port
      name     = local.database.name
      username = local.database.username
      auth_mechanism = local.database.mechanism
    }
    cache = {
      type  = "redis"
      nodes = local.redis_nodes
    }
    monitoring = {
      prometheus_enabled = local.monitoring.prometheus_enabled
      prometheus_port    = local.monitoring.prometheus_port
    }
  })

  depends_on = [null_resource.create_directories]
}

# Generate Docker Compose configuration
resource "local_file" "docker_compose" {
  filename = "${path.module}/generated/docker-compose.yml"
  content = templatefile("${path.module}/templates/docker-compose.yml.tpl", {
    app_name     = local.app.name
    app_version  = local.app.version
    backend_port = local.backend.port
    db_host      = local.database.host
    db_port      = local.database.port
    db_name      = local.database.name
    redis_master = local.redis_nodes[0]
    redis_slave  = local.redis_nodes[1]
    prometheus_enabled = local.monitoring.prometheus_enabled
    prometheus_port    = local.monitoring.prometheus_port
    grafana_enabled    = local.monitoring.grafana_enabled
    grafana_port       = local.monitoring.grafana_port
  })

  depends_on = [null_resource.create_directories]
}

# Generate Kubernetes manifests
resource "local_file" "k8s_configmap" {
  filename = "${path.module}/generated/k8s/configmap.yaml"
  content = templatefile("${path.module}/templates/k8s-configmap.yaml.tpl", {
    app_name    = local.app.name
    namespace   = "${local.app.name}-${local.app.environment}"
    config_data = {
      APP_NAME        = local.app.name
      APP_VERSION     = local.app.version
      APP_ENVIRONMENT = local.app.environment
      BACKEND_PORT    = tostring(local.backend.port)
      DB_HOST         = local.database.host
      DB_PORT         = tostring(local.database.port)
      DB_NAME         = local.database.name
      REDIS_MASTER    = "${local.redis_nodes[0].host}:${local.redis_nodes[0].port}"
      REDIS_SLAVE     = "${local.redis_nodes[1].host}:${local.redis_nodes[1].port}"
    }
  })

  depends_on = [null_resource.create_directories]
}

# Generate monitoring configuration
resource "local_file" "prometheus_config" {
  count    = local.monitoring.prometheus_enabled ? 1 : 0
  filename = "${path.module}/generated/monitoring/prometheus.yml"
  content = templatefile("${path.module}/templates/prometheus.yml.tpl", {
    scrape_interval = local.config["infrastructure.monitoring.prometheus.scrape_interval"]
    retention       = local.config["infrastructure.monitoring.prometheus.retention"]
    targets = [
      "${local.backend.port}${local.backend.health_check}",
      "${local.database.host}:${local.database.port}",
      "${local.redis_nodes[0].host}:${local.redis_nodes[0].port}",
      "${local.redis_nodes[1].host}:${local.redis_nodes[1].port}"
    ]
  })

  depends_on = [null_resource.create_directories]
}

# Generate backup scripts
resource "local_file" "backup_script" {
  count    = local.backup.db_enabled ? 1 : 0
  filename = "${path.module}/generated/scripts/backup.sh"
  content = templatefile("${path.module}/templates/backup.sh.tpl", {
    db_host     = local.database.host
    db_port     = local.database.port
    db_name     = local.database.name
    db_username = local.database.username
    schedule    = local.backup.db_schedule
    retention   = local.backup.db_retention
    bucket      = local.backup.storage_bucket
    region      = local.backup.storage_region
  })

  file_permission = "0755"
  depends_on = [null_resource.create_directories]
}

# Generate security configuration
resource "local_file" "nginx_ssl_config" {
  count    = local.security.ssl_enabled ? 1 : 0
  filename = "${path.module}/generated/nginx/ssl.conf"
  content = templatefile("${path.module}/templates/nginx-ssl.conf.tpl", {
    cert_path = local.security.cert_path
    key_path  = local.security.key_path
    app_name  = local.app.name
    backend_port = local.backend.port
  })

  depends_on = [null_resource.create_directories]
}

# Generate firewall rules
resource "local_file" "firewall_rules" {
  filename = "${path.module}/generated/security/firewall.rules"
  content = templatefile("${path.module}/templates/firewall.rules.tpl", {
    rules = [
      {
        name     = local.config["infrastructure.security.firewall.rules[0].name"]
        port     = tonumber(local.config["infrastructure.security.firewall.rules[0].port"])
        protocol = local.config["infrastructure.security.firewall.rules[0].protocol"]
        source   = local.config["infrastructure.security.firewall.rules[0].source"]
      },
      {
        name     = local.config["infrastructure.security.firewall.rules[1].name"]
        port     = tonumber(local.config["infrastructure.security.firewall.rules[1].port"])
        protocol = local.config["infrastructure.security.firewall.rules[1].protocol"]
        source   = local.config["infrastructure.security.firewall.rules[1].source"]
      },
      {
        name     = local.config["infrastructure.security.firewall.rules[2].name"]
        port     = tonumber(local.config["infrastructure.security.firewall.rules[2].port"])
        protocol = local.config["infrastructure.security.firewall.rules[2].protocol"]
        source   = local.config["infrastructure.security.firewall.rules[2].source"]
      }
    ]
  })

  depends_on = [null_resource.create_directories]
}

# Create directory structure
resource "null_resource" "create_directories" {
  provisioner "local-exec" {
    command = <<-EOT
      mkdir -p ${path.module}/generated/{frontend,backend,k8s,monitoring,scripts,nginx,security}
      mkdir -p ${path.module}/templates
    EOT
  }
}

# Create template files if they don't exist
resource "local_file" "frontend_env_template" {
  filename = "${path.module}/templates/frontend.env.tpl"
  content = <<-EOT
REACT_APP_API_URL=${api_url}
REACT_APP_CDN_URL=${cdn_url}
REACT_APP_VERSION=${version}
REACT_APP_ENVIRONMENT=${environment}
EOT
}

resource "local_file" "docker_compose_template" {
  filename = "${path.module}/templates/docker-compose.yml.tpl"
  content = <<-EOT
version: '3.8'
services:
  ${app_name}-backend:
    image: ${app_name}:${app_version}
    ports:
      - "${backend_port}:${backend_port}"
    environment:
      - DB_HOST=${db_host}
      - DB_PORT=${db_port}
      - DB_NAME=${db_name}
      - REDIS_MASTER=${redis_master.host}:${redis_master.port}
      - REDIS_SLAVE=${redis_slave.host}:${redis_slave.port}
    %{ if prometheus_enabled }
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "${prometheus_port}:9090"
    %{ endif }
    %{ if grafana_enabled }
  grafana:
    image: grafana/grafana:latest
    ports:
      - "${grafana_port}:3000"
    %{ endif }
EOT
}

resource "local_file" "k8s_configmap_template" {
  filename = "${path.module}/templates/k8s-configmap.yaml.tpl"
  content = <<-EOT
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${app_name}-config
  namespace: ${namespace}
data:
%{ for key, value in config_data ~}
  ${key}: "${value}"
%{ endfor ~}
EOT
}

# Validation resource to ensure configuration is complete
resource "null_resource" "validate_config" {
  provisioner "local-exec" {
    command = <<-EOT
      echo "Validating configuration..."
      echo "Application: ${local.app.name} v${local.app.version}"
      echo "Environment: ${local.app.environment}"
      echo "Backend port: ${local.backend.port}"
      echo "Database: ${local.database.type} at ${local.database.host}:${local.database.port}"
      echo "Redis nodes: ${length(local.redis_nodes)}"
      echo "Monitoring enabled: ${local.monitoring.prometheus_enabled}"
      echo "SSL enabled: ${local.security.ssl_enabled}"
      echo "Backup enabled: ${local.backup.db_enabled}"
      echo "Configuration validation complete!"
    EOT
  }

  depends_on = [
    local_file.frontend_env,
    local_file.backend_config,
    local_file.docker_compose,
    local_file.k8s_configmap
  ]
}

# Outputs
output "application_info" {
  description = "Application information from flattened YAML"
  value = {
    name        = local.app.name
    version     = local.app.version
    environment = local.app.environment
  }
}

output "service_endpoints" {
  description = "Service endpoints configuration"
  value = {
    backend  = "http://localhost:${local.backend.port}"
    database = "${local.database.host}:${local.database.port}"
    redis_master = "${local.redis_nodes[0].host}:${local.redis_nodes[0].port}"
    redis_slave  = "${local.redis_nodes[1].host}:${local.redis_nodes[1].port}"
    prometheus = local.monitoring.prometheus_enabled ? "http://localhost:${local.monitoring.prometheus_port}" : null
    grafana    = local.monitoring.grafana_enabled ? "http://localhost:${local.monitoring.grafana_port}" : null
  }
}

output "generated_files" {
  description = "List of generated configuration files"
  value = [
    local_file.frontend_env.filename,
    local_file.backend_config.filename,
    local_file.docker_compose.filename,
    local_file.k8s_configmap.filename,
  ]
}

output "flattened_keys_sample" {
  description = "Sample of available flattened configuration keys"
  value = [
    "application.name",
    "application.components.backend.port",
    "application.components.database.host",
    "infrastructure.monitoring.prometheus.enabled",
    "backup.database.schedule",
    "infrastructure.security.ssl.enabled"
  ]
}

output "configuration_summary" {
  description = "Summary of configuration extracted from YAML"
  value = {
    total_keys = length(keys(local.config))
    components = {
      frontend_configured  = contains(keys(local.config), "application.components.frontend.type")
      backend_configured   = contains(keys(local.config), "application.components.backend.type")
      database_configured  = contains(keys(local.config), "application.components.database.type")
      cache_configured     = contains(keys(local.config), "application.components.cache.type")
    }
    features = {
      monitoring_enabled = local.monitoring.prometheus_enabled
      ssl_enabled       = local.security.ssl_enabled
      backup_enabled    = local.backup.db_enabled
    }
  }
}
