// 37_warrant demonstrates GetWarrant to search for callable warrants
// on a Hong Kong underlying security (Tencent/00700).
//
// NOTE: GetWarrant only supports HK market warrants (stocks, CBBCs, inline warrants).
// US/China stocks do not have warrant data via this API.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	// Search for call warrants on Tencent (HK.00700), issuer doesn't matter
	warrantResult, err := client.GetWarrant(context.Background(), mc.Client,
		constant.Market_HK, "00700", // Tencent in HK market
		0, 20,                                  // begin, num (get up to 20 warrants)
		constant.WarrantSortField_EffectiveLeverage, true, // sort by effective leverage, ascending
		constant.WarrantType_Buy,               // buy (call) warrants only
		qotcommon.Issuer_Issuer_Unknow,          // all issuers
		constant.WarrantStatus_Normal,          // active trading status
	)
	if err != nil {
		log.Fatalf("GetWarrant failed: %v", err)
	}

	warrants := warrantResult.Items

	if len(warrants) == 0 {
		fmt.Println("No buy warrants found for HK.00700 (Tencent).")
		fmt.Println("Try a different stock or warrant type.")
		return
	}

	fmt.Printf("Found %d buy (call) warrants for HK.00700 (Tencent) (total: %d):\n\n", len(warrants), warrantResult.AllCount)
	fmt.Printf("%-12s %-25s %-8s %-8s %-8s %-10s %s\n",
		"Code", "Name", "Price", "Strike", "Maturity", "Eff Leverage", "Type")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────")

	for _, w := range warrants {
		wtype := "Buy"
		if w.Type == int32(constant.WarrantType_Sell) {
			wtype = "Sell"
		} else if w.Type == int32(constant.WarrantType_Bull) {
			wtype = "Bull"
		} else if w.Type == int32(constant.WarrantType_Bear) {
			wtype = "Bear"
		} else if w.Type == int32(constant.WarrantType_InLine) {
			wtype = "Inline"
		}
		fmt.Printf("%-12s %-25s %-8.3f %-8.3f %-8s %-10.2fx %s\n",
			w.Stock.GetCode(),
			truncate(w.Name, 25),
			w.CurPrice,
			w.StrikePrice,
			w.MaturityTime,
			w.EffectiveLeverage,
			wtype,
		)
	}

	fmt.Println()
	fmt.Println("Tip: Use warrant scanners to filter by premium, delta, or implied volatility.")

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(warrants)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-2] + ".."
}