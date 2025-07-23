# Integration Examples

This directory contains examples showing how to integrate the YAML Flattener provider with other Terraform resources and providers. These examples demonstrate real-world use cases where YAML configuration files need to be processed and used with cloud resources.

## Examples Overview

1. **AWS Integration** - Using flattened YAML with AWS resources
2. **Kubernetes Integration** - Processing Kubernetes configurations
3. **Azure Integration** - Azure resource configuration from YAML
4. **Multi-Cloud Integration** - Using the same YAML across multiple providers
5. **CI/CD Integration** - Processing CI/CD pipeline configurations

## Common Patterns

### Pattern 1: Configuration-Driven Infrastructure

Use YAML files to define infrastructure parameters and flatten them for use in Terraform resources:

```hcl
# config.yaml
infrastructure:
  aws:
    region: "us-west-2"
    instance_type: "t3.medium"
    vpc_cidr: "10.0.0.0/16"
  database:
    engine: "postgres"
    version: "13.7"
    instance_class: "db.t3.micro"

# main.tf
data "yamlflattener_flatten" "config" {
  yaml_file = "${path.module}/config.yaml"
}

resource "aws_instance" "app" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = data.yamlflattener_flatten.config.flattened["infrastructure.aws.instance_type"]

  tags = {
    Environment = data.yamlflattener_flatten.config.flattened["infrastructure.environment"]
  }
}
```

### Pattern 2: Multi-Environment Configuration

Use the same YAML structure across different environments:

```hcl
# environments/prod.yaml, environments/dev.yaml, etc.
data "yamlflattener_flatten" "env_config" {
  yaml_file = "${path.module}/environments/${var.environment}.yaml"
}

locals {
  config = data.yamlflattener_flatten.env_config.flattened
}
```

### Pattern 3: Application Configuration Processing

Process application configuration files for use in container deployments:

```hcl
data "yamlflattener_flatten" "app_config" {
  yaml_file = "${path.module}/app-config.yaml"
}

resource "kubernetes_config_map" "app_config" {
  metadata {
    name = "app-config"
  }

  data = {
    for key, value in data.yamlflattener_flatten.app_config.flattened :
    replace(key, ".", "_") => value
  }
}
```

## Running the Examples

Each example directory contains:
- `main.tf` - Main Terraform configuration
- `variables.tf` - Input variables
- `outputs.tf` - Output values
- `README.md` - Specific instructions
- Sample YAML files

To run an example:

```bash
cd examples/integration-examples/aws-integration
terraform init
terraform plan
terraform apply
```

## Prerequisites

Different examples may require:
- AWS CLI configured
- kubectl configured for Kubernetes
- Azure CLI configured
- Appropriate provider credentials

Check each example's README for specific requirements.
