package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	ctx := context.Background()

	fmt.Println("=== VWAP Executor Demo (US Simulated) ===")
	fmt.Println("Volume Weighted Average Price execution strategy")
	fmt.Println()

	// Get simulated US account
	accounts, err := client.GetAccountList(ctx, cli)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			for _, auth := range acc.TrdMarketAuthList {
				if auth == constant.TrdMarket_US.Int32() {
					accID = acc.AccID
					break
				}
			}
		}
		if accID != 0 {
			break
		}
	}
	if accID == 0 {
		accID = accounts[0].AccID
	}
	fmt.Printf("Using AccID=%d\n", accID)

	// Target execution
	targetSymbol := "AAPL"
	targetQty := 100.0
	startPrice := 180.0

	fmt.Printf("\n=== Execution Plan ===\n")
	fmt.Printf("Symbol: %s\n", targetSymbol)
	fmt.Printf("Quantity: %.0f shares\n", targetQty)
	fmt.Printf("Target Start Price: $%.2f\n", startPrice)

	// Step 1: Get order book for VWAP calculation
	fmt.Println("\n--- Step 1: Market Order Book Analysis ---")

	orderbook, err := client.GetOrderBook(ctx, cli, constant.Market_US, targetSymbol, 20)
	if err != nil {
		fmt.Printf("GetOrderBook failed: %v\n", err)
	} else {
		fmt.Printf("Order book levels: %d bids, %d asks\n",
			len(orderbook.BidItems), len(orderbook.AskItems))

		// Calculate VWAP from order book
		var bidVolume int64
		var bidValue float64
		for i, item := range orderbook.BidItems {
			if i >= 10 {
				break
			}
			bidVolume += int64(item.Volume)
			bidValue += float64(item.Volume) * item.Price
		}

		var askVolume int64
		var askValue float64
		for i, item := range orderbook.AskItems {
			if i >= 10 {
				break
			}
			askVolume += int64(item.Volume)
			askValue += float64(item.Volume) * item.Price
		}

		bidVWAP := bidValue / float64(bidVolume)
		askVWAP := askValue / float64(askVolume)

		fmt.Printf("Bid VWAP (10 levels): $%.2f (vol=%d)\n", bidVWAP, bidVolume)
		fmt.Printf("Ask VWAP (10 levels): $%.2f (vol=%d)\n", askVWAP, askVolume)
		fmt.Printf("Spread: $%.2f (%.2f%%)\n", askVWAP-bidVWAP, (askVWAP-bidVWAP)/bidVWAP*100)
	}

	// Step 2: Check account info
	fmt.Println("\n--- Step 2: Account Buying Power ---")

	funds, err := client.GetFunds(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetFunds failed: %v\n", err)
	} else {
		maxBuy := targetQty * startPrice
		if funds.Power >= maxBuy {
			fmt.Printf("✓ Sufficient power: $%.2f >= $%.2f needed\n", funds.Power, maxBuy)
		} else {
			fmt.Printf("✗ Insufficient power: $%.2f < $%.2f needed\n", funds.Power, maxBuy)
		}
	}

	// Step 3: Get maximum trading quantity
	fmt.Println("\n--- Step 3: Max Position Check ---")

	maxQtys, err := client.GetMaxTrdQtys(ctx, cli, accID, constant.TrdMarket_US,
		targetSymbol, constant.OrderType_Normal, startPrice, 0)
	if err != nil {
		fmt.Printf("GetMaxTrdQtys failed: %v\n", err)
	} else {
		fmt.Printf("Max buy qty: %.0f\n", maxQtys.MaxBuyQty)
		fmt.Printf("Max sell qty: %.0f\n", maxQtys.MaxSellQty)
	}

	// Step 4: Simulate VWAP execution
	fmt.Println("\n--- Step 4: Simulated VWAP Execution ---")

	// VWAP strategy: split order into chunks based on volume分布
	// For demo, we simulate 5 time slices
	slices := 5
	qtyPerSlice := targetQty / float64(slices)

	fmt.Printf("VWAP Execution Plan:\n")
	fmt.Printf("%-10s %15s %15s %15s\n", "Slice", "Price", "Qty", "Value")
	fmt.Println("─────────────────────────────────────────────────")

	for i := 1; i <= slices; i++ {
		// Simulate price based on order book
		execPrice := startPrice + float64(i-3)*0.5 // Slight variation
		execQty := qtyPerSlice
		execValue := execPrice * execQty
		fmt.Printf("%-10d $%14.2f %15.2f $%14.2f\n", i, execPrice, execQty, execValue)
	}

	totalValue := startPrice * targetQty
	fmt.Printf("%-10s %15s %15.2f $%14.2f\n", "TOTAL", "~VWAP", targetQty, totalValue)

	// Step 5: Check open orders
	fmt.Println("\n--- Step 5: Current Open Orders ---")

	orders, err := client.GetOrderList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetOrderList failed: %v\n", err)
	} else {
		fmt.Printf("Open orders: %d\n", len(orders))
		for _, o := range orders {
			fmt.Printf("  %d: %s %s Qty=%.0f Price=%.2f Status=%s\n",
				o.OrderID, o.Code, o.Side, o.Qty, o.Price, o.Status)
		}
	}

	fmt.Println("\n=== VWAP Executor Complete ===")
	fmt.Println("Note: This is a demonstration. In production, orders would be")
	fmt.Println("placed incrementally based on real-time volume and price.")
}