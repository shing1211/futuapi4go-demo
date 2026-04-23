# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2026-04-24

### Changed

- **main.go** — use `SubscribeKLines` for multi-period K-line routing instead of shared single channel
- **go.mod** — `futuapi4go` upgraded to **v0.9.7** with local `replace` directive pointing to `D:/github/futuapi4go`

## [0.2.0] - 2026-04-23

### Changed

- **README.md** — rewritten with vivid style, visual demo table, color output description, vibrant layout
- **examples/README.md** — rewritten with vivid style, friendly tone, clear troubleshooting section
- **AGENTS.md** — corrected script paths (`scripts\` → root), expanded SDK debugging section
- **CONTRIBUTING.md** — polished with local SDK setup instructions
- **build.bat / .sh**, **run.bat / .sh**, **clean.bat / .sh**, **upgrade.bat / .sh** — fixed `cd` path bug (`%cd%..` → `%cd%`)

### Dependencies

- **futuapi4go** upgraded to **v0.9.0** — removes local replace directive, pulls from `proxy.golang.org`

## [0.1.0] - 2026-04-22

### Added

- Interactive menu-driven demo with 10 categories covering all major APIs
- Connection & System demos (`GetGlobalState`, `GetUserInfo`)
- Market Data demos (quote, K-line, order book, tick, broker, snapshot)
- Market Analysis demos (plates, capital flow, stock filter)
- Options & Warrants demos (option chain, expiry dates, warrant data)
- Historical Data demos (K-line history, rehab, quota)
- Corporate Actions demos (IPO, splits, suspensions)
- Trading Operations demos (accounts, positions, orders, fills, flow)
- Watchlists & Alerts demos (groups, price reminders)
- Real-time Push Subscriptions (BasicQot, K-line, OrderBook, Ticker)
- Multi-market support (HK, US)
- `FUTU_ADDR` environment variable for custom OpenD addresses
- Apache 2.0 License, Contributing guidelines, Code of Conduct, Security policy
