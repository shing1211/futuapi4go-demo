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
go run ./examples/47_subscribe_quote
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

| # | Example | SDK Function | Category |
|---|---------|-------------|---------|
| 00 | `00_connect` | `client.Connect` | Connection |
| 01 | `01_quote` | `client.GetQuote` | Market Data |
| 02 | `02_ticker` | `chanpkg.SubscribeTicker` | Market Data |
| 03 | `03_orderbook` | `chanpkg.SubscribeOrderBook` | Market Data |
| 04 | `04_rt` | `chanpkg.SubscribeRT` | Market Data |
| 05 | `05_broker` | `chanpkg.SubscribeBroker` | Market Data |
| 06 | `06_kline_single` | `client.GetKLines` | Market Data |
| 07 | `07_kline_multi` | `chanpkg.SubscribeKLines` | Market Data |
| 08 | `08_orderbook_req` | `client.GetOrderBook` | Market Data |
| 09 | `09_ticker_req` | `client.GetTicker` | Market Data |
| 10 | `10_rt_req` | `client.GetRT` | Market Data |
| 11 | `11_broker_req` | `client.GetBroker` | Market Data |
| 12 | `12_capital_flow` | `client.GetCapitalFlow` | Market Analysis |
| 13 | `13_plate_set` | `client.GetPlateSet` | Market Analysis |
| 14 | `14_plate_stock` | `client.GetPlateSecurity` | Market Analysis |
| 15 | `15_history_kline` | `client.RequestHistoryKL` | Market Data |
| 16 | `16_market_state` | `client.GetMarketState` | Market Data |
| 17 | `17_snapshot` | `client.GetSecuritySnapshot` | Market Data |
| 18 | `18_global_state` | `client.GetGlobalState` | System |
| 19 | `19_account_list` | `client.GetAccountList` | Trading |
| 20 | `20_funds` | `client.GetFunds` | Trading |
| 21 | `21_positions` | `client.GetPositionList` | Trading |
| 22 | `22_unlock_trade` | `client.UnlockTrading` | Trading |
| 23 | `23_place_order` | `client.PlaceOrder` | Trading |
| 24 | `24_order_list` | `client.GetOrderList` | Trading |
| 25 | `25_cancel_order` | `client.ModifyOrder` | Trading |
| 26 | `26_history_order` | `client.GetHistoryOrderList` | Trading |
| 27 | `27_order_fill` | `client.GetOrderFillList` | Trading |
| 28 | `28_history_fill` | `client.GetHistoryOrderFillList` | Trading |
| 29 | `29_acc_trading_info` | `client.GetAccTradingInfo` | Trading |
| 30 | `30_trade_date` | `client.GetTradeDate` | Market Data |
| 31 | `31_price_reminder` | `client.GetPriceReminder` | Alerts |
| 32 | `32_owner_plate` | `client.GetOwnerPlate` | Market Analysis |
| 33 | `33_capital_distribution` | `client.GetCapitalDistribution` | Market Analysis |
| 34 | `34_stock_filter` | `client.StockFilter` | Market Analysis |
| 35 | `35_reference` | `client.GetReference` | Market Analysis |
| 36 | `36_ipo_list` | `client.GetIpoList` | Market Data |
| 37 | `37_future_info` | `client.GetFutureInfo` | Market Data |
| 38 | `38_suspend` | `client.GetSuspend` | Market Data |
| 39 | `39_holding_change` | `client.GetHoldingChangeList` | Market Analysis |
| 40 | `40_rehab` | `client.RequestRehab` | Market Data |
| 41 | `41_code_change` | `client.GetCodeChange` | Market Data |
| 42 | `42_warrant` | `client.GetWarrant` | Derivatives |
| 43 | `43_option_expiration` | `client.GetOptionExpirationDate` | Derivatives |
| 44 | `44_option_chain` | `client.GetOptionChain` | Derivatives |
| 45 | `45_user_security` | `client.GetUserSecurity` | User Data |
| 46 | `46_user_info` | `client.GetUserInfo` | User Data |
| 47 | `47_subscribe_quote` | `chanpkg.SubscribeQuote` | Market Data |
| 48 | `48_subscribe_kline_single` | `chanpkg.SubscribeKLine` | Market Data |
| 49 | `49_subscribe_price_reminder` | `chanpkg.SubscribePriceReminder` | Alerts |

## Project Layout

```
futuapi4go-demo/
├── examples/              # 50 standalone examples (00-49)
│   └── README.md          # Full example descriptions
├── docs/
│   └── FUTU_PROTO_REF.md
├── AGENTS.md
├── README.md
└── LICENSE
```

## SDK Reference

```go
import "github.com/shing1211/futuapi4go/pkg/constant"

int32(constant.Market_US)   // 11
int32(constant.Market_HK)   // 1
int32(constant.KLType_K_Day)
int32(constant.SubType_Quote)
int32(constant.TrdSide_Buy)
int32(constant.OrderType_Normal)
```

## Known Caveats

- **`GetDelayStatistics`** — skipped due to a proto2/proto3 wire-format mismatch between Go's protobuf library and OpenD's C++ parser. All other APIs work normally.
- **`GetTradeDate`** — may return an error on older OpenD versions if all required C2S fields aren't populated.
- **`subscribe first`** — `GetKLines` and `GetQuote` require `client.Subscribe` before calling.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
