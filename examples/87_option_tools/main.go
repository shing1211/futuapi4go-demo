package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/option"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()

	fmt.Println("=== Option Filtering & Analysis Toolkit Demo ===")
	fmt.Println()

	fmt.Println("--- Step 1: Fetch Option Chain ---")
	now := time.Now()
	end := now.AddDate(0, 0, 45)
	chains, err := client.GetOptionChain(ctx, mc.Client,
		constant.Market_US, "NVDA",
		constant.IndexOptionType_Standard,
		0, 0,
		now.Format("2006-01-02"), end.Format("2006-01-02"),
	)
	if err != nil {
		log.Fatalf("GetOptionChain failed: %v", err)
	}
	fmt.Printf("Found %d expiry dates\n", len(chains))

	quote, err := client.GetQuote(ctx, mc.Client, constant.Market_US, "NVDA")
	if err != nil {
		log.Fatalf("GetQuote failed: %v", err)
	}
	fmt.Printf("NVDA current price: $%.2f\n\n", quote.Price)

	fmt.Println("--- Step 2: Parse Option Codes (HK format) ---")
	examples := []string{"24011900700P", "24021600180C", "230930AAPL220P"}
	for _, c := range examples {
		parsed, err := option.ParseCode(c)
		if err != nil {
			fmt.Printf("  %s: %v\n", c, err)
		} else {
			fmt.Printf("  %-15s → expiry=%s strike=$%.2f type=%s\n",
				c, parsed.Expiry.Format("2006-01-02"), parsed.Strike, parsed.Type)
		}
	}

	fmt.Println()
	fmt.Println("--- Step 3: Build Option Objects from Chain ---")
	var opts []*option.OptionCode
	for _, chain := range chains {
		if chain == nil {
			continue
		}
		strikeTime := chain.StrikeTime
		expiry, err := time.Parse("2006-01-02", strikeTime)
		if err != nil {
			continue
		}
		for _, item := range chain.Option {
			if item == nil {
				continue
			}
			exData := item.Call.GetOptionExData()
			if exData != nil && exData.GetStrikePrice() > 0 {
				opts = append(opts, &option.OptionCode{
					Code:       fmt.Sprintf("NVDA%sC%.0f", expiry.Format("060102"), exData.GetStrikePrice()),
					Expiry:     expiry,
					Strike:     exData.GetStrikePrice(),
					Type:       option.OptionTypeCall,
					Underlying: "NVDA",
				})
			}
			exData = item.Put.GetOptionExData()
			if exData != nil && exData.GetStrikePrice() > 0 {
				opts = append(opts, &option.OptionCode{
					Code:       fmt.Sprintf("NVDA%sP%.0f", expiry.Format("060102"), exData.GetStrikePrice()),
					Expiry:     expiry,
					Strike:     exData.GetStrikePrice(),
					Type:       option.OptionTypePut,
					Underlying: "NVDA",
				})
			}
		}
	}
	fmt.Printf("Built %d option entries\n\n", len(opts))

	fmt.Println("--- Step 4: Filter by Nearest Expiry ---")
	dates := uniqueExpiries(opts)
	var nearest time.Time
	for _, d := range dates {
		if d.After(time.Now()) && (nearest.IsZero() || d.Before(nearest)) {
			nearest = d
		}
	}
	if !nearest.IsZero() {
		fmt.Printf("Nearest expiry: %s\n", nearest.Format("2006-01-02"))
		expiring := option.FilterByExpiry(opts, nearest)
		fmt.Printf("Options: %d\n", len(expiring))
		for _, o := range expiring[:min(8, len(expiring))] {
			fmt.Printf("  %-20s strike=$%.2f type=%s\n", o.Code, o.Strike, o.Type)
		}
	}

	fmt.Println()
	fmt.Println("--- Step 5: Filter by Strike Range ±5% ---")
	minS := quote.Price * 0.95
	maxS := quote.Price * 1.05
	nearMoney := option.FilterByStrikeRange(opts, minS, maxS)
	fmt.Printf("Range: $%.2f – $%.2f → %d options\n", minS, maxS, len(nearMoney))
	for _, o := range nearMoney[:min(8, len(nearMoney))] {
		dist := option.StrikeDistance(o.Strike, quote.Price)
		fmt.Printf("  %-20s strike=$%.2f type=%s (dist=$%.2f)\n",
			o.Code, o.Strike, o.Type, dist)
	}

	fmt.Println()
	fmt.Println("--- Step 6: Find ATM Options ---")
	atm := option.FindAtm(opts, quote.Price)
	if atm != nil {
		if atm.Call != nil {
			fmt.Printf("  ATM Call: %s (strike=$%.2f, expiry=%s)\n",
				atm.Call.Code, atm.Call.Strike, atm.Call.Expiry.Format("2006-01-02"))
		}
		if atm.Put != nil {
			fmt.Printf("  ATM Put:  %s (strike=$%.2f, expiry=%s)\n",
				atm.Put.Code, atm.Put.Strike, atm.Put.Expiry.Format("2006-01-02"))
		}
	}

	fmt.Println()
	fmt.Println("--- Strategy Reference ---")
	fmt.Println("  Covered Call:     Own 100 shares → sell 1 OTM call")
	fmt.Println("  Protective Put:   Own 100 shares → buy 1 ATM/OTM put")
	fmt.Println("  Cash-Secured Put: Hold cash → sell 1 ATM/OTM put")
	fmt.Println("  Call Vertical:    Buy call at strike A, sell call at strike B")
	fmt.Println()
	fmt.Println("Note: Requires FUTU_TRADE_PWD and real trading to execute.")

	fmt.Println("\n=== Option Analysis Complete ===")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(atm)
}

func uniqueExpiries(codes []*option.OptionCode) []time.Time {
	seen := make(map[string]bool)
	var result []time.Time
	for _, o := range codes {
		key := o.Expiry.Format("2006-01-02")
		if !seen[key] {
			seen[key] = true
			result = append(result, o.Expiry)
		}
	}
	return result
}
