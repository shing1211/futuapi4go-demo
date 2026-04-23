# futuapi4go-demo AGENTS.md

## Project

Go demo showcasing the futuapi4go SDK. Each example is a standalone `main.go` demonstrating one SDK function.

## Dev Commands

```bash
go run ./examples/00_connect        # Run an example
go build ./...                      # Build
go vet ./...                        # Lint
```

## OpenD Simulator (for testing without a real account)

```bash
# Terminal 1: run the simulator (in futuapi4go repo)
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2: run any example
go run ./examples/00_connect
```

## Project Structure

```
futuapi4go-demo/
├── examples/                  # 50 standalone examples (00-49)
│   ├── README.md              # Example descriptions & links
│   ├── 00_connect/           # client.Connect
│   ├── 01_quote/             # client.GetQuote
│   ├── 02_ticker/           # chanpkg.SubscribeTicker
│   ├── 03_orderbook/        # chanpkg.SubscribeOrderBook
│   ├── 04_rt/               # chanpkg.SubscribeRT
│   ├── 05_broker/           # chanpkg.SubscribeBroker
│   ├── 06_kline_single/     # client.GetKLines
│   ├── 07_kline_multi/      # chanpkg.SubscribeKLines
│   └── ... (40 more: 08-49)
├── docs/
│   └── FUTU_PROTO_REF.md
├── AGENTS.md
└── README.md
```

## Examples

Run any example:
```bash
go run ./examples/00_connect
go run ./examples/01_quote
go run ./examples/07_kline_multi
# ... 47 more in examples/
```

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
