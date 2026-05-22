# futuapi4go-demo AGENTS.md

## Project

Go demo showcasing the futuapi4go SDK. Each example is a standalone `main.go` demonstrating one SDK function.

## Dev Commands

```bash
go run ./examples/00_connect        # Run an example
go build ./...                      # Build
go vet ./...                        # Lint
```

## Running Examples

Run examples against a real OpenD instance:

```bash
# Set OpenD address (default: 127.0.0.1:11111)
export FUTU_ADDR="127.0.0.1:11111"

# Run any example
go run ./examples/00_connect
```

## Project Structure

```
futuapi4go-demo/
├── examples/                  # 104+ standalone examples (00-101)
│   ├── README.md              # Example descriptions & links
│   ├── 00_connect/           # client.Connect
│   ├── 00_connect_ha/        # HA (high-availability) connect
│   ├── 00_rsa_connect/       # RSA encrypted connect
│   ├── 01_quote/             # client.GetQuote
│   ├── 02_ticker/           # chanpkg.SubscribeTicker
│   ├── ... (up to 101)
│   └── pkg/                  # Shared packages (connect/, display/)
├── AGENTS.md
├── go.mod                    # SDK v0.9.2 + replace directive
└── README.md
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD server address | `127.0.0.1:11111` |
| `FUTU_TRADE_PWD` | MD5 hash of trading password (32 hex chars) | (not set) |
| `FUTU_RSA_PUBKEY` | RSA public key PEM for remote encrypted connections | (not set) |
| `FUTU_WS_ADDR` | WebSocket OpenD address | `127.0.0.1:11113` |
| `FUTU_WS_SECRET` | WebSocket secret key | (not set) |

## Trading Modes

The SDK defaults to **simulate trading** (`TrdEnv=0`). To use real trading:

```go
cli := client.New().WithTradeEnv(constant.TrdEnv_Real)
```

Real trading requires `FUTU_TRADE_PWD` environment variable with MD5 hash of your trading password.

## SDK Dependency

The futuapi4go SDK source is at `github.com/shing1211/futuapi4go`.

- Proto files: `api/proto/` in the SDK repo
- Generated Go protobuf code: `pkg/pb/`
- SDK source: `pkg/`

### Local `replace` Directive (Required)

`go.mod` includes a `replace` directive pointing to `../futuapi4go` because:

- The Go module proxy cached a stale version of **v0.9.2** (the tag was force-pushed)
- New SDK packages (`pkg/cache`, `pkg/history`, `pkg/option`, `pkg/tracing`, `pkg/tracing/otel`) exist in the local repo but aren't served by the proxy
- Examples 86-101 import these packages

If the SDK repo is at a different path, update the `replace` directive in `go.mod`:

```
replace github.com/shing1211/futuapi4go => /path/to/local/futuapi4go
```

After editing `go.mod`, re-download modules:

```
go mod download
```

## Known SDK Issues

### GetDelayStatistics — SDK fix applied (v0.5.13+)

Some OpenD C++ parsers reject Go's default packed encoding for `repeated int32` fields. The SDK (v0.5.13+) includes a custom proto2 marshaling workaround in `pkg/sys/system.go` (`marshalC2SProto2`). The request body is hand-encoded using proto2 non-packed varint format.

**Demo behavior:** Example 96 calls `GetDelayStatistics` and handles both success and failure gracefully. If OpenD rejects the request (older build or missing API support), the demo prints a clear explanation and exits cleanly.

### GetTradeDate — all C2S fields are required

`GetTradeDate` has all required fields in its C2S. If the SDK doesn't populate all required fields, OpenD returns "解析protobuf协议失败". Works correctly with OpenD v10.5.6508.

**Workaround in demo:** If this API fails, the demo exits with a red error.

**Proto reference:** https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md

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

- SDK: `github.com/shing1211/futuapi4go` (current: v0.12.0)
- Official Proto Doc: https://openapi.futunn.com/mds/Futu-API-Doc-zh-Proto.md
- OpenD Downloads: https://www.futunn.com/download/fetch-lasted-link?name=opend-windows

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **futuapi4go-demo** (1577 symbols, 3975 relationships, 133 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

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
