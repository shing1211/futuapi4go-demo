package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/futuapi4go/client"
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
	if err != nil {
		log.Fatalf("GetAccountList failed: %v", err)
	}
	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	_ = accID // unused in this example

	positions, err := client.GetPositionList(context.Background(), cli, accID)
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