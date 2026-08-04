package main

import (
	"context"
	"fmt"

	"github.com/shing1211/futuapi4go/client"
	predcommon "github.com/shing1211/futuapi4go/pkg/pb/common"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	qotcombolist "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractcombolist"
	qotcomborfq "github.com/shing1211/futuapi4go/pkg/pb/qotgeteventcontractcomborfq"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

// Event Contract Combo: combine YES/NO contracts across multiple events into a
// single combo quotation, then request a firm RFQ quote for the leg combo.
func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()
	cli := mc.Client

	fmt.Println("=== Event Contract Combo ===")

	// 1. Discover which events are combine-able into a Combo.
	fmt.Println("\n-- GetEventContractComboList (category=SPORTS) --")
	comboRsp, err := client.GetEventContractComboList(context.Background(), cli,
		&qotcombolist.C2S{Category: ptrStr("SPORTS")})
	if err != nil {
		fmt.Printf("  GetEventContractComboList: %v (expected if no EC combo matches)\n", err)
	} else {
		fmt.Printf("  MVC=%s  combos=%d\n", comboRsp.GetMvc(), len(comboRsp.GetComboEventList()))
		for _, e := range comboRsp.GetComboEventList() {
			fmt.Printf("    event %s  name=%s  contracts=%d\n",
				e.GetEventSecurity().GetCode(), e.GetEventName(), len(e.GetComboContracts()))
		}
	}

	// 2. Build a 2-leg combo from idealised contracts (YES buying team A, YES team B).
	leg1 := &qotcommon.ComboLeg{
		Security: client.NewECSecurity("EC.SPORTSFOOTBALL"),
		Side:     ptrInt32(int32(constant.TrdSide_Buy)),
		QtyRatio: ptrFloat64(1.0),
		PredSide: ptrInt32(int32(predcommon.PredSide_PredSide_Yes)),
	}
	leg2 := &qotcommon.ComboLeg{
		Security: client.NewECSecurity("EC.SPORTSFOOTBALL_2"),
		Side:     ptrInt32(int32(constant.TrdSide_Buy)),
		QtyRatio: ptrFloat64(1.0),
		PredSide: ptrInt32(int32(predcommon.PredSide_PredSide_Yes)),
	}

	// 3. Request a quote (RFQ) for the leg combo. The returned quoteId is
	// required by Trd_PlaceComboOrder when placing the order.
	mvc := ""
	if comboRsp != nil {
		mvc = comboRsp.GetMvc()
	}
	fmt.Println("\n-- GetEventContractComboRfq --")
	rfqRsp, err := client.GetEventContractComboRfq(context.Background(), cli,
		&qotcomborfq.C2S{
			ComboLegList: []*qotcommon.ComboLeg{leg1, leg2},
			Mvc:          ptrStr(mvc),
		})
	if err != nil {
		fmt.Printf("  GetEventContractComboRfq: %v\n", err)
	} else {
		fmt.Printf("  quoteId=%q bid=%.2f ask=%.2f shouldRetry=%v\n",
			rfqRsp.GetQuoteId(), rfqRsp.GetBidPrice(), rfqRsp.GetAskPrice(), rfqRsp.GetShouldRetry())
		fmt.Printf("  legs echoed: %d\n", len(rfqRsp.GetComboLegList()))
	}

	fmt.Println("\n── Full responses (JSON) ────────────────────────")
	if comboRsp != nil {
		display.PrintJSON(comboRsp)
	}
	if rfqRsp != nil {
		display.PrintJSON(rfqRsp)
	}
}

func ptrStr(v string) *string   { return &v }
func ptrInt32(v int32) *int32   { return &v }
func ptrFloat64(v float64) *float64 { return &v }