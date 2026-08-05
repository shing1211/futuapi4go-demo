// 107_option_screener demonstrates the four option screening engines:
//   - GetOptionEarningsScreener   (sorted by IV/HV around upcoming earnings)
//   - GetOptionZeroDteScreener    (0-DTE option candidates)
//   - GetOptionSellerScreener     (covered-call / cash-secured-put candidates)
//   - OptionScreen                (v10.6 multi-criteria option screen)
//
// All four target the US option market (OptionMarket_US_Security).
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	qotoptioncommon "github.com/shing1211/futuapi4go/pkg/pb/qotoptioncommon"
	qotgetoptionearningsscreener "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionearningsscreener"
	qotgetoptionsellerscreener "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionsellerscreener"
	qotgetoptionzerodtescreener "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionzerodtescreener"
	qot "github.com/shing1211/futuapi4go/pkg/qot"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	optMarketUS := ptrInt32(int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security))

	fmt.Println("=== Option Screener Suite (US) ===")

	fmt.Println("--- 1) GetOptionEarningsScreener (top 5 by IV) ---")
	rsp1, err := client.GetOptionEarningsScreener(ctx, mc.Client, &qotgetoptionearningsscreener.C2S{
		OptionMarket: optMarketUS,
		SortType:     ptrInt32(int32(qotoptioncommon.EarningsSortType_EarningsSortType_IV)),
		Count:        ptrInt32(5),
	})
	if err != nil {
		log.Fatalf("GetOptionEarningsScreener failed: %v", err)
	}
	fmt.Printf("Items: %d  |  AllCount: %d\n", len(rsp1.ItemList), rsp1.GetAllCount())
	for i, it := range rsp1.ItemList {
		if it == nil {
			continue
		}
		fmt.Printf("  #%d %-8s IV=%.2f%%  IVr=%.1f  HV=%.2f%%  earnTS=%.0f\n",
			i+1, it.GetName(), it.GetIv(), it.GetIvRank(), it.GetHv(), it.GetEarningsTimestamp())
	}
	fmt.Println()

	fmt.Println("--- 2) GetOptionZeroDteScreener (top 5 by Volume) ---")
	rsp2, err := client.GetOptionZeroDteScreener(ctx, mc.Client, &qotgetoptionzerodtescreener.C2S{
		OptionMarket: optMarketUS,
		SortType:     ptrInt32(int32(qotoptioncommon.ZeroDteSortType_ZeroDteSortType_Volume)),
		Count:        ptrInt32(5),
	})
	if err != nil {
		log.Fatalf("GetOptionZeroDteScreener failed: %v", err)
	}
	fmt.Printf("Items: %d\n", len(rsp2.ItemList))
	for i, it := range rsp2.ItemList {
		if it == nil {
			continue
		}
		fmt.Printf("  #%d %-8s IV=%.2f%%  vol=%d  OI=%d\n",
			i+1, it.GetName(), it.GetIv(), it.GetVolume(), it.GetOpenInterest())
	}
	fmt.Println()

	fmt.Println("--- 3) GetOptionSellerScreener (Covered Calls, top 5 by annualized return) ---")
	rsp3, err := client.GetOptionSellerScreener(ctx, mc.Client, &qotgetoptionsellerscreener.C2S{
		OptionMarket: optMarketUS,
		SellerType:   ptrInt32(int32(qotoptioncommon.SellerType_SellerType_CoveredCall)),
		SortType:     ptrInt32(int32(qotoptioncommon.SellerSortType_SellerSortType_AnnualizedReturn)),
	})
	if err != nil {
		log.Fatalf("GetOptionSellerScreener failed: %v", err)
	}
	fmt.Printf("Items: %d\n", len(rsp3.ItemList))
	for i, it := range rsp3.ItemList {
		if it == nil {
			continue
		}
		sec := it.GetOption()
		code := ""
		if sec != nil {
			code = sec.GetCode()
		}
		fmt.Printf("  #%d %-22s strike=%.2f prem=%.2f annRet=%.2f%% days=%d\n",
			i+1, code, it.GetStrikePrice(), it.GetPremium(),
			it.GetAnnualizedReturn(), it.GetLeftDays())
	}
	fmt.Println()

	fmt.Println("--- 4) OptionScreen (v10.6 — page 1, 10 results) ---")
	screenRsp, err := client.OptionScreen(ctx, mc.Client, &qot.OptionScreenRequest{
		MarketCategoryList: []int32{int32(qotoptioncommon.OptionMarket_OptionMarket_US_Security)},
		PageFrom:           0,
		PageCount:          10,
	})
	if err != nil {
		log.Fatalf("OptionScreen failed: %v", err)
	}
	fmt.Printf("AllCount: %d  |  LastPage: %v  |  Items: %d\n",
		screenRsp.AllCount, screenRsp.LastPage, len(screenRsp.DataList))
	for i, it := range screenRsp.DataList {
		if it == nil {
			continue
		}
		sec := it.Security
		code := ""
		if sec != nil {
			code = sec.GetCode()
		}
		fmt.Printf("  #%d %-22s IV=%.2f%%  OI=%d  daysLeft=%d\n",
			i+1, code, it.GetImpliedVolatility(), it.GetOpenInterest(), it.GetLeftDay())
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp1)
	display.PrintJSON(rsp2)
	display.PrintJSON(rsp3)
	display.PrintJSON(screenRsp)
}

func ptrInt32(v int32) *int32 { return &v }