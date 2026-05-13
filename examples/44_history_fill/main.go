package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
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

	fills, err := client.GetHistoryOrderFillList(context.Background(), mc.Client,
	accID, constant.TrdMarket_HK,
)
	if err != nil {
		log.Fatalf("GetHistoryOrderFillList failed: %v", err)
	}
	if len(fills) == 0 {
		fmt.Println("(no historical fills)")
	}
	for _, f := range fills {
		fmt.Printf("HIST FILL: id=%d %s price=%.2f qty=%.0f time=%s\n",
			f.FillID, f.Code, f.Price, f.Qty, f.CreateTime)
	}
}
