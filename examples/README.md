# Examples

> Copy-paste-ready Go programs. Each demonstrates one SDK function.

## Run

```powershell
# Simulator (no account needed)
go run github.com/shing1211/futuapi4go/cmd/examples/simulator
go run ./examples/00_connect

# Real OpenD
set FUTU_ADDR=127.0.0.1:11111
go run ./examples/01_quote
```

## All Examples

### Market Data — Requests

| # | Example | SDK Function |
|---|---------|-------------|
| 01 | [`01_quote`](./01_quote) | `client.GetQuote` |
| 06 | [`06_kline_single`](./06_kline_single) | `client.GetKLines` |
| 08 | [`08_orderbook_req`](./08_orderbook_req) | `client.GetOrderBook` |
| 09 | [`09_ticker_req`](./09_ticker_req) | `client.GetTicker` |
| 10 | [`10_rt_req`](./10_rt_req) | `client.GetRT` |
| 11 | [`11_broker_req`](./11_broker_req) | `client.GetBroker` |
| 15 | [`15_history_kline`](./15_history_kline) | `client.RequestHistoryKL` |
| 17 | [`17_snapshot`](./17_snapshot) | `client.GetSecuritySnapshot` |
| 25 | [`25_trade_date`](./25_trade_date) | `client.GetTradeDate` |
| 30 | [`30_stock_filter`](./30_stock_filter) | `client.StockFilter` |

### Market Data — Real-time Push

| # | Example | SDK Function |
|---|---------|-------------|
| 00 | [`00_connect`](./00_connect) | `client.Connect` |
| 02 | [`02_ticker`](./02_ticker) | `chanpkg.SubscribeTicker` |
| 03 | [`03_orderbook`](./03_orderbook) | `chanpkg.SubscribeOrderBook` |
| 04 | [`04_rt`](./04_rt) | `chanpkg.SubscribeRT` |
| 05 | [`05_broker`](./05_broker) | `chanpkg.SubscribeBroker` |
| 07 | [`07_kline_multi`](./07_kline_multi) | `chanpkg.SubscribeKLines` |

### Market Analysis

| # | Example | SDK Function |
|---|---------|-------------|
| 12 | [`12_capital_flow`](./12_capital_flow) | `client.GetCapitalFlow` |
| 13 | [`13_plate_set`](./13_plate_set) | `client.GetPlateSet` |
| 14 | [`14_plate_stock`](./14_plate_stock) | `client.GetPlateSecurity` |
| 16 | [`16_market_state`](./16_market_state) | `client.GetMarketState` |
| 28 | [`28_owner_plate`](./28_owner_plate) | `client.GetOwnerPlate` |
| 29 | [`29_capital_distribution`](./29_capital_distribution) | `client.GetCapitalDistribution` |

### Trading

| # | Example | SDK Function |
|---|---------|-------------|
| 17 | [`17_global_state`](./17_global_state) | `client.GetGlobalState` |
| 18 | [`18_account_list`](./18_account_list) | `client.GetAccountList` |
| 19 | [`19_funds`](./19_funds) | `client.GetFunds` |
| 20 | [`20_positions`](./20_positions) | `client.GetPositionList` |
| 21 | [`21_unlock_trade`](./21_unlock_trade) | `client.UnlockTrading` |
| 22 | [`22_place_order`](./22_place_order) | `client.PlaceOrder` |
| 23 | [`23_order_list`](./23_order_list) | `client.GetOrderList` |
| 26 | [`26_price_reminder`](./26_price_reminder) | `client.GetPriceReminder` |
| 27 | [`27_cancel_order`](./27_cancel_order) | `client.ModifyOrder` (cancel) |

## Quick Reference

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

## Troubleshooting

- **`connection refused`** — OpenD isn't running. Set `FUTU_ADDR=127.0.0.1:11111`
- **`no positions`** — normal for simulator
- **`subscribe first`** — `GetKLines` requires `client.Subscribe` first
