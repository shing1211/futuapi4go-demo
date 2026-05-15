package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Trailing Stop Order Demo (US Simulated) ===")
	fmt.Println("Demonstrates trailing stop order parameters and placement")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, mc.Client)
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

	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, mc.Client, pwd); err != nil {
			fmt.Printf("UnlockTrading warning: %v\n", err)
		}
	}

	fmt.Println("\n--- What is a Trailing Stop? ---")
	fmt.Println("A trailing stop follows the market price upward.")
	fmt.Println("If price rises, the stop rises with it by the trail amount.")
	fmt.Println("If price falls by the trail spread, the stop triggers.")
	fmt.Println()

	fmt.Println("--- Worked Example ---")
	fmt.Println("Buy 100 NVDA at $120.00")
	fmt.Println("Set trailing stop sell: TrailType=Ratio, TrailValue=2%, TrailSpread=$0.10")
	fmt.Println("  NVDA rises to $125.00 → stop trails to $122.50 (2% below peak)")
	fmt.Println("  NVDA rises to $130.00 → stop trails to $127.40")
	fmt.Println("  NVDA drops from $130.00 → triggers at $127.40")
	fmt.Println()

	trdAPI := mc.Client.Trade()

	fmt.Println("--- Building Trailing Stop Order ---")
	req, err := trdAPI.NewOrder(accID, constant.TrdMarket_US).
		Sell("US.NVDA", 100).
		Market().
		WithAuxPrice(118.00).
		WithTrailType(constant.TrailType_Ratio).
		WithTrailValue(2.0).
		WithSpread(0.10).
		AutoDetectMarket().
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	fmt.Printf("%-22s %s\n", "Code", req.Code)
	fmt.Printf("%-22s Sell\n", "Side")
	fmt.Printf("%-22s TrailingStop\n", "OrderType")
	fmt.Printf("%-22s %.0f\n", "Qty", req.Qty)
	fmt.Printf("%-22s $%.2f (activate below)\n", "AuxPrice", req.AuxPrice)
	fmt.Printf("%-22s Ratio (%d)\n", "TrailType", req.TrailType)
	fmt.Printf("%-22s %.1f%%\n", "TrailValue", req.TrailValue)
	fmt.Printf("%-22s $%.2f\n", "TrailSpread", req.TrailSpread)

	resp, err := trdAPI.PlaceOrder(ctx, req)
	if err != nil {
		fmt.Printf("\nPlaceOrder (simulate): %v\n", err)
	} else {
		fmt.Printf("\nPlaced: OrderID=%d\n", resp.OrderID)
	}

	fmt.Println("\n--- Order Verification ---")
	orders, err := client.GetOrderList(ctx, mc.Client, accID)
	if err != nil {
		fmt.Printf("GetOrderList failed: %v\n", err)
	} else {
		fmt.Printf("Open orders: %d\n", len(orders))
		for _, o := range orders {
			if o.OrderType != int32(constant.OrderType_TrailingStop) {
				continue
			}
			fmt.Printf("  OrderID=%d %s\n", o.OrderID, o.Code)
			fmt.Printf("    AuxPrice=$%.2f TrailType=%d TrailValue=%.1f TrailSpread=$%.2f\n",
				o.AuxPrice, o.TrailType, o.TrailValue, o.TrailSpread)
		}
	}

	fmt.Println("\n--- Trailing Stop Parameter Reference ---")
	fmt.Println("TrailType=Ratio (1):   TrailValue = percentage (e.g. 2 = 2%)")
	fmt.Println("TrailType=Amount (2):  TrailValue = dollar amount")
	fmt.Println("TrailSpread:            Minimum distance between stop and market")
	fmt.Println("AuxPrice:               Activation price (start trailing when reached)")
	fmt.Println()
	fmt.Println("Note: Trailing stop requires FUTU_TRADE_PWD and real trading.")
	fmt.Println("  Simulate mode displays parameters but cannot execute trailing orders.")
}
