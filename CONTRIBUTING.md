# Contributing to futuapi4go-demo

> Thanks for helping make the demo better! Every contribution counts.

## Code of Conduct

Be kind and professional. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## Development Setup

```bash
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo
go mod download
go build ./...
```

## Running the Demo

```bash
# With real OpenD
go run ./cmd/demo/main.go

# With mock simulator (no account needed)
# Terminal 1
go run github.com/shing1211/futuapi4go/cmd/examples/simulator
# Terminal 2
go run ./cmd/demo/main.go

# Custom OpenD address
set FUTU_ADDR=192.168.1.100:11111
go run ./cmd/demo/main.go
```

## Code Standards

- Run `go build ./...` and `go vet ./...` before submitting
- Keep `cmd/demo/main.go` as a single file — that's its character
- Helper functions (`sec()`, `ptrStr()`, `ptrInt32()`) live at the top of `main.go`
- Follow existing demo function patterns when adding new categories

## Adding a New Demo Category

1. Add a `func demoXXX(cli *client.Client)` function in `cmd/demo/main.go`
2. Print the section header with `section(n, "Title")`
3. Wire it into the menu switch in `main()`
4. Update the demo menu table in `README.md`

## Testing with a Local SDK

To test changes to the futuapi4go SDK alongside the demo, add a `replace` directive to `go.mod`:

```go
replace github.com/shing1211/futuapi4go => D:/github/futuapi4go
```

Then refresh the module cache:

```bash
go clean -modcache
go mod download
```

## Pull Request Process

1. Fork the repo and create a branch
2. Make your changes — keep commits atomic
3. Ensure `go build ./...` and `go vet ./...` pass
4. Update `README.md` if you add new API coverage
5. Open a PR with a clear description

## License

By contributing, you agree your work will be licensed under Apache License 2.0.
