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

### Running the Demo

```bash
# With local OpenD running on default 127.0.0.1:11111
go run ./cmd/demo/main.go

# With custom OpenD address
FUTU_ADDR=192.168.1.100:11111 go run ./cmd/demo/main.go

# With mock simulator (no real account needed)
# Terminal 1: start simulator (requires futuapi4go source at ../futuapi4go)
go run github.com/shing1211/futuapi4go/cmd/examples/simulator
# Terminal 2: run demo
go run ./cmd/demo/main.go
```

### Using a Local SDK Version

If you need to test changes to the [futuapi4go](https://github.com/shing1211/futuapi4go) SDK alongside this demo, add a `replace` directive to `go.mod`:

```go
replace github.com/shing1211/futuapi4go => D:/github/futuapi4go
```

After editing `go.mod`, clear the module cache:

```bash
go clean -modcache
go mod download
```

## Coding Standards

- Run `go build ./...` and `go vet ./...` before submitting.
- Keep the demo focused — it should remain a single-file, menu-driven showcase.
- New API additions should mirror the style of existing demo functions.
- Follow standard `gofmt` formatting.
- Helper functions (`sec()`, `ptrStr()`, `ptrInt32()`, etc.) are defined at the top of `main.go` for reuse.

## Adding a New Demo Category

1. Add a new `func demoXXX(cli *client.Client)` function in `cmd/demo/main.go`.
2. Use `section(n, "Title")` to print the category header.
3. Call the function from `main()` using the interactive menu switch.
4. Update the demo menu table in `README.md`.

## Pull Request Process

1. Ensure no build errors or vet warnings.
2. Update README.md if you add new API coverage or change the demo menu.
3. Your PR will be reviewed by a maintainer. Feedback may be provided.
4. Once approved, a maintainer will merge your PR.
