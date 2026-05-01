package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	fmt.Println("=== DCA Grid Bot Demo (US Simulated) ===")
	fmt.Println("Dollar Cost Averaging with Grid Strategy")
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

	funds, err := client.GetFunds(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetFunds failed: %v\n", err)
	} else {
		fmt.Printf("Available Power: $%.2f\n", funds.Power)
	}

	targetSymbol := "AAPL"
	basePrice := 180.0
	gridLevels := 5
	investmentPerLevel := funds.Power / float64(gridLevels*2)

	fmt.Printf("\n=== Strategy Configuration ===\n")
	fmt.Printf("Symbol: %s\n", targetSymbol)
	fmt.Printf("Base Price: $%.2f\n", basePrice)
	fmt.Printf("Grid Levels: %d\n", gridLevels)
	fmt.Printf("Investment per level: $%.2f\n", investmentPerLevel)

	quote, err := client.GetQuote(ctx, cli, constant.Market_US, targetSymbol)
	if err != nil {
		fmt.Printf("GetQuote failed: %v\n", err)
	} else {
		fmt.Printf("Current Price: $%.2f\n", quote.Price)
	}

	pwd := os.Getenv("FUTU_TRADE_PWD")
	if pwd != "" {
		if err := client.UnlockTrading(ctx, cli, pwd); err != nil {
			fmt.Printf("UnlockTrading warning: %v\n", err)
		} else {
			fmt.Println("Trading unlocked")
		}
	}

	fmt.Println("\n=== Simulated DCA Grid Orders ===")
	fmt.Printf("%-20s %15s %15s %10s\n", "Order Type", "Price", "Qty", "Status")
	fmt.Println("─────────────────────────────────────────────────")

	for i := 1; i <= gridLevels; i++ {
		buyPrice := basePrice - float64(i)*5.0
		buyQty := investmentPerLevel / buyPrice
		fmt.Printf("%-20s $%14.2f %15.2f %10s\n",
			fmt.Sprintf("Buy Grid %c", 'A'+i-1), buyPrice, buyQty, "SIMULATED")

		sellPrice := basePrice + float64(i)*5.0
		sellQty := investmentPerLevel / sellPrice
		fmt.Printf("%-20s $%14.2f %15.2f %10s\n",
			fmt.Sprintf("Sell Grid %c", 'A'+i-1), sellPrice, sellQty, "SIMULATED")
	}

	fmt.Println("\n=== Current Holdings ===")
	positions, err := client.GetPositionList(ctx, cli, accID)
	if err != nil {
		fmt.Printf("GetPositionList failed: %v\n", err)
	} else {
		for _, p := range positions {
			fmt.Printf("%s: Qty=%.0f Cost=%.2f Current=%.2f P/L=%.2f\n",
				p.Code, p.Quantity, p.CostPrice, p.CurPrice, p.PnL)
		}
	}

	fmt.Println("\n=== Grid Bot Summary ===")
	fmt.Println("In a live implementation, this bot would:")
	fmt.Println("1. Monitor price and execute buys when price drops")
	fmt.Println("2. Execute sells when price rises above grid levels")
	fmt.Println("3. Rebalance positions periodically")
	fmt.Println("4. Track performance and adjust grid spacing")
	fmt.Println("\nNote: This is a demonstration - no real orders placed")
	time.Sleep(1 * time.Second)
}