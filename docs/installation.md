# Installation Guide

This guide covers how to install and configure the YAML Flattener provider in your Terraform projects.

## Requirements

- Terraform >= 1.0
- Go >= 1.21 (for building from source)

## Installation Methods

### 1. Terraform Registry (Recommended)

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

Then run:

```bash
terraform init
```

### 2. Local Installation

For development or testing purposes, you can install the provider locally:

1. Download the appropriate binary for your platform from the [releases page](https://github.com/Perun-Engineering/terraform-provider-yamlflattener/releases)

2. Create the provider directory structure:
   ```bash
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/perun-engineering/yamlflattener/1.0.0/darwin_amd64
   ```

3. Move the binary to the directory:
   ```bash
   mv terraform-provider-yamlflattener ~/.terraform.d/plugins/registry.terraform.io/perun-engineering/yamlflattener/1.0.0/darwin_amd64/
   ```

4. Make it executable:
   ```bash
   chmod +x ~/.terraform.d/plugins/registry.terraform.io/perun-engineering/yamlflattener/1.0.0/darwin_amd64/terraform-provider-yamlflattener
   ```

### 3. Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/Perun-Engineering/terraform-provider-yamlflattener.git
   cd terraform-provider-yamlflattener
   ```

2. Build the provider:
   ```bash
   go build -o terraform-provider-yamlflattener
   ```

3. Install locally following the steps in method 2.

## Provider Configuration

The provider doesn't require any configuration parameters by default:

```hcl
provider "yamlflattener" {}
```

### Optional Configuration

While the provider works without configuration, you can optionally set these parameters:

```hcl
provider "yamlflattener" {
  # Maximum nesting depth (default: 100)
  max_depth = 50

  # Maximum file size in bytes (default: 10MB)
  max_file_size = 5242880  # 5MB

  # Maximum number of flattened keys (default: 100,000)
  max_keys = 50000
}
```

## Verification

To verify the provider is installed correctly, create a simple test configuration:

```hcl
terraform {
  required_providers {
    yamlflattener = {
      source = "Perun-Engineering/yamlflattener"
    }
  }
}

provider "yamlflattener" {}

data "yamlflattener_flatten" "test" {
  yaml_content = <<-EOT
    test:
      key: "value"
  EOT
}

output "test_result" {
  value = data.yamlflattener_flatten.test.flattened["test.key"]
}
```

Run:

```bash
terraform init
terraform plan
```

If successful, you should see the provider initialize and the plan show the expected output.

## Supported Platforms

The provider is built for the following platforms:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Upgrading

To upgrade to a newer version:

1. Update the version constraint in your configuration:
   ```hcl
   terraform {
     required_providers {
       yamlflattener = {
         source  = "Perun-Engineering/yamlflattener"
         version = "~> 1.1"
       }
     }
   }
   ```

2. Run:
   ```bash
   terraform init -upgrade
   ```

## Troubleshooting Installation

### Provider Not Found

If you get a "provider not found" error:

1. Check that the provider source is correct
2. Verify your Terraform version is >= 1.0
3. Run `terraform init` to download the provider

### Permission Errors

If you encounter permission errors during local installation:

1. Ensure the binary is executable: `chmod +x terraform-provider-yamlflattener`
2. Check directory permissions for `~/.terraform.d/plugins/`
3. On macOS, you may need to allow the binary in Security & Privacy settings

### Version Conflicts

If you have version conflicts:

1. Clear the provider cache: `rm -rf .terraform/`
2. Run `terraform init` again
3. Check for version constraints in your configuration

For more troubleshooting help, see the [Troubleshooting Guide](troubleshooting.md).
