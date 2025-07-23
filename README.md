# Terraform YAML Flattener Provider

A Terraform provider that flattens nested YAML structures into flat key-value maps with dot notation for objects and bracket notation for arrays.

## Features

- **Recursive Flattening**: Handles deeply nested YAML structures
- **Multiple Input Methods**: Accept YAML content directly or from files
- **Dual Interface**: Available as both data source and provider function
- **Type Preservation**: Converts all values to strings while maintaining semantic meaning
- **Error Handling**: Clear error messages for invalid YAML or file issues

## Installation

### Terraform Registry

```hcl
terraform {
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 1.0"
    }
  }
}
```

### Manual Installation

Download the appropriate binary for your platform from the [releases page](https://github.com/Perun-Engineering/terraform-provider-yamlflattener/releases) and place it in your Terraform plugins directory.

## Usage

### Data Source

```hcl
# Using YAML content directly
data "yamlflattener_flatten" "example" {
  yaml_content = <<EOF
alertmanager:
  config:
    global:
      slack_api_url: "your-encrypted-slack-webhook"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "your-encrypted-webhook-url"
EOF
}

# Using YAML file
data "yamlflattener_flatten" "from_file" {
  yaml_file = "config.yaml"
}

# Access flattened values
output "flattened_config" {
  value = data.yamlflattener_flatten.example.flattened
}
```

### Provider Function

```hcl
locals {
  yaml_content = <<EOF
database:
  host: "localhost"
  port: 5432
  credentials:
    username: "admin"
    password: "secret"
EOF

  flattened = provider::yamlflattener::flatten(local.yaml_content)
}

# Use flattened values in resources
resource "aws_ssm_parameter" "db_host" {
  name  = "/app/db/host"
  value = local.flattened["database.host"]
  type  = "String"
}

resource "aws_ssm_parameter" "db_port" {
  name  = "/app/db/port"
  value = local.flattened["database.port"]
  type  = "String"
}
```

## Examples

### Simple Object Flattening

**Input:**
```yaml
key1: value1
key2:
  nested: value2
```

**Output:**
```json
{
  "key1": "value1",
  "key2.nested": "value2"
}
```

### Array Flattening

**Input:**
```yaml
items:
  - name: item1
    value: 100
  - name: item2
    value: 200
```

**Output:**
```json
{
  "items[0].name": "item1",
  "items[0].value": "100",
  "items[1].name": "item2",
  "items[1].value": "200"
}
```

### Complex Mixed Structure

**Input:**
```yaml
alertmanager:
  config:
    global:
      slack_api_url: "webhook-url"
    receivers:
      - name: "slack-notifications"
        slack_configs:
          - api_url: "webhook-url"
            channel: "#alerts"
```

**Output:**
```json
{
  "alertmanager.config.global.slack_api_url": "webhook-url",
  "alertmanager.config.receivers[0].name": "slack-notifications",
  "alertmanager.config.receivers[0].slack_configs[0].api_url": "webhook-url",
  "alertmanager.config.receivers[0].slack_configs[0].channel": "#alerts"
}
```

## Data Types

The provider handles various YAML data types:

- **Strings**: Preserved as-is
- **Numbers**: Converted to string representation
- **Booleans**: Converted to "true" or "false"
- **Null values**: Represented as empty strings
- **Arrays**: Indexed with bracket notation `[0]`, `[1]`, etc.
- **Objects**: Flattened with dot notation

## Error Handling

The provider provides clear error messages for common issues:

- Invalid YAML syntax
- File not found or permission errors
- Conflicting input parameters
- Deeply nested structures that exceed limits

## Monitoring and Maintenance

The provider includes automated monitoring and maintenance workflows to ensure reliability and security:

### Automated Testing

- **Pull Request Validation**: Comprehensive testing on every PR including:
  - Unit tests with code coverage reporting
  - Integration tests with example configurations
  - Code quality checks with golangci-lint
  - Security scanning with govulncheck and gosec

### Security Monitoring

- **Weekly Vulnerability Scanning**: Automated checks for:
  - Known vulnerabilities in Go dependencies
  - Static code analysis for security issues
  - Supply chain security verification

- **Automated Security Updates**:
  - Monthly dependency audits
  - Automatic PRs for security-related dependency updates
  - Vulnerability reporting and tracking

### Build Notifications

- **Failure Alerts**: Automated notifications for workflow failures via:
  - Slack notifications (when configured)
  - Email alerts (when configured)
  - GitHub issues for critical failures

### Usage Analytics

- **Download Metrics Tracking**:
  - Monthly collection of Terraform Registry download statistics
  - Growth trend analysis
  - Usage reports for maintainers

### Configuration

To enable notifications:

1. **Slack Notifications**:
   - Add `SLACK_WEBHOOK_URL` secret to your GitHub repository

2. **Email Notifications**:
   - Add the following secrets to your GitHub repository:
     - `EMAIL_NOTIFICATION_ENABLED`: Set to "true"
     - `MAIL_SERVER`: SMTP server address
     - `MAIL_PORT`: SMTP server port
     - `MAIL_USERNAME`: SMTP username
     - `MAIL_PASSWORD`: SMTP password
     - `NOTIFICATION_EMAIL`: Recipient email address

## Development

### Prerequisites

- Go 1.21 or later
- Terraform 1.0 or later
- Git
- golangci-lint (for code quality checks)
- GNU Make (optional, for using the Makefile)

### Development Environment Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/Perun-Engineering/terraform-provider-yamlflattener.git
   cd terraform-provider-yamlflattener
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   ```

3. **Install Development Tools**

   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Install other tools as needed
   go install github.com/goreleaser/goreleaser@latest
   ```

4. **Build the Provider**

   ```bash
   go build -o terraform-provider-yamlflattener
   ```

5. **Install the Provider Locally for Testing**

   ```bash
   # For Terraform 0.13+
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/Perun-Engineering/yamlflattener/1.0.0/$(go env GOOS)_$(go env GOARCH)
   cp terraform-provider-yamlflattener ~/.terraform.d/plugins/registry.terraform.io/Perun-Engineering/yamlflattener/1.0.0/$(go env GOOS)_$(go env GOARCH)/
   ```

### Development Workflow

1. **Create a Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**

   - Follow the code structure and patterns in existing files
   - Add tests for new functionality
   - Update documentation as needed

3. **Run Tests**

   ```bash
   # Run unit tests
   go test ./...

   # Run specific tests
   go test ./internal/flattener -v

   # Run acceptance tests (requires Terraform)
   TF_ACC=1 go test ./internal/provider -v

   # Run tests with coverage
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out
   ```

4. **Run Linting**

   ```bash
   golangci-lint run
   ```

5. **Format Code**

   ```bash
   go fmt ./...
   ```

6. **Commit Your Changes**

   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   We follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `test:` for test changes
   - `refactor:` for code refactoring
   - `chore:` for routine tasks and maintenance

7. **Push Your Changes**

   ```bash
   git push origin feature/your-feature-name
   ```

8. **Create a Pull Request**

   - Go to the repository on GitHub
   - Click "New Pull Request"
   - Select your branch
   - Fill out the PR template with details about your changes

### Testing with Example Configurations

To test the provider with the example configurations:

```bash
# Build and install the provider locally
go build -o terraform-provider-yamlflattener

# Navigate to an example directory
cd examples/complete-example

# Initialize Terraform with local provider
terraform init

# Apply the configuration
terraform apply
```

### Debugging

For debugging the provider:

```bash
# Enable Terraform provider logs
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log

# Run Terraform with debugging enabled
terraform apply
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming and inclusive environment for all contributors.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- [GitHub Issues](https://github.com/Perun-Engineering/terraform-provider-yamlflattener/issues) for bug reports and feature requests
- [Terraform Registry Documentation](https://registry.terraform.io/providers/Perun-Engineering/yamlflattener/latest/docs) for detailed usage information

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.
