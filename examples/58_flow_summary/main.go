package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
)

func main() {
	cli := client.New()
	defer cli.Close()

	addr := os.Getenv("FUTU_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11111"
	}
	if err := cli.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	accounts, err := client.GetAccountList(context.Background(), cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	market := constant.TrdMarket(acc.TrdMarketAuthList[0])

	// direction=1 means cash inflow, direction=2 means outflow
	flows, err := client.GetFlowSummary(cli,
		accID,
		market,
		"",    // clearingDate: empty for today
		1,     // direction: 1=inflow
	)
	if err != nil {
		log.Fatalf("GetFlowSummary failed: %v", err)
	}
	for _, f := range flows {
		fmt.Printf("FLOW: id=%d date=%s type=%s amount=%.2f remark=%s\n",
			f.CashFlowID, f.ClearingDate, f.CashFlowType, f.CashFlowAmount, f.CashFlowRemark)
	}
}
