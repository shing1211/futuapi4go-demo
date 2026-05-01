package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
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

	fmt.Println("\n=== Check 1: Market State ===")
	state, err := client.GetMarketState(ctx, cli, constant.Market_US, "AAPL")
	if err != nil {
		fmt.Printf("  GetMarketState failed: %v\n", err)
	} else {
		canTrade := state == 1
		fmt.Printf("  US Market State: %d (%s)\n", state, marketStateString(state))
		if canTrade {
			fmt.Println("  ✓ Market is OPEN")
		} else {
			fmt.Println("  ✗ Market is NOT open for trading")
		}
	}

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

	fmt.Println("\n=== Check 3: Position Limits ===")
	positions, err := client.GetPositionList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("  GetPositionList failed: %v\n", err)
	} else {
		fmt.Printf("  Current positions: %d\n", len(positions))
		totalValue := 0.0
		for _, p := range positions {
			totalValue += p.MarketVal
		}
		fmt.Printf("  Total position value: $%.2f\n", totalValue)

		for _, p := range positions {
			if p.Code == "AAPL" {
				fmt.Printf("  ✗ Already holding %s (Qty=%.0f)\n", p.Code, p.Quantity)
			}
		}
	}

	fmt.Println("\n=== Check 4: Current Quote ===")
	if err := client.Subscribe(ctx, cli, constant.Market_US, "AAPL",
		[]constant.SubType{constant.SubType_Quote}); err != nil {
		fmt.Printf("  Subscribe failed: %v\n", err)
	} else {
		quote, err := client.GetQuote(ctx, cli, constant.Market_US, "AAPL")
		if err != nil {
			fmt.Printf("  GetQuote failed: %v\n", err)
		} else {
			fmt.Printf("  AAPL: Last=$%.2f Open=$%.2f High=$%.2f Low=$%.2f\n",
				quote.Price, quote.Open, quote.High, quote.Low)
		}
	}

	fmt.Println("\n=== Check 5: Security Snapshot ===")
	sec1 := qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("AAPL")}
	sec2 := qotcommon.Security{Market: ptrInt32(int32(constant.Market_US)), Code: ptrStr("TSLA")}
	snapshots, err := client.GetSecuritySnapshot(ctx, cli, []*qotcommon.Security{&sec1, &sec2})
	if err != nil {
		fmt.Printf("  GetSecuritySnapshot failed: %v\n", err)
	} else {
		for _, s := range snapshots {
			fmt.Printf("  %s: Price=%.2f Vol=%d 52W_High=%.2f 52W_Low=%.2f\n",
				s.Name, s.CurPrice, s.Volume, s.Highest52WeeksPrice, s.Lowest52WeeksPrice)
		}
	}

	fmt.Println("\n=== Pre-Trade Check Complete ===")
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

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string { return &v }