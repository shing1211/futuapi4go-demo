# Examples

> Copy-paste-ready code that actually works. Each example is a complete, runnable program.

## Quick Start

### With the Simulator (Recommended)

```powershell
# Terminal 1 — start the mock OpenD
go run github.com/shing1211/futuapi4go/cmd/examples/simulator

# Terminal 2 — run any example
go run ./examples/getting_started
go run ./examples/trading_demo
```

The simulator fires back realistic mock data so you can run the full stack without a Futu account or OpenD installed.

### With a Real OpenD

```powershell
# Make sure Futu OpenD is running (default: 127.0.0.1:11111)
set FUTU_ADDR=127.0.0.1:11111

go run ./examples/getting_started
go run ./examples/trading_demo
```

### Interactive Menu Demo

```powershell
go run ./cmd/demo/main.go
```

## Available Examples

| Example | What it Does |
|---------|-------------|
| `cmd/demo` | Interactive menu with all 10 demo categories |
| `examples/getting_started` | Connect → quote → K-line → subscribe |
| `examples/trading_demo` | Accounts → positions → orders → fills |

## Example Structure

Every example follows the same three-step pattern:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/shing1211/futuapi4go/client"
    "github.com/shing1211/futuapi4go/pkg/constant"
    "github.com/shing1211/futuapi4go/pkg/push"
    chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
)

func main() {
    // 1. Create & connect
    cli := client.New()
    defer cli.Close()

    addr := os.Getenv("FUTU_ADDR")
    if addr == "" {
        addr = "127.0.0.1:11111"
    }

    if err := cli.Connect(addr); err != nil {
        log.Fatalf("Connection failed: %v", err)
    }

    // 2. Subscribe to ALL data types for NVDA
    allSubTypes := []constant.SubType{
        constant.SubType_Quote,
        constant.SubType_OrderBook,
        constant.SubType_Ticker,
        constant.SubType_RT,
        constant.SubType_Broker,
        constant.SubType_K_1Min,
        constant.SubType_K_5Min,
        constant.SubType_K_15Min,
        constant.SubType_K_30Min,
        constant.SubType_K_60Min,
        constant.SubType_K_Day,
        constant.SubType_K_Week,
        constant.SubType_K_Month,
    }
    if err := client.Subscribe(cli, constant.Market_US, "NVDA", allSubTypes); err != nil {
        log.Fatalf("Subscribe failed: %v", err)
    }

    quote, err := client.GetQuote(context.Background(), cli, constant.Market_US, "NVDA")
    if err != nil {
        log.Fatalf("GetQuote failed: %v", err)
    }
    fmt.Printf("US.NVDA: price=%.2f open=%.2f high=%.2f low=%.2f vol=%d\n",
        quote.Price, quote.Open, quote.High, quote.Low, quote.Volume)

    // 3. Set up channel listeners for each data type
    quoteCh     := make(chan *push.UpdateBasicQot, 100)
    tickerCh    := make(chan *push.UpdateTicker, 100)
    orderBookCh := make(chan *push.UpdateOrderBook, 100)
    rtCh        := make(chan *push.UpdateRT, 100)
    brokerCh    := make(chan *push.UpdateBroker, 100)
    klCh        := make(chan *push.UpdateKL, 100)

    chanpkg.SubscribeQuote(cli, constant.Market_US, "NVDA", quoteCh)
    chanpkg.SubscribeTicker(cli, constant.Market_US, "NVDA", tickerCh)
    chanpkg.SubscribeOrderBook(cli, constant.Market_US, "NVDA", orderBookCh)
    chanpkg.SubscribeRT(cli, constant.Market_US, "NVDA", rtCh)
    chanpkg.SubscribeBroker(cli, constant.Market_US, "NVDA", brokerCh)
    chanpkg.SubscribeKLine(cli, constant.Market_US, "NVDA", constant.KLType_K_1Min, klCh)

    for {
        select {
        case q := <-quoteCh:
            fmt.Printf("QUOTE [%s]: price=%.2f vol=%d\n",
                q.Security.GetCode(), q.CurPrice, q.Volume)
        case t := <-tickerCh:
            if len(t.TickerList) > 0 {
                fmt.Printf("TICKER: price=%.2f vol=%d\n",
                    t.TickerList[0].GetPrice(), t.TickerList[0].GetVolume())
            }
        case ob := <-orderBookCh:
            if len(ob.OrderBookBidList) > 0 && len(ob.OrderBookAskList) > 0 {
                fmt.Printf("ORDERBOOK: bid=%.2f ask=%.2f\n",
                    ob.OrderBookBidList[0].GetPrice(), ob.OrderBookAskList[0].GetPrice())
            }
        case rt := <-rtCh:
            if len(rt.RTList) > 0 {
                fmt.Printf("RT: price=%.2f avg=%.2f\n",
                    rt.RTList[0].GetPrice(), rt.RTList[0].GetAvgPrice())
            }
        case b := <-brokerCh:
            if len(b.BidBrokerList) > 0 {
                fmt.Printf("BROKER: name=%s pos=%d\n",
                    b.BidBrokerList[0].GetName(), b.BidBrokerList[0].GetPos())
            }
        case kl := <-klCh:
            for _, bar := range kl.KLList {
                fmt.Printf("KL: %s O=%.2f H=%.2f L=%.2f C=%.2f V=%d\n",
                    *bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice, *bar.Volume)
            }
        }
    }
}
```

## High-Level Client API

The `client` package wraps the raw SDK into friendly functions:

```go
// Connect
cli := client.New()
cli.Connect("127.0.0.1:11111")

// Get a quote
quote, _ := client.GetQuote(ctx, cli, constant.Market_HK, "00700")

// Fetch K-lines
klines, _ := client.GetKLines(cli, constant.Market_HK, "00700", constant.KLType_K_Day, 100)

// Subscribe to real-time data
client.Subscribe(cli, constant.Market_HK, "00700", []constant.SubType{constant.SubType_Quote})

// List accounts
accounts, _ := client.GetAccountList(cli)

// Place a simulated order
result, _ := client.PlaceOrder(cli, accID, constant.Market_HK, "00700",
    constant.TrdSide_Buy, constant.OrderType_Normal, 350.0, 100)
```

## Helper Utilities

Every example uses these pointer helpers to keep proto struct literals clean:

```go
func ptrStr(s string) *string    { return &s }
func ptrInt32(v int32) *int32   { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrBool(v bool) *bool       { return &v }
```

## Simulator vs Real OpenD

| | Simulator | Real OpenD |
|---|---|---|
| Quote data | Realistic mock | Live market |
| Trading | Simulated fills | Real fills |
| API latency | Instant | Network-dependent |
| Account needed | No | Yes (logged in) |

Both work with the exact same code — swap the address and you're done.

## Trading Safety

- **Always test with the simulator first.** Switch to `constant.TrdEnv_Real` only when you're ready.
- The client defaults to **simulate mode** (`constant.TrdEnv_Simulate`) out of the box.
- If you accidentally run a live order in simulate mode, nothing actually trades.

## Troubleshooting

**`connection refused`**

Futu OpenD (or the simulator) isn't running on the specified address. Check the port:

```powershell
set FUTU_ADDR=127.0.0.1:11111
```

**`Found 0 positions`**

Normal for the simulator — it returns an empty account. With a real OpenD you'll see your actual positions.

**`no such host`**

Make sure your `FUTU_ADDR` doesn't have a trailing slash or spaces.

---

**Happy exploring!**
