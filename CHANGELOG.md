# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Example 65:** `65_smoke_test` — fast happy-path example for CI: `GetGlobalState` + `GetAccountList` + `GetQuote(NVDA)` with pass/fail output and nonzero exit on failure
- **`scripts/check-example-numbers.sh`** — validates examples/ numbering (name shape, no gaps, no unexpected duplicates, EXAMPLES.md count); wired into CI
- **`docs/TROUBLESHOOTING.md`** — consolidated error reference (common errors, simulate-mode limitations, known SDK issues, env vars)

### Changed

- `make run-all` now prints a pass/fail summary with failed-dir list and nonzero exit on failure
- `make clean` also removes 3-/4-digit example binaries (`111_*`, `100_*`, etc.)
- README example count → 114; "See Also" now links TROUBLESHOOTING.md
- `examples/pkg/connect` gained unit tests + benchmarks (see `connect_test.go`)

### Pending (tracked in `docs/OPTION_API_COVERAGE_PLAN.md`)

- **Examples 111–129** — Tiers 2–5 demo coverage for institutional/shareholder/flow suites, macro & research data, K-line/rehab variants, quick-trade, and stock-screen:
  - `111_institutional_flow`, `112_shareholder_data`, `113_top_brokers`, `114_short_data`
  - `115_macro_calendar`, `116_macro_indicators`, `117_earnings_dividends`, `118_research_ratings`, `119_ranks_distributions`, `120_search_discovery`, `121_user_security_single`
  - `122_kl_realtime`, `123_history_kl_variants`, `124_rehab_history`, `125_combo_order`, `126_quick_trade`, `127_today_orders_fills`, `128_system_verification`
  - `129_stock_screen`

## [v0.15.2] - 2026-08-12

### Changed

- **Upgraded SDK to v0.15.2** — bumped `github.com/shing1211/futuapi4go` from `v0.15.0` → `v0.15.2`; **dropped** the local `replace` directive — the SDK now resolves from the Go module proxy at the pinned tag. Picks up the SDK-side v0.15.1 additions (push parsers `ParseUpdateOptionEvent`/`ParsePushIndicatorCalc` and the Event Contract channel wrappers).

### Added

- **Example 106:** `GetOptionVolatility` + `GetOptionRank` + `GetOptionMarketStatistic` — implied/historical vol, option ranking by volume, market-wide call/put stat
- **Example 107:** `GetOptionEarningsScreener` + `GetOptionZeroDteScreener` + `GetOptionSellerScreener` + `OptionScreen` — the four option screening engines
- **Example 108:** `GetOptionExerciseProbability` + `GetOptionUnderlyingHisStatistic` + `GetOptionUnderlyingHisVolatility` + `GetOptionUnderlyingOverview` + `GetOptionUnderlyingRank` — option risk & underlying-level analytics
- **Example 109:** `GetOptionEvent` + `GetOptionEventAlert` + `SetOptionEventAlert` — unusual option events + alert CRUD lifecycle (add → deleteAll)
- **Example 110:** `GetIndicatorList` + `RequestIndicatorCalc` — MyLang indicator search and on-demand K-line indicator calculation
- **Documented orphaned 97a–97e sub-examples** in `docs/EXAMPLES.md` (option quote/strategy/analysis/spread, combo order)

## [v0.15.0] - 2026-08-04

### Changed

- **Upgraded SDK to v0.15.0** — Futu Protocol v10.9.6908 (Prediction Markets / Event Contract, 17 new APIs)

### Added

- **Example 102:** `GetEventContractCategory` + `FilterCompetition` + `GetEventContractSeriesList` + `GetEventContractEventList` + `GetEventContract` + `GetEventContractMilestoneList` — explore the full Event Contract hierarchy (category → series → event → contracts → milestones)
- **Example 103:** `GetEventContractSnapshot` + `GetEventContractOrderBook` + `GetEventContractTicker` + `GetEventContractKline` + `RequestHistoryEventContractKL` — real-time and historical EC quotes (YES/NO bid/ask, predSide, klineSource)
- **Example 104:** `SubEventContract` + `RegisterHandler` + `ParseUpdateEventContract{OrderBook,Kline,Ticker}` — streaming EC push subscription & parsing
- **Example 105:** `GetEventContractComboList` + `GetEventContractComboRfq` — combine EC contracts into a Combo and request a firm RFQ quote

## [v0.13.0] - 2026-06-05

### Changed

- **Upgraded SDK to v0.13.0** — Futu Protocol v10.7.6708 (combo options, 6 new APIs)

### Added

- **Example 97a:** `GetOptionQuote` — real-time quotes for option combo legs
- **Example 97b:** `GetOptionStrategy` — available option strategies for an underlying
- **Example 97c:** `GetOptionStrategyAnalysis` — P&L analysis (max profit/loss, Greeks, breakeven)
- **Example 97d:** `GetOptionStrategySpread` — available spread values for vertical spreads
- **Example 97e:** `GetComboMaxTrdQtys` + `PlaceComboOrder` — combo order max qty and placement

## [v0.12.0] - 2026-05-22

### Changed

- **Upgraded SDK to v0.12.0** — ProtoID mismatch fix (CRITICAL), 2 new earnings APIs (GetFinancialsEarningsPriceMove/History), proto safety hardening (30+ fixes), concurrency bug fixes, AES padding fix, breaker halfOpenMax enforcement

## [v0.10.0] - 2026-05-21

### Changed

- **Upgraded SDK to v0.10.0** — Phase 3 proto safety (1,200+ `GetXxx()` → direct field access), input validation, retry skip for trading ops
- **Removed `replace` directive** — demo now uses published `v0.10.0` from GitHub instead of local SDK

## [Unreleased]

(Previously: `replace` directive for local SDK development — now removed in v0.10.0)

## [v0.9.2] - 2026-05-19

## [v0.9.0] - 2026-05-18

### Changed

- **Upgraded SDK to v0.9.0** — `GetBrokerResult`, `WarrantResult`, `OwnerPlateEntry`, `ReconfirmOrderResult.JpAccType`, `GetUserSecurity` all fields

## [v0.8.9] - 2026-05-18

### Changed

- **Upgraded SDK to v0.8.9** — `OrderBookResult`, `TickerResult`, `RTResult` wrappers expose S2C security/name
- **GetOrderBook return type changed** — now returns `*client.OrderBookResult` (use `.Items[0]`)
- **GetTicker return type changed** — now returns `*client.TickerResult` (use `.Items`)
- **GetRT return type changed** — now returns `*client.RTResult` (use `.Items`)

### Added

- **Examples 08, 09, 10, 80, 91, 93** updated for new return types

## [v0.8.8] - 2026-05-18

### Changed

- **Upgraded SDK to v0.8.8** — extended proto field mappings for Quote, Snapshot, IpoList, StaticInfo, SubInfo, OwnerPlate, OptionExpirationDate, StockFilter
- **GetOwnerPlate return type changed** — now returns `[]*client.OwnerPlateInfo` with Code, Name, PlateType instead of just plate names

### Added

- **All request-response examples** updated to display newly available extended data fields

## [v0.8.7] - 2026-05-18

### Changed

- **Upgraded SDK to v0.8.7** — `RequestHistoryKLQuota` API removed from SDK; deleted example 65

### Added

- Compatible with futuapi4go SDK v0.8.7 (extended Quote, Snapshot, UserInfo, RehabInfo, DelayStatistics types)

### Fixed

- **examples/48_subscribe_kline_single** — updated for `SubscribeKLine` signature changes in SDK v0.8.7
- **examples/69_subscribe_handler** — updated for KLine push handler API changes in SDK v0.8.7

## [v0.8.6] - 2026-05-17

### Added

- **JSON dump to all request-response examples** — Added `display.PrintJSON()` output to all request-response examples (00-35, 36-60, 61-85, 86-99) for easy data inspection. Push-subscription examples (streaming) are unchanged.
- **New `examples/pkg/display` package** — Shared JSON pretty-print utility used by all examples.

### Changed

- **97 examples enriched** with JSON output section. Pre-existing push-subscription examples (47, 48, 49, 63, 69, 77, 83, 91) and void-return examples (50, 53, 54) are unchanged.
- **examples/93_smart_money** — removed erroneous `display.PrintJSON(book)` reference from out-of-scope variable.

## [v0.8.5] - 2026-05-17

### Changed

- **examples/12_capital_flow** — updated for SDK `GetCapitalFlow` return type change (`*client.CapitalFlowResponse` with `LastValidTime`/`LastValidTimestamp`).
- **examples/16_market_state** — updated for SDK `GetMarketState` return type change (`*client.MarketStateResult` with `Code`/`Name`/`State`).
- **examples/68_market_hours_check**, **76_pre_trade_checks**, **90_system_diagnostics** — updated for `GetMarketState` return type change (use `.State`).
- **examples/93_smart_money** — updated for `GetCapitalFlow` return type change (use `.Items`).

## [0.8.3] - 2026-05-17

### Changed

- **examples/13_plate_set** — now prints `PlateType` field (newly added in SDK v0.8.3).

## [0.8.2] - 2026-05-17

### Changed

- **SDK dependency bumped** — `github.com/shing1211/futuapi4go` from v0.8.0 to v0.8.2.
- **examples/07_kline_multi** — fixed `SubscribeKLines` call for updated SDK signature (typed `constant.Market` and `[]constant.KLType` parameters).

## [0.8.0] - 2026-05-16

### Added

- **State machine shutdown example** — New `examples/98_state_shutdown` demonstrating connection state monitoring with `WithOnStateChange` callback, `State()` query, and graceful `Shutdown()` with drain timeout.
- **Order audit validation example** — New `examples/99_audit_validation` demonstrating pre-flight `ValidateOrder` checks, `HasErrors` routing, and structured audit logging via `NewAuditLogger`.
- **KL cache example** — New `examples/100_kl_cache` demonstrating LRU+TTL K-line cache with hit/miss behavior, cleanup, and size monitoring.
- **OTel metrics example** — New `examples/101_otel_metrics` demonstrating all OTel meter record functions: connections, API calls, latency, push messages, rate limiting, retry, breaker state.

### Changed

- **New SDK APIs** — uses local SDK with `State()`, `Shutdown()`, `WithOnStateChange()`, `trd.ValidateOrder`, and `trd.NewAuditLogger`.
- **go.mod** — added `replace` directive pointing to local `../futuapi4go` for unreleased SDK features.
- **README.md** — updated example count (97 → 101), added Lifecycle & Audit category, updated Infrastructure & Tracing range.
- **docs/EXAMPLES.md** — added entries for examples 98 and 99.

## [0.7.0] - 2026-05-16

### Added

- **OTel tracing example** — New `examples/97_opentelemetry_tracing` showing OpenTelemetry setup with stdout exporter, TracerProvider, and auto-generated spans from SDK calls.

### Changed

- **Upgraded futuapi4go to v0.7.0** — Adds OpenTelemetry tracing (opt-in), goreleaser release automation, and trilingual package docs.
- **Documentation overhaul**: Rewrote README.md (~120 lines, -56%) with TOC, grouped example categories, and cleaner structure. Created `docs/EXAMPLES.md` with full 96-example reference. Removed redundant `examples/README.md` and outdated `examples/00_ws_connect/README.md`.
- **ARCHITECTURE.md** — Updated SDK version (v0.5.17 → v0.6.2 → v0.7.0), example count (81 → 96 → 97), added quant strategies (91–95) to area tables and Mermaid diagram.
- **SECURITY.md** — Added 0.6.x and 0.7.x to supported versions.
- **go.mod** — updated SDK dependency to v0.6.2 → v0.7.0

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