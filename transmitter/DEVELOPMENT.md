# Development Setup

## Prerequisites

- Go 1.21+
- [Lefthook](https://github.com/evilmartians/lefthook) for Git hooks
- [golangci-lint](https://golangci-lint.run/usage/install/) for linting
- [gosec](https://github.com/securecodewarrior/gosec) for security scanning

## Installation

1. Install lefthook:
   ```bash
   go install github.com/evilmartians/lefthook@latest
   ```

2. Install golangci-lint:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. Install gosec:
   ```bash
   go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
   ```

4. Install Git hooks:
   ```bash
   lefthook install
   ```

## Git Hooks

This project uses lefthook to run quality checks:

### Pre-commit hooks:
- `gofmt` - Format Go code
- `goimports` - Organize imports
- `go vet` - Static analysis
- `golangci-lint` - Comprehensive linting
- `go mod tidy` - Clean up module dependencies

### Pre-push hooks:
- `go test -v -race` - Run tests with race detection
- `go build -v` - Ensure code builds
- `gosec` - Security vulnerability scanning
- `go mod verify` - Verify module dependencies

These hooks mirror the same checks that run in CI, catching issues early in development.