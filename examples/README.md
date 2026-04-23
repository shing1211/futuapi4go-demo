# Examples

> Each `main.go` demonstrates one SDK function. Run with `go run ./examples/NN_name`.

## Market Data — Snapshot & History

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
| 30 | [`30_trade_date`](./30_trade_date) | `client.GetTradeDate` |
| 36 | [`36_ipo_list`](./36_ipo_list) | `client.GetIpoList` |
| 37 | [`37_future_info`](./37_future_info) | `client.GetFutureInfo` |
| 38 | [`38_suspend`](./38_suspend) | `client.GetSuspend` |
| 40 | [`40_rehab`](./40_rehab) | `client.RequestRehab` |
| 41 | [`41_code_change`](./41_code_change) | `client.GetCodeChange` |

## Market Data — Real-time Push

| # | Example | SDK Function |
|---|---------|-------------|
| 00 | [`00_connect`](./00_connect) | `client.Connect` |
| 02 | [`02_ticker`](./02_ticker) | `chanpkg.SubscribeTicker` |
| 03 | [`03_orderbook`](./03_orderbook) | `chanpkg.SubscribeOrderBook` |
| 04 | [`04_rt`](./04_rt) | `chanpkg.SubscribeRT` |
| 05 | [`05_broker`](./05_broker) | `chanpkg.SubscribeBroker` |
| 07 | [`07_kline_multi`](./07_kline_multi) | `chanpkg.SubscribeKLines` |
| 16 | [`16_market_state`](./16_market_state) | `client.GetMarketState` |
| 47 | [`47_subscribe_quote`](./47_subscribe_quote) | `chanpkg.SubscribeQuote` |
| 48 | [`48_subscribe_kline_single`](./48_subscribe_kline_single) | `chanpkg.SubscribeKLine` |

## Market Analysis

| # | Example | SDK Function |
|---|---------|-------------|
| 12 | [`12_capital_flow`](./12_capital_flow) | `client.GetCapitalFlow` |
| 13 | [`13_plate_set`](./13_plate_set) | `client.GetPlateSet` |
| 14 | [`14_plate_stock`](./14_plate_stock) | `client.GetPlateSecurity` |
| 32 | [`32_owner_plate`](./32_owner_plate) | `client.GetOwnerPlate` |
| 33 | [`33_capital_distribution`](./33_capital_distribution) | `client.GetCapitalDistribution` |
| 34 | [`34_stock_filter`](./34_stock_filter) | `client.StockFilter` |
| 35 | [`35_reference`](./35_reference) | `client.GetReference` |
| 39 | [`39_holding_change`](./39_holding_change) | `client.GetHoldingChangeList` |

## Trading

| # | Example | SDK Function |
|---|---------|-------------|
| 18 | [`18_global_state`](./18_global_state) | `client.GetGlobalState` |
| 19 | [`19_account_list`](./19_account_list) | `client.GetAccountList` |
| 20 | [`20_funds`](./20_funds) | `client.GetFunds` |
| 21 | [`21_positions`](./21_positions) | `client.GetPositionList` |
| 22 | [`22_unlock_trade`](./22_unlock_trade) | `client.UnlockTrading` |
| 23 | [`23_place_order`](./23_place_order) | `client.PlaceOrder` |
| 24 | [`24_order_list`](./24_order_list) | `client.GetOrderList` |
| 25 | [`25_cancel_order`](./25_cancel_order) | `client.ModifyOrder` (cancel) |
| 26 | [`26_history_order`](./26_history_order) | `client.GetHistoryOrderList` |
| 27 | [`27_order_fill`](./27_order_fill) | `client.GetOrderFillList` |
| 28 | [`28_history_fill`](./28_history_fill) | `client.GetHistoryOrderFillList` |
| 29 | [`29_acc_trading_info`](./29_acc_trading_info) | `client.GetAccTradingInfo` |

## Derivatives

| # | Example | SDK Function |
|---|---------|-------------|
| 42 | [`42_warrant`](./42_warrant) | `client.GetWarrant` |
| 43 | [`43_option_expiration`](./43_option_expiration) | `client.GetOptionExpirationDate` |
| 44 | [`44_option_chain`](./44_option_chain) | `client.GetOptionChain` |

## Alerts & User Data

| # | Example | SDK Function |
|---|---------|-------------|
| 31 | [`31_price_reminder`](./31_price_reminder) | `client.GetPriceReminder` |
| 45 | [`45_user_security`](./45_user_security) | `client.GetUserSecurity` |
| 46 | [`46_user_info`](./46_user_info) | `client.GetUserInfo` |
| 49 | [`49_subscribe_price_reminder`](./49_subscribe_price_reminder) | `chanpkg.SubscribePriceReminder` |

## SDK Functions Without Examples

| SDK Function | Reason |
|-------------|--------|
| `client.GetStaticInfo` | Not covered |
| `client.Unsubscribe` | Not covered |
| `client.UnsubscribeAll` | Not covered |
| `client.CancelAllOrder` | Not covered |
| `client.GetOrderFee` | Not covered |
| `client.GetMarginRatio` | Not covered |
| `client.SetPriceReminder` | Not covered |
| `client.GetDelayStatistics` | Known bug — proto2/proto3 wire-format mismatch |

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

## Troubleshooting

- **`connection refused`** — OpenD isn't running. Set `FUTU_ADDR=127.0.0.1:11111`
- **`subscribe first`** — `GetKLines`, `GetQuote`, `GetTicker`, `GetRT`, `GetBroker`, `GetOrderBook` require `client.Subscribe` first
- **`no positions`** — normal for simulator; real OpenD shows actual positions
