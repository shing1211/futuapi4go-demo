# futuapi4go-demo Architecture

## Overview

`futuapi4go-demo` is a repository of **96 standalone Go examples** demonstrating the [futuapi4go](https://github.com/shing1211/futuapi4go) SDK for connecting to Futu's OpenD trading gateway. Each example is a self-contained `main.go` that exercises one or more SDK functions.

The repository has three layers:

1. **Examples (00–96)** — the visible surface; each is a runnable program
2. **Shared helpers (`examples/pkg/`)** — reusable connection management
3. **SDK dependency** — `github.com/shing1211/futuapi4go` (at `v0.6.2`)

## Functional Areas

### Area 1 — Connection (`examples/00_*`)

| Example | Protocol | Description |
|---------|----------|-------------|
| `00_connect` | TCP plain | Basic `client.Connect` to local OpenD |
| `00_rsa_connect` | TCP + RSA | RSA-encrypted connection for remote OpenD |
| `00_ws_connect` | WebSocket | `client.ConnectWS` with secret key auth |

**Shared infrastructure**: `examples/pkg/connect/connect.go` provides HA connection management including:
- Parallel TCP probe of multiple hosts, sorted by latency
- Per-host RSA configuration
- Auto-reconnect with exponential backoff
- State machine: `Disconnected → Connecting → Connected → Reconnecting → Failed`
- Keep-alive monitoring every 30s

### Area 2 — Market Data — One-Shot (examples 01–17)

One-time request/response calls. No persistent subscription.

| # | Example | SDK Function |
|---|---------|-------------|
| 01 | `01_quote` | `client.GetQuote` |
| 08 | `08_orderbook_req` | `client.GetOrderBook` |
| 09 | `09_ticker_req` | `client.GetTicker` |
| 10 | `10_rt_req` | `client.GetRT` |
| 11 | `11_broker_req` | `client.GetBroker` |
| 06 | `06_kline_single` | `client.GetKLines` |
| 15 | `15_history_kline` | `client.RequestHistoryKL` |
| 12 | `12_capital_flow` | `client.GetCapitalFlow` |
| 13 | `13_plate_set` | `client.GetPlateSet` |
| 14 | `14_plate_stock` | `client.GetPlateSecurity` |
| 16 | `16_market_state` | `client.GetMarketState` |
| 17 | `17_global_state` | `client.GetGlobalState` |
| 24 | `24_snapshot` | `client.GetSecuritySnapshot` |
| 25 | `25_trade_date` | `client.GetTradeDate` |
| 28 | `28_owner_plate` | `client.GetOwnerPlate` |
| 29 | `29_capital_distribution` | `client.GetCapitalDistribution` |
| 30 | `30_stock_filter` | `client.StockFilter` |
| 32 | `32_future_info` | `client.GetFutureInfo` |
| 33 | `33_suspend` | `client.GetSuspend` |
| 35 | `35_rehab` | `client.RequestRehab` |
| 36 | `36_code_change` | `client.GetCodeChange` |
| 37 | `37_warrant` | `client.GetWarrant` |
| 38 | `38_option_chain` | `client.GetOptionChain` |
| 39 | `39_option_expiration` | `client.GetOptionExpirationDate` |
| 40 | `40_reference` | `client.GetReference` |
| 59 | `59_static_info` | `client.GetStaticInfo` |

**Pattern**: Subscribe (optional) → One-shot call → Result.

### Area 3 — Market Data — Streaming (examples 02–05, 47–49)

Push-based subscriptions with callback channels from `chanpkg`.

| # | Example | Channel Function |
|---|---------|-----------------|
| 02 | `02_ticker` | `chanpkg.SubscribeTicker` |
| 03 | `03_orderbook` | `chanpkg.SubscribeOrderBook` |
| 04 | `04_rt` | `chanpkg.SubscribeRT` |
| 05 | `05_broker` | `chanpkg.SubscribeBroker` |
| 47 | `47_subscribe_quote` | `chanpkg.SubscribeQuote` |
| 48 | `48_subscribe_kline_single` | `chanpkg.SubscribeKLine` |
| 49 | `49_subscribe_price_reminder` | `chanpkg.SubscribePriceReminder` |
| 07 | `07_kline_multi` | `chanpkg.SubscribeKLines` |

**Pattern**: `SubscribeXXX` returns a `stop()` function; call `defer stop()` to clean up.

### Area 4 — Trading & Accounts (examples 18–27, 42–46, 54–58, 70–75)

Trading operations require unlocked account with `FUTU_TRADE_PWD` set.

| Category | Examples | SDK Functions |
|----------|----------|---------------|
| Account listing | 18, 19, 70, 73 | `GetAccountList`, `GetAccList(TrdCategory_Future)` |
| Funds & positions | 20, 21, 71–75 | `GetFunds`, `GetPositionList` |
| Order placement | 22, 67 | `PlaceOrder` |
| Order management | 23, 27, 54 | `GetOrderList`, `ModifyOrder`, `CancelAllOrder` |
| Order fills | 43, 44 | `GetOrderFillList`, `GetHistoryOrderFillList` |
| Historical orders | 42 | `GetHistoryOrderList` |
| Trading info | 45, 46, 55–58 | `GetAccTradingInfo`, `GetUserInfo`, `GetMaxTrdQtys`, `GetOrderFee`, `GetMarginRatio`, `GetFlowSummary` |

**Key constraint**: Simulate trading (default) does not support fills, flow summary, or real order fees. Use `WithTradeEnv(constant.TrdEnv_Real)` for real trading.

### Area 5 — Subscription Management (examples 50–53, 61, 63)

| # | Example | SDK Function |
|---|---------|-------------|
| 50 | `50_unsubscribe` | `client.Unsubscribe` |
| 51 | `51_unsubscribe_all` | `client.UnsubscribeAll` |
| 52 | `52_query_subscription` | `client.QuerySubscription` |
| 53 | `53_reg_qot_push` | `client.RegQotPush` |
| 61 | `61_sub_info` | `client.GetSubInfo` |
| 63 | `63_sub_acc_push` | `client.SubAccPush` |

### Area 6 — Price Alerts (examples 26, 62)

| # | Example | SDK Function |
|---|---------|-------------|
| 26 | `26_price_reminder` | `client.GetPriceReminder` |
| 62 | `62_set_price_reminder` | `client.SetPriceReminder` |

### Area 7 — Watchlist (examples 41, 60)

| # | Example | SDK Function |
|---|---------|-------------|
| 41 | `41_user_security` | `GetUserSecurityGroup`, `GetUserSecurity` |
| 60 | `60_modify_user_security` | `client.ModifyUserSecurity` |

### Area 8 — Advanced Combo (examples 66–80)

Multi-function workflows that combine several SDK calls.

| # | Example | Composition |
|---|---------|-------------|
| 66 | `66_multi_symbol_kline` | Subscribe + `GetKLines` + `RequestHistoryKL` across multiple symbols |
| 67 | `67_order_lifecycle` | `PlaceOrder` → polling `GetOrderList` → `ModifyOrder` (cancel) |
| 68 | `68_market_hours_check` | `GetMarketState` + `GetTradeDate` to determine if market is open |
| 69 | `69_subscribe_handler` | `Subscribe` + push handlers for Ticker, KLine, OrderBook simultaneously |
| 76 | `76_pre_trade_checks` | Market state + account funds + position validation before placing order |
| 77 | `77_realtime_dashboard` | Ticker subscription + periodic snapshot refresh |
| 78 | `78_dca_grid_bot` | DCA buy orders + grid sell orders on price bands |
| 79 | `79_momentum_scanner` | `StockFilter` → `GetSecuritySnapshot` → `GetKLines` for momentum scoring |
| 80 | `80_vwap_executor` | `GetOrderBook` → VWAP calculation → split order execution |

### Area 9 — Quantitative Strategies (examples 91–95)

Multi-source signal fusion and automated trading logic.

| # | Example | SDK Functions |
|---|---------|---------------|
| 91 | `91_orderbook_imbalance` | `GetOrderBook` + `SubscribeOrderBook` — bid/ask imbalance scoring, iceberg detection |
| 92 | `92_pairs_trading` | `GetKLines` + Pearson correlation + spread z-score — stat arb pairs |
| 93 | `93_smart_money` | `GetCapitalFlow` + `GetCapitalDistribution` + `GetBroker` + `GetOrderBook` — institutional flow |
| 94 | `94_portfolio_rebalance` | `GetPositionList` + `GetFunds` + `GetMaxTrdQtys` + `PlaceOrder` — multi-asset rebalancer |
| 95 | `95_earnings_vol_strategy` | `GetOptionChain` + `GetKLines` — earnings straddle, implied vs historical vol |

## Key Execution Flows

### Flow 1 — Streaming Quote (Ticker)

```
main
  └─ connect.MustConnect         → ManagedConnection (HA)
  └─ chanpkg.SubscribeTicker     → TickerChannel (blocking read loop)
       └─ for ticker := range ch
            fmt.Printf("[ticker] %s %d @ %.2f\n", ...)
```

### Flow 2 — Place and Monitor Order

```
main
  └─ connect.MustConnect
  └─ client.GetAccountList       → find default account
  └─ client.PlaceOrder           → get OrderID
  └─ client.GetOrderList         → poll until filled/cancelled
  └─ client.ModifyOrder          → cancel if needed
```

### Flow 3 — Historical K-Line Retrieval

```
main
  └─ connect.MustConnect
  └─ client.Subscribe             → pre-subscribe K-line type
  └─ client.RequestHistoryKL      → paginated historical bars
  └─ client.GetKLines             → live bars (after subscribe)
```

### Flow 4 — VWAP Execution (example 80)

```
main
  └─ connect.MustConnect
  └─ client.GetOrderBook          → bid/ask levels
  └─ calculate VWAP from book
  └─ client.GetAccountList        → find account
  └─ client.PlaceOrder            → execute slices at VWAP
```

### Flow 5 — DCA + Grid Bot (example 78)

```
main
  └─ connect.MustConnect
  └─ client.GetPositionList      → check existing position
  └─ client.GetFunds             → check buying power
  └─ timer/price loop:
       ├─ if price crosses lower grid → buy DCA increment
       ├─ if price crosses upper grid → sell grid unit
       └─ update position and funds
```

## Architecture Diagram

```mermaid
flowchart TB
    subgraph examples["examples/ — 96 standalone programs"]
        direction TB
        subgraph connection["Connection (00_*)\n3 examples"]
            c0[00_connect\nplain TCP]
            crsa[00_rsa_connect\nTCP + RSA]
            cws[00_ws_connect\nWebSocket]
        end
        subgraph mkt_one["Market Data — One-Shot (01-40)\n~40 examples"]
            q1[01_quote\nGetQuote]
            kl[06 15 66\nKLines / History]
            ob[08_orderbook_req\nGetOrderBook]
            sf[30_stock_filter\nStockFilter]
        end
        subgraph mkt_stream["Market Data — Streaming (02-05, 07, 47-49)\n8 examples"]
            st[02_ticker\nSubscribeTicker]
            sb[03_orderbook\nSubscribeOrderBook]
            sk[07_kline_multi\nSubscribeKLines]
        end
        subgraph trade["Trading & Accounts (18-27, 42-46, 54-58, 70-75)\n~20 examples"]
            po[22_place_order\nPlaceOrder]
            pl[23_order_list\nGetOrderList]
            pos[21_positions\nGetPositionList]
            fu[20_funds\nGetFunds]
        end
        subgraph advanced["Advanced Combo (66-80)\n15 examples"]
            vwap[80_vwap_executor\nOrderBook + VWAP]
            dca[78_dca_grid_bot\nDCA + Grid]
            scan[79_momentum_scanner\nFilter + Snapshot + KL]
        end
        subgraph quant["Quant Strategies (91-95)\n5 examples"]
            imb[91_orderbook_imbalance\nbid/ask imbalance]
            pairs[92_pairs_trading\nstat arb pairs]
            smart[93_smart_money\ninstitutional flow]
        end
    end

    subgraph pkg["examples/pkg/"]
        connect_pkg["connect/\nMustConnect, ManagedConnection\nHA: probe → connect → keepalive → reconnect"]
    end

    subgraph sdk["github.com/shing1211/futuapi4go v0.6.2"]
        client_pkg["client/\nConnect, GetQuote, PlaceOrder, ..."]
        chanpkg["chanpkg/\nSubscribeTicker, SubscribeKLine, ..."]
        push["push/\nUpdateTicker, UpdateKL, ..."]
        constant["constant/\nMarket_US, TrdEnv_Real, ..."]
        pb["pb/\nProtobuf messages"]
    end

    subgraph opend["Futu OpenD\nlocalhost:11111 (default)"]
        tcp[TCP Server]
        ws[WebSocket Server]
        proto[Protobuf Protocol\n(Qot/Trd APIs)"]
    end

    c0 --> connect_pkg
    crsa --> connect_pkg
    cws --> client_pkg
    connect_pkg --> client_pkg
    mkt_one --> client_pkg
    mkt_stream --> chanpkg
    trade --> client_pkg
    advanced --> client_pkg
    advanced --> chanpkg
    quant --> client_pkg
    quant --> chanpkg

    client_pkg --> pb
    chanpkg --> push
    pb --> proto
    push --> proto

    client_pkg --> tcp
    client_pkg --> ws
    connect_pkg -.-> "RSA encryption\n(remote hosts)" .-> tcp
    cws --> ws

    style connect_pkg fill:#e1f5fe
    style sdk fill:#fff3e0
    style opend fill:#f3e5f5
```

## Environment & Configuration

| Environment | Default | Purpose |
|-------------|---------|---------|
| `FUTU_ADDR` | `127.0.0.1:11111` | OpenD TCP address |
| `FUTU_WS_ADDR` | `127.0.0.1:11113` | OpenD WebSocket address |
| `FUTU_WS_SECRET` | *(not set)* | WebSocket auth secret |
| `FUTU_TRADE_PWD` | *(not set)* | MD5 hex of trading password (for real trading) |
| `FUTU_RSA_KEY` | `/etc/futu/keys/private_key.pem` | RSA private key for encrypted connections |
| `FUTU_OPEND_HOSTS` | *(single host)* | Comma-separated list for HA: `host:port:rsa` |

The `examples/pkg/connect` package reads all environment variables via `os.Getenv` / `godotenv.Load()`, so no code changes are needed when switching between local/simulate/real environments.

## SDK Reference

The SDK is imported as `github.com/shing1211/futuapi4go` and provides three main packages used across all examples:

- **`client`** — synchronous request/response calls (`GetQuote`, `PlaceOrder`, `GetAccountList`, etc.)
- **`chanpkg`** — channel-based push subscriptions (`SubscribeTicker`, `SubscribeKLine`, etc.)
- **`constant`** — typed market/trade constants (`Market_US = 11`, `TrdEnv_Real = 1`, etc.)

All examples use `context.Background()` and `defer mc.Close()` for clean lifecycle management.

## Known Limitations

| Issue | Workaround |
|-------|------------|
| `GetDelayStatistics` — proto2/proto3 wire format mismatch | Skipped in all examples |
| Simulate trading — no fill data | Use `WithTradeEnv(TrdEnv_Real)` + `FUTU_TRADE_PWD` |
| US stocks — need `Subscribe` before `GetQuote` returns data | HK stocks do not require pre-subscribe |
| `ReconfirmOrder` — not implemented by OpenD | Skipped in example 64 |