# FutuAPI4Go Examples

This directory contains examples demonstrating how to use the futuapi4go SDK.

## Quick Start

### Using with Simulator (Recommended for Development)

1. **Start the simulator** (from the SDK directory):
```bash
cd D:\github\futuapi4go
go run .\cmd\examples\simulator
```

2. **Run the demo** (from the demo directory):
```bash
cd futuapi4go-demo
go run .\cmd\demo
```

Or run individual examples:
```bash
go run examples/getting_started
go run examples/trading_demo
```

### Using with Real Futu OpenD

1. **Ensure Futu OpenD is running** on `127.0.0.1:11111`

2. **Set environment variable** (optional):
```bash
set FUTU_ADDR=127.0.0.1:11111
```

3. **Run the demo**:
```bash
go run .\cmd\demo
```

## Available Examples

| Example | Description |
|---------|-------------|
| `cmd/demo` | Interactive menu demo — covers all major APIs |
| `getting_started` | Basic usage: connect, query quotes, fetch K-lines, subscribe |
| `trading_demo` | Trading operations: accounts, positions, orders, place/cancel |

## Example Structure

Each example follows a consistent pattern:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/shing1211/futuapi4go/client"
    "github.com/shing1211/futuapi4go/pkg/constant"
    "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
    "github.com/shing1211/futuapi4go/pkg/qot"
)

func main() {
    cli := client.New()
    defer cli.Close()

    addr := os.Getenv("FUTU_ADDR")
    if addr == "" {
        addr = "127.0.0.1:11111"
    }

    if err := cli.Connect(addr); err != nil {
        log.Fatalf("Connection failed: %v", err)
    }

    // Use public client API helpers or SDK functions with cli.Inner()
    hkMarket := int32(constant.Market_HK)
    sec := &qotcommon.Security{Market: &hkMarket, Code: ptrStr("00700")}

    quotes, err := qot.GetBasicQot(context.Background(), cli.Inner(), []*qotcommon.Security{sec})
    if err != nil {
        log.Fatalf("GetBasicQot failed: %v", err)
    }
    // ...
}
```

## Public Client API

The `github.com/shing1211/futuapi4go/client` package provides a high-level API:

```go
// Connect
cli := client.New()
cli.Connect("127.0.0.1:11111")

// GetQuote
quote, _ := client.GetQuote(ctx, cli, constant.Market_HK, "00700")

// GetKLines
klines, _ := client.GetKLines(cli, constant.Market_HK, "00700", constant.KLType_Day, 100)

// Subscribe
client.Subscribe(cli, constant.Market_HK, "00700", []int32{constant.SubType_Basic})

// GetAccountList
accounts, _ := client.GetAccountList(cli)

// PlaceOrder
result, _ := client.PlaceOrder(cli, accID, constant.Market_HK, "00700",
    constant.TrdSide_Buy, constant.OrderType_Normal, 350.0, 100.0)
```

## Helper Functions

All examples include helper functions for creating pointers:

```go
func ptrStr(s string) *string { return &s }
func ptrInt32(v int32) *int32 { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrBool(v bool) *bool { return &v }
```

## Notes

1. **Simulator vs Real OpenD**:
   - Simulator returns mock data for testing
   - Real OpenD provides live market data
   - Examples work with both!

2. **Trading Safety**:
   - Always test with simulator first
   - The public `client.New()` defaults to simulate trading environment
   - Use `cli.WithTradeEnv(constant.TrdEnv_Real)` for live trading

3. **API Coverage**:
   - All documented APIs are implemented
   - Some advanced APIs may return empty results in simulator

## Troubleshooting

### Connection Failed
```
Connection failed: dial: connection refused
```
**Solution**: Ensure Futu OpenD or simulator is running on the specified address.

### API Returns Empty Results
```
Found 0 positions
```
**Solution**: Normal for simulator; try with real OpenD for actual data.

---

**Happy Coding!**
