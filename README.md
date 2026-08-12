# futuapi4go-demo

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  <img src="https://img.shields.io/github/stars/shing1211/futuapi4go-demo" alt="Stars">
  <img src="https://img.shields.io/badge/futuapi4go-v0.15.0-00ADD8?style=flat-square" alt="SDK Version">
</p>

> **⚠️ Under Active Development**  
> This demo project is under active development alongside the futuapi4go SDK. Examples
> may reference APIs or types still being finalized. Not recommended for production use.

> **Production-ready Go examples for the [futuapi4go](https://github.com/shing1211/futuapi4go) SDK.** 113 standalone examples (00–105), covering all SDK functions and advanced trading strategies. All examples tested and verified against the OpenD simulator.

## Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Example Categories](#example-categories)
- [Environment Variables](#environment-variables)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)
- [See Also](#see-also)
- [License](#license)

## Quick Start

```bash
# Clone and run
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

# Run an example
go run ./examples/00_connect
# Output: Connected to local OpenD @ 127.0.0.1:11111

go run ./examples/01_quote
# Output: HK.00700 -> 380.20 / 380.40 | high=385.00 low=376.50 vol=15,234,567
```

### Real Trading

```bash
$env:FUTU_TRADE_PWD="32-char-md5-hex-string"
go run ./examples/54_cancel_all_order
```

## Prerequisites

- Go 1.26+
- A running [Futu OpenD](https://www.futunn.com/en/overview) instance (default: `127.0.0.1:11111`)
- For remote connections: RSA key pair (see `00_rsa_connect`)
- For real trading: account unlocked with trade password MD5 hash

## Project Structure

```
futuapi4go-demo/
├── examples/               # 108 standalone programs (00–105)
│   ├── 00_connect/        # client.Connect
│   ├── 00_rsa_connect/    # TCP + RSA encryption
│   ├── 00_ws_connect/     # WebSocket connection
│   ├── 01_quote/          # GetQuote
│   ├── ...                # (see docs/EXAMPLES.md for full list)
│   ├── 99_audit_validation/
│   ├── 100_kl_cache/      # K-line cache (LRU+TTL)
│   └── 102–105_ec_*/      # Prediction Market / Event Contract
└── examples/pkg/          # Shared helpers
    └── connect/           # MustConnect, ManagedConnection (HA)
```

See **[docs/EXAMPLES.md](docs/EXAMPLES.md)** for the complete 108-example reference.

## Example Categories

| Category | Count | Range | Purpose |
|----------|-------|-------|---------|
| Connection | 3 | `00_*` | Plain TCP, RSA-encrypted TCP, WebSocket |
| Basic Functions | 64 | `01–64` | One-shot & streaming market data, trading, subscriptions |
| Gap Fill | 4 | `66–69` | Multi-function workflows |
| Futures & Options | 6 | `70–75` | Futures/options accounts, positions, margin |
| Advanced Combo | 5 | `76–80` | Pre-trade checks, DCA grid bot, VWAP |
| Advanced Trading | 5 | `81–85` | Options trading, trailing stop, risk analysis |
| Infrastructure & Tracing | 8 | `86–90, 97, 100–101` | History downloader, quota manager, diagnostics, OTel tracing, KL cache, OTel metrics |
| Quant Strategies | 5 | `91–95` | Order book imbalance, pairs trading, smart money |
| Special | 1 | `96` | Delay statistics |
| Lifecycle & Audit | 2 | `98–99` | State machine shutdown, order audit validation |
| Prediction Markets | 4 | `102–105` | Event Contract discovery, market data, subscribe & push, combo RFQ |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD server address | `127.0.0.1:11111` |
| `FUTU_TRADE_PWD` | MD5 hash of trading password (32 chars) | (not set) |
| `FUTU_RSA_PUBKEY` | RSA public key PEM for remote encrypted connections | (not set) |
| `FUTU_WS_ADDR` | WebSocket OpenD address | `127.0.0.1:11113` |
| `FUTU_WS_SECRET` | WebSocket secret key | (not set) |

## Common Patterns

```go
// Create client (default: simulate trading)
cli := client.New()
defer cli.Close()
cli.Connect("127.0.0.1:11111")

// Real trading: use constant.TrdEnv_Real
cli := client.New().WithTradeEnv(constant.TrdEnv_Real)

// Market values — typed (no magic numbers)
constant.Market_US        // 11
constant.Market_HK        // 1
constant.TrdMarket_HK     // 1
constant.TrdMarket_Futures // 5

// Futures accounts: use TrdCategory_Future
futuresAccounts, _ := cli.Trade().GetAccList(ctx, TrdCategory_Future)

// One-shot request
client.GetQuote(ctx, cli, constant.Market_US, "NVDA")

// Subscribe: continuous stream, call stop() to unsubscribe
stop := chanpkg.SubscribeTicker(ctx, cli, constant.Market_US, "NVDA", tickerCh)
defer stop()

// Dynamic account selection (no hardcoded account numbers)
accounts, _ := client.GetAccountList(ctx, cli)
acc := cli.FindAccount(accounts)
accID := acc.AccID
```

## Troubleshooting

See [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) for the consolidated
error reference, simulate-trading limitations, and known SDK issues.

| Error | Cause |
|-------|-------|
| `connection refused` | OpenD not running. Set `FUTU_ADDR=127.0.0.1:11111` |
| no data from `GetKLines`, `GetQuote`, etc. | Call `client.Subscribe` first for push-type data |
| `账户购买力不足` | Simulate account has no buying power — expected |
| `模拟交易不支持` | Function not supported in simulate mode — use real trading |
| `未知的协议ID` | OpenD doesn't implement this API (e.g. ReconfirmOrder) |
| `没有解锁交易，请先解锁交易` | Need to unlock trading with `FUTU_TRADE_PWD` env var |
| `请求获取实时K线接口前，请先订阅` | Must subscribe to K-line type before calling `GetKLines` |
| `暂不提供美股 OTC 市场行情` | Some US stocks are OTC and not supported — skip with error handling |

### Known Caveats

- **`GetDelayStatistics`** — SDK v0.5.13+ includes a custom proto2 marshaling workaround for the repeated-int32 encoding issue. The demo (96) handles both success and failure gracefully. May still fail on older OpenD builds that don't implement the API.
- **US stocks** — require `client.Subscribe` before `GetQuote` returns data. HK stocks do not.
- **Simulate trading** — many order/flow APIs are not supported. Use real trading environment (`WithTradeEnv(constant.TrdEnv_Real)`) with `FUTU_TRADE_PWD` set.
- **Futures accounts** — use `GetAccList(TrdCategory_Future)`, not `GetAccountList` which only returns stock/options accounts.
- **secMarket required** — `PlaceOrder` and `GetMaxTrdQtys` require explicit `TrdSecMarket` parameter.

## See Also

- **[futuapi4go](https://github.com/shing1211/futuapi4go)** — the Go SDK this demo is built on
- [Full Example Reference](docs/EXAMPLES.md) — complete table of all 114 examples
- [Troubleshooting](docs/TROUBLESHOOTING.md) — common errors, sim-mode limitations, known SDK issues
- [Architecture](ARCHITECTURE.md) — design overview, execution flows, Mermaid diagram
- [CHANGELOG](CHANGELOG.md) — version history and release notes

## License

Apache License 2.0 — see [LICENSE](LICENSE).

> **Trading Disclaimer**: Trading financial instruments carries significant risk. Always test thoroughly in simulate mode before using real funds.
