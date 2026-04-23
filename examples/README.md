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

    // 2. Call APIs — US stocks require subscription first
    if err := client.Subscribe(cli, constant.Market_US, "NVDA",
        []int32{int32(constant.SubType_Quote), int32(constant.SubType_K_1Min)}); err != nil {
        log.Fatalf("Subscribe failed: %v", err)
    }

    quote, err := client.GetQuote(context.Background(), cli, constant.Market_US, "NVDA")
    if err != nil {
        log.Fatalf("GetQuote failed: %v", err)
    }
    fmt.Printf("US.NVDA: price=%.2f\n", quote.Price)

    // 3. Subscribe to real-time K-line updates via channel
    klCh := make(chan *push.UpdateKL, 100)
    stop := chanpkg.SubscribeKLine(cli, constant.Market_US, "NVDA", constant.KLType_K_1Min, klCh)
    defer stop()

    for kl := range klCh {
        for _, bar := range kl.KLList {
            fmt.Printf("KL: %s O=%.2f H=%.2f L=%.2f C=%.2f\n",
                *bar.Time, *bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice)
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
client.Subscribe(cli, constant.Market_HK, "00700", []int32{constant.SubType_Quote})

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
