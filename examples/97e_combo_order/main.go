package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	qotcommon "github.com/shing1211/futuapi4go/pkg/pb/qotcommon"
	trdcommon "github.com/shing1211/futuapi4go/pkg/pb/trdcommon"
	trd "github.com/shing1211/futuapi4go/pkg/trd"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/display"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := mc.Client.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID

	fmt.Println("=== GetComboMaxTrdQtys ===")
	fmt.Printf("Account ID: %d\n\n", accID)

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

	trdEnv := int32(constant.TrdEnv_Simulate)
	trdMarket := int32(constant.TrdMarket_US)
	header := &trdcommon.TrdHeader{
		AccID:     &accID,
		TrdMarket: &trdMarket,
		TrdEnv:    &trdEnv,
	}

	qtyReq := &trd.GetComboMaxTrdQtysRequest{
		Header:    header,
		ComboLegs: legs,
		Qty:       1,
		OrderType: int32(constant.OrderType_Normal),
	}

	qtyRsp, err := client.GetComboMaxTrdQtys(context.Background(), mc.Client, qtyReq)
	if err != nil {
		fmt.Printf("GetComboMaxTrdQtys: %v\n", err)
		fmt.Println("This is expected in simulate trading mode — requires real options account.")
	} else {
		fmt.Printf("Max combo qty: %+v\n", qtyRsp.MaxTrdQtys)
	}

	fmt.Println()
	fmt.Println("=== PlaceComboOrder ===")
	fmt.Println("SKIPPED: Combo order placement requires real account and options approval.")
	fmt.Println("See PlaceComboOrder in pkg/trd/trade.go for the request structure.")

	placeReq := &trd.PlaceComboOrderRequest{
		Header:    header,
		ComboLegs: legs,
		Qty:       1,
		OrderType: int32(constant.OrderType_Normal),
	}
	fmt.Println()

	fmt.Printf("Would place order with:\n")
	fmt.Printf("  Legs:       %d\n", len(placeReq.ComboLegs))
	fmt.Printf("  Qty:        %.0f\n", placeReq.Qty)
	fmt.Printf("  OrderType:  %d\n", placeReq.OrderType)

	fmt.Println("\n── Request (JSON) ────────────────────────")
	display.PrintJSON(qtyReq)
}

func ptrInt32(v int32) *int32       { return &v }
func ptrStr(v string) *string       { return &v }
func ptrFloat64(v float64) *float64 { return &v }
