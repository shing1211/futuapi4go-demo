# futuapi4go New Features Examples

This directory contains examples demonstrating new SDK features implemented in recent releases.

## WebSocket Connection (v0.3.1+)

Connect to Futu OpenD via WebSocket for better push handling:

```bash
go run ./examples/00_ws_connect
```

### Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	cli := client.New()
	defer cli.Close()

	// TCP connection (default)
	if err := cli.Connect("127.0.0.1:11111"); err != nil {
		log.Fatal(err)
	}

	// Or WebSocket connection
	// if err := cli.ConnectWS("127.0.0.1:11113"); err != nil {
	//     log.Fatal(err)
	// }

	// Get quote via WebSocket
	quote, err := client.GetQuote(context.Background(), cli, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("NVDA: %v\n", quote)
}
```

## Zero-Allocation Path (v0.3.1+)

The SDK now uses `sync.Pool` internally for request/response buffers - no user-facing API needed.

Performance improvement: ~10-30% reduction in GC pressure for high-frequency trading.

## Connection Pool O(1) Lookup (v0.3.1+)

Pool operations are now O(1) instead of O(n).

```go
pool := client.NewClientPool(config)
client, _ := pool.Get(ctx, client.PoolTypeMarketData)
// ... use client
pool.Put(client)
```

## Structured Logging (v0.3.1+)

```go
cli := client.New(
    client.WithSlogLevel(constant.LevelInfo),
    // Or with custom writer:
    // client.WithSlogWriter(os.Stderr, constant.LevelDebug),
)
```

## Historical Data Downloader (v0.3.1+)

```go
import "github.com/shing1211/futuapi4go/pkg/history"

downloader := history.NewDownloader(cli)
err := downloader.DownloadKLine(ctx, history.KLineRequest{
    Code:       "00700",
    Market:     constant.Market_HK,
    KLType:     constant.KLType_K_1Day,
    StartDate:  "2020-01-01",
    EndDate:    "2024-01-01",
})
```

## Market Hours Utility

```go
import "github.com/shing1211/futuapi4go/pkg/market"

isOpen := market.IsHKOpen(time.Now())
untilClose := market.UntilClose(time.Now())
```

## Option Chain Utilities

```go
import "github.com/shing1211/futuapi4go/pkg/option"

owner, expiry, strike, optType := option.ParseOptionCode("NVDA2406C200")
atmCalls := option.FindAtm(cli, "NVDA", expiry, constant.OptionType_Call)
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `FUTU_ADDR` | OpenD TCP address | `127.0.0.1:11111` |
| `FUTU_WS_ADDR` | OpenD WebSocket address | `127.0.0.1:11113` |
| `FUTU_TRADE_PWD` | Trading password MD5 | (not set) |