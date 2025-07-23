# Kubernetes Integration Example
# This example shows how to use the YAML Flattener provider with Kubernetes resources

terraform {
  required_version = ">= 1.0"
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 1.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "yamlflattener" {}

provider "kubernetes" {
  # Configure your Kubernetes provider here
  # config_path = "~/.kube/config"
}

# Flatten application configuration for Kubernetes deployment
data "yamlflattener_flatten" "app_config" {
  yaml_content = <<-EOT
    application:
      name: "web-app"
      version: "1.2.3"
      namespace: "production"

      image:
        repository: "nginx"
        tag: "1.21-alpine"
        pull_policy: "IfNotPresent"

      deployment:
        replicas: 3
        strategy:
          type: "RollingUpdate"
          rolling_update:
            max_surge: 1
            max_unavailable: 0

      resources:
        requests:
          cpu: "100m"
          memory: "128Mi"
        limits:
          cpu: "500m"
          memory: "512Mi"

      service:
        type: "ClusterIP"
        port: 80
        target_port: 80

      ingress:
        enabled: true
        host: "web-app.example.com"
        path: "/"
        tls:
          enabled: true
          secret_name: "web-app-tls"

      config:
        database:
          host: "postgres.database.svc.cluster.local"
          port: 5432
          name: "webapp"
          ssl_mode: "require"
        redis:
          host: "redis.cache.svc.cluster.local"
          port: 6379
          database: 0
        logging:
          level: "info"
          format: "json"
        features:
          authentication: true
          metrics: true
          tracing: false

      secrets:
        database:
          username: "webapp_user"
          # password would be injected from external secret management
        redis:
          password: ""
        jwt:
          secret_key: "your-jwt-secret-key"

      health_checks:
        liveness:
          path: "/health"
          port: 8080
          initial_delay: 30
          period: 10
        readiness:
          path: "/ready"
          port: 8080
          initial_delay: 5
          period: 5

      monitoring:
        prometheus:
          enabled: true
          port: 9090
          path: "/metrics"
        jaeger:
          enabled: false
          endpoint: "http://jaeger-collector:14268/api/traces"
  EOT
}

# Extract configuration into locals
locals {
  config = data.yamlflattener_flatten.app_config.flattened

  app_name      = local.config["application.name"]
  app_version   = local.config["application.version"]
  app_namespace = local.config["application.namespace"]

  # Image configuration
  image_repo        = local.config["application.image.repository"]
  image_tag         = local.config["application.image.tag"]
  image_pull_policy = local.config["application.image.pull_policy"]

  # Resource configuration
  replicas = tonumber(local.config["application.deployment.replicas"])

  # Service configuration
  service_type        = local.config["application.service.type"]
  service_port        = tonumber(local.config["application.service.port"])
  service_target_port = tonumber(local.config["application.service.target_port"])

  # Ingress configuration
  ingress_enabled    = tobool(local.config["application.ingress.enabled"])
  ingress_host       = local.config["application.ingress.host"]
  ingress_path       = local.config["application.ingress.path"]
  ingress_tls_enabled = tobool(local.config["application.ingress.tls.enabled"])
  ingress_tls_secret  = local.config["application.ingress.tls.secret_name"]

  # Health check configuration
  liveness_path          = local.config["application.health_checks.liveness.path"]
  liveness_port          = tonumber(local.config["application.health_checks.liveness.port"])
  liveness_initial_delay = tonumber(local.config["application.health_checks.liveness.initial_delay"])
  liveness_period        = tonumber(local.config["application.health_checks.liveness.period"])

  readiness_path          = local.config["application.health_checks.readiness.path"]
  readiness_port          = tonumber(local.config["application.health_checks.readiness.port"])
  readiness_initial_delay = tonumber(local.config["application.health_checks.readiness.initial_delay"])
  readiness_period        = tonumber(local.config["application.health_checks.readiness.period"])

  # Common labels
  common_labels = {
    app     = local.app_name
    version = local.app_version
  }
}

# Namespace
resource "kubernetes_namespace" "app" {
  metadata {
    name = local.app_namespace

    labels = merge(local.common_labels, {
      environment = "production"
    })
  }
}

# ConfigMap for application configuration
resource "kubernetes_config_map" "app_config" {
  metadata {
    name      = "${local.app_name}-config"
    namespace = kubernetes_namespace.app.metadata[0].name
    labels    = local.common_labels
  }

  data = {
    # Database configuration
    DATABASE_HOST     = local.config["application.config.database.host"]
    DATABASE_PORT     = local.config["application.config.database.port"]
    DATABASE_NAME     = local.config["application.config.database.name"]
    DATABASE_SSL_MODE = local.config["application.config.database.ssl_mode"]

    # Redis configuration
    REDIS_HOST     = local.config["application.config.redis.host"]
    REDIS_PORT     = local.config["application.config.redis.port"]
    REDIS_DATABASE = local.config["application.config.redis.database"]

    # Logging configuration
    LOG_LEVEL  = local.config["application.config.logging.level"]
    LOG_FORMAT = local.config["application.config.logging.format"]

    # Feature flags
    FEATURE_AUTHENTICATION = local.config["application.config.features.authentication"]
    FEATURE_METRICS        = local.config["application.config.features.metrics"]
    FEATURE_TRACING        = local.config["application.config.features.tracing"]

    # Monitoring configuration
    PROMETHEUS_ENABLED = local.config["application.monitoring.prometheus.enabled"]
    PROMETHEUS_PORT    = local.config["application.monitoring.prometheus.port"]
    PROMETHEUS_PATH    = local.config["application.monitoring.prometheus.path"]
  }
}

# Secret for sensitive configuration
resource "kubernetes_secret" "app_secrets" {
  metadata {
    name      = "${local.app_name}-secrets"
    namespace = kubernetes_namespace.app.metadata[0].name
    labels    = local.common_labels
  }

  data = {
    DATABASE_USERNAME = base64encode(local.config["application.secrets.database.username"])
    DATABASE_PASSWORD = base64encode("changeme123!") # In production, use external secret management
    REDIS_PASSWORD    = base64encode(local.config["application.secrets.redis.password"])
    JWT_SECRET_KEY    = base64encode(local.config["application.secrets.jwt.secret_key"])
  }

  type = "Opaque"
}

# Deployment
resource "kubernetes_deployment" "app" {
  metadata {
    name      = local.app_name
    namespace = kubernetes_namespace.app.metadata[0].name
    labels    = local.common_labels
  }

  spec {
    replicas = local.replicas

    selector {
      match_labels = local.common_labels
    }

    strategy {
      type = local.config["application.deployment.strategy.type"]

      rolling_update {
        max_surge       = local.config["application.deployment.strategy.rolling_update.max_surge"]
        max_unavailable = local.config["application.deployment.strategy.rolling_update.max_unavailable"]
      }
    }

    template {
      metadata {
        labels = local.common_labels

        annotations = {
          "prometheus.io/scrape" = local.config["application.monitoring.prometheus.enabled"]
          "prometheus.io/port"   = local.config["application.monitoring.prometheus.port"]
          "prometheus.io/path"   = local.config["application.monitoring.prometheus.path"]
        }
      }

      spec {
        container {
          name  = local.app_name
          image = "${local.image_repo}:${local.image_tag}"
          image_pull_policy = local.image_pull_policy

          port {
            container_port = local.service_target_port
            name          = "http"
          }

          port {
            container_port = local.liveness_port
            name          = "health"
          }

          # Environment variables from ConfigMap
          env_from {
            config_map_ref {
              name = kubernetes_config_map.app_config.metadata[0].name
            }
          }

          # Environment variables from Secret
          env_from {
            secret_ref {
              name = kubernetes_secret.app_secrets.metadata[0].name
            }
          }

          # Resource limits and requests
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

          # Liveness probe
          liveness_probe {
            http_get {
              path = local.liveness_path
              port = local.liveness_port
            }
            initial_delay_seconds = local.liveness_initial_delay
            period_seconds        = local.liveness_period
          }

          # Readiness probe
          readiness_probe {
            http_get {
              path = local.readiness_path
              port = local.readiness_port
            }
            initial_delay_seconds = local.readiness_initial_delay
            period_seconds        = local.readiness_period
          }
        }
      }
    }
  }
}

# Service
resource "kubernetes_service" "app" {
  metadata {
    name      = local.app_name
    namespace = kubernetes_namespace.app.metadata[0].name
    labels    = local.common_labels
  }

  spec {
    selector = local.common_labels

    port {
      name        = "http"
      port        = local.service_port
      target_port = local.service_target_port
      protocol    = "TCP"
    }

    type = local.service_type
  }
}

# Ingress (conditional based on configuration)
resource "kubernetes_ingress_v1" "app" {
  count = local.ingress_enabled ? 1 : 0

  metadata {
    name      = local.app_name
    namespace = kubernetes_namespace.app.metadata[0].name
    labels    = local.common_labels

    annotations = {
      "kubernetes.io/ingress.class"                = "nginx"
      "nginx.ingress.kubernetes.io/rewrite-target" = "/"
    }
  }

  spec {
    dynamic "tls" {
      for_each = local.ingress_tls_enabled ? [1] : []
      content {
        hosts       = [local.ingress_host]
        secret_name = local.ingress_tls_secret
      }
    }

    rule {
      host = local.ingress_host

      http {
        path {
          path      = local.ingress_path
          path_type = "Prefix"

          backend {
            service {
              name = kubernetes_service.app.metadata[0].name
              port {
                number = local.service_port
              }
            }
          }
        }
      }
    }
  }
}

# Horizontal Pod Autoscaler
resource "kubernetes_horizontal_pod_autoscaler_v2" "app" {
  metadata {
    name      = local.app_name
    namespace = kubernetes_namespace.app.metadata[0].name
    labels    = local.common_labels
  }

  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment.app.metadata[0].name
    }

    min_replicas = 2
    max_replicas = 10

    metric {
      type = "Resource"
      resource {
        name = "cpu"
        target {
          type                = "Utilization"
          average_utilization = 70
        }
      }
    }

    metric {
      type = "Resource"
      resource {
        name = "memory"
        target {
          type                = "Utilization"
          average_utilization = 80
        }
      }
    }
  }
}

# ServiceMonitor for Prometheus (if monitoring is enabled)
resource "kubernetes_manifest" "service_monitor" {
  count = tobool(local.config["application.monitoring.prometheus.enabled"]) ? 1 : 0

  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"

    metadata = {
      name      = local.app_name
      namespace = kubernetes_namespace.app.metadata[0].name
      labels    = local.common_labels
    }

    spec = {
      selector = {
        matchLabels = local.common_labels
      }

      endpoints = [
        {
          port = "http"
          path = local.config["application.monitoring.prometheus.path"]
        }
      ]
    }
  }
}

# Outputs
output "namespace" {
  description = "Application namespace"
  value       = kubernetes_namespace.app.metadata[0].name
}

output "service_name" {
  description = "Service name"
  value       = kubernetes_service.app.metadata[0].name
}

output "ingress_host" {
  description = "Ingress hostname"
  value       = local.ingress_enabled ? local.ingress_host : null
}

output "deployment_replicas" {
  description = "Number of deployment replicas"
  value       = local.replicas
}

output "flattened_config_sample" {
  description = "Sample of flattened configuration keys"
  value = {
    app_name           = local.config["application.name"]
    image_repository   = local.config["application.image.repository"]
    database_host      = local.config["application.config.database.host"]
    prometheus_enabled = local.config["application.monitoring.prometheus.enabled"]
    ingress_enabled    = local.config["application.ingress.enabled"]
  }
}

output "all_config_keys" {
  description = "All available configuration keys from flattened YAML"
  value       = keys(local.config)
}
