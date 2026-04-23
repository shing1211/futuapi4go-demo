# futuapi4go-demo

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  <img src="https://img.shields.io/github/stars/shing1211/futuapi4go-demo" alt="Stars">
</p>

> **Spin up the full Futu OpenAPI in your terminal.** This demo walks you through every market data, trading, and real-time push API in the futuapi4go SDK — no account needed to explore.

## What You'll Get

Walk through **10 interactive demo categories** that cover the complete Futu OpenAPI surface:

| # | Category | What's Inside |
|---|---------|-------------|
| 1 | **Connection & System** | OpenD state, user info, server version |
| 2 | **Market Data** | Real-time quotes, K-lines, order book, tick data, broker queue |
| 3 | **Market Analysis** | Sector plates, capital flow, owner plates, stock filter |
| 4 | **Stock Screening** | Multi-criteria screener with 20+ filter fields |
| 5 | **Options & Warrants** | Option chains, expiry dates, warrant data |
| 6 | **Historical Data** | K-line history, rehab data, API quota |
| 7 | **Corporate Actions** | IPO calendar, stock splits, trading halts |
| 8 | **Trading Operations** | Accounts, positions, orders, fills, fund flow |
| 9 | **Watchlists & Alerts** | Custom groups, price reminders |
| 10 | **Real-time Push** | Live quote, K-line, order book, ticker streams |
| 0 | **Run All** | Execute demos 1–9 in one shot |

## Prerequisites

- **Go 1.26+** — [golang.org/dl](https://golang.org/dl)
- **Futu OpenD** running on `127.0.0.1:11111` — [download here](https://www.futunn.com/download/fetch-lasted-link?name=opend-windows)

Don't have OpenD or a trading account? **No problem.** Use the built-in mock simulator instead (see below).

## Get Up and Running

```bash
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

go run ./cmd/demo/main.go
```

The interactive menu will greet you. Pick a category (or `0` to run everything) and watch the data flow.

### Custom OpenD Address

```bash
FUTU_ADDR=192.168.1.100:11111 go run ./cmd/demo/main.go
```

### No Account? Use the Simulator

```bash
# Terminal 1 — fire up the mock OpenD
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2 — run the demo
go run ./cmd/demo/main.go
```

The simulator returns realistic mock data for all APIs so you can explore without risking a single dollar.

## Build Scripts

```bash
# Build the binary
.\build.bat          # Windows
./build.sh           # Linux / macOS

# Run the demo
.\run.bat            # Windows
./run.sh             # Linux / macOS

# Clean up
.\clean.bat          # Windows
./clean.sh           # Linux / macOS

# Upgrade SDK to latest
.\upgrade.bat        # Windows
./upgrade.sh         # Linux / macOS
```

## Python SDK? The Constants Are Familiar

The demo (and the full SDK) mirror Python SDK naming conventions so the migration feels natural:

```go
import "github.com/shing1211/futuapi4go/pkg/constant"

// Markets: constant.Market_HK, constant.Market_US, constant.Market_SH
// K-Lines: constant.KLType_K_Day, constant.KLType_K_1Min
// Trading: constant.TrdEnv_Simulate, constant.TrdSide_Buy
// Subscriptions: constant.SubType_Quote, constant.SubType_K_1Min
```

See the full [Python Migration Guide](https://github.com/shing1211/futuapi4go/blob/main/PYTHON_MIGRATION_GUIDE.md) for side-by-side comparisons of every API.

## Color Terminal Output

Every demo prints with ANSI color coding for quick visual parsing:

- **Green** — buy side, positive change, success
- **Red** — sell side, negative change, error
- **Yellow** — warnings, skipped sections, partial data
- **Cyan** — section headers and key metrics

## Project Layout

```
futuapi4go-demo/
├── cmd/demo/main.go           # Single-file interactive demo (~1,500 lines)
├── examples/
│   ├── getting_started/       # First steps: connect → quote → K-line → subscribe
│   └── trading_demo/          # Full trading flow: accounts → positions → orders
├── docs/
│   └── FUTU_PROTO_REF.md     # Proto field reference for all APIs
├── build.bat / .sh            # Build binary to cmd/demo/
├── run.bat / .sh              # Run the demo
├── clean.bat / .sh            # Remove build artifacts
├── upgrade.bat / .sh          # Upgrade futuapi4go dependency
├── .github/                  # Issue templates, PR template
├── AGENTS.md                  # AI agent instructions
├── README.md
└── LICENSE                   # Apache 2.0
```

## API Coverage at a Glance

### Market Data
`GetBasicQot` · `GetKL` · `GetOrderBook` · `GetTicker` · `GetRT` · `GetBroker` · `GetSecuritySnapshot` · `GetTradeDate`

### Market Analysis
`GetPlateSet` · `GetPlateSecurity` · `GetCapitalFlow` · `GetCapitalDistribution` · `GetOwnerPlate` · `GetReference` · `GetStaticInfo` · `GetFutureInfo` · `StockFilter`

### Options & Warrants
`GetOptionExpirationDate` · `GetOptionChain` · `GetWarrant`

### Historical Data
`RequestHistoryKL` · `GetHistoryKL` · `RequestHistoryKLQuota` · `GetRehab`

### Corporate Actions
`GetIpoList` · `GetCodeChange` · `GetSuspend` · `GetHoldingChangeList`

### Trading
`GetAccList` · `GetFunds` · `GetPositionList` · `GetMaxTrdQtys` · `GetOrderList` · `GetOrderFillList` · `GetHistoryOrderList` · `PlaceOrder` · `ModifyOrder` · `GetFlowSummary`

### System & User
`GetGlobalState` · `GetUserInfo` · `GetUserSecurityGroup` · `GetUserSecurity` · `SetPriceReminder` · `GetPriceReminder`

### Real-time Push
`Subscribe` · `UnsubscribeAll` · `RegQotPush`

## Known Caveats

- **`GetDelayStatistics`** — skipped in the demo due to a proto2/proto3 wire-format mismatch between Go's protobuf library and OpenD's C++ parser. All other APIs work normally.
- **`GetTradeDate`** — may return an error on older OpenD versions if all required C2S fields aren't populated.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). All contributions welcome — new API demos, better output formatting, test coverage, you name it.

## License

Apache License 2.0 — see [LICENSE](LICENSE).
