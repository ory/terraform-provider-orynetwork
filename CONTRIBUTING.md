# Contributing to the Ory Terraform Provider

Thank you for your interest in contributing to the Ory Terraform Provider!

## Development Setup

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://www.terraform.io/downloads) >= 1.0
- An [Ory Network](https://console.ory.sh/) account for testing

### Building

```bash
# Clone the repository
git clone https://github.com/ory/terraform-provider-orynetwork.git
cd terraform-provider-orynetwork

# Install development tools and set up git hooks
make tools
make hooks

# Build
make build

# Install locally
make install
```

### Running Tests

```bash
# Unit tests (no credentials needed)
make test

# Acceptance tests (requires Ory credentials)
export ORY_WORKSPACE_API_KEY="ory_wak_..."
export ORY_PROJECT_API_KEY="ory_pat_..."
export ORY_PROJECT_ID="..."
export ORY_PROJECT_SLUG="..."
make testacc
```

### Using Local Provider

To use a locally built provider, create a `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "ory/terraform-provider-orynetwork" = "/path/to/terraform-provider-orynetwork"
  }
  direct {}
}
```

## Making Changes

### Adding a New Resource

1. Create a new package in `internal/resources/`
2. Implement the resource with these methods:
   - `Metadata()` - Resource type name
   - `Schema()` - Resource schema definition
   - `Configure()` - Provider configuration
   - `Create()` - Create the resource
   - `Read()` - Read the resource state
   - `Update()` - Update the resource
   - `Delete()` - Delete the resource
   - `ImportState()` - Import existing resources
3. Register the resource in `internal/provider/provider.go`
4. Add documentation in `docs/resources/`
5. Add examples in `examples/resources/`
6. Write acceptance tests

### Code Style

- Run `make fmt` before committing
- Run `make lint` to check for issues
- Follow existing patterns in the codebase
- Add meaningful comments for complex logic

### Commit Messages

Use clear, descriptive commit messages:

```
Add ory_foo resource for managing foos

- Implement CRUD operations
- Add acceptance tests
- Add documentation
```

### Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Submit a pull request

Please include:

- Description of the changes
- Link to any related issues
- Test results or screenshots if applicable

## Reporting Issues

When reporting issues, please include:

- Terraform version (`terraform version`)
- Provider version
- Relevant Terraform configuration (sanitized of secrets)
- Expected behavior
- Actual behavior
- Steps to reproduce

## Code of Conduct

Please be respectful and constructive in all interactions. We're all here to build something useful together.
