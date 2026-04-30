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

	fmt.Println("=== Options Cash/Buying Power ===")
	fmt.Println("Note: Options use same account as stocks, retrieved via GetFunds()")
	fmt.Println()

	// Get stock accounts
	resp, err := cli.Trade().GetAccList(ctx, constant.TrdCategory_Security, true)
	if err != nil {
		log.Fatalf("GetAccList(TrdCategory_Security) failed: %v", err)
	}

	for _, acc := range resp.AccList {
		if acc.TrdEnv != 0 {
			continue // skip real accounts
		}

		fmt.Printf("Account AccID=%d:\n", acc.AccID)

		// Get funds using GetAccountInfo (works for options-enabled accounts)
		for _, mkt := range acc.TrdMarketAuthList {
			market := constant.TrdMarket(mkt)
			funds, err := client.GetAccountInfo(ctx, cli, acc.AccID, market)
			if err != nil {
				fmt.Printf("  Market %s: GetAccountInfo failed: %v\n", market, err)
				continue
			}
			fmt.Printf("  Market=%s: Cash=%.2f Power=%.2f Frozen=%.2f\n",
				market, funds.Cash, funds.Power, funds.FrozenCash)
			fmt.Printf("           Assets=%.2f UnrealizedPL=%.2f RealizedPL=%.2f\n",
				funds.TotalAssets, funds.UnrealizedPL, funds.RealizedPL)
		}

		// Also try GetFunds for cash info
		funds2, err := client.GetFunds(ctx, cli, acc.AccID)
		if err != nil {
			fmt.Printf("  GetFunds failed: %v\n", err)
		} else {
			fmt.Printf("  GetFunds: Cash=%.2f Available=%.2f Frozen=%.2f\n",
				funds2.Cash, funds2.Power, funds2.FrozenCash)
		}
	}

	fmt.Println("\n=== GetAccTradingInfo for Options ===")
	fmt.Println("Note: GetAccTradingInfo can be used for options margin calculation")

	// Try to get options trading info for a sample option
	sampleOptionCode := "AAPL250516C00150000" // Example: AAPL May 16 2025 Call $150

	for _, acc := range resp.AccList {
		if acc.TrdEnv != 0 {
			continue
		}

		// Use US market for options
		info, err := client.GetAccTradingInfo(ctx, cli, acc.AccID,
			constant.TrdMarket_US, sampleOptionCode,
			constant.OrderType_Normal, 5.0) // $5 premium estimate
		if err != nil {
			fmt.Printf("AccID %d: GetAccTradingInfo failed: %v\n", acc.AccID, err)
			continue
		}

		fmt.Printf("Account AccID=%d (option %s):\n", acc.AccID, sampleOptionCode)
		fmt.Printf("  MaxCashBuy: %.2f\n", info.MaxCashBuy)
		fmt.Printf("  MaxCashAndMarginBuy: %.2f\n", info.MaxCashAndMarginBuy)
		fmt.Printf("  MaxPositionSell: %.2f\n", info.MaxPositionSell)
	}

	fmt.Println("\n=== Note ===")
	fmt.Println("Options buying power is calculated based on cash + margin enabled")
	fmt.Println("Use GetAccountInfo with specific market for options-enabled accounts")
}