# Contributing to Terraform YAML Flattener Provider

Thank you for your interest in contributing to the Terraform YAML Flattener Provider! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, please include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, Terraform version, provider version)
- Any relevant logs or error messages

### Suggesting Features

Feature requests are welcome! Please provide:

- A clear and descriptive title
- Detailed description of the proposed feature
- Use cases and benefits
- Any alternative solutions considered

### Development Setup

1. **Prerequisites**
   - Go 1.21 or later
   - Terraform 1.0 or later
   - Git

2. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/terraform-provider-yamlflattener.git
   cd terraform-provider-yamlflattener
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Build the Provider**
   ```bash
   go build -o terraform-provider-yamlflattener
   ```

### Development Workflow

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write code following Go best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**
   ```bash
   # Run unit tests
   go test ./...

   # Run acceptance tests
   TF_ACC=1 go test ./internal/provider -v

   # Run linting
   golangci-lint run
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

### Coding Standards

- Follow standard Go formatting (`go fmt`)
- Write clear, self-documenting code
- Add comments for complex logic
- Use meaningful variable and function names
- Follow the existing code structure and patterns

### Testing Guidelines

- Write unit tests for all new functionality
- Ensure tests cover edge cases and error conditions
- Use table-driven tests where appropriate
- Mock external dependencies
- Maintain high test coverage

### Documentation

- Update README.md for user-facing changes
- Add examples for new features
- Update inline code documentation
- Keep CHANGELOG.md updated

### Issue Management

We use GitHub issues to track bugs, features, and other work. To help organize our work:

1. **Issue Labels**
   - Issues are categorized with labels to indicate type, priority, and status
   - Common labels include: `bug`, `enhancement`, `documentation`, `good first issue`
   - See the full list of labels in `.github/labels.yml`

2. **Project Boards**
   - We use GitHub Projects to track work across milestones
   - Issues move through columns: To Do → In Progress → Review → Done
   - See `.github/project-template.md` for our project board structure

3. **Issue Templates**
   - Use the provided templates when creating new issues
   - Fill out all relevant sections to help maintainers understand your request

4. **Issue Assignment**
   - Feel free to ask for an issue to be assigned to you if you want to work on it
   - Issues labeled `good first issue` are great for new contributors

### Pull Request Process

1. Ensure all tests pass and code is properly formatted
2. Update documentation and examples as needed
3. Fill out the pull request template completely
4. Link any related issues
5. Request review from maintainers
6. Address any feedback from code reviews

### Review Process

- All submissions require review before merging
- Reviewers will check for code quality, tests, and documentation
- Address feedback promptly and professionally
- Maintainers may request changes or provide suggestions

### Release Process

Releases are handled by maintainers:

1. Version bumping follows semantic versioning
2. CHANGELOG.md is updated with release notes
3. GitHub releases are created with binaries
4. Terraform Registry is updated automatically

## Getting Help

- Check existing documentation and examples
- Search through existing issues
- Create a new issue for questions or problems
- Join discussions in pull requests

## Recognition

Contributors will be recognized in:
- CHANGELOG.md for significant contributions
- GitHub contributors list
- Release notes for major features

Thank you for contributing to make this project better!
