# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2025-01-XX

### Added
- Custom error types for better error categorization and handling
- Context propagation for cancellation and timeout control
- Fuzz testing for security validation
- Comprehensive Makefile for development tasks
- Consolidated CI workflow with coverage reporting
- Dependabot configuration for automated dependency updates
- Go module caching in all CI workflows

### Changed
- **BREAKING (internal)**: Renamed `NewFlattener()` to `Default()` for better clarity
- **BREAKING (internal)**: Renamed error constructors (removed "New" prefix): `ValidationError()`, `ParsingError()`, etc.
- **BREAKING (internal)**: Renamed `FlattenerError` to `flattener.Error` to avoid stuttering
- Improved key sanitization with character validation (removes control characters)
- Enhanced directory traversal check (now validates before filepath.Clean())
- Better symlink validation in file path security checks
- Optimized CI workflows (consolidated build.yml and pr-validation.yml into ci.yml)

### Fixed
- Directory traversal security vulnerability in file path validation
- PAT_TOKEN fallback in security-updates workflow

### Security
- Enhanced input sanitization for YAML keys (removes all control characters)
- Improved file path validation against directory traversal attacks
- Added symlink resolution validation

## [1.0.0] - TBD

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
