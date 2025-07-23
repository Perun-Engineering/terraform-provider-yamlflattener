#!/bin/bash

# Terraform YAML Flattener Provider - Cross-Platform Test Runner
# This script runs tests on multiple platforms using Docker

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

# Check if Docker is installed
check_docker() {
    print_status "Checking if Docker is installed..."

    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed or not in PATH"
        exit 1
    fi

    print_success "Docker is installed"
}

# Run tests on a specific platform
run_platform_tests() {
    local platform=$1
    local tag=$2

    print_status "Running tests on $platform ($tag)..."

    # Create a Docker container for the platform
    docker run --rm -v "$(pwd):/app" -w /app $tag bash -c "
        echo 'Building and testing on $platform...'
        go version
        go build -o terraform-provider-yamlflattener .
        TF_ACC=1 go test -v ./internal/provider/... -run 'TestFinalIntegration' -timeout 10m
    "

    if [ $? -eq 0 ]; then
        print_success "Tests passed on $platform"
        return 0
    else
        print_error "Tests failed on $platform"
        return 1
    fi
}

# Main function
main() {
    print_status "Starting cross-platform tests for Terraform YAML Flattener Provider..."

    # Check if Docker is installed
    check_docker

    # Define platforms to test on
    platforms=(
        "Linux/amd64:golang:1.21-bullseye"
        "Linux/arm64:arm64v8/golang:1.21-bullseye"
        "Windows:golang:1.21-windowsservercore"
    )

    # Track results
    results=()

    # Run tests on each platform
    for platform in "${platforms[@]}"; do
        IFS=':' read -r name tag <<< "$platform"

        if run_platform_tests "$name" "$tag"; then
            results+=("$name: PASSED")
        else
            results+=("$name: FAILED")
        fi
    done

    # Print results
    print_status "Cross-platform test results:"
    for result in "${results[@]}"; do
        if [[ $result == *"PASSED"* ]]; then
            print_success "$result"
        else
            print_error "$result"
        fi
    done
}

# Run the main function
main
