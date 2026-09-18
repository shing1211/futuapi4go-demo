# Testing Guide

> **Note:** Examples in this demo run against a live Futu OpenD instance.
> There are no mock servers or isolated unit tests — every example exercises
> real SDK ↔ OpenD communication.

## Running Examples

Examples require a running OpenD instance. They are **not** integration tests
in the CI sense — they are standalone programs you run manually or via the
`go run` command.

```bash
# Set OpenD address (default: 127.0.0.1:11111)
export FUTU_ADDR=127.0.0.1:11111   # Linux/macOS
# or
$env:FUTU_ADDR="127.0.0.1:11111"   # Windows PowerShell

# Run any example
go run ./examples/01_quote
go run ./examples/18_account_list
go run ./examples/22_place_order
```

## Prerequisites

| Requirement | Details |
|-------------|---------|
| OpenD running | TCP on `127.0.0.1:11111`, WebSocket on `127.0.0.1:11117` |
| OpenD logged in | User must be logged in to the Futu/moomoo app |
| For trading | `FUTU_TRADE_PWD` set to trading password MD5 hash |

### Verifying OpenD

```bash
# Test TCP connectivity
nc -zv 127.0.0.1 11111
# or Windows PowerShell
Test-NetConnection -ComputerName 127.0.0.1 -Port 11111
```

### Real Trading

Real trading examples require:
1. OpenD logged into a **real** (not simulate) account
2. Trading password set via environment variable:

```powershell
# Windows PowerShell
$env:FUTU_TRADE_PWD = "32-char-md5-hex-string"
```

```bash
# Linux/macOS
export FUTU_TRADE_PWD=32-char-md5-hex-string
```

## CI / Build Verification

The CI pipeline (`ci.yml`) runs on every push and PR. It verifies:

```bash
go vet ./...           # No vet errors
go build ./...         # All examples compile
go test -race ./...    # Tests with race detection
```

These checks do **not** require OpenD — they only verify that every
`main.go` in `examples/` compiles without error.

```bash
# Local equivalent of CI checks
go vet ./...
go build ./...
go test -race ./...
```

## Example Categories by OpenD Requirement

| Category | Examples | Needs Real Account |
|----------|----------|--------------------|
| Connection | `00_*` | No |
| Market data | `01–17` | No (simulate OK) |
| Basic trading | `18–54` | Simulate or real |
| Order/modify/cancel | `22, 27, 54, 64` | Real account needed |
| Futures/options | `70–75, 81, 97a–e` | Real account needed |
| Flow/history fills | `43, 44, 56, 57, 58` | Real account needed |

## Simulate Trading Limitations

The following APIs **do not work** in simulate mode:

| Example | Function | Error |
|---------|----------|-------|
| `43_order_fill` | GetOrderFillList | 模拟交易不支持成交数据 |
| `44_history_fill` | GetHistoryOrderFillList | 模拟交易不支持成交数据 |
| `56_order_fee` | GetOrderFee | 暂时不支持模拟交易 |
| `57_margin_ratio` | GetMarginRatio | 模拟账户不支持 |
| `58_flow_summary` | GetFlowSummary | 模拟账户不支持查询现金流水 |

For these, use real trading environment with `FUTU_TRADE_PWD` set.

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for the full error reference.

## Writing New Examples

When adding a new example:

1. Create `examples/XXX_name/main.go` where `XXX` is the next available number
2. Use `context.Context` as the first parameter to all API calls
3. Call `cli.Close()` via `defer` after successful connect
4. Print results to stdout — no external logging dependencies
5. Handle errors explicitly — do not silently ignore

```go
package main

import (
    "context"
    "log"

    "github.com/shing1211/futuapi4go/client"
    "github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
    cli := client.New()
    if err := cli.Connect("127.0.0.1:11111"); err != nil {
        log.Fatal(err)
    }
    defer cli.Close()

    ctx := context.Background()
    quote, err := client.GetQuote(ctx, cli, constant.Market_HK, "00700")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("00700: %.2f", quote.Price)
}
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `connection refused` | Start OpenD, check `FUTU_ADDR` |
| `没有解锁交易` | Set `FUTU_TRADE_PWD` env var |
| `请求获取实时K线接口前，请先订阅` | Call `client.Subscribe` before `GetKLines` |
| `模拟交易不支持` | Use real trading with `WithTradeEnv(constant.TrdEnv_Real)` |
| No data for US stock | US stocks require `Subscribe` first |

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for the full reference.
