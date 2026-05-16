# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Documentation overhaul**: Rewrote README.md (~120 lines, -56%) with TOC, grouped example categories, and cleaner structure. Created `docs/EXAMPLES.md` with full 96-example reference. Removed redundant `examples/README.md` and outdated `examples/00_ws_connect/README.md`.
- **ARCHITECTURE.md** — Updated SDK version (v0.5.17 → v0.6.2), example count (81 → 96), added quant strategies (91–95) to area tables and Mermaid diagram.
- **SECURITY.md** — Added 0.6.x to supported versions.

## [0.6.2] - 2026-05-16

### Changed

- **Upgraded futuapi4go to v0.6.2** — SHA1 validation test infrastructure,
  improved comments on WritePacketEncrypted (OpenD accepts both SHA1(plaintext)
  and SHA1(ciphertext)), StrictSHA1 mode in mock server.
  Full changelog: https://github.com/shing1211/futuapi4go/blob/main/CHANGELOG.md
- **go.mod** — updated SDK dependency to v0.6.2

## [0.5.12] - 2026-05-14

### Changed

- **Upgraded futuapi4go to v0.5.12** — FTAES_ECB AES encryption for API requests, pool buffer + sync.Cond bug fixes. Local replace directive removed — now using released SDK version.
- **README.md** — updated SDK badge to v0.5.12
- **examples/pkg/connect** — HA connect fix: return on first successful connection, don't retry same host with opposite RSA flag

## [0.5.11] - 2026-05-14

### Changed

- **Upgraded futuapi4go to v0.5.11** — SHA1 plaintext fix for RSA encrypted InitConnect ([#15](https://github.com/shing1211/futuapi4go/issues/15)). RSA connections to remote OpenD now work correctly.
  - SHA1 computed over plaintext before encryption (server verifies decrypted plaintext)
  - Removed incorrect `packetEncAlgo=0` re-marshal that corrupted body
  - Added `WritePacketWithSHA1` to both TCP and WebSocket implementations
- **00_rsa_connect** — updated to use v0.5.11 with SHA1 fix, tested against remote gateway
- **README.md** — updated SDK badge to v0.5.11, added 00_rsa_connect example, RSA env vars

## [0.5.7] - 2026-05-11

### Changed

- **Upgraded futuapi4go to v0.5.7** — Futu OpenD API v10.5.6508 support

## [0.5.6] - 2026-05-07

### Added

- **66_multi_symbol_kline** — `Subscribe` + `GetKLines` + `RequestHistoryKL` batch queries
- **67_order_lifecycle** — Full order workflow: `PlaceOrder` → `GetOrderList` → `ModifyOrder`
- **68_market_hours_check** — `GetMarketState` + `GetTradeDate` for market timing
- **69_subscribe_handler** — Push handlers for Ticker/KLine/OrderBook streams
- **70_futures_account_list** — Futures account support via `GetAccList(TrdCategory_Future)`
- **71_futures_cash** — Futures margin and cash queries
- **72_futures_positions** — `GetPositionList(TrdMarket_Futures)`
- **73_options_account_list** — Options rights check via `GetAccList`
- **74_options_cash** — Options buying power and margin
- **75_options_positions** — Stock + options combined positions
- **76_pre_trade_checks** — Market state + funds + position + quote + snapshot validation
- **77_realtime_dashboard** — Real-time price monitoring with ticker subscriptions
- **78_dca_grid_bot** — Dollar Cost Averaging + Grid trading strategy
- **79_momentum_scanner** — StockFilter + Snapshot + K-lines momentum analysis
- **80_vwap_executor** — OrderBook + VWAP calculation + execution planning
- **81 examples total** — now includes examples 66-80 for gap fill, advanced combos, futures & options
- **go.mod** — updated dependency to `github.com/shing1211/futuapi4go v0.5.6`

### Changed

- **README.md** — updated SDK version badge to v0.5.6, added v0.5.6 changelog
- **AGENTS.md** — corrected example count (66 → 81)

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