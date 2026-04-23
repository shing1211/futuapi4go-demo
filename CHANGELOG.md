# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- README.md: enriched with API coverage tables, badges, troubleshooting, known issues section
- AGENTS.md: corrected project structure, expanded SDK debugging documentation
- CONTRIBUTING.md: added local SDK setup, mock simulator usage, new demo category guide
- docs/FUTU_PROTO_REF.md: added table of contents, trading API proto definitions
- main.go: fixed 4 `go vet` warnings (wrong format verbs, redundant newlines)
- main.go: removed CNSH/CNSZ stocks (demo only uses HK and US now), removed unused `MarketCNSH`/`MarketCNSZ` constants
- main.go: removed `GetDelayStatistics` call from demo (known proto2 wire-format incompatibility with OpenD)
- scripts: added `run.bat` / `run.sh`, improved all scripts with proper `cd`, `setlocal`/`set -e`, correct exit codes

### Added

- README.md, CHANGELOG.md, CONTRIBUTING.md, CODE_OF_CONDUCT.md, LICENSE, SECURITY.md: moved from `docs/` to project root for standard layout
- docs/README.md: removed (merged into root README.md)
- docs/IMPLEMENTATION_PLAN.md: pre-existing implementation tracking document (unchanged)

## [0.1.0] - 2026-04-22

### Added

- Interactive menu-driven demo with 10 categories
- Connection & System demos (GetGlobalState, GetUserInfo, GetDelayStatistics)
- Market Data demos (GetBasicQot, GetKL, GetOrderBook, GetTicker, GetRT, GetBroker, GetSecuritySnapshot, GetTradeDate)
- Market Analysis demos (GetPlateSet, GetPlateSecurity, GetCapitalFlow, GetCapitalDistribution, GetOwnerPlate, GetReference, GetStaticInfo, GetFutureInfo, StockFilter)
- Options & Warrants demos (GetOptionExpirationDate, GetOptionChain, GetWarrant)
- Historical Data demos (RequestHistoryKL, GetHistoryKL, RequestHistoryKLQuota, GetRehab)
- Corporate Actions demos (GetIpoList, GetCodeChange, GetSuspend, GetHoldingChangeList)
- Trading Operations demos (GetAccList, GetFunds, GetPositionList, GetMaxTrdQtys, GetOrderList, GetOrderFillList, GetHistoryOrderList, PlaceOrder, GetFlowSummary)
- User Groups & Alerts demos (GetUserSecurityGroup, GetUserSecurity, SetPriceReminder, GetPriceReminder)
- Real-time Push Subscriptions (BasicQot, K-line, OrderBook, Ticker)
- Multi-market support (HK, US, China A-shares)
- Custom OpenD address via `FUTU_ADDR` environment variable
- Apache 2.0 License
- Contributing guidelines
- Code of Conduct
- Security policy

[0.1.0]: https://github.com/shing1211/futuapi4go-demo/releases/tag/v0.1.0
