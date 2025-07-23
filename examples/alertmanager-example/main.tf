terraform {
  required_providers {
    yamlflattener = {
      source = "local/yamlflattener"
    }
  }
}

provider "yamlflattener" {}

# Using the data source approach with inline YAML content
data "yamlflattener_flatten" "alertmanager_ds_inline" {
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
            send_resolved: true
  EOT
}

# Using the function approach with the same YAML content
locals {
  alertmanager_yaml = <<-EOT
alertmanager:
  config:
    global:
      slack_api_url: "your-encrypted-slack-webhook"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "your-encrypted-webhook-url"
            channel: "#alerts"
            send_resolved: true
  EOT

  alertmanager_flattened = yamlflattener_flatten(local.alertmanager_yaml)
}

# Output the entire flattened structure from both approaches
output "alertmanager_flattened_ds" {
  value = data.yamlflattener_flatten.alertmanager_ds_inline.flattened
  description = "The complete flattened Alertmanager configuration using the data source approach"
}

output "alertmanager_flattened_fn" {
  value = local.alertmanager_flattened
  description = "The complete flattened Alertmanager configuration using the function approach"
}

# Output specific values to demonstrate accessing individual flattened keys
output "slack_webhook_ds" {
  value = data.yamlflattener_flatten.alertmanager_ds_inline.flattened["alertmanager.config.global.slack_api_url"]
  description = "The Slack webhook URL from the data source approach"
}

output "slack_webhook_fn" {
  value = local.alertmanager_flattened["alertmanager.config.global.slack_api_url"]
  description = "The Slack webhook URL from the function approach"
}

output "receiver_name_ds" {
  value = data.yamlflattener_flatten.alertmanager_ds_inline.flattened["alertmanager.config.receivers[0].name"]
  description = "The receiver name from the data source approach"
}

output "receiver_name_fn" {
  value = local.alertmanager_flattened["alertmanager.config.receivers[0].name"]
  description = "The receiver name from the function approach"
}

output "slack_channel_ds" {
  value = data.yamlflattener_flatten.alertmanager_ds_inline.flattened["alertmanager.config.receivers[0].slack_configs[0].channel"]
  description = "The Slack channel from the data source approach"
}

output "slack_channel_fn" {
  value = local.alertmanager_flattened["alertmanager.config.receivers[0].slack_configs[0].channel"]
  description = "The Slack channel from the function approach"
}

output "send_resolved_ds" {
  value = data.yamlflattener_flatten.alertmanager_ds_inline.flattened["alertmanager.config.receivers[0].slack_configs[0].send_resolved"]
  description = "The send_resolved flag from the data source approach"
}

output "send_resolved_fn" {
  value = local.alertmanager_flattened["alertmanager.config.receivers[0].slack_configs[0].send_resolved"]
  description = "The send_resolved flag from the function approach"
}
