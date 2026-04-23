# futuapi4go-demo

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  <img src="https://img.shields.io/github/stars/shing1211/futuapi4go-demo" alt="Stars">
</p>

> **Copy-paste-ready Go examples for the futuapi4go SDK.** Each example is a standalone `main.go` demonstrating one SDK function. No account needed — run against the built-in simulator.

## Prerequisites

- **Go 1.26+** — [golang.org/dl](https://golang.org/dl)
- **Futu OpenD** on `127.0.0.1:11111` — [download](https://www.futunn.com/download/fetch-lasted-link?name=opend-windows), or use the simulator below

## Run an Example

```powershell
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

# Pick one:
go run ./examples/00_connect
go run ./examples/01_quote
go run ./examples/02_ticker
go run ./examples/03_orderbook
go run ./examples/04_rt
go run ./examples/05_broker
go run ./examples/06_kline_single
go run ./examples/07_kline_multi
```

### Custom OpenD Address

```powershell
set FUTU_ADDR=192.168.1.100:11111
go run ./examples/01_quote
```

### Simulator (No Account Needed)

```powershell
# Terminal 1
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2
go run ./examples/07_kline_multi
```

## Examples Overview

| # | Example | SDK Function | Description |
|---|---------|-------------|-------------|
| 00 | `00_connect` | `client.Connect` | Connect and disconnect |
| 01 | `01_quote` | `client.GetQuote` | Snapshot quote |
| 02 | `02_ticker` | `chanpkg.SubscribeTicker` | Real-time tick trades |
| 03 | `03_orderbook` | `chanpkg.SubscribeOrderBook` | Order book (bids & asks) |
| 04 | `04_rt` | `chanpkg.SubscribeRT` | Tick-by-tick time & sales |
| 05 | `05_broker` | `chanpkg.SubscribeBroker` | Broker queue |
| 06 | `06_kline_single` | `client.GetKLines` | Historical K-lines (one-shot) |
| 07 | `07_kline_multi` | `chanpkg.SubscribeKLines` | Live K-lines for multiple periods |

## Project Layout

```
futuapi4go-demo/
├── examples/
│   ├── README.md              # Example descriptions & usage guide
│   ├── 00_connect/
│   ├── 01_quote/
│   ├── 02_ticker/
│   ├── 03_orderbook/
│   ├── 04_rt/
│   ├── 05_broker/
│   ├── 06_kline_single/
│   └── 07_kline_multi/
├── docs/
│   └── FUTU_PROTO_REF.md      # Proto field reference
├── build.bat / .sh            # Build
├── run.bat / .sh              # Run (uses run.bat in examples/)
├── clean.bat / .sh            # Clean
├── upgrade.bat / .sh          # Upgrade SDK
├── AGENTS.md                  # AI agent instructions
├── README.md
└── LICENSE
```

## SDK Reference

The SDK mirrors Python SDK naming conventions:

```go
import "github.com/shing1211/futuapi4go/pkg/constant"

// Markets: constant.Market_HK, constant.Market_US, constant.Market_SH
// K-Lines: constant.KLType_K_Day, constant.KLType_K_1Min
// Subscriptions: constant.SubType_Quote, constant.SubType_K_1Min
// Trading: constant.TrdEnv_Simulate, constant.TrdSide_Buy
```

See the full [Python Migration Guide](https://github.com/shing1211/futuapi4go/blob/main/PYTHON_MIGRATION_GUIDE.md).

## Build & Scripts

```powershell
.\build.bat      # Build
.\run.bat        # Run (default: 00_connect)
.\clean.bat      # Clean artifacts
.\upgrade.bat    # Upgrade futuapi4go dependency
```

## Known Caveats

- **`GetDelayStatistics`** — skipped due to a proto2/proto3 wire-format mismatch between Go's protobuf library and OpenD's C++ parser. All other APIs work normally.
- **`GetTradeDate`** — may return an error on older OpenD versions if all required C2S fields aren't populated.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
