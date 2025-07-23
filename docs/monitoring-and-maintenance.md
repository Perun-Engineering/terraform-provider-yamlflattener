# Monitoring and Maintenance

This document describes the automated monitoring and maintenance workflows implemented for the Terraform YAML Flattener Provider.

## Automated Testing

### Pull Request Validation

Every pull request to the main branch triggers a comprehensive validation workflow:

- **Build and Test**: Compiles the provider and runs unit tests with coverage reporting
- **Linting**: Runs golangci-lint to ensure code quality
- **Security Scanning**: Checks for vulnerabilities using govulncheck and gosec
- **Integration Testing**: Validates the provider with example Terraform configurations

### Cross-Platform Testing

The build workflow tests the provider on multiple platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Security Monitoring

### Vulnerability Scanning

Weekly automated security scans check for:

- **Dependency Vulnerabilities**: Using govulncheck to identify known vulnerabilities in Go dependencies
- **Static Code Analysis**: Using gosec to detect potential security issues in the codebase
- **Supply Chain Security**: Verifying dependency integrity

### Security Updates

Monthly security update workflows:

1. **Security Audit**:
   - Scans dependencies for vulnerabilities using Nancy and govulncheck
   - Creates GitHub issues for detected vulnerabilities
   - Generates detailed reports for maintainer review

2. **Dependency Updates**:
   - Identifies outdated dependencies
   - Creates pull requests with security-related updates
   - Automatically runs tests on the updated dependencies

## Build Notifications

Automated notifications are sent when critical workflows fail:

### Notification Channels

- **Slack**: Real-time alerts with workflow details and links to logs
- **Email**: Detailed failure reports sent to designated recipients
- **GitHub Issues**: Automatically created for critical failures (release and security workflows)

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

## Usage Analytics

Monthly collection and analysis of provider usage metrics:

### Download Metrics

- **Collection**: Gathers download statistics from the Terraform Registry
- **Reporting**: Generates monthly reports with download counts and growth trends
- **Visualization**: Creates charts for download trends over time

### Metrics Dashboard

The metrics tracking workflow:
1. Fetches download data from the Terraform Registry
2. Generates CSV reports and visualization charts
3. Creates a GitHub issue with a summary of key metrics
4. Stores detailed reports as workflow artifacts

## Maintenance Workflows

### Automated Dependency Management

- **Dependabot**: Weekly checks for dependency updates
- **Auto-merge**: Non-breaking updates can be automatically merged
- **Security Prioritization**: Security-related updates are labeled for priority review

### Release Process

The automated release process:
1. Triggered by creating a version tag (v*)
2. Builds and signs provider binaries for all supported platforms
3. Generates checksums and signatures
4. Creates GitHub release with release notes from CHANGELOG.md
5. Publishes the provider to the Terraform Registry

## Workflow Files

The monitoring and maintenance automation is implemented in these GitHub Actions workflow files:

- `.github/workflows/build.yml`: Multi-architecture build workflow
- `.github/workflows/pr-validation.yml`: Pull request validation
- `.github/workflows/security.yml`: Security scanning
- `.github/workflows/release.yml`: Release and publishing
- `.github/workflows/notifications.yml`: Build failure notifications
- `.github/workflows/security-updates.yml`: Automated security updates
- `.github/workflows/metrics-tracking.yml`: Download metrics tracking
