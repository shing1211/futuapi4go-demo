# futuapi4go-demo

[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/shing1211/futuapi4go-demo)](https://goreportcard.com/report/github.com/shing1211/futuapi4go-demo)
[![Go Reference](https://pkg.go.dev/badge/github.com/shing1211/futuapi4go-demo.svg)](https://pkg.go.dev/github.com/shing1211/futuapi4go-demo)

> A lean, interactive demo showcasing every public API in the [futuapi4go](https://github.com/shing1211/futuapi4go) SDK — real-time quotes, K-lines, order books, trading, options, and more.

## Features

- **10 demo categories** covering all major SDK functions
- **Menu-driven interface** — run one demo or all at once
- **Real-time push subscriptions** with live quote, K-line, order book, and ticker updates
- **Multi-market** — Hong Kong, US, China A-shares

## Quick Start

```bash
# Clone the repo
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

# Prerequisites: Futu OpenD running on 127.0.0.1:11111
# Or start the mock simulator:
#   go run github.com/shing1211/futuapi4go/cmd/simulator

go run main.go
```

Set a custom address:

```bash
FUTU_ADDR=192.168.1.100:11111 go run main.go
```

## Demo Menu

| # | Category | What's Shown |
|---|----------|-------------|
| 1 | Connection & System | Global state, user info, delay stats, market status, trade dates |
| 2 | Market Data | Quotes, daily K-lines, order book, tick data, intraday RT, broker queue, snapshot |
| 3 | Market Analysis | Plate sets, plate members, capital flow/distribution, owner plates, reference securities, static info, futures info |
| 4 | Stock Screening | StockFilter with multi-criteria queries |
| 5 | Options & Warrants | Option expiry dates, option chain, warrant data |
| 6 | Historical Data | RequestHistoryKL (paginated + auto), GetHistoryKL, K-line quota, rehab |
| 7 | Corporate Actions | IPO list, code changes, suspension info, holder changes |
| 8 | Trading Operations | Account list, funds, positions, max qty, order list, fills, history, place order, flow summary |
| 9 | User Groups & Alerts | Security groups, price reminders |
| 10 | Real-time Push | Live BasicQot, K-line, OrderBook, Ticker subscriptions |
| 0 | Run All | Executes demos 1–9 in sequence |

## API Coverage

### Market Data (`pkg/qot`)
- `GetBasicQot` — Real-time quotes
- `GetKL` — K-line (candlestick) snapshot
- `GetOrderBook` — Bid/ask depth
- `GetTicker` — Tick-by-tick trades
- `GetRT` — Intraday time-share
- `GetBroker` — Broker queue
- `GetSecuritySnapshot` — Full snapshot
- `GetTradeDate` — Trading calendar

### Market Analysis (`pkg/qot`)
- `GetPlateSet` — Sector/plate list
- `GetPlateSecurity` — Stocks in a plate
- `GetCapitalFlow` — Money flow
- `GetCapitalDistribution` — Capital tier distribution
- `GetOwnerPlate` — Plates owning a stock
- `GetReference` — Related securities
- `GetStaticInfo` — Security metadata
- `GetFutureInfo` — Futures contract details
- `StockFilter` — Multi-criteria screener

### Historical Data (`pkg/qot`)
- `RequestHistoryKL` — Auto-paginated or manual pagination
- `GetHistoryKL` — Single time-range query
- `RequestHistoryKLQuota` — API quota check
- `GetRehab` — Dividend adjustment factors

### Corporate Actions (`pkg/qot`)
- `GetIpoList` — IPO calendar
- `GetCodeChange` — Stock splits / mergers
- `GetSuspend` — Halt information
- `GetHoldingChangeList` — Major holder changes

### Options & Warrants (`pkg/qot`)
- `GetOptionExpirationDate` — Option expiry dates
- `GetOptionChain` — Option chain
- `GetWarrant` — Warrant data

### Trading (`pkg/trd`)
- `GetAccList` — Account list
- `GetFunds` — Account balance
- `GetPositionList` — Open positions
- `GetOrderList` — Active orders
- `GetOrderFillList` — Recent fills
- `GetHistoryOrderList` — Historical orders
- `GetMaxTrdQtys` — Max tradable quantities
- `PlaceOrder` — Place order
- `GetFlowSummary` — Daily fund flow

### System (`pkg/sys`)
- `GetGlobalState` — Connection & server info
- `GetUserInfo` — User account details
- `GetDelayStatistics` — Performance metrics

### User Data (`pkg/qot`)
- `GetUserSecurityGroup` — Watchlist groups
- `GetUserSecurity` — Securities in a group
- `ModifyUserSecurity` — Add/remove from group
- `SetPriceReminder` — Create price alert
- `GetPriceReminder` — Fetch alerts

### Push Subscriptions (`pkg/qot` + `pkg/push`)
- `Subscribe` — Subscribe to real-time data
- `UnsubscribeAll` — Unsubscribe from all

## Project Structure

```
futuapi4go-demo/
├── main.go              # All demos, menu-driven
├── go.mod               # Module definition
├── go.sum               # Dependency checksums
├── .gitignore           # Git ignore rules
├── LICENSE              # Apache 2.0
├── README.md            # This file
├── CONTRIBUTING.md      # Contribution guidelines
├── CODE_OF_CONDUCT.md   # Community standards
├── SECURITY.md          # Security policy
├── CHANGELOG.md         # Release history
└── .github/
    ├── ISSUE_TEMPLATE/
    │   ├── bug_report.md
    │   └── feature_request.md
    └── PULL_REQUEST_TEMPLATE.md
```

## Requirements

- Go 1.26+
- Futu OpenD running on the target address (default `127.0.0.1:11111`)
- Or use the built-in mock simulator: `go run github.com/shing1211/futuapi4go/cmd/simulator`

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

For security concerns, please see [SECURITY.md](SECURITY.md).

## License

Licensed under the [Apache License 2.0](LICENSE).
