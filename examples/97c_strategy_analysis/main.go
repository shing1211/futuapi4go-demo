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

	legs := []*qotcommon.ComboLeg{
		{
			Security: &qotcommon.Security{
				Market: ptrInt32(int32(constant.Market_US)),
				Code:   ptrStr("NVDA"),
			},
			Side:     ptrInt32(int32(constant.TrdSide_Buy)),
			QtyRatio: ptrFloat64(1.0),
		},
		{
			Security: &qotcommon.Security{
				Market: ptrInt32(int32(constant.Market_US)),
				Code:   ptrStr("NVDA"),
			},
			Side:     ptrInt32(int32(constant.TrdSide_Sell)),
			QtyRatio: ptrFloat64(1.0),
		},
	}

	req := &qot.GetOptionStrategyAnalysisRequest{
		MultiLegs: legs,
	}

	rsp, err := client.GetOptionStrategyAnalysis(context.Background(), mc.Client, req)
	if err != nil {
		log.Fatalf("GetOptionStrategyAnalysis failed: %v", err)
	}

	fmt.Printf("Strategy Analysis: %s (%s)\n", rsp.Code, rsp.Name)
	fmt.Printf("  Type:          %d\n", rsp.OptionStrategy)
	fmt.Printf("  Bid:           $%.2f\n", rsp.Bid1)
	fmt.Printf("  Ask:           $%.2f\n", rsp.Ask1)
	fmt.Printf("  Max Profit:    $%.2f\n", rsp.MaxProfit)
	fmt.Printf("  Max Loss:      $%.2f\n", rsp.MaxLoss)
	fmt.Printf("  Prob of Profit: %.1f%%\n", rsp.ProbOfProfit)
	fmt.Printf("  Breakeven:     ")
	for i, bp := range rsp.BreakevenPoints {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("$%.2f", bp)
	}
	fmt.Println()
	fmt.Printf("  Greeks:\n")
	fmt.Printf("    Delta:  %.4f\n", rsp.Delta)
	fmt.Printf("    Theta:  %.4f\n", rsp.Theta)

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp)
}

func ptrInt32(v int32) *int32     { return &v }
func ptrStr(v string) *string     { return &v }
func ptrFloat64(v float64) *float64 { return &v }
