package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Options Positions ===")
	fmt.Println("Note: Options positions are returned by client.GetPositionList() with stock account")
	fmt.Println("Options use the same account and position API as stocks")
	fmt.Println()

	// Get stock accounts
	resp, err := mc.Client.Trade().GetAccList(ctx, constant.TrdCategory_Security, true)
	if err != nil {
		log.Fatalf("GetAccList(TrdCategory_Security) failed: %v", err)
	}

	for _, acc := range resp.AccList {
		if acc.TrdEnv != 0 {
			continue // skip real accounts
		}

		fmt.Printf("Account AccID=%d:\n", acc.AccID)

		positions, err := client.GetPositionList(ctx, mc.Client, acc.AccID)
		if err != nil {
			fmt.Printf("  GetPositionList failed: %v\n", err)
			continue
		}

		if len(positions) == 0 {
			fmt.Println("  (no positions)")
			continue
		}

		fmt.Printf("  %d total positions:\n", len(positions))
		for _, p := range positions {
			fmt.Printf("  %s %s: Qty=%.0f Cost=%.2f Cur=%.2f P/L=%.2f\n",
				p.Code, p.Name, p.Quantity, p.CostPrice, p.CurPrice, p.PnL)
		}
	}

	fmt.Println("\n=== Option Code Format ===")
	fmt.Println("US Options: AAPL250516C00150000 = AAPL + Expiry + Call/Put + Strike")
	fmt.Println("  Format: [Underlying][YYMMDD][C/P][StrikePrice]")
	fmt.Println("  Example: AAPL250516C00150000 = AAPL Call May 16 2025 $150.00")
	fmt.Println()
	fmt.Println("HK Options: Format varies by期权 type (standard/exotic)")
}