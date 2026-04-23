package main

import (
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

	accounts, err := client.GetAccountList(cli)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("GetAccountList failed: %v", err)
	}

	positions, err := client.GetPositionList(cli, accounts[0].AccID)
	if err != nil {
		log.Fatalf("GetPositionList failed: %v", err)
	}
	if len(positions) == 0 {
		fmt.Println("(no positions)")
	}
	for _, p := range positions {
		fmt.Printf("POS: %s %s qty=%.0f cost=%.2f cur=%.2f pnl=%.2f (%.2f%%)\n",
			p.Code, p.Name, p.Quantity, p.CostPrice, p.CurPrice, p.PnL, p.PnLRate*100)
	}
}
