package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go-demo/examples/pkg/connect"
)

func main() {
	mc := connect.MustConnect(context.Background())
	defer mc.Close()

	accounts, err := client.GetAccountList(context.Background(), mc.Client)
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	acc := mc.Client.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	_ = accID // unused in this example

	positions, err := client.GetPositionList(context.Background(), mc.Client, accID)
	if err != nil {
		log.Fatalf("GetPositionList failed: %v", err)
	}
	if len(positions) == 0 {
		fmt.Println("(no positions)")
	}
	for _, p := range positions {
		fmt.Printf("POS: %s qty=%.0f cost=%.2f cur=%.2f pnl=%.2f (%.2f%%)\n",
			p.Code, p.Quantity, p.CostPrice, p.CurPrice, p.PnL, p.PnLRate)
	}
}