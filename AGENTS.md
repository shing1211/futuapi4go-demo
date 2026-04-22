# futuapi4go-demo AGENTS.md

## Project

Go demo showcasing the futuapi4go SDK. Connects to a running Futu OpenD instance and demonstrates all major APIs via an interactive menu.

## Dev Commands

```bash
go run ./cmd/demo/main.go                    # Run demo (requires Futu OpenD on 127.0.0.1:11111)
FUTU_ADDR=192.168.1.100:11111 go run ./cmd/demo/main.go  # Custom OpenD address
go build ./...                    # Build
go vet ./...                     # Lint
```

## OpenD Simulator (for testing without a real account)

```bash
# Terminal 1: run the simulator (in futuapi4go repo)
cd D:\github\futuapi4go
go run ./cmd/examples/simulator

# Terminal 2: run the demo
cd D:\github\futuapi4go-demo
go run ./cmd/demo/main.go
```

## Project Structure

```
futuapi4go-demo/
├── cmd/demo/main.go     # Source code
├── docs/               # Documentation (README, LICENSE, etc.)
├── scripts/            # Build scripts (.bat for Windows, .sh for Linux/macOS)
├── .github/           # GitHub config (issue templates, PR template)
├── AGENTS.md         # This file
├── go.mod
└── go.sum
```

## Build Scripts

| Script | Platform | Description |
|--------|---------|-------------|
| `scripts/build.bat` / `.sh` | Win/Mac/Linux | Build binary to `cmd/demo/` |
| `scripts/clean.bat` / `.sh` | Win/Mac/Linux | Clean build artifacts |
| `scripts/upgrade.bat` / `.sh` | Win/Mac/Linux | Upgrade dependencies |

## SDK Debugging

The futuapi4go SDK is checked out at `D:\github\futuapi4go`. Proto files live in `api/proto/`. Generated Go code lives in `pkg/pb/`. To regenerate proto files:

```powershell
cd D:\github\futuapi4go
# Use the regen scripts in scripts/ (PowerShell or batch)
```

## SDK Bug: OpenD Rejects Empty C2S Requests

OpenD (serverVer=1003) rejects some API calls with "解析protobuf协议失败".

**Affected APIs:**
- `GetDelayStatistics` — C2S has all optional fields
- `GetTradeDate` — C2S has all required fields

**Workaround:** Demo logs yellow warnings and continues when these fail.

**Proto reference:** See `docs/FUTU_PROTO_REF.md` or https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md

## Related Repositories

- SDK: `github.com/shing1211/futuapi4go` (checked out at `D:\github\futuapi4go`)
- Official Proto Doc: https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md
- OpenD Downloads: https://www.futunn.com/download/fetch-lasted-link?name=opend-windows
