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

	fmt.Println("=== Pre-Trade Checks Demo (US Simulated) ===")
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

	// Step 1: Check Market State
	fmt.Println("\n=== Check 1: Market State ===")
	state, err := client.GetMarketState(ctx, cli, constant.Market_US, "AAPL")
	if err != nil {
		fmt.Printf("  GetMarketState failed: %v\n", err)
	} else {
		canTrade := state == 1 // Trading = 1
		fmt.Printf("  US Market State: %d (%s)\n", state, marketStateString(state))
		if canTrade {
			fmt.Println("  ✓ Market is OPEN")
		} else {
			fmt.Println("  ✗ Market is NOT open for trading")
		}
	}

	// Step 2: Check Account Funds
	fmt.Println("\n=== Check 2: Account Funds ===")
	funds, err := client.GetFunds(ctx, cli, accID)
	if err != nil {
		fmt.Printf("  GetFunds failed: %v\n", err)
	} else {
		fmt.Printf("  Cash: $%.2f | Power: $%.2f\n", funds.Cash, funds.Power)
		if funds.Power > 1000 {
			fmt.Println("  ✓ Sufficient funds available")
		} else {
			fmt.Println("  ✗ Insufficient funds")
		}
	}

	// Step 3: Check Current Positions
	fmt.Println("\n=== Check 3: Position Limits ===")
	positions, err := client.GetPositionList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("  GetPositionList failed: %v\n", err)
	} else {
		fmt.Printf("  Current positions: %d\n", len(positions))
		// Calculate position value
		totalValue := 0.0
		for _, p := range positions {
			totalValue += p.MarketVal
		}
		fmt.Printf("  Total position value: $%.2f\n", totalValue)

		// Check if already holding AAPL (for example)
		for _, p := range positions {
			if p.Code == "AAPL" {
				fmt.Printf("  ✗ Already holding %s (Qty=%.0f)\n", p.Code, p.Quantity)
			}
		}
	}

	// Step 4: Get Quote for Target Symbol
	fmt.Println("\n=== Check 4: Current Quote ===")
	quote, err := client.GetQuote(ctx, cli, constant.Market_US, "AAPL")
	if err != nil {
		fmt.Printf("  GetQuote failed: %v\n", err)
	} else {
		fmt.Printf("  AAPL: Last=$%.2f Bid=$%.2f Ask=$%.2f\n",
			quote.LastPrice, quote.BidPrice[0], quote.AskPrice[0])
	}

	// Step 5: Get Security Snapshot (comprehensive)
	fmt.Println("\n=== Check 5: Security Snapshot ===")
	snapshots, err := client.GetSecuritySnapshot(ctx, cli, []*constant.Security{
		{Market: constant.Market_US, Code: "AAPL"},
		{Market: constant.Market_US, Code: "TSLA"},
	})
	if err != nil {
		fmt.Printf("  GetSecuritySnapshot failed: %v\n", err)
	} else {
		for _, s := range snapshots {
			fmt.Printf("  %s: Price=%.2f Volume=%d 52W_High=%.2f 52W_Low=%.2f\n",
				s.Name, s.LastPrice, s.Volume, s.High52Week, s.Low52Week)
		}
	}

	fmt.Println("\n=== Pre-Trade Check Complete ===")
	fmt.Println("Summary: All checks completed. Ready to trade if all checks pass.")
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
		return "Unknown"
	}
}