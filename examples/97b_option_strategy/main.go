package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qot "github.com/shing1211/futuapi4go/pkg/qot"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	req := &qot.GetOptionStrategyRequest{
		Owner: &qotcommon.Security{
			Market: ptrInt32(int32(constant.Market_US)),
			Code:   ptrStr("NVDA"),
		},
		OptionStrategy:  int32(qotcommon.OptionStrategyType_OptionStrategyType_Straddle),
		IndexOptionType: 1,
	}

	rsp, err := client.GetOptionStrategy(context.Background(), mc.Client, req)
	if err != nil {
		log.Fatalf("GetOptionStrategy failed: %v", err)
	}

	fmt.Printf("Found %d option strategies\n\n", len(rsp.StrategyList))
	for i, item := range rsp.StrategyList {
		if item == nil {
			continue
		}
		fmt.Printf("Strategy #%d:\n", i+1)
		fmt.Printf("  Code:       %s\n", item.GetCode())
		fmt.Printf("  Name:       %s\n", item.GetName())
		fmt.Printf("  Type:       %d\n", item.GetOptionStrategy())
		fmt.Printf("  Legs:       %d\n", len(item.GetMultiLegs()))
		for j, leg := range item.GetMultiLegs() {
			if leg == nil {
				continue
			}
			sec := leg.GetSecurity()
			code := ""
			if sec != nil {
				code = sec.GetCode()
			}
			side := "Buy"
			if leg.GetSide() == 2 {
				side = "Sell"
			}
			fmt.Printf("    Leg %d: %s %s (ratio=%.0f)\n", j+1, side, code, leg.GetQtyRatio())
		}
		fmt.Println()
	}

	fmt.Println("── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp)
}

func ptrInt32(v int32) *int32   { return &v }
func ptrStr(v string) *string   { return &v }
