package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go/pkg/pb/trdflowsummary"
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
	_ = accID

	flows, err := client.GetFlowSummary(context.Background(), mc.Client,
	accID,
	constant.TrdMarket_HK,
	"",
	trdflowsummary.TrdCashFlowDirection(1),
)
	if err != nil {
		log.Fatalf("GetFlowSummary failed: %v", err)
	}
	for _, f := range flows {
		fmt.Printf("FLOW: id=%d date=%s type=%s amount=%.2f remark=%s\n",
			f.CashFlowID, f.ClearingDate, f.CashFlowType, f.CashFlowAmount, f.CashFlowRemark)
	}

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(flows)
}
