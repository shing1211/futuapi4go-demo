# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **66_multi_symbol_kline** — `Subscribe` + `GetKLines` + `RequestHistoryKL` batch queries
- **67_order_lifecycle** — Full order workflow: `PlaceOrder` → `GetOrderList` → `ModifyOrder`
- **68_market_hours_check** — `GetMarketState` + `GetTradeDate` for market timing
- **69_subscribe_handler** — Push handlers for Ticker/KLine/OrderBook streams
- **76_pre_trade_checks** — Market state + funds + position + quote + snapshot validation
- **77_realtime_dashboard** — Real-time price monitoring with ticker subscriptions
- **78_dca_grid_bot** — Dollar Cost Averaging + Grid trading strategy
- **79_momentum_scanner** — StockFilter + Snapshot + K-lines momentum analysis
- **80_vwap_executor** — OrderBook + VWAP calculation + execution planning
- **go.mod** — updated dependency to `github.com/shing1211/futuapi4go v0.5.4`

### Changed

- **README.md** — added examples 66-80 for gap fill, advanced combos, and v0.5.4 features

## [0.5.2] - 2026-04-28

### Changed

- **Updated for futuapi4go v0.5.2** — Fluent API, GetHistoryKLPoints, GetUsedQuota
- **go.mod** — updated dependency to `github.com/shing1211/futuapi4go v0.5.2`
- **README.md** — added v0.5.2 examples

## [0.4.0] - 2026-04-28

### Changed

- **Updated for futuapi4go v0.5.1** — all API calls now use context.Context as first parameter
- **go.mod** — updated dependency to `github.com/shing1211/futuapi4go v0.5.1`
- **Typed enums** — examples updated to use typed constants (constant.TrdMarket_HK, constant.TrdMarket_US, etc.)
- **Input validation** — examples updated with proper typed enum values

## [0.3.0] - 2026-04-24

### Changed

- **All trading examples** — now use `FindAccount()` helper for dynamic account selection (no hardcoded account numbers)
- **Trading examples** — use HK stock `00100` (Tencent) instead of US stock to match simulate account market
- **README.md** — completely rewritten with correct example numbers, categories, environment variables, troubleshooting section
- **AGENTS.md** — updated with `FUTU_TRADE_PWD` variable, simulate trading limitations table

### Fixed

- **54_cancel_all_order** — added `UnlockTrading` call + `FUTU_TRADE_PWD` requirement for real trading
- **55_max_trd_qtys** — added `secMarket` parameter (was missing, causing "缺少必要参数证券所属市场")
- **45_acc_trading_info** — changed stock from NVDA to `00100` (Tencent, HK)

## [0.2.1] - 2026-04-24

### Added

- **`chanpkg.SubscribeKLines`** — subscribe to multiple K-line periods with type-safe per-period callbacks (map[KLType]func)

### Fixed

- **`constant.KLType` enum values** — were scrambled (SubType values used instead of KLType values); 5min=2→6, 60min=5→9, Day=6→2, etc.

## [0.2.0] - 2026-04-23

### Changed

- **main.go** — use `SubscribeKLines` for multi-period K-line routing instead of shared single channel
- **go.mod** — `futuapi4go` upgraded to **v0.9.7** with local `replace` directive pointing to `D:/github/futuapi4go`

## [0.1.0] - 2026-04-22

### Added

- 66 standalone examples (00–65) covering all SDK functions
- Multi-market support (HK, US, CN)
- `FUTU_ADDR` environment variable for custom OpenD addresses
- Apache 2.0 License, Contributing guidelines, Code of Conduct, Security policy