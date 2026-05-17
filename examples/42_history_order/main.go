package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shing1211/futuapi4go/client"
	"github.com/shing1211/futuapi4go/pkg/constant"
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

	orders, err := client.GetHistoryOrderList(context.Background(), mc.Client,
	accID, constant.TrdMarket_HK,
	"2025-01-01", "2026-04-28",
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

	fmt.Println("\n── Result (JSON) ────────────────────────")
	display.PrintJSON(orders)
}
