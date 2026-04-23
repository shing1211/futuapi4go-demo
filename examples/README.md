# Examples

> Copy-paste-ready code that actually works. Each example is a complete, runnable program demonstrating one SDK function.

## Quick Start

### With the Simulator (Recommended)

```powershell
# Terminal 1 — start the mock OpenD
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2 — run any example
go run ./examples/00_connect
go run ./examples/01_quote
go run ./examples/02_ticker
go run ./examples/03_orderbook
go run ./examples/04_rt
go run ./examples/05_broker
go run ./examples/06_kline_single
go run ./examples/07_kline_multi
```

### With a Real OpenD

```powershell
set FUTU_ADDR=127.0.0.1:11111
go run ./examples/00_connect
# ...
```

## Available Examples

| Example | SDK Function | What it Does |
|---------|-------------|-------------|
| [`00_connect`](./00_connect) | `client.Connect` | Connect to OpenD and disconnect |
| [`01_quote`](./01_quote) | `client.GetQuote` | Snapshot quote (price, open, high, low, volume) |
| [`02_ticker`](./02_ticker) | `chanpkg.SubscribeTicker` | Real-time tick trades |
| [`03_orderbook`](./03_orderbook) | `chanpkg.SubscribeOrderBook` | Order book (bids & asks) |
| [`04_rt`](./04_rt) | `chanpkg.SubscribeRT` | Tick-by-tick time & sales |
| [`05_broker`](./05_broker) | `chanpkg.SubscribeBroker` | Broker queue (bid/ask queues) |
| [`06_kline_single`](./06_kline_single) | `client.GetKLines` | Historical K-lines (one-shot) |
| [`07_kline_multi`](./07_kline_multi) | `chanpkg.SubscribeKLines` | Live K-lines for multiple periods |

## Example Descriptions

### 00_connect — `client.Connect`
The simplest possible program — connect to OpenD and exit cleanly. Verifies your OpenD address and network are working.

```go
cli := client.New()
cli.Connect("127.0.0.1:11111")
```

### 01_quote — `client.GetQuote`
Fetches a snapshot quote for NVDA: last price, open, high, low, volume. One-shot request, no subscription needed.

```go
quote, _ := client.GetQuote(ctx, cli, int32(constant.Market_US), "NVDA")
fmt.Println(quote.Price, quote.Open, quote.High, quote.Low, quote.Volume)
```

### 02_ticker — `chanpkg.SubscribeTicker`
Subscribes to real-time tick trades for NVDA. Each tick includes price, volume, and trade direction. Fires continuously while the market is open.

```go
chanpkg.SubscribeTicker(cli, int32(constant.Market_US), "NVDA", tickerCh)
```

### 03_orderbook — `chanpkg.SubscribeOrderBook`
Subscribes to the real-time order book (limit order book) for NVDA. Shows top 5 bid and ask price levels with volumes. Updates as orders are placed, modified, and cancelled.

```go
chanpkg.SubscribeOrderBook(cli, int32(constant.Market_US), "NVDA", orderBookCh)
```

### 04_rt — `chanpkg.SubscribeRT`
Subscribes to tick-by-tick time & sales data for NVDA. Each record is an individual trade: timestamp, price, volume, and average price. No aggregation.

```go
chanpkg.SubscribeRT(cli, int32(constant.Market_US), "NVDA", rtCh)
```

### 05_broker — `chanpkg.SubscribeBroker`
Subscribes to broker queue data for NVDA. Shows the top-of-book broker (market maker) names and their positions on the bid and ask sides.

```go
chanpkg.SubscribeBroker(cli, int32(constant.Market_US), "NVDA", brokerCh)
```

### 06_kline_single — `client.GetKLines`
Fetches historical K-line (candlestick) data for NVDA as a one-shot request. Gets the last 10 daily bars. Supports any `KLType`: `KLType_K_1Min`, `KLType_K_5Min`, `KLType_K_15Min`, `KLType_K_30Min`, `KLType_K_60Min`, `KLType_K_Day`, `KLType_K_Week`, `KLType_K_Month`, `KLType_K_Quarter`, `KLType_K_Year`.

```go
klines, _ := client.GetKLines(cli, int32(constant.Market_US), "NVDA", int32(constant.KLType_K_Day), 10)
for _, bar := range klines {
    fmt.Printf("%s  O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
        bar.Time, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
}
```

### 07_kline_multi — `chanpkg.SubscribeKLines`
Subscribes to live K-line updates for NVDA across multiple periods simultaneously — 1 min, 5 min, and daily. Uses a single call with a map of period-to-callback functions. Each callback fires independently when that period's bar updates.

```go
stop := chanpkg.SubscribeKLines(cli, int32(constant.Market_US), "NVDA", map[constant.KLType]func(*push.UpdateKL){
    KLType_K_1Min: func(kl *push.UpdateKL) { fmt.Println("1min bar!") },
    KLType_K_5Min: func(kl *push.UpdateKL) { fmt.Println("5min bar!") },
    KLType_K_Day:  func(kl *push.UpdateKL) { fmt.Println("day bar!") },
})
defer stop()
```

## Common Patterns

**Market constant** — all APIs take `int32` for market:
```go
int32(constant.Market_HK)  // 1 — Hong Kong
int32(constant.Market_US)  // 11 — United States
int32(constant.Market_SH)  // 21 — Shanghai A
int32(constant.Market_SZ)  // 22 — Shenzhen A
```

**Request vs Subscribe**
- **Request** (`client.GetQuote`, `client.GetKLines`): one-shot call, returns data immediately.
- **Subscribe** (`chanpkg.Subscribe*`): registers a channel; data flows continuously. Call the returned `stop` function to unsubscribe.

**Stock codes**: `"00700"`, `"NVDA"`, `"AAPL"`, `"9988.HK"`.

## Troubleshooting

**`connection refused`** — OpenD isn't running. Check `FUTU_ADDR`:
```powershell
set FUTU_ADDR=127.0.0.1:11111
```

**`Found 0 positions`** — normal for the simulator; real OpenD shows actual positions.

**`no such host`** — `FUTU_ADDR` has a trailing slash or spaces.
