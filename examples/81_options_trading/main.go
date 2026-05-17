package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Options Trading Workflow Demo (US) ===")
	fmt.Println("Analyses option chain, check margin, place covered call")
	fmt.Println()

	fmt.Println("--- Step 1: Check Account ---")
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

	funds, err := client.GetFunds(ctx, mc.Client, accID)
	if err != nil {
		fmt.Printf("GetFunds failed: %v\n", err)
	} else {
		fmt.Printf("Buying Power: $%.2f | Cash: $%.2f\n", funds.Power, funds.Cash)
	}

	fmt.Println()
	fmt.Println("--- Step 2: Fetch Option Chain ---")
	now := time.Now()
	end := now.AddDate(0, 0, 45)
	beginStr := now.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	chains, err := client.GetOptionChain(ctx, mc.Client,
		constant.Market_US, "NVDA",
		constant.IndexOptionType_Standard,
		0,
		0,
		beginStr, endStr,
	)
	if err != nil {
		log.Fatalf("GetOptionChain failed: %v", err)
	}

	fmt.Printf("Expiry range: %s to %s\n", beginStr, endStr)
	if len(chains) == 0 {
		fmt.Println("No option data returned for the specified period.")
		return
	}
	fmt.Printf("Found %d expiry dates\n", len(chains))

	quote, err := client.GetQuote(ctx, mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		fmt.Printf("GetQuote failed: %v (using estimated $120)\n", err)
	}

	underlyingPrice := 120.0
	if quote != nil {
		underlyingPrice = quote.Price
	}
	fmt.Printf("Underlying price: $%.2f\n", underlyingPrice)

	fmt.Println()
	fmt.Println("--- Step 3: Analyze OTM Call Options ---")
	fmt.Println("Looking for OTM covered calls (strike above current price)")
	fmt.Println()

	fmt.Printf("%-12s %-12s %-12s %-12s %s\n",
		"Expiry", "Strike", "Option Code", "Type", "Days")
	fmt.Println(stringRepeat("─", 70))

	count := 0
	for _, chain := range chains {
		if count >= 8 {
			break
		}
		for _, opt := range chain.Option {
			if count >= 8 {
				break
			}
			if opt == nil || opt.Call == nil {
				continue
			}
			staticEx := opt.Call.GetOptionExData()
			if staticEx == nil {
				continue
			}
			strike := staticEx.GetStrikePrice()
			if strike <= underlyingPrice {
				continue
			}

			code := ""
			if sec := opt.Call.GetBasic().GetSecurity(); sec != nil {
				code = sec.GetCode()
			}

			strikeTime := chain.StrikeTime
			daysToExpiry := "?"
			if t, err := time.Parse("2006-01-02", strikeTime); err == nil {
				d := int(t.Sub(now).Hours() / 24)
				daysToExpiry = fmt.Sprintf("%dd", d)
			}

			fmt.Printf("%-12s $%-10.2f %-12s %-12s %s\n",
				strikeTime, strike, code, "Call", daysToExpiry)
			count++
		}
	}

	if count == 0 {
		fmt.Println("No OTM call options found in the date range.")
		fmt.Println("Try extending the expiry window.")
		return
	}

	fmt.Println()
	fmt.Println("--- Step 4: Covered Call Margin Check ---")
	maintenanceReq := underlyingPrice * 100 * 0.25
	fmt.Printf("Estimated maint. req. per contract: $%.2f (25%% of underlying x100)\n", maintenanceReq)
	if funds != nil {
		canSell := funds.Power / maintenanceReq
		fmt.Printf("Buying power allows ~%.0f covered call contract(s)\n", canSell)
	}

	fmt.Println()
	fmt.Println("--- Step 5: Options Order (Display Only) ---")
	fmt.Println("Selling a covered call requires:")
	fmt.Println("  1. Own 100 shares of NVDA (or margin equity to cover)")
	fmt.Println("  2. Real trading environment (TrdEnv_Real)")
	fmt.Println("  3. FUTU_TRADE_PWD set for trade unlock")
	fmt.Println()
	fmt.Println("Simulated broker does not support options orders.")
	fmt.Println("To test with real trading:")
	fmt.Println("  - Set FUTU_TRADE_PWD environment variable")
	fmt.Println("  - Use WithTradeEnv(constant.TrdEnv_Real)")
	fmt.Println("  - Options order uses OrderType_Normal, qty in contracts")

	fmt.Println()
	fmt.Println("=== Options Analysis Complete ===")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(quote)
}

func stringRepeat(s string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = s[0]
	}
	return string(b)
}
