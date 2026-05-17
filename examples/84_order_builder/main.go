package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/trd"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== OrderBuilder Fluent API Demo ===")
	fmt.Println("Showcasing the recommended way to construct orders")
	fmt.Println()

	accounts, err := client.GetAccountList(ctx, mc.Client)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	var accID uint64
	for _, acc := range accounts {
		if acc.TrdEnv == 0 {
			accID = acc.AccID
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

	trdAPI := mc.Client.Trade()

	var resp *trd.PlaceOrderResponse
	fmt.Println("\n--- Pattern A: Limit Buy Order ---")
	reqA, err := trdAPI.NewOrder(accID, constant.TrdMarket_US).
		Buy("US.NVDA", 100).
		At(120.50).
		WithRemark("limit buy demo").
		AutoDetectMarket().
		Build()
	if err != nil {
		fmt.Printf("Build failed: %v\n", err)
	} else {
		fmt.Printf("%-20s %s\n", "Code", reqA.Code)
		fmt.Printf("%-20s Buy\n", "Side")
		fmt.Printf("%-20s Normal (limit)\n", "OrderType")
		fmt.Printf("%-20s $%.2f\n", "Price", reqA.Price)
		fmt.Printf("%-20s %.0f\n", "Qty", reqA.Qty)
		fmt.Printf("%-20s %v\n", "TrdMarket", reqA.TrdMarket)
		fmt.Printf("%-20s %v\n", "SecMarket", reqA.SecMarket)
		fmt.Printf("%-20s %s\n", "Remark", reqA.Remark)

		resp, err = trdAPI.PlaceOrder(ctx, reqA)
		if err != nil {
			fmt.Printf("PlaceOrder (simulate): %v\n", err)
		} else {
			fmt.Printf("Placed: OrderID=%d\n", resp.OrderID)
		}
	}

	fmt.Println("\n--- Pattern B: Market Sell Order with Remark ---")
	reqB, err := trdAPI.NewOrder(accID, constant.TrdMarket_US).
		Sell("US.NVDA", 100).
		Market().
		WithRemark("market sell demo").
		AutoDetectMarket().
		Build()
	if err != nil {
		fmt.Printf("Build failed: %v\n", err)
	} else {
		fmt.Printf("%-20s %s\n", "Code", reqB.Code)
		fmt.Printf("%-20s Sell\n", "Side")
		fmt.Printf("%-20s Market\n", "OrderType")
		fmt.Printf("%-20s $%.2f\n", "Price", reqB.Price)
		fmt.Printf("%-20s %.0f\n", "Qty", reqB.Qty)

		resp, err := trdAPI.PlaceOrder(ctx, reqB)
		if err != nil {
			fmt.Printf("PlaceOrder (simulate): %v\n", err)
		} else {
			fmt.Printf("Placed: OrderID=%d\n", resp.OrderID)
		}
	}

	fmt.Println("\n--- Pattern C: GTC Limit Buy ---")
	reqC, err := trdAPI.NewOrder(accID, constant.TrdMarket_US).
		Buy("US.AAPL", 50).
		At(175.00).
		WithTimeInForce(constant.TimeInForce_GTC).
		AutoDetectMarket().
		Build()
	if err != nil {
		fmt.Printf("Build failed: %v\n", err)
	} else {
		fmt.Printf("%-20s %s\n", "Code", reqC.Code)
		fmt.Printf("%-20s Buy\n", "Side")
		fmt.Printf("%-20s Normal (limit)\n", "OrderType")
		fmt.Printf("%-20s $%.2f\n", "Price", reqC.Price)
		fmt.Printf("%-20s %.0f\n", "Qty", reqC.Qty)
		fmt.Printf("%-20s GTC\n", "TimeInForce")

		resp, err := trdAPI.PlaceOrder(ctx, reqC)
		if err != nil {
			fmt.Printf("PlaceOrder (simulate): %v\n", err)
		} else {
			fmt.Printf("Placed: OrderID=%d\n", resp.OrderID)
		}
	}

	fmt.Println("\n=== Builder Methods Reference ===")
	fmt.Println("NewOrder() -> Buy()/Sell() -> At() -> Market() -> Build()")
	fmt.Println("Chain modifiers: WithRemark(), AutoDetectMarket(), WithSecMarket(),")
	fmt.Println("  WithTimeInForce(), WithFillOutsideRTH(), WithAuxPrice()")
	fmt.Println("\nNote: Requires FUTU_TRADE_PWD and WithTradeEnv(TrdEnv_Real)")
	fmt.Println("  to execute orders. Simulate mode displays request fields only.")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(resp)
}
