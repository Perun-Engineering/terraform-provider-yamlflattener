terraform {
  required_providers {
    yamlflattener = {
      source = "Perun-Engineering/yamlflattener"
    }
  }
}

# Complex alertmanager configuration example (from requirements)
data "yamlflattener_flatten" "alertmanager" {
  yaml_content = <<-EOT
    alertmanager:
      config:
        global:
          slack_api_url: "your-encrypted-slack-webhook"
        receivers:
          - name: "slack-notifications"
            slack_configs:
              - api_url: "your-encrypted-webhook-url"
                channel: "#alerts"
                title: "Alert: {{ .GroupLabels.alertname }}"
          - name: "email-notifications"
            email_configs:
              - to: "admin@example.com"
                subject: "Alert: {{ .GroupLabels.alertname }}"
        route:
          group_by: ['alertname']
          group_wait: 10s
          group_interval: 10s
          repeat_interval: 1h
          receiver: 'slack-notifications'
  EOT
}

# Use flattened values in other resources
output "slack_webhook" {
  value = data.yamlflattener_flatten.alertmanager.flattened["alertmanager.config.global.slack_api_url"]
}

output "first_receiver_name" {
  value = data.yamlflattener_flatten.alertmanager.flattened["alertmanager.config.receivers[0].name"]
}

output "all_flattened" {
  value = data.yamlflattener_flatten.alertmanager.flattened
}
