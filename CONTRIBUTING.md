# Contributing to futuapi4go-demo

First off, thank you for considering contributing to futuapi4go-demo!

## Code of Conduct

This project is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

- Open a [bug report issue](https://github.com/shing1211/futuapi4go-demo/issues/new?template=bug_report.md).
- Include your Go version, OS, and steps to reproduce.
- Include relevant log output (redact any personal/account info).

### Suggesting Enhancements

- Open a [feature request issue](https://github.com/shing1211/futuapi4go-demo/issues/new?template=feature_request.md).
- Describe the use case and expected behavior.

### Pull Requests

1. **Fork the repo** and create a branch from `main`.
2. **Make your changes** — keep commits atomic and well-described.
3. **Test with a simulator or local OpenD** before submitting.
4. **Open a pull request** with a clear description linking to any related issue.

## Development Setup

```bash
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo
go mod download
go build ./...
```

### Running Tests

```bash
go build ./...
go vet ./...
```

## Coding Standards

- Run `go build ./...` and `go vet ./...` before submitting.
- Keep the demo focused — it should remain a single-file, menu-driven showcase.
- New API additions should mirror the style of existing demo functions.
- Follow standard `gofmt` formatting.

## Pull Request Process

1. Ensure no build errors or vet warnings.
2. Update README.md if you add new API coverage.
3. Your PR will be reviewed by a maintainer. Feedback may be provided.
4. Once approved, a maintainer will merge your PR.
