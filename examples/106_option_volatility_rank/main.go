// 106_option_volatility_rank demonstrates option volatility and ranking APIs:
//   - GetOptionVolatility       (implied/historical vol for an underlying)
//   - GetOptionRank             (option ranking by volume/OI/turnover/etc.)
//   - GetOptionMarketStatistic  (option market-wide stat over a window)
//
// Markets covered: US option market (OptionMarket_US_Security).
// All three calls take NVDA as the underlying.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qotoptioncommon "github.com/shing1211/futuapi4go/pkg/pb/qotoptioncommon"
	qotgetoptionmarketstatistic "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionmarketstatistic"
	qotgetoptionrank "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionrank"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	market := int32(constant.Market_US)
	underlying := "NVDA"

	fmt.Println("=== Option Volatility & Rank ===")
	fmt.Printf("Underlying: %s (US)\n\n", underlying)

	fmt.Println("--- 1) GetOptionVolatility ---")
	rsp, err := client.GetOptionVolatility(ctx, mc.Client, constant.Market(market), underlying,
		qotcommon.OptionVolatilityTimePeriodType_OptionVolatilityTimePeriodType_Month, 0)
	if err != nil {
		log.Fatalf("GetOptionVolatility failed: %v", err)
	}
	fmt.Printf("Average IV:    %.2f%%\n", rsp.AverageImpvol)
	fmt.Printf("Status code:   %d\n", rsp.ImpvolStatus)
	fmt.Printf("Analysis:      %s\n", rsp.Analysis)
	fmt.Printf("Volatility items: %d\n", len(rsp.ItemList))
	for i, it := range rsp.ItemList {
		if it == nil {
			continue
		}
		ts := time.Unix(it.GetTimestamp(), 0).Format("2006-01-02")
		fmt.Printf("  [%d] impliedVol=%.2f%%  hv=%.2f%%  date=%s\n",
			i+1, it.GetImpliedVolatility(), it.GetHistoryVolatility(), ts)
	}
	fmt.Println()

	fmt.Println("--- 2) GetOptionRank (top 10 by Volume) ---")
	rankRsp, err := client.GetOptionRank(ctx, mc.Client, &qotgetoptionrank.C2S{
		OptionMarket: ptrInt32(int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security)),
		SortType:     ptrInt32(int32(qotoptioncommon.OptionRankType_OptionRankType_Volume)),
		Count:        ptrInt32(10),
	})
	if err != nil {
		log.Fatalf("GetOptionRank failed: %v", err)
	}
	fmt.Printf("Total: %d  |  Items returned: %d\n", rankRsp.GetAllCount(), len(rankRsp.RankList))
	for i, item := range rankRsp.RankList {
		if item == nil {
			continue
		}
		sec := item.GetOption()
		code := ""
		if sec != nil {
			code = sec.GetCode()
		}
		fmt.Printf("  #%d %-22s vol=%d  OI=%d  IV=%.2f%%  type=%d\n",
			i+1, code, item.GetVolume(), item.GetOpenInterest(), item.GetIv(), item.GetOptionType())
	}
	fmt.Println()

	fmt.Println("--- 3) GetOptionMarketStatistic (Volume, last 5 trading days) ---")
	end := time.Now()
	begin := end.AddDate(0, 0, -5)
	statRsp, err := client.GetOptionMarketStatistic(ctx, mc.Client, &qotgetoptionmarketstatistic.C2S{
		OptionMarket: ptrInt32(int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security)),
		DataType:     ptrInt32(1),
		BeginTime:    ptrStr(begin.Format("2006-01-02")),
		EndTime:      ptrStr(end.Format("2006-01-02")),
	})
	if err != nil {
		log.Fatalf("GetOptionMarketStatistic failed: %v", err)
	}
	fmt.Printf("Items returned: %d\n", len(statRsp.StatisticList))
	for i, it := range statRsp.StatisticList {
		if it == nil {
			continue
		}
		fmt.Printf("  [%d] date=%s call=%d put=%d\n",
			i+1, time.Unix(int64(it.GetTimestamp()), 0).Format("2006-01-02"),
			it.GetCallValue(), it.GetPutValue())
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp)
	display.PrintJSON(rankRsp)
	display.PrintJSON(statRsp)
}

func ptrInt32(v int32) *int32 { return &v }
func ptrStr(v string) *string { return &v }