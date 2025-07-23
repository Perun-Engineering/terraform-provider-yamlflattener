# Alertmanager YAML Flattening Example

This example demonstrates how to use the YAML Flattener provider to flatten an Alertmanager configuration, which is a common use case for this provider. The example shows both the data source and function approaches.

## Alertmanager Configuration

The example uses this Alertmanager configuration:

```yaml
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
```

## Usage

```hcl
# Using the data source approach
data "yamlflattener_flatten" "alertmanager_ds" {
  yaml_file = "${path.module}/alertmanager.yaml"
}

# Using the function approach
locals {
  alertmanager_yaml = file("${path.module}/alertmanager.yaml")
  alertmanager_flattened = yamlflattener_flatten(local.alertmanager_yaml)
}

# Accessing values (both approaches produce identical results)
output "slack_webhook_ds" {
  value = data.yamlflattener_flatten.alertmanager_ds.flattened["alertmanager.config.global.slack_api_url"]
}

output "slack_webhook_fn" {
  value = local.alertmanager_flattened["alertmanager.config.global.slack_api_url"]
}
```

## Running the Example

To run this example, execute:

```bash
# Initialize Terraform
terraform init

# Apply the configuration
terraform apply
```

The outputs will show the flattened Alertmanager configuration and demonstrate how to access specific values from the flattened map.
