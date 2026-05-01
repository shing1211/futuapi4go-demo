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

	fmt.Println("=== Order Lifecycle Demo (US Simulated) ===")
	fmt.Println()

	// Get simulated US account
	accounts, err := client.GetAccountList(ctx, cli)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	for _, acc := range accounts {
		if acc.TrdEnv == 0 { // Simulated
			for _, auth := range acc.TrdMarketAuthList {
				if auth == constant.TrdMarket_US.Int32() || auth == constant.TrdMarket_HK.Int32() {
					accID = acc.AccID
					fmt.Printf("Using simulated AccID=%d (TrdEnv=%d)\n", accID, acc.TrdEnv)
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
		fmt.Printf("Using AccID=%d\n", accID)
	}

	// Unlock trading (required for simulated trading)
	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, cli, pwd); err != nil {
			fmt.Printf("UnlockTrading warning: %v\n", err)
		} else {
			fmt.Println("Trading unlocked successfully")
		}
	}

	fmt.Println("\n=== Step 1: Check Account Funds ===")
	funds, err := client.GetFunds(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetFunds failed: %v\n", err)
	} else {
		fmt.Printf("Cash: %.2f | Power: %.2f | Assets: %.2f\n",
			funds.Cash, funds.Power, funds.TotalAssets)
	}

	fmt.Println("\n=== Step 2: Check Current Positions ===")
	positions, err := client.GetPositionList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetPositionList failed: %v\n", err)
	} else {
		fmt.Printf("Current positions: %d\n", len(positions))
		for _, p := range positions {
			fmt.Printf("  %s: Qty=%.0f Cost=%.2f Cur=%.2f\n",
				p.Code, p.Quantity, p.CostPrice, p.CurPrice)
		}
	}

	fmt.Println("\n=== Step 3: Place a Limit Order (Demo) ===")
	// Place a demo order - this will likely fail if market is closed or price is invalid
	// In simulated mode, use a price within reasonable range
	stock := "AAPL"
	price := 180.0
	qty := 1.0

	fmt.Printf("Placing order: Buy %d share(s) of %s @ $%.2f\n", int(qty), stock, price)

	result, err := client.PlaceOrder(ctx, cli, accID, constant.TrdMarket_US,
		stock, constant.TrdSide_Buy, constant.OrderType_Normal, price, qty, 0)
	if err != nil {
		fmt.Printf("PlaceOrder failed: %v\n", err)
	} else {
		fmt.Printf("Order placed successfully! OrderID=%d\n", result.OrderID)
		fmt.Printf("  Status: %s\n", result.Status)
	}

	fmt.Println("\n=== Step 4: List Open Orders ===")
	orders, err := client.GetOrderList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetOrderList failed: %v\n", err)
	} else {
		fmt.Printf("Open orders: %d\n", len(orders))
		for _, o := range orders {
			fmt.Printf("  OrderID=%d %s %s Qty=%.0f Price=%.2f Status=%s\n",
				o.OrderID, o.Code, o.Side, o.Qty, o.Price, o.Status)
		}
	}

	fmt.Println("\n=== Step 5: Modify/Cancel Order (Demo) ===")
	// Note: In real usage, you'd cancel specific orders
	// Here we demonstrate the modify flow
	if len(orders) > 0 {
		orderID := orders[0].OrderID
		newPrice := 185.0
		fmt.Printf("Modifying order %d to price $%.2f...\n", orderID, newPrice)

		modResult, err := client.ModifyOrder(ctx, cli, accID, constant.TrdMarket_US,
			orderID, constant.ModifyOrderOp_Modify, newPrice, 0)
		if err != nil {
			fmt.Printf("ModifyOrder failed: %v\n", err)
		} else {
			fmt.Printf("Modify result: %s\n", modResult.ModifyStatus)
		}
	}

	fmt.Println("\n=== Order Lifecycle Complete ===")
	fmt.Println("Note: In simulated mode, orders are simulated but may fail if:")
	fmt.Println("  - Market is closed")
	fmt.Println("  - Price is outside valid range")
	fmt.Println("  - Insufficient funds")
}