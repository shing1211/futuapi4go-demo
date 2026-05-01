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

	fmt.Println("=== Market Hours Check (US + HK Stocks) ===")
	fmt.Println("Note: Simulator may return extended state values")
	fmt.Println()

	fmt.Println("--- US Market ---")
	state, err := client.GetMarketState(ctx, cli, constant.Market_US, "AAPL")
	if err != nil {
		fmt.Printf("GetMarketState failed: %v\n", err)
	} else {
		fmt.Printf("Market State: %d (%s)\n", state, marketStateString(state))
	}

	fmt.Println("\n--- US Trading Days ---")
	tradeDates, err := client.GetTradeDate(ctx, cli, constant.Market_US, "", "")
	if err != nil {
		fmt.Printf("GetTradeDate failed: %v (expected for simulator)\n", err)
	} else {
		fmt.Printf("Recent trading days: %d days available\n", len(tradeDates))
		if len(tradeDates) > 0 {
			fmt.Printf("Last 5: ")
			for i := len(tradeDates) - 5; i < len(tradeDates); i++ {
				if i >= 0 {
					fmt.Printf("%s ", tradeDates[i])
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("\n--- HK Market ---")
	hkState, err := client.GetMarketState(ctx, cli, constant.Market_HK, "00100")
	if err != nil {
		fmt.Printf("GetMarketState failed: %v\n", err)
	} else {
		fmt.Printf("Market State: %d (%s)\n", hkState, marketStateString(hkState))
	}

	hkDates, err := client.GetTradeDate(ctx, cli, constant.Market_HK, "", "")
	if err != nil {
		fmt.Printf("GetTradeDate failed: %v (expected for simulator)\n", err)
	} else {
		fmt.Printf("Recent trading days: %d days available\n", len(hkDates))
	}

	fmt.Println("\n=== Pre-Trade Check Helper ===")
	checkMarketReadiness(constant.Market_US, "AAPL", cli, ctx)
	checkMarketReadiness(constant.Market_HK, "00100", cli, ctx)
}

func checkMarketReadiness(market constant.Market, code string, cli *client.Client, ctx context.Context) {
	fmt.Printf("\nChecking %s market readiness for %s...\n", market, code)

	state, err := client.GetMarketState(ctx, cli, market, code)
	if err != nil {
		fmt.Printf("  Cannot get market state: %v\n", err)
		return
	}

	switch state {
	case 0:
		fmt.Println("  Status: PRE-MARKET (not tradeable)")
	case 1:
		fmt.Println("  Status: OPEN (tradeable)")
	case 2:
		fmt.Println("  Status: POST-MARKET (not tradeable)")
	case 3:
		fmt.Println("  Status: CLOSED (not tradeable)")
	default:
		fmt.Printf("  Status: EXTENDED (%d) - may indicate simulator mode\n", state)
	}
}

func marketStateString(state int32) string {
	switch state {
	case 0:
		return "Pre-Market"
	case 1:
		return "Trading"
	case 2:
		return "Post-Market"
	case 3:
		return "Closed"
	default:
		return fmt.Sprintf("Extended(%d)", state)
	}
}