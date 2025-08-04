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
- **BREAKING**: Provider function now returns `list(tuple([string, string]))` instead of `map(string)` to preserve order
- **BREAKING**: All flattening methods now return `*OrderedMap` instead of `map[string]string`
- **Unified Implementation**: Removed redundant methods, keeping only ordered versions
- **Provider Integration**: Data source and function now use ordered iteration to maintain key sequence

### Fixed
- **Critical Order Bug**: Fixed issue where `types.MapValue` was losing YAML key order due to Go's random map iteration
- **Terraform Display**: Keys now appear in correct order in Terraform output
- **Provider Function**: Now returns ordered list of key-value pairs that preserves exact YAML structure

### Performance
- **Single-Pass Processing**: Optimized flattening with unified ordered implementation
- **Memory Efficient**: OrderedMap provides better memory usage patterns

### Migration Guide

**For Terraform users:**
```hcl
# Before (v0.1.x)
locals {
  flattened = provider::yamlflattener::flatten(var.yaml)  # map(string)
}

# After (v0.2.0)
locals {
  flattened_pairs = provider::yamlflattener::flatten(var.yaml)  # list(tuple([string, string]))
  flattened_map = { for pair in local.flattened_pairs : pair[0] => pair[1] }  # convert when needed
}

# Use with Helm (preserves order!)
dynamic "set_sensitive" {
  for_each = local.flattened_pairs
  content {
    name  = set_sensitive.value[0]  # key
    value = set_sensitive.value[1]  # value
  }
}
```

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
