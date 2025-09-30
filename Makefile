.PHONY: help build install test test-unit test-integration test-acceptance test-all \
        lint fmt vet clean coverage coverage-html run-local release \
        install-tools install-local tidy fuzz security-audit

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
BINARY_NAME=terraform-provider-yamlflattener
COVERAGE_FILE=coverage.out

# Build settings
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

help: ## Display this help message
	@echo "Terraform YAML Flattener Provider - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the provider binary
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 $(GOBUILD) -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

install: build ## Build and install the provider locally
	@echo "Installing provider locally..."
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/Perun-Engineering/yamlflattener/$(VERSION)/$(shell go env GOOS)_$(shell go env GOARCH)
	cp $(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/Perun-Engineering/yamlflattener/$(VERSION)/$(shell go env GOOS)_$(shell go env GOARCH)/
	@echo "Provider installed to ~/.terraform.d/plugins/"

install-local: install ## Alias for install

test: test-unit ## Run unit tests (default)

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GOTEST) -v -race -timeout 5m ./...

test-integration: ## Run integration tests (requires TF_ACC=1)
	@echo "Running integration tests..."
	TF_ACC=1 $(GOTEST) -v -timeout 10m ./internal/provider/...

test-acceptance: ## Run acceptance tests (requires TF_ACC=1)
	@echo "Running acceptance tests..."
	TF_ACC=1 $(GOTEST) -v -timeout 20m -run "TestAcc" ./...

test-all: ## Run all tests including acceptance tests
	@echo "Running all tests..."
	TF_ACC=1 $(GOTEST) -v -race -timeout 20m ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report: coverage.html"

fuzz: ## Run fuzz tests
	@echo "Running fuzz tests..."
	$(GOTEST) -fuzz=Fuzz -fuzztime=30s ./internal/flattener/

lint: ## Run golangci-lint
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 5m; \
	else \
		echo "golangci-lint not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

security-audit: ## Run security vulnerability scan
	@echo "Running security audit..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not installed. Installing..."; \
		$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./...; \
	fi

tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	$(GOMOD) tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html
	rm -rf dist/
	@echo "Clean complete"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GOCMD) install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	$(GOCMD) install gotest.tools/gotestsum@latest
	@echo "Tools installed"

release: ## Run goreleaser in snapshot mode (local test)
	@echo "Running goreleaser in snapshot mode..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not installed. Install from https://goreleaser.com/install/"; \
		exit 1; \
	fi

docs: ## Generate provider documentation
	@echo "Generating provider documentation..."
	@if command -v tfplugindocs >/dev/null 2>&1; then \
		tfplugindocs generate; \
	else \
		echo "tfplugindocs not installed. Run 'make install-tools' first."; \
		exit 1; \
	fi

run-local: install ## Build, install, and show installation instructions
	@echo ""
	@echo "Provider installed locally!"
	@echo ""
	@echo "To use the local provider in your Terraform configuration:"
	@echo ""
	@echo "terraform {"
	@echo "  required_providers {"
	@echo "    yamlflattener = {"
	@echo "      source  = \"registry.terraform.io/Perun-Engineering/yamlflattener\""
	@echo "      version = \"$(VERSION)\""
	@echo "    }"
	@echo "  }"
	@echo "}"
	@echo ""

pre-commit: fmt vet lint test-unit ## Run pre-commit checks (format, vet, lint, test)
	@echo "Pre-commit checks complete"

ci: fmt vet lint coverage security-audit ## Run CI checks (format, vet, lint, coverage, security)
	@echo "CI checks complete"
