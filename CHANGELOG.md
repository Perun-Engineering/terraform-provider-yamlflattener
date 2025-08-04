# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [0.2.0] - 2025-01-04

### Added
- **True Order Preservation**: YAML key order is now perfectly preserved during flattening for consistent and predictable results
- **OrderedMap Implementation**: Internal ordered data structure ensures iteration order matches YAML order
- **Simplified API**: Single implementation that always preserves order

### Changed
- **BREAKING**: All flattening methods now return `*OrderedMap` instead of `map[string]string`
- **Unified Implementation**: Removed redundant methods, keeping only ordered versions
- **Provider Integration**: Data source and function now use ordered iteration to maintain key sequence

### Fixed
- **Order Consistency**: Fixed issue where Go's random map iteration was losing YAML key order
- **Terraform Display**: Keys now appear in correct order in Terraform output

### Performance
- **Single-Pass Processing**: Optimized flattening with unified ordered implementation
- **Memory Efficient**: OrderedMap provides better memory usage patterns

## [0.1.0] - 2024-12-XX

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
