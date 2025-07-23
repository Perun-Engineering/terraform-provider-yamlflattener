# Final Integration and Acceptance Testing

This document provides instructions for running the final integration and acceptance tests for the Terraform YAML Flattener Provider.

## Overview

The final integration and acceptance testing verifies that the provider meets all requirements and works correctly across different platforms. The tests focus on:

1. Running full provider tests with real Terraform configurations
2. Testing cross-platform compatibility
3. Verifying that both data source and function produce identical results for the same input
4. Testing the provider installation and usage workflow

## Prerequisites

- Go 1.21 or later
- Terraform 1.0 or later
- Docker (optional, for cross-platform testing)

## Running the Tests

### Complete Test Suite

To run the complete test suite:

```bash
./scripts/run-final-tests.sh
```

This will:
1. Check prerequisites
2. Run unit tests
3. Run integration tests
4. Run acceptance tests
5. Run final integration tests
6. Verify data source and function equivalence
7. Test cross-platform compatibility
8. Test provider installation workflow
9. Generate a comprehensive test report

### Specific Test Categories

To run specific test categories:

```bash
./scripts/run-final-tests.sh unit       # Run unit tests
./scripts/run-final-tests.sh integration # Run integration tests
./scripts/run-final-tests.sh acceptance  # Run acceptance tests
./scripts/run-final-tests.sh final       # Run final integration tests
./scripts/run-final-tests.sh equivalence # Test data source and function equivalence
./scripts/run-final-tests.sh cross-platform # Test cross-platform compatibility
./scripts/run-final-tests.sh install     # Test provider installation workflow
./scripts/run-final-tests.sh report      # Generate test report
```

## Cross-Platform Testing

The provider is tested on multiple platforms to ensure compatibility:

- Linux/amd64
- Linux/arm64
- macOS/amd64
- macOS/arm64
- Windows/amd64

If Docker is available, the cross-platform tests will use Docker containers to test on different platforms. Otherwise, the tests will run only on the current platform.

To run cross-platform tests specifically:

```bash
./scripts/run-cross-platform-tests.sh
```

## Test Report

A comprehensive test report is generated at the end of the testing process. The report includes:

- Test environment details
- Test results for each category
- Requirements coverage
- Cross-platform compatibility
- Test coverage metrics
- Conclusion and recommendations

The report is saved as `test-report-{timestamp}.md` in the project root directory.

## Requirements Verification

The tests verify that the provider meets all requirements specified in the requirements document. In particular, the tests focus on:

- Requirement 3.3: Return correct string values through both methods
- Requirement 3.4: Show expected flattened output in Terraform plan
- Requirement 6.5: Produce identical output for same input in both implementations

## Troubleshooting

If tests fail, check the following:

1. Ensure Go and Terraform are installed and in your PATH
2. Ensure you have the necessary permissions to run the tests
3. Check the test output for specific error messages
4. Verify that the provider builds successfully
5. Check for platform-specific issues if running on a non-Linux platform

For Docker-based cross-platform tests, ensure Docker is installed and running.
