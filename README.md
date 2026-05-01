# futuapi4go-demo

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  <img src="https://img.shields.io/github/stars/shing1211/futuapi4go-demo" alt="Stars">
  <img src="https://img.shields.io/badge/futuapi4go-v0.5.4-00ADD8?style=flat-square" alt="SDK Version">
</p>

> **Production-ready Go examples for the [futuapi4go](https://github.com/shing1211/futuapi4go) SDK.** 80 standalone examples (00–80), covering all SDK functions and advanced trading strategies. All examples tested and verified against the OpenD simulator.

## v0.5.4

```go
// Futures account support
cli.Trade().GetAccList(ctx, TrdCategory_Future) // Separate from stocks

// US stock momentum scanner
client.StockFilter(ctx, cli, Market_US, 0, 30)    // Screen stocks
client.GetSecuritySnapshot(ctx, cli, securities) // Enrich with snapshot
client.GetKLines(ctx, cli, Market_US, code, KLType_K_Day, 10) // K-line analysis
```

## Quick Start

```powershell
# Clone and run
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

# Run an example (80 examples: 00–80)
go run ./examples/00_connect
go run ./examples/01_quote
go run ./examples/79_momentum_scanner

# Custom OpenD address
$env:FUTU_ADDR="192.168.1.100:11111"
go run ./examples/01_quote
```

### Simulator (no account needed)

```powershell
# Terminal 1: start the simulator
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2: run any example
go run ./examples/07_kline_multi
```

### Real Trading (requires unlocked account)

```powershell
# Set trading password (MD5 hash of your trade password)
$env:FUTU_TRADE_PWD="32-char-md5-hex-string"
go run ./examples/54_cancel_all_order
```

## Project Structure

```
futuapi4go-demo/
├── examples/           # 80 standalone examples (00–80), one main.go each
│   ├── 00_connect/     # client.Connect
│   ├── 01_quote/       # client.GetQuote
│   ├── ...
│   ├── 66_multi_symbol_kline/   # Batch K-line retrieval
│   ├── 67_order_lifecycle/       # Full order workflow
│   ├── 70_futures_account_list/  # Futures accounts
│   ├── 79_momentum_scanner/      # Stock screening + analysis
│   └── 80_vwap_executor/          # VWAP execution strategy
└── docs/
    └── FUTU_PROTO_REF.md  # Futu OpenAPI protobuf reference
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD server address | `127.0.0.1:11111` |
| `FUTU_TRADE_PWD` | MD5 hash of trading password (32 chars) | (not set) |

## All Examples (00–80)

### Basic Function Examples (00-65)

| # | Example | SDK Function |
|---|---------|-------------|
| 00 | [`00_connect`](./examples/00_connect) | `client.Connect` |
| 01 | [`01_quote`](./examples/01_quote) | `client.GetQuote` |
| 02 | [`02_ticker`](./examples/02_ticker) | `chanpkg.SubscribeTicker` |
| 03 | [`03_orderbook`](./examples/03_orderbook) | `chanpkg.SubscribeOrderBook` |
| 04 | [`04_rt`](./examples/04_rt) | `chanpkg.SubscribeRT` |
| 05 | [`05_broker`](./examples/05_broker) | `chanpkg.SubscribeBroker` |
| 06 | [`06_kline_single`](./examples/06_kline_single) | `client.GetKLines` |
| 07 | [`07_kline_multi`](./examples/07_kline_multi) | `chanpkg.SubscribeKLines` |
| 08 | [`08_orderbook_req`](./examples/08_orderbook_req) | `client.GetOrderBook` |
| 09 | [`09_ticker_req`](./examples/09_ticker_req) | `client.GetTicker` |
| 10 | [`10_rt_req`](./examples/10_rt_req) | `client.GetRT` |
| 11 | [`11_broker_req`](./examples/11_broker_req) | `client.GetBroker` |
| 12 | [`12_capital_flow`](./examples/12_capital_flow) | `client.GetCapitalFlow` |
| 13 | [`13_plate_set`](./examples/13_plate_set) | `client.GetPlateSet` |
| 14 | [`14_plate_stock`](./examples/14_plate_stock) | `client.GetPlateSecurity` |
| 15 | [`15_history_kline`](./examples/15_history_kline) | `client.RequestHistoryKL` |
| 16 | [`16_market_state`](./examples/16_market_state) | `client.GetMarketState` |
| 17 | [`17_global_state`](./examples/17_global_state) | `client.GetGlobalState` |
| 18 | [`18_account_list`](./examples/18_account_list) | `client.GetAccountList` |
| 19 | [`19_account_list`](./examples/19_account_list) | `client.GetAccountInfo` |
| 20 | [`20_funds`](./examples/20_funds) | `client.GetFunds` |
| 21 | [`21_positions`](./examples/21_positions) | `client.GetPositionList` |
| 22 | [`22_place_order`](./examples/22_place_order) | `client.PlaceOrder` |
| 23 | [`23_order_list`](./examples/23_order_list) | `client.GetOrderList` |
| 24 | [`24_snapshot`](./examples/24_snapshot) | `client.GetSecuritySnapshot` |
| 25 | [`25_trade_date`](./examples/25_trade_date) | `client.GetTradeDate` |
| 26 | [`26_price_reminder`](./examples/26_price_reminder) | `client.GetPriceReminder` |
| 27 | [`27_cancel_order`](./examples/27_cancel_order) | `client.ModifyOrder` (cancel) |
| 28 | [`28_owner_plate`](./examples/28_owner_plate) | `client.GetOwnerPlate` |
| 29 | [`29_capital_distribution`](./examples/29_capital_distribution) | `client.GetCapitalDistribution` |
| 30 | [`30_stock_filter`](./examples/30_stock_filter) | `client.StockFilter` |
| 31 | [`31_ipo_list`](./examples/31_ipo_list) | `client.GetIpoList` |
| 32 | [`32_future_info`](./examples/32_future_info) | `client.GetFutureInfo` |
| 33 | [`33_suspend`](./examples/33_suspend) | `client.GetSuspend` |
| 34 | [`34_holding_change`](./examples/34_holding_change) | `client.GetHoldingChangeList` |
| 35 | [`35_rehab`](./examples/35_rehab) | `client.RequestRehab` |
| 36 | [`36_code_change`](./examples/36_code_change) | `client.GetCodeChange` |
| 37 | [`37_warrant`](./examples/37_warrant) | `client.GetWarrant` |
| 38 | [`38_option_chain`](./examples/38_option_chain) | `client.GetOptionChain` |
| 39 | [`39_option_expiration`](./examples/39_option_expiration) | `client.GetOptionExpirationDate` |
| 40 | [`40_reference`](./examples/40_reference) | `client.GetReference` |
| 41 | [`41_user_security`](./examples/41_user_security) | `client.GetUserSecurityGroup` |
| 42 | [`42_history_order`](./examples/42_history_order) | `client.GetHistoryOrderList` |
| 43 | [`43_order_fill`](./examples/43_order_fill) | `client.GetOrderFillList` |
| 44 | [`44_history_fill`](./examples/44_history_fill) | `client.GetHistoryOrderFillList` |
| 45 | [`45_acc_trading_info`](./examples/45_acc_trading_info) | `client.GetAccTradingInfo` |
| 46 | [`46_user_info`](./examples/46_user_info) | `client.GetUserInfo` |
| 47 | [`47_subscribe_quote`](./examples/47_subscribe_quote) | `chanpkg.SubscribeQuote` |
| 48 | [`48_subscribe_kline_single`](./examples/48_subscribe_kline_single) | `chanpkg.SubscribeKLine` |
| 49 | [`49_subscribe_price_reminder`](./examples/49_subscribe_price_reminder) | `chanpkg.SubscribePriceReminder` |
| 50 | [`50_unsubscribe`](./examples/50_unsubscribe) | `client.Unsubscribe` |
| 51 | [`51_unsubscribe_all`](./examples/51_unsubscribe_all) | `client.UnsubscribeAll` |
| 52 | [`52_query_subscription`](./examples/52_query_subscription) | `client.QuerySubscription` |
| 53 | [`53_reg_qot_push`](./examples/53_reg_qot_push) | `client.RegQotPush` |
| 54 | [`54_cancel_all_order`](./examples/54_cancel_all_order) | `client.CancelAllOrder` |
| 55 | [`55_max_trd_qtys`](./examples/55_max_trd_qtys) | `client.GetMaxTrdQtys` |
| 56 | [`56_order_fee`](./examples/56_order_fee) | `client.GetOrderFee` |
| 57 | [`57_margin_ratio`](./examples/57_margin_ratio) | `client.GetMarginRatio` |
| 58 | [`58_flow_summary`](./examples/58_flow_summary) | `client.GetFlowSummary` |
| 59 | [`59_static_info`](./examples/59_static_info) | `client.GetStaticInfo` |
| 60 | [`60_modify_user_security`](./examples/60_modify_user_security) | `client.ModifyUserSecurity` |
| 61 | [`61_sub_info`](./examples/61_sub_info) | `client.GetSubInfo` |
| 62 | [`62_set_price_reminder`](./examples/62_set_price_reminder) | `client.SetPriceReminder` |
| 63 | [`63_sub_acc_push`](./examples/63_sub_acc_push) | `client.SubAccPush` |
| 64 | [`64_reconfirm_order`](./examples/64_reconfirm_order) | `client.ReconfirmOrder` |
| 65 | [`65_history_kl_quota`](./examples/65_history_kl_quota) | `client.RequestHistoryKLQuota` |

### Gap Fill Examples (66-69)

| # | Example | SDK Functions |
|---|---------|---------------|
| 66 | [`66_multi_symbol_kline`](./examples/66_multi_symbol_kline) | `Subscribe` + `GetKLines` + `RequestHistoryKL` |
| 67 | [`67_order_lifecycle`](./examples/67_order_lifecycle) | `PlaceOrder` → `GetOrderList` → `ModifyOrder` |
| 68 | [`68_market_hours_check`](./examples/68_market_hours_check) | `GetMarketState` + `GetTradeDate` |
| 69 | [`69_subscribe_handler`](./examples/69_subscribe_handler) | `Subscribe` + push handlers (Ticker/KLine/OrderBook) |

### Futures & Options Examples (70-75)

| # | Example | Description |
|---|---------|-------------|
| 70 | [`70_futures_account_list`](./examples/70_futures_account_list) | `GetAccList(TrdCategory_Future)` for futures accounts |
| 71 | [`71_futures_cash`](./examples/71_futures_cash) | Futures margin and cash queries |
| 72 | [`72_futures_positions`](./examples/72_futures_positions) | `GetPositionList(TrdMarket_Futures)` |
| 73 | [`73_options_account_list`](./examples/73_options_account_list) | Options rights check via `GetAccList` |
| 74 | [`74_options_cash`](./examples/74_options_cash) | Options buying power and margin |
| 75 | [`75_options_positions`](./examples/75_options_positions) | Stock + options combined positions |

### Advanced Combo Examples (76-80)

| # | Example | Description |
|---|---------|-------------|
| 76 | [`76_pre_trade_checks`](./examples/76_pre_trade_checks) | Market state + account funds + position validation |
| 77 | [`77_realtime_dashboard`](./examples/77_realtime_dashboard) | Real-time monitoring with ticker subscriptions |
| 78 | [`78_dca_grid_bot`](./examples/78_dca_grid_bot) | Dollar Cost Averaging + Grid strategy |
| 79 | [`79_momentum_scanner`](./examples/79_momentum_scanner) | StockFilter + Snapshot + K-lines momentum analysis |
| 80 | [`80_vwap_executor`](./examples/80_vwap_executor) | OrderBook + VWAP calculation + execution |

## Common Patterns

```go
// Create client (default: simulate trading)
cli := client.New()
defer cli.Close()
cli.Connect("127.0.0.1:11111")

// Real trading: use WithTradeEnv(1)
cli := client.New().WithTradeEnv(1) // Real trading

// Market constant — typed constant (no cast needed)
constant.Market_US // 11
constant.Market_HK // 1
constant.TrdMarket_HK // 1 — HK trading market
constant.TrdMarket_US // 2 — US trading market
constant.TrdMarket_Futures // 5 — Futures market

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

## Known Caveats

- **`GetDelayStatistics`** — skipped. Known proto2/proto3 wire-format mismatch with OpenD serverVer=1003. See SDK's CHANGELOG.
- **`GetTradeDate`** — requires OpenD serverVer >= 1004 for proto2 field compatibility. Use `RequestTradeDate` as a fallback.
- **US stocks** — require `client.Subscribe` before `GetQuote` returns data. HK stocks do not.
- **Simulate trading** — many order/flow APIs are not supported. Use real trading environment (`WithTradeEnv(1)`) with `FUTU_TRADE_PWD` set.
- **Futures accounts** — use `GetAccList(TrdCategory_Future)`, not `GetAccountList` which only returns stock/options accounts.
- **secMarket required** — `PlaceOrder` and `GetMaxTrdQtys` require explicit `TrdSecMarket` parameter.

## License

Apache License 2.0 — see [LICENSE](LICENSE).