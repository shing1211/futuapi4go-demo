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
├── examples/                  # 81 standalone examples (00-80)
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

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **futuapi4go-demo** (1074 symbols, 1113 relationships, 1 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/futuapi4go-demo/context` | Codebase overview, check index freshness |
| `gitnexus://repo/futuapi4go-demo/clusters` | All functional areas |
| `gitnexus://repo/futuapi4go-demo/processes` | All execution flows |
| `gitnexus://repo/futuapi4go-demo/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
