# Security Policy

## Supported Versions

We currently support the following versions of the Terraform YAML Flattener Provider with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0.0 | :x:                |

## Reporting a Vulnerability

We take the security of our provider seriously. If you believe you've found a security vulnerability, please follow these steps:

1. **Do not disclose the vulnerability publicly**
2. **Email the maintainers** at security@example.com with details about the vulnerability
3. **Include the following information**:
   - Type of vulnerability
   - Full path of source file(s) related to the vulnerability
   - Steps to reproduce
   - Potential impact of the vulnerability
   - Suggested fix if available

## Response Process

When a vulnerability is reported, we will:

1. Confirm receipt of the vulnerability report within 48 hours
2. Assess the vulnerability and determine its impact
3. Develop and test a fix
4. Release a patch as soon as possible, depending on complexity
5. Publicly disclose the vulnerability after the fix has been released

## Security Updates

Security updates will be released as patch versions and announced in:
- GitHub releases
- CHANGELOG.md
- Security advisories on GitHub

## Best Practices

When using this provider, we recommend the following security best practices:

1. Keep the provider updated to the latest version
2. Use Terraform's state encryption features
3. Limit file access permissions for YAML files containing sensitive data
4. Use environment-specific variables for sensitive values
5. Regularly audit your Terraform configurations

## Security Measures

This provider implements the following security measures:

1. Input validation for all YAML content
2. File path validation to prevent directory traversal
3. Memory limits to prevent DoS through large YAML files
4. Dependency scanning in CI/CD pipeline
5. Static code analysis for security issues
