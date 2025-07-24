#!/bin/bash
set -e

VERSION=$(echo ${GITHUB_REF#refs/tags/v} || echo "0.1.0")

cat > terraform-provider-yamlflattener_${VERSION}_manifest.json << EOF
{
  "version": 1,
  "metadata": {
    "protocol_versions": ["5.0"]
  }
}
EOF
