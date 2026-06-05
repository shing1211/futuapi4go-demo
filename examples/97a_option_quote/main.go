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

	leg := &qotcommon.ComboLeg{
		Security: &qotcommon.Security{
			Market: ptrInt32(int32(constant.Market_US)),
			Code:   ptrStr("NVDA"),
		},
		Side:     ptrInt32(int32(constant.TrdSide_Buy)),
		QtyRatio: ptrFloat64(1.0),
	}

	req := &qot.GetOptionQuoteRequest{
		MultiLegs: []*qotcommon.ComboLeg{leg},
	}

	rsp, err := client.GetOptionQuote(context.Background(), mc.Client, req)
	if err != nil {
		log.Fatalf("GetOptionQuote failed: %v", err)
	}

	fmt.Printf("Got %d option quotes\n", len(rsp.OptionQuoteList))
	for i, q := range rsp.OptionQuoteList {
		if q == nil {
			continue
		}
		fmt.Printf("\n  Quote #%d:\n", i+1)
		fmt.Printf("    Price:     $%.2f\n", q.GetPrice())
		fmt.Printf("    Change:    $%.2f (%.2f%%)\n", q.GetChg(), q.GetChgRate())
		fmt.Printf("    Volume:    %d\n", q.GetVol())
		fmt.Printf("    IV:        %.2f%%\n", q.GetIV())
		fmt.Printf("    Delta:     %.4f\n", q.GetDelta())
		fmt.Printf("    Gamma:     %.4f\n", q.GetGamma())
		fmt.Printf("    Theta:     %.4f\n", q.GetTheta())
		fmt.Printf("    Vega:      %.4f\n", q.GetVega())
		fmt.Printf("    OI:        %d\n", q.GetOpenInterest())
		fmt.Printf("    Expiry:    %s\n", q.GetExpireTime())
		fmt.Printf("    Strike:    $%.2f\n", q.GetStrike())
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(rsp)
}

func ptrInt32(v int32) *int32     { return &v }
func ptrStr(v string) *string     { return &v }
func ptrFloat64(v float64) *float64 { return &v }
