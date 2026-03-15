# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2026-03-15

### Added
- Custom error types (`flattener.Error`) for structured error handling with categories: validation, parsing, depth_limit, size_limit, timeout, file_access, security
- Fuzz testing for security validation (`flattener_fuzz_test.go`)
- Comprehensive Makefile for development tasks (build, test, lint, coverage, fuzz, security-audit)
- Consolidated CI workflow (`ci.yml`) with coverage reporting
- Dependabot configuration for automated dependency updates
- Go module caching in all CI workflows

### Changed
- Renamed `NewFlattener()` to `Default()` for better API clarity
- Renamed internal package `internal/utils` to `internal/yamlutil`
- Improved key sanitization: now removes all control characters (not just null bytes) and trims whitespace
- Updated complete example (`examples/complete-example/main.tf`)
- Upgraded Go from 1.23 to 1.25.7

### Fixed
- Go directive version (1.25.8 → 1.25.7) to match latest available release
- golangci-lint v2 staticcheck and revive compatibility issues
- Security-updates workflow PAT_TOKEN authentication
- Complete example for integration tests

### Security
- Enhanced YAML key sanitization: removes all control characters (0x00-0x1F, 0x7F-0x9F)

### CI/CD
- Bumped actions/checkout v4 → v6
- Bumped actions/setup-go v5 → v6
- Bumped actions/upload-artifact v4 → v7
- Bumped goreleaser/goreleaser-action v6 → v7
- Bumped hashicorp/setup-terraform v3 → v4
- Bumped crazy-max/ghaction-import-gpg v6 → v7
- Bumped actions/github-script v7 → v8
- Bumped actions/labeler v5 → v6
- Bumped amannn/action-semantic-pull-request v5 → v6
- Bumped golangci/golangci-lint-action v8 → v9

### Dependencies
- terraform-plugin-framework v1.15.0 → v1.19.0
- terraform-plugin-go v0.25.0 → v0.31.0
- terraform-plugin-testing v1.12.0 → v1.15.0
- terraform-plugin-sdk/v2 v2.35.0 → v2.40.0
- cloudflare/circl v1.6.1 → v1.6.3
- Bump transitive dependencies (x/crypto, x/net, x/sys, x/text, grpc, go-crypto, etc.)

## [0.1.0] - 2025-07-24

### Added
- Initial release of terraform-provider-yamlflattener
- Core YAML flattening functionality
- Data source implementation
- Provider function implementation
- Documentation and examples
- CI/CD pipeline with multi-architecture builds
- Security scanning integration
- Terraform Registry publishing

---

## Release Template

Use this template for future releases:

## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Now removed features

### Fixed
- Bug fixes

### Security
- Security improvements
