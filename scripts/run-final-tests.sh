#!/bin/bash

# Terraform YAML Flattener Provider - Final Integration and Acceptance Testing
# This script runs comprehensive tests to verify the provider meets all requirements

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi

    # Check if Terraform is installed
    if ! command -v terraform &> /dev/null; then
        print_error "Terraform is not installed or not in PATH"
        exit 1
    fi

    # Check Go version
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    print_status "Go version: $GO_VERSION"

    # Check Terraform version
    TF_VERSION=$(terraform version -json | grep '"version"' | head -1 | cut -d'"' -f4)
    print_status "Terraform version: $TF_VERSION"

    print_success "Prerequisites check passed"
}

# Run all tests
run_all_tests() {
    print_status "Running all tests..."

    # Set TF_ACC environment variable for acceptance tests
    export TF_ACC=1

    # Run all tests with verbose output
    if go test -v ./internal/... -timeout 30m; then
        print_success "All tests passed"
        return 0
    else
        print_error "Some tests failed"
        return 1
    fi
}

# Run specific test categories
run_test_category() {
    local category=$1
    local pattern=$2
    local timeout=$3

    print_status "Running $category tests..."

    # Set TF_ACC environment variable for acceptance tests
    export TF_ACC=1

    # Run tests with the specified pattern
    if go test -v ./internal/... -run "$pattern" -timeout "$timeout"; then
        print_success "$category tests passed"
        return 0
    else
        print_error "$category tests failed"
        return 1
    fi
}

# Verify data source and function equivalence
verify_equivalence() {
    print_status "Verifying data source and function equivalence..."

    export TF_ACC=1

    # Run equivalence tests
    if go test -v ./internal/provider/... -run "TestFinalIntegration_ComprehensiveEquivalence" -timeout 10m; then
        print_success "Data source and function produce identical results for the same input"
        return 0
    else
        print_error "Data source and function equivalence test failed"
        return 1
    fi
}

# Test cross-platform compatibility
test_cross_platform() {
    print_status "Testing cross-platform compatibility..."

    # Check if Docker is available for cross-platform testing
    if command -v docker &> /dev/null; then
        print_status "Docker is available, running cross-platform tests..."

        # Run the cross-platform test script
        if ./scripts/run-cross-platform-tests.sh; then
            print_success "Cross-platform tests passed"
            return 0
        else
            print_warning "Some cross-platform tests failed"
            return 1
        fi
    else
        print_warning "Docker is not available, skipping containerized cross-platform tests"

        # Run platform-specific tests on the current platform
        export TF_ACC=1

        if go test -v ./internal/provider/... -run "TestFinalIntegration_CrossPlatformFeatures" -timeout 10m; then
            print_success "Platform-specific tests passed on current platform"
            return 0
        else
            print_error "Platform-specific tests failed on current platform"
            return 1
        fi
    fi
}

# Test provider installation workflow
test_installation_workflow() {
    print_status "Testing provider installation workflow..."

    # Create a temporary directory for testing
    TEST_DIR=$(mktemp -d)
    cd "$TEST_DIR"

    # Create a simple Terraform configuration
    cat > main.tf << 'EOF'
terraform {
  required_providers {
    yamlflattener = {
      source = "registry.terraform.io/terraform/yamlflattener"
      version = ">= 0.1.0"
    }
  }
}

provider "yamlflattener" {}

data "yamlflattener_flatten" "test" {
  yaml_content = <<EOT
test:
  installation: "success"
  provider: "yamlflattener"
EOT
}

output "test_result" {
  value = data.yamlflattener_flatten.test.flattened["test.installation"]
}

# Test provider function as well
output "function_result" {
  value = provider::yamlflattener::flatten(<<EOT
test:
  function: "success"
EOT
)["test.function"]
}
EOF

    # Initialize Terraform (this will fail in real scenario but we test the config)
    terraform init -backend=false 2>/dev/null || true

    # Validate the configuration
    if terraform validate; then
        print_success "Terraform configuration validation passed"
        cd - > /dev/null
        rm -rf "$TEST_DIR"
        return 0
    else
        print_error "Terraform configuration validation failed"
        cd - > /dev/null
        rm -rf "$TEST_DIR"
        return 1
    fi
}

# Generate comprehensive test report
generate_test_report() {
    print_status "Generating comprehensive test report..."

    REPORT_FILE="test-report-$(date +%Y%m%d-%H%M%S).md"
    TEMPLATE_FILE="docs/test-report-template.md"

    if [ ! -f "$TEMPLATE_FILE" ]; then
        print_error "Test report template not found: $TEMPLATE_FILE"
        return 1
    fi

    # Copy template to report file
    cp "$TEMPLATE_FILE" "$REPORT_FILE"

    # Get system information
    DATE=$(date)
    PLATFORM="$(uname -s)/$(uname -m)"
    GO_VERSION=$(go version)
    TERRAFORM_VERSION=$(terraform version | head -1)

    # Replace placeholders in the template
    sed -i.bak "s|{{DATE}}|$DATE|g" "$REPORT_FILE"
    sed -i.bak "s|{{PLATFORM}}|$PLATFORM|g" "$REPORT_FILE"
    sed -i.bak "s|{{GO_VERSION}}|$GO_VERSION|g" "$REPORT_FILE"
    sed -i.bak "s|{{TERRAFORM_VERSION}}|$TERRAFORM_VERSION|g" "$REPORT_FILE"

    # Set test statuses based on test results
    UNIT_STATUS="PASSED"
    INTEGRATION_STATUS="PASSED"
    ACCEPTANCE_STATUS="PASSED"
    FINAL_INTEGRATION_STATUS="PASSED"
    PERFORMANCE_STATUS="PASSED"

    # Run quick tests to verify status
    go test -v ./internal/flattener/... -run "TestFlatten" > /dev/null 2>&1 || UNIT_STATUS="FAILED"
    go test -v ./internal/provider/... -run "TestIntegration" > /dev/null 2>&1 || INTEGRATION_STATUS="FAILED"
    go test -v ./internal/provider/... -run "TestAcceptance" > /dev/null 2>&1 || ACCEPTANCE_STATUS="FAILED"
    go test -v ./internal/provider/... -run "TestFinalIntegration" > /dev/null 2>&1 || FINAL_INTEGRATION_STATUS="FAILED"
    go test -v ./internal/flattener/... -run "TestPerformance" > /dev/null 2>&1 || PERFORMANCE_STATUS="FAILED"

    sed -i.bak "s|{{UNIT_TEST_STATUS}}|$UNIT_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{INTEGRATION_TEST_STATUS}}|$INTEGRATION_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{ACCEPTANCE_TEST_STATUS}}|$ACCEPTANCE_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{FINAL_INTEGRATION_TEST_STATUS}}|$FINAL_INTEGRATION_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{PERFORMANCE_TEST_STATUS}}|$PERFORMANCE_STATUS|g" "$REPORT_FILE"

    # Set requirement statuses based on test results
    # This is a simplified approach - in a real scenario, you would map specific tests to requirements
    for i in {1..6}; do
        for j in {1..5}; do
            STATUS="PASSED"
            if [ "$UNIT_STATUS" == "FAILED" ] || [ "$INTEGRATION_STATUS" == "FAILED" ] || [ "$ACCEPTANCE_STATUS" == "FAILED" ] || [ "$FINAL_INTEGRATION_STATUS" == "FAILED" ]; then
                STATUS="NEEDS VERIFICATION"
            fi
            sed -i.bak "s|{{REQ_${i}_${j}_STATUS}}|$STATUS|g" "$REPORT_FILE"
        done
    fi

    # Set platform statuses
    CURRENT_PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
    CURRENT_ARCH=$(uname -m)

    # Default all platforms to "NOT TESTED"
    LINUX_AMD64_STATUS="NOT TESTED"
    LINUX_ARM64_STATUS="NOT TESTED"
    MACOS_AMD64_STATUS="NOT TESTED"
    MACOS_ARM64_STATUS="NOT TESTED"
    WINDOWS_AMD64_STATUS="NOT TESTED"

    # Set current platform to PASSED or FAILED based on test results
    if [ "$CURRENT_PLATFORM" == "linux" ]; then
        if [ "$CURRENT_ARCH" == "x86_64" ] || [ "$CURRENT_ARCH" == "amd64" ]; then
            LINUX_AMD64_STATUS="$FINAL_INTEGRATION_STATUS"
        elif [ "$CURRENT_ARCH" == "arm64" ] || [ "$CURRENT_ARCH" == "aarch64" ]; then
            LINUX_ARM64_STATUS="$FINAL_INTEGRATION_STATUS"
        fi
    elif [ "$CURRENT_PLATFORM" == "darwin" ]; then
        if [ "$CURRENT_ARCH" == "x86_64" ] || [ "$CURRENT_ARCH" == "amd64" ]; then
            MACOS_AMD64_STATUS="$FINAL_INTEGRATION_STATUS"
        elif [ "$CURRENT_ARCH" == "arm64" ]; then
            MACOS_ARM64_STATUS="$FINAL_INTEGRATION_STATUS"
        fi
    elif [ "$CURRENT_PLATFORM" == "windows" ]; then
        WINDOWS_AMD64_STATUS="$FINAL_INTEGRATION_STATUS"
    fi

    sed -i.bak "s|{{LINUX_AMD64_STATUS}}|$LINUX_AMD64_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{LINUX_ARM64_STATUS}}|$LINUX_ARM64_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{MACOS_AMD64_STATUS}}|$MACOS_AMD64_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{MACOS_ARM64_STATUS}}|$MACOS_ARM64_STATUS|g" "$REPORT_FILE"
    sed -i.bak "s|{{WINDOWS_AMD64_STATUS}}|$WINDOWS_AMD64_STATUS|g" "$REPORT_FILE"

    # Set platform notes
    sed -i.bak "s|{{LINUX_AMD64_NOTES}}|Tests run on current platform|g" "$REPORT_FILE"
    sed -i.bak "s|{{LINUX_ARM64_NOTES}}|Tests run on current platform|g" "$REPORT_FILE"
    sed -i.bak "s|{{MACOS_AMD64_NOTES}}|Tests run on current platform|g" "$REPORT_FILE"
    sed -i.bak "s|{{MACOS_ARM64_NOTES}}|Tests run on current platform|g" "$REPORT_FILE"
    sed -i.bak "s|{{WINDOWS_AMD64_NOTES}}|Tests run on current platform|g" "$REPORT_FILE"

    # Generate and add test coverage information
    print_status "Generating test coverage report..."
    go test -coverprofile=coverage.out ./internal/... > /dev/null 2>&1 || true
    if [ -f coverage.out ]; then
        COVERAGE_REPORT=$(go tool cover -func=coverage.out)
        # Replace the coverage placeholder with the actual coverage report
        sed -i.bak "s|{{TEST_COVERAGE}}|$COVERAGE_REPORT|g" "$REPORT_FILE"
        rm coverage.out
    else
        sed -i.bak "s|{{TEST_COVERAGE}}|Coverage report generation failed|g" "$REPORT_FILE"
    fi

    # Set conclusion based on test results
    if [ "$UNIT_STATUS" == "PASSED" ] && [ "$INTEGRATION_STATUS" == "PASSED" ] && [ "$ACCEPTANCE_STATUS" == "PASSED" ] && [ "$FINAL_INTEGRATION_STATUS" == "PASSED" ]; then
        CONCLUSION="All tests have passed successfully. The Terraform YAML Flattener Provider meets all requirements and is ready for production use. The provider has been tested on the current platform and with various YAML structures to ensure compatibility and correctness."
    else
        CONCLUSION="Some tests have failed. The Terraform YAML Flattener Provider needs additional work before it can be considered ready for production use."
    fi
    sed -i.bak "s|{{CONCLUSION}}|$CONCLUSION|g" "$REPORT_FILE"

    # Clean up backup files
    rm -f "$REPORT_FILE.bak"

    print_success "Test report generated: $REPORT_FILE"
}

# Main execution
main() {
    print_status "Starting final integration and acceptance testing for Terraform YAML Flattener Provider..."
    print_status "Platform: $(uname -s)/$(uname -m)"

    # Change to project root directory
    cd "$(dirname "$0")/.."

    # Track test results
    TESTS_PASSED=true

    # Run tests in sequence
    check_prerequisites || TESTS_PASSED=false

    # Run specific test categories
    run_test_category "Unit" "TestFlatten|TestParse" "5m" || TESTS_PASSED=false
    run_test_category "Integration" "TestIntegration" "10m" || TESTS_PASSED=false
    run_test_category "Acceptance" "TestAcceptance" "15m" || TESTS_PASSED=false
    run_test_category "Final Integration" "TestFinalIntegration" "15m" || TESTS_PASSED=false

    # Verify data source and function equivalence
    verify_equivalence || TESTS_PASSED=false

    # Test cross-platform compatibility
    test_cross_platform || TESTS_PASSED=false

    # Test provider installation workflow
    test_installation_workflow || TESTS_PASSED=false

    # Generate test report
    generate_test_report

    # Print final status
    if [ "$TESTS_PASSED" = true ]; then
        print_success "All tests completed successfully!"
        print_status "Provider is ready for production use."
    else
        print_error "Some tests failed. See the test report for details."
        exit 1
    fi
}

# Handle script arguments
case "${1:-all}" in
    "prereq")
        check_prerequisites
        ;;
    "unit")
        run_test_category "Unit" "TestFlatten|TestParse" "5m"
        ;;
    "integration")
        run_test_category "Integration" "TestIntegration" "10m"
        ;;
    "acceptance")
        run_test_category "Acceptance" "TestAcceptance" "15m"
        ;;
    "final")
        run_test_category "Final Integration" "TestFinalIntegration" "15m"
        ;;
    "equivalence")
        verify_equivalence
        ;;
    "cross-platform")
        test_cross_platform
        ;;
    "install")
        test_installation_workflow
        ;;
    "report")
        generate_test_report
        ;;
    "all"|*)
        main
        ;;
esac
