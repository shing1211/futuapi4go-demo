// 38_option_chain demonstrates GetOptionChain to retrieve US equity option chains
// for a given underlying (NVDA) with a specific expiry window.
//
// The API requires beginTime and endTime to be within 30 days of each other.
// This example fetches weekly options for the near-term expiry.
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

	// Calculate a date range within 30 days
	// Use the next Friday as the target expiry week
	now := time.Now()
	friday := nextFriday(now)
	end := friday.AddDate(0, 0, 14) // 2 weeks after that expiry
	begin := friday.AddDate(0, 0, -7)

	beginStr := begin.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	fmt.Printf("Fetching NVDA option chain from %s to %s\n", beginStr, endStr)
	fmt.Printf("Target expiry: %s (%s)\n\n", friday.Format("2006-01-02"), friday.Weekday())

	chains, err := client.GetOptionChain(context.Background(), mc.Client,
		constant.Market_US, "NVDA",
		1,       // indexOptionType: 1=Standard (US Equity)
		0,       // optType: 0 = All (calls and puts)
		0,       // condition: 0 = All
		beginStr,
		endStr,
	)
	if err != nil {
		log.Fatalf("GetOptionChain failed: %v", err)
	}

	if len(chains) == 0 {
		fmt.Println("No option chain data returned for the specified period.")
		fmt.Println("Try adjusting the date range to cover upcoming expiry dates.")
		return
	}

	fmt.Printf("Found %d option expiry dates:\n\n", len(chains))

	for _, chain := range chains {
		expiryDate := chain.StrikeTime
		if len(chain.Option) == 0 {
			continue
		}

		// Count calls and puts
		var calls, puts int
		for _, opt := range chain.Option {
			if opt == nil {
				continue
			}
			if opt.Call != nil {
				calls++
			}
			if opt.Put != nil {
				puts++
			}
		}

		fmt.Printf("📅 Expiry: %s | Calls: %d | Puts: %d\n", expiryDate, calls, puts)

		// Print sample strikes
		const maxShow = 5
		shown := 0
		for _, opt := range chain.Option {
			if opt == nil || shown >= maxShow {
				break
			}
			if opt.Call != nil {
				sec := opt.Call.GetBasic().GetSecurity()
				code := ""
				if sec != nil {
					code = sec.GetCode()
				}
				strike := ""
				if opt.Call.GetOptionExData() != nil && opt.Call.GetOptionExData().StrikePrice != nil {
					strike = fmt.Sprintf("%.2f", *opt.Call.GetOptionExData().StrikePrice)
				}
				fmt.Printf("  ✅ CALL %s @ strike %s\n", code, strike)
				shown++
			}
		}
		fmt.Println()
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(chains)
}

// nextFriday returns the next Friday from the given date (including if today is Friday).
func nextFriday(from time.Time) time.Time {
	daysUntilFriday := (5 - int(from.Weekday()) + 7) % 7
	if daysUntilFriday == 0 {
		daysUntilFriday = 7 // get next Friday if today is Friday
	}
	return from.AddDate(0, 0, daysUntilFriday)
}