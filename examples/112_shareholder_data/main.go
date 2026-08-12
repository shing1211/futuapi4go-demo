// 112_shareholder_data demonstrates the shareholder-composition suite:
//   - GetShareholdersOverview       (top holders by type: main, institutional, etc.)
//   - GetShareholdersHolderDetail   (detailed holding rows with change deltas)
//   - GetShareholdersHoldingChanges (reported holder change events)
//   - GetShareholdersInstitutional  (institutional holders only)
//
// All four are convenience wrappers taking (market, code, ...) directly.
// NVDA (US) has the richest shareholder data.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	market := constant.Market(constant.Market_US)
	code := "NVDA"

	fmt.Println("=== Shareholder Data (US.NVDA) ===")

	fmt.Println("--- 1) GetShareholdersOverview (period 1 = latest) ---")
	overview, err := client.GetShareholdersOverview(ctx, mc.Client, market, code, 1)
	if err != nil {
		log.Printf("GetShareholdersOverview failed: %v", err)
	} else {
		fmt.Printf("  Main holders: %d  |  Holder-type rows: %d  |  Holding periods: %d\n",
			len(overview.MainHolderInfoList),
			len(overview.HolderTypeInfoList),
			len(overview.HoldingPeriodList),
		)
	}
	fmt.Println()

	fmt.Println("--- 2) GetShareholdersHolderDetail (top holders by ratio) ---")
	detail, err := client.GetShareholdersHolderDetail(ctx, mc.Client, market, code,
		1 /*reqType=holder*/, "" /*nextKey*/, 10 /*num*/, 0 /*sortColumn*/, 0 /*sortType*/, 1 /*periodId*/, 0 /*holderId*/)
	if err != nil {
		log.Printf("GetShareholdersHolderDetail failed: %v", err)
	} else {
		fmt.Printf("  Detail rows: %d (nextKey=%q)\n", len(detail.ItemList), detail.NextKey)
	}
	fmt.Println()

	fmt.Println("--- 3) GetShareholdersHoldingChanges (recent events) ---")
	changes, err := client.GetShareholdersHoldingChanges(ctx, mc.Client, market, code,
		"" /*nextKey*/, 10 /*num*/, 0 /*sortType*/, 0 /*sortColumn*/, 0 /*filterType*/)
	if err != nil {
		log.Printf("GetShareholdersHoldingChanges failed: %v", err)
	} else {
		fmt.Printf("  Change events: %d (nextKey=%q)\n", len(changes.ItemList), changes.NextKey)
	}
	fmt.Println()

	fmt.Println("--- 4) GetShareholdersInstitutional ---")
	institutional, err := client.GetShareholdersInstitutional(ctx, mc.Client, market, code,
		"" /*nextKey*/, 10 /*num*/)
	if err != nil {
		log.Printf("GetShareholdersInstitutional failed: %v", err)
	} else {
		fmt.Printf("  Institutional holders: %d (nextKey=%q)\n", len(institutional.ItemList), institutional.NextKey)
	}
	fmt.Println()

	fmt.Println("\n── Result (JSON) ────────────────────────")
	for _, r := range []any{overview, detail, changes, institutional} {
		if r != nil {
			display.PrintJSON(r)
		}
	}
}
