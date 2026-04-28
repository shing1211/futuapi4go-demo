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

	fills, err := client.GetHistoryOrderFillList(context.Background(), cli,
		accID,
		market,
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
