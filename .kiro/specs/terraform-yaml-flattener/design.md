# Design Document

## Overview

The Terraform YAML Flattener provider will be implemented as a custom Terraform provider using the Terraform Plugin Framework. The provider will offer two main interfaces: a data source for declarative usage and a provider function for inline usage. Both will implement the same core flattening logic to convert nested YAML structures into flat key-value maps with dot notation for objects and bracket notation for arrays.

## Architecture

The provider will follow the standard Terraform provider architecture with these main components:

```
terraform-provider-yamlflattener/
├── internal/
│   ├── provider/
│   │   ├── provider.go           # Main provider configuration
│   │   ├── data_source_flatten.go # Data source implementation
│   │   └── function_flatten.go    # Provider function implementation
│   ├── flattener/
│   │   ├── flattener.go          # Core flattening logic
│   │   └── flattener_test.go     # Unit tests for flattening
│   └── utils/
│       └── yaml_parser.go        # YAML parsing utilities
├── examples/                     # Usage examples
├── docs/                        # Documentation
├── main.go                      # Provider entry point
└── go.mod                       # Go module definition
```

The provider will be built using:
- **Terraform Plugin Framework v1.x** for modern provider development
- **gopkg.in/yaml.v3** for robust YAML parsing
- **Go 1.21+** for implementation

## Components and Interfaces

### Core Flattener Component

The central flattening logic will be implemented as a reusable component:

```go
type Flattener struct {
    // Configuration options if needed in future
}

func (f *Flattener) FlattenYAML(yamlContent string) (map[string]string, error)
func (f *Flattener) flattenValue(value interface{}, prefix string) map[string]string
```

Key flattening rules:
- Objects: `parent.child.key`
- Arrays: `parent.array[0].key`
- Primitive values: converted to strings
- Null values: represented as empty strings
- Boolean values: "true"/"false"

### Data Source Implementation

```go
type flattenDataSource struct{}

// Schema defines the data source configuration
func (d *flattenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)

// Read performs the flattening operation
func (d *flattenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
```

Data source attributes:
- `yaml_content` (string, optional): Direct YAML content
- `yaml_file` (string, optional): Path to YAML file
- `flattened` (map[string]string, computed): Resulting flattened map

### Provider Function Implementation

```go
func NewFlattenFunction() function.Function {
    return &flattenFunction{}
}

type flattenFunction struct{}

func (f *flattenFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse)
func (f *flattenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse)
```

Function signature:
- Input: `yaml_content` (string) - YAML content to flatten
- Output: `map[string]string` - Flattened key-value pairs

### YAML Parser Utilities

```go
func ParseYAML(content string) (interface{}, error)
func ReadYAMLFile(filepath string) (string, error)
func ValidateYAML(content string) error
```

## Data Models

### Input Models

```go
// For data source
type FlattenDataSourceModel struct {
    YAMLContent types.String `tfsdk:"yaml_content"`
    YAMLFile    types.String `tfsdk:"yaml_file"`
    Flattened   types.Map    `tfsdk:"flattened"`
}

// For function (handled by framework)
// Input: string (YAML content)
// Output: map[string]string
```

### Internal Models

```go
// Represents the parsed YAML structure
type ParsedYAML interface{}

// Flattening result
type FlattenResult struct {
    Data map[string]string
    Error error
}
```

## Error Handling

The provider will implement comprehensive error handling:

### YAML Parsing Errors
- Invalid YAML syntax
- Unsupported YAML features
- File reading errors (permissions, not found)

### Validation Errors
- Missing required inputs
- Conflicting input parameters (both yaml_content and yaml_file provided)

### Runtime Errors
- Stack overflow protection for deeply nested structures
- Memory limits for large YAML files

Error messages will be clear and actionable:
```go
// Example error handling
if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
    return nil, fmt.Errorf("failed to parse YAML content: %w", err)
}

if yamlContent == "" && yamlFile == "" {
    return nil, fmt.Errorf("either yaml_content or yaml_file must be provided")
}
```

## Testing Strategy

### Unit Tests
- Core flattening logic with various YAML structures
- Edge cases: empty YAML, deeply nested structures, arrays with mixed types
- Error conditions: invalid YAML, malformed structures

### Integration Tests
- Data source functionality with Terraform configurations
- Provider function usage in expressions
- File reading capabilities

### Test Cases
1. **Simple Object Flattening**
   ```yaml
   key1: value1
   key2:
     nested: value2
   ```
   Expected: `{"key1": "value1", "key2.nested": "value2"}`

2. **Array Flattening**
   ```yaml
   items:
     - name: item1
     - name: item2
   ```
   Expected: `{"items[0].name": "item1", "items[1].name": "item2"}`

3. **Complex Mixed Structure** (as provided in requirements)
   ```yaml
   alertmanager:
     config:
       global:
         slack_api_url: "your-encrypted-slack-webhook"
       receivers:
         - name: "slack-notifications"
           slack_configs:
             - api_url: "your-encrypted-webhook-url"
   ```

4. **Error Cases**
   - Invalid YAML syntax
   - Non-existent file paths
   - Circular references (if possible in YAML)

### Acceptance Tests
- Full Terraform configuration tests using both data source and function
- Performance tests with large YAML files
- Cross-platform compatibility tests

## Implementation Notes

### Performance Considerations
- Efficient recursive traversal to avoid unnecessary memory allocation
- Streaming approach for large YAML files if needed
- Caching parsed YAML within the same Terraform run

### Security Considerations
- File path validation to prevent directory traversal attacks
- Memory limits to prevent DoS through large YAML files
- Input sanitization for YAML content

### Terraform Integration
- Proper state management for data sources
- Deterministic output to avoid unnecessary plan changes
- Support for Terraform's type system and validation

## Project Structure and Repository Setup

### GitHub Repository Structure

The project will be hosted under the "Perun-Engineering" GitHub organization with the following structure:

```
terraform-provider-yamlflattener/
├── .github/
│   ├── workflows/
│   │   ├── build.yml              # Multi-architecture build workflow
│   │   ├── security.yml           # Security scanning workflow
│   │   ├── release.yml            # Terraform Registry publishing
│   │   └── pr-validation.yml      # Pull request validation
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
├── internal/                      # Provider implementation (as above)
├── examples/                      # Usage examples
├── docs/                         # Documentation
├── scripts/                      # Build and utility scripts
├── .goreleaser.yml               # GoReleaser configuration
├── .golangci.yml                 # Linting configuration
├── README.md                     # Project documentation
├── LICENSE                       # License file
├── CHANGELOG.md                  # Release notes
├── main.go                       # Provider entry point
└── go.mod                        # Go module definition
```

### Branch Protection and Contribution Workflow

- **Main branch protection**: Require pull request reviews, status checks
- **Required status checks**: Build, security scan, tests
- **Merge strategy**: Squash and merge for clean history
- **Semantic versioning**: Following semver for releases

## CI/CD Pipeline Design

### Build Workflow Architecture

The build pipeline will support multiple architectures and platforms:

**Target Architectures:**
- `linux/amd64` - Primary Linux platform
- `linux/arm64` - ARM-based Linux systems
- `darwin/amd64` - Intel-based macOS
- `darwin/arm64` - Apple Silicon macOS
- `windows/amd64` - Windows platform

**Build Process:**
1. **Code Checkout**: Fetch source code and dependencies
2. **Go Setup**: Configure Go environment (version 1.21+)
3. **Dependency Management**: Download and cache Go modules
4. **Cross-compilation**: Build binaries for all target architectures
5. **Artifact Storage**: Store binaries as GitHub Actions artifacts
6. **Testing**: Run unit and integration tests across platforms

### Security Scanning Integration

**Security Scanning Components:**
1. **Dependency Scanning**:
   - Use `govulncheck` for Go vulnerability detection
   - Scan `go.mod` dependencies for known CVEs
   - Fail builds on high/critical vulnerabilities

2. **Static Code Analysis**:
   - `golangci-lint` for code quality and security issues
   - `gosec` for security-specific static analysis
   - Custom rules for Terraform provider best practices

3. **Supply Chain Security**:
   - Verify dependency checksums
   - SBOM (Software Bill of Materials) generation
   - Signed commits verification

### Release and Publishing Pipeline

**Terraform Registry Publishing Process:**
1. **Release Trigger**: Automated on Git tag creation (e.g., `v1.0.0`)
2. **Binary Building**: Cross-compile for all supported architectures
3. **Code Signing**: Sign binaries using GPG keys stored in GitHub secrets
4. **Checksum Generation**: Create SHA256 checksums for all binaries
5. **Registry Publishing**: Automated upload to Terraform Registry
6. **Documentation Update**: Sync provider documentation

**Release Artifacts:**
- Signed provider binaries for each architecture
- SHA256SUMS file with checksums
- SHA256SUMS.sig file with GPG signature
- Release notes from CHANGELOG.md

## GitHub Actions Workflow Specifications

### Build Workflow (`build.yml`)
```yaml
# Triggered on: push, pull_request
# Jobs:
#   - test: Run unit tests
#   - build: Cross-compile for all architectures
#   - validate: Run integration tests
```

### Security Workflow (`security.yml`)
```yaml
# Triggered on: push, pull_request, schedule (weekly)
# Jobs:
#   - vulnerability-scan: Check Go dependencies
#   - static-analysis: Run security-focused linters
#   - dependency-review: Analyze new dependencies in PRs
```

### Release Workflow (`release.yml`)
```yaml
# Triggered on: tag creation (v*)
# Jobs:
#   - build-and-sign: Create signed binaries
#   - publish-registry: Upload to Terraform Registry
#   - create-release: Generate GitHub release with artifacts
```

## Documentation and Examples

### Provider Documentation Structure
- **README.md**: Quick start guide and basic usage
- **docs/**: Detailed documentation for data sources and functions
- **examples/**: Real-world usage examples
- **CHANGELOG.md**: Version history and breaking changes

### Usage Examples
The repository will include comprehensive examples:
1. Basic YAML flattening with data source
2. Inline usage with provider function
3. Complex nested structure handling
4. Integration with other Terraform resources
5. Error handling scenarios

## Deployment and Distribution Strategy

### Terraform Registry Integration
- **Namespace**: `Perun-Engineering/yamlflattener`
- **Versioning**: Semantic versioning (semver)
- **Documentation**: Auto-generated from provider schema
- **Examples**: Included in registry listing

### Binary Distribution
- **GitHub Releases**: Primary distribution method
- **Architecture Support**: All major platforms
- **Verification**: GPG signatures for security
- **Automation**: Fully automated release process

### Monitoring and Maintenance
- **Download Metrics**: Track usage through Terraform Registry
- **Issue Tracking**: GitHub Issues for bug reports and feature requests
- **Security Updates**: Automated dependency updates via Dependabot
- **Community Support**: Clear contribution guidelines and issue templates
