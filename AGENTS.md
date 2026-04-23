# futuapi4go-demo AGENTS.md

## Project

Go demo showcasing the futuapi4go SDK. Connects to a running Futu OpenD instance and demonstrates all major APIs via an interactive menu.

## Dev Commands

```bash
.\run.bat                    # Run demo (requires Futu OpenD on 127.0.0.1:11111)
FUTU_ADDR=192.168.1.100:11111 .\run.bat  # Custom OpenD address
go build ./...                    # Build
go vet ./...                     # Lint
```

## OpenD Simulator (for testing without a real account)

```bash
# Terminal 1: run the simulator (in futuapi4go repo)
go run D:\github\futuapi4go\cmd\examples\simulator

# Terminal 2: run the demo
.\run.bat
```

## Project Structure

```
futuapi4go-demo/
├── examples/
│   ├── README.md             # Example descriptions
│   ├── 00_connect/          # Connect & disconnect
│   ├── 01_quote/            # GetQuote snapshot
│   ├── 02_ticker/          # SubscribeTicker
│   ├── 03_orderbook/        # SubscribeOrderBook
│   ├── 04_rt/               # SubscribeRT
│   ├── 05_broker/           # SubscribeBroker
│   ├── 06_kline_single/     # GetKLines (one-shot)
│   └── 07_kline_multi/      # SubscribeKLines (multi-period)
├── docs/                     # Supplementary docs (proto reference)
├── build.bat / .sh           # Build script
├── run.bat / .sh             # Run script (default: 00_connect)
├── clean.bat / .sh           # Clean script
└── upgrade.bat / .sh         # Upgrade dependencies
```

## Examples

Run any example:
```bash
go run ./examples/00_connect
go run ./examples/01_quote
go run ./examples/02_ticker
go run ./examples/03_orderbook
go run ./examples/04_rt
go run ./examples/05_broker
go run ./examples/06_kline_single
go run ./examples/07_kline_multi
```

## Scripts

| Script | Platform | Description |
|--------|---------|-------------|
| `build.bat` / `.sh` | Win/Mac/Linux | Build all packages |
| `run.bat` / `.sh` | Win/Mac/Linux | Run an example (default: 00_connect) |
| `clean.bat` / `.sh` | Win/Mac/Linux | Clean build artifacts |
| `upgrade.bat` / `.sh` | Win/Mac/Linux | Upgrade dependencies |

## SDK Debugging

## SDK Debugging

The futuapi4go SDK is checked out at `D:\github\futuapi4go`.

- Proto files: `D:\github\futuapi4go\api\proto\`
- Generated Go protobuf code: `D:\github\futuapi4go\pkg\pb\`
- SDK source: `D:\github\futuapi4go\pkg\`

**To use a local SDK version** (e.g., after fixing proto bugs), add a `replace` directive to `go.mod`:

```go
replace github.com/shing1211/futuapi4go => D:/github/futuapi4go
```

After editing `go.mod`, clear the module cache and re-download:

```powershell
go clean -modcache
go mod download
```

**To regenerate proto files:**

```powershell
cd D:\github\futuapi4go
# Use the regen scripts in scripts/ (PowerShell or batch)
```

## Known SDK Issues

### GetDelayStatistics — proto2 wire-format incompatibility (serverVer=1003)

OpenD rejects the `GetDelayStatistics` request with "解析protobuf协议失败". Root cause: `google.golang.org/protobuf` encodes `repeated int32` fields using proto3 packed wire format by default, but OpenD's C++ parser expects proto2 non-packed encoding. This is an SDK-level issue requiring a fix in `futuapi4go` itself.

**Workaround in demo:** The call is skipped with a printed note. All other APIs work normally.

### GetTradeDate — all C2S fields are required (serverVer=1003)

`GetTradeDate` has all required fields in its C2S. If the SDK doesn't populate all required fields, OpenD returns "解析protobuf协议失败". This may also be a proto2 wire-format issue.

**Workaround in demo:** If this API fails, the demo exits with a red error.

**Proto reference:** See `docs/FUTU_PROTO_REF.md` or https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md

## Related Repositories

- SDK: `github.com/shing1211/futuapi4go` (checked out at `D:\github\futuapi4go`)
- Official Proto Doc: https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md
- OpenD Downloads: https://www.futunn.com/download/fetch-lasted-link?name=opend-windows
