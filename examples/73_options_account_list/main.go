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

	fmt.Println("=== Options Trading Overview ===")
	fmt.Println("Note: Options use the same account as stocks (TrdCategory_Security)")
	fmt.Println("Options are NOT a separate TrdCategory like Futures")
	fmt.Println()

	fmt.Println("=== Stock Accounts (TrdCategory_Security) ===")
	resp, err := cli.Trade().GetAccList(ctx, constant.TrdCategory_Security, true)
	if err != nil {
		log.Fatalf("GetAccList(TrdCategory_Security) failed: %v", err)
	}
	if len(resp.AccList) == 0 {
		fmt.Println("(no stock accounts)")
	} else {
		for i, acc := range resp.AccList {
			env := "Real"
			if acc.TrdEnv == 0 {
				env = "Simulated"
			}
			fmt.Printf("Account %d: AccID=%d TrdEnv=%s AccType=%d Auth=%v\n",
				i, acc.AccID, env, acc.AccType, acc.TrdMarketAuthList)
		}
	}

	fmt.Println("\n=== Checking for Options-Enabled Accounts ===")
	for _, acc := range resp.AccList {
		// Check if account has options authorization
		// AcGrantRights_HKOption = 6, AcGrantRights_USOption = 8
		hasHKOption := false
		hasUSOption := false
		for _, right := range acc.TrdMarketAuthList {
			if right == 6 {
				hasHKOption = true
			}
			if right == 8 {
				hasUSOption = true
			}
		}
		if hasHKOption || hasUSOption {
			fmt.Printf("AccID=%d has options rights: HK=%v US=%v\n",
				acc.AccID, hasHKOption, hasUSOption)
		}
	}

	fmt.Println("\n=== Option Expiration Dates Query ===")
	fmt.Println("Use GetOptionExpirationDate to get available expiration dates for an underlying")

	// Example: Query options expiration dates for a US stock (AAPL)
	expirations, err := client.GetOptionExpirationDate(ctx, cli, constant.Market_US, "AAPL")
	if err != nil {
		fmt.Printf("GetOptionExpirationDate failed: %v\n", err)
	} else {
		fmt.Printf("Options expiration dates for AAPL: %d dates\n", len(expirations))
		for i, exp := range expirations {
			if i < 10 {
				fmt.Printf("  %s (%d days, %s)\n", exp.Date, exp.Days, exp.Desc)
			}
		}
		if len(expirations) > 10 {
			fmt.Printf("  ... and %d more\n", len(expirations)-10)
		}
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("Options accounts use TrdCategory_Security (same as stocks)")
	fmt.Println("Check TrdMarketAuthList for rights 6 (HK Option) or 8 (US Option)")
	fmt.Println("GetPositionList returns both stocks and options positions")
	fmt.Println("GetAccountInfo/GetFunds work for options accounts")
}