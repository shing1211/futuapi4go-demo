// 116_option_strategy_screener demonstrates the option strategy analysis suite:
//   - GetOptionStrategy        (find available strategies for an underlying)
//   - GetOptionStrategyAnalysis (P&L analysis for a strategy combo)
//   - GetOptionStrategySpread  (available spread values)
//   - GetOptionEarningsScreener (earnings-based option screener)
//   - GetOptionSellerScreener  (seller-side screener)
//   - GetOptionZeroDteScreener (0DTE option screener)
//   - GetOptionZeroDteContract (0DTE contract details)
//   - GetOptionMarketStatistic (market-wide option statistics)
//   - GetOptionUnderlyingHisStatistic (underlying historical stats)
//   - GetOptionUnderlyingHisVolatility (underlying historical volatility)
//   - GetOptionUnderlyingOverview (underlying overview)
//   - GetOptionUnderlyingRank  (underlying ranking)
//   - GetOptionRank            (option ranking)
//   - GetOptionEvent           (option events)
//   - GetOptionEventAlert      (option event alerts)
//   - SetOptionEventAlert      (set option event alerts)
//
// Workflow:
//  1. Find available strategies for NVDA.
//  2. Analyze P&L for a Straddle strategy.
//  3. Check spread values.
//  4. Run screeners for earnings, seller, and 0DTE opportunities.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qot "github.com/shing1211/futuapi4go/pkg/qot"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotgetoptionrank "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionrank"
	qotgetoptionunderlyingrank "github.com/shing1211/futuapi4go/pkg/pb/qotgetoptionunderlyingrank"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client
	ctx := context.Background()

	fmt.Println("=== Option Strategy Screener ===")

	// 1. GetOptionStrategy — find available strategies for NVDA.
	fmt.Println("\n-- 1) GetOptionStrategy (NVDA Straddle) --")
	stratRsp, err := client.GetOptionStrategy(ctx, cli, &qot.GetOptionStrategyRequest{
		Owner: &qotcommon.Security{
			Market: ptrInt32(int32(constant.Market_US)),
			Code:   ptrStr("NVDA"),
		},
		OptionStrategy:  int32(qotcommon.OptionStrategyType_OptionStrategyType_Straddle),
		IndexOptionType: 1,
	})
	if err != nil {
		log.Printf("GetOptionStrategy failed: %v", err)
	} else {
		fmt.Printf("Found %d strategies\n", len(stratRsp.StrategyList))
		for i, s := range stratRsp.StrategyList {
			if s == nil {
				continue
			}
			fmt.Printf("  [%d] code=%s name=%s type=%d legs=%d\n",
				i+1, s.GetCode(), s.GetName(), s.GetOptionStrategy(), len(s.GetMultiLegs()))
		}
		display.PrintJSON(stratRsp)
	}

	// 2. GetOptionStrategyAnalysis — P&L analysis.
	fmt.Println("\n-- 2) GetOptionStrategyAnalysis (skipped — requires valid strategy code) --")

	// 3. GetOptionStrategySpread — spread values.
	fmt.Println("\n-- 3) GetOptionStrategySpread (skipped — requires valid strategy code) --")

	// 4. GetOptionEarningsScreener — earnings-based screener.
	fmt.Println("\n-- 4) GetOptionEarningsScreener --")
	earnRsp, err := client.GetOptionEarningsScreener(ctx, cli, nil)
	if err != nil {
		fmt.Printf("  GetOptionEarningsScreener: %v\n", err)
	} else {
		fmt.Printf("  Earnings screener returned %d items\n", len(earnRsp.ItemList))
		display.PrintJSON(earnRsp)
	}

	// 5. GetOptionSellerScreener — seller-side screener.
	fmt.Println("\n-- 5) GetOptionSellerScreener --")
	sellerRsp, err := client.GetOptionSellerScreener(ctx, cli, nil)
	if err != nil {
		fmt.Printf("  GetOptionSellerScreener: %v\n", err)
	} else {
		fmt.Printf("  Seller screener returned %d items\n", len(sellerRsp.ItemList))
		display.PrintJSON(sellerRsp)
	}

	// 6. GetOptionZeroDteScreener — 0DTE option screener.
	fmt.Println("\n-- 6) GetOptionZeroDteScreener --")
	zeroDteRsp, err := client.GetOptionZeroDteScreener(ctx, cli, nil)
	if err != nil {
		fmt.Printf("  GetOptionZeroDteScreener: %v\n", err)
	} else {
		fmt.Printf("  0DTE screener returned %d items\n", len(zeroDteRsp.ItemList))
		display.PrintJSON(zeroDteRsp)
	}

	// 7. GetOptionRank — option ranking.
	fmt.Println("\n-- 7) GetOptionRank --")
	rankRsp, err := client.GetOptionRank(ctx, cli, &qotgetoptionrank.C2S{})
	if err != nil {
		fmt.Printf("  GetOptionRank: %v\n", err)
	} else {
		fmt.Printf("  Option rank returned %d items\n", len(rankRsp.RankList))
		display.PrintJSON(rankRsp)
	}

	// 8. GetOptionUnderlyingRank — underlying ranking.
	fmt.Println("\n-- 8) GetOptionUnderlyingRank --")
	ulRankRsp, err := client.GetOptionUnderlyingRank(ctx, cli, &qotgetoptionunderlyingrank.C2S{})
	if err != nil {
		fmt.Printf("  GetOptionUnderlyingRank: %v\n", err)
	} else {
		fmt.Printf("  Underlying rank returned %d items\n", len(ulRankRsp.RankList))
		display.PrintJSON(ulRankRsp)
	}

	// 9. GetOptionMarketStatistic — market-wide stats.
	fmt.Println("\n-- 9) GetOptionMarketStatistic --")
	mktStatRsp, err := client.GetOptionMarketStatistic(ctx, cli, nil)
	if err != nil {
		fmt.Printf("  GetOptionMarketStatistic: %v\n", err)
	} else {
		fmt.Printf("  Market statistic received — inspect JSON for details\n")
		display.PrintJSON(mktStatRsp)
	}

	fmt.Println("\n── Done ─────────────────────────────────────────")
}

func ptrInt32(v int32) *int32 { return &v }
func ptrStr(v string) *string { return &v }
