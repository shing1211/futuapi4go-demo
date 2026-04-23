# futuapi4go-demo

[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/shing1211/futuapi4go-demo)](https://goreportcard.com/report/github.com/shing1211/futuapi4go-demo)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logo=blue)](https://pkg.go.dev/github.com/shing1211/futuapi4go-demo)
[![GitHub stars](https://img.shields.io/github/stars/shing1211/futuapi4go-demo)](https://github.com/shing1211/futuapi4go-demo/stargazers)

> Interactive Go demo showcasing every public API in the [futuapi4go](https://github.com/shing1211/futuapi4go) SDK — market data, K-lines, order books, real-time push, trading, options, and more.

## Features

- **10 demo categories** covering all major SDK functions
- **Menu-driven interface** — run one demo or all at once
- **Real-time push subscriptions** — live quote, K-line, order book, and ticker updates
- **Multi-market** — Hong Kong, US
- **Color terminal output** — green/red/yellow for quick visual parsing
- **Zero config** — reads `FUTU_ADDR` env var for custom OpenD addresses

## Requirements

- Go 1.26+
- [Futu OpenD](https://www.futunn.com/download/fetch-lasted-link?name=opend-windows) running on `127.0.0.1:11111` (default)

## Quick Start

```bash
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

go run ./cmd/demo/main.go
```

Connect to a custom OpenD instance:

```bash
FUTU_ADDR=192.168.1.100:11111 go run ./cmd/demo/main.go
```

Run without a real account using the built-in mock simulator:

```bash
# Terminal 1: start the simulator (requires futuapi4go source at ../futuapi4go)
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2: run the demo
go run ./cmd/demo/main.go
```

## Demo Menu

| # | Category | APIs Demonstrated |
|---|----------|-------------------|
| 1 | Connection & System | `GetGlobalState`, `GetUserInfo` (GetDelayStatistics skipped — known OpenD incompatibility) |
| 2 | Market Data | `GetBasicQot`, `GetKL`, `GetOrderBook`, `GetTicker`, `GetRT`, `GetBroker`, `GetSecuritySnapshot` |
| 3 | Market Analysis | `GetPlateSet`, `GetPlateSecurity`, `GetCapitalFlow`, `GetCapitalDistribution`, `GetOwnerPlate`, `GetReference`, `GetStaticInfo`, `GetFutureInfo`, `StockFilter` |
| 4 | Stock Screening | `StockFilter` — multi-criteria queries |
| 5 | Options & Warrants | `GetOptionExpirationDate`, `GetOptionChain`, `GetWarrant` |
| 6 | Historical Data | `RequestHistoryKL`, `GetHistoryKL`, `RequestHistoryKLQuota`, `GetRehab` |
| 7 | Corporate Actions | `GetIpoList`, `GetCodeChange`, `GetSuspend`, `GetHoldingChangeList` |
| 8 | Trading Operations | `GetAccList`, `GetFunds`, `GetPositionList`, `GetMaxTrdQtys`, `GetOrderList`, `GetOrderFillList`, `GetHistoryOrderList`, `PlaceOrder`, `GetFlowSummary` |
| 9 | User Groups & Alerts | `GetUserSecurityGroup`, `GetUserSecurity`, `SetPriceReminder`, `GetPriceReminder` |
| 10 | Real-time Push | Live `BasicQot`, `K-line`, `OrderBook`, `Ticker` subscriptions |
| 0 | Run All | Executes demos 1–9 in sequence |

## API Coverage

### Market Data (`pkg/qot`)
| API | Description |
|-----|-------------|
| `GetBasicQot` | Real-time quotes for multiple securities |
| `GetKL` | K-line (candlestick) snapshot |
| `GetOrderBook` | Bid/ask depth (order book) |
| `GetTicker` | Tick-by-tick trade data |
| `GetRT` | Intraday time-share (分时数据) |
| `GetBroker` | Broker queue (港股买卖盘经纪) |
| `GetSecuritySnapshot` | Full snapshot with extended market data |
| `GetTradeDate` | Trading calendar |

### Market Analysis (`pkg/qot`)
| API | Description |
|-----|-------------|
| `GetPlateSet` | List sector/plate sets |
| `GetPlateSecurity` | Stocks belonging to a plate |
| `GetCapitalFlow` | Money flow (inflow/outflow) |
| `GetCapitalDistribution` | Capital tier distribution |
| `GetOwnerPlate` | Plates owning a given stock |
| `GetReference` | Related securities (futures, options, warrants) |
| `GetStaticInfo` | Security static metadata |
| `GetFutureInfo` | Futures contract details |
| `StockFilter` | Multi-criteria stock screener |

### Options & Warrants (`pkg/qot`)
| API | Description |
|-----|-------------|
| `GetOptionExpirationDate` | Option expiry dates for an underlying |
| `GetOptionChain` | Full option chain with strikes/expiries |
| `GetWarrant` | Warrant data with filtering/sorting |

### Historical Data (`pkg/qot`)
| API | Description |
|-----|-------------|
| `RequestHistoryKL` | Paginated history K-lines (auto-pagination supported) |
| `GetHistoryKL` | Single time-range K-line query |
| `RequestHistoryKLQuota` | Check remaining API quota |
| `GetRehab` | Dividend/rights adjustment factors |

### Corporate Actions (`pkg/qot`)
| API | Description |
|-----|-------------|
| `GetIpoList` | IPO calendar |
| `GetCodeChange` | Stock splits / mergers |
| `GetSuspend` | Trading suspension (停牌) information |
| `GetHoldingChangeList` | Major shareholder changes |

### Trading (`pkg/trd`)
| API | Description |
|-----|-------------|
| `GetAccList` | List all trading accounts |
| `GetFunds` | Account balance / funds |
| `GetPositionList` | Open positions |
| `GetOrderList` | Active (today's) orders |
| `GetOrderFillList` | Today's fills/executions |
| `GetHistoryOrderList` | Historical orders |
| `GetMaxTrdQtys` | Maximum tradable quantities |
| `PlaceOrder` | Place or modify orders |
| `GetFlowSummary` | Daily fund flow summary |

### System (`pkg/sys`)
| API | Description |
|-----|-------------|
| `GetGlobalState` | Connection state and server info |
| `GetUserInfo` | User account details |
| `GetDelayStatistics` | Performance metrics (skipped in demo — known proto2 wire-format incompatibility) |

### User Data (`pkg/qot`)
| API | Description |
|-----|-------------|
| `GetUserSecurityGroup` | Watchlist groups |
| `GetUserSecurity` | Securities in a group |
| `ModifyUserSecurity` | Add/remove securities from a group |
| `SetPriceReminder` | Create/modify/delete price alerts |
| `GetPriceReminder` | Fetch existing alerts |

### Real-time Push (`pkg/qot` + `pkg/push`)
| API | Description |
|-----|-------------|
| `Subscribe` | Subscribe to real-time data streams |
| `UnsubscribeAll` | Cancel all subscriptions |

## Known Issues

- **`GetDelayStatistics`** — skipped in demo due to a proto2/proto3 wire-format mismatch between the Go SDK and OpenD's C++ protobuf parser. `google.golang.org/protobuf` uses packed encoding for `repeated int32` by default, but OpenD may not support this for proto2 messages.
- **`GetTradeDate`** — may fail if the SDK doesn't populate all required C2S fields.

See [AGENTS.md](AGENTS.md) for SDK debugging details.

## Project Structure

```
futuapi4go-demo/
├── cmd/demo/main.go          # Source code (single-file, menu-driven)
├── docs/                     # Supplementary docs
│   └── FUTU_PROTO_REF.md    # Slimmed proto API reference
├── scripts/                  # Build & run scripts
│   ├── build.bat / .sh      # Build binary to cmd/demo/
│   ├── run.bat / .sh         # Run the demo
│   ├── clean.bat / .sh      # Clean build artifacts
│   └── upgrade.bat / .sh     # Upgrade Go dependencies
├── .github/                  # GitHub config
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
├── AGENTS.md               # AI agent instructions
├── README.md
├── LICENSE                 # Apache 2.0
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── go.mod
└── go.sum
```

## Build & Run

```bash
# Build
scripts\build.bat          # Windows
./scripts/build.sh         # Linux/macOS

# Run (defaults to 127.0.0.1:11111, override with FUTU_ADDR env var)
scripts\run.bat            # Windows
./scripts/run.sh           # Linux/macOS

# Or run directly with go run
FUTU_ADDR=127.0.0.1:11111 go run ./cmd/demo/main.go

# Lint
go build ./...
go vet ./...
```

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

For security concerns, please see [SECURITY.md](SECURITY.md). **Never commit credentials or `.env` files.**

## License

Licensed under the [Apache License 2.0](LICENSE).
