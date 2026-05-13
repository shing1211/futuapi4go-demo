package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Multi-Symbol K-Line Query (US Stocks) ===")
	fmt.Println()

	symbols := []string{"AAPL", "TSLA", "NVDA", "MSFT", "AMZN"}

	fmt.Println("--- Subscribing to K-line data ---")
	for _, symbol := range symbols {
		if err := client.Subscribe(ctx, mc.Client, constant.Market_US, symbol,
			[]constant.SubType{constant.SubType_K_Day}); err != nil {
			fmt.Printf("Subscribe %s failed: %v\n", symbol, err)
		} else {
			fmt.Printf("Subscribed to %s daily K-lines\n", symbol)
		}
	}

	fmt.Println("\n--- GetKLines (after subscription) ---")
	for _, symbol := range symbols {
		fmt.Printf("\n--- %s ---\n", symbol)

		klines, err := client.GetKLines(ctx, mc.Client, constant.Market_US, symbol,
			constant.KLType_K_Day, 5)
		if err != nil {
			fmt.Printf("  GetKLines failed: %v\n", err)
			continue
		}

		for _, k := range klines {
			fmt.Printf("  %s: O=%.2f H=%.2f L=%.2f C=%.2f Vol=%d\n",
				k.Time, k.Open, k.High, k.Low, k.Close, k.Volume)
		}
	}

	fmt.Println("\n=== Request Historical K-Lines (Batch) ===")
	fmt.Println()

	startTime := "2026-01-01"
	endTime := time.Now().Format("2006-01-02")

	for _, symbol := range symbols[:3] {
		fmt.Printf("Requesting historical K-lines for %s (%s to %s)...\n", symbol, startTime, endTime)

		klines, err := client.RequestHistoryKL(ctx, mc.Client, constant.Market_US, symbol,
			constant.KLType_K_Day, startTime, endTime)
		if err != nil {
			fmt.Printf("  RequestHistoryKL failed: %v\n", err)
			continue
		}

		fmt.Printf("  Got %d K-lines\n", len(klines))
		if len(klines) > 0 {
			latest := klines[len(klines)-1]
			fmt.Printf("  Latest: %s O=%.2f C=%.2f\n", latest.Time, latest.Open, latest.Close)
		}
	}

	fmt.Println("\n=== Note ===")
	fmt.Println("RequestHistoryKLWithLimit for paginated historical data")
	fmt.Println("GetKLines for recent/realtime K-line data")
}