# Domain Context

## Terms

- **Flattener** — The core module (`internal/flattener`) that transforms nested YAML structures into flat `map[string]string` with dot notation for objects and bracket notation for arrays. Accepts YAML as a string (`FlattenYAMLString`) or a file path (`FlattenYAMLFile`). File path handling includes security checks (directory traversal rejection) and size enforcement. Configured via exported fields (`MaxNestingDepth`, `MaxResultSize`, `MaxYAMLSize`). Instantiated with `flattener.New()`.

- **Flatten data source** — The Terraform data source (`yamlflattener_flatten`) that exposes flattening via `yaml_content` or `yaml_file` attributes. Receives a configured Flattener from the provider via `Configure()`.

- **Flatten function** — The Terraform provider function (`provider::yamlflattener::flatten`) that exposes flattening as a pure function call. Receives a configured Flattener via its constructor.
