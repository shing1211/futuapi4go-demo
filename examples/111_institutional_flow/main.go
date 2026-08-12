// 111_institutional_flow demonstrates the institutional-investor data suite:
//   - GetInstitutionList          (list institutions in a market, by holding value)
//   - GetInstitutionProfile       (institution overview)
//   - GetInstitutionHoldingList   (stocks an institution holds)
//   - GetInstitutionHoldingChange (recent holding changes)
//   - GetInstitutionDistribution  (holding distribution across holdings)
//
// Workflow:
//  1. List the top institutions in HK by position value.
//  2. Take the first institution ID and chain profile / holdings /
//     holding-change / distribution calls off it.
//
// The institution APIs work best on the HK market; use Market_HK (1).
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetinstitutiondistribution"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetinstitutionholdingchange"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetinstitutionholdinglist"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetinstitutionlist"
	"github.com/shing1211/futuapi4go/pkg/pb/qotgetinstitutionprofile"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	ctx := context.Background()
	market := int32(1) // Market_HK

	fmt.Println("=== Institutional Flow ===")

	fmt.Println("--- 1) GetInstitutionList (top HK institutions by position value) ---")
	listRsp, err := client.GetInstitutionList(ctx, mc.Client, &qotgetinstitutionlist.C2S{
		Market: ptrInt32(market),
		Count:  ptrInt32(5),
	})
	if err != nil {
		log.Fatalf("GetInstitutionList failed: %v", err)
	}
	fmt.Printf("AllCount=%d  currency=%s\n",
		listRsp.GetAllCount(), derefStr(listRsp.Currency))

	instID := int32(0)
	instName := ""
	for i, item := range listRsp.DataList {
		if item == nil {
			continue
		}
		fmt.Printf("  [%d] id=%d %s  value=%.0f (%+.0f)  holdings=%d\n",
			i+1,
			item.GetInstitutionId(),
			item.GetInstitutionName(),
			item.GetPositionValue(),
			item.GetPositionValueChange(),
			item.GetPositionCount(),
		)
		if instID == 0 {
			instID = item.GetInstitutionId()
			instName = item.GetInstitutionName()
		}
	}
	if instID == 0 {
		log.Fatalf("GetInstitutionList returned no institution; data may be sparse outside HK")
	}
	fmt.Println()

	fmt.Println("--- 2) GetInstitutionProfile ---")
	profRsp, err := client.GetInstitutionProfile(ctx, mc.Client, &qotgetinstitutionprofile.C2S{
		Market:        ptrInt32(market),
		InstitutionId: ptrInt32(instID),
	})
	if err != nil {
		log.Printf("GetInstitutionProfile failed: %v", err)
	} else {
		fmt.Printf("  Profile for %s (%d) — inspect JSON for details\n", instName, instID)
	}
	fmt.Println()

	fmt.Println("--- 3) GetInstitutionHoldingList (first page) ---")
	holdRsp, err := client.GetInstitutionHoldingList(ctx, mc.Client, &qotgetinstitutionholdinglist.C2S{
		Market:        ptrInt32(market),
		InstitutionId: ptrInt32(instID),
		Count:         ptrInt32(5),
	})
	if err != nil {
		log.Printf("GetInstitutionHoldingList failed: %v", err)
	} else {
		fmt.Printf("  Holdings: %d (nextPage=%q) — inspect JSON for the stock list\n",
			len(holdRsp.DataList), derefStr(holdRsp.NextPage))
	}
	fmt.Println()

	fmt.Println("--- 4) GetInstitutionHoldingChange (recent changes) ---")
	chgRsp, err := client.GetInstitutionHoldingChange(ctx, mc.Client, &qotgetinstitutionholdingchange.C2S{
		Market:        ptrInt32(market),
		InstitutionId: ptrInt32(instID),
		Count:         ptrInt32(5),
	})
	if err != nil {
		log.Printf("GetInstitutionHoldingChange failed: %v", err)
	} else {
		fmt.Printf("  Change items: %d — inspect JSON for add/reduce/steady events\n", len(chgRsp.DataList))
	}
	fmt.Println()

	fmt.Println("--- 5) GetInstitutionDistribution ---")
	distRsp, err := client.GetInstitutionDistribution(ctx, mc.Client, &qotgetinstitutiondistribution.C2S{
		Market:        ptrInt32(market),
		InstitutionId: ptrInt32(instID),
	})
	if err != nil {
		log.Printf("GetInstitutionDistribution failed: %v", err)
	} else {
		fmt.Printf("  Distribution received — inspect JSON for holding mix\n")
	}
	fmt.Println()

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(listRsp)
	for _, r := range []any{profRsp, holdRsp, chgRsp, distRsp} {
		if r != nil {
			display.PrintJSON(r)
		}
	}
}

func ptrInt32(v int32) *int32 { return &v }

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
