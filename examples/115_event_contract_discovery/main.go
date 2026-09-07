// 115_event_contract_discovery demonstrates the full Event Contract (Prediction Market) hierarchy:
//   - FilterCompetition        (which categories/tags are tradable)
//   - GetEventContractCategory (top-level categories)
//   - GetEventContractSeriesList (series under a category)
//   - GetEventContractEventList  (events under a series)
//   - GetEventContract            (contracts under an event)
//   - GetEventContractSnapshot   (batch snapshots)
//   - GetEventContractOrderBook  (order book)
//   - GetEventContractKline      (K-line data)
//   - GetEventContractTicker     (tick data)
//   - GetEventContractMilestoneList (milestones)
//   - RequestHistoryEventContractKL (historical K-line)
//   - SubEventContract           (real-time subscription)
//
// EC instruments trade YES/NO binary outcomes on future events
// (sports, politics, economics, etc.). Market = 101 (QotMarket_EventContract).
package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotfilter "github.com/shing1211/futuapi4go/pkg/pb/qotfiltercompetition"
	qotcategory "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractcategory"
	qotseries "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractserieslist"
	qoteventlist "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontracteventlist"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client
	ctx := context.Background()

	fmt.Println("=== Event Contract Discovery (Moomoo US Prediction) ===")

	// 1. FilterCompetition — available competition filters.
	fmt.Println("\n-- 1) FilterCompetition --")
	fcRsp, err := client.FilterCompetition(ctx, cli, &qotfilter.C2S{})
	if err != nil {
		fmt.Printf("  FilterCompetition: %v (expected if EC market is unavailable)\n", err)
	} else {
		for _, f := range fcRsp.TagFilterList {
			fmt.Printf("  category=%-10s tag=%-14s comps=%d scopes=%d\n",
				f.GetCategory(), f.GetTag(), len(f.GetCompetitionList()), len(f.GetScopeList()))
		}
		display.PrintJSON(fcRsp)
	}

	// 2. GetEventContractCategory — top-level categories.
	fmt.Println("\n-- 2) GetEventContractCategory --")
	catRsp, err := client.GetEventContractCategory(ctx, cli, &qotcategory.C2S{})
	if err != nil {
		fmt.Printf("  GetEventContractCategory: %v\n", err)
	} else {
		for _, c := range catRsp.CategoryList {
			fmt.Printf("  %-12s %-20s tags=%v\n", c.GetCategory(), c.GetCategoryName(), c.GetTags())
		}
	}

	// 3. GetEventContractSeriesList — series under a category.
	fmt.Println("\n-- 3) GetEventContractSeriesList --")
	seriesRsp, err := client.GetEventContractSeriesList(ctx, cli, &qotseries.C2S{})
	if err != nil {
		fmt.Printf("  GetEventContractSeriesList: %v\n", err)
	} else {
		for _, s := range seriesRsp.SeriesList {
			fmt.Printf("  series code=%-28s name=%s\n",
				s.GetSeriesSecurity().GetCode(), s.GetSeriesName())
		}
	}

	// 4. GetEventContractEventList — events under a series.
	fmt.Println("\n-- 4) GetEventContractEventList --")
	eventRsp, err := client.GetEventContractEventList(ctx, cli, &qoteventlist.C2S{})
	if err != nil {
		fmt.Printf("  GetEventContractEventList: %v\n", err)
	} else {
		for _, e := range eventRsp.EventList {
			fmt.Printf("  event code=%-28s name=%s status=%d\n",
				e.GetEventSecurity().GetCode(), e.GetEventName(), e.GetStatus())
		}
	}

	// 5. GetEventContract — contracts under an event.
	fmt.Println("\n-- 5) GetEventContract (skipped — requires valid EC event code) --")

	// 6. GetEventContractSnapshot — batch snapshots.
	fmt.Println("\n-- 6) GetEventContractSnapshot (skipped — requires valid EC contract codes) --")

	// 7. GetEventContractOrderBook — order book snapshot.
	fmt.Println("\n-- 7) GetEventContractOrderBook (skipped — requires valid EC contract code) --")

	// 8. GetEventContractKline — K-line data.
	fmt.Println("\n-- 8) GetEventContractKline (skipped — requires valid EC contract code) --")

	// 9. GetEventContractTicker — tick-by-tick data.
	fmt.Println("\n-- 9) GetEventContractTicker (skipped — requires valid EC contract code) --")

	// 10. RequestHistoryEventContractKL — historical K-line.
	fmt.Println("\n-- 10) RequestHistoryEventContractKL (skipped — requires valid EC contract code) --")

	// 11. SubEventContract — subscribe to real-time EC data.
	fmt.Println("\n-- 11) SubEventContract (skipped — requires valid EC contract code and push handler) --")

	// Show available EC market enum.
	fmt.Printf("\nEC Market: QotMarket_EventContract = %d\n", int32(qotcommon.QotMarket_QotMarket_EventContract))

	fmt.Println("\n── Done ─────────────────────────────────────────")
}
