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

	req := &qot.GetOptionStrategySpreadRequest{
		Owner: &qotcommon.Security{
			Market: ptrInt32(int32(constant.Market_US)),
			Code:   ptrStr("NVDA"),
		},
		OptionStrategy:  int32(qotcommon.OptionStrategyType_OptionStrategyType_Spread),
		IndexOptionType: 1,
	}

	rsp, err := client.GetOptionStrategySpread(context.Background(), mc.Client, req)
	if err != nil {
		log.Fatalf("GetOptionStrategySpread failed: %v", err)
	}

	fmt.Printf("Available spreads for NVDA Spread strategy:\n")
	if len(rsp.SpreadList) == 0 {
		fmt.Println("  (no spread data returned)")
	} else {
		for i, s := range rsp.SpreadList {
			fmt.Printf("  Spread %d: $%.2f\n", i+1, s)
		}
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp)
}

func ptrInt32(v int32) *int32 { return &v }
func ptrStr(v string) *string { return &v }
