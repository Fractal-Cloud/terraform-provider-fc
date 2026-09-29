# Contributing to the Fractal Cloud Terraform Provider

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.1
- [golangci-lint](https://golangci-lint.run/welcome/install/) for linting

### Setting Up Your Development Environment

1. Fork the repository on GitHub.

2. Clone your fork:
   ```sh
   git clone https://github.com/YOUR-USERNAME/terraform-provider.git
   cd terraform-provider
   ```

3. Add the upstream remote:
   ```sh
   git remote add upstream https://github.com/Fractal-Cloud/terraform-provider-fc.git
   ```

4. Install dependencies:
   ```sh
   go mod download
   ```

5. Verify your setup:
   ```sh
   make build
   make test
   ```

### Local Provider Testing

To test the provider locally against a Fractal Cloud account, add a dev override to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/fractal-cloud/fc" = "/path/to/your/go/bin"
  }
  direct {}
}
```

Then build and install:

```sh
make install
```

## How to Contribute

### Reporting Bugs

Before creating a bug report, please search [existing issues](https://github.com/Fractal-Cloud/terraform-provider-fc/issues) to avoid duplicates. When filing a bug report, use the **Bug Report** issue template and include:

- Terraform version (`terraform version`)
- Provider version
- Relevant Terraform configuration (redact any secrets)
- Expected behavior
- Actual behavior
- Steps to reproduce

### Suggesting Features

Feature requests are welcome. Please use the **Feature Request** issue template and describe:

- The problem your feature would solve
- Your proposed solution
- Any alternatives you've considered

### Submitting Changes

1. Create a feature branch from `main`:
   ```sh
   git checkout -b feature/my-feature
   ```

2. Make your changes, following the coding guidelines below.

3. Run formatting and linting:
   ```sh
   make fmt
   make lint
   ```

4. Run tests:
   ```sh
   make test
   ```

5. Commit your changes with a clear, descriptive commit message.

6. Push to your fork and open a pull request against `main`.

### Pull Request Guidelines

- Open an issue first for non-trivial changes to discuss the approach.
- Keep PRs focused -- one feature or fix per PR.
- Include tests for new functionality.
- Update documentation if behavior changes.
- Fill out the PR template completely.
- Ensure all CI checks pass.

## Coding Guidelines

### Project Structure

```
internal/
  client/     # API client layer (HTTP, models, auth)
  provider/   # Terraform provider implementation
    provider.go                  # Provider registration
    *_resource.go                # Resource CRUD implementations
    *_data_source.go             # Data source implementations
    function_*.go                # Provider function implementations
    component_functions.go       # Shared component builder helpers
    diagnostic_helpers.go        # Error handling helpers
```

### Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Run `make lint` before submitting.
- Use the existing code patterns for new resources, data sources, and functions.

### Adding a New Provider Function

Provider functions follow a consistent pattern. Each function file contains:

1. An interface assertion (`var _ function.Function = &MyFunction{}`)
2. A struct and constructor (`NewMyFunction`)
3. `Metadata` -- sets the function name
4. `Definition` -- defines parameters and return type
5. `Run` -- implements the function logic using `buildComponent()`

Register the new function in `provider.go` under the `Functions()` method.

### Adding a New Resource or Data Source

1. Create the resource/data source file following existing patterns.
2. Register it in `provider.go` under `Resources()` or `DataSources()`.
3. Add acceptance tests.
4. Add example HCL in the `examples/` directory.

### Error Handling

- Use `resp.Diagnostics.AddError()` for errors in resources/data sources.
- Always `return` after adding an error.
- Use `resp.State.RemoveResource(ctx)` on 404 in resource Read operations.
- In provider functions, use `function.ConcatFuncErrors()`.

### Testing

- Unit tests: `make test`
- Acceptance tests: `make testacc` (requires `FRACTAL_CLOUD_SERVICE_ACCOUNT_ID` and `FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET` environment variables)

## Release Process

Releases are managed by the maintainers using [GoReleaser](https://goreleaser.com/) and GitHub Actions. To trigger a release:

1. Update `CHANGELOG.md` with the new version's changes.
2. Create and push a tag: `git tag v2.0.1 && git push origin v2.0.1`.
3. The release workflow builds and publishes binaries automatically.

## Questions?

If you have questions about contributing, feel free to open an issue with the question label.
