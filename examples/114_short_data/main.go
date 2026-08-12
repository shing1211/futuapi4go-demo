// 114_short_data demonstrates the short-selling data suite (Tier 2):
//   - GetDailyShortVolume  (daily short-sold volume & ratio for a stock)
//   - GetShortInterest     (outstanding short-interest position)
//   - GetShortSellingRank  (market-wide short-selling ranking)
//
// Maps best to US stocks (they have the deepest short data).
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetshortsellingrank"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	market := constant.Market(constant.Market_US)
	code := "AAPL"

	fmt.Println("=== Short Data (US.AAPL) ===")

	fmt.Println("--- 1) GetDailyShortVolume ---")
	vol, err := client.GetDailyShortVolume(ctx, mc.Client, market, code, "" /*nextKey*/, 10 /*num*/)
	if err != nil {
		log.Panicf("GetDailyShortVolume failed: %v", err)
	}
	fmt.Printf("  US items: %d | HK items: %d (nextKey=%q)\n",
		len(vol.UsItemList), len(vol.HkItemList), vol.NextKey)
	if len(vol.UsItemList) > 0 {
		it := vol.UsItemList[0]
		fmt.Printf("  aggregatedShort=%d aggregatedShortRatio=%.1f%%\n",
			vol.AggregatedShort, vol.AggregatedShortRatio)
		_ = it
	}
	fmt.Println()

	fmt.Println("--- 2) GetShortInterest ---")
	interest, err := client.GetShortInterest(ctx, mc.Client, market, code, "" /*nextKey*/, 10 /*num*/)
	if err != nil {
		log.Panicf("GetShortInterest failed: %v", err)
	}
	fmt.Printf("  US items: %d | HK items: %d (nextKey=%q)\n",
		len(interest.UsItemList), len(interest.HkItemList), interest.NextKey)
	fmt.Println()

	fmt.Println("--- 3) GetShortSellingRank (US top 10 by short volume) ---")
	rank, err := client.GetShortSellingRank(ctx, mc.Client, &qotgetshortsellingrank.C2S{
		Market:    ptrInt32(int32(constant.Market_US)),
		SortField: ptrInt32(0), // ShortSellingSortField_ShortVolValue
		SortDir:   ptrInt32(0), // SortDir_Descend
		Count:     ptrInt32(10),
	})
	if err != nil {
		log.Panicf("GetShortSellingRank failed: %v", err)
	}
	fmt.Printf("  Rank rows: %d (all=%d)\n", len(rank.DataList), rank.GetAllCount())
	for i, it := range rank.DataList {
		if it == nil {
			continue
		}
		var s string
		if it.Security != nil {
			s = it.Security.GetCode()
		}
		fmt.Printf("  [%2d] %-8s  shortNumber=%d  ratio=%.2f%%  posRatio=%.2f%%\n",
			i+1, s, it.GetShortNumber(), it.GetShortRatio(), it.GetShortPositionRatio())
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	for _, r := range []any{vol, interest, rank} {
		if r != nil {
			display.PrintJSON(r)
		}
	}
}

func ptrInt32(v int32) *int32 { return &v }
