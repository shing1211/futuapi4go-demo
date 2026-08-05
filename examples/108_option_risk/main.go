// 108_option_risk demonstrates option risk and underlying-level analytics:
//   - GetOptionExerciseProbability        (per-strike exercise probability)
//   - GetOptionUnderlyingHisStatistic     (historical open-interest / volume stat)
//   - GetOptionUnderlyingHisVolatility    (historical IV series for an underlying)
//   - GetOptionUnderlyingOverview         (per-underlying overview)
//   - GetOptionUnderlyingRank             (ranking of underlyings by IV/OI/etc.)
//
// All four target NVDA in the US option market.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotoptioncommon "github.com/shing1211/futuapi4go/pkg/pb/qotoptioncommon"
	qotgetoptionunderlyinghisstatistic "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionunderlyinghisstatistic"
	qotgetoptionunderlyinghisvolatility "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionunderlyinghisvolatility"
	qotgetoptionunderlyingoverview "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionunderlyingoverview"
	qotgetoptionunderlyingrank "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionunderlyingrank"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	market := int32(constant.Market_US)
	ownerSec := &qotcommon.Security{Market: &market, Code: ptrStr("NVDA")}

	fmt.Println("=== Option Risk & Underlying Analytics (NVDA, US) ===")

	fmt.Println("--- 1) GetOptionExerciseProbability ---")
	rsp1, err := client.GetOptionExerciseProbability(ctx, mc.Client, constant.Market(market), "NVDA")
	if err != nil {
		log.Fatalf("GetOptionExerciseProbability failed: %v", err)
	}
	fmt.Printf("Strike-probability items: %d\n", len(rsp1.ItemList))
	for i, it := range rsp1.ItemList {
		if it == nil {
			continue
		}
		fmt.Printf("  [%d] ts=%s price=%.2f prob=%.2f%%\n",
			i+1, it.GetTimestampStr(), it.GetSecurityPrice(), it.GetStrikeProbability()*100)
		if i >= 4 {
			fmt.Println("  ... (truncated)")
			break
		}
	}
	fmt.Println()

	fmt.Println("--- 2) GetOptionUnderlyingHisStatistic (last 5 trading days) ---")
	end := time.Now()
	begin := end.AddDate(0, 0, -5)
	rsp2, err := client.GetOptionUnderlyingHisStatistic(ctx, mc.Client,
		&qotgetoptionunderlyinghisstatistic.C2S{
			Owner:          ownerSec,
			BeginTime:      ptrStr(begin.Format("2006-01-02")),
			EndTime:        ptrStr(end.Format("2006-01-02")),
			IndexOptionType: ptrInt32(0),
		})
	if err != nil {
		log.Fatalf("GetOptionUnderlyingHisStatistic failed: %v", err)
	}
	fmt.Printf("Underlying %s  |  Statistic items: %d\n", rsp2.GetName(), len(rsp2.StatisticList))
	for i, it := range rsp2.StatisticList {
		if it == nil {
			continue
		}
		fmt.Printf("  [%d] date=%s callVol=%d putVol=%d callOI=%d putOI=%d price=%.2f\n",
			i+1, time.Unix(int64(it.GetTimestamp()), 0).Format("2006-01-02"),
			it.GetCallVolume(), it.GetPutVolume(),
			it.GetCallOpenInterest(), it.GetPutOpenInterest(),
			it.GetUnderlyingPrice())
	}
	fmt.Println()

	fmt.Println("--- 3) GetOptionUnderlyingHisVolatility (last 5 trading days) ---")
	rsp3, err := client.GetOptionUnderlyingHisVolatility(ctx, mc.Client,
		&qotgetoptionunderlyinghisvolatility.C2S{
			Owner:          ownerSec,
			BeginTime:      ptrStr(begin.Format("2006-01-02")),
			EndTime:        ptrStr(end.Format("2006-01-02")),
			IndexOptionType: ptrInt32(0),
		})
	if err != nil {
		log.Fatalf("GetOptionUnderlyingHisVolatility failed: %v", err)
	}
	fmt.Printf("Underlying %s  |  Volatility items: %d\n", rsp3.GetName(), len(rsp3.VolatilityList))
	for i, it := range rsp3.VolatilityList {
		if it == nil {
			continue
		}
		fmt.Printf("  [%d] date=%s IV=%.2f%%  HV=%.2f%%  price=%.2f\n",
			i+1, time.Unix(int64(it.GetTimestamp()), 0).Format("2006-01-02"),
			it.GetIv(), it.GetHv(), it.GetUnderlyingPrice())
	}
	fmt.Println()

	fmt.Println("--- 4) GetOptionUnderlyingOverview ---")
	rsp4, err := client.GetOptionUnderlyingOverview(ctx, mc.Client,
		&qotgetoptionunderlyingoverview.C2S{
			OwnerList:       []*qotcommon.Security{ownerSec},
			IndexOptionType: ptrInt32(0),
		})
	if err != nil {
		log.Fatalf("GetOptionUnderlyingOverview failed: %v", err)
	}
	fmt.Printf("Items: %d\n", len(rsp4.UnderlyingDataList))
	for i, it := range rsp4.UnderlyingDataList {
		if it == nil {
			continue
		}
		fmt.Printf("  [%d] %s  IV=%.2f%%  IVRank=%.1f  IVPctl=%.1f  callVol=%d putVol=%d callOI=%d putOI=%d\n",
			i+1, it.GetCode(), it.GetIv(), it.GetIvRank(), it.GetIvPercentile(),
			it.GetCallVolume(), it.GetPutVolume(),
			it.GetCallOpenInterest(), it.GetPutOpenInterest())
	}
	fmt.Println()

	fmt.Println("--- 5) GetOptionUnderlyingRank (top 5 by IV) ---")
	rsp5, err := client.GetOptionUnderlyingRank(ctx, mc.Client,
		&qotgetoptionunderlyingrank.C2S{
			OptionMarket: ptrInt32(int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security)),
			SortType:     ptrInt32(int32(qotoptioncommon.UnderlyingRankSortType_UnderlyingRankSortType_IV)),
			Count:        ptrInt32(5),
		})
	if err != nil {
		log.Fatalf("GetOptionUnderlyingRank failed: %v", err)
	}
	fmt.Printf("AllCount: %d  |  Items: %d\n", rsp5.GetAllCount(), len(rsp5.RankList))
	for i, it := range rsp5.RankList {
		if it == nil {
			continue
		}
		fmt.Printf("  #%d %-8s IV=%.2f%%  IVRank=%.1f  HV=%.2f%%  vol=%d  mktcap=%.0f\n",
			i+1, it.GetName(), it.GetIv(), it.GetIvRank(), it.GetHv(),
			it.GetTotalVolume(), it.GetMarketCap())
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp1)
	display.PrintJSON(rsp2)
	display.PrintJSON(rsp3)
	display.PrintJSON(rsp4)
	display.PrintJSON(rsp5)
}

func ptrInt32(v int32) *int32 { return &v }
func ptrStr(v string) *string { return &v }