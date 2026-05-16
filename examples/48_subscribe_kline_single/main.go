// 48_subscribe_kline_single demonstrates SubscribeKLine for receiving
// real-time K-line (OHLCV bar) updates via a Go channel.
//
// This example subscribes to 5-minute K-lines for NVDA and displays
// each bar update as it arrives from the server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ch := make(chan *push.UpdateKL, 100)
	stop, err := chanpkg.SubscribeKLine(context.Background(), mc.Client, constant.Market_US, "NVDA", constant.KLType_K_5Min, ch)
	if err != nil {
		fmt.Printf("SubscribeKLine: %v\n", err)
		return
	}
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening for NVDA 5-minute K-lines (Ctrl+C to exit)...")
	fmt.Println("K-lines push every 5 minutes during market hours.")
	fmt.Println()

	// Wait for K-line data or timeout
	timeout := time.After(10 * time.Second)
	received := false

	for {
		select {
		case kl := <-ch:
			received = true
			fmt.Printf("📊 KLINE [5m]: %s\n", time.Now().Format("15:04:05"))
			for _, bar := range kl.KLList {
				if bar == nil {
					continue
				}
				if bar.Time != nil {
					fmt.Printf("   Time:   %s\n", *bar.Time)
				}
				if bar.OpenPrice != nil {
					fmt.Printf("   O: %.2f  H: %.2f  L: %.2f  C: %.2f\n",
						*bar.OpenPrice, *bar.HighPrice, *bar.LowPrice, *bar.ClosePrice)
				}
				if bar.Volume != nil {
					fmt.Printf("   Vol:    %d\n", *bar.Volume)
				}
			}
			fmt.Println()

			// Got data, exit demo successfully
			fmt.Println("✓ Received K-line data successfully. Demo complete.")
			return

		case <-timeout:
			if !received {
				fmt.Println("⏱️  No K-line data received within 10 seconds.")
				fmt.Println("   This is normal if market is closed or US stocks are in pre-market.")
				fmt.Println("   K-lines only update during active trading hours.")
				fmt.Println("   Try again during US market hours (9:30-16:00 ET).")
			}
			return

		case <-sig:
			fmt.Println("\nInterrupted. Cleaning up...")
			return
		}
	}
}