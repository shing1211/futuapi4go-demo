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
├── examples/                  # 66 standalone examples (00-65)
│   ├── README.md              # Example descriptions & links
│   ├── 00_connect/           # client.Connect
│   ├── 01_quote/             # client.GetQuote
│   ├── 02_ticker/           # chanpkg.SubscribeTicker
│   ├── 03_orderbook/        # chanpkg.SubscribeOrderBook
│   ├── 04_rt/               # chanpkg.SubscribeRT
│   ├── 05_broker/           # chanpkg.SubscribeBroker
│   ├── 06_kline_single/     # client.GetKLines
│   ├── 07_kline_multi/      # chanpkg.SubscribeKLines
│   └── ... (59 more: 08-65)
├── docs/
│   └── FUTU_PROTO_REF.md
├── AGENTS.md
└── README.md
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD server address | `127.0.0.1:11111` |
| `FUTU_TRADE_PWD` | MD5 hash of trading password (32 hex chars) | (not set) |

## Trading Modes

The SDK defaults to **simulate trading** (`TrdEnv=0`). To use real trading:

```go
cli := client.New().WithTradeEnv(1) // Real trading
```

Real trading requires `FUTU_TRADE_PWD` environment variable with MD5 hash of your trading password.

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

### GetDelayStatistics — may have proto2 wire-format incompatibility with certain OpenD versions

OpenD may reject the `GetDelayStatistics` request with "解析protobuf协议失败". Root cause: `google.golang.org/protobuf` encodes `repeated int32` fields using proto3 packed wire format by default, but some OpenD C++ parsers expect proto2 non-packed encoding.

**Workaround in demo:** The call is skipped with a printed note. All other APIs work normally with OpenD v10.4.6408.

### GetTradeDate — all C2S fields are required

`GetTradeDate` has all required fields in its C2S. If the SDK doesn't populate all required fields, OpenD returns "解析protobuf协议失败". Works correctly with OpenD v10.4.6408.

**Workaround in demo:** If this API fails, the demo exits with a red error.

**Proto reference:** See `docs/FUTU_PROTO_REF.md` or https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md

## Simulate Trading Limitations

The following APIs are **not supported** in simulate trading mode:

| Example | Function | Error |
|---------|----------|-------|
| 43_order_fill | GetOrderFillList | 模拟交易不支持成交数据 |
| 44_history_fill | GetHistoryOrderFillList | 模拟交易不支持成交数据 |
| 56_order_fee | GetOrderFee | 暂时不支持模拟交易 |
| 57_margin_ratio | GetMarginRatio | 模拟账户不支持 |
| 58_flow_summary | GetFlowSummary | 模拟账户不支持查询现金流水 |
| 64_reconfirm_order | ReconfirmOrder | 未知的协议ID (OpenD doesn't implement) |

For these, use real trading environment (`WithTradeEnv(1)`) with `FUTU_TRADE_PWD` set.

## Related Repositories

- SDK: `github.com/shing1211/futuapi4go` (checked out at `D:\github\futuapi4go`)
- Official Proto Doc: https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md
- OpenD Downloads: https://www.futunn.com/download/fetch-lasted-link?name=opend-windows
