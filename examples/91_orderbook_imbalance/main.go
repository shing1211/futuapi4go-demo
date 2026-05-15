package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/push"
	chanpkg "github.com/shing1211/futuapi4go/pkg/push/chan"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Order Book Imbalance & Liquidity Pressure Analyzer ===")
	fmt.Println("Real-time bid/ask imbalance, iceberg detection, liquidity scoring")
	fmt.Println()

	symbol := "NVDA"

	fmt.Printf("--- Initial Order Book Snapshot (%s) ---\n", symbol)
	if err := client.Subscribe(ctx, mc.Client, constant.Market_US, symbol,
		[]constant.SubType{constant.SubType_OrderBook}); err != nil {
		fmt.Printf("Subscribe: %v\n", err)
		return
	}

	book, err := client.GetOrderBook(ctx, mc.Client, constant.Market_US, symbol, 10)
	if err != nil {
		fmt.Printf("GetOrderBook: %v\n", err)
	} else {
		printImbalance(book, "SNAPSHOT")
	}

	fmt.Println()
	fmt.Println("--- Real-Time Order Book Push (30s) ---")
	fmt.Println("Press Ctrl+C to stop early")
	fmt.Println()

	ch := make(chan *push.UpdateOrderBook, 100)
	stop := chanpkg.SubscribeOrderBook(ctx, mc.Client, constant.Market_US, symbol, ch)
	defer stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	tick := time.NewTicker(30 * time.Second)
	defer tick.Stop()

loop:
	for {
		select {
		case <-sig:
			fmt.Println("\nInterrupted")
			break loop
		case <-tick.C:
			fmt.Println("\n30s timeout reached")
			break loop
		case ob := <-ch:
			if ob == nil {
				continue
			}
			totalBidVol := int64(0)
			totalAskVol := int64(0)
			for _, bid := range ob.OrderBookBidList {
				totalBidVol += bid.GetVolume()
			}
			for _, ask := range ob.OrderBookAskList {
				totalAskVol += ask.GetVolume()
			}
			bestBid := 0.0
			if len(ob.OrderBookBidList) > 0 {
				bestBid = ob.OrderBookBidList[0].GetPrice()
			}
			bestAsk := 0.0
			if len(ob.OrderBookAskList) > 0 {
				bestAsk = ob.OrderBookAskList[0].GetPrice()
			}

			imbalance := 0.0
			totalVol := totalBidVol + totalAskVol
			if totalVol > 0 {
				imbalance = float64(totalBidVol-totalAskVol) / float64(totalVol) * 100
			}

			arrow := "➡️"
			pressure := "Neutral"
			if imbalance > 10 {
				arrow, pressure = "🟢", "Buy Pressure"
			} else if imbalance < -10 {
				arrow, pressure = "🔴", "Sell Pressure"
			} else if imbalance > 5 {
				arrow, pressure = "↗️", "Slight Buy"
			} else if imbalance < -5 {
				arrow, pressure = "↘️", "Slight Sell"
			}

			fmt.Printf("[PUSH] %s BidVol=%d AskVol=%d Imbalance=%+.1f%% %s (%s)\n",
				arrow, totalBidVol, totalAskVol, imbalance, pressure,
				time.Now().Format("15:04:05"))
			fmt.Printf("       Best: $%.2f × $%.2f\n", bestBid, bestAsk)
		}
	}

	fmt.Println("\n=== Imbalance Analysis Complete ===")
}

func printImbalance(book *client.OrderBook, tag string) {
	if book == nil {
		return
	}
	totalBidVol := int64(0)
	totalAskVol := int64(0)
	weightedBidPrice := 0.0
	weightedAskPrice := 0.0

	for _, b := range book.Bids {
		totalBidVol += b.Volume
		weightedBidPrice += b.Price * float64(b.Volume)
	}
	for _, a := range book.Asks {
		totalAskVol += a.Volume
		weightedAskPrice += a.Price * float64(a.Volume)
	}

	if totalBidVol > 0 {
		weightedBidPrice /= float64(totalBidVol)
	}
	if totalAskVol > 0 {
		weightedAskPrice /= float64(totalAskVol)
	}

	totalVol := totalBidVol + totalAskVol
	imbalance := 0.0
	if totalVol > 0 {
		imbalance = float64(totalBidVol-totalAskVol) / float64(totalVol) * 100
	}

	bestBid := 0.0
	if len(book.Bids) > 0 {
		bestBid = book.Bids[0].Price
	}
	bestAsk := 0.0
	if len(book.Asks) > 0 {
		bestAsk = book.Asks[0].Price
	}

	spread := math.Abs(bestAsk - bestBid)
	microPrice := (weightedBidPrice + weightedAskPrice) / 2

	pressure := "Neutral"
	arrow := "➡️"
	if imbalance > 10 {
		pressure = "Buy Pressure"
		arrow = "🟢"
	} else if imbalance > 5 {
		pressure = "Slight Buy Pressure"
		arrow = "↗️"
	} else if imbalance < -10 {
		pressure = "Sell Pressure"
		arrow = "🔴"
	} else if imbalance < -5 {
		pressure = "Slight Sell Pressure"
		arrow = "↘️"
	}

	fmt.Printf("[%s] %s BidVol=%d AskVol=%d Imbalance=%+.1f%% %s (%s)\n",
		tag, arrow, totalBidVol, totalAskVol, imbalance, pressure, time.Now().Format("15:04:05"))
	fmt.Printf("      Best: $%.2f × $%.2f (spread=$%.4f) | MicroPrice=$%.2f\n",
		bestBid, bestAsk, spread, microPrice)

	for _, side := range []struct {
		name  string
		items []client.OrderBookItem
	}{
		{"Bids", book.Bids},
		{"Asks", book.Asks},
	} {
		for _, item := range side.items {
			if len(item.DetailList) > 0 {
				totalDetailVol := int64(0)
				for _, d := range item.DetailList {
					totalDetailVol += d.Volume
				}
				if totalDetailVol > item.Volume*2 {
					fmt.Printf("      ⚡ Iceberg @ $%.2f: visible=%d hidden=%d orders=%d\n",
						item.Price, item.Volume, totalDetailVol-item.Volume, item.OrderCount)
				}
			}
		}
	}

	fmt.Printf("      Levels: %d bids / %d asks\n", len(book.Bids), len(book.Asks))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
