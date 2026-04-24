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

	acc := cli.FindAccount(accounts)
	if acc == nil {
		log.Fatal("no account found")
	}
	accID := acc.AccID
	market := acc.TrdMarketAuthList[0]

	orders, err := client.GetHistoryOrderList(cli,
		accID,
		market,
		"2026-01-01", "2026-04-24",
	)
	if err != nil {
		log.Fatalf("GetHistoryOrderList failed: %v", err)
	}
	if len(orders) == 0 {
		fmt.Println("(no historical orders)")
	}
	for _, o := range orders {
		fmt.Printf("HIST ORDER: id=%d %s %s qty=%.0f price=%.2f status=%d\n",
			o.OrderID, o.Code, o.Name, o.Qty, o.Price, o.OrderStatus)
	}
}
