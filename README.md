# futuapi4go-demo

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  <img src="https://img.shields.io/github/stars/shing1211/futuapi4go-demo" alt="Stars">
</p>

> **Copy-paste-ready Go examples for the futuapi4go SDK.** Each example is a standalone `main.go` demonstrating one SDK function. Run against the simulator or a real OpenD.

## Prerequisites

- **Go 1.26+** — [golang.org/dl](https://golang.org/dl)
- **Futu OpenD** on `127.0.0.1:11111` — [download](https://www.futunn.com/download/fetch-lasted-link?name=opend-windows), or use the simulator below

## Run an Example

```powershell
git clone https://github.com/shing1211/futuapi4go-demo.git
cd futuapi4go-demo

# Pick one (00-49):
go run ./examples/00_connect
go run ./examples/07_kline_multi
go run ./examples/23_place_order
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

## All Examples (00–49)

### Market Data — Snapshot & History

| # | Example | SDK Function |
|---|---------|-------------|
| 01 | [`01_quote`](./examples/01_quote) | `client.GetQuote` |
| 06 | [`06_kline_single`](./examples/06_kline_single) | `client.GetKLines` |
| 08 | [`08_orderbook_req`](./examples/08_orderbook_req) | `client.GetOrderBook` |
| 09 | [`09_ticker_req`](./examples/09_ticker_req) | `client.GetTicker` |
| 10 | [`10_rt_req`](./examples/10_rt_req) | `client.GetRT` |
| 11 | [`11_broker_req`](./examples/11_broker_req) | `client.GetBroker` |
| 15 | [`15_history_kline`](./examples/15_history_kline) | `client.RequestHistoryKL` |
| 17 | [`17_snapshot`](./examples/17_snapshot) | `client.GetSecuritySnapshot` |
| 30 | [`30_trade_date`](./examples/30_trade_date) | `client.GetTradeDate` |
| 36 | [`36_ipo_list`](./examples/36_ipo_list) | `client.GetIpoList` |
| 37 | [`37_future_info`](./examples/37_future_info) | `client.GetFutureInfo` |
| 38 | [`38_suspend`](./examples/38_suspend) | `client.GetSuspend` |
| 40 | [`40_rehab`](./examples/40_rehab) | `client.RequestRehab` |
| 41 | [`41_code_change`](./examples/41_code_change) | `client.GetCodeChange` |

### Market Data — Real-time Push

| # | Example | SDK Function |
|---|---------|-------------|
| 00 | [`00_connect`](./examples/00_connect) | `client.Connect` |
| 02 | [`02_ticker`](./examples/02_ticker) | `chanpkg.SubscribeTicker` |
| 03 | [`03_orderbook`](./examples/03_orderbook) | `chanpkg.SubscribeOrderBook` |
| 04 | [`04_rt`](./examples/04_rt) | `chanpkg.SubscribeRT` |
| 05 | [`05_broker`](./examples/05_broker) | `chanpkg.SubscribeBroker` |
| 07 | [`07_kline_multi`](./examples/07_kline_multi) | `chanpkg.SubscribeKLines` |
| 16 | [`16_market_state`](./examples/16_market_state) | `client.GetMarketState` |
| 47 | [`47_subscribe_quote`](./examples/47_subscribe_quote) | `chanpkg.SubscribeQuote` |
| 48 | [`48_subscribe_kline_single`](./examples/48_subscribe_kline_single) | `chanpkg.SubscribeKLine` |

### Market Analysis

| # | Example | SDK Function |
|---|---------|-------------|
| 12 | [`12_capital_flow`](./examples/12_capital_flow) | `client.GetCapitalFlow` |
| 13 | [`13_plate_set`](./examples/13_plate_set) | `client.GetPlateSet` |
| 14 | [`14_plate_stock`](./examples/14_plate_stock) | `client.GetPlateSecurity` |
| 32 | [`32_owner_plate`](./examples/32_owner_plate) | `client.GetOwnerPlate` |
| 33 | [`33_capital_distribution`](./examples/33_capital_distribution) | `client.GetCapitalDistribution` |
| 34 | [`34_stock_filter`](./examples/34_stock_filter) | `client.StockFilter` |
| 35 | [`35_reference`](./examples/35_reference) | `client.GetReference` |
| 39 | [`39_holding_change`](./examples/39_holding_change) | `client.GetHoldingChangeList` |

### Trading

| # | Example | SDK Function |
|---|---------|-------------|
| 18 | [`18_global_state`](./examples/18_global_state) | `client.GetGlobalState` |
| 19 | [`19_account_list`](./examples/19_account_list) | `client.GetAccountList` |
| 20 | [`20_funds`](./examples/20_funds) | `client.GetFunds` |
| 21 | [`21_positions`](./examples/21_positions) | `client.GetPositionList` |
| 22 | [`22_unlock_trade`](./examples/22_unlock_trade) | `client.UnlockTrading` |
| 23 | [`23_place_order`](./examples/23_place_order) | `client.PlaceOrder` |
| 24 | [`24_order_list`](./examples/24_order_list) | `client.GetOrderList` |
| 25 | [`25_cancel_order`](./examples/25_cancel_order) | `client.ModifyOrder` (cancel) |
| 26 | [`26_history_order`](./examples/26_history_order) | `client.GetHistoryOrderList` |
| 27 | [`27_order_fill`](./examples/27_order_fill) | `client.GetOrderFillList` |
| 28 | [`28_history_fill`](./examples/28_history_fill) | `client.GetHistoryOrderFillList` |
| 29 | [`29_acc_trading_info`](./examples/29_acc_trading_info) | `client.GetAccTradingInfo` |

### Derivatives

| # | Example | SDK Function |
|---|---------|-------------|
| 42 | [`42_warrant`](./examples/42_warrant) | `client.GetWarrant` |
| 43 | [`43_option_expiration`](./examples/43_option_expiration) | `client.GetOptionExpirationDate` |
| 44 | [`44_option_chain`](./examples/44_option_chain) | `client.GetOptionChain` |

### Alerts & User Data

| # | Example | SDK Function |
|---|---------|-------------|
| 31 | [`31_price_reminder`](./examples/31_price_reminder) | `client.GetPriceReminder` |
| 45 | [`45_user_security`](./examples/45_user_security) | `client.GetUserSecurity` |
| 46 | [`46_user_info`](./examples/46_user_info) | `client.GetUserInfo` |
| 49 | [`49_subscribe_price_reminder`](./examples/49_subscribe_price_reminder) | `chanpkg.SubscribePriceReminder` |

## Common Patterns

```go
// Market constant — all APIs take int32
int32(constant.Market_US)  // 11
int32(constant.Market_HK)  // 1

// Request: one-shot, returns immediately
client.GetQuote(ctx, cli, int32(constant.Market_US), "NVDA")

// Subscribe: continuous stream, call returned stop() to unsubscribe
stop := chanpkg.SubscribeTicker(cli, int32(constant.Market_US), "NVDA", tickerCh)
defer stop()
```

## Known Caveats

- **`GetDelayStatistics`** — skipped due to a proto2/proto3 wire-format mismatch. All other APIs work normally.
- **`subscribe first`** — `GetKLines`, `GetQuote`, `GetTicker`, `GetRT`, `GetBroker`, `GetOrderBook` require `client.Subscribe` first.
- **`no positions`** — normal for simulator; real OpenD shows actual positions.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
