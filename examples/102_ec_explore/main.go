package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	qotcategory "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractcategory"
	qotcontract "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontract"
	qoteventlist "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontracteventlist"
	qotfilter "github.com/shing1211/futuapi4go/pkg/pb/qotfiltercompetition"
	qotmilestone "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractmilestonelist"
	qotseries "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractserieslist"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

// Explore the Prediction Market / Event Contract (EC) hierarchy:
//   category -> series -> event -> contracts -> milestones.
//
// EC instruments (Moomoo US Prediction) trade YES / NO binary outcomes on future
// events (sports, politics, economics, ...). Market = 101 (QotMarket_EventContract).
func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client

	fmt.Println("=== Event Contract Discovery (Moomoo US Prediction) ===")

	// 1. Top-level competition filters (which sports / tags are tradable).
	fmt.Println("\n-- FilterCompetition (all categories) --")
	fcRsp, err := client.FilterCompetition(context.Background(), cli, &qotfilter.C2S{})
	if err != nil {
		fmt.Printf("  FilterCompetition: %v (expected if EC market is unavailable)\n", err)
	} else {
		for _, f := range fcRsp.TagFilterList {
			fmt.Printf("  category=%-10s tag=%-14s comps=%d scopes=%d\n",
				f.GetCategory(), f.GetTag(), len(f.GetCompetitionList()), len(f.GetScopeList()))
		}
	}

	// 2. Top-level EC categories.
	fmt.Println("\n-- GetEventContractCategory --")
	catRsp, err := client.GetEventContractCategory(context.Background(), cli, &qotcategory.C2S{})
	if err != nil {
		fmt.Printf("  GetEventContractCategory: %v\n", err)
	} else {
		for _, c := range catRsp.CategoryList {
			fmt.Printf("  %-12s %-20s tags=%v\n", c.GetCategory(), c.GetCategoryName(), c.GetTags())
		}
	}

	// 3. Series under a specific category (e.g. "SPORTS").
	fmt.Println("\n-- GetEventContractSeriesList (category=SPORTS) --")
	seriesRsp, err := client.GetEventContractSeriesList(context.Background(), cli,
		&qotseries.C2S{Category: ptrStr("SPORTS")})
	if err != nil {
		fmt.Printf("  GetEventContractSeriesList: %v\n", err)
	} else {
		for _, s := range seriesRsp.SeriesList {
			fmt.Printf("  series code=%-28s name=%s tags=%v\n",
				s.GetSeriesSecurity().GetCode(), s.GetSeriesName(), s.GetTags())
		}
	}

	// 4. Events under a specific Series.
	fmt.Println("\n-- GetEventContractEventList (series=EC.SPORTSFOOTBALL) --")
	evRsp, err := client.GetEventContractEventList(context.Background(), cli,
		&qoteventlist.C2S{
			Series: client.NewECSecurity("EC.SPORTSFOOTBALL"),
			Count:  ptrInt32(10),
		})
	if err != nil {
		fmt.Printf("  GetEventContractEventList: %v\n", err)
	} else {
		for _, e := range evRsp.EventList {
			fmt.Printf("  event code=%-28s name=%s status=%d\n",
				e.GetEventSecurity().GetCode(), e.GetEventName(), e.GetStatus())
		}
	}

	// 5. Contracts under an event.
	fmt.Println("\n-- GetEventContract (event=EC.SPORTSFOOTBALL) --")
	ctRsp, err := client.GetEventContract(context.Background(), cli,
		&qotcontract.C2S{Event: client.NewECSecurity("EC.SPORTSFOOTBALL")})
	if err != nil {
		fmt.Printf("  GetEventContract: %v\n", err)
	} else {
		for _, c := range ctRsp.ContractList {
			fmt.Printf("  contract code=%s type=%d status=%d result=%s\n",
				c.GetContractSecurity().GetCode(), c.GetContractType(), c.GetStatus(), c.GetResult())
		}
	}

	// 6. Milestones that affect EC outcomes (optional discovery).
	fmt.Println("\n-- GetEventContractMilestoneList (category=SPORTS) --")
	msRsp, err := client.GetEventContractMilestoneList(context.Background(), cli,
		&qotmilestone.C2S{Category: ptrStr("SPORTS")})
	if err != nil {
		fmt.Printf("  GetEventContractMilestoneList: %v\n", err)
	} else {
		for _, m := range msRsp.MilestoneList {
			fmt.Printf("  milestone code=%s title=%s start=%s\n",
				m.GetMilestoneSecurity().GetCode(), m.GetTitle(), m.GetStartDate())
		}
	}

	fmt.Println("\n── Full responses (JSON) ────────────────────────")
	if catRsp != nil {
		display.PrintJSON(catRsp)
	}
}

func ptrStr(v string) *string { return &v }
func ptrInt32(v int32) *int32 { return &v }
